<!--
  管理后台顶部栏组件

  功能：显示页面标题与面包屑导航
  职责：提供管理后台的顶部导航，采用毛玻璃风格固定于页面顶部

  主要接口：
  - 无 props（独立组件，从 route.meta 获取数据）
-->
<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAdminBreadcrumbs } from '@/composables/useAdminBreadcrumbs'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/ui/breadcrumb'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

const props = defineProps<{
  titleOverride?: string
  activeTitle?: string
  hideUserInfo?: boolean
}>()

const router = useRouter()
const authStore = useAuthStore()

const displayName = computed(() => authStore.user?.name || authStore.user?.username || '用户')
const tenantName = computed(() => authStore.user?.tenant_name || '')

const { breadcrumbs, pageTitle } = useAdminBreadcrumbs({
  titleOverride: () => props.titleOverride,
  activeTitle: () => props.activeTitle,
})
</script>

<template>
  <header class="fixed top-0 left-0 right-0 z-30 h-16 flex items-center px-8 gap-4 bg-background/80 backdrop-blur-md">
    <!-- 左侧固定：页面标题 + 面包屑导航，max-w 防止长文本撑开导航栏 -->
    <div class="flex items-center gap-4 flex-none max-w-[50%] overflow-hidden">
      <Tooltip>
        <TooltipTrigger as-child>
          <h1 class="text-xl font-normal shrink-0 max-w-[180px] truncate cursor-default">{{ pageTitle }}</h1>
        </TooltipTrigger>
        <TooltipContent side="bottom">{{ pageTitle }}</TooltipContent>
      </Tooltip>

      <Breadcrumb v-if="breadcrumbs.length > 1" class="min-w-0 overflow-hidden">
        <BreadcrumbList class="flex-nowrap">
          <template v-for="(item, index) in breadcrumbs" :key="index">
            <BreadcrumbItem class="min-w-0">
              <Tooltip v-if="item.path">
                <TooltipTrigger as-child>
                  <BreadcrumbLink class="cursor-pointer inline-block max-w-[100px] truncate" @click.prevent="router.push(item.path!)">
                    {{ item.name }}
                  </BreadcrumbLink>
                </TooltipTrigger>
                <TooltipContent>{{ item.name }}</TooltipContent>
              </Tooltip>
              <Tooltip v-else>
                <TooltipTrigger as-child>
                  <BreadcrumbPage class="inline-block max-w-[130px] truncate">{{ item.name }}</BreadcrumbPage>
                </TooltipTrigger>
                <TooltipContent>{{ item.name }}</TooltipContent>
              </Tooltip>
            </BreadcrumbItem>
            <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" class="shrink-0" />
          </template>
        </BreadcrumbList>
      </Breadcrumb>
      <slot name="title-suffix" />
    </div>

    <!-- 中间弹性区：可注入工具栏内容（编辑器等页面使用） -->
    <div class="flex-1 flex items-center min-w-0">
      <slot name="toolbar" />
    </div>

    <!-- 右侧固定：可注入操作区 + 用户信息（始终可见） -->
    <div class="flex items-center gap-3 flex-none">
      <slot name="actions" />
      <span v-if="!hideUserInfo" class="text-base text-muted-foreground whitespace-nowrap">
        <template v-if="tenantName">{{ tenantName }} · </template>{{ displayName }}
      </span>
    </div>

  </header>
</template>
