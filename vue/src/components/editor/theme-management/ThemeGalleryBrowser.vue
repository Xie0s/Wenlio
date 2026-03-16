<!--
  ThemeGalleryBrowser.vue - 主题画廊侧边栏浏览器（只读/浏览版）
  职责：读者端主题画廊侧边栏，提供分类树浏览（点击筛选）和标签多选筛选功能，无 CRUD 操作
  设计：对齐 CategoryTagSidebar 的视觉风格（区块折叠、标题字体、分隔线、颜色体系）
  对外暴露：
    Props: categories, categoriesLoading, selectedCategoryId, tags, tagsLoading, selectedTagSlugs
    Emits: select-category(id: string), toggle-tag(slug: string)
-->
<script setup lang="ts">
import { ref } from 'vue'
import type { Category, Tag } from '@/utils/types'
import GalleryCategoryItem from './categorytag/GalleryCategoryItem.vue'
import { Loader2, ChevronRight } from 'lucide-vue-next'

defineProps<{
  categories: Category[]
  categoriesLoading: boolean
  selectedCategoryId: string
  tags: Tag[]
  tagsLoading: boolean
  selectedTagSlugs: string[]
}>()

const emit = defineEmits<{
  'select-category': [id: string]
  'toggle-tag': [slug: string]
}>()

const catExpanded = ref(false)
const tagExpanded = ref(false)
</script>

<template>
  <nav class="h-full select-none overflow-y-auto">

    <!-- ════════ 分类区块 ════════ -->
    <div>

      <!-- 标题行（可折叠） -->
      <div class="flex items-center justify-between px-3 py-3 sticky top-0 z-10 transition-colors cursor-pointer"
        :class="catExpanded ? 'bg-accent' : 'bg-background hover:bg-accent'" @click="catExpanded = !catExpanded">
        <div class="flex items-center gap-2.5">
          <span class="text-xl font-light text-foreground">全部分类</span>
          <span v-if="categories.length > 0" class="text-sm text-muted-foreground tabular-nums">
            {{ categories.length }}
          </span>
        </div>
        <ChevronRight class="h-4 w-4 text-muted-foreground transition-transform duration-200"
          :class="catExpanded ? 'rotate-90' : ''" />
      </div>

      <!-- 分类内容（展开时） -->
      <div v-show="catExpanded" class="py-1 pb-3 transition-all duration-200"
        :class="catExpanded ? 'opacity-100' : 'opacity-0 h-0 overflow-hidden'">
        <div v-if="categoriesLoading" class="flex justify-center py-8">
          <Loader2 class="h-3.5 w-3.5 animate-spin text-muted-foreground/40" />
        </div>
        <p v-else-if="categories.length === 0" class="pt-5 text-center text-xs text-muted-foreground/35">
          暂无分类
        </p>
        <template v-else>
          <GalleryCategoryItem v-for="cat in categories" :key="cat.id" :cat="cat"
            :selected-category-id="selectedCategoryId" @select-category="emit('select-category', $event)" />
        </template>
      </div>
    </div>

    <!-- 分隔线 -->
    <div class="flex items-center mx-4 my-2">
      <div class="flex-1 h-px bg-gradient-to-l from-border/70 to-transparent" />
      <div class="w-1 h-1 rounded-full bg-border/60 mx-2" />
      <div class="flex-1 h-px bg-gradient-to-r from-border/70 to-transparent" />
    </div>

    <!-- ════════ 标签区块 ════════ -->
    <div>

      <!-- 标题行（可折叠） -->
      <div class="flex items-center justify-between px-3 py-3 sticky top-0 z-10 transition-colors cursor-pointer"
        :class="tagExpanded ? 'bg-accent' : 'bg-background hover:bg-accent'" @click="tagExpanded = !tagExpanded">
        <div class="flex items-center gap-2.5">
          <span class="text-xl font-light text-foreground">全部标签</span>
          <span v-if="tags.length > 0" class="text-sm text-muted-foreground tabular-nums">
            {{ tags.length }}
          </span>
        </div>
        <ChevronRight class="h-4 w-4 text-muted-foreground transition-transform duration-200"
          :class="tagExpanded ? 'rotate-90' : ''" />
      </div>

      <!-- 标签内容（展开时） -->
      <div v-show="tagExpanded" class="px-3 pb-3 transition-all duration-200"
        :class="tagExpanded ? 'opacity-100' : 'opacity-0 h-0 overflow-hidden'">
        <div v-if="tagsLoading" class="flex justify-center py-6">
          <Loader2 class="h-3.5 w-3.5 animate-spin text-muted-foreground/40" />
        </div>
        <p v-else-if="tags.length === 0" class="pt-5 text-center text-xs text-muted-foreground/35">
          暂无标签
        </p>
        <div v-else class="pt-2.5">
          <div class="flex flex-wrap gap-2">
            <button v-for="tag in tags" :key="tag.id"
              class="inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-colors"
              :class="selectedTagSlugs.includes(tag.slug)
                ? 'border-emerald-500/50 bg-emerald-500/15 text-emerald-700 font-medium dark:text-emerald-300'
                : 'border-emerald-200/60 bg-emerald-50/70 text-emerald-700 hover:bg-emerald-100/80 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-300 dark:hover:bg-emerald-500/20'"
              @click="emit('toggle-tag', tag.slug)">
              {{ tag.name }}
              <span v-if="tag.usage_count" class="tabular-nums text-xs opacity-60">{{ tag.usage_count }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

  </nav>
</template>
