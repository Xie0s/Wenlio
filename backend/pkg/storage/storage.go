// Package storage 存储策略抽象层
//
// 职责：定义统一的文件存储接口，支持本地和 S3 兼容存储的热插拔切换
// 对外接口：Provider 接口、UploadInput/UploadResult 数据结构
package storage

import (
	"context"
	"io"
)

// Provider 存储后端接口，所有实现必须满足此契约
type Provider interface {
	// Upload 上传文件，返回可访问 URL 和存储 Key
	Upload(ctx context.Context, input UploadInput) (UploadResult, error)
	// Delete 按存储 Key 删除文件
	Delete(ctx context.Context, key string) error
	// ListKeys 列出指定前缀下的所有存储 Key（用于审计对照）
	ListKeys(ctx context.Context, prefix string) ([]string, error)
	// Type 返回存储类型标识（"local" | "cloud"）
	Type() string
}

// UploadInput 上传请求参数
type UploadInput struct {
	Key         string    // 目标存储路径（含目录前缀）
	Content     io.Reader // 文件内容流
	Size        int64     // 文件大小（字节）
	ContentType string    // MIME 类型
}

// UploadResult 上传响应
type UploadResult struct {
	URL string // 公开可访问 URL
	Key string // 存储 Key（删除时使用）
}
