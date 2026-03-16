/**
 * theme-list.ts - 主题列表业务逻辑层
 *
 * 职责：封装主题列表的状态管理与 API 操作（加载、创建、编辑、删除）
 * 对外暴露：useThemeList() composable 函数
 *
 * 注意：表单状态（form、showCreate 等）已下沉到组件层管理，
 *       createTheme/updateTheme 直接接受数据参数。
 */
import { ref } from 'vue'
import { http, type ApiResponse } from '@/utils/http'
import type { Theme, CreateThemeReq, UpdateThemeReq } from '@/utils/types'
import { toast } from 'vue-sonner'

export function useThemeList() {
  const themes = ref<Theme[]>([])
  const loading = ref(false)

  // 记住上次筛选参数，CRUD 后用相同参数重新加载，避免丢失筛选视图
  let lastCategoryId: string | undefined
  let lastTagIds: string[] | undefined
  // 请求版本号：只有最新请求的结果才更新 UI，防止竞态覆盖
  let requestSeq = 0

  async function loadThemes(categoryId?: string, tagIds?: string[]) {
    lastCategoryId = categoryId
    lastTagIds = tagIds
    const seq = ++requestSeq
    loading.value = true
    const params = new URLSearchParams()
    if (categoryId) params.set('category_id', categoryId)
    if (tagIds?.length) tagIds.forEach(id => params.append('tag_id', id))
    const qs = params.toString()
    const url = qs ? `/tenant/themes?${qs}` : '/tenant/themes'
    const res = await http.get<Theme[]>(url)
    if (seq !== requestSeq) return
    loading.value = false
    if (res.code === 0 && res.data) {
      themes.value = res.data
    }
  }

  function reloadThemes() {
    return loadThemes(lastCategoryId, lastTagIds)
  }

  async function createTheme(data: CreateThemeReq): Promise<boolean> {
    const res = await http.post<Theme>('/tenant/themes', data)
    if (res.code === 0 && res.data) {
      // 自动创建初始草稿版本，确保进入编辑器时版本已存在
      await http.post(`/tenant/themes/${res.data.id}/versions`, { name: 'v1.0', label: '' })
      toast.success('主题创建成功')
      reloadThemes()
      return true
    }
    toast.error(res.message)
    return false
  }

  async function deleteTheme(id: string, opts?: { silentError?: boolean }): Promise<ApiResponse<null>> {
    const res = await http.delete(`/tenant/themes/${id}`)
    if (res.code === 0) {
      toast.success('已删除')
      await reloadThemes()
    } else if (!opts?.silentError) {
      toast.error(res.message)
    }
    return res as ApiResponse<null>
  }

  async function deleteThemeCascade(id: string): Promise<ApiResponse<null>> {
    const res = await http.delete(`/tenant/themes/${id}?cascade=true`)
    if (res.code === 0) {
      toast.success('已删除主题及其版本')
      await reloadThemes()
    } else {
      toast.error(res.message)
    }
    return res as ApiResponse<null>
  }

  async function updateTheme(id: string, data: UpdateThemeReq): Promise<boolean> {
    const res = await http.patch(`/tenant/themes/${id}`, data)
    if (res.code === 0) {
      toast.success('主题已更新')
      reloadThemes()
      return true
    }
    toast.error(res.message)
    return false
  }

  function formatDate(dateStr: string) {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }

  return {
    themes,
    loading,
    loadThemes,
    createTheme,
    deleteTheme,
    deleteThemeCascade,
    updateTheme,
    formatDate,
  }
}
