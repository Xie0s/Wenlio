/**
 * composables/useAutoSave.ts - 通用自动保存
 * 职责：定时检查脏状态并触发保存回调，自动随作用域销毁
 * 对外暴露：useAutoSave(options) → { start, stop }
 */

import { onScopeDispose } from 'vue'

export interface AutoSaveOptions {
  save: () => Promise<void>
  isDirty: () => boolean
  interval?: number
}

export function useAutoSave(options: AutoSaveOptions) {
  const { save, isDirty, interval = 30_000 } = options
  let timer: ReturnType<typeof setInterval> | null = null

  function start() {
    stop()
    timer = setInterval(async () => {
      if (isDirty()) await save()
    }, interval)
  }

  function stop() {
    if (timer) { clearInterval(timer); timer = null }
  }

  onScopeDispose(stop)

  return { start, stop }
}
