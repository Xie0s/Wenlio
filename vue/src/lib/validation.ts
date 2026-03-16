/**
 * lib/validation.ts - 路由参数校验工具
 * 职责：与 router/index.ts 中的路由正则保持同步，提供前端表单校验函数和提示文案
 * 对外暴露：SLUG_PATTERN, VERSION_PATTERN, validateSlug, validateVersionName, slugHint, versionHint
 */

/** 主题 slug / 页面 slug：小写字母数字，允许中间连字符，至少 1 字符 */
export const SLUG_PATTERN = /^[a-z0-9][a-z0-9-]*$/

/** 版本名：不含路径分隔符 / */
export const VERSION_PATTERN = /^[^/]+$/

export const slugHint = '仅允许小写字母、数字和连字符（-），以字母或数字开头，如 quick-start'
export const versionHint = '不能包含 /，推荐格式如 v1.0、V 1.0、latest'

export function validateSlug(value: string): string {
  if (!value) return '必填'
  if (!SLUG_PATTERN.test(value)) return slugHint
  return ''
}

export function validateVersionName(value: string): string {
  if (!value) return '必填'
  if (!VERSION_PATTERN.test(value)) return versionHint
  return ''
}

export function normalizeBrowserIconUrl(value: string): string {
  const trimmed = value.trim()
  if (!trimmed) return ''
  if (/^(https?:\/\/|\/|data:image\/|blob:)/i.test(trimmed)) return trimmed
  if (trimmed.startsWith('www.')) return `https://${trimmed}`
  return trimmed
}

export function validateBrowserIconUrl(value: string): string {
  const normalized = normalizeBrowserIconUrl(value)
  if (!normalized) return ''
  if (/^(\/|data:image\/|blob:)/i.test(normalized)) return ''
  try {
    const url = new URL(normalized)
    if (url.protocol === 'http:' || url.protocol === 'https:') return ''
  } catch {
    return '请输入有效的图标地址，支持 https://、/path、data:image 或 blob: 地址'
  }
  return '请输入有效的图标地址，支持 https://、/path、data:image 或 blob: 地址'
}
