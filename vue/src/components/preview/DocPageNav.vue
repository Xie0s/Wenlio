<!-- components/DocPage/DocPageNav.vue
     职责：文档底部上一页 / 下一页翻页导航（Nextra / VitePress 标志性设计）
     对外暴露事件：navigate(slug) -->

<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

defineProps<{
  prevPage: { title: string; slug: string } | null
  nextPage: { title: string; slug: string } | null
}>()

const emit = defineEmits<{
  navigate: [slug: string]
}>()
</script>

<template>
  <div class="mt-14 flex items-stretch justify-between gap-4 border-t border-border pt-8">
    <!-- 上一页 -->
    <button
      v-if="prevPage"
      class="group flex flex-col items-start gap-1 rounded-2xl border border-border px-5 py-4 text-left
             transition-all duration-200 hover:border-primary/40 hover:bg-accent/40 min-w-0 flex-1 max-w-[48%]"
      @click="emit('navigate', prevPage.slug)"
    >
      <span class="flex items-center gap-1 text-xs text-muted-foreground">
        <ChevronLeft class="h-3.5 w-3.5 transition-transform duration-200 group-hover:-translate-x-0.5" />
        上一页
      </span>
      <span class="mt-0.5 text-sm font-medium text-foreground group-hover:text-primary transition-colors duration-200 line-clamp-2 w-full">
        {{ prevPage.title }}
      </span>
    </button>
    <div v-else class="flex-1" />

    <!-- 下一页 -->
    <button
      v-if="nextPage"
      class="group flex flex-col items-end gap-1 rounded-2xl border border-border px-5 py-4 text-right
             transition-all duration-200 hover:border-primary/40 hover:bg-accent/40 min-w-0 flex-1 max-w-[48%]"
      @click="emit('navigate', nextPage.slug)"
    >
      <span class="flex items-center gap-1 text-xs text-muted-foreground">
        下一页
        <ChevronRight class="h-3.5 w-3.5 transition-transform duration-200 group-hover:translate-x-0.5" />
      </span>
      <span class="mt-0.5 text-sm font-medium text-foreground group-hover:text-primary transition-colors duration-200 line-clamp-2 w-full">
        {{ nextPage.title }}
      </span>
    </button>
    <div v-else class="flex-1" />
  </div>
</template>
