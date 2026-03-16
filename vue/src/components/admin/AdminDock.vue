<!--
  管理后台底部 Dock 导航栏组件

  功能：提供底部浮动导航栏，替代传统侧边栏
  职责：根据用户角色（超管/租户管理员）显示对应导航项，采用毛玻璃胶囊风格
  设计：底部居中浮动，图标按钮 + Tooltip 提示，当前路由高亮
        右侧集成主题切换 + 用户头像（点击展开用户中心下拉菜单）

  主要接口：
  - 无 props（独立组件，从 authStore 和 route 获取数据）
-->
<script setup lang="ts">
import { computed, markRaw, type Component } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import UserMenu from '@/components/auth/UserMenu.vue'
import {
  LayoutDashboard,
  Building2,
  Users,
  BookOpen,
  MessageSquare,
  ArrowUpRight,
  PanelTop,
  Settings2,
  Images,
} from 'lucide-vue-next'

interface NavItem {
  icon: Component
  label: string
  to: string
  exact?: boolean
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isSuperAdmin = computed(() => authStore.user?.role === 'super_admin')

const navItems = computed<NavItem[]>(() => {
  if (isSuperAdmin.value) {
    return [
      { icon: markRaw(LayoutDashboard), label: '仪表盘', to: '/admin', exact: true },
      { icon: markRaw(Building2), label: '租户管理', to: '/admin/tenants' },
      { icon: markRaw(Users), label: '用户管理', to: '/admin/users' },
    ]
  }
  return [
    { icon: markRaw(LayoutDashboard), label: '仪表盘', to: '/admin', exact: true },
    { icon: markRaw(BookOpen), label: '主题管理', to: '/admin/themes' },
    { icon: markRaw(MessageSquare), label: '评论管理', to: '/admin/comments' },
    { icon: markRaw(Users), label: '用户管理', to: '/admin/tenant-users' },
    { icon: markRaw(Images), label: '媒体管理', to: '/admin/media' },
    { icon: markRaw(PanelTop), label: '首页编辑', to: '/admin/homepage' },
    { icon: markRaw(Settings2), label: '租户设置', to: '/admin/settings' },
  ]
})

function isActive(item: NavItem): boolean {
  if (item.exact) return route.path === item.to
  return route.path.startsWith(item.to)
}

function goToTenantHome() {
  const tenantId = authStore.user?.tenant_id?.trim()
  if (tenantId) {
    window.open(`/${tenantId}/themes`, '_blank')
    return
  }
  window.open('/', '_blank')
}

</script>

<template>
  <div class="fixed bottom-4 left-1/2 -translate-x-1/2 z-40 pointer-events-none">
    <nav class="glass rounded-full px-2.5 py-2 flex items-center gap-1.5 pointer-events-auto">
      <!-- 导航项 -->
      <Tooltip v-for="item in navItems" :key="item.to">
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="h-12 w-12 rounded-full transition-colors"
            :class="isActive(item) ? 'bg-primary/15 text-primary' : 'text-foreground'" @click="router.push(item.to)">
            <component :is="item.icon" class="size-[22px]" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">{{ item.label }}</TooltipContent>
      </Tooltip>

      <!-- 分隔线 -->
      <div class="mx-1 h-6 w-px bg-border" />

      <!-- 前往文档集 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="h-12 w-12 rounded-full text-foreground transition-colors"
            @click="goToTenantHome">
            <ArrowUpRight class="size-[22px]" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">前往文档集</TooltipContent>
      </Tooltip>

      <!-- 主题切换 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <ThemeToggle button-size="size-12" class="text-foreground" icon-size="size-[22px]" />
        </TooltipTrigger>
        <TooltipContent side="top">切换主题</TooltipContent>
      </Tooltip>

      <!-- 用户头像 + 下拉菜单 -->
      <UserMenu side="top" align="end" :side-offset="12" trigger-class="h-12 w-12" avatar-class="h-10 w-10" />
    </nav>
  </div>
</template>
