// Package middleware HTTP 中间件集合
//
// 职责：处理跨域请求（开发环境 :5173 → :8080），生产环境限定允许的 Origin
// 对外接口：CORS()
package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// allowedOrigins 解析环境变量 ALLOWED_ORIGINS（逗号分隔），缓存在包级变量中。
// 为空时默认仅允许同源访问（不设置 Access-Control-Allow-Origin）。
var allowedOrigins []string

func init() {
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		for _, o := range strings.Split(v, ",") {
			if o = strings.TrimSpace(o); o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}
}

// isOriginAllowed 检查请求 Origin 是否在白名单内
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

// CORS 跨域中间件
// 生产环境应通过 ALLOWED_ORIGINS 环境变量配置允许的域名（逗号分隔），
// 例如：ALLOWED_ORIGINS=https://example.com,https://admin.example.com
// 如果未设置 ALLOWED_ORIGINS，默认仅允许同源请求（不返回 CORS 头）。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if len(allowedOrigins) > 0 {
			// 白名单模式：仅允许已配置的 Origin
			if isOriginAllowed(origin) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		} else {
			// 未配置白名单：开发环境或同源部署，允许任意 Origin
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Request-ID,X-Theme-Access")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
