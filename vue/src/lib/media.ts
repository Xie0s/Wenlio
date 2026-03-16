/**
 * lib/media.ts - 媒体文件业务逻辑层
 *
 * 职责：封装媒体文件列表查询与删除操作的状态管理与 API 调用
 * 对外暴露：useMedia() composable 函数
 */
import { ref } from 'vue'
import { http } from '@/utils/http'
import { toast } from 'vue-sonner'

const MAX_FILE_SIZE = 50 * 1024 * 1024
const IMAGE_MAX_SIZE = 5 * 1024 * 1024
const IMAGE_ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']

/**
 * 规范化后端返回的文件 URL
 * 处理 Windows 风格路径、补齐根路径前缀等边界情况
 */
export function normalizeUploadedUrl(raw: string): string {
  const value = raw.trim()
  if (!value) return value
  if (/^https?:\/\//i.test(value) || value.startsWith('//')) return value
  const webPath = value.replace(/\\/g, '/')
  return webPath.startsWith('/') ? webPath : `/${webPath}`
}

export interface MediaItem {
  id: string
  file_name: string
  file_url: string
  file_size: number
  mime_type: string
  storage_type: 'local' | 'cloud'
  created_at: string
}

export interface AuditResult {
  orphan_keys: string[]
  missing_keys: string[]
  match_count: number
}

export interface MediaUsageRef {
  theme_id: string
  theme_name: string
  theme_slug: string
  theme_deleted: boolean
  version_id: string
  version_name: string
  page_id: string
  page_title: string
  page_slug: string
}

export type MediaUsageMap = Record<string, MediaUsageRef[]>

export function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

export function useMedia() {
  const loading = ref(false)
  const deleting = ref<string | null>(null)
  const items = ref<MediaItem[]>([])

  async function loadMedia() {
    loading.value = true
    const res = await http.get<MediaItem[]>('/tenant/media')
    loading.value = false
    if (res.code === 0) {
      items.value = res.data ?? []
    } else {
      toast.error(res.message || '加载失败')
    }
  }

  async function deleteMedia(id: string, force = false): Promise<boolean> {
    deleting.value = id
    const url = force ? `/tenant/media/${id}?force=true` : `/tenant/media/${id}`
    const res = await http.delete(url)
    deleting.value = null
    if (res.code === 0) {
      items.value = items.value.filter(m => m.id !== id)
      toast.success('已删除')
      return true
    }
    toast.error(res.message || '删除失败')
    return false
  }

  const uploading = ref(false)

  async function uploadFile(file: File): Promise<MediaItem | null> {
    if (file.size > MAX_FILE_SIZE) {
      toast.error('文件大小不能超过 50MB')
      return null
    }
    uploading.value = true
    const formData = new FormData()
    formData.append('file', file)
    const res = await http.upload<MediaItem>('/tenant/media/upload', formData)
    uploading.value = false
    if (res.code === 0 && res.data) {
      res.data.file_url = normalizeUploadedUrl(res.data.file_url)
      items.value.unshift(res.data)
      return res.data
    }
    toast.error(res.message || '上传失败')
    return null
  }

  // ── 存储审计 ──
  const auditing = ref(false)
  const auditResult = ref<AuditResult | null>(null)

  async function auditStorage(storageType: 'cloud' | 'local' = 'cloud') {
    auditing.value = true
    auditResult.value = null
    const res = await http.get<AuditResult>(`/tenant/media/audit?storage_type=${storageType}`)
    auditing.value = false
    if (res.code === 0 && res.data) {
      auditResult.value = {
        ...res.data,
        orphan_keys: res.data.orphan_keys ?? [],
        missing_keys: res.data.missing_keys ?? [],
      }
    } else {
      toast.error(res.message || '审计失败')
    }
  }

  // ── 使用来源 ──
  const usageMap = ref<MediaUsageMap>({})
  const loadingUsage = ref(false)
  const usageLoaded = ref(false)

  async function loadUsage() {
    loadingUsage.value = true
    const res = await http.get<{ usages: MediaUsageMap }>('/tenant/media/usage')
    loadingUsage.value = false
    if (res.code === 0 && res.data) {
      usageMap.value = res.data.usages ?? {}
      usageLoaded.value = true
    } else {
      toast.error(res.message || '加载使用来源失败')
    }
  }

  // ── 批量清理未使用文件 ──
  const cleaningUnused = ref(false)

  async function cleanupUnused(): Promise<number> {
    cleaningUnused.value = true
    const res = await http.post<{ deleted_count: number }>('/tenant/media/cleanup-unused', {})
    cleaningUnused.value = false
    if (res.code === 0 && res.data) {
      const count = res.data.deleted_count
      if (count > 0) {
        toast.success(`已清理 ${count} 个未使用文件`)
        await loadMedia()
        usageLoaded.value = false
        usageMap.value = {}
      } else {
        toast.success('没有需要清理的文件')
      }
      return count
    }
    toast.error(res.message || '批量清理失败')
    return 0
  }

  const deletingOrphan = ref<string | null>(null)

  async function deleteOrphan(storageType: string, key: string): Promise<boolean> {
    deletingOrphan.value = key
    const res = await http.post('/tenant/media/orphan/delete', { storage_type: storageType, key })
    deletingOrphan.value = null
    if (res.code === 0) {
      if (auditResult.value) {
        auditResult.value.orphan_keys = auditResult.value.orphan_keys.filter(k => k !== key)
      }
      toast.success('孤立文件已删除')
      return true
    }
    toast.error(res.message || '删除失败')
    return false
  }

  async function uploadImage(file: File): Promise<MediaItem | null> {
    if (!IMAGE_ALLOWED_TYPES.includes(file.type)) { toast.error('不支持的图片格式'); return null }
    if (file.size > IMAGE_MAX_SIZE) { toast.error('图片大小不能超过 5MB'); return null }
    return uploadFile(file)
  }

  return {
    loading, uploading, deleting, items,
    loadMedia, uploadFile, uploadImage, deleteMedia, formatBytes,
    auditing, auditResult, auditStorage,
    deletingOrphan, deleteOrphan,
    usageMap, loadingUsage, usageLoaded, loadUsage,
    cleaningUnused, cleanupUnused,
  }
}
