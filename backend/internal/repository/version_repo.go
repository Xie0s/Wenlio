// Package repository MongoDB 数据访问层
//
// 职责：封装版本集合的数据库操作
// 对外接口：VersionRepo
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

// VersionRepo 版本数据访问
type VersionRepo struct {
	*BaseRepo
}

// NewVersionRepo 创建 VersionRepo
func NewVersionRepo() *VersionRepo {
	return &VersionRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("versions"))}
}

// Create 创建版本
func (r *VersionRepo) Create(ctx context.Context, version *entity.Version) error {
	now := time.Now().UTC()
	version.ID = primitive.NewObjectID()
	version.CreatedAt = now
	version.UpdatedAt = now
	_, err := r.InsertOne(ctx, version)
	return err
}

// FindByIDTyped 按 ID 查询版本
func (r *VersionRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Version, error) {
	var version entity.Version
	err := r.FindByID(ctx, id, &version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// FindByIDsMap 批量查询多个版本，返回 id→version 映射。
// 单次 $in 查询代替 N 次单文档查询，减少 DB 往返次数。
func (r *VersionRepo) FindByIDsMap(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]*entity.Version, error) {
	if len(ids) == 0 {
		return map[primitive.ObjectID]*entity.Version{}, nil
	}
	cursor, err := r.Collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[primitive.ObjectID]*entity.Version, len(ids))
	for cursor.Next(ctx) {
		var v entity.Version
		if err := cursor.Decode(&v); err == nil {
			vCopy := v
			result[v.ID] = &vCopy
		}
	}
	return result, cursor.Err()
}

// ListByTheme 按主题查询版本列表
func (r *VersionRepo) ListByTheme(ctx context.Context, tenantID string, themeID primitive.ObjectID) ([]*entity.Version, error) {
	filter := bson.M{"tenant_id": tenantID, "theme_id": themeID}
	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []*entity.Version
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// ListByThemeIDs 按主题 ID 列表批量查询版本（用于主题列表统计，只投影必要字段）
func (r *VersionRepo) ListByThemeIDs(ctx context.Context, tenantID string, themeIDs []primitive.ObjectID) ([]*entity.Version, error) {
	if len(themeIDs) == 0 {
		return []*entity.Version{}, nil
	}
	filter := bson.M{
		"tenant_id": tenantID,
		"theme_id":  bson.M{"$in": themeIDs},
	}
	opts := options.Find().
		SetSort(bson.D{
			{Key: "theme_id", Value: 1},
			{Key: "is_default", Value: -1},
			{Key: "created_at", Value: -1},
		}).
		SetProjection(bson.M{
			"_id": 1, "theme_id": 1, "name": 1,
			"label": 1, "status": 1, "is_default": 1,
		})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []*entity.Version
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// ListPublishedByThemeIDs 批量查询多个主题下的已发布版本（利用 tenant_id+theme_id+status 复合索引，只投影必要字段）
func (r *VersionRepo) ListPublishedByThemeIDs(ctx context.Context, tenantID string, themeIDs []primitive.ObjectID) ([]*entity.Version, error) {
	if len(themeIDs) == 0 {
		return []*entity.Version{}, nil
	}
	filter := bson.M{
		"tenant_id": tenantID,
		"theme_id":  bson.M{"$in": themeIDs},
		"status":    "published",
	}
	opts := options.Find().
		SetSort(bson.D{
			{Key: "theme_id", Value: 1},
			{Key: "is_default", Value: -1},
			{Key: "created_at", Value: -1},
		}).
		SetProjection(bson.M{
			"_id": 1, "theme_id": 1, "name": 1,
			"label": 1, "status": 1, "is_default": 1,
		})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var versions []*entity.Version
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// ListPublishedThemeIDs 查询租户下所有有已发布版本的主题 ID（精简投影，用于分类/标签列表统计）
func (r *VersionRepo) ListPublishedThemeIDs(ctx context.Context, tenantID string, themeIDs []primitive.ObjectID) ([]primitive.ObjectID, error) {
	if len(themeIDs) == 0 {
		return []primitive.ObjectID{}, nil
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"tenant_id": tenantID,
			"theme_id":  bson.M{"$in": themeIDs},
			"status":    "published",
		}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$theme_id"}}},
	}
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
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

// ListPublishedByTheme 查询主题下已发布版本（公开接口，归档和草稿版本对读者不可见）
func (r *VersionRepo) ListPublishedByTheme(ctx context.Context, themeID primitive.ObjectID) ([]*entity.Version, error) {
	filter := bson.M{
		"theme_id": themeID,
		"status":   "published",
	}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []*entity.Version
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

// FindDefaultByTheme 按主题查询默认版本
func (r *VersionRepo) FindDefaultByTheme(ctx context.Context, themeID primitive.ObjectID) (*entity.Version, error) {
	var version entity.Version
	err := r.FindOne(ctx, bson.M{"theme_id": themeID, "status": "published", "is_default": true}, &version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// FindByThemeAndName 按主题 ID + 版本标签名查询（用于 Raw Markdown slug 链式定位）
func (r *VersionRepo) FindByThemeAndName(ctx context.Context, themeID primitive.ObjectID, name string) (*entity.Version, error) {
	var version entity.Version
	// 优先匹配 name 字段（URL 友好），其次匹配 label
	err := r.FindOne(ctx, bson.M{
		"theme_id": themeID,
		"status":   "published",
		"$or": bson.A{
			bson.M{"name": name},
			bson.M{"label": name},
		},
	}, &version)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// ClearDefault 清除同主题下所有版本的默认标记
func (r *VersionRepo) ClearDefault(ctx context.Context, tenantID string, themeID primitive.ObjectID) error {
	filter := bson.M{"tenant_id": tenantID, "theme_id": themeID, "is_default": true}
	_, err := r.UpdateMany(ctx, filter, bson.M{"is_default": false, "updated_at": time.Now().UTC()})
	return err
}

// CountByTheme 统计主题下版本数量
func (r *VersionRepo) CountByTheme(ctx context.Context, tenantID string, themeID primitive.ObjectID) (int64, error) {
	return r.Count(ctx, bson.M{"tenant_id": tenantID, "theme_id": themeID})
}

// DeleteByTheme 删除主题下所有版本
func (r *VersionRepo) DeleteByTheme(ctx context.Context, tenantID string, themeID primitive.ObjectID) error {
	_, err := r.DeleteMany(ctx, bson.M{"tenant_id": tenantID, "theme_id": themeID})
	return err
}

// EnsureIndexes 创建版本集合索引
func (r *VersionRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "theme_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "theme_id", Value: 1}, {Key: "is_default", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
