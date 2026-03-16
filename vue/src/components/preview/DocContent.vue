<!-- components/DocPage/DocContent.vue
     职责：文档正文渲染区域（标题 + Markdown HTML 输出）
     纯展示组件，不含加载状态，由父组件控制显隐；正文内图片禁止右键保存、拖拽、复制
     对外暴露：无（slot-less 纯渲染） -->

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { PROSE_CLASSES } from '@/lib/prose'
import { toast } from 'vue-sonner'

const ICON_COPY = '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>'
const ICON_CHECK = '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/><path d="m12 15 2 2 4-4"/></svg>'

const props = defineProps<{
  html: string
}>()

const contentEl = ref<HTMLElement | null>(null)

function injectCopyButtons(container: HTMLElement) {
  container.querySelectorAll('pre').forEach(pre => {
    if (pre.querySelector('.code-copy-btn')) return
    pre.style.position = 'relative'
    const btn = document.createElement('button')
    btn.className = 'code-copy-btn'
    btn.innerHTML = ICON_COPY
    btn.addEventListener('click', () => {
      const code = pre.querySelector('code')?.textContent ?? pre.textContent ?? ''
      navigator.clipboard.writeText(code).then(() => {
        btn.innerHTML = ICON_CHECK
        toast.success('复制成功')
        setTimeout(() => { btn.innerHTML = ICON_COPY }, 1500)
      })
    })
    pre.appendChild(btn)
  })
}

watch(() => props.html, () => {
  nextTick(() => { if (contentEl.value) injectCopyButtons(contentEl.value) })
}, { immediate: true })

/** 正文图片防护：禁止右键另存、拖拽保存、复制到剪贴板 */
function onContextMenu(e: MouseEvent) {
  if ((e.target as HTMLElement)?.closest?.('img')) e.preventDefault()
}

function onDragStart(e: DragEvent) {
  if ((e.target as HTMLElement)?.closest?.('img')) e.preventDefault()
}

function onCopy(e: ClipboardEvent) {
  const container = (e.target as HTMLElement)?.closest?.('.doc-content-reader')
  if (!container) return
  const sel = document.getSelection()
  if (!sel?.rangeCount) return
  const imgs = container.querySelectorAll('img')
  for (let i = 0; i < sel.rangeCount; i++) {
    const range = sel.getRangeAt(i)
    for (const img of imgs) {
      try {
        if (range.intersectsNode(img)) {
          e.preventDefault()
          if (e.clipboardData) e.clipboardData.clearData()
          return
        }
      } catch {
        /* ignore */
      }
    }
  }
}
</script>

<template>
  <article class="doc-content-reader reader-img-protect" @contextmenu="onContextMenu" @dragstart="onDragStart"
    @copy="onCopy">
    <!-- Markdown 渲染内容 -->
    <div ref="contentEl" :class="PROSE_CLASSES" v-html="html" />
  </article>
</template>

<style scoped>
.reader-img-protect :deep(img) {
  user-select: none;
  -webkit-user-select: none;
  -webkit-user-drag: none;
  pointer-events: auto;
}

:deep(.code-copy-btn) {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid oklch(0 0 0 / 9%);
  background: oklch(1 0 0 / 72%);
  color: var(--color-muted-foreground);
  backdrop-filter: blur(8px);
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.dark :deep(.code-copy-btn) {
  background: oklch(0.2 0 0 / 60%);
  border-color: oklch(1 0 0 / 10%);
}

:deep(.code-copy-btn:hover) {
  background: oklch(0.95 0 0);
  color: var(--color-foreground);
}

.dark :deep(.code-copy-btn:hover) {
  background: oklch(0.3 0 0);
}
</style>
