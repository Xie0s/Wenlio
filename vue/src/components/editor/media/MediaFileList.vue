<!--
  MediaFileList.vue - 媒体文件列表组件
  职责：以紧凑列表形式展示非图片类媒体文件，提供操作按钮和使用状态指示
  对外暴露：Props: items, deleting, usageLoaded, usageMap
            Events: delete, open-usage
-->
<script setup lang="ts">
import type { MediaItem, MediaUsageMap } from '@/lib/media'
import { formatBytes } from '@/lib/media'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Separator } from '@/components/ui/separator'
import {
  Trash2, ExternalLink, Cloud, HardDrive,
  FileText, FileCode, FileArchive, File,
  Unlink, BookOpen,
} from 'lucide-vue-next'

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

function openInNewTab(url: string) {
  window.open(url, '_blank', 'noopener')
}
</script>

<template>
  <div class="space-y-0.5">
    <div v-if="items.length === 0" class="flex flex-col items-center py-10 text-muted-foreground">
      <File class="size-8 mb-2 opacity-30" :stroke-width="1" />
      <p class="text-sm">暂无文件</p>
    </div>

    <template v-for="(item, index) in items" :key="item.id">
    <Separator v-if="index > 0" class="mx-4" />
    <div
      class="group flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-muted/50 transition-colors"
    >
      <!-- 文件图标 -->
      <div class="size-9 rounded-lg bg-muted flex items-center justify-center shrink-0">
        <component :is="fileIcon(item.mime_type)" class="size-4 text-muted-foreground" :stroke-width="1.5" />
      </div>

      <!-- 文件信息 -->
      <div class="flex-1 min-w-0">
        <p class="text-sm truncate" :title="item.file_name">{{ item.file_name }}</p>
        <div class="flex items-center gap-1.5 flex-wrap text-[10px] text-muted-foreground mt-0.5">
          <span>{{ formatBytes(item.file_size) }}</span>
          <span class="opacity-40">·</span>
          <span class="font-mono">{{ fileExt(item.file_name) }}</span>
          <Badge variant="secondary" class="text-[10px] px-1 py-0 h-4 gap-0.5">
            <Cloud v-if="item.storage_type === 'cloud'" class="size-2" :stroke-width="1.5" />
            <HardDrive v-else class="size-2" :stroke-width="1.5" />
            {{ item.storage_type === 'cloud' ? '云' : '本地' }}
          </Badge>
          <template v-if="usageLoaded">
            <button
              v-if="usageMap[item.id]?.length"
              class="flex items-center gap-0.5 text-primary hover:underline"
              @click="emit('open-usage', item)"
            >
              <BookOpen class="size-2.5" :stroke-width="1.5" />
              {{ usageMap[item.id]?.length }}处引用
            </button>
            <span v-else class="flex items-center gap-0.5 text-amber-500">
              <Unlink class="size-2.5" :stroke-width="1.5" />
              未使用
            </span>
          </template>
        </div>
      </div>

      <!-- 操作按钮（悬浮显示） -->
      <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost" size="icon"
              class="size-7 rounded-full"
              @click="openInNewTab(item.file_url)"
            >
              <ExternalLink class="size-3.5" :stroke-width="1.5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>新标签页打开</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost" size="icon"
              class="size-7 rounded-full text-destructive hover:text-destructive hover:bg-destructive/10"
              :disabled="deleting === item.id"
              @click="emit('delete', item.id)"
            >
              <Trash2 class="size-3.5" :stroke-width="1.5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>删除文件</TooltipContent>
        </Tooltip>
      </div>
    </div>
    </template>
  </div>
</template>
