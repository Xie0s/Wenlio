/**
 * utils/types.ts - 全局类型定义
 * 职责：定义与后端 API 对应的 TypeScript 类型，作为前端类型系统的唯一来源。
 * 对外暴露：所有业务实体类型和枚举常量
 */

// ============================================================
// 用户与认证
// ============================================================

export interface UserInfo {
  id: string
  tenant_id: string
  tenant_name: string
  username: string
  name: string
  role: 'super_admin' | 'tenant_admin'
  avatar_url: string
  bio: string
  profile_bg_url: string
}

export interface LoginReq {
  username: string
  password: string
  captcha_token: string
}

export interface RegisterTenantReq {
  tenant_id: string
  tenant_name: string
  admin_username: string
  admin_password: string
  admin_name: string
  captcha_token: string
}

export type CaptchaScene = 'login' | 'register' | 'reset_password'

export interface CaptchaChallengeReq {
  scene: CaptchaScene
}

export interface CaptchaChallengeResp {
  challenge_id: string
  scene: CaptchaScene
  mode: string
  prompt: string
  expires_at: number
  min_decision_ms: number
}

export interface CaptchaSignalSummary {
  dwell_ms: number
  visible_ms: number
  focused_ms: number
  visibility_changes: number
  focus_changes: number
  pointer_events: number
  key_events: number
  trusted_click: boolean
  language: string
  platform: string
  screen_width: number
  screen_height: number
  timezone_offset: number
  touch_points: number
  hardware_concurrency: number
  webdriver: boolean
}

export interface VerifyCaptchaReq {
  scene: CaptchaScene
  challenge_id: string
  duration_ms: number
  signals: CaptchaSignalSummary
}

export interface VerifyCaptchaResp {
  captcha_token: string
  expires_at: number
}

export interface LoginResp {
  access_token: string
  expires_in: number
  user: UserInfo
}

// ============================================================
// 租户
// ============================================================

export interface Tenant {
  id: string
  name: string
  logo_url: string
  browser_title?: string
  browser_icon_url?: string
  status: 'active' | 'suspended' | 'deleting'
  created_at: string
  updated_at: string
  user_count?: number
  admin_username?: string
  admin_name?: string
  admin_avatar_url?: string
}

export interface CreateTenantReq {
  id: string
  name: string
  logo_url: string
  admin_username: string
  admin_password: string
  admin_name: string
  admin_email?: string
}

export interface UpdateTenantReq {
  name: string
  logo_url: string
}

// ============================================================
// 主题
// ============================================================

export type AccessMode = 'public' | 'login' | 'code'

export interface Theme {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  sort_order: number
  category_id?: string
  tag_ids?: string[]
  created_at: string
  updated_at?: string
  version_count?: number
  section_count?: number
  page_count?: number
  current_version?: ThemeCurrentVersion | null
  access_mode?: AccessMode
}

export interface UpdateThemeReq {
  name: string
  slug: string
  description: string
  category_id?: string
  tag_ids?: string[]
  access_mode?: AccessMode
  access_code?: string
}

export interface ThemeCurrentVersion {
  id: string
  name: string
  label: string
  status: VersionStatus
  is_default: boolean
}

export interface CreateThemeReq {
  name: string
  slug: string
  description: string
  category_id?: string
  tag_ids?: string[]
}

// ============================================================
// 版本
// ============================================================

export type VersionStatus = 'draft' | 'published' | 'archived'

export interface Version {
  id: string
  tenant_id: string
  theme_id: string
  name: string
  label: string
  status: VersionStatus
  is_default: boolean
  published_at: string
  created_at: string
}

export interface CreateVersionReq {
  name: string
  label: string
}

// ============================================================
// 章节
// ============================================================

export interface Section {
  id: string
  tenant_id: string
  version_id: string
  title: string
  sort_order: number
}

// ============================================================
// 文档页
// ============================================================

export type PageStatus = 'draft' | 'published'

export interface DocPage {
  id: string
  tenant_id: string
  version_id: string
  section_id: string
  title: string
  slug: string
  content: string
  status: PageStatus
  sort_order: number
  published_at: string
  created_at: string
  updated_at: string
}

export interface CreatePageReq {
  title: string
  slug: string
  content: string
}

export interface UpdatePageReq {
  title: string
  slug: string
  content: string
  section_id?: string
}

// ============================================================
// 评论
// ============================================================

export type CommentStatus = 'pending' | 'approved' | 'rejected'

export interface Comment {
  id: string
  tenant_id: string
  page_id: string
  parent_id: string
  author: { name: string; email?: string }
  content: string
  status: CommentStatus
  created_at: string
}

// ============================================================
// 文档树（公开接口）
// ============================================================

export interface SectionTree {
  id: string
  title: string
  sort_order: number
  pages: PageMeta[]
}

export interface PageMeta {
  id: string
  title: string
  slug: string
  sort_order: number
}

// ============================================================
// 搜索
// ============================================================

export interface SearchResult {
  page_id: string
  title: string
  snippet: string
  theme_name: string
  version_name: string
  path: string
}

export interface SearchResponse {
  total: number
  items: SearchResult[]
}

// ============================================================
// 常量映射
// ============================================================

export const VERSION_STATUS_LABEL: Record<VersionStatus, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档',
}

export const PAGE_STATUS_LABEL: Record<PageStatus, string> = {
  draft: '草稿',
  published: '已发布',
}

export const COMMENT_STATUS_LABEL: Record<CommentStatus, string> = {
  pending: '待审核',
  approved: '已批准',
  rejected: '已拒绝',
}

export const VERSION_STATUS_COLOR: Record<VersionStatus, string> = {
  draft: 'bg-gray-100 text-gray-700',
  published: 'bg-green-100 text-green-700',
  archived: 'bg-amber-100 text-amber-700',
}

export const PAGE_STATUS_COLOR: Record<PageStatus, string> = {
  draft: 'bg-gray-100 text-gray-700',
  published: 'bg-green-100 text-green-700',
}

export const COMMENT_STATUS_COLOR: Record<CommentStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-700',
  approved: 'bg-green-100 text-green-700',
  rejected: 'bg-red-100 text-red-700',
}

export const ACCESS_MODE_LABEL: Record<AccessMode, string> = {
  public: '公开',
  login: '登录可见',
  code: '需验证码',
}

export const ACCESS_MODE_COLOR: Record<AccessMode, string> = {
  public: 'bg-green-100 text-green-700',
  login: 'bg-blue-100 text-blue-700',
  code: 'bg-amber-100 text-amber-700',
}

// ============================================================
// 租户首页个性化（详细类型定义见 components/personalization/types.ts）
// ============================================================

export type { TenantHomepageConfig, HomepageLayout } from '@/components/personalization/types'

// ============================================================
// 租户设置
// ============================================================

export interface StorageSettings {
  enabled: boolean
  provider: string
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  custom_domain: string
  default_target: 'local' | 'cloud'
  force_path_style: boolean
}

export interface UpdateStorageReq extends StorageSettings {
  secret_access_key?: string
}

export interface StorageUsage {
  used_bytes: number
  limit_bytes: number
  used_mb: number
  limit_mb: number
  percent: number
}

export interface AIEndpoint {
  base_url: string
  model_id: string
}

export interface AISettings {
  enabled: boolean
  chat?: AIEndpoint
  embedding?: AIEndpoint
  reranker?: AIEndpoint
}

export interface AccessSettings {
  maintenance_mode: boolean
  gallery_login_required: boolean
}

export interface TenantSettings {
  tenant_id: string
  storage?: StorageSettings
  ai?: AISettings
  access?: AccessSettings
  created_at: string
  updated_at: string
}

// ============================================================
// 分类
// ============================================================

export interface Category {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  parent_id: string
  sort_order: number
  level: number
  children?: Category[]
  theme_count?: number
  created_at: string
  updated_at: string
}

export interface CreateCategoryReq {
  name: string
  slug: string
  description: string
  parent_id?: string
}

export interface UpdateCategoryReq {
  name: string
  slug: string
  description: string
}

// ============================================================
// 标签
// ============================================================

export interface Tag {
  id: string
  tenant_id: string
  name: string
  slug: string
  description: string
  usage_count: number
  created_at: string
  updated_at: string
}

export interface CreateTagReq {
  name: string
  slug: string
  description: string
}

export interface UpdateTagReq {
  name: string
  slug: string
  description: string
}
