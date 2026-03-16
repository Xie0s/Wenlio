// Package service 业务逻辑层
//
// 职责：实现章节管理业务逻辑（CRUD、排序、级联删除）
// 对外接口：SectionService
package service

import (
	"context"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SectionService 章节业务
type SectionService struct {
	sectionRepo    *repository.SectionRepo
	pageRepo       *repository.PageRepo
	versionRepo    *repository.VersionRepo
	themeRepo      *repository.ThemeRepo
	versionService *VersionService
}

// NewSectionService 创建 SectionService
func NewSectionService() *SectionService {
	return &SectionService{
		sectionRepo:    repository.NewSectionRepo(),
		pageRepo:       repository.NewPageRepo(),
		versionRepo:    repository.NewVersionRepo(),
		themeRepo:      repository.NewThemeRepo(),
		versionService: NewVersionService(),
	}
}

// Create 创建章节
func (s *SectionService) Create(ctx context.Context, tenantID string, versionID primitive.ObjectID, req *dto.CreateSectionReq) (*entity.Section, *errcode.AppError) {
	if appErr := s.versionService.CheckEditable(ctx, versionID); appErr != nil {
		return nil, appErr
	}
	version, err := s.versionRepo.FindByIDTyped(ctx, versionID)
	if err != nil {
		return nil, errcode.ErrVersionNotFound
	}
	if version.TenantID != tenantID {
		return nil, errcode.ErrTenantMismatch
	}
	theme, err := s.themeRepo.FindByIDAndTenant(ctx, version.ThemeID, tenantID)
	if err != nil {
		return nil, errcode.ErrThemeNotFound
	}
	if theme.Deleting {
		return nil, errcode.ErrThemeDeleting
	}

	maxSort, _ := s.sectionRepo.MaxSortOrder(ctx, versionID)
	section := &entity.Section{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		VersionID:   versionID,
		Title:       req.Title,
		SortOrder:   maxSort + 1,
	}
	if err := s.sectionRepo.Create(ctx, section); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return section, nil
}

// GetByID 获取章节详情
func (s *SectionService) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Section, *errcode.AppError) {
	section, err := s.sectionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, errcode.ErrSectionNotFound
	}
	return section, nil
}

// ListByVersion 版本下章节列表
func (s *SectionService) ListByVersion(ctx context.Context, tenantID string, versionID primitive.ObjectID) ([]*entity.Section, *errcode.AppError) {
	sections, err := s.sectionRepo.ListByVersion(ctx, tenantID, versionID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return sections, nil
}

// Update 更新章节
func (s *SectionService) Update(ctx context.Context, id primitive.ObjectID, req *dto.UpdateSectionReq) *errcode.AppError {
	section, err := s.sectionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrSectionNotFound
	}
	if appErr := s.versionService.CheckEditable(ctx, section.VersionID); appErr != nil {
		return errcode.ErrSectionVersionArchived
	}
	version, err := s.versionRepo.FindByIDTyped(ctx, section.VersionID)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	theme, err := s.themeRepo.FindByIDAndTenant(ctx, version.ThemeID, version.TenantID)
	if err != nil {
		return errcode.ErrThemeNotFound
	}
	if theme.Deleting {
		return errcode.ErrThemeDeleting
	}
	return toAppError(s.sectionRepo.UpdateByID(ctx, id, map[string]interface{}{"title": req.Title}))
}

// Delete 删除章节（级联删除文档页）
func (s *SectionService) Delete(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	section, err := s.sectionRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrSectionNotFound
	}
	if appErr := s.versionService.CheckEditable(ctx, section.VersionID); appErr != nil {
		return errcode.ErrSectionVersionArchived
	}

	// 级联删除章节下所有文档页
	_ = s.pageRepo.DeleteBySection(ctx, id)
	return toAppError(s.sectionRepo.DeleteByID(ctx, id))
}

// Sort 批量更新章节排序
func (s *SectionService) Sort(ctx context.Context, items []dto.SortItem) *errcode.AppError {
	for _, item := range items {
		oid, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			continue
		}
		_ = s.sectionRepo.UpdateByID(ctx, oid, map[string]interface{}{"sort_order": item.SortOrder})
	}
	return nil
}
