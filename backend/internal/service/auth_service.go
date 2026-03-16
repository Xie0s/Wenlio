// Package service 业务逻辑层
//
// 职责：实现认证相关业务逻辑（登录、密码修改）
// 对外接口：AuthService
package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"docplatform/internal/config"
	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/captcha"
	"docplatform/pkg/constants"
	"docplatform/pkg/errcode"
	"docplatform/pkg/jwt"
	"docplatform/pkg/logger"
	"docplatform/pkg/storage"
	"docplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// AuthService 认证业务
type AuthService struct {
	userRepo       *repository.UserRepo
	tenantRepo     *repository.TenantRepo
	captchaManager *captcha.Manager
}

// NewAuthService 创建 AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:       repository.NewUserRepo(),
		tenantRepo:     repository.NewTenantRepo(),
		captchaManager: captcha.DefaultManager(),
	}
}

// fetchTenantName 按 tenantID 查租户名，查不到时静默返回空字符串
func (s *AuthService) fetchTenantName(ctx context.Context, tenantID string) string {
	if tenantID == "" {
		return ""
	}
	tenant, err := s.tenantRepo.FindByIDTyped(ctx, tenantID)
	if err != nil {
		return ""
	}
	return tenant.Name
}

func (s *AuthService) buildLoginResp(ctx context.Context, user *entity.User) (*dto.LoginResp, *errcode.AppError) {
	token, err := jwt.GenerateToken(user.ID.Hex(), user.TenantID, user.Role)
	if err != nil {
		logger.L().Error("签发 Token 失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	tenantName := s.fetchTenantName(ctx, user.TenantID)

	return &dto.LoginResp{
		AccessToken: token,
		ExpiresIn:   86400,
		User: dto.UserInfo{
			ID:           user.ID.Hex(),
			Name:         user.Name,
			Username:     user.Username,
			Role:         user.Role,
			TenantID:     user.TenantID,
			TenantName:   tenantName,
			AvatarURL:    user.AvatarURL,
			Bio:          user.Bio,
			ProfileBgURL: user.ProfileBgURL,
		},
	}, nil
}

func (s *AuthService) consumeCaptchaToken(scene, token string, clientCtx captcha.ClientContext) *errcode.AppError {
	err := s.captchaManager.ConsumeToken(scene, token, clientCtx)
	if err == nil {
		return nil
	}
	if errors.Is(err, captcha.ErrTokenExpired) {
		return errcode.ErrCaptchaExpired
	}
	if errors.Is(err, captcha.ErrTokenInvalid) {
		return errcode.ErrCaptchaInvalid
	}
	logger.L().Error("消费验证码 token 失败", zap.Error(err), zap.String("scene", scene))
	return errcode.ErrInternalServer
}

// CreateCaptchaChallenge 创建安全检查 challenge
func (s *AuthService) CreateCaptchaChallenge(_ context.Context, clientCtx captcha.ClientContext, req *dto.CaptchaChallengeReq) (*dto.CaptchaChallengeResp, *errcode.AppError) {
	challenge, err := s.captchaManager.CreateChallenge(req.Scene, clientCtx)
	if err != nil {
		switch {
		case errors.Is(err, captcha.ErrInvalidScene):
			return nil, errcode.ErrInvalidParam
		case errors.Is(err, captcha.ErrTooManyRequests):
			return nil, errcode.ErrCaptchaTooFrequent
		}
		logger.L().Error("创建安全检查 challenge 失败", zap.Error(err), zap.String("scene", req.Scene))
		return nil, errcode.ErrInternalServer
	}

	return &dto.CaptchaChallengeResp{
		ChallengeID:   challenge.ID,
		Scene:         challenge.Scene,
		Mode:          challenge.Mode,
		Prompt:        challenge.Prompt,
		ExpiresAt:     challenge.ExpiresAt,
		MinDecisionMs: challenge.MinDecisionMs,
	}, nil
}

// VerifyCaptcha 完成安全检查并签发一次性验证码 token
func (s *AuthService) VerifyCaptcha(_ context.Context, clientCtx captcha.ClientContext, req *dto.VerifyCaptchaReq) (*dto.VerifyCaptchaResp, *errcode.AppError) {
	result, err := s.captchaManager.Verify(&captcha.VerifyInput{
		ChallengeID: req.ChallengeID,
		Scene:       req.Scene,
		DurationMs:  req.DurationMs,
		Signals: captcha.SignalSummary{
			DwellMs:             req.Signals.DwellMs,
			VisibleMs:           req.Signals.VisibleMs,
			FocusedMs:           req.Signals.FocusedMs,
			VisibilityChanges:   req.Signals.VisibilityChanges,
			FocusChanges:        req.Signals.FocusChanges,
			PointerEvents:       req.Signals.PointerEvents,
			KeyEvents:           req.Signals.KeyEvents,
			TrustedClick:        req.Signals.TrustedClick,
			Language:            req.Signals.Language,
			Platform:            req.Signals.Platform,
			ScreenWidth:         req.Signals.ScreenWidth,
			ScreenHeight:        req.Signals.ScreenHeight,
			TimezoneOffset:      req.Signals.TimezoneOffset,
			TouchPoints:         req.Signals.TouchPoints,
			HardwareConcurrency: req.Signals.HardwareConcurrency,
			Webdriver:           req.Signals.Webdriver,
		},
	}, clientCtx)
	if err != nil {
		switch {
		case errors.Is(err, captcha.ErrTooManyRequests), errors.Is(err, captcha.ErrClientBlocked):
			return nil, errcode.ErrCaptchaTooFrequent
		case errors.Is(err, captcha.ErrChallengeExpired):
			return nil, errcode.ErrCaptchaExpired
		case errors.Is(err, captcha.ErrChallengeNotFound), errors.Is(err, captcha.ErrChallengeFailed), errors.Is(err, captcha.ErrInvalidScene):
			return nil, errcode.ErrCaptchaInvalid
		default:
			logger.L().Error("安全检查校验失败", zap.Error(err), zap.String("scene", req.Scene))
			return nil, errcode.ErrInternalServer
		}
	}

	return &dto.VerifyCaptchaResp{
		CaptchaToken: result.Token,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, clientCtx captcha.ClientContext, req *dto.LoginReq) (*dto.LoginResp, *errcode.AppError) {
	if appErr := s.consumeCaptchaToken(captcha.SceneLogin, req.CaptchaToken, clientCtx); appErr != nil {
		return nil, appErr
	}

	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errcode.ErrLoginFailed
	}

	// 检查账号锁定
	if time.Now().Before(user.LockedUntil) {
		return nil, errcode.ErrAccountLocked
	}

	// 检查账号状态
	if user.Status != constants.UserStatusActive {
		return nil, errcode.ErrAccountDisabled
	}
	if user.TenantID != "" {
		tenant, tenantErr := s.tenantRepo.FindByIDTyped(ctx, user.TenantID)
		if tenantErr != nil {
			return nil, errcode.ErrTenantNotFound
		}
		if tenant.Status == constants.TenantStatusDeleting {
			return nil, errcode.ErrTenantDeleting
		}
		if tenant.Status != constants.TenantStatusActive {
			return nil, errcode.ErrTenantSuspended
		}
	}

	// 校验密码
	if !utils.CheckPassword(req.Password, user.Password) {
		failCount := user.LoginFailCount + 1
		var lockedUntil time.Time
		if failCount >= 5 {
			lockedUntil = time.Now().Add(15 * time.Minute)
		}
		_ = s.userRepo.UpdateLoginFail(ctx, user.ID, failCount, lockedUntil)
		return nil, errcode.ErrLoginFailed
	}

	// 登录成功，重置失败计数
	_ = s.userRepo.ResetLoginFail(ctx, user.ID)

	return s.buildLoginResp(ctx, user)
}

func (s *AuthService) Register(ctx context.Context, clientCtx captcha.ClientContext, req *dto.RegisterReq) (*dto.LoginResp, *errcode.AppError) {
	if appErr := s.consumeCaptchaToken(captcha.SceneRegister, req.CaptchaToken, clientCtx); appErr != nil {
		return nil, appErr
	}

	if !isStrongPassword(req.AdminPassword) {
		return nil, errcode.ErrPasswordTooWeak
	}

	createReq := &dto.CreateTenantReq{
		ID:            req.TenantID,
		Name:          req.TenantName,
		LogoURL:       req.LogoURL,
		AdminUsername: req.AdminUsername,
		AdminPassword: req.AdminPassword,
		AdminName:     req.AdminName,
		AdminEmail:    req.AdminEmail,
	}

	tenantService := NewTenantService()
	if _, appErr := tenantService.Create(ctx, createReq); appErr != nil {
		return nil, appErr
	}

	user, err := s.userRepo.FindByUsername(ctx, req.AdminUsername)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	logger.L().Info("用户注册并创建租户成功",
		zap.String("tenant_id", req.TenantID),
		zap.String("user_id", user.ID.Hex()),
		zap.String("username", user.Username),
	)

	return s.buildLoginResp(ctx, user)
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*dto.UserInfo, *errcode.AppError) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}
	user, err := s.userRepo.FindByIDTyped(ctx, oid)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	if user.TenantID != "" {
		tenant, tenantErr := s.tenantRepo.FindByIDTyped(ctx, user.TenantID)
		if tenantErr != nil {
			return nil, errcode.ErrTenantNotFound
		}
		if tenant.Status == constants.TenantStatusDeleting {
			return nil, errcode.ErrTenantDeleting
		}
		if tenant.Status != constants.TenantStatusActive {
			return nil, errcode.ErrTenantSuspended
		}
	}
	tenantName := s.fetchTenantName(ctx, user.TenantID)

	return &dto.UserInfo{
		ID:           user.ID.Hex(),
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		TenantID:     user.TenantID,
		TenantName:   tenantName,
		AvatarURL:    user.AvatarURL,
		Bio:          user.Bio,
		ProfileBgURL: user.ProfileBgURL,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordReq) *errcode.AppError {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errcode.ErrTokenInvalid
	}
	user, err := s.userRepo.FindByIDTyped(ctx, oid)
	if err != nil {
		return errcode.ErrUserNotFound
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errcode.ErrPasswordIncorrect
	}

	if !isStrongPassword(req.NewPassword) {
		return errcode.ErrPasswordTooWeak
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errcode.ErrInternalServer
	}

	return toAppError(s.userRepo.UpdateByID(ctx, oid, map[string]interface{}{"password": hashed}))
}

// isStrongPassword 校验密码强度（至少8位，含字母和数字）
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter, hasDigit := false, false
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

// toAppError 将普通 error 转为 AppError（数据库错误统一处理）
func toAppError(err error) *errcode.AppError {
	if err == nil {
		return nil
	}
	return errcode.ErrDatabase.Wrap(err)
}

// ToAppErrorNotFound 将查询 error 转为 NotFound AppError
func ToAppErrorNotFound(err error, notFoundErr *errcode.AppError) *errcode.AppError {
	if err == nil {
		return nil
	}
	if err.Error() == "mongo: no documents in result" {
		return notFoundErr
	}
	return errcode.ErrDatabase.Wrap(err)
}

// UserFromEntity 从 entity.User 转换为 dto.UserInfo（去除敏感字段，TenantName 需调用方填充）
func UserFromEntity(u *entity.User) *dto.UserInfo {
	return &dto.UserInfo{
		ID:           u.ID.Hex(),
		Name:         u.Name,
		Username:     u.Username,
		Role:         u.Role,
		TenantID:     u.TenantID,
		AvatarURL:    u.AvatarURL,
		Bio:          u.Bio,
		ProfileBgURL: u.ProfileBgURL,
	}
}

// UpdateProfile 更新当前用户资料（name / bio 可选更新）
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileReq) (*dto.UserInfo, *errcode.AppError) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}
	fields := map[string]interface{}{}
	if req.Name != nil {
		fields["name"] = *req.Name
	}
	if req.Bio != nil {
		fields["bio"] = *req.Bio
	}
	if len(fields) > 0 {
		if dbErr := s.userRepo.UpdateByID(ctx, oid, fields); dbErr != nil {
			return nil, errcode.ErrDatabase.Wrap(dbErr)
		}
	}
	user, err := s.userRepo.FindByIDTyped(ctx, oid)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	tenantName := s.fetchTenantName(ctx, user.TenantID)
	return &dto.UserInfo{
		ID:           user.ID.Hex(),
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		TenantID:     user.TenantID,
		TenantName:   tenantName,
		AvatarURL:    user.AvatarURL,
		Bio:          user.Bio,
		ProfileBgURL: user.ProfileBgURL,
	}, nil
}

// uploadUserAsset 上传用户个人资产文件（头像 / 背景），自动删除旧文件
// assetType: "avatar" | "profile_bg"
func (s *AuthService) uploadUserAsset(ctx context.Context, userID primitive.ObjectID, file *multipart.FileHeader, assetType string) (string, *errcode.AppError) {
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		return "", errcode.ErrUploadFileTooLarge
	}
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if !strings.HasPrefix(mimeType, "image/") {
		return "", errcode.ErrUploadFileTypeInvalid
	}

	cfg := config.Get()
	provider := storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")

	ext := strings.ToLower(filepath.Ext(file.Filename))
	key := fmt.Sprintf("_user_assets/%s/%s/%s%s", userID.Hex(), assetType, primitive.NewObjectID().Hex(), ext)

	src, err := file.Open()
	if err != nil {
		return "", errcode.ErrInternalServer.Wrap(err)
	}
	defer src.Close()

	result, err := provider.Upload(ctx, storage.UploadInput{
		Key:         key,
		Content:     src,
		Size:        file.Size,
		ContentType: mimeType,
	})
	if err != nil {
		return "", errcode.ErrUploadStorageFail.Wrap(err)
	}
	return result.URL, nil
}

// deleteOldUserAsset 删除用户旧的本地资产文件（静默失败）
func deleteOldUserAsset(oldURL string) {
	if oldURL == "" || !strings.HasPrefix(oldURL, "/uploads/") {
		return
	}
	cfg := config.Get()
	key := strings.TrimPrefix(oldURL, "/uploads/")
	provider := storage.NewLocalProvider(cfg.Storage.LocalPath, "/uploads")
	if err := provider.Delete(context.Background(), key); err != nil {
		logger.L().Warn("删除旧用户资产失败", zap.String("key", key), zap.Error(err))
	}
}

// UploadAvatar 上传用户头像，替换旧文件，返回更新后的 UserInfo
func (s *AuthService) UploadAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (*dto.UserInfo, *errcode.AppError) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}
	user, err := s.userRepo.FindByIDTyped(ctx, oid)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	newURL, appErr := s.uploadUserAsset(ctx, oid, file, "avatar")
	if appErr != nil {
		return nil, appErr
	}
	oldURL := user.AvatarURL
	if dbErr := s.userRepo.UpdateByID(ctx, oid, map[string]interface{}{"avatar_url": newURL}); dbErr != nil {
		return nil, errcode.ErrDatabase.Wrap(dbErr)
	}
	deleteOldUserAsset(oldURL)
	user.AvatarURL = newURL
	tenantName := s.fetchTenantName(ctx, user.TenantID)
	return &dto.UserInfo{
		ID:           user.ID.Hex(),
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		TenantID:     user.TenantID,
		TenantName:   tenantName,
		AvatarURL:    user.AvatarURL,
		Bio:          user.Bio,
		ProfileBgURL: user.ProfileBgURL,
	}, nil
}

// UploadProfileBg 上传用户主页背景，替换旧文件，返回更新后的 UserInfo
func (s *AuthService) UploadProfileBg(ctx context.Context, userID string, file *multipart.FileHeader) (*dto.UserInfo, *errcode.AppError) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}
	user, err := s.userRepo.FindByIDTyped(ctx, oid)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}
	newURL, appErr := s.uploadUserAsset(ctx, oid, file, "profile_bg")
	if appErr != nil {
		return nil, appErr
	}
	oldURL := user.ProfileBgURL
	if dbErr := s.userRepo.UpdateByID(ctx, oid, map[string]interface{}{"profile_bg_url": newURL}); dbErr != nil {
		return nil, errcode.ErrDatabase.Wrap(dbErr)
	}
	deleteOldUserAsset(oldURL)
	user.ProfileBgURL = newURL
	tenantName := s.fetchTenantName(ctx, user.TenantID)
	return &dto.UserInfo{
		ID:           user.ID.Hex(),
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		TenantID:     user.TenantID,
		TenantName:   tenantName,
		AvatarURL:    user.AvatarURL,
		Bio:          user.Bio,
		ProfileBgURL: user.ProfileBgURL,
	}, nil
}
