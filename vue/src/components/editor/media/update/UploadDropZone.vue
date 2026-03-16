<!--
  UploadDropZone.vue - 拖拽上传区域组件
  职责：包裹内容区域，监听文件拖入事件，触发全局上传 store 添加任务
  对外暴露：Props: accept（文件类型过滤）、onFilesAdded（可选回调）
  使用方式：<UploadDropZone> <slot /> </UploadDropZone>
-->
<template>
  <div
    ref="root"
    class="relative"
    @dragenter.prevent="onDragEnter"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
  >
    <slot />

    <!-- 拖拽遮罩 -->
    <Transition name="fade">
      <div
        v-if="isDragging"
        class="absolute inset-0 z-40 flex items-center justify-center rounded-3xl pointer-events-none"
      >
        <div class="absolute inset-0 bg-primary/5 border-2 border-dashed border-primary/40 rounded-3xl" />
        <div class="relative flex flex-col items-center gap-2 text-primary">
          <Upload class="h-8 w-8" />
          <span class="text-sm font-medium">松开以上传文件</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, useTemplateRef } from 'vue'
import { Upload } from 'lucide-vue-next'
import { useUploadStore, type AddTaskOptions } from '@/stores/upload'

const props = withDefaults(defineProps<{
  accept?: string
  uploadOptions?: AddTaskOptions
}>(), {
  accept: '',
})

const rootRef = useTemplateRef<HTMLDivElement>('root')
const uploadStore = useUploadStore()
const isDragging = ref(false)
let dragCounter = 0

function resetDragState() {
  isDragging.value = false
  dragCounter = 0
}

onMounted(() => {
  // 在捕获阶段注册，先于子元素的 stopPropagation 执行，确保 isDragging 总能被重置
  rootRef.value?.addEventListener('drop', resetDragState, true)
})

onUnmounted(() => {
  rootRef.value?.removeEventListener('drop', resetDragState, true)
})

function isFileTransfer(e: DragEvent): boolean {
  return e.dataTransfer?.types.includes('Files') ?? false
}

function onDragEnter(e: DragEvent) {
  if (!isFileTransfer(e)) return
  dragCounter++
  isDragging.value = true
}

function onDragOver(e: DragEvent) {
  if (!isFileTransfer(e)) return
  e.dataTransfer!.dropEffect = 'copy'
}

function onDragLeave(e: DragEvent) {
  if (!isFileTransfer(e)) return
  dragCounter--
  if (dragCounter <= 0) {
    dragCounter = 0
    isDragging.value = false
  }
}

function filterFiles(files: File[]): File[] {
  if (!props.accept) return files
  const patterns = props.accept.split(',').map(s => s.trim().toLowerCase())
  return files.filter(file => {
    const type = file.type.toLowerCase()
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()
    return patterns.some(p => {
      if (p.endsWith('/*')) return type.startsWith(p.replace('/*', '/'))
      if (p.startsWith('.')) return ext === p
      return type === p
    })
  })
}

function onDrop(e: DragEvent) {
  const rawFiles = Array.from(e.dataTransfer?.files ?? [])
  if (rawFiles.length === 0) return
  const files = props.accept ? filterFiles(rawFiles) : rawFiles
  if (files.length === 0) return
  uploadStore.requestUploadWithPreview(files, props.uploadOptions)
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
