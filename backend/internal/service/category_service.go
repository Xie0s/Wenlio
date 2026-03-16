// Package service 业务逻辑层
//
// 职责：实现分类管理业务逻辑（CRUD、树构建、排序）
// 对外接口：CategoryService
package service

import (
	"context"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CategoryService 分类业务
type CategoryService struct {
	categoryRepo *repository.CategoryRepo
	themeRepo    *repository.ThemeRepo
	versionRepo  *repository.VersionRepo
	tenantRepo   *repository.TenantRepo
}

// NewCategoryService 创建 CategoryService
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepo(),
		themeRepo:    repository.NewThemeRepo(),
		versionRepo:  repository.NewVersionRepo(),
		tenantRepo:   repository.NewTenantRepo(),
	}
}

// List 获取租户分类树（管理端，含所有主题计数，包括草稿）
func (s *CategoryService) List(ctx context.Context, tenantID string) ([]*dto.CategoryTreeNode, *errcode.AppError) {
	categories, err := s.categoryRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 统计每个分类下的主题总数（含草稿）
	themes, err := s.themeRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	countByCategoryID := make(map[primitive.ObjectID]int64)
	for _, t := range themes {
		countByCategoryID[t.CategoryID]++
	}
	return buildTree(categories, primitive.NilObjectID, countByCategoryID), nil
}

// ListPublic 获取读者端分类树（含各分类直接挂载的已发布主题数）
// isAuthenticated: 未登录时排除 access_mode="login" 的主题计数
func (s *CategoryService) ListPublic(ctx context.Context, tenantID string, isAuthenticated bool) ([]*dto.CategoryTreeNode, *errcode.AppError) {
	categories, err := s.categoryRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 精简投影：只需 _id 和 category_id
	themes, err := s.themeRepo.ListMinimalByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 聚合查询：只返回有已发布版本的主题 ID 集合（无需加载完整版本文档）
	themeIDs := make([]primitive.ObjectID, 0, len(themes))
	for _, t := range themes {
		themeIDs = append(themeIDs, t.ID)
	}
	publishedIDs, err := s.versionRepo.ListPublishedThemeIDs(ctx, tenantID, themeIDs)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	publishedSet := make(map[primitive.ObjectID]struct{}, len(publishedIDs))
	for _, id := range publishedIDs {
		publishedSet[id] = struct{}{}
	}
	// 按 category_id 统计已发布主题数（未登录时排除 login 类型主题）
	countByCategoryID := make(map[primitive.ObjectID]int64)
	for _, t := range themes {
		if _, ok := publishedSet[t.ID]; ok {
			if !isAuthenticated && t.AccessMode == "login" {
				continue
			}
			countByCategoryID[t.CategoryID]++
		}
	}
	return buildTree(categories, primitive.NilObjectID, countByCategoryID), nil
}

// Create 创建分类
func (s *CategoryService) Create(ctx context.Context, tenantID string, req *dto.CreateCategoryReq) (*entity.Category, *errcode.AppError) {
	if !utils.ValidateSlug(req.Slug) {
		return nil, errcode.ErrPageSlugInvalid
	}
	if existing, _ := s.categoryRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil {
		return nil, errcode.ErrCategorySlugExists
	}

	var parentID primitive.ObjectID
	var level int
	if req.ParentID != "" {
		oid, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errcode.ErrInvalidParam
		}
		parent, err := s.categoryRepo.FindByIDTyped(ctx, oid)
		if err != nil {
			return nil, errcode.ErrCategoryNotFound
		}
		if parent.TenantID != tenantID {
			return nil, errcode.ErrCategoryNotFound
		}
		if parent.Level >= 1 {
			return nil, errcode.ErrCategoryMaxDepth
		}
		parentID = oid
		level = parent.Level + 1
	}

	category := &entity.Category{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    parentID,
		SortOrder:   req.SortOrder,
		Level:       level,
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return category, nil
}

// Update 更新分类
func (s *CategoryService) Update(ctx context.Context, id primitive.ObjectID, tenantID string, req *dto.UpdateCategoryReq) *errcode.AppError {
	category, err := s.categoryRepo.FindByIDTyped(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errcode.ErrCategoryNotFound
		}
		return errcode.ErrDatabase.Wrap(err)
	}
	if category.TenantID != tenantID {
		return errcode.ErrCategoryNotFound
	}

	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Slug != "" && req.Slug != category.Slug {
		if !utils.ValidateSlug(req.Slug) {
			return errcode.ErrPageSlugInvalid
		}
		if existing, _ := s.categoryRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil && existing.ID != category.ID {
			return errcode.ErrCategorySlugExists
		}
		update["slug"] = req.Slug
	}
	if req.SortOrder != nil {
		update["sort_order"] = *req.SortOrder
	}
	if len(update) == 0 {
		return nil
	}
	return toAppError(s.categoryRepo.UpdateByIDAndTenant(ctx, id, tenantID, update))
}

// Delete 删除分类（需无子分类且无挂载主题）
func (s *CategoryService) Delete(ctx context.Context, id primitive.ObjectID, tenantID string) *errcode.AppError {
	category, err := s.categoryRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrCategoryNotFound
	}
	if category.TenantID != tenantID {
		return errcode.ErrCategoryNotFound
	}
	childCount, err := s.categoryRepo.CountChildren(ctx, tenantID, id)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if childCount > 0 {
		return errcode.ErrCategoryHasChildren
	}
	themeCount, err := s.themeRepo.CountByCategoryID(ctx, tenantID, id)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if themeCount > 0 {
		return errcode.ErrCategoryHasThemes
	}
	return toAppError(s.categoryRepo.DeleteByIDAndTenant(ctx, id, tenantID))
}

// Sort 批量更新分类排序
func (s *CategoryService) Sort(ctx context.Context, tenantID string, items []dto.SortItem) *errcode.AppError {
	for _, item := range items {
		oid, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			continue
		}
		_ = s.categoryRepo.UpdateByIDAndTenant(ctx, oid, tenantID, map[string]interface{}{"sort_order": item.SortOrder})
	}
	return nil
}

// ============================================================
// 内部辅助：构建分类树
// ============================================================

// buildTree 将平面列表（已按 level+sort_order 排序）构建为树结构
// countByCategoryID 为 nil 时不填充 ThemeCount（管理端调用路径）
func buildTree(categories []*entity.Category, parentID primitive.ObjectID, countByCategoryID map[primitive.ObjectID]int64) []*dto.CategoryTreeNode {
	var nodes []*dto.CategoryTreeNode
	for _, c := range categories {
		if c.ParentID == parentID {
			node := &dto.CategoryTreeNode{
				ID:        c.ID,
				TenantID:  c.TenantID,
				Name:      c.Name,
				Slug:      c.Slug,
				ParentID:  c.ParentID,
				SortOrder: c.SortOrder,
				Level:     c.Level,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			}
			if countByCategoryID != nil {
				node.ThemeCount = countByCategoryID[c.ID]
			}
			node.Children = buildTree(categories, c.ID, countByCategoryID)
			if node.Children == nil {
				node.Children = []*dto.CategoryTreeNode{}
			}
			nodes = append(nodes, node)
		}
	}
	if nodes == nil {
		return []*dto.CategoryTreeNode{}
	}
	return nodes
}
