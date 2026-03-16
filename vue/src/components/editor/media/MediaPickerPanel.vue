<!--
  MediaPickerPanel.vue - 媒体文件选择面板（左右分栏）
  职责：在拖拽浮窗内展示全量媒体文件，支持分类过滤与多选；图片以缩略图网格显示
  对外暴露：Props: items(MediaItem[]), selected(Set<string>)
            Events: toggle(id: string)
-->
<script setup lang="ts">
import { ref, computed } from 'vue'
import type { MediaItem } from '@/lib/media'
import { formatBytes } from '@/lib/media'
import { Badge } from '@/components/ui/badge'
import { Check, Cloud, HardDrive, FileText, FileCode, FileArchive, File, Layers, Image } from 'lucide-vue-next'

const props = defineProps<{
  items: MediaItem[]
  selected: Set<string>
}>()

const emit = defineEmits<{
  toggle: [id: string]
}>()

type CategoryKey = 'all' | 'image' | 'document' | 'code' | 'archive' | 'other'

const activeCategory = ref<CategoryKey>('all')

function isImage(item: MediaItem) {
  return item.mime_type?.startsWith('image/') ?? false
}

function isDocument(item: MediaItem) {
  const m = item.mime_type
  return (
    m.startsWith('text/') ||
    m.includes('pdf') ||
    m.includes('word') ||
    m.includes('document') ||
    m.includes('presentation') ||
    m.includes('spreadsheet') ||
    m.includes('excel') ||
    m.includes('powerpoint')
  )
}

function isCode(item: MediaItem) {
  const m = item.mime_type
  return (
    m.includes('javascript') ||
    m.includes('json') ||
    m.includes('xml') ||
    m.includes('html') ||
    m.includes('css') ||
    m.includes('yaml') ||
    m.includes('markdown') ||
    m.includes('typescript')
  )
}

function isArchive(item: MediaItem) {
  const m = item.mime_type
  return (
    m.includes('zip') ||
    m.includes('tar') ||
    m.includes('rar') ||
    m.includes('7z') ||
    m.includes('gzip') ||
    m.includes('bzip') ||
    m.includes('compress')
  )
}

function fileIcon(mime: string) {
  if (mime.startsWith('text/') || mime.includes('document') || mime.includes('pdf')) return FileText
  if (mime.includes('zip') || mime.includes('tar') || mime.includes('rar') || mime.includes('7z') || mime.includes('gzip')) return FileArchive
  if (mime.includes('javascript') || mime.includes('json') || mime.includes('xml') || mime.includes('html') || mime.includes('css')) return FileCode
  return File
}

function fileExt(name: string) {
  const parts = name.split('.')
  return parts.length > 1 ? (parts[parts.length - 1] ?? 'FILE').toUpperCase() : 'FILE'
}

const categories = computed(() => [
  { key: 'all' as CategoryKey, label: '全部', icon: Layers, count: props.items.length },
  { key: 'image' as CategoryKey, label: '图片', icon: Image, count: props.items.filter(isImage).length },
  { key: 'document' as CategoryKey, label: '文档', icon: FileText, count: props.items.filter(isDocument).length },
  { key: 'code' as CategoryKey, label: '代码', icon: FileCode, count: props.items.filter(isCode).length },
  { key: 'archive' as CategoryKey, label: '归档', icon: FileArchive, count: props.items.filter(isArchive).length },
  {
    key: 'other' as CategoryKey,
    label: '其他',
    icon: File,
    count: props.items.filter(i => !isImage(i) && !isDocument(i) && !isCode(i) && !isArchive(i)).length,
  },
])

const filteredItems = computed(() => {
  switch (activeCategory.value) {
    case 'image':    return props.items.filter(isImage)
    case 'document': return props.items.filter(isDocument)
    case 'code':     return props.items.filter(isCode)
    case 'archive':  return props.items.filter(isArchive)
    case 'other':    return props.items.filter(i => !isImage(i) && !isDocument(i) && !isCode(i) && !isArchive(i))
    default:         return props.items
  }
})

const isImageCategory = computed(() => activeCategory.value === 'image')
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <!-- 左侧分类栏 -->
    <aside class="w-44 shrink-0 border-r flex flex-col overflow-hidden">
      <div class="flex-1 overflow-y-auto px-2 pt-2 pb-2 space-y-0.5">
        <button
          v-for="cat in categories"
          :key="cat.key"
          class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-sm transition-colors"
          :class="activeCategory === cat.key
            ? 'bg-primary/10 text-primary font-medium'
            : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'"
          @click="activeCategory = cat.key"
        >
          <component :is="cat.icon" class="size-4 shrink-0" :stroke-width="1.5" />
          <span class="flex-1 text-left truncate">{{ cat.label }}</span>
          <span class="text-xs opacity-60 tabular-nums">{{ cat.count }}</span>
        </button>
      </div>
      <p class="shrink-0 py-2 text-center text-xs">文件类型</p>
    </aside>

    <!-- 右侧内容区 -->
    <div class="flex-1 overflow-y-auto">
      <!-- 空态 -->
      <div v-if="filteredItems.length === 0" class="flex flex-col items-center justify-center h-full py-16 text-muted-foreground">
        <File class="size-10 mb-3 opacity-20" :stroke-width="1" />
        <p class="text-sm">该分类暂无文件</p>
      </div>

      <!-- 图片分类：缩略图网格 -->
      <div v-else-if="isImageCategory" class="p-3 grid grid-cols-3 gap-2.5">
        <button
          v-for="item in filteredItems"
          :key="item.id"
          class="relative group rounded-xl overflow-hidden aspect-square bg-muted focus:outline-none"
          :class="selected.has(item.id) ? 'ring-2 ring-primary ring-offset-2' : ''"
          @click="emit('toggle', item.id)"
        >
          <img
            :src="item.file_url"
            :alt="item.file_name"
            class="w-full h-full object-cover"
            loading="lazy"
          />
          <!-- 悬浮遮罩 -->
          <div
            class="absolute inset-0 transition-colors"
            :class="selected.has(item.id)
              ? 'bg-primary/30'
              : 'bg-black/0 group-hover:bg-black/20'"
          />
          <!-- 选中标记 -->
          <div
            v-if="selected.has(item.id)"
            class="absolute top-2 right-2 size-6 rounded-full bg-primary flex items-center justify-center shadow"
          >
            <Check class="size-3.5 text-primary-foreground" :stroke-width="2.5" />
          </div>
          <!-- 文件名（悬浮显示） -->
          <div class="absolute bottom-0 left-0 right-0 px-2 py-1.5 bg-black/50 backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity">
            <p class="text-xs text-white truncate">{{ item.file_name }}</p>
          </div>
        </button>
      </div>

      <!-- 其他分类：列表 -->
      <div v-else class="p-3 space-y-1">
        <button
          v-for="item in filteredItems"
          :key="item.id"
          class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors text-left group relative"
          :class="selected.has(item.id)
            ? 'bg-primary/10 ring-1 ring-primary/30'
            : 'hover:bg-muted/50'"
          @click="emit('toggle', item.id)"
        >
          <!-- 文件类型图标 -->
          <div
            class="size-9 rounded-lg flex items-center justify-center shrink-0 transition-colors"
            :class="selected.has(item.id) ? 'bg-primary/20' : 'bg-muted'"
          >
            <component :is="fileIcon(item.mime_type)" class="size-4 text-muted-foreground" :stroke-width="1.5" />
          </div>

          <!-- 文件信息 -->
          <div class="flex-1 min-w-0">
            <p class="text-sm truncate" :title="item.file_name">{{ item.file_name }}</p>
            <div class="flex items-center gap-1.5 text-[11px] text-muted-foreground mt-0.5">
              <span>{{ formatBytes(item.file_size) }}</span>
              <span class="opacity-40">·</span>
              <span class="font-mono uppercase">{{ fileExt(item.file_name) }}</span>
              <Badge variant="secondary" class="text-[10px] px-1 py-0 h-4 gap-0.5 ml-0.5">
                <Cloud v-if="item.storage_type === 'cloud'" class="size-2" :stroke-width="1.5" />
                <HardDrive v-else class="size-2" :stroke-width="1.5" />
                {{ item.storage_type === 'cloud' ? '云' : '本地' }}
              </Badge>
            </div>
          </div>

          <!-- 选中标记 -->
          <div
            class="size-5 rounded-full border-2 flex items-center justify-center shrink-0 transition-all"
            :class="selected.has(item.id)
              ? 'bg-primary border-primary'
              : 'border-border group-hover:border-primary/50'"
          >
            <Check v-if="selected.has(item.id)" class="size-3 text-primary-foreground" :stroke-width="2.5" />
          </div>
        </button>
      </div>
    </div>
  </div>
</template>
