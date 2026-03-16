// Package dto 请求/响应数据传输对象
//
// 职责：定义主题管理相关的请求和响应结构体
// 对外接口：CreateThemeReq, UpdateThemeReq, SortItem
package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateThemeReq 创建主题请求
type CreateThemeReq struct {
	Name        string `json:"name"        binding:"required"`
	Slug        string `json:"slug"        binding:"required"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id" binding:"required"`
}

// UpdateThemeReq 更新主题请求
type UpdateThemeReq struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description *string  `json:"description"`
	CategoryID  string   `json:"category_id"`
	TagIDs      []string `json:"tag_ids"`
	AccessMode  *string  `json:"access_mode"` // "public" / "login" / "code"（nil 不更新）
	AccessCode  *string  `json:"access_code"` // 6位验证码（nil 不更新，仅 access_mode=code 时有意义）
}

// SortReq 批量排序请求
type SortReq struct {
	Items []SortItem `json:"items" binding:"required,dive"`
}

// SortItem 排序项
type SortItem struct {
	ID        string `json:"id"         binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// ThemeCurrentVersion 当前版本信息
type ThemeCurrentVersion struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Label     string             `json:"label"`
	Status    string             `json:"status"`
	IsDefault bool               `json:"is_default"`
}

// ThemeListItem 主题列表项（含统计信息）
type ThemeListItem struct {
	ID             primitive.ObjectID   `json:"id"`
	TenantID       string               `json:"tenant_id"`
	Name           string               `json:"name"`
	Slug           string               `json:"slug"`
	Description    string               `json:"description"`
	SortOrder      int                  `json:"sort_order"`
	CategoryID     *primitive.ObjectID  `json:"category_id,omitempty"`
	TagIDs         []primitive.ObjectID `json:"tag_ids"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	VersionCount   int64                `json:"version_count"`
	SectionCount   int64                `json:"section_count"`
	PageCount      int64                `json:"page_count"`
	CurrentVersion *ThemeCurrentVersion `json:"current_version"`
	AccessMode     string               `json:"access_mode"` // 访问模式: "public" / "login" / "code"（空串等价于 public）
}
