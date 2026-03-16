// Package handler HTTP 请求处理器
//
// 职责：处理媒体文件上传、列表、删除、存储审计相关的 HTTP 请求
// 对外接口：MediaHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/service"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MediaHandler 媒体文件处理器
type MediaHandler struct {
	mediaService *service.MediaService
}

// NewMediaHandler 创建 MediaHandler
func NewMediaHandler() *MediaHandler {
	return &MediaHandler{mediaService: service.NewMediaService()}
}

// Upload 上传媒体文件（图片）
func (h *MediaHandler) Upload(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	tenantID := middleware.GetTenantID(c)
	userID, err := primitive.ObjectIDFromHex(auth.UserID)
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, errcode.ErrMissingParam)
		return
	}

	media, appErr := h.mediaService.Upload(c.Request.Context(), tenantID, userID, file)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, media)
}

// List 列出租户所有媒体文件
func (h *MediaHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	items, appErr := h.mediaService.List(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, items)
}

// Delete 删除媒体文件（同步清理存储）
// 查询参数 force=true 跳过引用检查，强制删除
func (h *MediaHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	mediaID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	force := c.DefaultQuery("force", "false") == "true"
	if appErr := h.mediaService.Delete(c.Request.Context(), tenantID, mediaID, force); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Audit 审计存储：对照 DB 记录与实际存储中的文件
func (h *MediaHandler) Audit(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	storageType := c.DefaultQuery("storage_type", "cloud")
	if storageType != "cloud" && storageType != "local" {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	result, appErr := h.mediaService.AuditStorage(c.Request.Context(), tenantID, storageType)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, result)
}

// Usage 查询媒体文件使用来源（扫描所有文档页内容中的引用）
func (h *MediaHandler) Usage(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	result, appErr := h.mediaService.GetUsage(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, result)
}

// CleanupUnused 批量清理所有未被任何文档页引用的媒体文件
func (h *MediaHandler) CleanupUnused(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	result, appErr := h.mediaService.CleanupUnused(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, result)
}

// DeleteOrphan 删除存储中的孤立文件
func (h *MediaHandler) DeleteOrphan(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req struct {
		StorageType string `json:"storage_type" binding:"required"`
		Key         string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if req.StorageType != "cloud" && req.StorageType != "local" {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.mediaService.DeleteOrphan(c.Request.Context(), tenantID, req.StorageType, req.Key); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
