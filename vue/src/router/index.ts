/**
 * router/index.ts - 路由配置
 * 职责：定义管理后台和文档阅读界面的路由规则，含中间件链式导航守卫
 * 对外暴露：router 实例
 *
 * 中间件链（beforeEach 执行顺序）：
 *   1. 启动加载态
 *   2. 保留路径拦截（reserved IDs → 404）
 *   3. 登录鉴权（admin 路由）
 *   4. 角色守卫（super_admin 专属 / tenant_admin 专属双向隔离）
 */
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    /** 仅超管可访问 */
    role?: 'super_admin' | 'tenant_admin'
    /** 仅租户管理员可访问（super_admin 无 tenant_id，访问会导致 API 异常） */
    tenantOnly?: boolean
    title?: string
    fullWidth?: boolean
  }
}
import { startRouteLoading, stopRouteLoading } from '@/lib/loading'
import { resetBrowserBranding } from '@/lib/browser-branding'
import { useAuthStore } from '@/stores/auth'
import { useReaderStore } from '@/stores/reader'
import { RESERVED_TENANT_IDS } from '@/constants/reserved'
import { hasToken } from '@/utils/http'
import { toast } from 'vue-sonner'

// ── 路由参数正则 ──────────────────────────────────────────────
// tenantId：3+ 字符小写字母数字，允许中间连字符
const TENANT_RE = '[a-z0-9][a-z0-9-]+[a-z0-9]'
// slug（主题/页面）：1+ 字符小写字母数字，允许连字符
const SLUG_RE = '[a-z0-9][a-z0-9-]*'
// version：允许任何非路径分隔符字符（兼容含空格、点、连字符的版本号，如 "V 1.0"）
const VERSION_RE = '[^/]+'

/**
 * 解析文档路由：加载 主题 → 版本 → 文档树，返回第一篇页面的完整路径
 * ThemeHome 和 VersionHome 的 beforeEnter 共用此函数
 * @param versionName 未指定时使用默认版本
 */
async function resolveFirstPagePath(
  tenantId: string,
  themeSlug: string,
  versionName?: string,
): Promise<string> {
  const store = useReaderStore()

  // 进入前已确保当前租户 themes 已加载
  await store.ensureThemesLoaded(tenantId)
  const theme = store.findThemeBySlug(themeSlug)
  if (!theme) {
    toast.error('文档主题不存在')
    return `/${tenantId}`
  }

  // 主题级访问控制
  if (theme.access_mode === 'login' && !hasToken()) {
    return `/${tenantId}/${themeSlug}/login-required`
  }
  if (theme.access_mode === 'code' && !hasToken()) {
    const storedToken = localStorage.getItem(`theme_access_${theme.id}`)
    if (!storedToken) {
      return `/${tenantId}/${themeSlug}/verify`
    }
  }

  // 受保护主题已登录：签发 theme_access token（用于 raw 链接，避免泄露完整 JWT）
  if ((theme.access_mode === 'login' || theme.access_mode === 'code') && hasToken()) {
    await store.ensureThemeAccessToken(theme.id)
  }

  await store.loadVersions(theme.id)
  const version = versionName
    ? store.findVersionByName(versionName)
    : store.getDefaultVersion()
  if (!version) {
    toast.error(versionName ? '版本不存在' : '该主题暂无可用版本')
    return `/${tenantId}`
  }

  await store.loadTree(version.id)
  const firstSlug = store.getFirstPageSlug()
  if (!firstSlug) {
    toast.error('该版本暂无文档')
    return `/${tenantId}`
  }

  return `/${tenantId}/${themeSlug}/${version.name}/${firstSlug}`
}

const routes: RouteRecordRaw[] = [
  // ============================================================
  // 首页（精确匹配，必须在动态路由之前）
  // ============================================================
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home.vue'),
  },

  // 组件演示页：所有环境可访问
  {
    path: '/demo',
    name: 'Demo',
    component: () => import('@/views/demo.vue'),
  },

  // ============================================================
  // 管理后台路由
  // ============================================================
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/auth/AuthPage.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/admin/register',
    name: 'AdminRegister',
    component: () => import('@/views/auth/AuthPage.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/admin',
    component: () => import('@/components/admin/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/DashboardPage.vue'),
        meta: { title: '仪表盘' },
      },
      // 超管专属：租户管理
      {
        path: 'tenants',
        name: 'TenantList',
        component: () => import('@/views/admin/TenantListPage.vue'),
        meta: { role: 'super_admin', title: '租户管理' },
      },
      // 超管专属：用户管理
      {
        path: 'users',
        name: 'AdminUserList',
        component: () => import('@/views/admin/UserListPage.vue'),
        meta: { role: 'super_admin', title: '用户管理' },
      },
      // 租户管理员专属：主题列表
      {
        path: 'themes',
        name: 'ThemeList',
        component: () => import('@/views/admin/ThemeListPage.vue'),
        meta: { tenantOnly: true, title: '主题管理', sidebarLayout: true },
      },
      // 租户管理员专属：主题编辑器（版本/章节/文档/编辑器合一）
      {
        path: 'themes/:themeId',
        name: 'ThemeEditor',
        component: () => import('@/views/admin/ThemeEditorPage.vue'),
        meta: { tenantOnly: true, title: '主题编辑', fullWidth: true },
      },
      // 租户管理员专属：评论管理
      {
        path: 'comments',
        name: 'CommentList',
        component: () => import('@/views/admin/CommentListPage.vue'),
        meta: { tenantOnly: true, title: '评论管理' },
      },
      // 租户管理员专属：首页编辑器
      {
        path: 'homepage',
        name: 'HomepageEditor',
        component: () => import('@/views/admin/HomepageEditorPage.vue'),
        meta: { tenantOnly: true, title: '首页编辑', fullWidth: true },
      },
      // 租户管理员专属：租户用户管理
      {
        path: 'tenant-users',
        name: 'TenantUserList',
        component: () => import('@/views/admin/TenantUserListPage.vue'),
        meta: { tenantOnly: true, title: '用户管理' },
      },
      // 租户管理员专属：媒体文件管理
      {
        path: 'media',
        name: 'MediaManage',
        component: () => import('@/views/admin/MediaPage.vue'),
        meta: { tenantOnly: true, title: '媒体管理', sidebarLayout: true },
      },
      // 租户管理员专属：租户设置
      {
        path: 'settings',
        name: 'TenantSettings',
        component: () => import('@/views/admin/TenantSettingsPage.vue'),
        meta: { tenantOnly: true, title: '租户设置' },
      },
      // 通用：个人中心
      {
        path: 'profile',
        name: 'AdminProfile',
        component: () => import('@/views/admin/ProfilePage.vue'),
        meta: { title: '个人中心' },
      },
    ],
  },

  // ============================================================
  // 文档阅读路由（公开，无需登录）
  // tenantId 正则排除系统保留路径；themeSlug/version/pageSlug
  // 添加约束防止特殊字符注入，保留路径拦截在 beforeEach 二次把关
  // ============================================================
  {
    path: `/:tenantId(${TENANT_RE})`,
    name: 'TenantHome',
    component: () => import('@/views/reader/TenantHomePage.vue'),
  },
  // 主题画廊：支持分类/标签筛选（公开）
  {
    path: `/:tenantId(${TENANT_RE})/themes`,
    name: 'ThemeGallery',
    component: () => import('@/views/reader/ThemeGalleryPage.vue'),
  },
  // 维护模式页面
  {
    path: `/:tenantId(${TENANT_RE})/maintenance`,
    name: 'Maintenance',
    component: () => import('@/views/reader/MaintenancePage.vue'),
  },
  // 主题验证码页面
  {
    path: `/:tenantId(${TENANT_RE})/:themeSlug(${SLUG_RE})/verify`,
    name: 'ThemeVerifyCode',
    component: () => import('@/views/reader/ThemeVerifyCodePage.vue'),
  },
  // 主题登录提示页面
  {
    path: `/:tenantId(${TENANT_RE})/:themeSlug(${SLUG_RE})/login-required`,
    name: 'ThemeLoginRequired',
    component: () => import('@/views/reader/ThemeLoginRequiredPage.vue'),
  },
  {
    // ThemeHome：beforeEnter 始终重定向到第一篇页面，component 仅为满足 Vue Router API
    path: `/:tenantId(${TENANT_RE})/:themeSlug(${SLUG_RE})`,
    name: 'ThemeHome',
    component: () => import('@/views/reader/TenantHomePage.vue'),
    beforeEnter: async (to) => {
      try {
        return await resolveFirstPagePath(
          to.params.tenantId as string,
          to.params.themeSlug as string,
        )
      } catch {
        toast.error('加载失败，请重试')
        return `/${to.params.tenantId}`
      }
    },
  },
  {
    // VersionHome：beforeEnter 始终重定向到第一篇页面（版本切换 / 直接访问均触发）
    path: `/:tenantId(${TENANT_RE})/:themeSlug(${SLUG_RE})/:version(${VERSION_RE})`,
    name: 'VersionHome',
    component: () => import('@/views/reader/DocReaderPage.vue'),
    beforeEnter: async (to) => {
      try {
        return await resolveFirstPagePath(
          to.params.tenantId as string,
          to.params.themeSlug as string,
          to.params.version as string,
        )
      } catch {
        toast.error('加载失败，请重试')
        return `/${to.params.tenantId}`
      }
    },
  },
  {
    path: `/:tenantId(${TENANT_RE})/:themeSlug(${SLUG_RE})/:version(${VERSION_RE})/:pageSlug(${SLUG_RE})`,
    name: 'DocPageView',
    component: () => import('@/views/reader/DocReaderPage.vue'),
  },

  // ============================================================
  // 404 兜底（必须放在最后）
  // ============================================================
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundPage.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// ── 中间件链 ──────────────────────────────────────────────────
router.beforeEach(async (to) => {
  startRouteLoading()

  // ① 保留路径拦截：tenantId 命中后端保留词时导向 404
  const tenantId = to.params.tenantId as string | undefined
  if (tenantId && RESERVED_TENANT_IDS.has(tenantId)) {
    return { name: 'NotFound' }
  }

  // ② 读者端访问控制：维护模式 / 画廊登录拦截
  if (tenantId && !to.path.startsWith('/admin') && to.name !== 'Maintenance' && to.name !== 'ThemeVerifyCode' && to.name !== 'ThemeLoginRequired') {
    const readerStore = useReaderStore()
    if (!readerStore.accessSettings) {
      await readerStore.loadAccessSettings(tenantId)
    }
    const access = readerStore.accessSettings
    if (access) {
      const isLoggedIn = hasToken()
      // 维护模式：未登录重定向到维护页
      if (access.maintenance_mode && !isLoggedIn) {
        return `/${tenantId}/maintenance`
      }
      // 画廊登录可见：未登录访问主题列表时重定向登录
      if (access.gallery_login_required && !isLoggedIn && to.name === 'ThemeGallery') {
        toast.info('该页面需要登录后访问')
        return '/admin/login'
      }
    }

    // 主题级访问控制：针对带有 themeSlug 的路由（DocPageView 等直接访问场景）
    const themeSlug = to.params.themeSlug as string | undefined
    if (themeSlug && to.name !== 'ThemeHome' && to.name !== 'VersionHome') {
      await readerStore.ensureThemesLoaded(tenantId)
      const theme = readerStore.themes.find(t => t.slug === themeSlug)
      if (theme) {
        if (theme.access_mode === 'login' && !hasToken()) {
          return `/${tenantId}/${themeSlug}/login-required`
        }
        if (theme.access_mode === 'code' && !hasToken()) {
          const storedToken = localStorage.getItem(`theme_access_${theme.id}`)
          if (!storedToken) {
            return `/${tenantId}/${themeSlug}/verify`
          }
        }
      }
    }
  }

  // ③ 登录鉴权：/admin/* 路由（login 页自身除外）
  if (to.meta.requiresAuth !== false && to.path.startsWith('/admin')) {
    const authStore = useAuthStore()
    if (!hasToken()) {
      authStore.clearAuthState()
      return '/admin/login'
    }

    if (!authStore.user) {
      const ok = await authStore.fetchCurrentUser()
      if (!ok) {
        return '/admin/login'
      }
    }

    const role = authStore.user?.role

    // ④ 角色守卫：super_admin 专属路由
    if (to.meta.role && role !== to.meta.role) {
      return '/admin'
    }

    // ⑤ 角色守卫（反向）：tenant_admin 专属路由禁止 super_admin 访问
    //    super_admin 无 tenant_id，访问租户专属功能会导致 API 异常
    if (to.meta.tenantOnly && role === 'super_admin') {
      return '/admin'
    }
  }
})

router.afterEach(() => {
  stopRouteLoading()
})

router.afterEach((to) => {
  if (!to.params.tenantId) {
    resetBrowserBranding()
  }
})

router.onError(() => {
  stopRouteLoading()
})

export default router
