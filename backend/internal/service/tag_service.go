// Package service 业务逻辑层
//
// 职责：实现标签管理业务逻辑（CRUD）
// 对外接口：TagService
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

// TagService 标签业务
type TagService struct {
	tagRepo     *repository.TagRepo
	themeRepo   *repository.ThemeRepo
	versionRepo *repository.VersionRepo
	tenantRepo  *repository.TenantRepo
}

// NewTagService 创建 TagService
func NewTagService() *TagService {
	return &TagService{
		tagRepo:     repository.NewTagRepo(),
		themeRepo:   repository.NewThemeRepo(),
		versionRepo: repository.NewVersionRepo(),
		tenantRepo:  repository.NewTenantRepo(),
	}
}

// List 获取租户标签列表（管理端，实时统计各标签被主题使用的次数）
func (s *TagService) List(ctx context.Context, tenantID string) ([]*dto.TagListItem, *errcode.AppError) {
	tags, err := s.tagRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	themes, err := s.themeRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	countByTagID := make(map[primitive.ObjectID]int64)
	for _, t := range themes {
		for _, tagID := range t.TagIDs {
			countByTagID[tagID]++
		}
	}
	items := make([]*dto.TagListItem, 0, len(tags))
	for _, t := range tags {
		items = append(items, &dto.TagListItem{
			ID:         t.ID,
			TenantID:   t.TenantID,
			Name:       t.Name,
			Slug:       t.Slug,
			Color:      t.Color,
			UsageCount: countByTagID[t.ID],
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		})
	}
	return items, nil
}

// ListPublic 获取读者端标签列表（实时统计已发布主题数，过滤计数为 0 的标签）
// isAuthenticated: 未登录时排除 access_mode="login" 的主题计数
func (s *TagService) ListPublic(ctx context.Context, tenantID string, isAuthenticated bool) ([]*dto.TagListItem, *errcode.AppError) {
	tags, err := s.tagRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 精简投影：只需 _id 和 tag_ids
	themes, err := s.themeRepo.ListMinimalByTenant(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	// 聚合查询：只返回有已发布版本的主题 ID 集合
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
	// 按 tag_id 统计已发布主题数（未登录时排除 login 类型主题）
	countByTagID := make(map[primitive.ObjectID]int64)
	for _, t := range themes {
		if _, ok := publishedSet[t.ID]; ok {
			if !isAuthenticated && t.AccessMode == "login" {
				continue
			}
			for _, tagID := range t.TagIDs {
				countByTagID[tagID]++
			}
		}
	}
	items := make([]*dto.TagListItem, 0, len(tags))
	for _, t := range tags {
		count := countByTagID[t.ID]
		if count == 0 {
			continue
		}
		items = append(items, &dto.TagListItem{
			ID:         t.ID,
			TenantID:   t.TenantID,
			Name:       t.Name,
			Slug:       t.Slug,
			Color:      t.Color,
			UsageCount: count,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		})
	}
	return items, nil
}

// Create 创建标签
func (s *TagService) Create(ctx context.Context, tenantID string, req *dto.CreateTagReq) (*entity.Tag, *errcode.AppError) {
	if !utils.ValidateSlug(req.Slug) {
		return nil, errcode.ErrPageSlugInvalid
	}
	if existing, _ := s.tagRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil {
		return nil, errcode.ErrTagSlugExists
	}
	tag := &entity.Tag{
		TenantModel: entity.TenantModel{TenantID: tenantID},
		Name:        req.Name,
		Slug:        req.Slug,
		Color:       req.Color,
	}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return tag, nil
}

// Update 更新标签
func (s *TagService) Update(ctx context.Context, id primitive.ObjectID, tenantID string, req *dto.UpdateTagReq) *errcode.AppError {
	tag, err := s.tagRepo.FindByIDTyped(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errcode.ErrTagNotFound
		}
		return errcode.ErrDatabase.Wrap(err)
	}
	if tag.TenantID != tenantID {
		return errcode.ErrTagNotFound
	}

	update := map[string]interface{}{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Slug != "" && req.Slug != tag.Slug {
		if !utils.ValidateSlug(req.Slug) {
			return errcode.ErrPageSlugInvalid
		}
		if existing, _ := s.tagRepo.FindBySlug(ctx, tenantID, req.Slug); existing != nil && existing.ID != tag.ID {
			return errcode.ErrTagSlugExists
		}
		update["slug"] = req.Slug
	}
	if req.Color != "" {
		update["color"] = req.Color
	}
	if len(update) == 0 {
		return nil
	}
	return toAppError(s.tagRepo.UpdateByIDAndTenant(ctx, id, tenantID, update))
}

// Delete 删除标签，并级联清理所有主题 tag_ids 中的该标签引用。
// force=false 时若标签仍被主题使用则返回 ErrTagInUse，由前端二步确认后以 force=true 重试。
func (s *TagService) Delete(ctx context.Context, id primitive.ObjectID, tenantID string, force bool) *errcode.AppError {
	tag, err := s.tagRepo.FindByIDTyped(ctx, id)
	if err != nil {
		return errcode.ErrTagNotFound
	}
	if tag.TenantID != tenantID {
		return errcode.ErrTagNotFound
	}
	if !force {
		themes, err := s.themeRepo.ListByFilter(ctx, tenantID, nil, []primitive.ObjectID{id})
		if err != nil {
			return errcode.ErrDatabase.Wrap(err)
		}
		if len(themes) > 0 {
			return errcode.ErrTagInUse
		}
	}
	if appErr := toAppError(s.tagRepo.DeleteByIDAndTenant(ctx, id, tenantID)); appErr != nil {
		return appErr
	}
	_ = s.themeRepo.RemoveTagFromAll(ctx, tenantID, id)
	return nil
}
