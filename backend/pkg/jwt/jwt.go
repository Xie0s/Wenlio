// Package jwt JWT 工具
//
// 职责：JWT Token 的签发与解析，封装 Claims 结构
// 对外接口：GenerateToken() 签发 Token，ParseToken() 解析 Token
package jwt

import (
	"time"

	"docplatform/internal/config"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   string `json:"sub"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
	jwtv5.RegisteredClaims
}

// GenerateToken 签发 JWT Token
func GenerateToken(userID, tenantID, role string) (string, error) {
	cfg := config.Get()
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)),
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	cfg := config.Get()
	token, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwtv5.ErrTokenInvalidClaims
}

// ThemeAccessClaims 主题访问验证码通过后签发的短期 Token Claims
// 用途：access_mode="code" 的主题，验证码校验成功后颁发，前端存于 sessionStorage
type ThemeAccessClaims struct {
	ThemeID string `json:"theme_id"`
	jwtv5.RegisteredClaims
}

// GenerateThemeAccessToken 签发主题访问 Token（24小时有效）
func GenerateThemeAccessToken(themeID string) (string, error) {
	cfg := config.Get()
	now := time.Now()
	claims := ThemeAccessClaims{
		ThemeID: themeID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(24 * time.Hour)),
			Subject:   "theme_access",
		},
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseThemeAccessToken 解析主题访问 Token，返回 theme_id
func ParseThemeAccessToken(tokenStr string) (string, error) {
	cfg := config.Get()
	token, err := jwtv5.ParseWithClaims(tokenStr, &ThemeAccessClaims{}, func(t *jwtv5.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*ThemeAccessClaims); ok && token.Valid && claims.Subject == "theme_access" {
		return claims.ThemeID, nil
	}
	return "", jwtv5.ErrTokenInvalidClaims
}
