package jwt

import (
	"testing"
	"time"

	"docplatform/internal/config"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// setupTestConfig 为测试初始化最小配置
func setupTestConfig() {
	// 直接通过包级变量注入测试配置（避免读取文件）
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret-for-jwt-unit-tests"
	cfg.JWT.ExpireHours = 24
	config.SetForTest(cfg)
}

func TestGenerateAndParseToken(t *testing.T) {
	setupTestConfig()

	token, err := GenerateToken("user123", "tenant456", "admin")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != "user123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user123")
	}
	if claims.TenantID != "tenant456" {
		t.Errorf("TenantID = %q, want %q", claims.TenantID, "tenant456")
	}
	if claims.Role != "admin" {
		t.Errorf("Role = %q, want %q", claims.Role, "admin")
	}
	if claims.Issuer != "docplatform" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "docplatform")
	}
}

func TestParseTokenRejectsThemeAccessToken(t *testing.T) {
	setupTestConfig()

	// 签发一个 ThemeAccess token
	themeToken, err := GenerateThemeAccessToken("theme789")
	if err != nil {
		t.Fatalf("GenerateThemeAccessToken() error = %v", err)
	}

	// 尝试用 ParseToken（主认证解析器）解析 ThemeAccess token，应被拒绝
	_, err = ParseToken(themeToken)
	if err == nil {
		t.Fatal("ParseToken() should reject ThemeAccess token, but got nil error")
	}
}

func TestParseThemeAccessTokenRejectsRegularToken(t *testing.T) {
	setupTestConfig()

	// 签发一个常规用户 token
	userToken, err := GenerateToken("user123", "tenant456", "admin")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// 尝试用 ParseThemeAccessToken 解析常规 token，应被拒绝
	_, err = ParseThemeAccessToken(userToken)
	if err == nil {
		t.Fatal("ParseThemeAccessToken() should reject regular user token, but got nil error")
	}
}

func TestThemeAccessTokenRoundTrip(t *testing.T) {
	setupTestConfig()

	themeID := "theme_abc123"
	token, err := GenerateThemeAccessToken(themeID)
	if err != nil {
		t.Fatalf("GenerateThemeAccessToken() error = %v", err)
	}

	parsed, err := ParseThemeAccessToken(token)
	if err != nil {
		t.Fatalf("ParseThemeAccessToken() error = %v", err)
	}
	if parsed != themeID {
		t.Errorf("ParseThemeAccessToken() = %q, want %q", parsed, themeID)
	}
}

func TestParseTokenRejectsExpired(t *testing.T) {
	setupTestConfig()
	cfgVal := config.Get()

	// 手动构造一个过期的 token
	now := time.Now()
	claims := Claims{
		UserID:   "user123",
		TenantID: "tenant456",
		Role:     "admin",
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    "docplatform",
			IssuedAt:  jwtv5.NewNumericDate(now.Add(-48 * time.Hour)),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(-24 * time.Hour)),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfgVal.JWT.Secret))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	_, err = ParseToken(tokenStr)
	if err == nil {
		t.Fatal("ParseToken() should reject expired token")
	}
}

func TestParseTokenRejectsTamperedToken(t *testing.T) {
	setupTestConfig()

	token, err := GenerateToken("user123", "tenant456", "admin")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// 篡改 token 的最后一个字符
	tampered := token[:len(token)-1] + "X"
	_, err = ParseToken(tampered)
	if err == nil {
		t.Fatal("ParseToken() should reject tampered token")
	}
}
