// Package dto 请求/响应数据传输对象
//
// 职责：定义租户管理相关的请求和响应结构体
// 对外接口：CreateTenantReq, UpdateTenantReq, TenantListItemDTO
package dto

import (
	"time"
)

// TenantListItemDTO 租户列表条目（含用户数统计及管理员信息）
type TenantListItemDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	LogoURL        string    `json:"logo_url"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserCount      int64     `json:"user_count"`
	AdminUsername  string    `json:"admin_username"`
	AdminName      string    `json:"admin_name"`
	AdminAvatarURL string    `json:"admin_avatar_url"`
}

type PublicTenantDTO struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	LogoURL        string    `json:"logo_url"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	BrowserTitle   string    `json:"browser_title"`
	BrowserIconURL string    `json:"browser_icon_url"`
}

// CreateTenantReq 创建租户请求
type CreateTenantReq struct {
	ID            string `json:"id"             binding:"required,min=3,max=32"`
	Name          string `json:"name"           binding:"required"`
	LogoURL       string `json:"logo_url"`
	AdminUsername string `json:"admin_username" binding:"required"`
	AdminPassword string `json:"admin_password" binding:"required,min=8"`
	AdminName     string `json:"admin_name"     binding:"required"`
	AdminEmail    string `json:"admin_email"`
}

// UpdateTenantReq 更新租户请求
type UpdateTenantReq struct {
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}
