// Package handler HTTP 请求处理器
//
// 职责：处理认证相关的 HTTP 请求（登录、登出、获取当前用户、修改密码）
// 对外接口：AuthHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/captcha"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: service.NewAuthService()}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	resp, appErr := h.authService.Login(c.Request.Context(), buildCaptchaClientContext(c), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) CreateCaptchaChallenge(c *gin.Context) {
	var req dto.CaptchaChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	resp, appErr := h.authService.CreateCaptchaChallenge(c.Request.Context(), buildCaptchaClientContext(c), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) VerifyCaptcha(c *gin.Context) {
	var req dto.VerifyCaptchaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	resp, appErr := h.authService.VerifyCaptcha(c.Request.Context(), buildCaptchaClientContext(c), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	resp, appErr := h.authService.Register(c.Request.Context(), buildCaptchaClientContext(c), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, resp)
}

// Logout 登出（客户端清 Token 即可）
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// Me 获取当前用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	info, appErr := h.authService.GetCurrentUser(c.Request.Context(), auth.UserID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, info)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	if appErr := h.authService.ChangePassword(c.Request.Context(), auth.UserID, &req); appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, nil)
}

// UpdateProfile 更新个人资料（name / bio 可选更新）
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	var req dto.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	info, appErr := h.authService.UpdateProfile(c.Request.Context(), auth.UserID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, info)
}

// UploadAvatar 上传用户头像
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, errcode.ErrMissingParam)
		return
	}
	info, appErr := h.authService.UploadAvatar(c.Request.Context(), auth.UserID, file)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, info)
}

// UploadProfileBg 上传个人主页背景图
func (h *AuthHandler) UploadProfileBg(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, errcode.ErrMissingParam)
		return
	}
	info, appErr := h.authService.UploadProfileBg(c.Request.Context(), auth.UserID, file)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, info)
}

func buildCaptchaClientContext(c *gin.Context) captcha.ClientContext {
	return captcha.ClientContext{
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
}
