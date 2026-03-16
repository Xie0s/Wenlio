// Package handler HTTP 请求处理器
//
// 职责：处理文档页和章节管理相关的 HTTP 请求
// 对外接口：PageHandler, SectionHandler
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

// SectionHandler 章节处理器
type SectionHandler struct {
	sectionService *service.SectionService
}

// NewSectionHandler 创建 SectionHandler
func NewSectionHandler() *SectionHandler {
	return &SectionHandler{sectionService: service.NewSectionService()}
}

// ListSections 章节列表
func (h *SectionHandler) ListSections(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	sections, appErr := h.sectionService.ListByVersion(c.Request.Context(), middleware.GetTenantID(c), versionID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, sections)
}

// CreateSection 创建章节
func (h *SectionHandler) CreateSection(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.CreateSectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	section, appErr := h.sectionService.Create(c.Request.Context(), middleware.GetTenantID(c), versionID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, section)
}

// UpdateSection 更新章节
func (h *SectionHandler) UpdateSection(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateSectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.sectionService.Update(c.Request.Context(), id, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// DeleteSection 删除章节
func (h *SectionHandler) DeleteSection(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.sectionService.Delete(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// SortSections 批量排序章节
func (h *SectionHandler) SortSections(c *gin.Context) {
	var req dto.SortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.sectionService.Sort(c.Request.Context(), req.Items); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// PageHandler 文档页处理器
type PageHandler struct {
	pageService *service.PageService
}

// NewPageHandler 创建 PageHandler
func NewPageHandler() *PageHandler {
	return &PageHandler{pageService: service.NewPageService()}
}

// ListPages 章节下文档页列表
func (h *PageHandler) ListPages(c *gin.Context) {
	sectionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	pages, appErr := h.pageService.ListBySection(c.Request.Context(), middleware.GetTenantID(c), sectionID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, pages)
}

// CreatePage 创建文档页
func (h *PageHandler) CreatePage(c *gin.Context) {
	sectionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.CreatePageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	userID, _ := primitive.ObjectIDFromHex(auth.UserID)

	// 从章节获取版本 ID
	sectionSvc := service.NewSectionService()
	section, appErr := sectionSvc.GetByID(c.Request.Context(), sectionID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}

	page, appErr := h.pageService.Create(c.Request.Context(), middleware.GetTenantID(c), sectionID, section.VersionID, userID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, page)
}

// GetPage 文档页详情
func (h *PageHandler) GetPage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	page, appErr := h.pageService.GetByID(c.Request.Context(), id)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, page)
}

// UpdatePage 全量更新文档页
func (h *PageHandler) UpdatePage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdatePageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.pageService.Update(c.Request.Context(), id, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// PatchPage 局部更新文档页（自动保存）
func (h *PageHandler) PatchPage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.PatchPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.pageService.Patch(c.Request.Context(), id, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// DeletePage 删除文档页
func (h *PageHandler) DeletePage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.pageService.Delete(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// PublishPage 发布文档页
func (h *PageHandler) PublishPage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.pageService.Publish(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// UnpublishPage 下线文档页
func (h *PageHandler) UnpublishPage(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.pageService.Unpublish(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// SortPages 批量排序文档页
func (h *PageHandler) SortPages(c *gin.Context) {
	var req dto.SortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.pageService.Sort(c.Request.Context(), req.Items); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
