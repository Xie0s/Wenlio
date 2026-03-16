<!--
  src/components/personalization/theme-list/ThemeRowCompact.vue
  职责：主题列表密集行子组件，渲染紧凑单行主题信息
  对外暴露：Props(theme, showIcon, showDescription, showSlug, showDate, loading)、Emits(click)
-->
<script setup lang="ts">
import type { Theme } from '@/utils/types'
import { ArrowUpRight, Loader2 } from 'lucide-vue-next'

defineProps<{
  theme: Theme
  showDescription: boolean
  showSlug: boolean
  showDate: boolean
  loading: boolean
}>()

defineEmits<{ click: [] }>()

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}
</script>

<template>
  <div
    class="group flex items-center gap-2 px-4 py-3 cursor-pointer transition-colors hover:bg-muted/40 active:bg-muted/60 sm:py-2.5"
    :class="{ 'pointer-events-none opacity-60': loading }"
    @click="$emit('click')"
  >
    <!-- 左侧：名称 + slug + 描述（单行，更紧凑） -->
    <div class="min-w-0 flex-1 flex flex-col items-start gap-1 overflow-hidden sm:flex-row sm:items-center sm:gap-2">
      <span class="min-w-0 max-w-full truncate text-sm font-medium leading-tight sm:text-base sm:leading-none">{{ theme.name }}</span>
      <code
        v-if="showSlug"
        class="shrink-0 text-xs font-mono border border-border/50 rounded px-1.5 py-0.5 leading-none text-muted-foreground/65 hidden sm:inline"
      >{{ theme.slug }}</code>
      <span
        v-if="showDescription && theme.description"
        class="w-full min-w-0 truncate text-xs text-muted-foreground/60 leading-none"
      >{{ theme.description }}</span>
    </div>

    <!-- 右侧：日期 + 箭头 -->
    <div class="flex items-center gap-2 shrink-0">
      <span
        v-if="showDate && theme.created_at"
        class="text-xs text-muted-foreground/55 tabular-nums hidden sm:inline"
      >{{ formatDate(theme.created_at) }}</span>
      <Loader2 v-if="loading" class="h-3.5 w-3.5 animate-spin text-muted-foreground/40" />
      <ArrowUpRight
        v-else
        class="h-3.5 w-3.5 text-muted-foreground/40 transition-colors sm:text-transparent sm:group-hover:text-muted-foreground/50"
      />
    </div>
  </div>
</template>
