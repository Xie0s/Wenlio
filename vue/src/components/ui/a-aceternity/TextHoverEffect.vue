<!--
  TextHoverEffect.vue - SVG 文字悬停光效组件
  移植自 https://ui.aceternity.com/components/text-hover-effect
  职责：鼠标悬停时在 SVG 文字上呈现跟随鼠标的彩色径向渐变光效，并在挂载时播放描边动画。
  对外暴露：
    props: text (string) - 展示的文字
            duration (number, 可选) - 径向渐变跟随平滑度，0 为即时跟随，默认 0
  移植来源：https://ui.aceternity.com/components/text-hover-effect
-->
<script setup lang="ts">
import { ref, onMounted, onUnmounted, getCurrentInstance } from 'vue'
import { animate } from 'motion'

const props = withDefaults(defineProps<{
  text: string
  duration?: number
}>(), {
  duration: 0,
})

// 每个实例使用唯一 ID 避免多实例 SVG id 冲突
const uid = getCurrentInstance()?.uid ?? Math.random().toString(36).slice(2)
const gradientId = `textGradient-${uid}`
const maskId = `revealMask-${uid}`
const textMaskId = `textMask-${uid}`

const svgRef = ref<SVGSVGElement | null>(null)
const animatedTextRef = ref<SVGTextElement | null>(null)
const radialGradRef = ref<SVGRadialGradientElement | null>(null)

const hovered = ref(false)

// 使用非响应式变量 + setAttribute，避免 rAF 频繁触发 Vue 渲染
let currentCx = 50
let currentCy = 50
let targetCx = 50
let targetCy = 50
let rafId = 0
// 指数平滑系数，仅在 duration 变化时重算（避免每帧重复计算）
let smoothFactor = 0

function computeFactor(d: number) {
  smoothFactor = d > 0 ? 1 - Math.exp(-8 / (d * 60)) : 0
}

function tick() {
  currentCx += (targetCx - currentCx) * smoothFactor
  currentCy += (targetCy - currentCy) * smoothFactor

  const grad = radialGradRef.value
  if (grad) {
    grad.setAttribute('cx', `${currentCx}%`)
    grad.setAttribute('cy', `${currentCy}%`)
  }

  // 收敛后自动停止循环，避免空转
  const converged = Math.abs(targetCx - currentCx) < 0.01 && Math.abs(targetCy - currentCy) < 0.01
  if (!converged) {
    rafId = requestAnimationFrame(tick)
  } else {
    rafId = 0
  }
}

function startLoop() {
  if (rafId === 0 && smoothFactor > 0) {
    rafId = requestAnimationFrame(tick)
  }
}

function stopLoop() {
  if (rafId !== 0) {
    cancelAnimationFrame(rafId)
    rafId = 0
  }
}

function setGradientPos(cx: number, cy: number) {
  const grad = radialGradRef.value
  if (grad) {
    grad.setAttribute('cx', `${cx}%`)
    grad.setAttribute('cy', `${cy}%`)
  }
}

function onMouseEnter() {
  hovered.value = true
}
function onMouseLeave() {
  hovered.value = false
  // 平滑模式下 mouseleave 后不再需要追踪，停止循环
  if (smoothFactor > 0) stopLoop()
}
function onMouseMove(e: MouseEvent) {
  if (!svgRef.value) return
  const rect = svgRef.value.getBoundingClientRect()
  targetCx = ((e.clientX - rect.left) / rect.width) * 100
  targetCy = ((e.clientY - rect.top) / rect.height) * 100

  if (smoothFactor === 0) {
    // duration=0：直接更新，无需 rAF 循环
    currentCx = targetCx
    currentCy = targetCy
    setGradientPos(currentCx, currentCy)
  } else {
    // duration>0：启动平滑循环（若未在跑）
    startLoop()
  }
}

onMounted(() => {
  computeFactor(props.duration ?? 0)

  // 描边动画：stroke-dashoffset 1000 → 0
  if (animatedTextRef.value) {
    animate(
      animatedTextRef.value,
      { strokeDashoffset: [1000, 0] } as any,
      { duration: 4, ease: 'easeInOut' } as any,
    )
  }
})

onUnmounted(() => {
  stopLoop()
})
</script>

<template>
  <svg
    ref="svgRef"
    width="100%"
    height="100%"
    viewBox="0 0 300 100"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
    @mousemove="onMouseMove"
  >
    <defs>
      <!-- 彩色线性渐变，悬停时激活 -->
      <linearGradient
        :id="gradientId"
        gradientUnits="userSpaceOnUse"
        cx="50%"
        cy="50%"
        r="25%"
      >
        <template v-if="hovered">
          <stop offset="0%" stop-color="#eab308" />
          <stop offset="25%" stop-color="#ef4444" />
          <stop offset="50%" stop-color="#3b82f6" />
          <stop offset="75%" stop-color="#06b6d4" />
          <stop offset="100%" stop-color="#8b5cf6" />
        </template>
      </linearGradient>

      <!-- 径向渐变遮罩，cx/cy 由 rAF 直接驱动 -->
      <radialGradient
        ref="radialGradRef"
        :id="maskId"
        gradientUnits="userSpaceOnUse"
        r="20%"
        cx="50%"
        cy="50%"
      >
        <stop offset="0%" stop-color="white" />
        <stop offset="100%" stop-color="black" />
      </radialGradient>

      <mask :id="textMaskId">
        <rect x="0" y="0" width="100%" height="100%" :fill="`url(#${maskId})`" />
      </mask>
    </defs>

    <!-- 层 1：悬停时显示的静态描边文字 -->
    <text
      x="50%"
      y="50%"
      text-anchor="middle"
      dominant-baseline="middle"
      stroke-width="0.3"
      class="fill-transparent stroke-neutral-200 font-[helvetica] text-7xl font-bold dark:stroke-neutral-800"
      :style="{ opacity: hovered ? 0.7 : 0, transition: 'opacity 0.3s' }"
    >{{ text }}</text>

    <!-- 层 2：挂载时播放描边动画的文字 -->
    <text
      ref="animatedTextRef"
      x="50%"
      y="50%"
      text-anchor="middle"
      dominant-baseline="middle"
      stroke-width="0.3"
      stroke-dasharray="1000"
      stroke-dashoffset="1000"
      class="fill-transparent stroke-neutral-200 font-[helvetica] text-7xl font-bold dark:stroke-neutral-800"
    >{{ text }}</text>

    <!-- 层 3：彩色渐变文字，跟随鼠标的径向遮罩控制显示区域 -->
    <text
      x="50%"
      y="50%"
      text-anchor="middle"
      dominant-baseline="middle"
      :stroke="`url(#${gradientId})`"
      stroke-width="0.3"
      :mask="`url(#${textMaskId})`"
      class="fill-transparent font-[helvetica] text-7xl font-bold"
    >{{ text }}</text>
  </svg>
</template>
