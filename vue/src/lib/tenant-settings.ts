/**
 * tenant-settings.ts - 租户设置业务逻辑层
 *
 * 职责：封装租户设置（存储/AI）的状态管理与 API 操作
 * 对外暴露：useTenantSettings() composable 函数
 */
import { ref, reactive } from 'vue'
import { http } from '@/utils/http'
import type {
  TenantSettings,
  UpdateStorageReq,
  StorageUsage,
  AISettings,
  AccessSettings,
} from '@/utils/types'
import { toast } from 'vue-sonner'

const defaultStorageForm = (): UpdateStorageReq => ({
  enabled: false,
  provider: 'custom',
  endpoint: '',
  region: '',
  bucket: '',
  access_key_id: '',
  secret_access_key: '',
  custom_domain: '',
  default_target: 'local',
  force_path_style: false,
})

export function useTenantSettings() {
  const settings = ref<TenantSettings | null>(null)
  const loading = ref(false)
  const saving = ref(false)
  const testing = ref(false)

  // 存储设置表单
  const storageForm = reactive<UpdateStorageReq>(defaultStorageForm())
  const usage = ref<StorageUsage>({ used_bytes: 0, limit_bytes: 104857600, used_mb: 0, limit_mb: 100, percent: 0 })

  // AI 设置表单
  const aiForm = reactive<AISettings>({
    enabled: false,
    chat: { base_url: '', model_id: '' },
    embedding: { base_url: '', model_id: '' },
    reranker: { base_url: '', model_id: '' },
  })

  // 访问控制设置表单
  const accessForm = reactive<AccessSettings>({
    maintenance_mode: false,
    gallery_login_required: false,
  })

  async function loadSettings() {
    loading.value = true
    const res = await http.get<TenantSettings>('/tenant/settings')
    loading.value = false
    if (res.code === 0 && res.data) {
      settings.value = res.data
      // 同步到表单（不含 secret_access_key，后端不回传）
      if (res.data.storage) {
        const s = res.data.storage
        Object.assign(storageForm, {
          enabled: s.enabled,
          provider: s.provider || 'custom',
          endpoint: s.endpoint || '',
          region: s.region || '',
          bucket: s.bucket || '',
          access_key_id: s.access_key_id || '',
          secret_access_key: '',
          custom_domain: s.custom_domain || '',
          default_target: s.default_target || 'local',
          force_path_style: s.force_path_style || false,
        })
      }
      if (res.data.ai) {
        const a = res.data.ai
        Object.assign(aiForm, {
          enabled: a.enabled,
          chat: a.chat ? { ...a.chat } : { base_url: '', model_id: '' },
          embedding: a.embedding ? { ...a.embedding } : { base_url: '', model_id: '' },
          reranker: a.reranker ? { ...a.reranker } : { base_url: '', model_id: '' },
        })
      }
      if (res.data.access) {
        Object.assign(accessForm, res.data.access)
      }
    }
  }

  async function loadUsage() {
    const res = await http.get<StorageUsage>('/tenant/settings/storage/usage')
    if (res.code === 0 && res.data) {
      usage.value = res.data
    }
  }

  async function saveStorage() {
    saving.value = true
    // 若 secret_access_key 为空则不传（避免覆盖已有密钥）
    const payload: Partial<UpdateStorageReq> = { ...storageForm }
    if (!payload.secret_access_key) delete payload.secret_access_key
    const res = await http.patch('/tenant/settings/storage', payload)
    saving.value = false
    if (res.code === 0) {
      toast.success('存储设置已保存')
      loadSettings()
    } else {
      toast.error(res.message)
    }
  }

  async function testS3Connection() {
    testing.value = true
    const payload: Partial<UpdateStorageReq> = { ...storageForm }
    if (!payload.secret_access_key) delete payload.secret_access_key
    const res = await http.post('/tenant/settings/storage/test', payload)
    testing.value = false
    if (res.code === 0) {
      toast.success('云存储连接成功')
    } else {
      toast.error(`连接失败：${res.message}`)
    }
  }

  async function saveAI() {
    saving.value = true
    const res = await http.patch('/tenant/settings/ai', aiForm)
    saving.value = false
    if (res.code === 0) {
      toast.success('AI 设置已保存')
      loadSettings()
    } else {
      toast.error(res.message)
    }
  }

  async function saveAccess() {
    saving.value = true
    const res = await http.patch('/tenant/settings/access', accessForm)
    saving.value = false
    if (res.code === 0) {
      toast.success('访问控制设置已保存')
      loadSettings()
    } else {
      toast.error(res.message)
    }
  }

  function formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
  }

  return {
    settings,
    loading,
    saving,
    testing,
    storageForm,
    aiForm,
    usage,
    loadSettings,
    loadUsage,
    saveStorage,
    testS3Connection,
    saveAI,
    accessForm,
    saveAccess,
    formatBytes,
  }
}
