// Package service 业务逻辑层
//
// 职责：实现租户功能设置管理（云存储、AI 等），提供脱敏读取和安全更新
// 对外接口：TenantSettingsService
package service

import (
	"context"

	"docplatform/internal/model/dto"
	"docplatform/internal/model/entity"
	"docplatform/internal/repository"
	"docplatform/pkg/errcode"
	"docplatform/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
)

// TenantSettingsService 租户设置业务
type TenantSettingsService struct {
	settingsRepo *repository.TenantSettingsRepo
	tenantRepo   *repository.TenantRepo
	mediaRepo    *repository.MediaRepo
}

// NewTenantSettingsService 创建 TenantSettingsService
func NewTenantSettingsService() *TenantSettingsService {
	return &TenantSettingsService{
		settingsRepo: repository.NewTenantSettingsRepo(),
		tenantRepo:   repository.NewTenantRepo(),
		mediaRepo:    repository.NewMediaRepo(),
	}
}

// GetSettings 获取完整设置（脱敏响应）
func (s *TenantSettingsService) GetSettings(ctx context.Context, tenantID string) (*dto.TenantSettingsResp, *errcode.AppError) {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return nil, appErr
	}
	settings, err := s.settingsRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	resp := &dto.TenantSettingsResp{}
	if settings != nil {
		resp.Storage = toStorageResp(settings.Storage)
		resp.AI = toAIResp(settings.AI)
		resp.Access = toAccessResp(settings.Access)
	}
	return resp, nil
}

// UpdateStorage 更新云存储设置
func (s *TenantSettingsService) UpdateStorage(ctx context.Context, tenantID string, req *dto.UpdateStorageReq) *errcode.AppError {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return appErr
	}

	// 如果 SecretAccessKey 为空，保留原有值
	secret := req.SecretAccessKey
	if secret == "" {
		existing, err := s.settingsRepo.GetStorageConfig(ctx, tenantID)
		if err != nil {
			return errcode.ErrDatabase.Wrap(err)
		}
		if existing != nil {
			secret = existing.SecretAccessKey
		}
	}

	defaultTarget := req.DefaultTarget
	if defaultTarget == "" {
		defaultTarget = "local"
	}

	storage := entity.StorageSettings{
		Enabled:         req.Enabled,
		Provider:        req.Provider,
		Endpoint:        req.Endpoint,
		Region:          req.Region,
		Bucket:          req.Bucket,
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: secret,
		CustomDomain:    req.CustomDomain,
		DefaultTarget:   defaultTarget,
		ForcePathStyle:  req.ForcePathStyle,
	}

	return toAppError(s.settingsRepo.Upsert(ctx, tenantID, bson.M{"storage": storage}))
}

// GetStorageConfig 获取云存储配置（含密钥，内部使用）
func (s *TenantSettingsService) GetStorageConfig(ctx context.Context, tenantID string) (*entity.StorageSettings, *errcode.AppError) {
	cfg, err := s.settingsRepo.GetStorageConfig(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return cfg, nil
}

// GetStorageUsage 获取本地存储用量
func (s *TenantSettingsService) GetStorageUsage(ctx context.Context, tenantID string) (*dto.StorageUsageResp, *errcode.AppError) {
	const localStorageLimitBytes int64 = 100 * 1024 * 1024 // 100MB

	used, err := s.mediaRepo.SumFileSizeByTenant(ctx, tenantID, "local")
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}

	percent := float64(0)
	if localStorageLimitBytes > 0 {
		percent = float64(used) / float64(localStorageLimitBytes) * 100
		if percent > 100 {
			percent = 100
		}
	}

	return &dto.StorageUsageResp{
		UsedBytes:  used,
		LimitBytes: localStorageLimitBytes,
		UsedMB:     float64(used) / 1024 / 1024,
		LimitMB:    float64(localStorageLimitBytes) / 1024 / 1024,
		Percent:    percent,
	}, nil
}

// UpdateAI 更新 AI 设置
func (s *TenantSettingsService) UpdateAI(ctx context.Context, tenantID string, req *dto.UpdateAIReq) *errcode.AppError {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return appErr
	}

	// 获取已有 AI 配置（用于保留未提交的密钥）
	existingAI, err := s.settingsRepo.GetAIConfig(ctx, tenantID)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}

	ai := entity.AISettings{
		Enabled:   req.Enabled,
		Chat:      mergeAIEndpoint(req.Chat, getExistingEndpoint(existingAI, "chat")),
		Embedding: mergeAIEndpoint(req.Embedding, getExistingEndpoint(existingAI, "embedding")),
		Reranker:  mergeAIEndpoint(req.Reranker, getExistingEndpoint(existingAI, "reranker")),
	}

	return toAppError(s.settingsRepo.Upsert(ctx, tenantID, bson.M{"ai": ai}))
}

// TestS3Connection 测试 S3 连通性（使用已保存的配置发起真实 Ping）
func (s *TenantSettingsService) TestS3Connection(ctx context.Context, tenantID string) *errcode.AppError {
	cfg, err := s.settingsRepo.GetStorageConfig(ctx, tenantID)
	if err != nil {
		return errcode.ErrDatabase.Wrap(err)
	}
	if cfg == nil || !cfg.Enabled {
		return errcode.ErrCloudStorageNotEnabled
	}
	// Bucket / AK / SK 必须填写；AWS 原生（endpoint 为空）时 Region 也是必填
	if cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return errcode.ErrS3ConnectionFailed
	}
	if cfg.Endpoint == "" && cfg.Region == "" {
		return errcode.ErrS3ConnectionFailed
	}
	provider, err := storage.NewS3Provider(storage.S3Config{
		Endpoint:       cfg.Endpoint,
		Region:         cfg.Region,
		Bucket:         cfg.Bucket,
		AccessKeyID:    cfg.AccessKeyID,
		SecretKey:      cfg.SecretAccessKey,
		CustomDomain:   cfg.CustomDomain,
		ForcePathStyle: cfg.ForcePathStyle,
	})
	if err != nil {
		return errcode.ErrS3ConnectionFailed.Wrap(err)
	}
	if pingErr := provider.Ping(ctx); pingErr != nil {
		return errcode.ErrS3ConnectionFailed.Wrap(pingErr)
	}
	return nil
}

// TestAIConnection 测试 AI 连通性（预留，P5 阶段实现）
func (s *TenantSettingsService) TestAIConnection(ctx context.Context, tenantID string) *errcode.AppError {
	// TODO(P5): 调用 AI 端点的 /models 或 /chat/completions 接口验证连通性
	return nil
}

// UpdateAccess 更新站点级访问控制设置（维护模式 / 画廊登录可见）
func (s *TenantSettingsService) UpdateAccess(ctx context.Context, tenantID string, req *dto.UpdateAccessReq) *errcode.AppError {
	if appErr := (&TenantService{tenantRepo: s.tenantRepo}).EnsureActive(ctx, tenantID); appErr != nil {
		return appErr
	}
	access := entity.AccessSettings{
		MaintenanceMode:      req.MaintenanceMode,
		GalleryLoginRequired: req.GalleryLoginRequired,
	}
	return toAppError(s.settingsRepo.Upsert(ctx, tenantID, bson.M{"access": access}))
}

// GetAccessSettings 获取站点级访问控制（公开接口调用，不做租户活跃校验——由调用方负责）
func (s *TenantSettingsService) GetAccessSettings(ctx context.Context, tenantID string) (*entity.AccessSettings, *errcode.AppError) {
	settings, err := s.settingsRepo.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	if settings == nil || settings.Access == nil {
		return &entity.AccessSettings{}, nil
	}
	return settings.Access, nil
}

// GetAIConfig 获取 AI 配置（含密钥，内部使用）
func (s *TenantSettingsService) GetAIConfig(ctx context.Context, tenantID string) (*entity.AISettings, *errcode.AppError) {
	cfg, err := s.settingsRepo.GetAIConfig(ctx, tenantID)
	if err != nil {
		return nil, errcode.ErrDatabase.Wrap(err)
	}
	return cfg, nil
}

// ============================================================
// 内部辅助：脱敏转换
// ============================================================

func toAccessResp(a *entity.AccessSettings) *dto.AccessSettingsResp {
	if a == nil {
		return &dto.AccessSettingsResp{}
	}
	return &dto.AccessSettingsResp{
		MaintenanceMode:      a.MaintenanceMode,
		GalleryLoginRequired: a.GalleryLoginRequired,
	}
}

func toStorageResp(s *entity.StorageSettings) *dto.StorageSettingsResp {
	if s == nil {
		return &dto.StorageSettingsResp{DefaultTarget: "local"}
	}
	return &dto.StorageSettingsResp{
		Enabled:        s.Enabled,
		Provider:       s.Provider,
		Endpoint:       s.Endpoint,
		Region:         s.Region,
		Bucket:         s.Bucket,
		AccessKeyID:    s.AccessKeyID,
		HasSecret:      s.SecretAccessKey != "",
		CustomDomain:   s.CustomDomain,
		DefaultTarget:  s.DefaultTarget,
		ForcePathStyle: s.ForcePathStyle,
	}
}

func toAIResp(a *entity.AISettings) *dto.AISettingsResp {
	if a == nil {
		return &dto.AISettingsResp{}
	}
	return &dto.AISettingsResp{
		Enabled:   a.Enabled,
		Chat:      toAIEndpointResp(a.Chat),
		Embedding: toAIEndpointResp(a.Embedding),
		Reranker:  toAIEndpointResp(a.Reranker),
	}
}

func toAIEndpointResp(e *entity.AIEndpoint) *dto.AIEndpointResp {
	if e == nil {
		return nil
	}
	return &dto.AIEndpointResp{
		BaseURL: e.BaseURL,
		HasKey:  e.APIKey != "",
		ModelID: e.ModelID,
	}
}

func mergeAIEndpoint(req *dto.AIEndpointReq, existing *entity.AIEndpoint) *entity.AIEndpoint {
	if req == nil {
		return nil
	}
	apiKey := req.APIKey
	if apiKey == "" && existing != nil {
		apiKey = existing.APIKey
	}
	return &entity.AIEndpoint{
		BaseURL: req.BaseURL,
		APIKey:  apiKey,
		ModelID: req.ModelID,
	}
}

func getExistingEndpoint(ai *entity.AISettings, kind string) *entity.AIEndpoint {
	if ai == nil {
		return nil
	}
	switch kind {
	case "chat":
		return ai.Chat
	case "embedding":
		return ai.Embedding
	case "reranker":
		return ai.Reranker
	}
	return nil
}
