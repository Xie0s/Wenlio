/**
 * utils/http.ts - HTTP 客户端封装
 * 职责：基于 fetch 封装统一的 API 请求工具，自动处理 Token 注入和错误响应。
 * 对外暴露：http 对象（get/post/put/patch/delete 方法），Token 管理函数
 */

const BASE_URL = '/api/v1'

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  request_id: string
}

export interface RequestConfig {
  suppressErrorToast?: boolean
  throwHttpError?: boolean
}

export interface PageData<T = any> {
  list: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Token 管理
function getAccessToken(): string | null {
  return localStorage.getItem('access_token')
}

export function hasToken(): boolean {
  return !!getAccessToken()
}

export function setToken(token: string) {
  localStorage.setItem('access_token', token)
}

export function clearToken() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('auth')
}

// 核心请求方法
async function request<T = any>(
  url: string,
  options: RequestInit = {},
  config: RequestConfig = {},
): Promise<ApiResponse<T>> {
  const token = getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${url}`, { ...options, headers, cache: 'no-store' })

  // Token 过期或无效
  if (res.status === 401) {
    // 登录接口本身返回 401 属于正常业务逻辑（密码错误等），直接返回响应体
    if (url === '/auth/login' && !config.throwHttpError) {
      return res.json()
    }
    clearToken()
    // 只在 admin 路由才重定向到登录页，避免 admin 组件的遗留请求在用户已离开 admin
    // 后触发，导致用户从首页等公开页面被错误踢出
    const { default: router } = await import('@/router')
    if (router.currentRoute.value.path.startsWith('/admin')) {
      router.replace('/admin/login')
    }
    return { code: 401001, message: '请先登录', data: null as any, request_id: '' }
  }

  let body: any
  try {
    body = await res.json()
  } catch {
    // 代理错误（502/504 等）返回 HTML，无法解析为 JSON
    if (config.throwHttpError) {
      const error = new Error(`网络错误 (${res.status})`) as Error & { status?: number; code?: number }
      error.status = res.status
      throw error
    }
    return { code: res.status, message: `网络错误 (${res.status})`, data: null as any, request_id: '' }
  }
  if (!res.ok && config.throwHttpError) {
    const error = new Error(body?.message || '请求失败') as Error & { status?: number; code?: number }
    error.status = res.status
    error.code = body?.code
    throw error
  }
  return body
}

// 文件上传请求（不设置 Content-Type，让浏览器自动处理 multipart boundary）
async function uploadRequest<T = any>(
  url: string,
  formData: FormData,
): Promise<ApiResponse<T>> {
  const token = getAccessToken()
  const headers: Record<string, string> = {}
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${url}`, {
    method: 'POST',
    headers,
    body: formData,
  })

  if (res.status === 401) {
    clearToken()
    const { default: router } = await import('@/router')
    if (router.currentRoute.value.path.startsWith('/admin')) {
      router.replace('/admin/login')
    }
    return { code: 401001, message: '请先登录', data: null as any, request_id: '' }
  }

  try {
    return await res.json()
  } catch {
    return { code: res.status, message: `网络错误 (${res.status})`, data: null as any, request_id: '' }
  }
}

// 带进度回调的文件上传（基于 XMLHttpRequest）
export interface UploadProgressOptions {
  onProgress?: (percent: number, loaded: number, total: number) => void
  signal?: AbortSignal
}

function uploadWithProgress<T = any>(
  url: string,
  formData: FormData,
  options?: UploadProgressOptions,
): Promise<ApiResponse<T>> {
  return new Promise((resolve) => {
    const token = getAccessToken()
    const xhr = new XMLHttpRequest()

    if (options?.signal) {
      options.signal.addEventListener('abort', () => xhr.abort())
    }

    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable && options?.onProgress) {
        const percent = Math.round((e.loaded / e.total) * 100)
        options.onProgress(percent, e.loaded, e.total)
      }
    })

    xhr.addEventListener('load', () => {
      if (xhr.status === 401) {
        clearToken()
        import('@/router').then(({ default: router }) => {
          if (router.currentRoute.value.path.startsWith('/admin')) {
            router.replace('/admin/login')
          }
        })
        resolve({ code: 401001, message: '请先登录', data: null as any, request_id: '' })
        return
      }
      try {
        resolve(JSON.parse(xhr.responseText))
      } catch {
        resolve({ code: xhr.status, message: `网络错误 (${xhr.status})`, data: null as any, request_id: '' })
      }
    })

    xhr.addEventListener('error', () => {
      resolve({ code: 0, message: '网络错误', data: null as any, request_id: '' })
    })

    xhr.addEventListener('abort', () => {
      resolve({ code: -1, message: '上传已取消', data: null as any, request_id: '' })
    })

    xhr.open('POST', `${BASE_URL}${url}`)
    if (token) xhr.setRequestHeader('Authorization', `Bearer ${token}`)
    xhr.send(formData)
  })
}

// 导出便捷方法
export const http = {
  get<T = any>(url: string, params?: Record<string, any>) {
    const query = params
      ? '?' + new URLSearchParams(
          Object.entries(params)
            .filter(([, v]) => v !== undefined && v !== null && v !== '')
            .reduce((acc, [k, v]) => {
              if (Array.isArray(v)) {
                v.forEach(item => acc.append(k, String(item)))
              } else {
                acc.append(k, String(v))
              }
              return acc
            }, new URLSearchParams()),
        ).toString()
      : ''
    return request<T>(url + query)
  },

  post<T = any>(url: string, data?: any, config?: RequestConfig) {
    return request<T>(url, { method: 'POST', body: data ? JSON.stringify(data) : undefined }, config)
  },

  put<T = any>(url: string, data?: any, config?: RequestConfig) {
    return request<T>(url, { method: 'PUT', body: data ? JSON.stringify(data) : undefined }, config)
  },

  patch<T = any>(url: string, data?: any, config?: RequestConfig) {
    return request<T>(url, { method: 'PATCH', body: data ? JSON.stringify(data) : undefined }, config)
  },

  delete<T = any>(url: string, config?: RequestConfig) {
    return request<T>(url, { method: 'DELETE' }, config)
  },

  upload<T = any>(url: string, formData: FormData) {
    return uploadRequest<T>(url, formData)
  },

  uploadProgress<T = any>(url: string, formData: FormData, options?: UploadProgressOptions) {
    return uploadWithProgress<T>(url, formData, options)
  },
}
