// Package repository MongoDB 数据访问层
//
// 职责：封装标签集合的数据库操作
// 对外接口：TagRepo
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

// TagRepo 标签数据访问
type TagRepo struct {
	*BaseRepo
}

// NewTagRepo 创建 TagRepo
func NewTagRepo() *TagRepo {
	return &TagRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("tags"))}
}

// Create 创建标签
func (r *TagRepo) Create(ctx context.Context, tag *entity.Tag) error {
	now := time.Now().UTC()
	tag.ID = primitive.NewObjectID()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	_, err := r.InsertOne(ctx, tag)
	return err
}

// FindByIDTyped 按 ObjectID 查询标签
func (r *TagRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.FindByID(ctx, id, &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindBySlug 按租户和 slug 查询
func (r *TagRepo) FindBySlug(ctx context.Context, tenantID, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.FindOne(ctx, bson.M{"tenant_id": tenantID, "slug": slug}, &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// ListByTenant 查询租户所有标签（按 usage_count desc + name 排序）
func (r *TagRepo) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Tag, error) {
	opts := options.Find().SetSort(bson.D{
		{Key: "usage_count", Value: -1},
		{Key: "name", Value: 1},
	})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*entity.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// FindBySlugs 批量按 slug 查询标签（单次 $in 查询替代 N 次 FindBySlug）
func (r *TagRepo) FindBySlugs(ctx context.Context, tenantID string, slugs []string) ([]*entity.Tag, error) {
	if len(slugs) == 0 {
		return []*entity.Tag{}, nil
	}
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID, "slug": bson.M{"$in": slugs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*entity.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByIDs 批量按 ID 查询（用于主题关联标签名称展示）
func (r *TagRepo) FindByIDs(ctx context.Context, tenantID string, ids []primitive.ObjectID) ([]*entity.Tag, error) {
	if len(ids) == 0 {
		return []*entity.Tag{}, nil
	}
	opts := options.Find().SetProjection(bson.M{"name": 1, "slug": 1, "color": 1})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID, "_id": bson.M{"$in": ids}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tags []*entity.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateByIDAndTenant 按 ID 和 tenantID 更新
func (r *TagRepo) UpdateByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string, update bson.M) error {
	update["updated_at"] = time.Now().UTC()
	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "tenant_id": tenantID},
		bson.M{"$set": update},
	)
	return err
}

// IncrUsageCount 标签使用计数 +delta（delta 可为负数）
func (r *TagRepo) IncrUsageCount(ctx context.Context, id primitive.ObjectID, delta int64) error {
	_, err := r.Collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"usage_count": delta}},
	)
	return err
}

// DeleteByIDAndTenant 按 ID 和 tenantID 删除
func (r *TagRepo) DeleteByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

// EnsureIndexes 创建标签集合索引
func (r *TagRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "slug", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "usage_count", Value: -1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
