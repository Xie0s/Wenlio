// Package dto 请求/响应数据传输对象
//
// 职责：定义搜索功能专用的请求和响应结构体
// 对外接口：SearchReq, SearchResult, SearchResponse
package dto

// SearchReq 全文搜索请求参数
type SearchReq struct {
	Query     string `form:"q"          binding:"required"`
	TenantID  string `form:"tenant_id"  binding:"required"`
	VersionID string `form:"version_id"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

// SearchResult 单条搜索结果
type SearchResult struct {
	PageID      string `json:"page_id"`
	Title       string `json:"title"`
	Snippet     string `json:"snippet"`
	ThemeName   string `json:"theme_name"`
	VersionName string `json:"version_name"`
	Path        string `json:"path"`
}

// SearchResponse 搜索响应（含分页元信息）
type SearchResponse struct {
	Total int64          `json:"total"`
	Items []SearchResult `json:"items"`
}
