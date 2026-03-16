<!--
  AnimatedThemeToggler.vue - Magic UI 风格的主题切换根组件
  职责：提供基于 View Transition API 的明暗主题切换动画按钮
  对外接口：
    - Props:
      - duration?: number 主题扩散动画时长（毫秒），默认 400
      - class?: string 自定义样式类
    - Emits: 无
    - 依赖：@/stores/theme（统一管理主题状态与持久化）
-->
<script setup lang="ts">
import { computed, nextTick, ref, useAttrs } from 'vue'
import { storeToRefs } from 'pinia'
import type { ButtonHTMLAttributes } from 'vue'
import { Moon, Sun } from 'lucide-vue-next'
import { cn } from '@/utils'
import { useThemeStore } from '@/stores/theme'

interface Props {
  duration?: number
  class?: ButtonHTMLAttributes['class']
  buttonSize?: string
  iconSize?: string
}

const props = withDefaults(defineProps<Props>(), {
  duration: 500,
  buttonSize: 'size-9',
  iconSize: 'size-5',
})

const attrs = useAttrs()
const themeStore = useThemeStore()
const { resolvedTheme } = storeToRefs(themeStore)
const buttonRef = ref<HTMLButtonElement | null>(null)
const isDark = computed(() => resolvedTheme.value === 'dark')
const isAnimating = ref(false)

// 与官方一致：toggleTheme 核心逻辑（增加防连点与 VT API 降级）
async function toggleTheme() {
  const btn = buttonRef.value
  if (!btn || isAnimating.value) return

  // 无 View Transition API 时直接切换（兼容旧浏览器）
  if (!(document as any).startViewTransition) {
    themeStore.toggleTheme()
    return
  }

  isAnimating.value = true

  // Vue 中用 async + nextTick 替代 React 的 flushSync，确保 VT 截取到新 DOM 状态
  try {
    await (document as any).startViewTransition(async () => {
      themeStore.toggleTheme()
      await nextTick()
    }).ready
  } catch {
    isAnimating.value = false
    return
  }

  // 与官方一致：在 ready 后获取按钮位置，计算圆形扩散半径
  const { top, left, width, height } = btn.getBoundingClientRect()
  const x = left + width / 2
  const y = top + height / 2
  const maxRadius = Math.hypot(
    Math.max(left, window.innerWidth - left),
    Math.max(top, window.innerHeight - top),
  )

  // 与官方一致：fire-and-forget，不 await .finished
  const anim = document.documentElement.animate(
    {
      clipPath: [
        `circle(0px at ${x}px ${y}px)`,
        `circle(${maxRadius}px at ${x}px ${y}px)`,
      ],
    },
    {
      duration: props.duration,
      easing: 'ease-in-out',
      pseudoElement: '::view-transition-new(root)',
    },
  )

  // 动画结束后释放防连点锁
  anim.onfinish = () => { isAnimating.value = false }
  anim.oncancel = () => { isAnimating.value = false }
}
</script>

<template>
  <button ref="buttonRef" type="button" :class="cn(
    'inline-flex items-center justify-center rounded-full outline-none transition-all focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground',
    props.buttonSize,
    props.class,
  )" v-bind="attrs" @click="toggleTheme">
    <component :is="isDark ? Sun : Moon" :class="props.iconSize" :stroke-width="1.5" />
    <span class="sr-only">Toggle theme</span>
  </button>
</template>
