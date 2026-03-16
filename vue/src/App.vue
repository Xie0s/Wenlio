<!-- App.vue - 应用根组件
     职责：全局布局容器、全局组件注入（Toaster / 版本更新提示）、认证状态恢复 -->
<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import LoadingIndicator from '@/components/common/LoadingIndicator.vue'
import VersionUpdatePrompt from '@/components/common/VersionUpdatePrompt.vue'
import { Toaster } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'
import { stopAppBootLoading, useGlobalLoading } from '@/lib/loading'
import { useAuthStore } from '@/stores/auth'
import { hasToken } from '@/utils/http'
import router from '@/router'
import { setupAriaHiddenGuard } from '@/lib/overlay-focus'

const { isGlobalLoading } = useGlobalLoading()
const authStore = useAuthStore()

let cleanupAriaGuard: (() => void) | undefined
onUnmounted(() => cleanupAriaGuard?.())

onMounted(async () => {
  cleanupAriaGuard = setupAriaHiddenGuard()
  try {
    // 页面刷新后恢复登录状态：Token 存在时重新获取用户信息
    if (hasToken()) {
      await authStore.fetchCurrentUser()
    } else {
      authStore.clearAuthState()
    }
    await router.isReady()
  } finally {
    stopAppBootLoading()
  }
})
</script>

<template>
  <TooltipProvider>
    <router-view />
    <LoadingIndicator v-if="isGlobalLoading" :is-running="isGlobalLoading" />
    <VersionUpdatePrompt />
    <Toaster position="bottom-right" :duration="3000" />
  </TooltipProvider>
</template>
