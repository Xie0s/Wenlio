// Package constants 系统常量定义
//
// 职责：定义租户 ID 保留词列表
// 对外接口：ReservedTenantIDs
package constants

// ReservedTenantIDs 系统保留的路径段，禁止用作 tenant_id
var ReservedTenantIDs = map[string]struct{}{
	"admin":       {},
	"api":         {},
	"assets":      {},
	"static":      {},
	"health":      {},
	"favicon.ico": {},
	"robots.txt":  {},
	"sitemap":     {},
}
