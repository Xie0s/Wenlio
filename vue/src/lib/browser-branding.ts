/**
 * lib/browser-branding.ts - 浏览器标签品牌信息控制工具
 *
 * 职责：统一管理页面 title 与 favicon 的默认值、租户值应用与恢复逻辑。
 * 功能边界：只处理浏览器标签文本与图标，不参与路由、业务状态或页面内容渲染。
 * 对外暴露：DEFAULT_BROWSER_TITLE、DEFAULT_FAVICON_HREF、applyBrowserBranding、resetBrowserBranding、applyTenantBrowserBranding
 */

import { normalizeBrowserIconUrl, validateBrowserIconUrl } from '@/lib/validation'

export const DEFAULT_BROWSER_TITLE = '微讯云信息'
export const DEFAULT_FAVICON_HREF = '/favicon.ico'

function ensureFaviconLink(): HTMLLinkElement {
  const existing = document.querySelector("link[rel='icon']") as HTMLLinkElement | null
  if (existing) {
    return existing
  }

  const link = document.createElement('link')
  link.rel = 'icon'
  document.head.appendChild(link)
  return link
}

export function applyBrowserBranding(options?: {
  title?: string | null
  faviconHref?: string | null
}) {
  const title = options?.title?.trim() || DEFAULT_BROWSER_TITLE
  const rawFaviconHref = normalizeBrowserIconUrl(options?.faviconHref ?? '')
  const faviconHref = validateBrowserIconUrl(rawFaviconHref) ? DEFAULT_FAVICON_HREF : (rawFaviconHref || DEFAULT_FAVICON_HREF)

  document.title = title
  ensureFaviconLink().href = faviconHref
}

export function resetBrowserBranding() {
  applyBrowserBranding()
}

export function applyTenantBrowserBranding(options: {
  tenantName?: string | null
  browserTitle?: string | null
  browserIconUrl?: string | null
}) {
  const title = options.browserTitle?.trim() || options.tenantName?.trim() || DEFAULT_BROWSER_TITLE
  const faviconHref = options.browserIconUrl?.trim() || DEFAULT_FAVICON_HREF

  applyBrowserBranding({
    title,
    faviconHref,
  })
}
