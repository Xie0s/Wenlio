/**
 * composables/useAdminBreadcrumbs.ts - 管理后台面包屑导航 composable
 * 职责：根据当前路由路径自动生成面包屑条目，支持 ObjectId 段标题覆盖与末尾追加动态叶节点标题
 * 对外暴露：
 *   - useAdminBreadcrumbs(options?) → { breadcrumbs, pageTitle }
 *   - BreadcrumbItemData（类型）
 *   - ADMIN_NAME_MAP（路由段→中文名映射）
 *   - isAdminObjectId（ObjectId 检测工具函数）
 */

import { computed, toValue, type MaybeRefOrGetter } from 'vue'
import { useRoute } from 'vue-router'

export interface BreadcrumbItemData {
  name: string
  path?: string
}

export const ADMIN_NAME_MAP: Record<string, string> = {
  admin: '管理后台',
  tenants: '租户管理',
  users: '用户管理',
  themes: '主题管理',
  versions: '版本管理',
  docs: '文档管理',
  editor: '编辑器',
  comments: '评论管理',
  'tenant-users': '用户管理',
}

export const isAdminObjectId = (str: string) => /^[0-9a-fA-F]{24}$/.test(str)

export function useAdminBreadcrumbs(options?: {
  titleOverride?: MaybeRefOrGetter<string | undefined>
  activeTitle?: MaybeRefOrGetter<string | undefined>
}) {
  const route = useRoute()

  const breadcrumbs = computed<BreadcrumbItemData[]>(() => {
    const titleOverride = toValue(options?.titleOverride)
    const activeTitle = toValue(options?.activeTitle)
    const hasActiveTitle = !!activeTitle

    const items: BreadcrumbItemData[] = [{ name: '仪表盘', path: '/admin' }]
    const pathSegments = route.path.split('/').filter(Boolean)

    let currentPath = ''
    for (let i = 1; i < pathSegments.length; i++) {
      const segment = pathSegments[i] as string
      currentPath = '/' + pathSegments.slice(0, i + 1).join('/')
      const isLast = i === pathSegments.length - 1

      if (isAdminObjectId(segment)) {
        if (isLast) items.push({ name: titleOverride || (route.meta.title as string) || segment })
        continue
      }

      const name = ADMIN_NAME_MAP[segment] || (route.meta.title as string) || segment
      items.push({
        name,
        path: isLast && !hasActiveTitle ? undefined : currentPath,
      })
    }

    if (hasActiveTitle) {
      items.push({ name: activeTitle! })
    }

    return items
  })

  const pageTitle = computed(() => {
    if (breadcrumbs.value.length <= 1) return '仪表盘'
    return breadcrumbs.value[breadcrumbs.value.length - 1]?.name || '管理后台'
  })

  return { breadcrumbs, pageTitle }
}
