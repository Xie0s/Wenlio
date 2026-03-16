/**
 * HeroUI CardFooter 组件
 * 职责：卡片的底部区域，通常用于放置按钮、链接或次要信息。
 */
<script setup lang="ts">
import { computed, type HTMLAttributes } from 'vue'
import { cn } from '@/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
  isBlurred?: boolean
}>()

const footerClasses = computed(() => {
  return cn(
    'flex p-3 h-auto w-full items-center overflow-hidden color-inherit subpixel-antialiased',
    {
      'absolute bottom-0 z-10 border-t border-white/20 bg-black/10 backdrop-blur-md backdrop-saturate-[1.8]': props.isBlurred,
    },
    props.class
  )
})
</script>

<template>
  <div :class="footerClasses">
    <slot />
  </div>
</template>
