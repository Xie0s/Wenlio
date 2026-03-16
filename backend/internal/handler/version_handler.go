// Package handler HTTP 请求处理器
//
// 职责：处理版本管理相关的 HTTP 请求
// 对外接口：VersionHandler
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

// VersionHandler 版本处理器
type VersionHandler struct {
	versionService *service.VersionService
}

// NewVersionHandler 创建 VersionHandler
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{versionService: service.NewVersionService()}
}

// List 版本列表
func (h *VersionHandler) List(c *gin.Context) {
	themeID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	versions, appErr := h.versionService.ListByTheme(c.Request.Context(), middleware.GetTenantID(c), themeID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, versions)
}

// Create 创建版本
func (h *VersionHandler) Create(c *gin.Context) {
	themeID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.CreateVersionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	userID, _ := primitive.ObjectIDFromHex(auth.UserID)
	version, appErr := h.versionService.Create(c.Request.Context(), middleware.GetTenantID(c), themeID, userID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, version)
}

// Get 版本详情
func (h *VersionHandler) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	version, appErr := h.versionService.GetByID(c.Request.Context(), id)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, version)
}

// Update 更新版本
func (h *VersionHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateVersionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.versionService.Update(c.Request.Context(), id, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Publish 发布版本
func (h *VersionHandler) Publish(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.Publish(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Unpublish 取消发布版本（published → draft）
func (h *VersionHandler) Unpublish(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.Unpublish(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Unarchive 取消归档版本（archived → published）
func (h *VersionHandler) Unarchive(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.Unarchive(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Archive 归档版本
func (h *VersionHandler) Archive(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.Archive(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// SetDefault 设为默认版本
func (h *VersionHandler) SetDefault(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.SetDefault(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Delete 删除版本
func (h *VersionHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.versionService.Delete(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Clone 克隆版本
func (h *VersionHandler) Clone(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.CloneVersionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	userID, _ := primitive.ObjectIDFromHex(auth.UserID)
	version, appErr := h.versionService.Clone(c.Request.Context(), id, userID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, version)
}
