// Package repository MongoDB 数据访问层
//
// 职责：封装分类集合的数据库操作
// 对外接口：CategoryRepo
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

// CategoryRepo 分类数据访问
type CategoryRepo struct {
	*BaseRepo
}

// NewCategoryRepo 创建 CategoryRepo
func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("categories"))}
}

// Create 创建分类
func (r *CategoryRepo) Create(ctx context.Context, category *entity.Category) error {
	now := time.Now().UTC()
	category.ID = primitive.NewObjectID()
	category.CreatedAt = now
	category.UpdatedAt = now
	_, err := r.InsertOne(ctx, category)
	return err
}

// FindByIDTyped 按 ObjectID 查询分类
func (r *CategoryRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Category, error) {
	var category entity.Category
	err := r.FindByID(ctx, id, &category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindBySlug 按租户和 slug 查询
func (r *CategoryRepo) FindBySlug(ctx context.Context, tenantID, slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.FindOne(ctx, bson.M{"tenant_id": tenantID, "slug": slug}, &category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ListByTenant 查询租户所有分类（按 level+sort_order 排序，用于构建树）
func (r *CategoryRepo) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Category, error) {
	opts := options.Find().SetSort(bson.D{
		{Key: "level", Value: 1},
		{Key: "sort_order", Value: 1},
	})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var categories []*entity.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// CountChildren 计算某分类的直接子分类数量
func (r *CategoryRepo) CountChildren(ctx context.Context, tenantID string, parentID primitive.ObjectID) (int64, error) {
	return r.Count(ctx, bson.M{"tenant_id": tenantID, "parent_id": parentID})
}

// ListChildIDs 查询某分类的直接子分类 ID 列表
func (r *CategoryRepo) ListChildIDs(ctx context.Context, tenantID string, parentID primitive.ObjectID) ([]primitive.ObjectID, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID, "parent_id": parentID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var ids []primitive.ObjectID
	for cursor.Next(ctx) {
		var doc struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&doc); err == nil {
			ids = append(ids, doc.ID)
		}
	}
	return ids, cursor.Err()
}

// UpdateByIDAndTenant 按 ID 和 tenantID 更新
func (r *CategoryRepo) UpdateByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string, update bson.M) error {
	update["updated_at"] = time.Now().UTC()
	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "tenant_id": tenantID},
		bson.M{"$set": update},
	)
	return err
}

// DeleteByIDAndTenant 按 ID 和 tenantID 删除
func (r *CategoryRepo) DeleteByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id, "tenant_id": tenantID})
	return err
}

// EnsureIndexes 创建分类集合索引
func (r *CategoryRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(false)},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "parent_id", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "level", Value: 1}, {Key: "sort_order", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
