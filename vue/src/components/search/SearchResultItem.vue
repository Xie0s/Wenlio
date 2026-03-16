<!--
  SearchResultItem.vue - 搜索结果项组件
  职责：展示单条文档搜索结果，包含文档图标、高亮标题、摘要片段、主题/版本标签
  对外暴露：Props: item (SearchResult), keyword (用于高亮); Emits: click
-->
<template>
  <button
    class="w-full flex items-center gap-4 px-5 py-4 text-left hover:bg-muted/50 transition-colors focus:outline-none focus-visible:bg-muted/50"
    @click="$emit('click', item)"
  >
    <!-- 文档图标 -->
    <div class="shrink-0 w-11 h-11 rounded-lg bg-primary/10 flex items-center justify-center">
      <FileText class="w-5 h-5 text-primary" />
    </div>

    <!-- 文档信息 -->
    <div class="flex-1 min-w-0">
      <!-- 标题（高亮关键词） -->
      <p class="text-base font-medium truncate" v-html="highlightedTitle" />
      <!-- 摘要 -->
      <p
        class="text-sm text-muted-foreground mt-1 line-clamp-2 leading-relaxed"
        v-html="highlightedSnippet"
      />
    </div>

    <!-- 右侧：主题 / 版本 -->
    <div class="shrink-0 text-right space-y-0.5">
      <p v-if="item.theme_name" class="text-sm text-muted-foreground truncate max-w-[100px]">
        {{ item.theme_name }}
      </p>
      <p v-if="item.version_name" class="text-sm text-muted-foreground/60 truncate max-w-[100px]">
        {{ item.version_name }}
      </p>
    </div>

    <!-- 箭头 -->
    <ChevronRight class="shrink-0 w-5 h-5 text-muted-foreground" />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { FileText, ChevronRight } from 'lucide-vue-next'
import type { SearchResult } from '@/utils/types'

const props = defineProps<{
  item: SearchResult
  keyword?: string
}>()

defineEmits<{
  click: [item: SearchResult]
}>()

/** 转义 HTML 特殊字符，防止 XSS */
function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

/** 在文本中高亮关键词 */
function highlight(text: string, keyword: string): string {
  if (!keyword || !keyword.trim()) return escapeHtml(text)
  const escaped = escapeHtml(text)
  const kw = keyword.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return escaped.replace(
    new RegExp(kw, 'gi'),
    match => `<mark class="bg-yellow-200 dark:bg-yellow-900/50 text-yellow-900 dark:text-yellow-100 px-0.5 rounded not-prose">${match}</mark>`,
  )
}

const highlightedTitle = computed(() => highlight(props.item.title, props.keyword ?? ''))

const highlightedSnippet = computed(() => highlight(props.item.snippet, props.keyword ?? ''))
</script>
