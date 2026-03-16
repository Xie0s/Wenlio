<!--
  MediaSidebar.vue - 媒体文件右侧栏组件
  职责：作为页面右侧面板，内嵌 MediaFileList 展示非图片类文件
  对外暴露：Props: items, deleting, usageLoaded, usageMap
            Events: delete, open-usage
-->
<script setup lang="ts">
import type { MediaItem, MediaUsageMap } from '@/lib/media'
import MediaFileList from './MediaFileList.vue'

defineProps<{
  items: MediaItem[]
  deleting: string | null
  usageLoaded: boolean
  usageMap: MediaUsageMap
}>()

const emit = defineEmits<{
  delete: [id: string]
  'open-usage': [item: MediaItem]
}>()
</script>

<template>
  <aside class="w-80 shrink-0 flex flex-col overflow-hidden">
    <!-- 顶部标题 -->
    <p class="shrink-0 py-2 text-center text-sm">文件 ({{ items.length }})</p>
    <!-- 文件列表（flex-1 填满剩余高度，独立滚动） -->
    <div class="flex-1 min-h-0">
      <div class="h-full overflow-y-auto overflow-x-hidden pl-4 pr-2 pt-2 scrollbar-visible">
        <MediaFileList
          :items="items"
          :deleting="deleting"
          :usage-loaded="usageLoaded"
          :usage-map="usageMap"
          @delete="emit('delete', $event)"
          @open-usage="emit('open-usage', $event)"
        />
      </div>
    </div>
  </aside>
</template>
