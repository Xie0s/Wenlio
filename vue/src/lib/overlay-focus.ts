/**
 * lib/overlay-focus.ts - 覆盖层打开前的焦点守卫工具
 * 职责：
 *   1) 在打开 Dialog/ConfirmDialog 等覆盖层前，安全释放当前焦点
 *   2) 对编辑器焦点做多帧兜底清理，避免 aria-hidden 与聚焦元素冲突
 *   3) 全局 MutationObserver 守卫：监听 aria-hidden 设置，自动 blur 被隐藏容器内的焦点元素
 * 功能边界：
 *   - 仅处理焦点释放与打开时序，不负责具体业务状态管理
 * 对外暴露：
 *   - blurActiveElement()
 *   - releaseOverlayFocus(options)
 *   - openOverlaySafely(openOverlay, options)
 *   - setupAriaHiddenGuard() → cleanup fn
 */

export interface OverlayFocusOptions {
  /** 编辑器根 DOM（如 ProseMirror 根节点） */
  editorElement?: HTMLElement | null
  /** 编辑器 blur 回调（如 tiptap 的 commands.blur） */
  blurEditor?: (() => void) | null
  /** 焦点释放后的等待帧数，默认 2 帧 */
  settleFrameCount?: number
}

function nextFrame(): Promise<void> {
  return new Promise((resolve) => {
    window.requestAnimationFrame(() => resolve())
  })
}

function isFocusInsideEditor(editorElement?: HTMLElement | null): boolean {
  if (!(editorElement instanceof HTMLElement)) return false
  const activeEl = document.activeElement
  return activeEl instanceof HTMLElement
    && (activeEl === editorElement || editorElement.contains(activeEl))
}

export function blurActiveElement(): void {
  const activeEl = document.activeElement
  if (activeEl instanceof HTMLElement) activeEl.blur()
}

export function releaseOverlayFocus(options: OverlayFocusOptions = {}): void {
  const { editorElement, blurEditor } = options

  // 先释放当前活跃焦点（覆盖编辑器 / 按钮 / 输入框）
  blurActiveElement()

  // 再显式清理编辑器 DOM 焦点
  if (editorElement instanceof HTMLElement) editorElement.blur()

  // 最后触发编辑器命令级 blur（兜底）
  blurEditor?.()
}

/**
 * 全局 aria-hidden 守卫：监听任意元素被设置 aria-hidden="true" 的时机，
 * 若此时焦点仍在被隐藏的元素内，立即 blur，从根本上消除浏览器警告。
 * 在 App.vue 的 onMounted 中调用一次，返回清理函数。
 */
export function setupAriaHiddenGuard(): () => void {
  if (typeof MutationObserver === 'undefined' || typeof document === 'undefined') {
    return () => {}
  }

  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (
        mutation.type === 'attributes' &&
        mutation.attributeName === 'aria-hidden' &&
        (mutation.target as Element).getAttribute('aria-hidden') === 'true'
      ) {
        const target = mutation.target as HTMLElement
        const active = document.activeElement
        if (active instanceof HTMLElement && active !== target && target.contains(active)) {
          active.blur()
        }
      }
    }
  })

  observer.observe(document.body, {
    subtree: true,
    attributes: true,
    attributeFilter: ['aria-hidden'],
  })

  return () => observer.disconnect()
}

export async function openOverlaySafely(
  openOverlay: () => void,
  options: OverlayFocusOptions = {},
): Promise<void> {
  if (typeof window === 'undefined') {
    openOverlay()
    return
  }

  const frameCount = Math.max(1, options.settleFrameCount ?? 2)

  releaseOverlayFocus(options)

  // 多帧兜底：如果焦点仍在编辑器内，继续释放，避免 aria-hidden 触发浏览器警告
  for (let i = 0; i < frameCount; i += 1) {
    if (!isFocusInsideEditor(options.editorElement)) break
    await nextFrame()
    releaseOverlayFocus(options)
  }

  // 至少等待一帧，让 blur 状态稳定后再打开覆盖层
  await nextFrame()
  openOverlay()
  // openOverlay 之后立即再次 blur：Vue 的 DOM 更新（Reka UI 设置 aria-hidden）
  // 是通过 microtask 异步执行的，此处同步 blur 确保更新前编辑器已失焦
  releaseOverlayFocus(options)
}
