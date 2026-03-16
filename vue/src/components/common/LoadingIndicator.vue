<!--
  加载指示器组件
  功能：参考 Google Material Design 风格的极简加载动画
  对外暴露：作为全局/局部加载状态显示使用
-->
<template>
  <div class="loader-container">
    <div class="loader-content">
      <!-- Logo 动态加载 -->
      <div class="logo-wrapper">
        <img :src="logoUrl" alt="Logo" class="logo" />
      </div>

      <!-- Google 风格线性进度指示器 -->
      <div class="progress-wrapper">
        <div class="progress-track">
          <div class="progress-bar"></div>
        </div>
      </div>

      <!-- 状态文字 -->
      <div class="loader-text">
        <span class="text-content">{{ loadingText }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'

interface Props {
  isRunning?: boolean
  loadingTexts?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  isRunning: true,
  loadingTexts: () => ['正在加载', '准备资源', '即将完成']
})

const logoUrl = '/logo.gif'
const safeLoadingTexts = computed(() =>
  props.loadingTexts.length > 0 ? props.loadingTexts : ['正在加载'],
)
const loadingText = ref(safeLoadingTexts.value[0])
let textIndex = 0
let textInterval: number | null = null

function startAnimation() {
  if (!props.isRunning) return

  stopAnimation()

  const texts = safeLoadingTexts.value
  textIndex = 0
  loadingText.value = texts[0]

  if (texts.length <= 1) {
    return
  }

  textInterval = window.setInterval(() => {
    textIndex = (textIndex + 1) % texts.length
    loadingText.value = texts[textIndex]
  }, 1500)
}

function stopAnimation() {
  if (textInterval) {
    clearInterval(textInterval)
    textInterval = null
  }
}

watch(
  [() => props.isRunning, safeLoadingTexts],
  ([running]) => {
    if (running) {
      startAnimation()
    } else {
      stopAnimation()
    }
  },
  {
    immediate: true,
  }
)

onUnmounted(() => {
  stopAnimation()
})
</script>

<style scoped>
.loader-container {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100dvh;
  width: 100vw;
  background: var(--background);
}

.loader-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 24px;
}

/* Logo 样式 */
.logo-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo {
  width: 280px;
  max-width: 70vw;
  height: auto;
  object-fit: contain;
}

.dark .logo {
  filter: brightness(0) invert(1);
}

/* Google 风格线性进度条 */
.progress-wrapper {
  width: 280px;
  max-width: 70vw;
}

.progress-track {
  width: 100%;
  height: 4px;
  background: var(--muted);
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.progress-bar {
  position: absolute;
  height: 100%;
  width: 30%;
  background: linear-gradient(90deg, #4285f4, #34a853, #fbbc05, #ea4335);
  background-size: 300% 100%;
  border-radius: 2px;
  animation: indeterminate 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite, gradient-shift 3s ease infinite;
  will-change: transform;
}

/* 状态文字 - 使用 HarmonyOS Sans 字体 */
.loader-text {
  font-size: 14px;
  color: var(--muted-foreground);
  font-weight: 400;
  letter-spacing: 0.1px;
  font-family: 'HarmonyOS Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}

.text-content {
  opacity: 0.9;
}

/* 版权信息 */
.copyright {
  margin-top: -8px;
  text-align: center;
}

.copyright-brand {
  font-size: 22px;
  color: var(--foreground);
  font-weight: 400;
  letter-spacing: 0.5px;
  font-family: 'HarmonyOS Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
}

.copyright-company {
  font-size: 16px;
  color: var(--foreground);
  margin-top: 8px;
  font-weight: 400;
}

/* Google 风格动画 - 不确定进度条 */
@keyframes indeterminate {
  0% {
    transform: translateX(-100%);
    left: 0;
  }
  100% {
    transform: translateX(433%);
    left: 0;
  }
}

/* 渐变色流动效果 */
@keyframes gradient-shift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
</style>
