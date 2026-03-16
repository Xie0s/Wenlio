// Package service 业务逻辑层
//
// 职责：实现版本管理业务逻辑（创建、发布、归档、克隆、设为默认）
// 对外接口：VersionService
package service

import (
	"context"
	"time"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/constants"
	"docplatform/pkg/errcode"
	"docplatform/pkg/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// VersionService 版本业务
type VersionService struct {
	versionRepo *repository.VersionRepo
	themeRepo   *repository.ThemeRepo
	sectionRepo *repository.SectionRepo
	pageRepo    *repository.PageRepo
	tenantRepo  *repository.TenantRepo
}

// NewVersionService 创建 VersionService
func NewVersionService() *VersionService {
	return &VersionService{
		versionRepo: repository.NewVersionRepo(),
		themeRepo:   repository.NewThemeRepo(),
		sectionRepo: repository.NewSectionRepo(),
		pageRepo:    repository.NewPageRepo(),
		tenantRepo:  repository.NewTenantRepo(),
	}
}

// Create 创建版本
func (s *VersionService) Create(ctx context.Context, tenantID string, themeID, userID primitive.ObjectID, req *dto.CreateVersionReq) (*entity.Version, *errcode.AppError) {
	theme, err := s.themeRepo.FindByIDAndTenant(ctx, themeID, tenantID)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	if theme.Deleting {
		return nil, errcode.ErrThemeDeleting
	}

	version := &entity.Version{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		ThemeID:     themeID,
		Name:        req.Name,
		Label:       req.Label,
		Status:      constants.VersionStatusDraft,
		IsDefault:   false,
		CreatedBy:   userID,
	}
	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return version, nil
}

// GetByID 获取版本详情
func (s *VersionService) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Version, *errcode.AppError) {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, errcode.ErrVersionNotFound
	}
	return version, nil
}

// ListByTheme 主题下版本列表
func (s *VersionService) ListByTheme(ctx context.Context, tenantID string, themeID primitive.ObjectID) ([]*entity.Version, *errcode.AppError) {
	versions, err := s.versionRepo.ListByTheme(ctx, tenantID, themeID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return versions, nil
}

// Update 更新版本信息
func (s *VersionService) Update(ctx context.Context, id primitive.ObjectID, req *dto.UpdateVersionReq) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status == constants.VersionStatusArchived {
		return errcode.ErrVersionArchived
	}
	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Label != "" {
		update["label"] = req.Label
	}
	return toAppError(s.versionRepo.UpdateByID(ctx, id, update))
}

// Publish 发布版本（将 draft → published，并批量发布文档页）
func (s *VersionService) Publish(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status != constants.VersionStatusDraft {
		return errcode.ErrVersionNotDraft
	}

	now := time.Now().UTC()
	if appErr := toAppError(s.versionRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status":       constants.VersionStatusPublished,
		"published_at": now,
	})); appErr != nil {
		return appErr
	}

	// 批量发布该版本下所有 draft 文档页
	count, _ := s.pageRepo.BatchPublish(ctx, id)
	logger.L().Info("版本发布成功",
		zap.String("version_id", id.Hex()),
		zap.Int64("pages_published", count),
	)
	return nil
}

// Unpublish 取消发布版本（published → draft），同步将所有已发布文档页回退为 draft
func (s *VersionService) Unpublish(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status == constants.VersionStatusDraft {
		return errcode.ErrVersionAlreadyDraft
	}
	if version.Status == constants.VersionStatusArchived {
		return errcode.ErrVersionArchived
	}
	if appErr := toAppError(s.versionRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.VersionStatusDraft,
	})); appErr != nil {
		return appErr
	}
	_, _ = s.pageRepo.BatchUnpublish(ctx, id)
	return nil
}

// Unarchive 取消归档版本（archived → published）
func (s *VersionService) Unarchive(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status != constants.VersionStatusArchived {
		return errcode.ErrVersionNotArchived
	}
	return toAppError(s.versionRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.VersionStatusPublished,
	}))
}

// Archive 归档版本（published → archived）
func (s *VersionService) Archive(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status != constants.VersionStatusPublished {
		return errcode.ErrVersionNotPublished
	}
	return toAppError(s.versionRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.VersionStatusArchived,
	}))
}

// SetDefault 设为默认版本
func (s *VersionService) SetDefault(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status != constants.VersionStatusPublished {
		return errcode.ErrVersionSetDefaultNotPub
	}

	// 清除同主题下其他版本的默认标记
	_ = s.versionRepo.ClearDefault(ctx, version.TenantID, version.ThemeID)

	return toAppError(s.versionRepo.UpdateByID(ctx, id, map[string]interface{}{
		"is_default": true,
	}))
}

// Clone 克隆版本（深拷贝章节 + 文档页）
func (s *VersionService) Clone(ctx context.Context, sourceID primitive.ObjectID, userID primitive.ObjectID, req *dto.CloneVersionReq) (*entity.Version, *errcode.AppError) {
	source, err := s.versionRepo.FindByIDTyped(ctx, sourceID)
	if err != nil {
		return nil, errcode.ErrVersionNotFound
	}

	// 创建新版本
	newVersion := &entity.Version{
		TenantModel: entity.TenantModel{TenantID: source.TenantID},
		ThemeID:     source.ThemeID,
		Name:        req.Name,
		Label:       req.Label,
		Status:      constants.VersionStatusDraft,
		IsDefault:   false,
		CreatedBy:   userID,
	}
	if err := s.versionRepo.Create(ctx, newVersion); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	// 查询源版本下所有章节
	sections, err := s.sectionRepo.FindByVersionID(ctx, sourceID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	// 逐章节复制
	for _, section := range sections {
		newSection := &entity.Section{
			TenantModel: entity.TenantModel{TenantID: section.TenantID},
			VersionID:   newVersion.ID,
			Title:       section.Title,
			SortOrder:   section.SortOrder,
		}
		if err := s.sectionRepo.Create(ctx, newSection); err != nil {
			continue
		}

		pages, err := s.pageRepo.FindBySectionID(ctx, section.ID)
		if err != nil {
			continue
		}

		for _, page := range pages {
			newPage := &entity.Page{
				TenantModel: entity.TenantModel{TenantID: page.TenantID},
				VersionID:   newVersion.ID,
				SectionID:   newSection.ID,
				Title:       page.Title,
				Slug:        page.Slug,
				Content:     page.Content,
				Status:      constants.PageStatusDraft,
				SortOrder:   page.SortOrder,
				CreatedBy:   userID,
			}
			_ = s.pageRepo.Create(ctx, newPage)
		}
	}

	logger.L().Info("版本克隆成功",
		zap.String("source_id", sourceID.Hex()),
		zap.String("new_id", newVersion.ID.Hex()),
	)
	return newVersion, nil
}

// Delete 删除版本（级联删除章节和文档页）
func (s *VersionService) Delete(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.IsDefault {
		return errcode.ErrVersionDefaultCannotDel
	}

	// 查询版本下所有章节，级联删除文档页
	sections, _ := s.sectionRepo.FindByVersionID(ctx, id)
	for _, section := range sections {
		_ = s.pageRepo.DeleteBySection(ctx, section.ID)
	}
	_ = s.sectionRepo.DeleteByVersion(ctx, id)

	return toAppError(s.versionRepo.DeleteByID(ctx, id))
}

// ListPublishedByTheme 查询主题下已发布/已归档版本（公开接口）
func (s *VersionService) ListPublishedByTheme(ctx context.Context, themeID primitive.ObjectID) ([]*entity.Version, *errcode.AppError) {
	theme, err := s.themeRepo.FindByIDTyped(ctx, themeID)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, theme.TenantID); appErr != nil {
		return nil, appErr
	}
	versions, err := s.versionRepo.ListPublishedByTheme(ctx, themeID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return versions, nil
}

// FindByThemeAndName 按主题 ID 和版本名查询已发布版本（Raw Markdown 链式定位用）
func (s *VersionService) FindByThemeAndName(ctx context.Context, themeID primitive.ObjectID, name string) (*entity.Version, *errcode.AppError) {
	version, err := s.versionRepo.FindByThemeAndName(ctx, themeID, name)
	if err != nil {
		return nil, errcode.ErrVersionNotFound
	}
	return version, nil
}

// CheckEditable 校验版本是否可编辑（非 archived）
func (s *VersionService) CheckEditable(ctx context.Context, versionID primitive.ObjectID) *errcode.AppError {
	version, err := s.versionRepo.FindByIDTyped(ctx, versionID)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status == constants.VersionStatusArchived {
		return errcode.ErrVersionArchived
	}
	return nil
}
