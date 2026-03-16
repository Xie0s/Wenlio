/**
 * composables/personalization/useHomepageConfig.ts - 首页配置 API 对接
 *
 * 职责：封装首页配置的加载（公开/管理端）、保存草稿、发布操作
 * 对外暴露：useHomepageConfig() composable 函数
 */

import { ref } from 'vue'
import { http } from '@/utils/http'
import { toast } from 'vue-sonner'
import type {
  TenantHomepageConfig,
  HomepageLayout,
} from '@/components/personalization/types'

export function useHomepageConfig() {
  const config = ref<TenantHomepageConfig | null>(null)
  const loading = ref(false)
  const saving = ref(false)
  const publishing = ref(false)

  /**
   * 加载已发布首页配置（读者端）
   * GET /public/tenants/:tenantId/homepage
   */
  async function loadPublished(tenantId: string): Promise<HomepageLayout | null> {
    loading.value = true
    try {
      const res = await http.get<TenantHomepageConfig>(`/public/tenants/${tenantId}/homepage`)
      if (res.code === 0 && res.data) {
        config.value = res.data
        return res.data.published
      }
      return null
    } finally {
      loading.value = false
    }
  }

  /**
   * 加载草稿配置（管理端）
   * GET /tenant/homepage
   */
  async function loadDraft(): Promise<TenantHomepageConfig | null> {
    loading.value = true
    try {
      const res = await http.get<TenantHomepageConfig>('/tenant/homepage')
      if (res.code === 0 && res.data) {
        config.value = res.data
        return res.data
      }
      return null
    } finally {
      loading.value = false
    }
  }

  /**
   * 保存草稿
   * PUT /tenant/homepage
   */
  async function saveDraft(
    layout: HomepageLayout,
    options?: { silentSuccess?: boolean },
  ): Promise<boolean> {
    saving.value = true
    try {
      const res = await http.put('/tenant/homepage', layout)
      if (res.code === 0) {
        if (!options?.silentSuccess) {
          toast.success('草稿已保存')
        }
        return true
      }
      toast.error(res.message || '保存失败')
      return false
    } finally {
      saving.value = false
    }
  }

  /**
   * 发布首页（将 draft → published）
   * POST /tenant/homepage/publish
   */
  async function publish(): Promise<boolean> {
    publishing.value = true
    try {
      const res = await http.post('/tenant/homepage/publish')
      if (res.code === 0) {
        toast.success('首页已发布')
        return true
      }
      toast.error(res.message || '发布失败')
      return false
    } finally {
      publishing.value = false
    }
  }

  return {
    config,
    loading,
    saving,
    publishing,
    loadPublished,
    loadDraft,
    saveDraft,
    publish,
  }
}
