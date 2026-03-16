// Package handler HTTP 请求处理器
//
// 职责：处理租户首页个性化相关的管理端 HTTP 请求
// 对外接口：HomepageHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// HomepageHandler 首页个性化处理器（租户管理端）
type HomepageHandler struct {
	tenantService *service.TenantService
}

// NewHomepageHandler 创建 HomepageHandler
func NewHomepageHandler() *HomepageHandler {
	return &HomepageHandler{tenantService: service.NewTenantService()}
}

// GetDraft 获取首页草稿 + 已发布配置
func (h *HomepageHandler) GetDraft(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	hp, appErr := h.tenantService.GetHomepageDraft(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	if hp == nil {
		response.Success(c, &dto.HomepageDraftResp{})
		return
	}
	response.Success(c, &dto.HomepageDraftResp{
		Published: hp.Published,
		Draft:     hp.Draft,
		UpdatedAt: hp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// SaveDraft 保存首页草稿
func (h *HomepageHandler) SaveDraft(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.SaveHomepageDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.tenantService.SaveHomepageDraft(c.Request.Context(), tenantID, req.ToLayout()); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Publish 发布首页（draft → published）
func (h *HomepageHandler) Publish(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	if appErr := h.tenantService.PublishHomepage(c.Request.Context(), tenantID); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
