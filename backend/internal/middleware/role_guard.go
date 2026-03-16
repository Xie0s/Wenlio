// Package middleware HTTP 中间件集合
//
// 职责：角色权限校验，不匹配则返回 403
// 对外接口：RoleGuard()
package middleware

import (
	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleGuard 角色校验中间件工厂
func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := GetAuthContext(c)
		if auth == nil {
			response.Fail(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		for _, role := range allowedRoles {
			if auth.Role == role {
				c.Next()
				return
			}
		}
		response.Fail(c, errcode.ErrForbidden)
		c.Abort()
	}
}
