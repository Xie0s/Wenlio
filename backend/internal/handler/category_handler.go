// Package handler HTTP 请求处理器
//
// 职责：处理分类管理相关的 HTTP 请求
// 对外接口：CategoryHandler
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

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 创建 CategoryHandler
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// List 获取分类树
func (h *CategoryHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	tree, appErr := h.categoryService.List(c.Request.Context(), tenantID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tree)
}

// Create 创建分类
func (h *CategoryHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	category, appErr := h.categoryService.Create(c.Request.Context(), tenantID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, category)
}

// Update 更新分类
func (h *CategoryHandler) Update(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.categoryService.Update(c.Request.Context(), id, tenantID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "更新成功"})
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.categoryService.Delete(c.Request.Context(), id, tenantID); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// Sort 批量排序
func (h *CategoryHandler) Sort(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	var req dto.CategorySortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.categoryService.Sort(c.Request.Context(), tenantID, req.Items); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"message": "排序已保存"})
}
