// Package middleware HTTP 中间件集合
//
// 职责：panic 恢复，防止单个请求崩溃导致服务停止
// 对外接口：Recovery()
package middleware

import (
	"net/http"

	"docplatform/pkg/errcode"
	"docplatform/pkg/logger"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L().Error("panic recovered",
					zap.Any("error", r),
					zap.String("request_id", c.GetString("request_id")),
					zap.Stack("stack"),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:      errcode.ErrInternalServer.Code,
					Message:   errcode.ErrInternalServer.Message,
					RequestID: c.GetString("request_id"),
				})
			}
		}()
		c.Next()
	}
}
