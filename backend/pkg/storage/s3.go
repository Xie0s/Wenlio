// Package storage 存储策略抽象层
//
// 职责：S3 兼容存储实现（支持 AWS S3 / MinIO / Aliyun OSS / Cloudflare R2 等）
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Provider S3 兼容存储后端
type S3Provider struct {
	client         *s3.Client
	bucket         string
	customDomain   string // 可选，CDN 或自定义域名前缀
	endpoint       string
	region         string
	forcePathStyle bool
}

// S3Config S3 连接配置
type S3Config struct {
	Endpoint       string
	Region         string
	Bucket         string
	AccessKeyID    string
	SecretKey      string
	CustomDomain   string
	ForcePathStyle bool
}

// NewS3Provider 按配置创建 S3Provider
func NewS3Provider(cfg S3Config) (*S3Provider, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID, cfg.SecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("加载 S3 配置失败: %w", err)
	}

	clientOpts := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = cfg.ForcePathStyle
			// 关闭默认的 CRC32 Checksum 追加行为，仅在 API 要求时才计算
			// 否则 Cloudflare R2 / 阿里云 OSS / 腾讯 COS / MinIO 等兼容服务会因
			// 收到不认识的 x-amz-checksum-crc32 头而返回错误
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		},
	}
	if cfg.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)
	return &S3Provider{
		client:         client,
		bucket:         cfg.Bucket,
		customDomain:   cfg.CustomDomain,
		endpoint:       cfg.Endpoint,
		region:         cfg.Region,
		forcePathStyle: cfg.ForcePathStyle,
	}, nil
}

// Upload 上传文件到 S3 兼容存储
func (p *S3Provider) Upload(ctx context.Context, input UploadInput) (UploadResult, error) {
	data, err := io.ReadAll(input.Content)
	if err != nil {
		return UploadResult{}, fmt.Errorf("读取文件流失败: %w", err)
	}

	size := int64(len(data))
	_, err = p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(p.bucket),
		Key:           aws.String(input.Key),
		Body:          bytes.NewReader(data),
		ContentLength: &size,
		ContentType:   aws.String(input.ContentType),
	})
	if err != nil {
		return UploadResult{}, fmt.Errorf("上传到 S3 失败: %w", err)
	}

	url := p.buildURL(input.Key)
	return UploadResult{URL: url, Key: input.Key}, nil
}

// Delete 从 S3 兼容存储删除文件
func (p *S3Provider) Delete(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	return err
}

// ListKeys 列出指定前缀下的所有存储 Key
func (p *S3Provider) ListKeys(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	var continuationToken *string
	for {
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(p.bucket),
			Prefix: aws.String(prefix),
		}
		if continuationToken != nil {
			input.ContinuationToken = continuationToken
		}
		output, err := p.client.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("列举 S3 对象失败: %w", err)
		}
		for _, obj := range output.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
		if output.IsTruncated == nil || !*output.IsTruncated {
			break
		}
		continuationToken = output.NextContinuationToken
	}
	return keys, nil
}

// Type 返回存储类型标识
func (p *S3Provider) Type() string { return "cloud" }

// Ping 简单连通性探测（列出 1 个对象即可验证连接和权限）
func (p *S3Provider) Ping(ctx context.Context) error {
	maxKeys := int32(1)
	_, err := p.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(p.bucket),
		MaxKeys: &maxKeys,
	})
	return err
}

// buildURL 构建文件的公开访问 URL
//
// 优先级：自定义域名 > Path Style > Virtual-Hosted Style > AWS 原生
func (p *S3Provider) buildURL(key string) string {
	// 1. 优先使用自定义域名（CDN 等）
	if p.customDomain != "" {
		return strings.TrimRight(p.customDomain, "/") + "/" + key
	}

	// 2. 指定了自定义 Endpoint
	if p.endpoint != "" {
		base := strings.TrimRight(p.endpoint, "/")
		if p.forcePathStyle {
			// Path Style: endpoint/bucket/key（MinIO / Cloudflare R2 等）
			return base + "/" + p.bucket + "/" + key
		}
		// Virtual-Hosted Style: 将 bucket 注入 hostname
		// e.g. https://oss-cn-hz.aliyuncs.com → https://bucket.oss-cn-hz.aliyuncs.com/key
		if u, err := url.Parse(base); err == nil && u.Host != "" {
			u.Host = p.bucket + "." + u.Host
			u.Path = "/" + key
			return u.String()
		}
		// fallback：无法解析时退化为 Path Style
		return base + "/" + p.bucket + "/" + key
	}

	// 3. AWS S3 原生（无自定义 Endpoint），按 Region 构建标准 URL
	if p.region != "" {
		if p.forcePathStyle {
			return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", p.region, p.bucket, key)
		}
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, p.region, key)
	}

	return ""
}

// DetectContentType 按内容探测 MIME 类型（最多读取 512 字节）
func DetectContentType(data []byte) string {
	if len(data) > 512 {
		data = data[:512]
	}
	return http.DetectContentType(data)
}
