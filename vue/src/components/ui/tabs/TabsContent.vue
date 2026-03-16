<script setup lang="ts">
/**
 * iOS 胶囊风格 TabsContent 组件
 * 提供内容切换的左滑/右滑过渡动画
 */
import type { TabsContentProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { ref, nextTick, onMounted, onUnmounted, inject } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsContent } from "reka-ui"
import { cn } from "@/utils"
import { TABS_SLIDE_DIRECTION_KEY } from "./keys"

const props = defineProps<TabsContentProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")

// 从父级 Tabs 注入滑动方向
const slideDirection = inject(TABS_SLIDE_DIRECTION_KEY, ref<'left' | 'right'>('right'))

// 动画状态
const isAnimating = ref(false)
const animationDirection = ref<'left' | 'right'>('right')
const contentRef = ref<HTMLElement | null>(null)
const hasActivated = ref(false)
let observer: MutationObserver | null = null

// 监听 data-state 变化来触发动画
onMounted(() => {
  nextTick(() => {
    if (contentRef.value) {
      // 获取实际 DOM 元素
      const el = (contentRef.value as any)?.$el || contentRef.value
      
      // 使用 MutationObserver 监听 data-state 变化
      observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          if (mutation.attributeName === 'data-state') {
            const target = mutation.target as HTMLElement
            const newState = target.getAttribute('data-state')
            
            if (newState === 'active') {
              // 首次激活只做稳定渲染，不执行位移动画，避免初次进入造成滚动抖动
              if (!hasActivated.value) {
                hasActivated.value = true
                isAnimating.value = false
                return
              }

              // 使用从父级注入的滑动方向
              animationDirection.value = slideDirection.value
              
              // 触发动画
              isAnimating.value = false
              nextTick(() => {
                isAnimating.value = true
              })
            }
          }
        })
      })
      
      observer.observe(el, { attributes: true, attributeFilter: ['data-state'] })
    }
  })
})

onUnmounted(() => {
  observer?.disconnect()
  observer = null
})
</script>

<template>
  <TabsContent
    ref="contentRef"
    data-slot="tabs-content"
    :class="cn(
      'tabs-content-base',
      isAnimating && (animationDirection === 'left' ? 'tabs-slide-from-left' : 'tabs-slide-from-right'),
      props.class
    )"
    v-bind="delegatedProps"
  >
    <slot />
  </TabsContent>
</template>

<style>
/* 基础样式 - 非 scoped 确保正确应用 */
.tabs-content-base {
  outline: none;
  width: 100%;
  overflow-x: hidden;
}

/* 非活动状态隐藏 */
.tabs-content-base[data-state="inactive"] {
  display: none;
}

/* 活动状态显示 */
.tabs-content-base[data-state="active"] {
  display: block;
}

/* 从右侧滑入动画 */
.tabs-slide-from-right {
  animation: slideFromRight 0.35s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

@keyframes slideFromRight {
  from {
    opacity: 0;
    transform: translateX(24px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 从左侧滑入动画 */
.tabs-slide-from-left {
  animation: slideFromLeft 0.35s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

@keyframes slideFromLeft {
  from {
    opacity: 0;
    transform: translateX(-24px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
