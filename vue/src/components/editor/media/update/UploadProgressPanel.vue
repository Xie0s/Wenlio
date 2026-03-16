<!--
  UploadProgressPanel.vue - 上传进度浮窗面板
  职责：右下角固定浮窗，显示全局上传任务列表与进度，支持拖拽移动、最小化、关闭
  对外暴露：无 props/emits，通过 upload store 驱动
-->
<template>
  <Teleport to="body">
    <TooltipProvider>
    <!-- 最小化气泡 -->
    <Transition name="fade-scale">
      <div
        v-if="uploadStore.showPanel && uploadStore.tasks.length > 0 && uploadStore.isMinimized"
        class="fixed z-50 bottom-6 right-6 cursor-pointer"
        @click="uploadStore.toggleMinimize()"
      >
        <Tooltip>
          <TooltipTrigger as-child>
            <div class="relative">
              <Button variant="outline" size="icon" class="rounded-full h-11 w-11 shadow-lg border-primary/30">
                <Upload class="h-4.5 w-4.5 text-primary" />
              </Button>
              <!-- 进行中数量角标 -->
              <span
                v-if="uploadStore.hasActiveTasks"
                class="absolute -top-1 -right-1 flex items-center justify-center h-5 min-w-5 px-1 text-[10px] font-medium bg-primary text-primary-foreground rounded-full"
              >
                {{ uploadStore.activeTasks.length + uploadStore.pendingTasks.length }}
              </span>
            </div>
          </TooltipTrigger>
          <TooltipContent side="left">
            <span v-if="uploadStore.hasActiveTasks">上传中 {{ uploadStore.totalProgress }}%</span>
            <span v-else>上传面板</span>
          </TooltipContent>
        </Tooltip>
      </div>
    </Transition>

    <!-- 展开面板 -->
    <Transition name="fade-scale">
      <div
        v-if="uploadStore.showPanel && uploadStore.tasks.length > 0 && !uploadStore.isMinimized"
        class="fixed z-50 w-80 glass rounded-3xl overflow-hidden"
        :style="panelStyle"
      >
        <!-- 拖拽区域 -->
        <div class="h-2 cursor-move select-none" @mousedown="onMouseDown" />

        <!-- 任务列表 -->
        <div class="overflow-y-auto max-h-[320px] px-2 scrollbar-visible">
          <div
            v-for="task in uploadStore.tasks"
            :key="task.id"
            class="group flex items-center gap-2 px-2 py-2 rounded-2xl hover:bg-muted/40 transition-colors"
          >
            <!-- 缩略图/图标 -->
            <div class="flex-shrink-0 w-8 h-8 rounded-lg overflow-hidden bg-muted flex items-center justify-center">
              <img
                v-if="task.file.type.startsWith('image/') && taskThumbs[task.id]"
                :src="taskThumbs[task.id]"
                class="w-full h-full object-cover"
              />
              <FileIcon v-else class="w-4 h-4 text-muted-foreground" />
            </div>

            <!-- 文件信息 -->
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-2">
                <p class="text-sm truncate" :title="task.file.name">{{ task.file.name }}</p>
                <!-- 状态图标/按钮 -->
                <div class="flex-shrink-0">
                  <CheckCircle2 v-if="task.status === 'success'" class="w-4 h-4 text-green-500" />
                  <Button
                    v-else-if="task.status === 'error'"
                    variant="ghost" size="icon-sm"
                    class="h-6 w-6 text-destructive"
                    @click.stop="uploadStore.retryTask(task.id)"
                    title="重试"
                  >
                    <RotateCw class="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    v-else
                    variant="ghost" size="icon-sm"
                    class="h-6 w-6 opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-destructive"
                    @click.stop="uploadStore.cancelTask(task.id)"
                    title="取消"
                  >
                    <X class="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
              <!-- 进度行 -->
              <div class="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
                <span>{{ formatFileSize(task.file.size) }}</span>
                <span class="opacity-50">·</span>
                <span :class="getStatusColor(task.status)">{{ getStatusText(task) }}</span>
                <template v-if="task.status === 'uploading' && task.speed > 0">
                  <span class="opacity-50">·</span>
                  <span class="tabular-nums">{{ formatSpeed(task.speed) }}</span>
                </template>
              </div>
              <!-- 错误原因 -->
              <p v-if="task.status === 'error' && task.error" class="text-xs text-destructive mt-0.5 truncate" :title="task.error">{{ task.error }}</p>
              <!-- 进度条 -->
              <div v-if="task.status === 'uploading'" class="h-1 w-full bg-muted rounded-full overflow-hidden mt-1">
                <div class="h-full bg-primary transition-all duration-300" :style="{ width: `${task.progress}%` }" />
              </div>
            </div>
          </div>
        </div>

        <!-- 底部操作按钮组 -->
        <div class="px-3 pb-3 pt-2 flex justify-center items-center gap-1.5">
          <!-- 进度指示 -->
          <Tooltip v-if="uploadStore.hasActiveTasks">
            <TooltipTrigger as-child>
              <div class="text-xs text-primary font-medium tabular-nums mr-1 cursor-default">
                {{ uploadStore.totalProgress }}%
              </div>
            </TooltipTrigger>
            <TooltipContent>总上传进度</TooltipContent>
          </Tooltip>
          <!-- 任务数量 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <div class="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded mr-auto cursor-default">
                {{ uploadStore.activeTasks.length + uploadStore.pendingTasks.length }}/{{ uploadStore.tasks.length }}
              </div>
            </TooltipTrigger>
            <TooltipContent>进行中 / 总任务数</TooltipContent>
          </Tooltip>
          <!-- 清除已完成 -->
          <Tooltip v-if="uploadStore.completedTasks.length > 0 || uploadStore.failedTasks.length > 0">
            <TooltipTrigger as-child>
              <Button
                variant="ghost" size="icon-sm"
                class="rounded-full border border-border"
                @click="uploadStore.removeCompletedTasks()"
              >
                <Eraser class="h-3.5 w-3.5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>清除已完成和失败的任务</TooltipContent>
          </Tooltip>
          <!-- 最小化 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                variant="ghost" size="icon-sm"
                class="rounded-full border border-border"
                @click="uploadStore.toggleMinimize()"
              >
                <Minus class="h-3.5 w-3.5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>最小化面板</TooltipContent>
          </Tooltip>
          <!-- 关闭 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                variant="ghost" size="icon-sm"
                class="rounded-full border border-border text-destructive hover:text-destructive"
                @click="handleClose"
              >
                <X class="h-3.5 w-3.5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>关闭（上传中时最小化）</TooltipContent>
          </Tooltip>
        </div>
      </div>
    </Transition>
    </TooltipProvider>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import {
  X, Minus, Eraser, CheckCircle2, RotateCw, Upload, FileIcon,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider } from '@/components/ui/tooltip'
import { useUploadStore, type UploadTask } from '@/stores/upload'
import { formatBytes } from '@/lib/media'

const uploadStore = useUploadStore()

// ── 面板定位（固定右下角，可拖拽） ──

const bottom = ref(24)
const right = ref(24)
let dragging = false
let startX = 0
let startY = 0
let startRight = 0
let startBottom = 0

const panelStyle = computed(() => ({
  bottom: `${bottom.value}px`,
  right: `${right.value}px`,
}))

function onMouseDown(e: MouseEvent) {
  dragging = true
  startX = e.clientX
  startY = e.clientY
  startRight = right.value
  startBottom = bottom.value
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  e.preventDefault()
}

function onMouseMove(e: MouseEvent) {
  if (!dragging) return
  right.value = Math.max(0, startRight - (e.clientX - startX))
  bottom.value = Math.max(0, startBottom + (startY - e.clientY))
}

function onMouseUp() {
  dragging = false
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
}

onUnmounted(() => {
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
})

// ── 图片缩略图（用 URL.createObjectURL） ──

const taskThumbs = ref<Record<string, string>>({})

watch(
  () => uploadStore.tasks.map(t => t.id),
  (ids) => {
    const existing = new Set(Object.keys(taskThumbs.value))
    for (const task of uploadStore.tasks) {
      if (!existing.has(task.id) && task.file.type.startsWith('image/')) {
        taskThumbs.value[task.id] = URL.createObjectURL(task.file)
      }
    }
    // 清理已移除任务的缩略图
    const currentIds = new Set(ids)
    for (const id of existing) {
      if (!currentIds.has(id) && taskThumbs.value[id]) {
        URL.revokeObjectURL(taskThumbs.value[id]!)
        delete taskThumbs.value[id]
      }
    }
  },
  { deep: true },
)

onUnmounted(() => {
  Object.values(taskThumbs.value).forEach(url => URL.revokeObjectURL(url))
})

// ── 工具函数 ──

function formatFileSize(bytes: number): string {
  return formatBytes(bytes)
}

function getStatusText(task: UploadTask): string {
  switch (task.status) {
    case 'pending': return '等待中'
    case 'uploading': return `${task.progress}%`
    case 'success': return '完成'
    case 'error': return '失败'
    case 'cancelled': return '已取消'
    default: return ''
  }
}

function getStatusColor(status: string): string {
  switch (status) {
    case 'uploading': return 'text-primary'
    case 'success': return 'text-green-600'
    case 'error': return 'text-destructive'
    case 'cancelled': return 'text-muted-foreground'
    default: return 'text-muted-foreground'
  }
}

function formatSpeed(bytesPerSecond: number): string {
  if (bytesPerSecond === 0) return '0 B/s'
  if (bytesPerSecond < 1024) return `${bytesPerSecond.toFixed(0)} B/s`
  if (bytesPerSecond < 1024 * 1024) return `${(bytesPerSecond / 1024).toFixed(1)} KB/s`
  return `${(bytesPerSecond / 1024 / 1024).toFixed(1)} MB/s`
}

function handleClose() {
  if (uploadStore.hasActiveTasks) {
    uploadStore.toggleMinimize()
  } else {
    uploadStore.closePanel()
    if (uploadStore.completedTasks.length > 0) {
      uploadStore.removeCompletedTasks()
    }
  }
}
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
