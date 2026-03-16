/*
 * theme.ts - 主题状态仓库（Pinia）
 * 职责：
 * 1) 统一管理主题模式（light/dark/system）与派生主题
 * 2) 处理系统主题监听、DOM 同步与本地持久化兼容
 * 3) 暴露 store 与兼容函数 API，供应用初始化与组件消费
 * 对外暴露：
 * - useThemeStore()
 * - initTheme()/setTheme()/toggleTheme()/getThemeMode()/getResolvedTheme()/subscribeThemeChange()
 */

import { defineStore } from 'pinia'
import type { Pinia } from 'pinia'
import 'pinia-plugin-persistedstate'

export type ThemeMode = 'light' | 'dark' | 'system'
export type ResolvedTheme = 'light' | 'dark'
type ThemeSnapshot = { mode: ThemeMode; resolvedTheme: ResolvedTheme }
type ThemeListener = (snapshot: ThemeSnapshot) => void

type ThemeState = {
  mode: ThemeMode
  systemTheme: ResolvedTheme
  initialized: boolean
}

const THEME_LEGACY_STORAGE_KEY = 'theme'
const THEME_STORE_STORAGE_KEY = 'theme-store'
let detachSystemListener: (() => void) | null = null

function canUseDOM() {
  return typeof window !== 'undefined' && typeof document !== 'undefined'
}

function isThemeMode(value: string): value is ThemeMode {
  return value === 'light' || value === 'dark' || value === 'system'
}

function getLegacyStoredThemeMode(): ThemeMode | null {
  if (!canUseDOM()) {
    return null
  }

  try {
    const raw = window.localStorage.getItem(THEME_LEGACY_STORAGE_KEY)
    if (raw && isThemeMode(raw)) {
      return raw
    }
  } catch {
    // 忽略 localStorage 不可用场景
  }

  return null
}

function getSystemTheme(): ResolvedTheme {
  if (!canUseDOM()) {
    return 'light'
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export const useThemeStore = defineStore('theme', {
  state: (): ThemeState => ({
    mode: 'system',
    systemTheme: 'light',
    initialized: false,
  }),

  getters: {
    resolvedTheme: (state): ResolvedTheme => (state.mode === 'system' ? state.systemTheme : state.mode),
  },

  actions: {
    applyResolvedThemeToDOM() {
      if (!canUseDOM()) {
        return
      }

      const root = document.documentElement
      const resolvedTheme = this.resolvedTheme
      root.classList.toggle('dark', resolvedTheme === 'dark')
      root.setAttribute('data-theme', resolvedTheme)
    },

    syncLegacyThemeStorage() {
      if (!canUseDOM()) {
        return
      }

      try {
        window.localStorage.setItem(THEME_LEGACY_STORAGE_KEY, this.mode)
      } catch {
        // 忽略 localStorage 不可用场景
      }
    },

    syncSystemTheme() {
      this.systemTheme = getSystemTheme()
    },

    setupSystemThemeListener() {
      if (detachSystemListener) {
        detachSystemListener()
        detachSystemListener = null
      }

      if (!canUseDOM() || this.mode !== 'system') {
        return
      }

      const query = window.matchMedia('(prefers-color-scheme: dark)')
      const handleChange = () => {
        if (this.mode !== 'system') {
          return
        }
        this.syncSystemTheme()
        this.applyResolvedThemeToDOM()
      }

      if (typeof query.addEventListener === 'function') {
        query.addEventListener('change', handleChange)
        detachSystemListener = () => query.removeEventListener('change', handleChange)
      } else {
        query.addListener(handleChange)
        detachSystemListener = () => query.removeListener(handleChange)
      }
    },

    init(): ResolvedTheme {
      if (this.initialized) {
        this.applyResolvedThemeToDOM()
        return this.resolvedTheme
      }

      const legacyMode = getLegacyStoredThemeMode()
      if (legacyMode) {
        this.mode = legacyMode
      }

      this.syncSystemTheme()
      this.applyResolvedThemeToDOM()
      this.setupSystemThemeListener()
      this.syncLegacyThemeStorage()
      this.initialized = true

      return this.resolvedTheme
    },

    setTheme(mode: ThemeMode): ResolvedTheme {
      if (!this.initialized) {
        this.init()
      }

      this.mode = mode
      this.syncSystemTheme()
      this.applyResolvedThemeToDOM()
      this.setupSystemThemeListener()
      this.syncLegacyThemeStorage()

      return this.resolvedTheme
    },

    toggleTheme(): ResolvedTheme {
      if (!this.initialized) {
        this.init()
      }

      const nextMode: ThemeMode = this.resolvedTheme === 'dark' ? 'light' : 'dark'
      return this.setTheme(nextMode)
    },
  },

  persist: {
    key: THEME_STORE_STORAGE_KEY,
    pick: ['mode'],
  },
})

function getThemeStore(pinia?: Pinia) {
  return useThemeStore(pinia)
}

function ensureThemeStoreInitialized(pinia?: Pinia) {
  const store = getThemeStore(pinia)
  if (!store.initialized) {
    store.init()
  }
  return store
}

export function getThemeMode(pinia?: Pinia): ThemeMode {
  return ensureThemeStoreInitialized(pinia).mode
}

export function getResolvedTheme(pinia?: Pinia): ResolvedTheme {
  return ensureThemeStoreInitialized(pinia).resolvedTheme
}

export function initTheme(pinia?: Pinia): ResolvedTheme {
  return getThemeStore(pinia).init()
}

export function setTheme(mode: ThemeMode, pinia?: Pinia): ResolvedTheme {
  return ensureThemeStoreInitialized(pinia).setTheme(mode)
}

export function toggleTheme(pinia?: Pinia): ResolvedTheme {
  return ensureThemeStoreInitialized(pinia).toggleTheme()
}

export function subscribeThemeChange(listener: ThemeListener, pinia?: Pinia): () => void {
  const store = ensureThemeStoreInitialized(pinia)

  listener({
    mode: store.mode,
    resolvedTheme: store.resolvedTheme,
  })

  const unsubscribe = store.$subscribe(
    () => {
      listener({
        mode: store.mode,
        resolvedTheme: store.resolvedTheme,
      })
    },
    { detached: true },
  )

  return unsubscribe
}
