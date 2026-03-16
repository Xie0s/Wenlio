// Package storage 存储策略抽象层
//
// 职责：本地文件系统存储实现
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalProvider 本地文件系统存储
type LocalProvider struct {
	BasePath string // 本地根目录，如 "./uploads"
	BaseURL  string // 可访问 URL 前缀，如 "/uploads"
}

// NewLocalProvider 创建本地存储 Provider
func NewLocalProvider(basePath, baseURL string) *LocalProvider {
	return &LocalProvider{BasePath: basePath, BaseURL: baseURL}
}

// Upload 将文件写入本地文件系统
func (p *LocalProvider) Upload(ctx context.Context, input UploadInput) (UploadResult, error) {
	fullPath := filepath.Join(p.BasePath, filepath.FromSlash(input.Key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return UploadResult{}, fmt.Errorf("创建目录失败: %w", err)
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return UploadResult{}, fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()
	if _, err = io.Copy(f, input.Content); err != nil {
		return UploadResult{}, fmt.Errorf("写入文件失败: %w", err)
	}
	return UploadResult{
		URL: p.BaseURL + "/" + input.Key,
		Key: input.Key,
	}, nil
}

// Delete 删除本地文件（文件不存在时不报错）
func (p *LocalProvider) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(p.BasePath, filepath.FromSlash(key))
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// ListKeys 列出指定前缀下的所有本地文件 Key
func (p *LocalProvider) ListKeys(ctx context.Context, prefix string) ([]string, error) {
	root := filepath.Join(p.BasePath, filepath.FromSlash(prefix))
	var keys []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(p.BasePath, path)
		if relErr != nil {
			return relErr
		}
		keys = append(keys, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("遍历本地目录失败: %w", err)
	}
	return keys, nil
}

// Type 返回存储类型标识
func (p *LocalProvider) Type() string { return "local" }
