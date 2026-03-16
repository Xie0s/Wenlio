<!--
  src/components/personalization/theme-list/ThemeRowHeadline.vue
  职责：主题列表大标题行子组件，突出展示主题名称与摘要信息
  对外暴露：Props(theme, showIcon, showDescription, showSlug, showDate, loading)、Emits(click)
-->
<script setup lang="ts">
import type { Theme } from '@/utils/types'
import { ArrowRight, Loader2 } from 'lucide-vue-next'

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
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<template>
  <div
    class="group flex items-center gap-3 py-4 border-b border-border/40 last:border-0 cursor-pointer active:bg-muted/40"
    :class="{ 'pointer-events-none opacity-60': loading }"
    @click="$emit('click')"
  >
    <!-- 左侧：名称（更大字号）+ slug + 描述（单行） -->
    <div class="min-w-0 flex-1 flex flex-col items-start gap-1.5 overflow-hidden sm:flex-row sm:items-center sm:gap-2.5">
      <h3 class="min-w-0 max-w-full truncate text-base sm:text-lg font-semibold leading-tight sm:leading-none group-hover:text-primary transition-colors duration-200">
        {{ theme.name }}
      </h3>
      <code
        v-if="showSlug"
        class="shrink-0 text-xs font-mono border border-border/50 rounded px-1.5 py-0.5 leading-none text-muted-foreground/65 hidden sm:inline"
      >{{ theme.slug }}</code>
      <span
        v-if="showDescription && theme.description"
        class="w-full min-w-0 truncate text-xs sm:text-sm text-muted-foreground/65 leading-none"
      >{{ theme.description }}</span>
    </div>

    <!-- 右侧：日期 + 箭头 -->
    <div class="flex items-center gap-3 shrink-0">
      <span
        v-if="showDate && theme.created_at"
        class="text-xs text-muted-foreground/55 tabular-nums hidden sm:inline"
      >{{ formatDate(theme.created_at) }}</span>
      <Loader2 v-if="loading" class="h-4 w-4 animate-spin text-muted-foreground" />
      <ArrowRight
        v-else
        class="h-4 w-4 shrink-0 text-muted-foreground/35 sm:text-muted-foreground/25 group-hover:text-primary/60 transition-colors duration-200"
      />
    </div>
  </div>
</template>
