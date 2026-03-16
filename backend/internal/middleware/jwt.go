// Package middleware HTTP 中间件集合
//
// 职责：JWT Token 解析与验证，将认证上下文注入 Gin Context
// 对外接口：JWTAuth(), AuthContext
package middleware

import (
	"strings"

	"docplatform/internal/repository"
	"docplatform/pkg/constants"
	"docplatform/pkg/ctxutil"
	"docplatform/pkg/errcode"
	"docplatform/pkg/jwt"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthContext 认证上下文，由 JWTAuth 中间件注入 Gin Context
type AuthContext struct {
	UserID   string // 当前用户 ID（ObjectID hex）
	TenantID string // 租户 ID（字符串，super_admin 为空）
	Role     string // 角色：super_admin / tenant_admin
}

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, errcode.ErrTokenInvalid)
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, errcode.ErrTokenExpired)
			c.Abort()
			return
		}

		if claims.TenantID != "" {
			tenantRepo := repository.NewTenantRepo()
			tenant, tenantErr := tenantRepo.FindByIDTyped(c.Request.Context(), claims.TenantID)
			if tenantErr != nil {
				response.Fail(c, errcode.ErrUnauthorized)
				c.Abort()
				return
			}
			if tenant.Status == constants.TenantStatusDeleting {
				response.Fail(c, errcode.ErrUnauthorized)
				c.Abort()
				return
			}
			if tenant.Status != constants.TenantStatusActive {
				response.Fail(c, errcode.ErrUnauthorized)
				c.Abort()
				return
			}
			// 缓存已验证的租户状态到 Go 标准 context，供 EnsureActive 复用，避免重复查库
			c.Request = c.Request.WithContext(ctxutil.WithTenantStatus(c.Request.Context(), tenant.Status))
		}

		c.Set("auth_context", &AuthContext{
			UserID:   claims.UserID,
			TenantID: claims.TenantID,
			Role:     claims.Role,
		})
		c.Next()
	}
}

// OptionalJWTAuth 可选 JWT 认证中间件
// 与 JWTAuth 区别：无 Token 或 Token 无效时不拒绝请求，仅在 Token 合法时注入 auth_context
// 用途：公开接口需要根据登录状态返回不同数据（如主题列表按 access_mode 过滤）
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			// Token 无效/过期，静默放行，视为未登录
			c.Next()
			return
		}
		// Token 合法，注入认证上下文（不做租户状态校验，公开接口不依赖租户活跃状态）
		c.Set("auth_context", &AuthContext{
			UserID:   claims.UserID,
			TenantID: claims.TenantID,
			Role:     claims.Role,
		})
		c.Next()
	}
}

// GetAuthContext 从 Gin Context 获取认证上下文
func GetAuthContext(c *gin.Context) *AuthContext {
	val, exists := c.Get("auth_context")
	if !exists {
		return nil
	}
	return val.(*AuthContext)
}
