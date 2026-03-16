// Package handler HTTP 请求处理器
//
// 职责：处理用户管理相关的 HTTP 请求
// 对外接口：UserHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"
	"docplatform/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建 UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{userService: service.NewUserService()}
}

// ListAll 全平台用户列表（超管）
func (h *UserHandler) ListAll(c *gin.Context) {
	p := utils.ParsePagination(c)
	users, total, appErr := h.userService.ListAll(c.Request.Context(), p.Page, p.PageSize, p.Keyword)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessWithPage(c, users, response.Pagination{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: utils.CalcTotalPages(total, p.PageSize),
	})
}

// CreateSuperAdmin 创建超管
func (h *UserHandler) CreateSuperAdmin(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	user, appErr := h.userService.CreateSuperAdmin(c.Request.Context(), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, service.UserFromEntity(user))
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	tenantID := middleware.GetTenantID(c)
	if appErr := h.userService.ResetPassword(c.Request.Context(), id, req.NewPassword, tenantID, auth.Role); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// ListByTenant 租户用户列表
func (h *UserHandler) ListByTenant(c *gin.Context) {
	p := utils.ParsePagination(c)
	users, total, appErr := h.userService.ListByTenant(c.Request.Context(), middleware.GetTenantID(c), p.Page, p.PageSize)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessWithPage(c, users, response.Pagination{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: utils.CalcTotalPages(total, p.PageSize),
	})
}

// CreateTenantAdmin 创建租户管理员
func (h *UserHandler) CreateTenantAdmin(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	user, appErr := h.userService.CreateTenantAdmin(c.Request.Context(), middleware.GetTenantID(c), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, service.UserFromEntity(user))
}

// Update 更新用户信息
func (h *UserHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	if appErr := h.userService.Update(c.Request.Context(), id, &req, middleware.GetTenantID(c), auth.Role); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Deactivate 禁用用户
func (h *UserHandler) Deactivate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	if appErr := h.userService.Deactivate(c.Request.Context(), id, auth.UserID, middleware.GetTenantID(c), auth.Role); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Activate 启用用户
func (h *UserHandler) Activate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	if appErr := h.userService.Activate(c.Request.Context(), id, middleware.GetTenantID(c), auth.Role); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	tenantID := middleware.GetTenantID(c)
	if appErr := h.userService.Delete(c.Request.Context(), id, auth.UserID, tenantID, auth.Role); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
