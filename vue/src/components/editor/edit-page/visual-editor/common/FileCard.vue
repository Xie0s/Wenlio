<!--
  FileCard.vue - 文件附件卡片基础组件
  职责：纯展示型文件卡片，根据 mimeType 显示对应图标、文件名、格式化大小
  对外暴露：Props: fileName, fileSize, mimeType, selected
-->
<script setup lang="ts">
import { computed } from 'vue'
import {
  File,
  FileText,
  FileImage,
  FileVideo,
  FileAudio,
  FileArchive,
  FileCode,
  FileSpreadsheet,
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  fileName: string
  fileSize?: number | null
  mimeType?: string | null
  selected?: boolean
}>(), {
  fileSize: null,
  mimeType: null,
  selected: false,
})

function resolveIcon(mime: string | null) {
  if (!mime) return File
  if (mime.startsWith('image/')) return FileImage
  if (mime.startsWith('video/')) return FileVideo
  if (mime.startsWith('audio/')) return FileAudio
  if (mime === 'application/pdf') return FileText
  if (
    mime.includes('zip') || mime.includes('tar') || mime.includes('gz') ||
    mime.includes('rar') || mime.includes('7z') ||
    mime === 'application/x-compressed' || mime === 'application/x-archive'
  ) return FileArchive
  if (
    mime.includes('spreadsheet') || mime.includes('excel') || mime === 'text/csv'
  ) return FileSpreadsheet
  if (
    mime.includes('javascript') || mime.includes('json') || mime.includes('html') ||
    mime.includes('xml') || mime.includes('css') || mime.includes('sql')
  ) return FileCode
  if (mime.startsWith('text/') || mime.includes('document') || mime.includes('word')) return FileText
  return File
}

function formatSize(bytes: number | null): string {
  if (bytes == null) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

function resolveExt(name: string): string {
  const dot = name.lastIndexOf('.')
  return dot > 0 ? name.slice(dot + 1).toUpperCase() : ''
}

const iconComponent = computed(() => resolveIcon(props.mimeType))
const sizeLabel = computed(() => formatSize(props.fileSize))
const extLabel = computed(() => resolveExt(props.fileName))
</script>

<template>
  <div class="file-card" :class="{ 'file-card--selected': selected }">
    <!-- 左侧：类型图标 + 扩展名 -->
    <div class="file-card__icon-area">
      <component :is="iconComponent" class="file-card__icon" />
      <span v-if="extLabel" class="file-card__ext">{{ extLabel }}</span>
    </div>

    <!-- 右侧：文件名 + 大小 -->
    <div class="file-card__body">
      <span class="file-card__name" :title="fileName">{{ fileName }}</span>
      <span v-if="sizeLabel" class="file-card__size">{{ sizeLabel }}</span>
    </div>
  </div>
</template>

<style scoped>
.file-card {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border);
  border-radius: 0.625rem;
  background-color: var(--muted);
  color: var(--foreground);
  cursor: default;
  user-select: none;
  min-width: 0;
}

.file-card--selected {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

.file-card__icon-area {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.125rem;
  width: 2rem;
}

.file-card__icon {
  width: 1.25rem;
  height: 1.25rem;
  color: var(--muted-foreground);
}

.file-card__ext {
  font-size: 0.5625rem;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--muted-foreground);
  line-height: 1;
  letter-spacing: 0.02em;
}

.file-card__body {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  min-width: 0;
  flex: 1;
}

.file-card__name {
  font-size: 0.8125rem;
  font-weight: 500;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--foreground);
}

.file-card__size {
  font-size: 0.6875rem;
  color: var(--muted-foreground);
  line-height: 1.2;
}
</style>
