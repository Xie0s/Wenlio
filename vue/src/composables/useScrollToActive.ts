/**
 * composables/useScrollToActive.ts - TOC 高亮条目自动滚动
 * 职责：监听 activeId 与条目列表变化，将高亮按钮自动滚动到容器可视区域中央
 * 对外暴露：useScrollToActive(containerRef, activeId, items)
 */

import { type Ref, watch, nextTick } from 'vue'
import type { TocItem } from '@/composables/useToc'

/**
 * 将指定 data-toc-id 的按钮滚动到容器可视区中央
 */
function scrollItemIntoView(container: HTMLElement, id: string) {
  if (!id) return
  const button = container.querySelector<HTMLButtonElement>(`button[data-toc-id="${id}"]`)
  if (!button) return

  const containerHeight = container.clientHeight
  const containerRect = container.getBoundingClientRect()
  const buttonRect = button.getBoundingClientRect()
  const itemTop = buttonRect.top - containerRect.top + container.scrollTop
  const itemHeight = button.offsetHeight
  const maxScrollTop = container.scrollHeight - containerHeight
  if (maxScrollTop <= 0) return

  const centeredScrollTop = itemTop - (containerHeight - itemHeight) / 2
  container.scrollTo({ top: Math.min(Math.max(centeredScrollTop, 0), maxScrollTop), behavior: 'auto' })
}

/**
 * @param containerRef  包含 data-toc-id 按钮的可滚动容器引用
 * @param activeId      当前高亮条目的 id
 * @param items         TOC 条目列表（用于监听列表变化时也触发同步）
 */
export function useScrollToActive(
  containerRef: Ref<HTMLElement | null>,
  activeId: Ref<string>,
  items: Ref<TocItem[]>,
) {
  function sync() {
    const id = activeId.value
    const container = containerRef.value
    if (!id || !container) return
    nextTick(() => scrollItemIntoView(container, id))
  }

  watch(activeId, sync, { flush: 'post', immediate: true })
  watch(() => items.value.length, sync, { flush: 'post' })
}
