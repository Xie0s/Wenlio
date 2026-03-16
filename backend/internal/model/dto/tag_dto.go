// Package dto 请求/响应数据传输对象
//
// 职责：定义标签管理相关的请求和响应结构体
// 对外接口：CreateTagReq, UpdateTagReq
package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTagReq 创建标签请求
type CreateTagReq struct {
	Name  string `json:"name"  binding:"required"`
	Slug  string `json:"slug"  binding:"required"`
	Color string `json:"color"`
}

// UpdateTagReq 更新标签请求
type UpdateTagReq struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}

// TagListItem 标签列表项
type TagListItem struct {
	ID         primitive.ObjectID `json:"id"`
	TenantID   string             `json:"tenant_id"`
	Name       string             `json:"name"`
	Slug       string             `json:"slug"`
	Color      string             `json:"color"`
	UsageCount int64              `json:"usage_count"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
