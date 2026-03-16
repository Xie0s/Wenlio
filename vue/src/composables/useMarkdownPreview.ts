/**
 * composables/useMarkdownPreview.ts - 防抖 Markdown 预览渲染
 * 职责：监听内容变化，延迟渲染 HTML，支持开关控制
 * 对外暴露：useMarkdownPreview(content, options?) → { html, enabled }
 */

import { ref, watch, type Ref, onScopeDispose } from 'vue'
import { renderMarkdown } from '@/lib/markdown'

export function useMarkdownPreview(content: Ref<string>, options?: { delay?: number }) {
  const html = ref('')
  const enabled = ref(true)
  let timer: ReturnType<typeof setTimeout> | null = null

  function update() {
    if (!enabled.value) return
    if (timer) clearTimeout(timer)
    timer = setTimeout(async () => {
      html.value = await renderMarkdown(content.value || '')
    }, options?.delay ?? 500)
  }

  watch(content, update)
  watch(enabled, (v) => { if (v) update() })

  onScopeDispose(() => { if (timer) clearTimeout(timer) })

  return { html, enabled }
}
