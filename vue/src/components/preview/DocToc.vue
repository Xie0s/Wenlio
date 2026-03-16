<!-- components/DocPage/DocToc.vue
     职责：文档阅读页右侧目录导航（桌面端，xl+ 可见）
     设计风格参考 Nextra："本页目录" 标题，层级缩进，激活高亮
     对外暴露事件：scroll-to(id) -->

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import type { TocItem } from '@/composables/useToc'
import { useScrollToActive } from '@/composables/useScrollToActive'
import { levelPaddingClass } from '@/lib/toc'

const props = defineProps<{
  items: TocItem[]
  activeId: string
}>()

const emit = defineEmits<{
  'scroll-to': [id: string]
}>()

const navRef = ref<HTMLElement | null>(null)

useScrollToActive(navRef, computed(() => props.activeId), computed(() => props.items))

const route = useRoute()
const searchKeyword = computed(() => (route.query.q as string) || '')

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function highlight(text: string, keyword: string): string {
  if (!keyword.trim()) return escapeHtml(text)
  const escaped = escapeHtml(text)
  const kw = keyword.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return escaped.replace(
    new RegExp(kw, 'gi'),
    match => `<mark class="bg-yellow-200 dark:bg-yellow-900/50 text-yellow-900 dark:text-yellow-100 px-0.5 rounded">${match}</mark>`,
  )
}
</script>

<template>
  <aside class="sticky top-4 hidden h-[calc(100vh-28rem)] w-72 shrink-0 self-start xl:flex xl:flex-col">
    <nav ref="navRef" class="toc-scroll min-h-0 flex-1 overflow-y-auto py-2 pr-3">
      <ul class="space-y-1">
        <li v-for="item in items" :key="item.id">
          <button :data-toc-id="item.id"
            class="group relative block w-full border-l-2 py-1.5 pr-2 text-left text-base leading-6 transition-colors duration-150"
            :class="[
              levelPaddingClass(item.level),
              activeId === item.id
                ? 'border-primary text-primary font-normal'
                : 'border-transparent text-foreground/70 hover:text-foreground hover:border-muted-foreground/40',
            ]" @click="emit('scroll-to', item.id)">
            <span class="block break-words" v-html="highlight(item.text, searchKeyword)" />
          </button>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<style scoped>
.toc-scroll {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.toc-scroll::-webkit-scrollbar {
  display: none;
}
</style>
