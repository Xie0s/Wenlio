<!--
  NoiseBackground - 噪点渐变背景效果组件 (Aceternity UI)
  
  功能：提供带有动态渐变光斑 + 噪点纹理叠加的容器效果
  职责：作为装饰性容器组件，包裹内容并添加动态噪点渐变视觉增强
  主要接口：<NoiseBackground> 组件，通过 slot 包裹子内容
  
  参考：https://ui.aceternity.com/components/noise-background
-->
<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { cn } from '@/utils'
import { useThemeStore } from '@/stores/theme'

defineOptions({ inheritAttrs: false })

const props = withDefaults(defineProps<{
  /** 内容区域自定义类名 */
  class?: string
  /** 容器自定义类名 */
  containerClass?: string
  /** 容器自定义类名（兼容官方 API 命名） */
  containerClassName?: string
  /** 渐变颜色数组（建议 3 个） */
  gradientColors?: string[]
  /** 噪点强度 0~1 */
  noiseIntensity?: number
  /** 动画移动速度 */
  speed?: number
  /** 是否启用背景模糊后处理 */
  backdropBlur?: boolean
  /** 是否启用动画 */
  animating?: boolean
}>(), {
  gradientColors: () => [
    'rgb(255, 100, 150)',
    'rgb(100, 150, 255)',
    'rgb(255, 200, 100)',
  ],
  noiseIntensity: 0.2,
  speed: 0.1,
  backdropBlur: false,
  animating: true,
})

const themeStore = useThemeStore()
const { resolvedTheme } = storeToRefs(themeStore)

// ========== DOM Refs ==========
const containerRef = ref<HTMLElement | null>(null)
const gradient1Ref = ref<HTMLElement | null>(null)
const gradient2Ref = ref<HTMLElement | null>(null)
const gradient3Ref = ref<HTMLElement | null>(null)
const topStripRef = ref<HTMLElement | null>(null)

// ========== 动画状态（非响应式，直接操作 DOM 以保证 60fps） ==========
let posX = 0, posY = 0
let springX = 0, springY = 0
let springVelX = 0, springVelY = 0
let velX = 0, velY = 0
let lastDirectionChange = 0
let animFrameId: number | null = null
let isDarkMode = false

/** 检测暗色模式 */
function checkDarkMode() {
  isDarkMode = resolvedTheme.value === 'dark'
}

/** 生成随机方向速度 */
function generateRandomVelocity() {
  const angle = Math.random() * Math.PI * 2
  const magnitude = props.speed * (0.5 + Math.random() * 0.5)
  velX = Math.cos(angle) * magnitude
  velY = Math.sin(angle) * magnitude
}

/** 直接写入 DOM style，避免 Vue 响应式开销 */
function updateGradients() {
  checkDarkMode()
  const colors = props.gradientColors
  // 暗色模式下提高渐变透明度以增强可见性
  const opacity1 = isDarkMode ? 0.55 : 0.4
  const opacity2 = isDarkMode ? 0.45 : 0.3
  const opacity3 = isDarkMode ? 0.35 : 0.25
  if (gradient1Ref.value) {
    gradient1Ref.value.style.opacity = String(opacity1)
    gradient1Ref.value.style.background =
      `radial-gradient(circle at ${springX}px ${springY}px, ${colors[0]} 0%, transparent 50%)`
  }
  if (gradient2Ref.value) {
    gradient2Ref.value.style.opacity = String(opacity2)
    gradient2Ref.value.style.background =
      `radial-gradient(circle at ${springX * 0.7}px ${springY * 0.7}px, ${colors[1]} 0%, transparent 50%)`
  }
  if (gradient3Ref.value) {
    const c3 = colors[2] || colors[0]
    gradient3Ref.value.style.opacity = String(opacity3)
    gradient3Ref.value.style.background =
      `radial-gradient(circle at ${springX * 1.2}px ${springY * 1.2}px, ${c3} 0%, transparent 50%)`
  }
  if (topStripRef.value && props.animating) {
    topStripRef.value.style.transform = `translateX(${springX * 0.1 - 50}px)`
  }
}

/** requestAnimationFrame 驱动的动画循环 */
function animate(time: number) {
  if (!props.animating || !containerRef.value) {
    animFrameId = null
    return
  }

  const rect = containerRef.value.getBoundingClientRect()
  const maxX = rect.width
  const maxY = rect.height
  const padding = 20

  // 每 1.5~3 秒随机换方向
  if (time - lastDirectionChange > 1500 + Math.random() * 1500) {
    generateRandomVelocity()
    lastDirectionChange = time
  }

  const deltaTime = 16
  let newX = posX + velX * deltaTime
  let newY = posY + velY * deltaTime

  // 碰边反弹：生成全新随机方向
  if (newX < padding || newX > maxX - padding || newY < padding || newY > maxY - padding) {
    generateRandomVelocity()
    lastDirectionChange = time
    newX = Math.max(padding, Math.min(maxX - padding, newX))
    newY = Math.max(padding, Math.min(maxY - padding, newY))
  }

  posX = newX
  posY = newY

  // 弹簧物理平滑（stiffness ≈ 100, damping ≈ 30）
  const stiffness = 0.003
  const damping = 0.93
  springVelX = (springVelX + (posX - springX) * stiffness) * damping
  springVelY = (springVelY + (posY - springY) * stiffness) * damping
  springX += springVelX
  springY += springVelY

  updateGradients()
  animFrameId = requestAnimationFrame(animate)
}

onMounted(() => {
  if (containerRef.value) {
    const rect = containerRef.value.getBoundingClientRect()
    posX = rect.width / 2
    posY = rect.height / 2
    springX = posX
    springY = posY
  }
  generateRandomVelocity()
  updateGradients()
  animFrameId = requestAnimationFrame(animate)
})

onUnmounted(() => {
  if (animFrameId != null) cancelAnimationFrame(animFrameId)
})

watch(() => props.speed, () => generateRandomVelocity())

// animating 恢复时重新启动 RAF 循环
watch(() => props.animating, (val) => {
  if (val && animFrameId == null) {
    animFrameId = requestAnimationFrame(animate)
  }
})
</script>

<template>
  <div
    ref="containerRef"
    :class="cn(
      'group relative overflow-hidden rounded-2xl bg-neutral-200 p-2 backdrop-blur-sm dark:bg-neutral-800',
      'shadow-[0px_0.5px_1px_0px_var(--color-neutral-400)_inset,0px_1px_0px_0px_var(--color-neutral-100)]',
      'dark:shadow-[0px_1px_0px_0px_var(--color-neutral-950)_inset,0px_1px_0px_0px_var(--color-neutral-800)]',
      backdropBlur && `after:absolute after:inset-0 after:h-full after:w-full after:backdrop-blur-lg after:content-['']`,
      containerClassName ?? containerClass,
    )"
    :style="{ '--noise-opacity': noiseIntensity } as any"
  >
    <!-- 动态渐变光斑层 -->
    <div ref="gradient1Ref" class="absolute inset-0" />
    <div ref="gradient2Ref" class="absolute inset-0" />
    <div ref="gradient3Ref" class="absolute inset-0" />

    <!-- 顶部渐变条 -->
    <div
      ref="topStripRef"
      class="absolute inset-x-0 top-0 h-1 rounded-t-2xl opacity-80 blur-sm"
      :style="{ background: `linear-gradient(to right, ${gradientColors.join(', ')})` }"
    />

    <!-- 噪点纹理叠加层 -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <img
        src="/common/noise.webp"
        alt=""
        class="h-full w-full object-cover dark:invert dark:opacity-[0.12]"
        :style="{ opacity: noiseIntensity, mixBlendMode: 'soft-light' }"
      />
    </div>

    <!-- 内容插槽 -->
    <div :class="cn('relative z-10', props.class)">
      <slot />
    </div>
  </div>
</template>
