// Package repository MongoDB 数据访问层
//
// 职责：封装租户功能设置集合的数据库操作
// 对外接口：TenantSettingsRepo
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

// TenantSettingsRepo 租户设置数据访问
type TenantSettingsRepo struct {
	*BaseRepo
}

// NewTenantSettingsRepo 创建 TenantSettingsRepo
func NewTenantSettingsRepo() *TenantSettingsRepo {
	return &TenantSettingsRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("tenant_settings"))}
}

// FindByTenantID 按租户 ID 获取设置（不存在时返回 nil, nil）
func (r *TenantSettingsRepo) FindByTenantID(ctx context.Context, tenantID string) (*entity.TenantSettings, error) {
	var settings entity.TenantSettings
	err := r.Collection.FindOne(ctx, bson.M{"_id": tenantID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

// Upsert 原子更新设置（不存在则创建）
// updates 为 bson.M，仅更新指定字段
func (r *TenantSettingsRepo) Upsert(ctx context.Context, tenantID string, updates bson.M) error {
	now := time.Now().UTC()
	updates["updated_at"] = now

	setOnInsert := bson.M{
		"_id":        tenantID,
		"created_at": now,
	}

	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": tenantID},
		bson.M{
			"$set":         updates,
			"$setOnInsert": setOnInsert,
		},
		options.Update().SetUpsert(true),
	)
	return err
}

// GetStorageConfig 仅获取 storage 字段（减少数据传输）
func (r *TenantSettingsRepo) GetStorageConfig(ctx context.Context, tenantID string) (*entity.StorageSettings, error) {
	var result struct {
		Storage *entity.StorageSettings `bson:"storage"`
	}
	opts := options.FindOne().SetProjection(bson.M{"storage": 1})
	err := r.Collection.FindOne(ctx, bson.M{"_id": tenantID}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result.Storage, nil
}

// GetAIConfig 仅获取 ai 字段
func (r *TenantSettingsRepo) GetAIConfig(ctx context.Context, tenantID string) (*entity.AISettings, error) {
	var result struct {
		AI *entity.AISettings `bson:"ai"`
	}
	opts := options.FindOne().SetProjection(bson.M{"ai": 1})
	err := r.Collection.FindOne(ctx, bson.M{"_id": tenantID}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result.AI, nil
}

// EnsureIndexes 创建集合索引（_id 已自动索引，此处无额外需求）
func (r *TenantSettingsRepo) EnsureIndexes(ctx context.Context) error {
	return nil
}
