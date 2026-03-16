<!--
  src/components/personalization/theme-list/ThemeRowTag.vue
  职责：主题列表标签云子组件，渲染药丸样式主题标签
  对外暴露：Props(theme, showIcon, showDescription, showSlug, showDate, loading)、Emits(click)
-->
<script setup lang="ts">
import type { Theme } from '@/utils/types'
import { Loader2 } from 'lucide-vue-next'

defineProps<{
  theme: Theme
  showDescription: boolean
  showSlug: boolean
  showDate: boolean
  loading: boolean
}>()

defineEmits<{ click: [] }>()
</script>

<template>
  <div
    class="inline-flex items-center gap-2 px-3 py-2 rounded-full border border-border/70 bg-card
           hover:border-primary/50 hover:bg-primary/5 hover:shadow-sm active:bg-primary/10
           text-sm sm:text-base font-medium cursor-pointer transition-all duration-200 select-none"
    :class="{ 'pointer-events-none opacity-50': loading }"
    @click="$emit('click')"
  >
    <Loader2 v-if="loading" class="h-3.5 w-3.5 animate-spin text-muted-foreground" />
    <span class="leading-none">{{ theme.name }}</span>
    <!-- slug 以无背景 outline 标签跟在名称后 -->
    <code
      v-if="showSlug"
      class="shrink-0 text-xs font-mono border border-border/50 rounded px-1 py-0.5 leading-none text-muted-foreground/60 hidden sm:inline"
    >{{ theme.slug }}</code>
  </div>
</template>
