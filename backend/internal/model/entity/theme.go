// Package entity MongoDB 文档映射结构体
//
// 职责：定义文档主题集合的文档结构
// 对外接口：Theme
package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

// Theme 文档主题模型，对应 MongoDB themes 集合
type Theme struct {
	TenantModel `bson:",inline"`
	Name        string               `bson:"name"         json:"name"`
	Slug        string               `bson:"slug"         json:"slug"`
	Description string               `bson:"description"  json:"description"`
	Deleting    bool                 `bson:"deleting"     json:"deleting"`
	SortOrder   int                  `bson:"sort_order"   json:"sort_order"`
	CreatedBy   primitive.ObjectID   `bson:"created_by"   json:"created_by"`
	CategoryID  primitive.ObjectID   `bson:"category_id"  json:"category_id"` // 所属分类（零值表示未分类）
	TagIDs      []primitive.ObjectID `bson:"tag_ids"      json:"tag_ids"`     // 关联标签列表
	AccessMode  string               `bson:"access_mode"  json:"access_mode"` // 访问模式: "" / "public" / "login" / "code"（空串等价于 public）
	AccessCode  string               `bson:"access_code"  json:"-"`           // 6位验证码（access_mode=code 时生效），敏感字段不回传前端
}
