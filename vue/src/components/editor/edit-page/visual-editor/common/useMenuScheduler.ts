/**
 * useMenuScheduler.ts - 浮动菜单共享调度器
 * 职责：为同一 editor 实例的多个浮动菜单提供统一的 RAF 调度和事件绑定
 *       共享 scroll/resize/editor 事件监听，每帧只计算一次 viewport rect
 * 对外暴露：useMenuScheduler, isRectInView, isTopInView
 */

import { onMounted, onBeforeUnmount } from 'vue'
import type { Editor } from '@tiptap/vue-3'

/** 查找最近的可滚动祖先元素 */
function findScrollParent(el: HTMLElement): HTMLElement | null {
  let current = el.parentElement
  while (current) {
    const { overflowY } = getComputedStyle(current)
    if (overflowY === 'auto' || overflowY === 'scroll') return current
    current = current.parentElement
  }
  return null
}

/** 判断目标矩形是否与 viewport 有重叠 */
export function isRectInView(rect: DOMRect, viewport: DOMRect): boolean {
  return rect.bottom > viewport.top && rect.top < viewport.bottom
    && rect.right > viewport.left && rect.left < viewport.right
}

/** 判断菜单 top 位置是否在 viewport 可见范围内 */
export function isTopInView(top: number, viewport: DOMRect): boolean {
  return top >= viewport.top && top <= viewport.bottom - 8
}

// ── 共享调度器：同一 editor 的所有菜单共用一个实例 ──

interface Scheduler {
  callbacks: Set<() => void>
  scheduleUpdate: () => void
  scrollParent: HTMLElement | null
  viewport: DOMRect
  mountCount: number
  rafId: number | null
  scheduled: boolean
}

const schedulers = new WeakMap<Editor, Scheduler>()

function ensureScheduler(editor: Editor): Scheduler {
  let s = schedulers.get(editor)
  if (s) return s

  s = {
    callbacks: new Set(),
    scheduleUpdate: null!,
    scrollParent: null,
    viewport: new DOMRect(0, 0, window.innerWidth, window.innerHeight),
    mountCount: 0,
    rafId: null,
    scheduled: false,
  }

  const scheduler = s
  function scheduleUpdate() {
    if (scheduler.scheduled) return
    scheduler.scheduled = true
    scheduler.rafId = requestAnimationFrame(() => {
      scheduler.scheduled = false

      // 懒发现：onMounted 时 DOM 可能未就绪，在此重试
      if (!scheduler.scrollParent) {
        const dom = editor.view?.dom
        if (dom instanceof HTMLElement) {
          scheduler.scrollParent = findScrollParent(dom)
          scheduler.scrollParent?.addEventListener('scroll', scheduleUpdate, { passive: true })
        }
      }

      scheduler.viewport = scheduler.scrollParent
        ? scheduler.scrollParent.getBoundingClientRect()
        : new DOMRect(0, 0, window.innerWidth, window.innerHeight)
      scheduler.callbacks.forEach(cb => cb())
    })
  }

  s.scheduleUpdate = scheduleUpdate
  schedulers.set(editor, s)
  return s
}

/**
 * 注册浮动菜单的 computeMenu 到共享调度器
 * 同一 editor 的多个菜单共享事件监听和 RAF 调度
 *
 * @param editor Tiptap Editor 实例
 * @param computeMenu 位置计算回调
 */
export function useMenuScheduler(editor: Editor, computeMenu: () => void) {
  const s = ensureScheduler(editor)
  s.callbacks.add(computeMenu)

  onMounted(() => {
    s.mountCount++
    if (s.mountCount === 1) {
      editor.on('selectionUpdate', s.scheduleUpdate)
      editor.on('transaction', s.scheduleUpdate)
      const dom = editor.view?.dom
      if (dom instanceof HTMLElement) {
        s.scrollParent = findScrollParent(dom)
        s.scrollParent?.addEventListener('scroll', s.scheduleUpdate, { passive: true })
      }
      window.addEventListener('resize', s.scheduleUpdate, { passive: true })
    }
    s.scheduleUpdate()
  })

  onBeforeUnmount(() => {
    s.callbacks.delete(computeMenu)
    s.mountCount--
    if (s.mountCount === 0) {
      editor.off('selectionUpdate', s.scheduleUpdate)
      editor.off('transaction', s.scheduleUpdate)
      if (s.rafId !== null) cancelAnimationFrame(s.rafId)
      s.scrollParent?.removeEventListener('scroll', s.scheduleUpdate)
      window.removeEventListener('resize', s.scheduleUpdate)
      schedulers.delete(editor)
    }
  })

  return {
    scheduleUpdate: s.scheduleUpdate,
    getViewport: () => s.viewport,
  }
}
