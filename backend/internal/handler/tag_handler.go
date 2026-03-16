// Package handler HTTP 请求处理器
//
// 职责：处理标签管理相关的 HTTP 请求
// 对外接口：TagHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler 创建 TagHandler
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: service.NewTagService(),
	}
}

// List 获取标签列表
func (h *TagHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	tags, appErr := h.tagService.List(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tags)
}

// Create 创建标签
func (h *TagHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	tag, appErr := h.tagService.Create(c.Request.Context(), tenantID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tag)
}

// Update 更新标签
func (h *TagHandler) Update(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.tagService.Update(c.Request.Context(), id, tenantID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "更新成功"})
}

// Delete 删除标签（?force=true 时跳过使用检查，直接级联清理）
func (h *TagHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	force := c.Query("force") == "true"
	if appErr := h.tagService.Delete(c.Request.Context(), id, tenantID, force); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
