// Package repository MongoDB 数据访问层
//
// 职责：封装媒体文件记录集合的数据库操作
// 对外接口：MediaRepo
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

// MediaRepo 媒体文件数据访问
type MediaRepo struct {
	*BaseRepo
}

// NewMediaRepo 创建 MediaRepo
func NewMediaRepo() *MediaRepo {
	return &MediaRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("media"))}
}

// Create 创建媒体记录
func (r *MediaRepo) Create(ctx context.Context, media *entity.Media) error {
	now := time.Now().UTC()
	media.ID = primitive.NewObjectID()
	media.CreatedAt = now
	media.UpdatedAt = now
	_, err := r.InsertOne(ctx, media)
	return err
}

// SumFileSizeByTenant 统计指定租户+存储类型的文件总大小（字节）
// storageType 为空时统计所有类型
func (r *MediaRepo) SumFileSizeByTenant(ctx context.Context, tenantID string, storageType string) (int64, error) {
	filter := bson.M{"tenant_id": tenantID}
	if storageType != "" {
		filter["storage_type"] = storageType
	}
	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$file_size"}}},
	}
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	var result []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}

// FindByIDAndTenant 按 ID + 租户 ID 查询单条媒体记录（租户隔离）
func (r *MediaRepo) FindByIDAndTenant(ctx context.Context, id primitive.ObjectID, tenantID string) (*entity.Media, error) {
	var media entity.Media
	err := r.FindOne(ctx, bson.M{"_id": id, "tenant_id": tenantID}, &media)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

// ListByTenant 查询租户所有媒体记录，按创建时间倒序
func (r *MediaRepo) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Media, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*entity.Media
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindLocalByTenant 查询租户本地存储的所有媒体记录（含 storage_key 用于删除）
func (r *MediaRepo) FindLocalByTenant(ctx context.Context, tenantID string) ([]*entity.Media, error) {
	filter := bson.M{"tenant_id": tenantID, "storage_type": "local"}
	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*entity.Media
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// FindCloudByTenant 查询租户云存储的所有媒体记录（含 storage_key 用于 S3 删除）
func (r *MediaRepo) FindCloudByTenant(ctx context.Context, tenantID string) ([]*entity.Media, error) {
	filter := bson.M{"tenant_id": tenantID, "storage_type": "cloud"}
	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*entity.Media
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// EnsureIndexes 创建媒体集合索引
func (r *MediaRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "storage_type", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
