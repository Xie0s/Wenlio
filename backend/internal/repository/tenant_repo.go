// Package repository MongoDB 数据访问层
//
// 职责：封装租户集合的数据库操作
// 对外接口：TenantRepo
package repository

import (
	"context"
	"time"

	"docplatform/internal/model/entity"
	mongopkg "docplatform/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TenantRepo 租户数据访问
type TenantRepo struct {
	*BaseRepo
}

// NewTenantRepo 创建 TenantRepo
func NewTenantRepo() *TenantRepo {
	return &TenantRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("tenants"))}
}

// FindByIDTyped 按字符串 ID 查询租户
func (r *TenantRepo) FindByIDTyped(ctx context.Context, id string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.FindByStringID(ctx, id, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Create 创建租户
func (r *TenantRepo) Create(ctx context.Context, tenant *entity.Tenant) error {
	now := time.Now().UTC()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now
	_, err := r.InsertOne(ctx, tenant)
	return err
}

// List 分页查询租户列表
func (r *TenantRepo) List(ctx context.Context, page, pageSize int, keyword string) ([]*entity.Tenant, int64, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"_id": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize)).SetSort(bson.M{"created_at": -1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tenants []*entity.Tenant
	if err = cursor.All(ctx, &tenants); err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

// ListDeleting 查询所有处于删除中的租户，用于服务启动恢复异步删除任务
func (r *TenantRepo) ListDeleting(ctx context.Context) ([]*entity.Tenant, error) {
	filter := bson.M{"status": "deleting"}
	opts := options.Find().SetSort(bson.M{"updated_at": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tenants []*entity.Tenant
	if err = cursor.All(ctx, &tenants); err != nil {
		return nil, err
	}
	return tenants, nil
}

// Exists 检查租户 ID 是否已存在
func (r *TenantRepo) Exists(ctx context.Context, id string) bool {
	count, _ := r.Count(ctx, bson.M{"_id": id})
	return count > 0
}

// MarkDeletingIfAllowed 将租户原子标记为 deleting，仅当当前状态不是 deleting 时生效
func (r *TenantRepo) MarkDeletingIfAllowed(ctx context.Context, id string) (bool, error) {
	now := time.Now().UTC()
	result, err := r.Collection.UpdateOne(ctx, bson.M{
		"_id":    id,
		"status": bson.M{"$ne": "deleting"},
	}, bson.M{"$set": bson.M{
		"status":     "deleting",
		"updated_at": now,
	}})
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

// FindByIDWithoutHomepage 查询租户（排除 homepage 字段，减少数据传输）
func (r *TenantRepo) FindByIDWithoutHomepage(ctx context.Context, id string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	opts := options.FindOne().SetProjection(bson.M{"homepage": 0})
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) FindPublicByID(ctx context.Context, id string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	opts := options.FindOne().SetProjection(bson.M{
		"name":       1,
		"logo_url":   1,
		"status":     1,
		"created_at": 1,
		"updated_at": 1,
		"homepage.published.global.browser_title":    1,
		"homepage.published.global.browser_icon_url": 1,
	})
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetHomepage 获取租户首页配置（仅 homepage 字段）
func (r *TenantRepo) GetHomepage(ctx context.Context, id string) (*entity.TenantHomepage, error) {
	var result struct {
		Homepage *entity.TenantHomepage `bson:"homepage"`
	}
	opts := options.FindOne().SetProjection(bson.M{"homepage": 1})
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Homepage, nil
}

// SaveHomepageDraft 保存首页草稿
func (r *TenantRepo) SaveHomepageDraft(ctx context.Context, id string, layout *entity.HomepageLayout) error {
	now := time.Now().UTC()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"homepage.draft":      layout,
			"homepage.updated_at": now,
			"updated_at":          now,
		},
	})
	return err
}

// PublishHomepage 发布首页（将 draft 复制到 published）
func (r *TenantRepo) PublishHomepage(ctx context.Context, id string, layout *entity.HomepageLayout) error {
	now := time.Now().UTC()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"homepage.published":  layout,
			"homepage.updated_at": now,
			"updated_at":          now,
		},
	})
	return err
}

// EnsureIndexes 创建租户集合索引
func (r *TenantRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
