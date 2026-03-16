// Package handler HTTP 请求处理器
//
// 职责：处理评论管理相关的 HTTP 请求
// 对外接口：CommentHandler
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

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler 创建 CommentHandler
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{commentService: service.NewCommentService()}
}

// List 评论列表（管理后台）
func (h *CommentHandler) List(c *gin.Context) {
	p := utils.ParsePagination(c)
	status := c.Query("status")
	comments, total, appErr := h.commentService.ListByTenant(c.Request.Context(), middleware.GetTenantID(c), status, p.Page, p.PageSize)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessWithPage(c, comments, response.Pagination{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: utils.CalcTotalPages(total, p.PageSize),
	})
}

// Approve 批准评论
func (h *CommentHandler) Approve(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.commentService.Approve(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Reject 拒绝评论
func (h *CommentHandler) Reject(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.commentService.Reject(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	if appErr := h.commentService.Delete(c.Request.Context(), id); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// Reply 管理员回复评论
func (h *CommentHandler) Reply(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.ReplyCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	auth := middleware.GetAuthContext(c)
	authSvc := service.NewAuthService()
	userInfo, _ := authSvc.GetCurrentUser(c.Request.Context(), auth.UserID)
	adminName := "管理员"
	if userInfo != nil {
		adminName = userInfo.Name
	}

	reply, appErr := h.commentService.Reply(c.Request.Context(), id, middleware.GetTenantID(c), adminName, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, reply)
}
