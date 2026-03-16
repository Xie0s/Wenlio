// Package dto 请求/响应数据传输对象
//
// 职责：定义用户管理相关的请求结构体
// 对外接口：CreateUserReq, UpdateUserReq, ResetPasswordReq
package dto

// CreateUserReq 创建用户请求（超管创建超管 / 租户管理员邀请管理员）
type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
