// Package entity MongoDB 文档映射结构体
//
// 职责：定义租户功能设置集合的文档结构
// 对外接口：TenantSettings, StorageSettings, AISettings, AIEndpoint
package entity

import "time"

// TenantSettings 租户功能设置，对应 MongoDB tenant_settings 集合
// _id 直接使用 tenant_id 字符串，一租户一条记录
type TenantSettings struct {
	TenantID  string           `bson:"_id"        json:"tenant_id"`
	Storage   *StorageSettings `bson:"storage"    json:"storage"`
	AI        *AISettings      `bson:"ai"         json:"ai"`
	Access    *AccessSettings  `bson:"access"     json:"access"` // 站点级访问控制（维护模式 / 画廊登录）
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_at"`
}

// AccessSettings 站点级访问控制配置
// 维护模式：整站仅登录可见；画廊登录：主题列表页需登录
type AccessSettings struct {
	MaintenanceMode      bool `bson:"maintenance_mode"       json:"maintenance_mode"`
	GalleryLoginRequired bool `bson:"gallery_login_required" json:"gallery_login_required"`
}

// StorageSettings 云存储 S3 配置（租户级）
type StorageSettings struct {
	Enabled         bool   `bson:"enabled"           json:"enabled"`
	Provider        string `bson:"provider"          json:"provider"` // aws / minio / aliyun / cloudflare / custom
	Endpoint        string `bson:"endpoint"          json:"endpoint"`
	Region          string `bson:"region"            json:"region"`
	Bucket          string `bson:"bucket"            json:"bucket"`
	AccessKeyID     string `bson:"access_key_id"     json:"access_key_id"`
	SecretAccessKey string `bson:"secret_access_key" json:"-"`                // 敏感字段，永不回传前端
	CustomDomain    string `bson:"custom_domain"     json:"custom_domain"`    // 可选 CDN/自定义域名
	DefaultTarget   string `bson:"default_target"    json:"default_target"`   // "local" | "cloud"
	ForcePathStyle  bool   `bson:"force_path_style"  json:"force_path_style"` // MinIO 等需要 path style
}

// AISettings AI 服务配置（租户级），预留用于后续 agent 扩展
type AISettings struct {
	Enabled   bool        `bson:"enabled"   json:"enabled"`
	Chat      *AIEndpoint `bson:"chat"      json:"chat"`
	Embedding *AIEndpoint `bson:"embedding" json:"embedding"`
	Reranker  *AIEndpoint `bson:"reranker"  json:"reranker"`
}

// AIEndpoint 单个 AI 服务端点（OpenAI 兼容接口）
type AIEndpoint struct {
	BaseURL string `bson:"base_url" json:"base_url"`
	APIKey  string `bson:"api_key"  json:"-"` // 敏感字段，永不回传前端
	ModelID string `bson:"model_id" json:"model_id"`
}
