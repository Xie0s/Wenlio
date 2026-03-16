// Package entity MongoDB 文档映射结构体
//
// 职责：定义用户集合的文档结构
// 对外接口：User
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型，对应 MongoDB users 集合
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID       string             `bson:"tenant_id"      json:"tenant_id"`
	Username       string             `bson:"username"       json:"username"`
	Password       string             `bson:"password"       json:"-"`
	Name           string             `bson:"name"           json:"name"`
	Email          string             `bson:"email"          json:"email"`
	Role           string             `bson:"role"           json:"role"`
	Status         string             `bson:"status"         json:"status"`
	AvatarURL      string             `bson:"avatar_url"     json:"avatar_url"`
	Bio            string             `bson:"bio"            json:"bio"`
	ProfileBgURL   string             `bson:"profile_bg_url" json:"profile_bg_url"`
	LastLoginAt    time.Time          `bson:"last_login_at"  json:"last_login_at"`
	LoginFailCount int                `bson:"login_fail_count" json:"-"`
	LockedUntil    time.Time          `bson:"locked_until"   json:"-"`
	CreatedAt      time.Time          `bson:"created_at"     json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"     json:"updated_at"`
}
