// Package repository MongoDB 数据访问层
//
// 职责：封装 MongoDB 基础 CRUD 操作，供各业务 Repository 复用
// 对外接口：BaseRepo
package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseRepo 基础数据访问封装
type BaseRepo struct {
	Collection *mongo.Collection
}

// NewBaseRepo 创建 BaseRepo
func NewBaseRepo(coll *mongo.Collection) *BaseRepo {
	return &BaseRepo{Collection: coll}
}

// InsertOne 插入单条文档，自动设置 created_at/updated_at
func (r *BaseRepo) InsertOne(ctx context.Context, doc interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, doc)
}

// FindByID 按 ObjectID 查询单条
func (r *BaseRepo) FindByID(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	return r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(result)
}

// FindByStringID 按字符串主键查询（用于 tenants 集合）
func (r *BaseRepo) FindByStringID(ctx context.Context, id string, result interface{}) error {
	return r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(result)
}

// FindOne 按条件查询单条
func (r *BaseRepo) FindOne(ctx context.Context, filter bson.M, result interface{}) error {
	return r.Collection.FindOne(ctx, filter).Decode(result)
}

// Find 按条件查询多条
func (r *BaseRepo) Find(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return r.Collection.Find(ctx, filter, opts...)
}

// UpdateByID 按 ObjectID 更新
func (r *BaseRepo) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now().UTC()
	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// UpdateByStringID 按字符串主键更新（用于 tenants 集合）
func (r *BaseRepo) UpdateByStringID(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now().UTC()
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

// UpdateMany 批量更新
func (r *BaseRepo) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error) {
	result, err := r.Collection.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// DeleteByID 按 ObjectID 物理删除
func (r *BaseRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteByStringID 按字符串主键物理删除
func (r *BaseRepo) DeleteByStringID(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteMany 批量物理删除
func (r *BaseRepo) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// Count 按条件计数
func (r *BaseRepo) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.Collection.CountDocuments(ctx, filter)
}
