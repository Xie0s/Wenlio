// Package repository MongoDB 数据访问层
//
// 职责：封装文档主题集合的数据库操作
// 对外接口：ThemeRepo
package repository

import (
	"context"
	"time"

	"docplatform/internal/model/entity"
	mongopkg "docplatform/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ThemeRepo 主题数据访问
type ThemeRepo struct {
	*BaseRepo
}

// NewThemeRepo 创建 ThemeRepo
func NewThemeRepo() *ThemeRepo {
	return &ThemeRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("themes"))}
}

// Create 创建主题
func (r *ThemeRepo) Create(ctx context.Context, theme *entity.Theme) error {
	now := time.Now().UTC()
	theme.ID = primitive.NewObjectID()
	theme.CreatedAt = now
	theme.UpdatedAt = now
	_, err := r.InsertOne(ctx, theme)
	return err
}

// FindByIDTyped 按 ID 查询主题
func (r *ThemeRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Theme, error) {
	var theme entity.Theme
	err := r.FindByID(ctx, id, &theme)
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// FindByIDAndTenant 按 ID + tenant_id 查询主题（租户隔离）
func (r *ThemeRepo) FindByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string) (*entity.Theme, error) {
	var theme entity.Theme
	err := r.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, &theme)
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// FindByIDsMap 批量查询多个主题，返回 id→theme 映射。
func (r *ThemeRepo) FindByIDsMap(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]*entity.Theme, error) {
	if len(ids) == 0 {
		return map[primitive.ObjectID]*entity.Theme{}, nil
	}
	cursor, err := r.Collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[primitive.ObjectID]*entity.Theme, len(ids))
	for cursor.Next(ctx) {
		var t entity.Theme
		if err := cursor.Decode(&t); err == nil {
			tCopy := t
			result[t.ID] = &tCopy
		}
	}
	return result, cursor.Err()
}

// ListByTenant 按租户查询主题列表
func (r *ThemeRepo) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Theme, error) {
	filter := bson.M{"tenant_id": tenantID}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var themes []*entity.Theme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, err
	}
	return themes, nil
}

// ListMinimalByTenant 查询租户所有主题的精简信息（仅 _id/category_id/tag_ids/access_mode，用于分类/标签列表统计）
func (r *ThemeRepo) ListMinimalByTenant(ctx context.Context, tenantID string) ([]*entity.Theme, error) {
	opts := options.Find().
		SetProjection(bson.M{"_id": 1, "category_id": 1, "tag_ids": 1, "access_mode": 1}).
		SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var themes []*entity.Theme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, err
	}
	return themes, nil
}

// FindBySlug 按租户和 slug 查询
func (r *ThemeRepo) FindBySlug(ctx context.Context, tenantID, slug string) (*entity.Theme, error) {
	var theme entity.Theme
	err := r.FindOne(ctx, bson.M{"tenant_id": tenantID, "slug": slug}, &theme)
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// UpdateByIDAndTenant 按 ID + tenant_id 更新主题（租户隔离）
func (r *ThemeRepo) UpdateByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string, update bson.M) error {
	update["updated_at"] = time.Now().UTC()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, bson.M{"$set": update})
	return err
}

// DeleteByIDAndTenant 按 ID + tenant_id 删除主题（租户隔离）
func (r *ThemeRepo) DeleteByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

// CountByCategoryID 统计指定分类下的主题数量
func (r *ThemeRepo) CountByCategoryID(ctx context.Context, tenantID string, categoryID primitive.ObjectID) (int64, error) {
	return r.Count(ctx, bson.M{"tenant_id": tenantID, "category_id": categoryID})
}

// ListByFilter 支持分类/标签筛选的主题列表查询
// categoryIDs 为空时不筛选分类；传入多个 ID 时使用 $in（覆盖父分类+子分类场景）
func (r *ThemeRepo) ListByFilter(ctx context.Context, tenantID string, categoryIDs []primitive.ObjectID, tagIDs []primitive.ObjectID) ([]*entity.Theme, error) {
	filter := bson.M{"tenant_id": tenantID}
	if len(categoryIDs) == 1 {
		filter["category_id"] = categoryIDs[0]
	} else if len(categoryIDs) > 1 {
		filter["category_id"] = bson.M{"$in": categoryIDs}
	}
	if len(tagIDs) > 0 {
		filter["tag_ids"] = bson.M{"$all": tagIDs}
	}
	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var themes []*entity.Theme
	if err = cursor.All(ctx, &themes); err != nil {
		return nil, err
	}
	return themes, nil
}

// RemoveTagFromAll 从租户所有主题的 tag_ids 中移除指定标签
func (r *ThemeRepo) RemoveTagFromAll(ctx context.Context, tenantID string, tagID primitive.ObjectID) error {
	_, err := r.Collection.UpdateMany(ctx,
		bson.M{"tenant_id": tenantID, "tag_ids": tagID},
		bson.M{"$pull": bson.M{"tag_ids": tagID}},
	)
	return err
}

// EnsureIndexes 创建主题集合索引
func (r *ThemeRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "sort_order", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "category_id", Value: 1}, {Key: "sort_order", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "tag_ids", Value: 1}, {Key: "sort_order", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
