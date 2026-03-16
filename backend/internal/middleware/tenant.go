// Package middleware HTTP 中间件集合
//
// 职责：租户数据隔离，从 JWT 提取 tenant_id 并注入 Gin Context
// 对外接口：TenantIsolation()
package middleware

import (
	"docplatform/pkg/constants"
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// TenantIsolation 租户数据隔离中间件
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := GetAuthContext(c)
		if auth == nil {
			response.Fail(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		// super_admin 不注入 tenant_id（平台管理接口专属）
		if auth.Role == constants.RoleSuperAdmin {
			c.Next()
			return
		}
		if auth.TenantID == "" {
			response.Fail(c, errcode.ErrForbidden)
			c.Abort()
			return
		}
		c.Set("tenant_id", auth.TenantID)
		c.Next()
	}
}

// GetTenantID 从 Gin Context 获取当前租户 ID
func GetTenantID(c *gin.Context) string {
	return c.GetString("tenant_id")
}
