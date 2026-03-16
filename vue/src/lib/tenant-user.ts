/**
 * tenant-user.ts - 租户用户管理业务逻辑层
 *
 * 职责：封装租户用户列表的状态管理与 API 操作（加载、创建、禁用/启用、重置密码）
 * 对外暴露：useTenantUserList() composable 函数
 */
import { ref } from 'vue'
import { http } from '@/utils/http'
import type { PageData } from '@/utils/http'
import { toast } from 'vue-sonner'

export interface TenantUser {
  id: string
  tenant_id: string
  username: string
  name: string
  email: string
  role: string
  status: string
  avatar_url: string
  bio: string
  profile_bg_url: string
  last_login_at: string
  created_at: string
  updated_at: string
}

export const ROLE_LABEL: Record<string, string> = {
  super_admin: '超级管理员',
  tenant_admin: '租户管理员',
}

export function formatDate(dateStr: string) {
  if (!dateStr || dateStr.startsWith('0001')) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

export function formatDateTime(dateStr: string) {
  if (!dateStr || dateStr.startsWith('0001')) return '从未'
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

export function useTenantUserList() {
  const users = ref<TenantUser[]>([])
  const loading = ref(false)

  // 创建弹窗
  const showCreate = ref(false)
  const form = ref({ username: '', password: '', name: '', email: '' })

  // 编辑弹窗
  const showEdit = ref(false)
  const editUserId = ref('')
  const editForm = ref({ name: '', email: '' })

  // 重置密码弹窗
  const showResetPwd = ref(false)
  const resetUserId = ref('')
  const resetUserName = ref('')
  const newPassword = ref('')

  async function loadUsers() {
    loading.value = true
    const res = await http.get<PageData<TenantUser>>('/tenant/users', { page: 1, page_size: 100 })
    loading.value = false
    if (res.code === 0 && res.data) { users.value = res.data.list || [] }
  }

  async function createUser() {
    const res = await http.post('/tenant/users', form.value)
    if (res.code === 0) {
      toast.success('用户创建成功')
      showCreate.value = false
      form.value = { username: '', password: '', name: '', email: '' }
      loadUsers()
    } else { toast.error(res.message) }
  }

  function openEditUser(user: TenantUser) {
    editUserId.value = user.id
    editForm.value = {
      name: user.name || '',
      email: user.email || '',
    }
    showEdit.value = true
  }

  async function updateUser() {
    const res = await http.patch(`/tenant/users/${editUserId.value}`, {
      name: editForm.value.name,
      email: editForm.value.email,
    })
    if (res.code === 0) {
      toast.success('用户信息已更新')
      showEdit.value = false
      loadUsers()
    } else {
      toast.error(res.message)
    }
  }

  async function deleteUser(user: TenantUser) {
    const res = await http.delete(`/tenant/users/${user.id}`)
    if (res.code === 0) {
      toast.success('用户已删除')
      loadUsers()
    } else {
      toast.error(res.message)
    }
  }

  async function toggleUserStatus(user: TenantUser) {
    const action = user.status === 'active' ? 'deactivate' : 'activate'
    const res = await http.post(`/tenant/users/${user.id}/${action}`)
    if (res.code === 0) {
      toast.success(action === 'deactivate' ? '已禁用' : '已启用')
      loadUsers()
    } else {
      toast.error(res.message)
    }
  }

  function openResetPassword(user: TenantUser) {
    resetUserId.value = user.id
    resetUserName.value = user.name
    newPassword.value = ''
    showResetPwd.value = true
  }

  async function resetPassword() {
    if (!newPassword.value) {
      toast.error('请输入新密码')
      return
    }
    const res = await http.post(`/tenant/users/${resetUserId.value}/reset-password`, {
      new_password: newPassword.value,
    })
    if (res.code === 0) {
      toast.success('密码已重置')
      showResetPwd.value = false
    } else {
      toast.error(res.message)
    }
  }

  return {
    // 状态
    users,
    loading,
    // 创建弹窗
    showCreate,
    form,
    // 编辑弹窗
    showEdit,
    editForm,
    // 重置密码弹窗
    showResetPwd,
    resetUserName,
    newPassword,
    // 操作
    loadUsers,
    createUser,
    openEditUser,
    updateUser,
    deleteUser,
    toggleUserStatus,
    openResetPassword,
    resetPassword,
  }
}
