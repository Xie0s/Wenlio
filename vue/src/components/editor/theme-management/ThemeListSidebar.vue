<!-- ThemeListSidebar.vue - 主题列表页侧边栏
     职责：提供固定定位的侧边栏容器，包裹分类/标签管理面板 + 底部工具栏（护眼 + 主题切换）
     对外暴露：
       Props: catLib(CategoryListLib), tagLib(TagListLib), eyeCareMode(boolean)
       Emits: toggle-eye-care, refresh -->
<script setup lang="ts">
import type { CategoryListLib, TagListLib } from '@/lib/category-tag'
import CategoryTagSidebar from './categorytag/CategoryTagSidebar.vue'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Leaf, Loader } from 'lucide-vue-next'

defineProps<{
  catLib: CategoryListLib
  tagLib: TagListLib
  eyeCareMode: boolean
}>()

defineEmits<{
  'toggle-eye-care': []
  refresh: []
}>()
</script>

<template>
  <aside class="relative z-10 w-80 shrink-0 h-full flex flex-col overflow-visible">
    <!-- 分类/标签管理面板 -->
    <div class="relative flex-1 min-h-0 overflow-visible">
      <div class="h-full overflow-y-auto pr-5">
        <CategoryTagSidebar :cat-lib="catLib" :tag-lib="tagLib" />
      </div>
    </div>

    <!-- 底部工具栏 -->
    <div class="shrink-0 border-t py-4 flex items-center justify-center gap-3">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="rounded-full h-11 w-11 transition-transform hover:scale-110"
            @click="$emit('refresh')">
            <Loader class="size-6" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">刷新数据</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="rounded-full h-11 w-11 transition-transform hover:scale-110"
            :class="eyeCareMode ? 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-950/30' : ''"
            @click="$emit('toggle-eye-care')">
            <Leaf class="size-6" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">{{ eyeCareMode ? '关闭护眼' : '护眼模式' }}</TooltipContent>
      </Tooltip>
    </div>
  </aside>
</template>
