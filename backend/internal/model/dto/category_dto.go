// Package dto 请求/响应数据传输对象
//
// 职责：定义分类管理相关的请求和响应结构体
// 对外接口：CreateCategoryReq, UpdateCategoryReq, CategoryTreeNode
package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCategoryReq 创建分类请求
type CreateCategoryReq struct {
	Name      string `json:"name"       binding:"required"`
	Slug      string `json:"slug"       binding:"required"`
	ParentID  string `json:"parent_id"` // 空字符串表示根节点
	SortOrder int    `json:"sort_order"`
}

// UpdateCategoryReq 更新分类请求
type UpdateCategoryReq struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	SortOrder *int   `json:"sort_order"`
}

// CategoryTreeNode 分类树节点（含子节点，用于前端渲染树结构）
type CategoryTreeNode struct {
	ID         primitive.ObjectID  `json:"id"`
	TenantID   string              `json:"tenant_id"`
	Name       string              `json:"name"`
	Slug       string              `json:"slug"`
	ParentID   primitive.ObjectID  `json:"parent_id"`
	SortOrder  int                 `json:"sort_order"`
	Level      int                 `json:"level"`
	ThemeCount int64               `json:"theme_count"` // 直接挂载在该分类下的已发布主题数
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Children   []*CategoryTreeNode `json:"children"` // 递归树结构
}

// CategorySortReq 批量排序请求
type CategorySortReq struct {
	Items []SortItem `json:"items" binding:"required,dive"`
}
