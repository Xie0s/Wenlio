package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	tests := []struct {
		header string
		want   string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"X-XSS-Protection", "1; mode=block"},
	}

	for _, tt := range tests {
		got := w.Header().Get(tt.header)
		if got != tt.want {
			t.Errorf("header %s = %q, want %q", tt.header, got, tt.want)
		}
	}

	pp := w.Header().Get("Permissions-Policy")
	if pp == "" {
		t.Error("Permissions-Policy header not set")
	}
}

func TestCORSWithoutAllowedOrigins(t *testing.T) {
	// 当 ALLOWED_ORIGINS 未设置时（测试环境默认），
	// 请求带 Origin 头时仍应返回 CORS 头（开发兼容模式）
	gin.SetMode(gin.TestMode)

	// 保存并清空 allowedOrigins
	saved := allowedOrigins
	allowedOrigins = nil
	defer func() { allowedOrigins = saved }()

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	r.ServeHTTP(w, req)

	acao := w.Header().Get("Access-Control-Allow-Origin")
	if acao != "http://localhost:5173" {
		t.Errorf("CORS without whitelist: Access-Control-Allow-Origin = %q, want %q", acao, "http://localhost:5173")
	}
}

func TestCORSWithAllowedOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 设置白名单
	saved := allowedOrigins
	allowedOrigins = []string{"https://example.com", "https://admin.example.com"}
	defer func() { allowedOrigins = saved }()

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 允许的 origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	acao := w.Header().Get("Access-Control-Allow-Origin")
	if acao != "https://example.com" {
		t.Errorf("Allowed origin: Access-Control-Allow-Origin = %q, want %q", acao, "https://example.com")
	}

	// 不允许的 origin
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://evil.com")
	r.ServeHTTP(w2, req2)

	acao2 := w2.Header().Get("Access-Control-Allow-Origin")
	if acao2 != "" {
		t.Errorf("Blocked origin: Access-Control-Allow-Origin = %q, want empty", acao2)
	}
}

func TestCORSPreflightReturns204(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := allowedOrigins
	allowedOrigins = nil
	defer func() { allowedOrigins = saved }()

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestIPRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	// 允许每10秒最多3次请求
	r.Use(IPRateLimit(3, 10*time.Second))
	r.POST("/verify", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 发送4次请求，第4次应该被限速
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/verify", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Request %d: status = %d, want 200", i+1, w.Code)
		}
	}

	// 第4次请求应被限速
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/verify", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Request 4 (rate limited): status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}
