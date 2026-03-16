/**
 * HeroUI CardHeader 组件
 * 职责：卡片的头部区域，通常用于放置标题、描述或头像。
 */
<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div :class="cn('flex p-3 z-10 w-full justify-start items-center shrink-0 overflow-inherit color-inherit subpixel-antialiased', props.class)">
    <slot />
  </div>
</template>
