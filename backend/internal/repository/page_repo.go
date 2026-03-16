// Package repository MongoDB 数据访问层
//
// 职责：封装文档页集合的数据库操作，包含全文搜索
// 对外接口：PageRepo
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

// PageRepo 文档页数据访问
type PageRepo struct {
	*BaseRepo
}

// NewPageRepo 创建 PageRepo
func NewPageRepo() *PageRepo {
	return &PageRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("pages"))}
}

// Create 创建文档页
func (r *PageRepo) Create(ctx context.Context, page *entity.Page) error {
	now := time.Now().UTC()
	page.ID = primitive.NewObjectID()
	page.CreatedAt = now
	page.UpdatedAt = now
	_, err := r.InsertOne(ctx, page)
	return err
}

// FindByIDTyped 按 ID 查询文档页
func (r *PageRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Page, error) {
	var page entity.Page
	err := r.FindByID(ctx, id, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// FindByVersionAndSlug 按版本和 slug 查询
func (r *PageRepo) FindByVersionAndSlug(ctx context.Context, versionID primitive.ObjectID, slug string) (*entity.Page, error) {
	var page entity.Page
	err := r.FindOne(ctx, bson.M{"version_id": versionID, "slug": slug}, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// ListBySection 按章节查询文档页列表
func (r *PageRepo) ListBySection(ctx context.Context, tenantID string, sectionID primitive.ObjectID) ([]*entity.Page, error) {
	filter := bson.M{"tenant_id": tenantID, "section_id": sectionID}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pages []*entity.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

// FindBySectionID 按章节 ID 查询所有文档页（用于克隆）
func (r *PageRepo) FindBySectionID(ctx context.Context, sectionID primitive.ObjectID) ([]*entity.Page, error) {
	filter := bson.M{"section_id": sectionID}
	opts := options.Find().SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pages []*entity.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

// FindPublishedMetaByVersion 查询版本下已发布文档页元数据（不含 content）
func (r *PageRepo) FindPublishedMetaByVersion(ctx context.Context, versionID primitive.ObjectID) ([]*entity.Page, error) {
	filter := bson.M{"version_id": versionID, "status": "published"}
	projection := bson.M{"content": 0}
	opts := options.Find().SetProjection(projection).SetSort(bson.M{"sort_order": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pages []*entity.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

// BatchPublish 批量发布版本下所有 draft 文档页
func (r *PageRepo) BatchPublish(ctx context.Context, versionID primitive.ObjectID) (int64, error) {
	now := time.Now().UTC()
	filter := bson.M{"version_id": versionID, "status": "draft"}
	update := bson.M{"status": "published", "published_at": now, "updated_at": now}
	return r.UpdateMany(ctx, filter, update)
}

// BatchUnpublish 批量将版本下所有 published 文档页回退为 draft
func (r *PageRepo) BatchUnpublish(ctx context.Context, versionID primitive.ObjectID) (int64, error) {
	now := time.Now().UTC()
	filter := bson.M{"version_id": versionID, "status": "published"}
	update := bson.M{"status": "draft", "updated_at": now}
	return r.UpdateMany(ctx, filter, update)
}

// MaxSortOrder 获取章节下最大排序值
func (r *PageRepo) MaxSortOrder(ctx context.Context, sectionID primitive.ObjectID) (int, error) {
	opts := options.FindOne().SetSort(bson.M{"sort_order": -1}).SetProjection(bson.M{"sort_order": 1})
	var result struct {
		SortOrder int `bson:"sort_order"`
	}
	err := r.Collection.FindOne(ctx, bson.M{"section_id": sectionID}, opts).Decode(&result)
	if err != nil {
		return 0, nil
	}
	return result.SortOrder, nil
}

// DeleteBySection 删除章节下所有文档页
func (r *PageRepo) DeleteBySection(ctx context.Context, sectionID primitive.ObjectID) error {
	_, err := r.DeleteMany(ctx, bson.M{"section_id": sectionID})
	return err
}

// DeleteByVersion 删除版本下所有文档页
func (r *PageRepo) DeleteByVersion(ctx context.Context, versionID primitive.ObjectID) error {
	_, err := r.DeleteMany(ctx, bson.M{"version_id": versionID})
	return err
}

// CountByVersionIDs 批量统计版本下文档页数量
func (r *PageRepo) CountByVersionIDs(ctx context.Context, tenantID string, versionIDs []primitive.ObjectID) (map[primitive.ObjectID]int64, error) {
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

// SearchByQuery 基于正则表达式的全文搜索，支持中文。
// 使用聚合管道：$match 筛选 → $addFields 计算相关性分 → $facet 同时返回总数与分页数据。
// 标题命中权重(2) > 仅正文命中权重(1)，相同分值按 sort_order 升序排列。
func (r *PageRepo) SearchByQuery(ctx context.Context, tenantID, query string, versionID primitive.ObjectID, page, pageSize int) ([]*entity.Page, int64, error) {
	escaped := primitive.Regex{Pattern: escapeRegex(query), Options: "i"}

	matchStage := bson.D{
		{Key: "tenant_id", Value: tenantID},
		{Key: "status", Value: "published"},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: escaped}},
			bson.D{{Key: "content", Value: escaped}},
		}},
	}
	if !versionID.IsZero() {
		matchStage = append(matchStage, bson.E{Key: "version_id", Value: versionID})
	}

	skip := int64((page - 1) * pageSize)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "_score", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$regexMatch", Value: bson.D{
						{Key: "input", Value: "$title"},
						{Key: "regex", Value: escapeRegex(query)},
						{Key: "options", Value: "i"},
					}}},
					2,
					1,
				}},
			}},
		}}},
		{{Key: "$facet", Value: bson.D{
			{Key: "count", Value: bson.A{
				bson.D{{Key: "$count", Value: "n"}},
			}},
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{
					{Key: "_score", Value: -1},
					{Key: "sort_order", Value: 1},
				}}},
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: int64(pageSize)}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "title", Value: 1},
					{Key: "slug", Value: 1},
					{Key: "content", Value: 1},
					{Key: "version_id", Value: 1},
				}}},
			}},
		}}},
	}

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var facetResult []struct {
		Count []struct {
			N int64 `bson:"n"`
		} `bson:"count"`
		Data []*entity.Page `bson:"data"`
	}
	if err = cursor.All(ctx, &facetResult); err != nil {
		return nil, 0, err
	}
	if len(facetResult) == 0 {
		return []*entity.Page{}, 0, nil
	}

	var total int64
	if len(facetResult[0].Count) > 0 {
		total = facetResult[0].Count[0].N
	}
	return facetResult[0].Data, total, nil
}

// escapeRegex 转义正则特殊字符，防止查询字符串被解释为正则元字符
func escapeRegex(s string) string {
	const special = `\.+*?()|[]{}^$`
	result := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		for j := 0; j < len(special); j++ {
			if c == special[j] {
				result = append(result, '\\')
				break
			}
		}
		result = append(result, c)
	}
	return string(result)
}

// TextSearch 全文搜索
func (r *PageRepo) TextSearch(ctx context.Context, filter bson.D, page, pageSize int) ([]*entity.Page, int64, error) {
	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	projection := bson.D{
		{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}},
	}
	sort := bson.D{
		{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}},
	}
	opts := options.Find().SetProjection(projection).SetSort(sort).SetSkip(skip).SetLimit(int64(pageSize))

	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var pages []*entity.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, 0, err
	}
	return pages, total, nil
}

// ListMinimalByTenant 查询租户下所有文档页的精简信息（用于媒体引用扫描）
// 只投影 id/version_id/title/slug/content，避免传输多余字段
func (r *PageRepo) ListMinimalByTenant(ctx context.Context, tenantID string) ([]*entity.Page, error) {
	opts := options.Find().SetProjection(bson.M{
		"_id":        1,
		"version_id": 1,
		"title":      1,
		"slug":       1,
		"content":    1,
	})
	cursor, err := r.Find(ctx, bson.M{"tenant_id": tenantID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var pages []*entity.Page
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

// EnsureIndexes 创建文档页集合索引
func (r *PageRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "version_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "version_id", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "section_id", Value: 1}, {Key: "sort_order", Value: 1}}},
		{
			Keys: bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
			Options: options.Index().SetWeights(bson.D{
				{Key: "title", Value: 10},
				{Key: "content", Value: 1},
			}),
		},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
