/**
 * HeroUI Card 组件
 * 职责：作为卡片容器，提供阴影、圆角、悬停、点击等基础样式和交互状态。
 */
<script setup lang="ts">
import { computed, type HTMLAttributes } from 'vue'
import { cn } from '@/utils'

interface CardProps {
  class?: HTMLAttributes['class']
  radius?: 'none' | 'sm' | 'md' | 'lg' | '2xl' | 'full'
  isHoverable?: boolean
  isPressable?: boolean
  isBlurred?: boolean
  isFooterBlurred?: boolean
  isDisabled?: boolean
  fullWidth?: boolean
}

const props = withDefaults(defineProps<CardProps>(), {
  radius: '2xl',
  isHoverable: false,
  isPressable: false,
  isBlurred: false,
  isFooterBlurred: false,
  isDisabled: false,
  fullWidth: false,
})

const cardClasses = computed(() => {
  return cn(
    'relative flex flex-col overflow-hidden height-auto bg-background text-foreground box-border border border-border/50 outline outline-1 outline-foreground/5 -outline-offset-1',
    // 圆角 (根据全局规范，默认使用 2xl)
    {
      'rounded-none': props.radius === 'none',
      'rounded-sm': props.radius === 'sm',
      'rounded-md': props.radius === 'md',
      'rounded-lg': props.radius === 'lg',
      'rounded-2xl': props.radius === '2xl',
      'rounded-full': props.radius === 'full',
    },
    // 交互状态
    {
      'transition-transform-background motion-reduce:transition-none': props.isHoverable || props.isPressable,
      'hover:bg-accent/50 cursor-pointer': props.isHoverable,
      'active:scale-[0.97] cursor-pointer': props.isPressable,
      'opacity-50 pointer-events-none': props.isDisabled,
      'w-full': props.fullWidth,
    },
    // 模糊效果
    {
      'backdrop-blur-md backdrop-saturate-[1.8] bg-background/70': props.isBlurred,
    },
    props.class
  )
})
</script>

<template>
  <div :class="cardClasses" role="section">
    <slot />
  </div>
</template>
