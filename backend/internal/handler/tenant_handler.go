// Package handler HTTP 请求处理器
//
// 职责：处理租户管理相关的 HTTP 请求
// 对外接口：TenantHandler
package handler

import (
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/response"
	"docplatform/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TenantHandler 租户处理器
type TenantHandler struct {
	tenantService *service.TenantService
}

// NewTenantHandler 创建 TenantHandler
func NewTenantHandler() *TenantHandler {
	return &TenantHandler{tenantService: service.NewTenantService()}
}

// List 租户列表
func (h *TenantHandler) List(c *gin.Context) {
	p := utils.ParsePagination(c)
	tenants, total, appErr := h.tenantService.List(c.Request.Context(), p.Page, p.PageSize, p.Keyword)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessWithPage(c, tenants, response.Pagination{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: utils.CalcTotalPages(total, p.PageSize),
	})
}

// Create 创建租户
func (h *TenantHandler) Create(c *gin.Context) {
	var req dto.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	tenant, appErr := h.tenantService.Create(c.Request.Context(), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, tenant)
}

// Get 租户详情
func (h *TenantHandler) Get(c *gin.Context) {
	id := c.Param("id")
	tenant, appErr := h.tenantService.GetByID(c.Request.Context(), id)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tenant)
}

// Update 更新租户
func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.tenantService.Update(c.Request.Context(), id, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Suspend 封禁租户
func (h *TenantHandler) Suspend(c *gin.Context) {
	if appErr := h.tenantService.Suspend(c.Request.Context(), c.Param("id")); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Activate 解封租户
func (h *TenantHandler) Activate(c *gin.Context) {
	if appErr := h.tenantService.Activate(c.Request.Context(), c.Param("id")); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Delete 删除租户
func (h *TenantHandler) Delete(c *gin.Context) {
	if appErr := h.tenantService.Delete(c.Request.Context(), c.Param("id")); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
