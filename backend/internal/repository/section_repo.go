// Package repository MongoDB 数据访问层
//
// 职责：封装章节集合的数据库操作
// 对外接口：SectionRepo
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

// SectionRepo 章节数据访问
type SectionRepo struct {
	*BaseRepo
}

// NewSectionRepo 创建 SectionRepo
func NewSectionRepo() *SectionRepo {
	return &SectionRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("sections"))}
}

// Create 创建章节
func (r *SectionRepo) Create(ctx context.Context, section *entity.Section) error {
	now := time.Now().UTC()
	section.ID = primitive.NewObjectID()
	section.CreatedAt = now
	section.UpdatedAt = now
	_, err := r.InsertOne(ctx, section)
	return err
}

// FindByIDTyped 按 ID 查询章节
func (r *SectionRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Section, error) {
	var section entity.Section
	err := r.FindByID(ctx, id, &section)
	if err != nil {
		return nil, err
	}
	return &section, nil
}

// ListByVersion 按版本查询章节列表（按 sort_order 排序）
func (r *SectionRepo) ListByVersion(ctx context.Context, tenantID string, versionID primitive.ObjectID) ([]*entity.Section, error) {
	filter := bson.M{"tenant_id": tenantID, "version_id": versionID}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sections []*entity.Section
	if err = cursor.All(ctx, &sections); err != nil {
		return nil, err
	}
	return sections, nil
}

// FindByVersionID 按版本 ID 查询所有章节（无租户过滤，用于克隆）
func (r *SectionRepo) FindByVersionID(ctx context.Context, versionID primitive.ObjectID) ([]*entity.Section, error) {
	filter := bson.M{"version_id": versionID}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sections []*entity.Section
	if err = cursor.All(ctx, &sections); err != nil {
		return nil, err
	}
	return sections, nil
}

// MaxSortOrder 获取版本下最大排序值
func (r *SectionRepo) MaxSortOrder(ctx context.Context, versionID primitive.ObjectID) (int, error) {
	opts := options.FindOne().SetSort(bson.M{"sort_order": -1}).SetProjection(bson.M{"sort_order": 1})
	var result struct {
		SortOrder int `bson:"sort_order"`
	}
	err := r.Collection.FindOne(ctx, bson.M{"version_id": versionID}, opts).Decode(&result)
	if err != nil {
		return 0, nil
	}
	return result.SortOrder, nil
}

// DeleteByVersion 删除版本下所有章节
func (r *SectionRepo) DeleteByVersion(ctx context.Context, versionID primitive.ObjectID) error {
	_, err := r.DeleteMany(ctx, bson.M{"version_id": versionID})
	return err
}

// CountByVersionIDs 批量统计版本下章节数量
func (r *SectionRepo) CountByVersionIDs(ctx context.Context, tenantID string, versionIDs []primitive.ObjectID) (map[primitive.ObjectID]int64, error) {
	if len(versionIDs) == 0 {
		return map[primitive.ObjectID]int64{}, nil
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"tenant_id":  tenantID,
			"version_id": bson.M{"$in": versionIDs},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$version_id",
			"count": bson.M{"$sum": 1},
		}}},
	}
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[primitive.ObjectID]int64, len(versionIDs))
	for cursor.Next(ctx) {
		var item struct {
			ID    primitive.ObjectID `bson:"_id"`
			Count int64              `bson:"count"`
		}
		if decodeErr := cursor.Decode(&item); decodeErr != nil {
			continue
		}
		counts[item.ID] = item.Count
	}
	return counts, nil
}

// EnsureIndexes 创建章节集合索引
func (r *SectionRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "version_id", Value: 1}, {Key: "sort_order", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
