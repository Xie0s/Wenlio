/**
 * HeroUI CardBody 组件
 * 职责：卡片的主体内容区域，提供默认的内边距和滚动支持。
 */
<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div :class="cn('relative flex flex-col flex-auto p-3 w-full h-auto break-words overflow-y-auto subpixel-antialiased', props.class)">
    <slot />
  </div>
</template>
