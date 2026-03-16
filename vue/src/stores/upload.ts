/**
 * stores/upload.ts - 全局上传任务管理 Store
 * 职责：管理文件上传任务队列、进度追踪、并发控制、取消/重试
 * 对外暴露：useUploadStore()、UploadTask 类型
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { http, type UploadProgressOptions } from '@/utils/http'
import { normalizeUploadedUrl, formatBytes, type MediaItem } from '@/lib/media'

const MAX_FILE_SIZE = 50 * 1024 * 1024 // 与后端保持一致：50 MB

// ── 类型定义 ──

export type UploadStatus = 'pending' | 'uploading' | 'success' | 'error' | 'cancelled'

export interface UploadTask {
  id: string
  file: File
  status: UploadStatus
  progress: number
  speed: number
  result: MediaItem | null
  error: string
  abortController: AbortController | null
  onComplete?: (media: MediaItem) => void
}

export interface AddTaskOptions {
  onComplete?: (media: MediaItem) => void
}

// ── 工具函数 ──

let taskIdCounter = 0
function genTaskId(): string {
  return `upload_${Date.now()}_${++taskIdCounter}`
}

// ── Store 定义 ──

export const useUploadStore = defineStore('upload', () => {
  const tasks = ref<UploadTask[]>([])
  const showPanel = ref(false)
  const isMinimized = ref(false)

  const pendingPreviewFiles = ref<File[]>([])
  const pendingPreviewOptions = ref<AddTaskOptions>({})
  const pendingAppendFiles = ref<File[]>([])
  const showPreviewDialog = ref(false)

  const MAX_CONCURRENT = 3

  // ── 计算属性 ──

  const activeTasks = computed(() =>
    tasks.value.filter(t => t.status === 'uploading'),
  )

  const pendingTasks = computed(() =>
    tasks.value.filter(t => t.status === 'pending'),
  )

  const completedTasks = computed(() =>
    tasks.value.filter(t => t.status === 'success'),
  )

  const failedTasks = computed(() =>
    tasks.value.filter(t => t.status === 'error'),
  )

  const hasActiveTasks = computed(() =>
    activeTasks.value.length > 0 || pendingTasks.value.length > 0,
  )

  const totalProgress = computed(() => {
    const all = tasks.value.filter(t => ['uploading', 'success', 'pending'].includes(t.status))
    if (all.length === 0) return 0
    const sum = all.reduce((acc, t) => acc + t.progress, 0)
    return Math.round(sum / all.length)
  })

  // ── 核心方法 ──

  function requestUploadWithPreview(files: File[], options?: AddTaskOptions) {
    if (showPreviewDialog.value) {
      pendingAppendFiles.value = [...files]
      return
    }
    pendingPreviewFiles.value = files
    pendingPreviewOptions.value = options ?? {}
    showPreviewDialog.value = true
  }

  function confirmPreview(confirmedFiles: File[]) {
    addTasks(confirmedFiles, pendingPreviewOptions.value)
    pendingPreviewFiles.value = []
    pendingPreviewOptions.value = {}
  }

  function cancelPreview() {
    pendingPreviewFiles.value = []
    pendingPreviewOptions.value = {}
    pendingAppendFiles.value = []
  }

  function addTasks(files: File[], options?: AddTaskOptions): UploadTask[] {
    const newTasks: UploadTask[] = files.map((file) => {
      const tooLarge = file.size > MAX_FILE_SIZE
      return {
        id: genTaskId(),
        file,
        status: (tooLarge ? 'error' : 'pending') as UploadStatus,
        progress: 0,
        speed: 0,
        result: null,
        error: tooLarge ? `文件大小超出限制（最大 50 MB，当前 ${formatBytes(file.size)}）` : '',
        abortController: null,
        onComplete: options?.onComplete,
      }
    })
    tasks.value.push(...newTasks)
    showPanel.value = true
    isMinimized.value = false
    processQueue()
    return newTasks
  }

  function processQueue() {
    const running = activeTasks.value.length
    const available = MAX_CONCURRENT - running
    if (available <= 0) return

    const pending = tasks.value.filter(t => t.status === 'pending')
    const toStart = pending.slice(0, available)
    toStart.forEach(task => executeUpload(task))
  }

  async function executeUpload(task: UploadTask) {
    const t = tasks.value.find(t => t.id === task.id)
    if (!t) return

    t.status = 'uploading'
    t.progress = 0
    t.speed = 0

    const abortController = new AbortController()
    t.abortController = abortController

    let lastLoaded = 0
    let lastTime = Date.now()

    const formData = new FormData()
    formData.append('file', t.file)

    const progressOptions: UploadProgressOptions = {
      signal: abortController.signal,
      onProgress: (percent, loaded) => {
        const now = Date.now()
        const elapsed = (now - lastTime) / 1000
        if (elapsed > 0.3) {
          t.speed = (loaded - lastLoaded) / elapsed
          lastLoaded = loaded
          lastTime = now
        }
        t.progress = percent
      },
    }

    const res = await http.uploadProgress<MediaItem>(
      '/tenant/media/upload',
      formData,
      progressOptions,
    )

    t.abortController = null

    if (res.code === -1) {
      t.status = 'cancelled'
      t.progress = 0
      t.speed = 0
    } else if (res.code === 0 && res.data) {
      res.data.file_url = normalizeUploadedUrl(res.data.file_url)
      t.status = 'success'
      t.progress = 100
      t.speed = 0
      t.result = res.data
      t.onComplete?.(res.data)
    } else {
      t.status = 'error'
      t.error = res.message || '上传失败'
      t.speed = 0
    }

    processQueue()
  }

  function cancelTask(taskId: string) {
    const t = tasks.value.find(t => t.id === taskId)
    if (!t) return
    if (t.status === 'uploading' && t.abortController) {
      t.abortController.abort()
    } else if (t.status === 'pending') {
      t.status = 'cancelled'
    }
  }

  function retryTask(taskId: string) {
    const t = tasks.value.find(t => t.id === taskId)
    if (!t || !['error', 'cancelled'].includes(t.status)) return
    t.status = 'pending'
    t.progress = 0
    t.speed = 0
    t.error = ''
    t.result = null
    processQueue()
  }

  function removeTask(taskId: string) {
    const idx = tasks.value.findIndex(t => t.id === taskId)
    if (idx === -1) return
    const t = tasks.value[idx]!
    if (t.status === 'uploading' && t.abortController) {
      t.abortController.abort()
    }
    tasks.value.splice(idx, 1)
  }

  function removeCompletedTasks() {
    tasks.value = tasks.value.filter(
      t => !['success', 'error', 'cancelled'].includes(t.status),
    )
  }

  function clearAllTasks() {
    tasks.value.forEach(t => {
      if (t.status === 'uploading' && t.abortController) {
        t.abortController.abort()
      }
    })
    tasks.value = []
  }

  function toggleMinimize() {
    isMinimized.value = !isMinimized.value
  }

  function closePanel() {
    if (hasActiveTasks.value) {
      isMinimized.value = true
    } else {
      showPanel.value = false
    }
  }

  function openPanel() {
    showPanel.value = true
    isMinimized.value = false
  }

  return {
    tasks,
    showPanel,
    isMinimized,
    pendingPreviewFiles,
    pendingPreviewOptions,
    pendingAppendFiles,
    showPreviewDialog,
    activeTasks,
    pendingTasks,
    completedTasks,
    failedTasks,
    hasActiveTasks,
    totalProgress,
    requestUploadWithPreview,
    confirmPreview,
    cancelPreview,
    addTasks,
    cancelTask,
    retryTask,
    removeTask,
    removeCompletedTasks,
    clearAllTasks,
    toggleMinimize,
    closePanel,
    openPanel,
  }
})
