// Package middleware HTTP 中间件集合
//
// 职责：注入安全响应头，防止点击劫持、MIME 嗅探、信息泄露等攻击
// 对外接口：SecurityHeaders()
package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders 安全响应头中间件
// 在每个 HTTP 响应中添加以下安全头：
//   - X-Content-Type-Options: nosniff   → 禁止浏览器 MIME 嗅探
//   - X-Frame-Options: DENY             → 防止页面被嵌入 iframe（点击劫持）
//   - Referrer-Policy: strict-origin-when-cross-origin → 跨域时仅发送 origin
//   - Permissions-Policy: camera=(), microphone=(), geolocation=() → 禁用敏感权限
//   - X-XSS-Protection: 1; mode=block   → IE/旧浏览器 XSS 过滤器
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
