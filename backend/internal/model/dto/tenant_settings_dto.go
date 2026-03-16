// Package dto 请求/响应数据传输对象
//
// 职责：定义租户设置相关的请求和响应结构体
// 对外接口：UpdateStorageReq, StorageUsageResp, UpdateAIReq, TenantSettingsResp
package dto

// UpdateStorageReq 更新云存储设置请求
type UpdateStorageReq struct {
	Enabled         bool   `json:"enabled"`
	Provider        string `json:"provider"`
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"` // 为空时保留原有值
	CustomDomain    string `json:"custom_domain"`
	DefaultTarget   string `json:"default_target" binding:"omitempty,oneof=local cloud"`
	ForcePathStyle  bool   `json:"force_path_style"`
}

// StorageSettingsResp 云存储设置响应（脱敏）
type StorageSettingsResp struct {
	Enabled        bool   `json:"enabled"`
	Provider       string `json:"provider"`
	Endpoint       string `json:"endpoint"`
	Region         string `json:"region"`
	Bucket         string `json:"bucket"`
	AccessKeyID    string `json:"access_key_id"`
	HasSecret      bool   `json:"has_secret"` // true 表示密钥已配置，但不回传值
	CustomDomain   string `json:"custom_domain"`
	DefaultTarget  string `json:"default_target"`
	ForcePathStyle bool   `json:"force_path_style"`
}

// StorageUsageResp 本地存储用量响应
type StorageUsageResp struct {
	UsedBytes  int64   `json:"used_bytes"`
	LimitBytes int64   `json:"limit_bytes"`
	UsedMB     float64 `json:"used_mb"`
	LimitMB    float64 `json:"limit_mb"`
	Percent    float64 `json:"percent"` // 0-100
}

// UpdateAIReq 更新 AI 设置请求
type UpdateAIReq struct {
	Enabled   bool           `json:"enabled"`
	Chat      *AIEndpointReq `json:"chat"`
	Embedding *AIEndpointReq `json:"embedding"`
	Reranker  *AIEndpointReq `json:"reranker"`
}

// AIEndpointReq 单个 AI 端点请求
type AIEndpointReq struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"` // 为空时保留原有值
	ModelID string `json:"model_id"`
}

// AISettingsResp AI 设置响应（脱敏）
type AISettingsResp struct {
	Enabled   bool            `json:"enabled"`
	Chat      *AIEndpointResp `json:"chat"`
	Embedding *AIEndpointResp `json:"embedding"`
	Reranker  *AIEndpointResp `json:"reranker"`
}

// AIEndpointResp 单个 AI 端点响应（脱敏）
type AIEndpointResp struct {
	BaseURL string `json:"base_url"`
	HasKey  bool   `json:"has_key"` // true 表示密钥已配置
	ModelID string `json:"model_id"`
}

// UpdateAccessReq 更新站点级访问控制请求
type UpdateAccessReq struct {
	MaintenanceMode      bool `json:"maintenance_mode"`
	GalleryLoginRequired bool `json:"gallery_login_required"`
}

// AccessSettingsResp 站点级访问控制响应
type AccessSettingsResp struct {
	MaintenanceMode      bool `json:"maintenance_mode"`
	GalleryLoginRequired bool `json:"gallery_login_required"`
}

// VerifyAccessCodeReq 验证主题访问码请求
type VerifyAccessCodeReq struct {
	Code string `json:"code" binding:"required"`
}

// VerifyAccessCodeResp 验证码校验成功响应
type VerifyAccessCodeResp struct {
	AccessToken string `json:"access_token"` // 短期 JWT（包含 theme_id + exp）
	ExpiresIn   int64  `json:"expires_in"`   // 秒
}

// TenantSettingsResp 完整设置响应
type TenantSettingsResp struct {
	Storage *StorageSettingsResp `json:"storage"`
	AI      *AISettingsResp      `json:"ai"`
	Access  *AccessSettingsResp  `json:"access"`
}
