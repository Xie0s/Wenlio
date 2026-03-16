<!--
  UploadPreviewDialog.vue - 上传文件预览确认弹窗
  职责：居中浮窗，用户选择文件后预览列表，支持删除不需要的项，确认后开始上传
  对外暴露：Props: open (v-model), accept（文件类型过滤）
             Emits: confirm(files)
-->
<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-150 ease-out"
      leave-active-class="transition-all duration-100 ease-in"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        ref="panelRef"
        class="fixed z-50 flex flex-col bg-background border rounded-3xl shadow-2xl overflow-hidden"
        style="width: 520px; max-height: 580px;"
        :style="floatingStyle"
      >
        <!-- 拖拽标题栏 -->
        <div
          class="flex items-center gap-2.5 px-5 py-3.5 border-b cursor-move select-none shrink-0"
          @mousedown="onMouseDown"
        >
          <Upload class="size-4 text-muted-foreground shrink-0" :stroke-width="1.5" />
          <h3 class="text-sm font-medium flex-1">{{ title }}</h3>
          <span v-if="pendingFiles.length > 0" class="text-xs text-primary font-medium">
            {{ pendingFiles.length }} 个文件
          </span>
          <button
            class="size-7 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted/60 transition-colors"
            @click="handleClose"
          >
            <X class="size-4" :stroke-width="1.5" />
          </button>
        </div>

        <!-- 文件列表 -->
        <div class="flex-1 overflow-y-auto px-3 py-2 scrollbar-visible min-h-0">
          <!-- 空状态：提示选择文件 -->
          <div
            v-if="pendingFiles.length === 0"
            class="flex flex-col items-center justify-center py-12 text-muted-foreground cursor-pointer"
            @click="triggerFileInput"
          >
            <div class="w-16 h-16 rounded-2xl border-2 border-dashed border-muted-foreground/30 flex items-center justify-center mb-3">
              <Plus class="size-6 text-muted-foreground/50" />
            </div>
            <p class="text-sm">点击选择文件</p>
            <p class="text-xs mt-1 opacity-60">或拖拽文件到此处</p>
          </div>

          <!-- 文件列表 -->
          <div v-else class="divide-y divide-border">
            <div
              v-for="(item, index) in pendingFiles"
              :key="item.id"
              class="group flex items-center gap-3 px-3 py-1.5 hover:bg-muted/40 transition-colors"
            >
              <!-- 缩略图/图标 -->
              <div class="flex-shrink-0 w-8 h-8 rounded-lg overflow-hidden bg-muted flex items-center justify-center">
                <img
                  v-if="item.thumb"
                  :src="item.thumb"
                  class="w-full h-full object-cover"
                />
                <FileIcon v-else class="w-4 h-4 text-muted-foreground" />
              </div>

              <!-- 文件信息 -->
              <div class="min-w-0 flex-1">
                <p class="text-sm truncate leading-tight" :title="item.file.name">{{ item.file.name }}</p>
                <div class="flex items-center gap-1 mt-0.5">
                  <span class="text-xs" :class="item.file.size > MAX_FILE_SIZE ? 'text-destructive' : 'text-muted-foreground'">{{ formatBytes(item.file.size) }}</span>
                  <template v-if="item.file.size > MAX_FILE_SIZE">
                    <AlertTriangle class="w-3 h-3 text-destructive" />
                    <span class="text-xs text-destructive">超出限制（50 MB）</span>
                  </template>
                </div>
              </div>

              <!-- 删除按钮 -->
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button
                    variant="ghost" size="icon-sm"
                    class="h-7 w-7 rounded-full text-muted-foreground hover:text-destructive"
                    @click="removeFile(index)"
                  >
                    <X class="w-3.5 h-3.5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>移除</TooltipContent>
              </Tooltip>
            </div>
          </div>
        </div>

        <!-- 底部操作栏 -->
        <div class="flex items-center px-5 py-3 border-t shrink-0">
          <span class="flex-1 text-xs text-muted-foreground">
            <template v-if="pendingFiles.length > 0">
              共 {{ formatBytes(totalSize) }}
            </template>
          </span>
          <div class="flex items-center gap-2">
            <!-- 继续添加文件 -->
            <Tooltip v-if="pendingFiles.length > 0">
              <TooltipTrigger as-child>
                <Button variant="outline" size="icon" class="size-8 rounded-full" @click="triggerFileInput">
                  <Plus class="size-4" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>继续添加</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="outline" size="icon" class="size-8 rounded-full" @click="handleClose">
                  <X class="size-4" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>取消</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button size="icon" class="size-8 rounded-full" :disabled="pendingFiles.length === 0" @click="handleConfirm">
                  <Check class="size-4" :stroke-width="2" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>开始上传</TooltipContent>
            </Tooltip>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 隐藏文件输入框 -->
    <input
      ref="fileInputRef"
      type="file"
      :accept="accept"
      multiple
      class="hidden"
      @change="handleFileSelect"
    />
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { Upload, X, Plus, Check, FileIcon, AlertTriangle } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import { useFloatingPanel } from '@/composables/useFloatingPanel'
import { formatBytes } from '@/lib/media'

interface PendingFile {
  id: string
  file: File
  thumb: string
}

const props = withDefaults(defineProps<{
  accept?: string
  title?: string
  autoOpenPicker?: boolean
  initialFiles?: File[]
  appendFiles?: File[]
}>(), {
  accept: '',
  title: '上传文件',
  autoOpenPicker: true,
  initialFiles: () => [],
  appendFiles: () => [],
})

const MAX_FILE_SIZE = 50 * 1024 * 1024

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  confirm: [files: File[]]
}>()

const { panelRef, style: floatingStyle, onMouseDown, reset } = useFloatingPanel()
const fileInputRef = ref<HTMLInputElement | null>(null)
const pendingFiles = ref<PendingFile[]>([])
let idCounter = 0

// ── 打开时重置状态并自动弹出文件选择器 ──

watch(open, (val) => {
  if (val) {
    ;(document.activeElement as HTMLElement)?.blur()
    clearPendingFiles()
    reset()
    if (props.initialFiles.length > 0) {
      addFiles([...props.initialFiles])
    } else if (props.autoOpenPicker) {
      setTimeout(() => triggerFileInput(), 100)
    }
  }
}, { flush: 'sync' })

// 对话框已开时，外部追加文件
watch(() => props.appendFiles, (newFiles) => {
  if (!open.value || !newFiles?.length) return
  addFiles([...newFiles])
})

// ── 拖拽支持 ──

// 面板自身也支持拖入文件
watch(panelRef, (el) => {
  if (el) {
    el.addEventListener('dragover', onPanelDragOver)
    el.addEventListener('drop', onPanelDrop)
  }
})

function onPanelDragOver(e: DragEvent) {
  if (e.dataTransfer?.types.includes('Files')) {
    e.preventDefault()
    e.dataTransfer!.dropEffect = 'copy'
  }
}

function onPanelDrop(e: DragEvent) {
  e.preventDefault()
  const files = Array.from(e.dataTransfer?.files ?? [])
  addFiles(files)
}

// ── 文件操作 ──

function triggerFileInput() {
  fileInputRef.value?.click()
}

function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  addFiles(files)
  input.value = ''
}

function addFiles(files: File[]) {
  for (const file of files) {
    const id = `pf_${++idCounter}`
    let thumb = ''
    if (file.type.startsWith('image/')) {
      thumb = URL.createObjectURL(file)
    }
    pendingFiles.value.push({ id, file, thumb })
  }
}

function removeFile(index: number) {
  const item = pendingFiles.value[index]
  if (item?.thumb) URL.revokeObjectURL(item.thumb)
  pendingFiles.value.splice(index, 1)
}

function clearPendingFiles() {
  pendingFiles.value.forEach(item => {
    if (item.thumb) URL.revokeObjectURL(item.thumb)
  })
  pendingFiles.value = []
}

const totalSize = computed(() =>
  pendingFiles.value.reduce((sum, item) => sum + item.file.size, 0),
)

// ── 确认 / 关闭 ──

function handleConfirm() {
  const files = pendingFiles.value.map(item => item.file)
  emit('confirm', files)
  open.value = false
}

function handleClose() {
  clearPendingFiles()
  open.value = false
}

onUnmounted(() => {
  clearPendingFiles()
})
</script>
