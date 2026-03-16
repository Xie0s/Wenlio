/**
 * public-gallery.ts - 读者端主题画廊业务逻辑层
 *
 * 职责：封装主题画廊的公开 API 调用（分类/标签/主题）和 URL query params 双向同步
 * 对外暴露：usePublicGallery() composable 函数
 */
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { http, hasToken } from '@/utils/http'
import type { Category, Tag, Theme } from '@/utils/types'

export function usePublicGallery() {
  const route = useRoute()
  const router = useRouter()

  const tenantId = computed(() => route.params.tenantId as string)

  // ── 数据 ────────────────────────────────────────────────────
  const categories = ref<Category[]>([])
  const categoriesLoading = ref(false)
  const tags = ref<Tag[]>([])
  const tagsLoading = ref(false)
  const themes = ref<Theme[]>([])
  const loading = ref(false)

  // 画廊层可见主题：未登录时过滤 access_mode="login" 的主题
  const visibleThemes = computed(() => {
    if (hasToken()) return themes.value
    return themes.value.filter(t => t.access_mode !== 'login')
  })

  // ── URL 参数（筛选状态源） ─────────────────────────────────
  const selectedCategorySlug = computed(() => (route.query.category as string) || '')

  const selectedTagSlugs = computed<string[]>(() => {
    const t = route.query.tag
    if (!t) return []
    return Array.isArray(t) ? (t as string[]) : [t as string]
  })

  const hasFilters = computed(
    () => selectedCategorySlug.value !== '' || selectedTagSlugs.value.length > 0,
  )

  // ── 分类树扁平化（用于 slug ↔ id 互查） ─────────────────────
  function flattenCategories(list: Category[]): Category[] {
    const result: Category[] = []
    for (const cat of list) {
      result.push(cat)
      if (cat.children?.length) result.push(...flattenCategories(cat.children))
    }
    return result
  }

  // 当前选中分类的 ID（从 slug 反查，供 ThemeCard catLib 使用）
  const selectedCategoryId = computed(() => {
    if (!selectedCategorySlug.value) return ''
    return flattenCategories(categories.value).find(c => c.slug === selectedCategorySlug.value)?.id ?? ''
  })

  // 当前选中分类的名称（供侧边栏筛选徽章显示）
  const selectedCategoryName = computed(() => {
    if (!selectedCategorySlug.value) return ''
    return flattenCategories(categories.value).find(c => c.slug === selectedCategorySlug.value)?.name ?? ''
  })

  // ── 筛选操作 ────────────────────────────────────────────────
  function selectCategory(id: string) {
    const cat = flattenCategories(categories.value).find(c => c.id === id)
    const slug = cat?.slug ?? ''
    const isSame = selectedCategorySlug.value === slug
    router.replace({
      query: { ...route.query, category: isSame ? undefined : slug },
    })
  }

  function toggleTag(slug: string) {
    const current = [...selectedTagSlugs.value]
    const idx = current.indexOf(slug)
    if (idx >= 0) current.splice(idx, 1)
    else current.push(slug)
    router.replace({
      query: { ...route.query, tag: current.length ? current : undefined },
    })
  }

  function clearFilters() {
    router.replace({ query: {} })
  }

  // ── API 调用 ─────────────────────────────────────────────────
  // 请求版本号：只有最新请求的结果才更新 UI，防止竞态覆盖
  let themeRequestSeq = 0

  async function loadFilters() {
    categoriesLoading.value = true
    tagsLoading.value = true
    const [catRes, tagRes] = await Promise.all([
      http.get<Category[]>(`/public/tenants/${tenantId.value}/categories`),
      http.get<Tag[]>(`/public/tenants/${tenantId.value}/tags`),
    ])
    categoriesLoading.value = false
    tagsLoading.value = false
    if (catRes.code === 0 && catRes.data) categories.value = catRes.data
    if (tagRes.code === 0 && tagRes.data) tags.value = tagRes.data
  }

  async function loadThemes() {
    const seq = ++themeRequestSeq
    loading.value = true
    const params: Record<string, string | string[]> = {}
    if (selectedCategorySlug.value) params.category = selectedCategorySlug.value
    if (selectedTagSlugs.value.length) params.tag = selectedTagSlugs.value

    const url = hasFilters.value
      ? `/public/tenants/${tenantId.value}/themes/filter`
      : `/public/tenants/${tenantId.value}/themes`
    const res = await http.get<Theme[]>(url, hasFilters.value ? params : undefined)
    if (seq !== themeRequestSeq) return
    loading.value = false
    if (res.code === 0 && res.data) themes.value = res.data
  }

  // ── catLib/tagLib 适配（供 ThemeCard 展示分类名/标签名） ─────
  const catLib = { categories }
  const tagLib = { tags }

  return {
    tenantId,
    // 数据
    categories,
    categoriesLoading,
    tags,
    tagsLoading,
    themes,
    visibleThemes,
    loading,
    // 筛选状态
    selectedCategorySlug,
    selectedCategoryId,
    selectedCategoryName,
    selectedTagSlugs,
    hasFilters,
    // 操作
    selectCategory,
    toggleTag,
    clearFilters,
    // API
    loadFilters,
    loadThemes,
    // ThemeCard 适配器
    catLib,
    tagLib,
  }
}
