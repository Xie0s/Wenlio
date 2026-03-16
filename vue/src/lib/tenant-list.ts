/**
 * tenant-list.ts - 租户管理业务逻辑层（超管专属）
 *
 * 职责：封装租户列表的状态管理与 API 操作（加载、搜索、创建、封禁/解封）
 * 对外暴露：useTenantList() composable 函数
 */
import { ref } from 'vue'
import { http } from '@/utils/http'
import type { PageData } from '@/utils/http'
import type { Tenant, CreateTenantReq, UpdateTenantReq } from '@/utils/types'
import { toast } from 'vue-sonner'

const defaultForm = (): CreateTenantReq => ({
  id: '', name: '', logo_url: '',
  admin_username: '', admin_password: '', admin_name: '',
})

export function useTenantList() {
  const tenants = ref<Tenant[]>([])
  const total = ref(0)
  const page = ref(1)
  const keyword = ref('')
  const loading = ref(false)

  // 创建弹窗
  const showCreate = ref(false)
  const form = ref<CreateTenantReq>(defaultForm())

  const showEdit = ref(false)
  const editTenantId = ref('')
  const editForm = ref<UpdateTenantReq>({ name: '', logo_url: '' })

  async function loadTenants() {
    loading.value = true
    const res = await http.get<PageData<Tenant>>('/admin/tenants', { page: page.value, page_size: 20, keyword: keyword.value })
    loading.value = false
    if (res.code === 0 && res.data) {
      tenants.value = res.data.list || []
      total.value = res.data.pagination.total
    }
  }

  async function createTenant() {
    const res = await http.post<Tenant>('/admin/tenants', form.value)
    if (res.code === 0) {
      toast.success('租户创建成功')
      showCreate.value = false
      form.value = defaultForm()
      loadTenants()
    } else {
      toast.error(res.message)
    }
  }

  function openEditTenant(tenant: Tenant) {
    if (tenant.status === 'deleting') {
      toast.error('租户正在删除中，暂不支持编辑')
      return
    }
    editTenantId.value = tenant.id
    editForm.value = {
      name: tenant.name || '',
      logo_url: tenant.logo_url || '',
    }
    showEdit.value = true
  }

  async function updateTenant() {
    const res = await http.patch(`/admin/tenants/${editTenantId.value}`, editForm.value)
    if (res.code === 0) {
      toast.success('租户信息已更新')
      showEdit.value = false
      loadTenants()
    } else {
      toast.error(res.message)
    }
  }

  async function toggleStatus(tenant: Tenant) {
    if (tenant.status === 'deleting') {
      toast.error('租户正在删除中，无法变更状态')
      return
    }
    const action = tenant.status === 'active' ? 'suspend' : 'activate'
    const res = await http.post(`/admin/tenants/${tenant.id}/${action}`)
    if (res.code === 0) {
      toast.success(action === 'suspend' ? '已封禁' : '已解封')
      loadTenants()
    } else {
      toast.error(res.message)
    }
  }

  async function deleteTenant(tenant: Tenant) {
    if (tenant.status === 'deleting') {
      toast.info('租户删除任务已在执行中，请稍后刷新列表')
      return
    }
    const res = await http.delete(`/admin/tenants/${tenant.id}`)
    if (res.code === 0) {
      toast.success('已提交删除任务，后台将异步清理租户数据')
      loadTenants()
    } else {
      toast.error(res.message)
    }
  }

  return {
    // 状态
    tenants,
    total,
    page,
    keyword,
    loading,
    // 创建弹窗
    showCreate,
    form,
    showEdit,
    editForm,
    // 操作
    loadTenants,
    createTenant,
    openEditTenant,
    updateTenant,
    toggleStatus,
    deleteTenant,
  }
}
