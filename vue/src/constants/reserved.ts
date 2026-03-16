/**
 * constants/reserved.ts - 前端保留路径常量
 * 职责：与后端 pkg/constants/reserved.go 保持同步，防止保留路径被误识别为 tenantId
 * 对外暴露：RESERVED_TENANT_IDS
 */

// 注：favicon.ico / robots.txt 含特殊字符，无法通过 tenantId 正则，无需列入
export const RESERVED_TENANT_IDS = new Set([
  'admin',
  'api',
  'assets',
  'static',
  'health',
  'sitemap',
])
