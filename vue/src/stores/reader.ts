/**
 * stores/reader.ts - 文档阅读数据缓存 store
 * 职责：
 * 1) 管理读者端的租户、主题、版本、文档树、页面内容、评论数据
 * 2) 提供客户端缓存，相同页面不重复请求
 * 3) 对外暴露统一的数据加载方法供 ReaderLayout 和 DocPage 消费
 * 对外暴露：useReaderStore
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { http, hasToken } from '@/utils/http'
import type {
  Tenant, Theme, Version, SectionTree, DocPage, Comment,
  HomepageLayout, AccessSettings,
} from '@/utils/types'

export const useReaderStore = defineStore('reader', () => {
  // 核心数据
  const tenant = ref<Tenant | null>(null)
  const themes = ref<Theme[]>([])
  const loadedThemesTenantId = ref('')
  const currentTheme = ref<Theme | null>(null)
  const versions = ref<Version[]>([])
  const currentVersion = ref<Version | null>(null)
  const tree = ref<SectionTree[]>([])
  const currentPage = ref<DocPage | null>(null)
  const comments = ref<Comment[]>([])

  // 页面内容缓存 Map<pageId, DocPage>
  const pageCache = new Map<string, DocPage>()

  // 个性化首页
  const homepageLayout = ref<HomepageLayout | null>(null)
  const loadingHomepage = ref(false)

  // 站点级访问控制
  const accessSettings = ref<AccessSettings | null>(null)

  // 加载状态
  const loadingTenant = ref(false)
  const loadingTree = ref(false)
  const loadingPage = ref(false)
  const loadingComments = ref(false)

  // 派生状态
  const tenantName = computed(() => tenant.value?.name || '文档平台')

  // 加载租户信息
  // 注：tenantId 为 URL slug，不能与 tenant.value.id（ObjectID）比较，此处无缓存跳过
  async function loadTenant(tenantId: string) {
    loadingTenant.value = true
    tenant.value = null
    try {
      const res = await http.get<Tenant>(`/public/tenants/${tenantId}`)
      if (res.code === 0 && res.data) {
        tenant.value = res.data
      } else {
        // 租户不存在（404 等），主动置 null，让调用方可通过 store.tenant 判断
        tenant.value = null
      }
    } finally {
      loadingTenant.value = false
    }
  }

  // 加载租户首页配置（已发布）
  async function loadHomepage(tenantId: string) {
    loadingHomepage.value = true
    homepageLayout.value = null
    try {
      const res = await http.get<{ published: HomepageLayout | null }>(`/public/tenants/${tenantId}/homepage`)
      if (res.code === 0 && res.data) {
        homepageLayout.value = res.data.published
      } else {
        homepageLayout.value = null
      }
    } finally {
      loadingHomepage.value = false
    }
  }

  // 加载站点级访问控制设置
  async function loadAccessSettings(tenantId: string) {
    const res = await http.get<AccessSettings>(`/public/tenants/${tenantId}/access`)
    if (res.code === 0 && res.data) {
      accessSettings.value = res.data
    }
  }

  // 加载主题列表
  async function loadThemes(tenantId: string) {
    themes.value = []
    currentTheme.value = null
    loadedThemesTenantId.value = ''
    const res = await http.get<Theme[]>(`/public/tenants/${tenantId}/themes`)
    if (res.code === 0 && res.data) {
      themes.value = res.data
    }
    loadedThemesTenantId.value = tenantId
  }

  function hasThemesForTenant(tenantId: string): boolean {
    return loadedThemesTenantId.value === tenantId
  }

  async function ensureThemesLoaded(tenantId: string) {
    if (hasThemesForTenant(tenantId)) return
    await loadThemes(tenantId)
  }

  // 根据 slug 查找主题并设为当前
  function findThemeBySlug(slug: string): Theme | null {
    const found = themes.value.find(t => t.slug === slug) || null
    currentTheme.value = found
    return found
  }

  // 获取当前主题的访问 token 参数（code/login 主题附带 theme_token）
  function getThemeAccessParams(): Record<string, string> | undefined {
    const theme = currentTheme.value
    if (theme?.access_mode === 'code' || theme?.access_mode === 'login') {
      const token = localStorage.getItem(`theme_access_${theme.id}`)
      if (token) return { theme_token: token }
    }
    return undefined
  }

  // 为 login 类型主题签发并缓存 theme_access token（已登录时调用）
  async function ensureThemeAccessToken(themeId: string): Promise<void> {
    const existing = localStorage.getItem(`theme_access_${themeId}`)
    if (existing) return
    if (!hasToken()) return
    const res = await http.post<{ access_token: string }>(`/public/themes/${themeId}/issue-token`)
    if (res.code === 0 && res.data?.access_token) {
      localStorage.setItem(`theme_access_${themeId}`, res.data.access_token)
    }
  }

  // 加载版本列表
  async function loadVersions(themeId: string) {
    versions.value = []
    const res = await http.get<Version[]>(`/public/themes/${themeId}/versions`, getThemeAccessParams())
    if (res.code === 0 && res.data) {
      versions.value = res.data
    }
  }

  // 根据版本名找到版本并设为当前
  function findVersionByName(name: string): Version | null {
    const found = versions.value.find(v => v.name === name) || null
    currentVersion.value = found
    return found
  }

  // 加载文档树
  async function loadTree(versionId: string) {
    loadingTree.value = true
    tree.value = []
    try {
      const res = await http.get<{ sections: SectionTree[] }>(`/public/versions/${versionId}/tree`, getThemeAccessParams())
      if (res.code === 0 && res.data) {
        tree.value = res.data.sections || []
      }
    } finally {
      loadingTree.value = false
    }
  }

  // 从树中按 slug 查找页面 ID
  function findPageIdBySlug(slug: string): string | null {
    for (const section of tree.value) {
      const found = section.pages.find(p => p.slug === slug)
      if (found) return found.id
    }
    return null
  }

  // 加载文档页内容（带缓存）
  async function loadPage(pageId: string) {
    const cached = pageCache.get(pageId)
    if (cached) {
      currentPage.value = cached
      return
    }
    loadingPage.value = true
    currentPage.value = null
    try {
      const res = await http.get<DocPage>(`/public/pages/${pageId}`, getThemeAccessParams())
      if (res.code === 0 && res.data) {
        pageCache.set(pageId, res.data)
        currentPage.value = res.data
      }
    } finally {
      loadingPage.value = false
    }
  }

  // 加载评论列表
  async function loadComments(pageId: string) {
    loadingComments.value = true
    comments.value = []
    try {
      const res = await http.get<Comment[]>(`/public/pages/${pageId}/comments`, getThemeAccessParams())
      if (res.code === 0 && res.data) {
        comments.value = res.data
      } else {
        comments.value = []
      }
    } finally {
      loadingComments.value = false
    }
  }

  // 提交评论
  async function submitComment(
    pageId: string,
    author: { name: string; email?: string },
    content: string,
  ): Promise<{ success: boolean; message?: string }> {
    const res = await http.post(`/public/pages/${pageId}/comments`, { author, content })
    if (res.code === 0) {
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  // 清除缓存（版本切换时调用）
  function clearCache() {
    pageCache.clear()
    currentPage.value = null
    tree.value = []
    comments.value = []
  }

  // 清除当前上下文（跳回租户首页时调用）
  function clearContext() {
    currentTheme.value = null
    currentVersion.value = null
    versions.value = []
    clearCache()
  }

  // 获取第一篇文档的 slug
  function getFirstPageSlug(): string | null {
    for (const section of tree.value) {
      const first = section.pages[0]
      if (first) return first.slug
    }
    return null
  }

  // 获取默认版本
  function getDefaultVersion(): Version | null {
    return versions.value.find(v => v.is_default) || versions.value[0] || null
  }

  return {
    tenant,
    themes,
    homepageLayout,
    loadingHomepage,
    currentTheme,
    versions,
    currentVersion,
    tree,
    currentPage,
    comments,
    loadingTenant,
    loadingTree,
    loadingPage,
    loadingComments,
    tenantName,
    accessSettings,
    loadAccessSettings,
    loadTenant,
    loadHomepage,
    loadThemes,
    hasThemesForTenant,
    ensureThemesLoaded,
    findThemeBySlug,
    loadVersions,
    findVersionByName,
    loadTree,
    findPageIdBySlug,
    loadPage,
    loadComments,
    submitComment,
    clearCache,
    clearContext,
    getFirstPageSlug,
    getDefaultVersion,
    ensureThemeAccessToken,
  }
})
