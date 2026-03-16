// Package service 业务逻辑层
//
// 职责：实现文档页管理业务逻辑（CRUD、发布/下线、自动保存、排序）
// 对外接口：PageService
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
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// PageService 文档页业务
type PageService struct {
	pageRepo       *repository.PageRepo
	versionRepo    *repository.VersionRepo
	themeRepo      *repository.ThemeRepo
	tenantRepo     *repository.TenantRepo
	versionService *VersionService
}

// NewPageService 创建 PageService
func NewPageService() *PageService {
	return &PageService{
		pageRepo:       repository.NewPageRepo(),
		versionRepo:    repository.NewVersionRepo(),
		themeRepo:      repository.NewThemeRepo(),
		tenantRepo:     repository.NewTenantRepo(),
		versionService: NewVersionService(),
	}
}

// Create 创建文档页
func (s *PageService) Create(ctx context.Context, tenantID string, sectionID primitive.ObjectID, versionID primitive.ObjectID, userID primitive.ObjectID, req *dto.CreatePageReq) (*entity.Page, *errcode.AppError) {
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

	if !utils.ValidateSlug(req.Slug) {
		return nil, errcode.ErrPageSlugInvalid
	}
	// Slug 唯一性校验（同版本内）
	if existing, _ := s.pageRepo.FindByVersionAndSlug(ctx, versionID, req.Slug); existing != nil {
		return nil, errcode.ErrPageSlugExists
	}

	maxSort, _ := s.pageRepo.MaxSortOrder(ctx, sectionID)
	page := &entity.Page{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		VersionID:   versionID,
		SectionID:   sectionID,
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Status:      constants.PageStatusDraft,
		SortOrder:   maxSort + 1,
		CreatedBy:   userID,
	}
	if err := s.pageRepo.Create(ctx, page); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return page, nil
}

// GetByID 获取文档页详情
func (s *PageService) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Page, *errcode.AppError) {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return nil, errcode.ErrPageNotFound
	}
	return page, nil
}

// GetPublished 获取已发布文档页（公开接口）
func (s *PageService) GetPublished(ctx context.Context, id primitive.ObjectID) (*entity.Page, *errcode.AppError) {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil || page.Status != constants.PageStatusPublished {
		return nil, errcode.ErrPageNotFound
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, page.TenantID); appErr != nil {
		return nil, appErr
	}
	return page, nil
}

// GetPublishedBySlug 按版本 ID + slug 获取已发布文档页（Raw Markdown 用）
func (s *PageService) GetPublishedBySlug(ctx context.Context, versionID primitive.ObjectID, slug string) (*entity.Page, *errcode.AppError) {
	page, err := s.pageRepo.FindByVersionAndSlug(ctx, versionID, slug)
	if err != nil || page.Status != constants.PageStatusPublished {
		return nil, errcode.ErrPageNotFound
	}
	return page, nil
}

// ListBySection 章节下文档页列表
func (s *PageService) ListBySection(ctx context.Context, tenantID string, sectionID primitive.ObjectID) ([]*entity.Page, *errcode.AppError) {
	pages, err := s.pageRepo.ListBySection(ctx, tenantID, sectionID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return pages, nil
}

// Update 全量更新文档页
func (s *PageService) Update(ctx context.Context, id primitive.ObjectID, req *dto.UpdatePageReq) *errcode.AppError {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrPageNotFound
	}
	if appErr := s.versionService.CheckEditable(ctx, page.VersionID); appErr != nil {
		return errcode.ErrPageVersionArchived
	}
	version, err := s.versionRepo.FindByIDTyped(ctx, page.VersionID)
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

	// Slug 变更时校验唯一性
	if req.Slug != page.Slug {
		if !utils.ValidateSlug(req.Slug) {
			return errcode.ErrPageSlugInvalid
		}
		if existing, _ := s.pageRepo.FindByVersionAndSlug(ctx, page.VersionID, req.Slug); existing != nil && existing.ID != page.ID {
			return errcode.ErrPageSlugExists
		}
	}

	update := map[string]interface{}{
		"title":   req.Title,
		"slug":    req.Slug,
		"content": req.Content,
	}
	if req.SectionID != "" {
		sectionOID, err := primitive.ObjectIDFromHex(req.SectionID)
		if err == nil {
			update["section_id"] = sectionOID
		}
	}
	return toAppError(s.pageRepo.UpdateByID(ctx, id, update))
}

// Patch 局部更新文档页（自动保存）
func (s *PageService) Patch(ctx context.Context, id primitive.ObjectID, req *dto.PatchPageReq) *errcode.AppError {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrPageNotFound
	}
	if appErr := s.versionService.CheckEditable(ctx, page.VersionID); appErr != nil {
		return errcode.ErrPageVersionArchived
	}
	version, err := s.versionRepo.FindByIDTyped(ctx, page.VersionID)
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

	update := map[string]interface{}{}
	if req.Content != "" {
		update["content"] = req.Content
	}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Slug != "" {
		update["slug"] = req.Slug
	}
	if len(update) == 0 {
		return nil
	}
	return toAppError(s.pageRepo.UpdateByID(ctx, id, update))
}

// Delete 删除文档页
func (s *PageService) Delete(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrPageNotFound
	}
	if appErr := s.versionService.CheckEditable(ctx, page.VersionID); appErr != nil {
		return errcode.ErrPageVersionArchived
	}
	return toAppError(s.pageRepo.DeleteByID(ctx, id))
}

// Publish 发布单页
func (s *PageService) Publish(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrPageNotFound
	}
	if page.Status == constants.PageStatusPublished {
		return errcode.ErrPageAlreadyPublished
	}

	// 校验所属版本必须为 published
	version, err := s.versionService.versionRepo.FindByIDTyped(ctx, page.VersionID)
	if err != nil {
		return errcode.ErrVersionNotFound
	}
	if version.Status != constants.VersionStatusPublished {
		return errcode.ErrPageVersionNotPublished
	}

	now := time.Now().UTC()
	update := map[string]interface{}{
		"status": constants.PageStatusPublished,
	}
	// 仅首次发布时设置 published_at
	if page.PublishedAt.IsZero() {
		update["published_at"] = now
	}

	logger.L().Info("文档页发布",
		zap.String("page_id", id.Hex()),
		zap.String("title", page.Title),
	)
	return toAppError(s.pageRepo.UpdateByID(ctx, id, update))
}

// Unpublish 下线单页
func (s *PageService) Unpublish(ctx context.Context, id primitive.ObjectID) *errcode.AppError {
	page, err := s.pageRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrPageNotFound
	}
	if page.Status == constants.PageStatusDraft {
		return errcode.ErrPageAlreadyDraft
	}
	return toAppError(s.pageRepo.UpdateByID(ctx, id, map[string]interface{}{
		"status": constants.PageStatusDraft,
	}))
}

// Sort 批量更新文档页排序
func (s *PageService) Sort(ctx context.Context, items []dto.SortItem) *errcode.AppError {
	for _, item := range items {
		oid, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			continue
		}
		_ = s.pageRepo.UpdateByID(ctx, oid, map[string]interface{}{"sort_order": item.SortOrder})
	}
	return nil
}

// GetVersionTree 获取版本文档树（公开接口，侧边栏数据）
func (s *PageService) GetVersionTree(ctx context.Context, versionID primitive.ObjectID) ([]dto.SectionTree, *errcode.AppError) {
	version, err := s.versionRepo.FindByIDTyped(ctx, versionID)
	if err != nil {
		return nil, errcode.ErrVersionNotFound
	}
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, version.TenantID); appErr != nil {
		return nil, appErr
	}
	sectionRepo := repository.NewSectionRepo()
	sections, err := sectionRepo.FindByVersionID(ctx, versionID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	pages, err := s.pageRepo.FindPublishedMetaByVersion(ctx, versionID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	// 按 section_id 分组
	pageMap := make(map[primitive.ObjectID][]dto.PageMeta)
	for _, p := range pages {
		pageMap[p.SectionID] = append(pageMap[p.SectionID], dto.PageMeta{
			ID:        p.ID.Hex(),
			Title:     p.Title,
			Slug:      p.Slug,
			SortOrder: p.SortOrder,
		})
	}

	var tree []dto.SectionTree
	for _, sec := range sections {
		node := dto.SectionTree{
			ID:        sec.ID.Hex(),
			Title:     sec.Title,
			SortOrder: sec.SortOrder,
			Pages:     pageMap[sec.ID],
		}
		if node.Pages == nil {
			node.Pages = []dto.PageMeta{}
		}
		tree = append(tree, node)
	}
	if tree == nil {
		tree = []dto.SectionTree{}
	}
	return tree, nil
}
