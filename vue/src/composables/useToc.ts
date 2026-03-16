/**
 * composables/useToc.ts - 文档目录（TOC）提取
 * 职责：从渲染后的 Markdown HTML 中提取标题层级，生成目录数据
 * 对外暴露：useToc()
 */

import { ref, watch, nextTick, onScopeDispose } from 'vue'

export interface TocItem {
  id: string
  text: string
  level: number
}

/**
 * 从渲染后的文档容器中提取 h2~h4 标题生成 TOC
 */
export function useToc(renderedHtml: () => string) {
  const tocItems = ref<TocItem[]>([])
  const activeId = ref('')
  const headingElements = ref<HTMLElement[]>([])

  function updateActiveIdByPosition() {
    const headings = headingElements.value
    if (headings.length === 0) {
      activeId.value = ''
      return
    }

    const markerTop = 96
    let current = headings[0]!

    for (const heading of headings) {
      if (heading.getBoundingClientRect().top - markerTop <= 0) {
        current = heading
      } else {
        break
      }
    }

    activeId.value = current.id
  }

  let observer: IntersectionObserver | null = null

  function cleanup() {
    if (observer) observer.disconnect()
    observer = null
  }

  function setupFromDom() {
    nextTick(() => {
      const container = document.querySelector('.prose')
      if (!container) {
        tocItems.value = []
        activeId.value = ''
        headingElements.value = []
        cleanup()
        return
      }

      const headings = Array.from(
        container.querySelectorAll<HTMLElement>('h2[id], h3[id], h4[id]'),
      )
      headingElements.value = headings

      const idCount: Record<string, number> = {}
      tocItems.value = headings.map((el) => {
        const rawId = el.id
        idCount[rawId] = (idCount[rawId] ?? 0) + 1
        const uniqueId = idCount[rawId] === 1 ? rawId : `${rawId}-${idCount[rawId]}`
        const clone = el.cloneNode(true) as HTMLElement
        clone.querySelectorAll('.header-anchor').forEach(a => a.remove())
        return {
          id: uniqueId,
          text: clone.textContent?.trim() || '',
          level: parseInt(el.tagName[1]!, 10),
        }
      })

      cleanup()
      if (headings.length === 0) return

      observer = new IntersectionObserver(
        () => {
          updateActiveIdByPosition()
        },
        { rootMargin: '-80px 0px -60% 0px', threshold: 0 },
      )
      headings.forEach((el) => observer!.observe(el))

      updateActiveIdByPosition()
    })
  }

  // 监听渲染内容变化时：一次性完成 TOC 提取 + 滚动高亮绑定
  watch(renderedHtml, setupFromDom, { flush: 'post' })

  /**
   * 查找最近的可滚动祖先元素（overflow auto/scroll）
   * 对源码编辑器预览区和可视化编辑器内部滚动容器均有效
   */
  function findScrollContainer(el: HTMLElement): HTMLElement | null {
    let parent = el.parentElement
    while (parent && parent !== document.documentElement) {
      const { overflow, overflowY } = getComputedStyle(parent)
      if (/(auto|scroll)/.test(overflow) || /(auto|scroll)/.test(overflowY)) {
        return parent
      }
      parent = parent.parentElement
    }
    return null
  }

  function scrollTo(id: string) {
    const el = document.getElementById(id)
    if (!el) return

    activeId.value = id

    const PADDING = 16
    const scrollContainer = findScrollContainer(el)

    if (scrollContainer) {
      // 滚动内部容器（编辑器预览区 / 可视化编辑器）
      const containerRect = scrollContainer.getBoundingClientRect()
      const elRect = el.getBoundingClientRect()
      const relativeTop = elRect.top - containerRect.top + scrollContainer.scrollTop - PADDING
      scrollContainer.scrollTo({ top: relativeTop, behavior: 'smooth' })
    } else {
      // 页面滚动（阅读页等场景）
      const HEADER_HEIGHT = 56
      const top = el.getBoundingClientRect().top + window.scrollY - HEADER_HEIGHT - PADDING
      window.scrollTo({ top, behavior: 'smooth' })
    }

    // contenteditable 内（可视化编辑器）跳过 DOM 高亮操作，避免破坏 ProseMirror 节点状态
    if (el.closest('[contenteditable]')) return

    // 用 span 包裹标题文字，实现文字宽度的黄色高亮
    const highlightSpan = document.createElement('span')
    highlightSpan.className = 'toc-heading-highlight'
    while (el.firstChild) {
      highlightSpan.appendChild(el.firstChild)
    }
    el.appendChild(highlightSpan)

    setTimeout(() => {
      highlightSpan.classList.add('toc-heading-highlight-fade')
      setTimeout(() => {
        while (highlightSpan.firstChild) {
          el.insertBefore(highlightSpan.firstChild, highlightSpan)
        }
        highlightSpan.remove()
      }, 800)
    }, 1800)
  }

  onScopeDispose(cleanup)

  return { tocItems, activeId, scrollTo, cleanup, refresh: setupFromDom }
}
