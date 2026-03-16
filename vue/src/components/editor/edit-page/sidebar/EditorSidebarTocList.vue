<!--
  EditorSidebarTocList.vue - 编辑器侧边栏 TOC 列表
  职责：负责渲染 TOC 导航与高亮同步滚动；不处理章节树、CRUD 表单和编辑器主状态
  对外暴露：props(items, activeId)；emits(scrollTo)
-->
<script setup lang="ts">
import { computed, ref } from 'vue'
import type { TocItem } from '@/composables/useToc'
import { useScrollToActive } from '@/composables/useScrollToActive'
import { levelPaddingClass } from '@/lib/toc'

const props = defineProps<{
  items: TocItem[]
  activeId: string
}>()

const emit = defineEmits<{
  scrollTo: [id: string]
}>()

const containerRef = ref<HTMLElement | null>(null)

useScrollToActive(
  containerRef,
  computed(() => props.activeId),
  computed(() => props.items),
)

function onScrollTo(id: string) {
  emit('scrollTo', id)
}
</script>

<template>
  <div ref="containerRef" class="h-full overflow-y-auto py-3 px-2">
    <nav class="space-y-0.5">
      <button
        v-for="item in items"
        :key="item.id"
        :data-toc-id="item.id"
        class="block w-full border-l-2 py-1.5 pr-2 text-left text-base leading-6 transition-colors duration-150"
        :class="[
          levelPaddingClass(item.level),
          activeId === item.id
            ? 'border-primary text-primary font-normal'
            : 'border-transparent text-foreground/70 font-light hover:text-foreground hover:border-muted-foreground/40',
        ]"
        @click="onScrollTo(item.id)"
      >
        <span class="block break-words">
          {{ item.text }}
        </span>
      </button>
    </nav>
  </div>
</template>
