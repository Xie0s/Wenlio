/**
 * stores/auth.ts - 认证状态管理
 * 职责：管理用户登录状态、Token、当前用户信息，提供登录/登出操作。
 * 对外暴露：useAuthStore
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { http, hasToken, setToken, clearToken } from '@/utils/http'
import type { UserInfo, LoginReq, LoginResp, RegisterTenantReq } from '@/utils/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null)
  const isLoggedIn = computed(() => !!user.value)
  const isSuperAdmin = computed(() => user.value?.role === 'super_admin')
  const isTenantAdmin = computed(() => user.value?.role === 'tenant_admin')

  // 登录
  async function login(req: LoginReq): Promise<{ success: boolean; message?: string }> {
    const res = await http.post<LoginResp>('/auth/login', req)
    if (res.code === 0 && res.data) {
      setToken(res.data.access_token)
      user.value = res.data.user
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  async function register(req: RegisterTenantReq): Promise<{ success: boolean; message?: string }> {
    const res = await http.post<LoginResp>('/auth/register', req)
    if (res.code === 0 && res.data) {
      setToken(res.data.access_token)
      user.value = res.data.user
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  // 登出
  async function logout() {
    try {
      await http.post('/auth/logout')
    } catch {
      // 忽略登出接口错误
    }
    clearToken()
    user.value = null
  }

  function clearAuthState() {
    clearToken()
    user.value = null
  }

  // 获取当前用户信息（页面刷新后恢复状态）
  async function fetchCurrentUser(): Promise<boolean> {
    if (!hasToken()) {
      clearAuthState()
      return false
    }
    const res = await http.get<UserInfo>('/auth/me')
    if (res.code === 0 && res.data) {
      user.value = res.data
      return true
    }
    clearAuthState()
    return false
  }

  // 修改密码
  async function changePassword(oldPassword: string, newPassword: string): Promise<{ success: boolean; message?: string }> {
    const res = await http.patch('/auth/me/password', {
      old_password: oldPassword,
      new_password: newPassword,
    })
    if (res.code === 0) {
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  // 更新个人资料（name / bio 可选）
  async function updateProfile(fields: { name?: string; bio?: string }): Promise<{ success: boolean; message?: string }> {
    const res = await http.patch<UserInfo>('/auth/me/profile', fields)
    if (res.code === 0 && res.data) {
      user.value = res.data
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  // 上传头像
  async function uploadAvatar(file: File): Promise<{ success: boolean; message?: string }> {
    const form = new FormData()
    form.append('file', file)
    const res = await http.upload<UserInfo>('/auth/me/avatar', form)
    if (res.code === 0 && res.data) {
      user.value = res.data
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  // 上传个人主页背景
  async function uploadProfileBg(file: File): Promise<{ success: boolean; message?: string }> {
    const form = new FormData()
    form.append('file', file)
    const res = await http.upload<UserInfo>('/auth/me/profile-bg', form)
    if (res.code === 0 && res.data) {
      user.value = res.data
      return { success: true }
    }
    return { success: false, message: res.message }
  }

  return {
    user,
    isLoggedIn,
    isSuperAdmin,
    isTenantAdmin,
    login,
    register,
    logout,
    clearAuthState,
    fetchCurrentUser,
    changePassword,
    updateProfile,
    uploadAvatar,
    uploadProfileBg,
  }
}, {
  persist: {
    pick: ['user'],
  },
})
