// Package middleware HTTP 中间件集合
//
// 职责：记录请求/响应的结构化日志
// 对外接口：Logger()
package middleware

import (
	"fmt"
	"time"

	"docplatform/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	clReset  = "\033[0m"
	clRed    = "\033[31m"
	clGreen  = "\033[32m"
	clYellow = "\033[33m"
	clBlue   = "\033[34m"
	clPurple = "\033[35m"
	clCyan   = "\033[36m"
	clDim    = "\033[2m"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	noCaller := logger.L().WithOptions(zap.WithCaller(false))
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		noCaller.Info(httpLogLine(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
		))
	}
}

// httpLogLine 格式化 HTTP 请求日志行（带 ANSI 颜色）
func httpLogLine(method, path string, status int, latency time.Duration, ip string) string {
	// 方法颜色
	mc := clBlue
	switch method {
	case "POST":
		mc = clGreen
	case "PUT", "PATCH":
		mc = clYellow
	case "DELETE":
		mc = clRed
	case "OPTIONS":
		mc = clPurple
	}

	// 状态码颜色
	sc := clGreen
	switch {
	case status >= 500:
		sc = clRed
	case status >= 400:
		sc = clYellow
	case status >= 300:
		sc = clCyan
	}

	// 延迟格式化
	var lat string
	switch {
	case latency < time.Millisecond:
		lat = fmt.Sprintf("%dµs", latency.Microseconds())
	case latency < time.Second:
		lat = fmt.Sprintf("%dms", latency.Milliseconds())
	default:
		lat = fmt.Sprintf("%.2fs", latency.Seconds())
	}

	return fmt.Sprintf("%s%-7s%s %-45s %s%d%s %7s  %s%s%s",
		mc, method, clReset,
		path,
		sc, status, clReset,
		lat,
		clDim, ip, clReset,
	)
}
