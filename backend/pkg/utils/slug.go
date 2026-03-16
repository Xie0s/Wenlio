// Package utils 工具函数集合
//
// 职责：Slug 生成与校验工具
// 对外接口：GenerateSlug() 从标题生成 Slug，ValidateSlug() 校验 Slug 格式
package utils

import (
	"regexp"
	"strings"
)

var (
	slugPattern   = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)
	nonAlphaNum   = regexp.MustCompile(`[^a-z0-9-]+`)
	multiDash     = regexp.MustCompile(`-+`)
)

// GenerateSlug 从标题生成 URL 友好的 slug
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = nonAlphaNum.ReplaceAllString(slug, "-")
	slug = multiDash.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// ValidateSlug 校验 slug 格式是否合法（小写字母、数字、连字符，长度 2-128）
func ValidateSlug(slug string) bool {
	if len(slug) < 2 || len(slug) > 128 {
		return false
	}
	return slugPattern.MatchString(slug)
}
