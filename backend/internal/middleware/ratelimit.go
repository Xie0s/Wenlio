// Package middleware HTTP 中间件集合
//
// 职责：基于 IP 的简单速率限制中间件，防止暴力破解主题访问码等敏感接口
// 对外接口：IPRateLimit()
package middleware

import (
	"net/http"
	"sync"
	"time"

	"docplatform/pkg/errcode"
	"docplatform/pkg/response"

	"github.com/gin-gonic/gin"
)

type ipRateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

// IPRateLimit 基于客户端 IP 的速率限制中间件
// maxAttempts: 时间窗口内允许的最大请求次数
// window: 滑动窗口时长
func IPRateLimit(maxAttempts int, window time.Duration) gin.HandlerFunc {
	limiter := &ipRateLimiter{
		attempts: make(map[string][]time.Time),
	}
	// 定期清理过期条目，避免内存泄漏
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.prune(window)
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip, maxAttempts, window) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Code:      errcode.ErrCaptchaTooFrequent.Code,
				Message:   "请求过于频繁，请稍后再试",
				RequestID: c.GetString("request_id"),
			})
			return
		}
		c.Next()
	}
}

func (l *ipRateLimiter) allow(key string, maxAttempts int, window time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// 保留窗口内的请求记录
	existing := l.attempts[key]
	valid := existing[:0]
	for _, t := range existing {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= maxAttempts {
		l.attempts[key] = valid
		return false
	}

	l.attempts[key] = append(valid, now)
	return true
}

func (l *ipRateLimiter) prune(window time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-window)
	for key, times := range l.attempts {
		valid := times[:0]
		for _, t := range times {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(l.attempts, key)
		} else {
			l.attempts[key] = valid
		}
	}
}
