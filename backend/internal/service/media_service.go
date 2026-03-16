// Package service 业务逻辑层
//
// 职责：实现媒体文件上传业务逻辑，支持本地/S3 双轨存储策略
// 对外接口：MediaService
package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"docplatform/internal/config"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"
	"docplatform/pkg/logger"
	"docplatform/pkg/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// 本地存储上限：100MB
const localStorageLimitBytes int64 = 100 * 1024 * 1024

// 最大单文件大小 50MB
const maxFileSize = 50 * 1024 * 1024

// detectFileMIME 通过文件内容的 magic bytes 检测真实 MIME 类型。
// 优先使用 http.DetectContentType（基于文件头512字节），
// 如果检测结果为 application/octet-stream 且客户端提供了有效 Content-Type 则回退使用客户端声明。
func detectFileMIME(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 读取文件头 512 字节用于 MIME 检测
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && n == 0 {
		return "", err
	}

	detected := http.DetectContentType(buf[:n])
	// http.DetectContentType 对未知类型返回 application/octet-stream，
	// 此时回退使用客户端声明的 Content-Type（如有）
	if detected == "application/octet-stream" {
		if ct := file.Header.Get("Content-Type"); ct != "" {
			return ct, nil
		}
	}
	return detected, nil
}

// MediaService 媒体文件业务
type MediaService struct {
	mediaRepo    *repository.MediaRepo
	settingsRepo *repository.TenantSettingsRepo
	pageRepo     *repository.PageRepo
	versionRepo  *repository.VersionRepo
	themeRepo    *repository.ThemeRepo
}

// NewMediaService 创建 MediaService
func NewMediaService() *MediaService {
	return &MediaService{
		mediaRepo:    repository.NewMediaRepo(),
		settingsRepo: repository.NewTenantSettingsRepo(),
		pageRepo:     repository.NewPageRepo(),
		versionRepo:  repository.NewVersionRepo(),
		themeRepo:    repository.NewThemeRepo(),
	}
}

// Upload 上传文件（自动选择存储策略）
// 策略：1) 优先按租户 DefaultTarget 决定 local/cloud
//
//  2. 若 DefaultTarget=local 但本地已满(>=100MB) → 自动切换 cloud（需已启用）
//  3. 若无法上传则返回对应业务错误
func (s *MediaService) Upload(ctx context.Context, tenantID string, userID primitive.ObjectID, file *multipart.FileHeader) (*entity.Media, *errcode.AppError) {
	if file.Size > maxFileSize {
		return nil, errcode.ErrUploadFileTooLarge
	}

	// 通过读取文件头部 magic bytes 检测真实 MIME 类型，防止 Content-Type 伪造
	mimeType, detectErr := detectFileMIME(file)
	if detectErr != nil {
		return nil, errcode.ErrInternalServer.Wrap(detectErr)
	}

	// 获取租户存储配置
	storageCfg, err := s.settingsRepo.GetStorageConfig(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	// 决定使用哪个存储后端
	useCloud := false
	if storageCfg != nil && storageCfg.Enabled {
		if storageCfg.DefaultTarget == "cloud" {
			useCloud = true
		} else {
			// DefaultTarget=local，检查本地用量
			used, dbErr := s.mediaRepo.SumFileSizeByTenant(ctx, tenantID, "local")
			if dbErr != nil {
				return nil, errcode.ErrDatabase.Wrap(dbErr)
			}
			if used+file.Size > localStorageLimitBytes {
				// 本地已满，自动切换云
				useCloud = true
			}
		}
	} else {
		// 未启用云存储，检查本地限额
		used, dbErr := s.mediaRepo.SumFileSizeByTenant(ctx, tenantID, "local")
		if dbErr != nil {
			return nil, errcode.ErrDatabase.Wrap(dbErr)
		}
		if used+file.Size > localStorageLimitBytes {
			return nil, errcode.ErrLocalStorageFull
		}
	}

	if useCloud {
		// 云存储配置必须完整
		if storageCfg == nil || !storageCfg.Enabled {
			return nil, errcode.ErrCloudStorageNotEnabled
		}
		return s.uploadToCloud(ctx, tenantID, userID, file, mimeType, storageCfg)
	}
	return s.uploadToLocal(ctx, tenantID, userID, file, mimeType)
}

// uploadToLocal 上传到本地文件系统
func (s *MediaService) uploadToLocal(ctx context.Context, tenantID string, userID primitive.ObjectID, file *multipart.FileHeader, mimeType string) (*entity.Media, *errcode.AppError) {
	cfg := config.Get()
	now := time.Now().UTC()
	dirSegment := fmt.Sprintf("%s/%s", tenantID, now.Format("2006-01"))

	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%s%s", primitive.NewObjectID().Hex(), ext)
	key := fmt.Sprintf("%s/%s", dirSegment, fileName)

	src, err := file.Open()
	if err != nil {
		return nil, errcode.ErrInternalServer.Wrap(err)
	}
	defer src.Close()

	provider := storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")
	result, err := provider.Upload(ctx, storage.UploadInput{
		Key:         key,
		Content:     src,
		Size:        file.Size,
		ContentType: mimeType,
	})
	if err != nil {
		return nil, errcode.ErrUploadStorageFail.Wrap(err)
	}

	media := &entity.Media{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		FileName:    file.Filename,
		FileURL:     result.URL,
		FileSize:    file.Size,
		MimeType:    mimeType,
		UploadedBy:  userID,
		StorageType: "local",
		StorageKey:  result.Key,
	}
	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return media, nil
}

// uploadToCloud 上传到 S3 兼容存储
func (s *MediaService) uploadToCloud(ctx context.Context, tenantID string, userID primitive.ObjectID, file *multipart.FileHeader, mimeType string, cfg *entity.StorageSettings) (*entity.Media, *errcode.AppError) {
	now := time.Now().UTC()
	dirSegment := fmt.Sprintf("%s/%s", tenantID, now.Format("2006-01"))
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%s%s", primitive.NewObjectID().Hex(), ext)
	key := fmt.Sprintf("%s/%s", dirSegment, fileName)

	src, err := file.Open()
	if err != nil {
		return nil, errcode.ErrInternalServer.Wrap(err)
	}
	defer src.Close()

	provider, err := storage.NewS3Provider(storage.S3Config{
		Endpoint:       cfg.Endpoint,
		Region:         cfg.Region,
		Bucket:         cfg.Bucket,
		AccessKeyID:    cfg.AccessKeyID,
		SecretKey:      cfg.SecretAccessKey,
		CustomDomain:   cfg.CustomDomain,
		ForcePathStyle: cfg.ForcePathStyle,
	})
	if err != nil {
		return nil, errcode.ErrS3ConnectionFailed.Wrap(err)
	}

	result, err := provider.Upload(ctx, storage.UploadInput{
		Key:         key,
		Content:     src,
		Size:        file.Size,
		ContentType: mimeType,
	})
	if err != nil {
		logger.L().Error("S3 上传失败",
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
			zap.String("bucket", cfg.Bucket),
			zap.String("endpoint", cfg.Endpoint),
			zap.String("region", cfg.Region),
			zap.String("mime", mimeType),
			zap.Int64("size", file.Size),
			zap.Error(err),
		)
		return nil, errcode.ErrUploadStorageFail.Wrap(err)
	}

	media := &entity.Media{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		FileName:    file.Filename,
		FileURL:     result.URL,
		FileSize:    file.Size,
		MimeType:    mimeType,
		UploadedBy:  userID,
		StorageType: "cloud",
		StorageKey:  result.Key,
	}
	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return media, nil
}

// List 列出租户所有媒体文件，按创建时间倒序
func (s *MediaService) List(ctx context.Context, tenantID string) ([]*entity.Media, *errcode.AppError) {
	items, err := s.mediaRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return items, nil
}

// Delete 删除媒体记录并同步删除实际存储文件（本地文件 / S3 对象）
// force=false 时先检查文档引用，有引用则拒绝删除；force=true 跳过检查强制删除
func (s *MediaService) Delete(ctx context.Context, tenantID string, mediaID primitive.ObjectID, force bool) *errcode.AppError {
	media, err := s.mediaRepo.FindByIDAndTenant(ctx, mediaID, tenantID)
	if err != nil {
		return errcode.ErrResourceNotFound
	}

	// 非强制模式：检查该媒体是否仍被文档引用
	if !force && media.FileURL != "" {
		pages, pErr := s.pageRepo.ListMinimalByTenant(ctx, tenantID)
		if pErr != nil {
			return errcode.ErrDatabase.Wrap(pErr)
		}
		for _, p := range pages {
			if strings.Contains(p.Content, media.FileURL) {
				return errcode.ErrMediaStillInUse
			}
		}
	}

	if media.StorageType == "cloud" {
		storageCfg, dbErr := s.settingsRepo.GetStorageConfig(ctx, tenantID)
		if dbErr != nil {
			return errcode.ErrDatabase.Wrap(dbErr)
		}
		if storageCfg != nil && storageCfg.Enabled && media.StorageKey != "" {
			provider, provErr := storage.NewS3Provider(storage.S3Config{
				Endpoint:       storageCfg.Endpoint,
				Region:         storageCfg.Region,
				Bucket:         storageCfg.Bucket,
				AccessKeyID:    storageCfg.AccessKeyID,
				SecretKey:      storageCfg.SecretAccessKey,
				CustomDomain:   storageCfg.CustomDomain,
				ForcePathStyle: storageCfg.ForcePathStyle,
			})
			if provErr == nil {
				if delErr := provider.Delete(ctx, media.StorageKey); delErr != nil {
					logger.L().Warn("S3 对象删除失败，继续删除 DB 记录",
						zap.String("key", media.StorageKey),
						zap.Error(delErr),
					)
				}
			}
		}
	} else if media.StorageType == "local" && media.StorageKey != "" {
		cfg := config.Get()
		provider := storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")
		if delErr := provider.Delete(ctx, media.StorageKey); delErr != nil {
			logger.L().Warn("本地文件删除失败，继续删除 DB 记录",
				zap.String("key", media.StorageKey),
				zap.Error(delErr),
			)
		}
	}

	if err := s.mediaRepo.DeleteByID(ctx, mediaID); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	return nil
}

// AuditResult 存储审计结果
type AuditResult struct {
	OrphanKeys  []string `json:"orphan_keys"`  // 存储中存在但 DB 无记录的 key
	MissingKeys []string `json:"missing_keys"` // DB 有记录但存储中不存在的 key
	MatchCount  int      `json:"match_count"`  // DB 与存储匹配的文件数
}

// AuditStorage 对照 DB 媒体记录与实际存储文件，检测孤立文件和缺失文件
func (s *MediaService) AuditStorage(ctx context.Context, tenantID string, storageType string) (*AuditResult, *errcode.AppError) {
	// 1. 获取 DB 中该租户指定存储类型的所有媒体记录
	allMedia, err := s.mediaRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	dbKeySet := make(map[string]struct{})
	for _, m := range allMedia {
		if m.StorageType == storageType && m.StorageKey != "" {
			dbKeySet[m.StorageKey] = struct{}{}
		}
	}

	// 2. 获取实际存储中的 key 列表
	var provider storage.Provider
	if storageType == "cloud" {
		storageCfg, cfgErr := s.settingsRepo.GetStorageConfig(ctx, tenantID)
		if cfgErr != nil {
			return nil, errcode.ErrDatabase.Wrap(cfgErr)
		}
		if storageCfg == nil || !storageCfg.Enabled {
			return nil, errcode.ErrCloudStorageNotEnabled
		}
		s3Provider, provErr := storage.NewS3Provider(storage.S3Config{
			Endpoint:       storageCfg.Endpoint,
			Region:         storageCfg.Region,
			Bucket:         storageCfg.Bucket,
			AccessKeyID:    storageCfg.AccessKeyID,
			SecretKey:      storageCfg.SecretAccessKey,
			CustomDomain:   storageCfg.CustomDomain,
			ForcePathStyle: storageCfg.ForcePathStyle,
		})
		if provErr != nil {
			return nil, errcode.ErrS3ConnectionFailed.Wrap(provErr)
		}
		provider = s3Provider
	} else {
		cfg := config.Get()
		provider = storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")
	}

	storageKeys, listErr := provider.ListKeys(ctx, tenantID+"/")
	if listErr != nil {
		logger.L().Error("审计：列举存储 key 失败",
			zap.String("tenant_id", tenantID),
			zap.String("storage_type", storageType),
			zap.Error(listErr),
		)
		return nil, errcode.ErrInternalServer.Wrap(listErr)
	}

	storageKeySet := make(map[string]struct{}, len(storageKeys))
	for _, k := range storageKeys {
		storageKeySet[k] = struct{}{}
	}

	// 3. 对比
	result := &AuditResult{}
	for k := range storageKeySet {
		if _, exists := dbKeySet[k]; exists {
			result.MatchCount++
		} else {
			result.OrphanKeys = append(result.OrphanKeys, k)
		}
	}
	for k := range dbKeySet {
		if _, exists := storageKeySet[k]; !exists {
			result.MissingKeys = append(result.MissingKeys, k)
		}
	}
	return result, nil
}

// MediaUsageRef 媒体文件被某个文档页引用的来源信息
type MediaUsageRef struct {
	ThemeID      string `json:"theme_id"`
	ThemeName    string `json:"theme_name"`
	ThemeSlug    string `json:"theme_slug"`
	ThemeDeleted bool   `json:"theme_deleted"` // 主题已被删除（pages/versions 残留）
	VersionID    string `json:"version_id"`
	VersionName  string `json:"version_name"`
	PageID       string `json:"page_id"`
	PageTitle    string `json:"page_title"`
	PageSlug     string `json:"page_slug"`
}

// MediaUsageResult 媒体使用来源查询结果
type MediaUsageResult struct {
	Usages map[string][]MediaUsageRef `json:"usages"` // media_id → []引用来源
}

// GetUsage 扫描租户所有页面内容，找出每个媒体文件的引用来源
// 查询流程：media_urls → pages（内容扫描）→ versions（批量）→ themes（批量）
func (s *MediaService) GetUsage(ctx context.Context, tenantID string) (*MediaUsageResult, *errcode.AppError) {
	allMedia, err := s.mediaRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	result := &MediaUsageResult{Usages: map[string][]MediaUsageRef{}}
	if len(allMedia) == 0 {
		return result, nil
	}

	// url → mediaID 映射，用于内容扫描定位
	urlToMediaID := make(map[string]string, len(allMedia))
	for _, m := range allMedia {
		if m.FileURL != "" {
			urlToMediaID[m.FileURL] = m.ID.Hex()
		}
	}

	pages, pErr := s.pageRepo.ListMinimalByTenant(ctx, tenantID)
	if pErr != nil {
		return nil, errcode.ErrDatabase.Wrap(pErr)
	}

	// 内容扫描：每个 page 检测是否包含各媒体 URL
	type pageRef struct {
		versionID primitive.ObjectID
		pageID    string
		pageTitle string
		pageSlug  string
	}
	mediaUsages := make(map[string][]pageRef)
	for _, p := range pages {
		for url, mediaID := range urlToMediaID {
			if strings.Contains(p.Content, url) {
				mediaUsages[mediaID] = append(mediaUsages[mediaID], pageRef{
					versionID: p.VersionID,
					pageID:    p.ID.Hex(),
					pageTitle: p.Title,
					pageSlug:  p.Slug,
				})
			}
		}
	}

	if len(mediaUsages) == 0 {
		return result, nil
	}

	// 批量查询涉及的 versions
	versionIDSet := make(map[primitive.ObjectID]struct{})
	for _, refs := range mediaUsages {
		for _, r := range refs {
			versionIDSet[r.versionID] = struct{}{}
		}
	}
	versionIDs := make([]primitive.ObjectID, 0, len(versionIDSet))
	for id := range versionIDSet {
		versionIDs = append(versionIDs, id)
	}
	versionMap, vErr := s.versionRepo.FindByIDsMap(ctx, versionIDs)
	if vErr != nil {
		return nil, errcode.ErrDatabase.Wrap(vErr)
	}

	// 批量查询涉及的 themes
	themeIDSet := make(map[primitive.ObjectID]struct{})
	for _, v := range versionMap {
		themeIDSet[v.ThemeID] = struct{}{}
	}
	themeIDs := make([]primitive.ObjectID, 0, len(themeIDSet))
	for id := range themeIDSet {
		themeIDs = append(themeIDs, id)
	}
	themeMap, tErr := s.themeRepo.FindByIDsMap(ctx, themeIDs)
	if tErr != nil {
		return nil, errcode.ErrDatabase.Wrap(tErr)
	}

	// 组装结果
	for mediaID, refs := range mediaUsages {
		var usageRefs []MediaUsageRef
		for _, ref := range refs {
			v, ok := versionMap[ref.versionID]
			if !ok {
				continue
			}
			ur := MediaUsageRef{
				VersionID:   ref.versionID.Hex(),
				VersionName: v.Name,
				PageID:      ref.pageID,
				PageTitle:   ref.pageTitle,
				PageSlug:    ref.pageSlug,
				ThemeID:     v.ThemeID.Hex(),
			}
			theme, exists := themeMap[v.ThemeID]
			if exists {
				ur.ThemeName = theme.Name
				ur.ThemeSlug = theme.Slug
				ur.ThemeDeleted = false
			} else {
				ur.ThemeDeleted = true
			}
			usageRefs = append(usageRefs, ur)
		}
		if len(usageRefs) > 0 {
			result.Usages[mediaID] = usageRefs
		}
	}
	return result, nil
}

// CleanupUnusedResult 批量清理未使用媒体文件结果
type CleanupUnusedResult struct {
	DeletedCount int `json:"deleted_count"`
}

// CleanupUnused 扫描并删除所有未被任何文档页引用的媒体文件（DB + 存储）
func (s *MediaService) CleanupUnused(ctx context.Context, tenantID string) (*CleanupUnusedResult, *errcode.AppError) {
	usageResult, appErr := s.GetUsage(ctx, tenantID)
	if appErr != nil {
		return nil, appErr
	}
	allMedia, err := s.mediaRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	result := &CleanupUnusedResult{}
	for _, m := range allMedia {
		if _, used := usageResult.Usages[m.ID.Hex()]; !used {
			if delErr := s.Delete(ctx, tenantID, m.ID, true); delErr == nil {
				result.DeletedCount++
			}
		}
	}
	return result, nil
}

// DeleteOrphan 删除存储中的孤立文件（仅删除存储对象，不涉及 DB 记录）
func (s *MediaService) DeleteOrphan(ctx context.Context, tenantID string, storageType string, key string) *errcode.AppError {
	if storageType == "cloud" {
		storageCfg, cfgErr := s.settingsRepo.GetStorageConfig(ctx, tenantID)
		if cfgErr != nil {
			return errcode.ErrDatabase.Wrap(cfgErr)
		}
		if storageCfg == nil || !storageCfg.Enabled {
			return errcode.ErrCloudStorageNotEnabled
		}
		provider, provErr := storage.NewS3Provider(storage.S3Config{
			Endpoint:       storageCfg.Endpoint,
			Region:         storageCfg.Region,
			Bucket:         storageCfg.Bucket,
			AccessKeyID:    storageCfg.AccessKeyID,
			SecretKey:      storageCfg.SecretAccessKey,
			CustomDomain:   storageCfg.CustomDomain,
			ForcePathStyle: storageCfg.ForcePathStyle,
		})
		if provErr != nil {
			return errcode.ErrS3ConnectionFailed.Wrap(provErr)
		}
		if delErr := provider.Delete(ctx, key); delErr != nil {
			return errcode.ErrInternalServer.Wrap(delErr)
		}
	} else {
		cfg := config.Get()
		provider := storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")
		if delErr := provider.Delete(ctx, key); delErr != nil {
			return errcode.ErrInternalServer.Wrap(delErr)
		}
	}
	return nil
}
