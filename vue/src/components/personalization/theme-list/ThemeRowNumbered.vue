<!--
  src/components/personalization/theme-list/ThemeRowNumbered.vue
  职责：主题列表编号行子组件，渲染带序号的主题条目
  对外暴露：Props(theme, showIcon, showDescription, showSlug, showDate, loading, index)、Emits(click)
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
  index: number
}>()

defineEmits<{ click: [] }>()

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<template>
  <div
    class="group flex items-center gap-3 px-4 py-3.5 cursor-pointer transition-colors hover:bg-muted/50 active:bg-muted/60 sm:px-5"
    :class="{ 'pointer-events-none opacity-60': loading }"
    @click="$emit('click')"
  >
    <!-- 序号 -->
    <span class="shrink-0 w-7 text-right text-sm font-mono tabular-nums text-muted-foreground/35 group-hover:text-muted-foreground/60 transition-colors select-none">
      {{ String(index).padStart(2, '0') }}
    </span>

    <!-- 左侧：名称 + slug + 描述（单行） -->
    <div class="min-w-0 flex-1 flex flex-col items-start gap-1 overflow-hidden sm:flex-row sm:items-center sm:gap-2">
      <span class="min-w-0 max-w-full truncate text-sm font-semibold leading-tight sm:text-base sm:leading-none">{{ theme.name }}</span>
      <code
        v-if="showSlug"
        class="shrink-0 text-xs font-mono border border-border/50 rounded px-1.5 py-0.5 leading-none text-muted-foreground/65 hidden sm:inline"
      >{{ theme.slug }}</code>
      <span
        v-if="showDescription"
        class="w-full min-w-0 truncate text-xs text-muted-foreground/65 leading-none"
      >{{ theme.description || '暂无描述' }}</span>
    </div>

    <!-- 右侧：日期 + loading -->
    <div class="flex items-center gap-2.5 shrink-0">
      <span
        v-if="showDate && theme.created_at"
        class="text-xs text-muted-foreground/55 tabular-nums hidden sm:inline"
      >{{ formatDate(theme.created_at) }}</span>
      <Loader2 v-if="loading" class="h-4 w-4 animate-spin text-muted-foreground/50" />
    </div>
  </div>
</template>
