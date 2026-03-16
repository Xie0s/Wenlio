// Package service 业务逻辑层
//
// 职责：实现租户管理相关业务逻辑（CRUD、封禁/解封）
// 对外接口：TenantService
package service

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"docplatform/internal/config"
	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/constants"
	"docplatform/pkg/ctxutil"
	"docplatform/pkg/errcode"
	"docplatform/pkg/logger"
	"docplatform/pkg/storage"
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// TenantService 租户业务
type TenantService struct {
	tenantRepo   *repository.TenantRepo
	userRepo     *repository.UserRepo
	themeRepo    *repository.ThemeRepo
	versionRepo  *repository.VersionRepo
	sectionRepo  *repository.SectionRepo
	pageRepo     *repository.PageRepo
	commentRepo  *repository.CommentRepo
	mediaRepo    *repository.MediaRepo
	settingsRepo *repository.TenantSettingsRepo
}

var tenantDeleteTasks = struct {
	mu      sync.Mutex
	running map[string]struct{}
}{
	running: map[string]struct{}{},
}

// NewTenantService 创建 TenantService
func NewTenantService() *TenantService {
	return &TenantService{
		tenantRepo:   repository.NewTenantRepo(),
		userRepo:     repository.NewUserRepo(),
		themeRepo:    repository.NewThemeRepo(),
		versionRepo:  repository.NewVersionRepo(),
		sectionRepo:  repository.NewSectionRepo(),
		pageRepo:     repository.NewPageRepo(),
		commentRepo:  repository.NewCommentRepo(),
		mediaRepo:    repository.NewMediaRepo(),
		settingsRepo: repository.NewTenantSettingsRepo(),
	}
}

var tenantIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$`)

// Create 创建租户（含初始管理员）
func (s *TenantService) Create(ctx context.Context, req *dto.CreateTenantReq) (*entity.Tenant, *errcode.AppError) {
	// 保留词检查
	if _, reserved := constants.ReservedTenantIDs[req.ID]; reserved {
		return nil, errcode.ErrTenantIDReserved
	}
	// 格式检查
	if !tenantIDPattern.MatchString(req.ID) {
		return nil, errcode.ErrTenantIDInvalid
	}
	// 唯一性检查
	if s.tenantRepo.Exists(ctx, req.ID) {
		return nil, errcode.ErrTenantIDExists
	}
	if existing, _ := s.userRepo.FindByUsername(ctx, req.AdminUsername); existing != nil {
		return nil, errcode.ErrUsernameExists
	}

	// 创建租户
	tenant := &entity.Tenant{
		ID:      req.ID,
		Name:    req.Name,
		LogoURL: req.LogoURL,
		Status:  constants.TenantStatusActive,
	}
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	// 创建初始管理员
	hashed, err := utils.HashPassword(req.AdminPassword)
	if err != nil {
		return nil, errcode.ErrInternalServer
	}
	admin := &entity.User{
		TenantID: req.ID,
		Username: req.AdminUsername,
		Password: hashed,
		Name:     req.AdminName,
		Email:    req.AdminEmail,
		Role:     constants.RoleTenantAdmin,
		Status:   constants.UserStatusActive,
	}
	if err := s.userRepo.Create(ctx, admin); err != nil {
		if rollbackErr := s.tenantRepo.DeleteByStringID(ctx, req.ID); rollbackErr != nil {
			logger.L().Error("租户管理员创建失败，回滚租户失败",
				zap.String("tenant_id", req.ID),
				zap.Error(rollbackErr),
			)
		}
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	logger.L().Info("租户创建成功",
		zap.String("tenant_id", req.ID),
		zap.String("admin_username", req.AdminUsername),
	)
	return tenant, nil
}

// GetByID 获取租户详情（排除 homepage 字段，避免数据膨胀）
func (s *TenantService) GetByID(ctx context.Context, id string) (*entity.Tenant, *errcode.AppError) {
	tenant, err := s.tenantRepo.FindByIDWithoutHomepage(ctx, id)
	if err != nil {
		return nil, errcode.ErrTenantNotFound
	}
	return tenant, nil
}

func (s *TenantService) GetPublicByID(ctx context.Context, id string) (*dto.PublicTenantDTO, *errcode.AppError) {
	tenant, err := s.tenantRepo.FindPublicByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return nil, errcode.ErrTenantDeleting
	}
	if tenant.Status != constants.TenantStatusActive {
		return nil, errcode.ErrTenantSuspended
	}
	resp := &dto.PublicTenantDTO{
		ID:        tenant.ID,
		Name:      tenant.Name,
		LogoURL:   tenant.LogoURL,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}
	if tenant.Homepage != nil && tenant.Homepage.Published != nil {
		resp.BrowserTitle = tenant.Homepage.Published.Global.BrowserTitle
		resp.BrowserIconURL = tenant.Homepage.Published.Global.BrowserIconURL
	}
	return resp, nil
}

// List 租户列表（含各租户用户数统计）
func (s *TenantService) List(ctx context.Context, page, pageSize int, keyword string) ([]*dto.TenantListItemDTO, int64, *errcode.AppError) {
	tenants, total, err := s.tenantRepo.List(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, 0, errcode.ErrDatabase.Wrap(err)
	}

	ids := make([]string, len(tenants))
	for i, t := range tenants {
		ids[i] = t.ID
	}
	counts, _ := s.userRepo.CountByTenantIDs(ctx, ids)
	admins, _ := s.userRepo.FindFirstAdminByTenantIDs(ctx, ids)

	items := make([]*dto.TenantListItemDTO, len(tenants))
	for i, t := range tenants {
		adminInfo := admins[t.ID]
		items[i] = &dto.TenantListItemDTO{
			ID:             t.ID,
			Name:           t.Name,
			LogoURL:        t.LogoURL,
			Status:         t.Status,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
			UserCount:      counts[t.ID],
			AdminUsername:  adminInfo[0],
			AdminName:      adminInfo[1],
			AdminAvatarURL: adminInfo[2],
		}
	}
	return items, total, nil
}

// Update 更新租户信息
func (s *TenantService) Update(ctx context.Context, id string, req *dto.UpdateTenantReq) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.LogoURL != "" {
		update["logo_url"] = req.LogoURL
	}
	if len(update) == 0 {
		return nil
	}
	return toAppError(s.tenantRepo.UpdateByStringID(ctx, id, update))
}

// Suspend 封禁租户
func (s *TenantService) Suspend(ctx context.Context, id string) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	if tenant.Status == constants.TenantStatusSuspended {
		return errcode.ErrTenantSuspended
	}
	return toAppError(s.tenantRepo.UpdateByStringID(ctx, id, map[string]interface{}{
		"status": constants.TenantStatusSuspended,
	}))
}

// Activate 解封租户
func (s *TenantService) Activate(ctx context.Context, id string) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	if tenant.Status == constants.TenantStatusActive {
		return errcode.ErrTenantAlreadyActive
	}
	return toAppError(s.tenantRepo.UpdateByStringID(ctx, id, map[string]interface{}{
		"status": constants.TenantStatusActive,
	}))
}

func (s *TenantService) Delete(ctx context.Context, id string) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		s.enqueueDeleteTask(tenant.ID)
		return nil
	}
	marked, err := s.tenantRepo.MarkDeletingIfAllowed(ctx, tenant.ID)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if !marked {
		s.enqueueDeleteTask(tenant.ID)
		return nil
	}
	s.enqueueDeleteTask(tenant.ID)
	return nil
}

func (s *TenantService) RecoverDeleting(ctx context.Context) *errcode.AppError {
	tenants, err := s.tenantRepo.ListDeleting(ctx)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	for _, tenant := range tenants {
		s.enqueueDeleteTask(tenant.ID)
	}
	return nil
}

func (s *TenantService) enqueueDeleteTask(tenantID string) {
	tenantDeleteTasks.mu.Lock()
	if _, exists := tenantDeleteTasks.running[tenantID]; exists {
		tenantDeleteTasks.mu.Unlock()
		return
	}
	tenantDeleteTasks.running[tenantID] = struct{}{}
	tenantDeleteTasks.mu.Unlock()

	go func() {
		defer func() {
			tenantDeleteTasks.mu.Lock()
			delete(tenantDeleteTasks.running, tenantID)
			tenantDeleteTasks.mu.Unlock()
		}()

		deleteCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if appErr := s.executeDeleteTask(deleteCtx, tenantID); appErr != nil {
			logger.L().Error("异步删除租户失败",
				zap.String("tenant_id", tenantID),
				zap.Int("code", appErr.Code),
				zap.String("message", appErr.Message),
				zap.Error(appErr.Internal),
			)
			return
		}
		logger.L().Info("异步删除租户完成", zap.String("tenant_id", tenantID))
	}()
}

func (s *TenantService) executeDeleteTask(ctx context.Context, tenantID string) *errcode.AppError {
	logger.L().Info("开始异步删除租户", zap.String("tenant_id", tenantID))

	filters := bson.M{"tenant_id": tenantID}
	if _, err := s.userRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if _, err := s.themeRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if _, err := s.versionRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if _, err := s.sectionRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if _, err := s.pageRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if _, err := s.commentRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	// 清理 S3 云存储文件（best-effort，出错只记录日志，不中断删除流程）
	s.cleanupCloudFiles(ctx, tenantID)

	if _, err := s.mediaRepo.DeleteMany(ctx, filters); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if err := s.tenantRepo.DeleteByStringID(ctx, tenantID); err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}

	uploadDir := filepath.Join(config.Get().Storage.LocalPath, tenantID)
	if err := os.RemoveAll(uploadDir); err != nil {
		logger.L().Error("删除租户上传目录失败",
			zap.String("tenant_id", tenantID),
			zap.String("upload_dir", uploadDir),
			zap.Error(err),
		)
	}

	return nil
}

// cleanupCloudFiles 清理租户在 S3 中的所有云存储文件（best-effort）
func (s *TenantService) cleanupCloudFiles(ctx context.Context, tenantID string) {
	cfg, err := s.settingsRepo.GetStorageConfig(ctx, tenantID)
	if err != nil || cfg == nil || !cfg.Enabled {
		return
	}
	if cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return
	}

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
		logger.L().Warn("清理云存储：创建 S3 Provider 失败",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return
	}

	cloudMedia, err := s.mediaRepo.FindCloudByTenant(ctx, tenantID)
	if err != nil {
		logger.L().Warn("清理云存储：查询云媒体记录失败",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return
	}

	var failCount int
	for _, media := range cloudMedia {
		if media.StorageKey == "" {
			continue
		}
		if delErr := provider.Delete(ctx, media.StorageKey); delErr != nil {
			failCount++
			logger.L().Warn("清理云存储：删除对象失败",
				zap.String("tenant_id", tenantID),
				zap.String("key", media.StorageKey),
				zap.Error(delErr),
			)
		}
	}

	logger.L().Info("清理云存储完成",
		zap.String("tenant_id", tenantID),
		zap.Int("total", len(cloudMedia)),
		zap.Int("failed", failCount),
	)
}

// ============================================================
// 首页个性化
// ============================================================

// GetHomepagePublished 获取已发布首页配置（读者端）
func (s *TenantService) GetHomepagePublished(ctx context.Context, tenantID string) (*entity.HomepageLayout, *errcode.AppError) {
	if appErr := s.EnsureActive(ctx, tenantID); appErr != nil {
		return nil, appErr
	}
	hp, err := s.tenantRepo.GetHomepage(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrTenantNotFound
	}
	if hp == nil {
		return nil, nil
	}
	return hp.Published, nil
}

// GetHomepageDraft 获取首页草稿 + 已发布配置（管理端）
func (s *TenantService) GetHomepageDraft(ctx context.Context, tenantID string) (*entity.TenantHomepage, *errcode.AppError) {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return nil, errcode.ErrTenantDeleting
	}
	hp, err := s.tenantRepo.GetHomepage(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrTenantNotFound
	}
	return hp, nil
}

// SaveHomepageDraft 保存首页草稿
func (s *TenantService) SaveHomepageDraft(ctx context.Context, tenantID string, layout *entity.HomepageLayout) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, tenantID)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	return toAppError(s.tenantRepo.SaveHomepageDraft(ctx, tenantID, layout))
}

// PublishHomepage 发布首页（将 draft 复制到 published）
func (s *TenantService) PublishHomepage(ctx context.Context, tenantID string) *errcode.AppError {
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, tenantID)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	hp, err := s.tenantRepo.GetHomepage(ctx, tenantID)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if hp == nil || hp.Draft == nil {
		return errcode.ErrHomepageNoDraft
	}
	return toAppError(s.tenantRepo.PublishHomepage(ctx, tenantID, hp.Draft))
}

func (s *TenantService) EnsureActive(ctx context.Context, tenantID string) *errcode.AppError {
	// 优先读取 JWTAuth 中间件已缓存的租户状态，命中则跳过 DB 查询
	if status, ok := ctxutil.TenantStatus(ctx); ok {
		if status == constants.TenantStatusDeleting {
			return errcode.ErrTenantDeleting
		}
		if status != constants.TenantStatusActive {
			return errcode.ErrTenantSuspended
		}
		return nil
	}
	// 未命中缓存（公开接口无 JWT），回退查库
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, tenantID)
	if err != nil {
		return errcode.ErrTenantNotFound
	}
	if tenant.Status == constants.TenantStatusDeleting {
		return errcode.ErrTenantDeleting
	}
	if tenant.Status != constants.TenantStatusActive {
		return errcode.ErrTenantSuspended
	}
	return nil
}
