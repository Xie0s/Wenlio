<!--
  GalleryCategoryItem.vue - 主题画廊分类节点（只读/浏览版）
  职责：渲染单个分类节点，支持展开/折叠子分类和点击筛选，无 CRUD 操作
  设计：对齐 CategorySidebarItem 的视觉风格（字体、缩进、选中状态、VS Code 竖线）
  对外暴露：
    Props: cat(Category), selectedCategoryId(string)
    Emits: select-category(id: string)
-->
<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Category } from '@/utils/types'
import { ChevronRight } from 'lucide-vue-next'

const props = defineProps<{
  cat: Category
  selectedCategoryId: string
}>()

const emit = defineEmits<{
  'select-category': [id: string]
}>()

const hasChildren = computed(() => (props.cat.children?.length ?? 0) > 0)
const isMainCategory = computed(() => (props.cat.level ?? 0) === 0)
const isOpen = ref(false)

function handleSelect(id: string) {
  emit('select-category', id)
  if (isMainCategory.value && hasChildren.value) {
    isOpen.value = true
  }
}

const totalCount = computed(() => {
  const own = props.cat.theme_count ?? 0
  if ((props.cat.level ?? 0) === 0 && props.cat.children?.length) {
    return own + props.cat.children.reduce((sum, c) => sum + (c.theme_count ?? 0), 0)
  }
  return own
})
</script>

<template>
  <div>
    <!-- 分类行 -->
    <div class="group/cat relative flex items-center gap-1.5 px-3 py-2.5 transition-colors cursor-pointer" :class="[
      (cat.level ?? 0) === 0 ? 'hover:bg-accent/50' : '',
      selectedCategoryId === cat.id ? 'bg-accent/40' : (isOpen && hasChildren && (cat.level ?? 0) === 0 ? 'bg-accent/50' : '')
    ]" @click="hasChildren && (isOpen = !isOpen)">
      <button class="flex-1 min-w-0 flex items-start gap-1.5 text-lg font-thin text-left transition-colors" :class="selectedCategoryId === cat.id
        ? 'text-primary font-extralight'
        : 'text-foreground hover:text-primary/80'" @click.stop="handleSelect(cat.id)">
        <span class="whitespace-pre-wrap break-words">{{ cat.name }}</span>
        <span v-if="totalCount > 0"
          class="shrink-0 inline-flex items-center justify-center h-5 min-w-5 px-1.5 rounded-full bg-muted text-xs tabular-nums text-muted-foreground">
          {{ totalCount }}
        </span>
      </button>

      <button v-if="hasChildren"
        class="p-2 rounded-full text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
        @click.stop="isOpen = !isOpen">
        <ChevronRight class="h-4 w-4 transition-transform duration-150" :class="{ 'rotate-90': isOpen }" />
      </button>
    </div>

    <!-- 子分类（CSS Grid 展开/折叠动画） -->
    <div v-if="hasChildren" class="grid transition-[grid-template-rows] duration-200 ease-in-out"
      :class="isOpen ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'">
      <div class="overflow-hidden">
        <div class="relative pl-4 space-y-0.5 transition-opacity duration-200"
          :class="isOpen ? 'opacity-100' : 'opacity-0'">
          <div class="absolute left-[9px] top-0 bottom-0 w-px bg-border/60" />
          <div v-for="child in cat.children" :key="child.id" class="relative -ml-[7px] transition-colors duration-150"
            :class="selectedCategoryId === child.id ? 'bg-accent/40' : 'hover:bg-accent/60'">
            <span v-if="selectedCategoryId === child.id"
              class="absolute left-0 inset-y-0 w-px bg-primary pointer-events-none" />
            <GalleryCategoryItem :cat="child" :selected-category-id="selectedCategoryId"
              @select-category="emit('select-category', $event)" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
