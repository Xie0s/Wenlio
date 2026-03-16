// Package ctxutil 请求上下文工具
//
// 职责：在中间件和业务层之间传递已验证的上下文数据，避免重复查库
// 对外接口：WithTenantStatus / TenantStatus
package ctxutil

import "context"

type tenantStatusKey struct{}

// WithTenantStatus 将已验证的租户状态写入 context（由 JWTAuth 中间件调用）
func WithTenantStatus(ctx context.Context, status string) context.Context {
	return context.WithValue(ctx, tenantStatusKey{}, status)
}

// TenantStatus 从 context 读取已缓存的租户状态（未命中时 ok=false，需回退查库）
func TenantStatus(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(tenantStatusKey{}).(string)
	return s, ok
}
