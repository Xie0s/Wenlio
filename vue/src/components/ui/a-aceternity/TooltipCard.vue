<!--
 * TooltipCard.vue - Aceternity UI 风格鼠标跟随 Tooltip 卡片
 * 职责：跟随鼠标位置展示富内容 Tooltip，支持视口边界检测与弹簧动画
 * 对外暴露：
 *   - Props: containerClassName?: string
 *   - Slots: default（触发元素），content（Tooltip 内容，支持富文本）
 * 移植来源：https://ui.aceternity.com/components/tooltip-card
-->
<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { cn } from '@/utils'

defineProps<{
  containerClassName?: string
}>()

const isVisible = ref(false)
const mouse = ref({ x: 0, y: 0 })
const position = ref({ x: 0, y: 0 })
const contentRef = ref<HTMLDivElement | null>(null)

// --- Transition JS hooks：弹簧高度动画 ---
const SPRING_EASING = 'cubic-bezier(0.34, 1.56, 0.64, 1)'
const EASE_OUT = 'cubic-bezier(0.25, 0.46, 0.45, 0.94)'

const onBeforeEnter = (el: Element) => {
  const htmlEl = el as HTMLElement
  htmlEl.style.height = '0px'
  htmlEl.style.overflow = 'hidden'
}

const onEnter = (el: Element, done: () => void) => {
  const htmlEl = el as HTMLElement
  const targetH = htmlEl.scrollHeight
  void htmlEl.offsetHeight // force reflow
  htmlEl.style.transition = `height 0.45s ${SPRING_EASING}`
  htmlEl.style.height = `${targetH}px`
  const cleanup = () => {
    htmlEl.style.height = 'auto'
    htmlEl.style.transition = ''
    done()
  }
  htmlEl.addEventListener('transitionend', cleanup, { once: true })
}

const onLeave = (el: Element, done: () => void) => {
  const htmlEl = el as HTMLElement
  htmlEl.style.height = `${htmlEl.scrollHeight}px`
  htmlEl.style.overflow = 'hidden'
  void htmlEl.offsetHeight // force reflow
  htmlEl.style.transition = `height 0.3s ${EASE_OUT}`
  htmlEl.style.height = '0px'
  htmlEl.addEventListener('transitionend', done, { once: true })
}

const calculatePosition = (clientX: number, clientY: number) => {
  const tooltipWidth = 288 // w-[18rem] = 288px
  const tooltipHeight = contentRef.value?.scrollHeight ?? 120

  let finalX = clientX + 12
  let finalY = clientY + 12

  if (finalX + tooltipWidth > window.innerWidth) finalX = clientX - tooltipWidth - 12
  if (finalX < 0) finalX = 12
  if (finalY + tooltipHeight > window.innerHeight) finalY = clientY - tooltipHeight - 12
  if (finalY < 0) finalY = 12

  return { x: finalX, y: finalY }
}

const updateMousePosition = (clientX: number, clientY: number) => {
  mouse.value = { x: clientX, y: clientY }
  position.value = calculatePosition(clientX, clientY)
}

// isVisible 变更时重算位置
watch(isVisible, async (vis) => {
  if (vis) {
    await nextTick()
    position.value = calculatePosition(mouse.value.x, mouse.value.y)
  }
})

const handleMouseEnter = (e: MouseEvent) => {
  isVisible.value = true
  updateMousePosition(e.clientX, e.clientY)
}

const handleMouseLeave = () => {
  mouse.value = { x: 0, y: 0 }
  position.value = { x: 0, y: 0 }
  isVisible.value = false
}

const handleMouseMove = (e: MouseEvent) => {
  if (!isVisible.value) return
  updateMousePosition(e.clientX, e.clientY)
}

const handleTouchStart = (e: TouchEvent) => {
  const touch = e.touches[0]
  if (!touch) return
  updateMousePosition(touch.clientX, touch.clientY)
  isVisible.value = true
}

const handleTouchEnd = () => {
  setTimeout(() => {
    isVisible.value = false
    mouse.value = { x: 0, y: 0 }
    position.value = { x: 0, y: 0 }
  }, 2000)
}

const handleClick = (e: MouseEvent) => {
  if (window.matchMedia('(hover: none)').matches) {
    e.preventDefault()
    if (isVisible.value) {
      isVisible.value = false
      mouse.value = { x: 0, y: 0 }
      position.value = { x: 0, y: 0 }
    } else {
      updateMousePosition(e.clientX, e.clientY)
      isVisible.value = true
    }
  }
}
</script>

<template>
  <div :class="cn('relative inline-block', containerClassName)" @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave" @mousemove="handleMouseMove" @touchstart.passive="handleTouchStart"
    @touchend="handleTouchEnd" @click="handleClick">
    <slot />

    <Transition :css="false" @before-enter="onBeforeEnter" @enter="onEnter" @leave="onLeave">
      <div v-if="isVisible"
        class="glass pointer-events-none fixed z-[999] w-[18rem] overflow-hidden rounded-3xl ring-[0.5px] ring-inset ring-black/[0.07] dark:ring-white/[0.08]"
        :style="{ top: `${position.y}px`, left: `${position.x}px` }">
        <div ref="contentRef" class="p-2 text-sm text-neutral-600 md:p-4 dark:text-neutral-400">
          <slot name="content" />
        </div>
      </div>
    </Transition>
  </div>
</template>
