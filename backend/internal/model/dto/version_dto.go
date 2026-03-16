// Package dto 请求/响应数据传输对象
//
// 职责：定义版本管理相关的请求结构体
// 对外接口：CreateVersionReq, UpdateVersionReq, CloneVersionReq
package dto

// CreateVersionReq 创建版本请求
type CreateVersionReq struct {
	Name  string `json:"name"  binding:"required"`
	Label string `json:"label"`
}

// UpdateVersionReq 更新版本请求
type UpdateVersionReq struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// CloneVersionReq 克隆版本请求
type CloneVersionReq struct {
	Name  string `json:"name"  binding:"required"`
	Label string `json:"label"`
}
