// Package handler HTTP 请求处理器
//
// 职责：处理主题管理相关的 HTTP 请求
// 对外接口：ThemeHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ThemeHandler 主题处理器
type ThemeHandler struct {
	themeService *service.ThemeService
}

// NewThemeHandler 创建 ThemeHandler
func NewThemeHandler() *ThemeHandler {
	return &ThemeHandler{themeService: service.NewThemeService()}
}

// List 主题列表，支持 ?category_id=xxx 和 ?tag_id=xxx 筛选（tag_id 可多次传）
func (h *ThemeHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	categoryID := c.Query("category_id")
	tagIDs := c.QueryArray("tag_id")
	themes, appErr := h.themeService.List(c.Request.Context(), tenantID, categoryID, tagIDs)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, themes)
}

// Create 创建主题
func (h *ThemeHandler) Create(c *gin.Context) {
	var req dto.CreateThemeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	userID, _ := primitive.ObjectIDFromHex(auth.UserID)
	theme, appErr := h.themeService.Create(c.Request.Context(), middleware.GetTenantID(c), userID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, theme)
}

// Get 主题详情
func (h *ThemeHandler) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	theme, appErr := h.themeService.GetByID(c.Request.Context(), id, middleware.GetTenantID(c))
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, theme)
}

// Update 更新主题
func (h *ThemeHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.UpdateThemeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.themeService.Update(c.Request.Context(), id, middleware.GetTenantID(c), &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Delete 删除主题
func (h *ThemeHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	tenantID := middleware.GetTenantID(c)
	cascade := false
	if rawCascade := c.Query("cascade"); rawCascade != "" {
		parsed, parseErr := strconv.ParseBool(rawCascade)
		if parseErr != nil {
			response.Fail(c, errcode.ErrInvalidParam)
			return
		}
		cascade = parsed
	}
	var appErr *errcode.AppError
	if cascade {
		appErr = h.themeService.DeleteCascade(c.Request.Context(), id, tenantID)
	} else {
		appErr = h.themeService.Delete(c.Request.Context(), id, tenantID)
	}
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Sort 批量排序
func (h *ThemeHandler) Sort(c *gin.Context) {
	var req dto.SortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.themeService.Sort(c.Request.Context(), middleware.GetTenantID(c), req.Items); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}
