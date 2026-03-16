// Package repository MongoDB 数据访问层
//
// 职责：封装用户集合的数据库操作
// 对外接口：UserRepo
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

// UserRepo 用户数据访问
type UserRepo struct {
	*BaseRepo
}

// NewUserRepo 创建 UserRepo
func NewUserRepo() *UserRepo {
	return &UserRepo{BaseRepo: NewBaseRepo(mongopkg.Collection("users"))}
}

// FindByUsername 按用户名查询
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.FindOne(ctx, bson.M{"username": username}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDTyped 按 ObjectID 查询并返回 User
func (r *UserRepo) FindByIDTyped(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	var user entity.User
	err := r.FindByID(ctx, id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	now := time.Now().UTC()
	user.ID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.InsertOne(ctx, user)
	return err
}

// ListByTenant 按租户查询用户列表
func (r *UserRepo) ListByTenant(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.User, int64, error) {
	filter := bson.M{"tenant_id": tenantID}
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

	var users []*entity.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// ListAll 查询所有用户（超管用）
func (r *UserRepo) ListAll(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, error) {
	filter := bson.M{}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": keyword, "$options": "i"}},
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
			{"tenant_id": bson.M{"$regex": keyword, "$options": "i"}},
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

	var users []*entity.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// UpdateLoginFail 更新登录失败信息
func (r *UserRepo) UpdateLoginFail(ctx context.Context, userID primitive.ObjectID, failCount int, lockedUntil time.Time) error {
	return r.UpdateByID(ctx, userID, bson.M{
		"login_fail_count": failCount,
		"locked_until":     lockedUntil,
	})
}

// ResetLoginFail 重置登录失败计数
func (r *UserRepo) ResetLoginFail(ctx context.Context, userID primitive.ObjectID) error {
	return r.UpdateByID(ctx, userID, bson.M{
		"login_fail_count": 0,
		"locked_until":     time.Time{},
		"last_login_at":    time.Now().UTC(),
	})
}

// FindFirstAdminByTenantIDs 批量查询各租户的首个管理员（按创建时间最早，一次聚合）
// 返回 map[tenant_id]{username, name, avatar_url}
func (r *UserRepo) FindFirstAdminByTenantIDs(ctx context.Context, tenantIDs []string) (map[string][3]string, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"tenant_id": bson.M{"$in": tenantIDs},
			"role":      "tenant_admin",
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: 1}}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":        "$tenant_id",
			"username":   bson.M{"$first": "$username"},
			"name":       bson.M{"$first": "$name"},
			"avatar_url": bson.M{"$first": "$avatar_url"},
		}}},
	}
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string][3]string, len(tenantIDs))
	for cursor.Next(ctx) {
		var item struct {
			ID        string `bson:"_id"`
			Username  string `bson:"username"`
			Name      string `bson:"name"`
			AvatarURL string `bson:"avatar_url"`
		}
		if decodeErr := cursor.Decode(&item); decodeErr != nil {
			continue
		}
		result[item.ID] = [3]string{item.Username, item.Name, item.AvatarURL}
	}
	return result, nil
}

// CountByTenantIDs 批量统计各租户的用户数（聚合管道，一次查询）
func (r *UserRepo) CountByTenantIDs(ctx context.Context, tenantIDs []string) (map[string]int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"tenant_id": bson.M{"$in": tenantIDs}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$tenant_id", "count": bson.M{"$sum": 1}}}},
	}
	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	counts := make(map[string]int64, len(tenantIDs))
	for cursor.Next(ctx) {
		var item struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if decodeErr := cursor.Decode(&item); decodeErr != nil {
			continue
		}
		counts[item.ID] = item.Count
	}
	return counts, nil
}

// EnsureIndexes 创建用户集合索引
func (r *UserRepo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tenant_id", Value: 1}, {Key: "role", Value: 1}}},
	}
	_, err := r.Collection.Indexes().CreateMany(ctx, indexes)
	return err
}
