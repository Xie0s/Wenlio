<!-- AdminLayout.vue - 管理后台布局组件
     职责：组合顶部导航栏、底部 Dock 导航和主内容区，提供管理后台的整体布局结构
     设计：采用顶部固定导航 + 底部浮动 Dock，移除传统侧边栏
     对外暴露：<AdminLayout /> 组件，由路由配置引用 -->
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import AdminDock from '@/components/admin/AdminDock.vue'

const route = useRoute()
const isFullWidth = computed(() => !!route.meta.fullWidth)
const isWideLayout = computed(() => !!route.meta.wideLayout)
const isSidebarLayout = computed(() => !!route.meta.sidebarLayout)
</script>

<template>
  <!-- fullWidth 页面（编辑器）自己完整管理布局，不渲染全局 Header/Dock -->
  <router-view v-if="isFullWidth" />

  <!-- 普通页面：使用全局 Header + Dock 布局 -->
  <div v-else class="min-h-svh bg-background">
    <AdminHeader />
    <template v-if="isWideLayout">
      <router-view />
    </template>
    <!-- sidebarLayout：固定视口高度，页面内部自管理双栏滚动 -->
    <main v-else-if="isSidebarLayout" class="pt-16 h-screen overflow-hidden">
      <div class="h-full">
        <router-view />
      </div>
    </main>

    <!-- 普通页面：整页自然滚动 -->
    <main v-else class="pt-16 pb-24 min-h-svh">
      <div class="mx-auto max-w-screen-2xl px-6 py-6">
        <router-view />
      </div>
    </main>
    <AdminDock />
  </div>
</template>
