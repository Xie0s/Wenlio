/**
 * composables/useDocReader.ts - 文档阅读器路由导航 composable
 * 职责：封装阅读器页面的路由参数提取与页面跳转逻辑，供 DocReaderPage 及相关组件复用
 * 对外暴露：
 *   - useDocReader() → { tenantId, themeSlug, versionName, pageSlug,
 *                        navigateTo, navigateToPage, navigateToVersion, redirectTo }
 */

import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export function useDocReader() {
  const route = useRoute()
  const router = useRouter()

  const tenantId = computed(() => route.params.tenantId as string)
  const themeSlug = computed(() => route.params.themeSlug as string)
  const versionName = computed(() => route.params.version as string)
  const pageSlug = computed(() => (route.params.pageSlug as string) || '')

  function navigateTo(path: string) {
    router.push(path)
  }

  function navigateToPage(slug: string) {
    router.push(`/${tenantId.value}/${themeSlug.value}/${versionName.value}/${slug}`)
  }

  function navigateToVersion(newVersion: string) {
    router.push(`/${tenantId.value}/${themeSlug.value}/${newVersion}`)
  }

  function redirectTo(path: string) {
    router.replace(path)
  }

  return {
    tenantId,
    themeSlug,
    versionName,
    pageSlug,
    navigateTo,
    navigateToPage,
    navigateToVersion,
    redirectTo,
  }
}
