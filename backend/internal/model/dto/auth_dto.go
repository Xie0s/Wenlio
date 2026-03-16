// Package dto 请求/响应数据传输对象
//
// 职责：定义认证相关的请求和响应结构体
// 对外接口：LoginReq, LoginResp, ChangePasswordReq, UserInfo
package dto

// LoginReq 登录请求
type LoginReq struct {
	Username     string `json:"username"      binding:"required"`
	Password     string `json:"password"      binding:"required"`
	CaptchaToken string `json:"captcha_token" binding:"required"`
}

type RegisterReq struct {
	TenantID      string `json:"tenant_id"       binding:"required,min=3,max=32"`
	TenantName    string `json:"tenant_name"     binding:"required"`
	LogoURL       string `json:"logo_url"`
	AdminUsername string `json:"admin_username"  binding:"required"`
	AdminPassword string `json:"admin_password"  binding:"required,min=8"`
	AdminName     string `json:"admin_name"      binding:"required"`
	AdminEmail    string `json:"admin_email"`
	CaptchaToken  string `json:"captcha_token"   binding:"required"`
}

type CaptchaChallengeReq struct {
	Scene string `json:"scene" binding:"required,oneof=login register reset_password"`
}

type CaptchaChallengeResp struct {
	ChallengeID   string `json:"challenge_id"`
	Scene         string `json:"scene"`
	Mode          string `json:"mode"`
	Prompt        string `json:"prompt"`
	ExpiresAt     int64  `json:"expires_at"`
	MinDecisionMs int64  `json:"min_decision_ms"`
}

type CaptchaSignalSummary struct {
	DwellMs             int64  `json:"dwell_ms" binding:"min=0"`
	VisibleMs           int64  `json:"visible_ms" binding:"min=0"`
	FocusedMs           int64  `json:"focused_ms" binding:"min=0"`
	VisibilityChanges   int    `json:"visibility_changes" binding:"min=0"`
	FocusChanges        int    `json:"focus_changes" binding:"min=0"`
	PointerEvents       int    `json:"pointer_events" binding:"min=0"`
	KeyEvents           int    `json:"key_events" binding:"min=0"`
	TrustedClick        bool   `json:"trusted_click"`
	Language            string `json:"language"`
	Platform            string `json:"platform"`
	ScreenWidth         int    `json:"screen_width" binding:"min=0"`
	ScreenHeight        int    `json:"screen_height" binding:"min=0"`
	TimezoneOffset      int    `json:"timezone_offset"`
	TouchPoints         int    `json:"touch_points" binding:"min=0"`
	HardwareConcurrency int    `json:"hardware_concurrency" binding:"min=0"`
	Webdriver           bool   `json:"webdriver"`
}

type VerifyCaptchaReq struct {
	Scene       string               `json:"scene"        binding:"required,oneof=login register reset_password"`
	ChallengeID string               `json:"challenge_id" binding:"required"`
	DurationMs  int64                `json:"duration_ms"  binding:"required,min=0"`
	Signals     CaptchaSignalSummary `json:"signals"      binding:"required"`
}

type VerifyCaptchaResp struct {
	CaptchaToken string `json:"captcha_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken string   `json:"access_token"`
	ExpiresIn   int      `json:"expires_in"`
	User        UserInfo `json:"user"`
}

// UserInfo 用户基础信息（用于 Token 响应和 /auth/me）
type UserInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	TenantID     string `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	AvatarURL    string `json:"avatar_url"`
	Bio          string `json:"bio"`
	ProfileBgURL string `json:"profile_bg_url"`
}

// UpdateProfileReq 更新个人资料请求（单字段可选更新）
type UpdateProfileReq struct {
	Name *string `json:"name"`
	Bio  *string `json:"bio"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
