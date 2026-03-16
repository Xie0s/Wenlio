// Package handler HTTP 请求处理器
//
// 职责：处理租户功能设置相关的 HTTP 请求（云存储、AI 等）
// 对外接口：TenantSettingsHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// TenantSettingsHandler 租户设置处理器
type TenantSettingsHandler struct {
	settingsService *service.TenantSettingsService
}

// NewTenantSettingsHandler 创建 TenantSettingsHandler
func NewTenantSettingsHandler() *TenantSettingsHandler {
	return &TenantSettingsHandler{
		settingsService: service.NewTenantSettingsService(),
	}
}

// Get 获取完整设置（脱敏）
func (h *TenantSettingsHandler) Get(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	resp, appErr := h.settingsService.GetSettings(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, resp)
}

// UpdateStorage 更新云存储设置
func (h *TenantSettingsHandler) UpdateStorage(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.UpdateStorageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.settingsService.UpdateStorage(c.Request.Context(), tenantID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "云存储设置已保存"})
}

// TestS3 测试 S3 连通性
func (h *TenantSettingsHandler) TestS3(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if appErr := h.settingsService.TestS3Connection(c.Request.Context(), tenantID); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "云存储连接成功"})
}

// GetStorageUsage 获取本地存储用量
func (h *TenantSettingsHandler) GetStorageUsage(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	usage, appErr := h.settingsService.GetStorageUsage(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, usage)
}

// UpdateAI 更新 AI 设置
func (h *TenantSettingsHandler) UpdateAI(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.UpdateAIReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.settingsService.UpdateAI(c.Request.Context(), tenantID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "AI 设置已保存"})
}

// UpdateAccess 更新站点级访问控制设置（维护模式 / 画廊登录可见）
func (h *TenantSettingsHandler) UpdateAccess(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.UpdateAccessReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.settingsService.UpdateAccess(c.Request.Context(), tenantID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "访问控制设置已保存"})
}

// TestAI 测试 AI 连通性（预留）
func (h *TenantSettingsHandler) TestAI(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if appErr := h.settingsService.TestAIConnection(c.Request.Context(), tenantID); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "AI 连接成功"})
}
