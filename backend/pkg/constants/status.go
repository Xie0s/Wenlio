// Package constants 系统常量定义
//
// 职责：定义各实体状态常量
// 对外接口：版本/文档页/评论/租户/用户状态常量
package constants

// 版本状态
const (
	VersionStatusDraft     = "draft"
	VersionStatusPublished = "published"
	VersionStatusArchived  = "archived"
)

// 文档页状态
const (
	PageStatusDraft     = "draft"
	PageStatusPublished = "published"
)

// 评论状态
const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
)

// 租户状态
const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusDeleting  = "deleting"
)

// 用户状态
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)
