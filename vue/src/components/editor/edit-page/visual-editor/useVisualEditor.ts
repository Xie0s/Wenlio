/**
 * useVisualEditor.ts - 可视化编辑器 Tiptap 实例管理 composable
 * 职责：
 *   1) 创建并配置 Tiptap 编辑器实例
 *   2) 处理 Markdown ↔ ProseMirror 双向同步（含自定义容器预处理/后处理）
 *   3) 管理编辑器焦点、销毁等生命周期
 * 对外暴露：useVisualEditor() → { editor, isReady }
 */

import { ref, watch, onBeforeUnmount, type Ref } from 'vue'
import { useEditor } from '@tiptap/vue-3'
import { createExtensions } from './extensions'
import { CONTAINER_TITLES, type ContainerType } from './extensions/CustomContainer'

/**
 * 将 :::type 语法转为 HTML，以便 tiptap-markdown 通过 HTML 解析为自定义节点
 * 支持 :::type 或 :::type 自定义标题 两种格式
 */
function preprocessContainers(markdown: string): string {
  return markdown.replace(
    /^:::(tip|warning|danger|info)\s*(.*?)$\n([\s\S]*?)^:::$/gm,
    (_match, type: string, customTitle: string, content: string) => {
      const title = customTitle.trim() || CONTAINER_TITLES[type as ContainerType] || type
      const trimmedContent = content.trim()
      return `<div class="custom-container ${type}"><p class="custom-container-title">${title}</p><div class="custom-container-content">\n\n${trimmedContent}\n\n</div></div>`
    },
  )
}

/**
 * 将 tiptap-markdown 输出中的 HTML 容器块还原为 ::: 语法
 */
function postprocessContainers(markdown: string): string {
  // 匹配 tiptap-markdown 序列化输出的容器格式
  // tiptap-markdown 会将自定义节点通过 storage.markdown.serialize 输出
  // 如果 serialize 正常工作，输出已经是 ::: 格式，无需后处理
  // 此函数作为兜底，处理 HTML 形式的容器
  return markdown.replace(
    /<div class="custom-container (tip|warning|danger|info)">\s*<p class="custom-container-title">(.*?)<\/p>\s*<div class="custom-container-content">([\s\S]*?)<\/div>\s*<\/div>/g,
    (_match, type: string, title: string, content: string) => {
      const defaultTitle = CONTAINER_TITLES[type as ContainerType]
      const titlePart = title && title !== defaultTitle ? ` ${title}` : ''
      return `:::${type}${titlePart}\n${content.trim()}\n:::`
    },
  )
}

export interface VisualEditorOptions {
  /** 双向绑定的 markdown 内容 ref */
  content: Ref<string>
  /** 编辑器占位文本 */
  placeholder?: string
  /** 内容变更防抖延迟（ms） */
  debounceMs?: number
}

export function useVisualEditor(options: VisualEditorOptions) {
  const { content, placeholder, debounceMs = 300 } = options
  const isReady = ref(false)

  // 防止 onUpdate 回写时触发 watch → setContent 循环
  let isUpdatingFromEditor = false
  let debounceTimer: ReturnType<typeof setTimeout> | null = null

  const editor = useEditor({
    extensions: createExtensions({ placeholder }),
    // 初始内容通过 preprocessContainers 转换后以 markdown 格式加载
    content: preprocessContainers(content.value),
    // 解析格式：tiptap-markdown 会自动识别
    onUpdate: ({ editor: e }: { editor: any }) => {
      if (debounceTimer) clearTimeout(debounceTimer)
      debounceTimer = setTimeout(() => {
        isUpdatingFromEditor = true
        const md = getMarkdownFromEditor(e)
        content.value = md
        isUpdatingFromEditor = false
      }, debounceMs)
    },
    onCreate: () => {
      isReady.value = true
    },
  })

  /**
   * 从 Tiptap 编辑器实例获取 Markdown 字符串
   */
  function getMarkdownFromEditor(e: any): string {
    try {
      const storage = e.storage as Record<string, any>
      const md = storage.markdown?.getMarkdown?.()
      if (typeof md === 'string') {
        return postprocessContainers(md)
      }
    } catch (error) {
      if (import.meta.env.DEV) {
        console.warn('[visual-editor:getMarkdownFromEditor] markdown serialize failed, fallback to getText()', error)
      }
      // fallback: 返回纯文本
    }
    return e.getText()
  }

  // 外部 content 变化时同步到编辑器（如切换文档页）
  watch(content, (newContent) => {
    if (isUpdatingFromEditor) return
    if (!editor.value) return

    const currentMd = getMarkdownFromEditor(editor.value)
    // 仅当内容实质性变化时才更新（避免不必要的 setContent）
    if (currentMd.trim() === newContent.trim()) return

    editor.value.commands.setContent(preprocessContainers(newContent))
  })

  onBeforeUnmount(() => {
    if (debounceTimer) clearTimeout(debounceTimer)
  })

  return {
    editor,
    isReady,
  }
}
