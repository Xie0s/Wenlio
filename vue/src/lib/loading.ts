/*
 * loading.ts - 全局加载状态管理模块
 * 职责：
 * 1) 管理应用启动（页面刷新）与路由切换两类真实加载状态
 * 2) 对外提供统一状态读取与状态切换方法，供 App 与 router 复用
 * 对外暴露：
 * - useGlobalLoading(): 返回全局共享加载状态
 * - startAppBootLoading()/stopAppBootLoading(): 控制应用启动加载
 * - startRouteLoading()/stopRouteLoading(): 控制路由切换加载
 */

import { computed, readonly, ref } from 'vue'

const appBootLoading = ref(true)
const routeLoading = ref(false)

const isGlobalLoading = computed(() => appBootLoading.value || routeLoading.value)

export function useGlobalLoading() {
  return {
    appBootLoading: readonly(appBootLoading),
    routeLoading: readonly(routeLoading),
    isGlobalLoading,
  }
}

export function startAppBootLoading() {
  appBootLoading.value = true
}

export function stopAppBootLoading() {
  appBootLoading.value = false
}

export function startRouteLoading() {
  routeLoading.value = true
}

export function stopRouteLoading() {
  routeLoading.value = false
}
