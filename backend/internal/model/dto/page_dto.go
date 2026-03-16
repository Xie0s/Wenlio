// Package dto 请求/响应数据传输对象
//
// 职责：定义文档页管理相关的请求和响应结构体
// 对外接口：CreatePageReq, UpdatePageReq, PatchPageReq, SectionTree, PageMeta
package dto

// CreatePageReq 创建文档页请求
type CreatePageReq struct {
	Title   string `json:"title"   binding:"required"`
	Slug    string `json:"slug"    binding:"required"`
	Content string `json:"content"`
}

// UpdatePageReq 更新文档页请求（全量）
type UpdatePageReq struct {
	Title     string `json:"title"      binding:"required"`
	Slug      string `json:"slug"       binding:"required"`
	Content   string `json:"content"`
	SectionID string `json:"section_id"`
}

// PatchPageReq 局部更新文档页请求（自动保存 / 设置面板）
type PatchPageReq struct {
	Content string `json:"content"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
}

// CreateSectionReq 创建章节请求
type CreateSectionReq struct {
	Title string `json:"title" binding:"required"`
}

// UpdateSectionReq 更新章节请求
type UpdateSectionReq struct {
	Title string `json:"title" binding:"required"`
}

// SectionTree 版本文档树中的章节节点（公开接口响应）
type SectionTree struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	SortOrder int        `json:"sort_order"`
	Pages     []PageMeta `json:"pages"`
}

// PageMeta 文档页元数据（不含 content，用于侧边栏）
type PageMeta struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	SortOrder int    `json:"sort_order"`
}
