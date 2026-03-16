// Package repository MongoDB 数据访问层
//
// 职责：封装评论集合的数据库操作
// 对外接口：CommentRepo
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

// CommentRepo 评论数据访问
type CommentRepo struct {
	*BaseRepo
}

// NewCommentRepo 创建 CommentRepo
func NewCommentRepo() *CommentRepo {
	return &CommentRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("comments"))}
}

// Create 创建评论
func (r *CommentRepo) Create(ctx context.Context, comment *entity.Comment) error {
	now := time.Now().UTC()
	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	_, err := r.InsertOne(ctx, comment)
	return err
}

// FindByIDTyped 按 ID 查询评论
func (r *CommentRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.FindByID(ctx, id, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByPage 按文档页查询已审核评论（公开接口）
func (r *CommentRepo) ListByPage(ctx context.Context, pageID primitive.ObjectID) ([]*entity.Comment, error) {
	filter := bson.M{"page_id": pageID, "status": "approved"}
	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*entity.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

// ListByTenant 按租户查询评论列表（管理后台，支持状态筛选）
func (r *CommentRepo) ListByTenant(ctx context.Context, tenantID, status string, page, pageSize int) ([]*entity.Comment, int64, error) {
	filter := bson.M{"tenant_id": tenantID}
	if status != "" {
		filter["status"] = status
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

	var comments []*entity.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// DeleteWithChildren 删除评论及其子回复
func (r *CommentRepo) DeleteWithChildren(ctx context.Context, commentID primitive.ObjectID) error {
	_, _ = r.DeleteMany(ctx, bson.M{"parent_id": commentID})
	return r.DeleteByID(ctx, commentID)
}

// EnsureIndexes 创建评论集合索引
func (r *CommentRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "page_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
