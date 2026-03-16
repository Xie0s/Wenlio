<!--
  VisualEditorContent.vue - 可视化编辑器内容区
  职责：渲染 Tiptap 编辑器实例，应用与阅读页一致的 prose + markdown.css 样式
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { type Editor, EditorContent } from '@tiptap/vue-3'
import { TextSelection } from '@tiptap/pm/state'
import { useUploadStore } from '@/stores/upload'
import { type MediaItem } from '@/lib/media'
import { PROSE_CLASSES } from '@/lib/prose'

const props = defineProps<{ editor: Editor | undefined }>()

const uploadStore = useUploadStore()

/** 将上传完成的媒体插入编辑器（图片用 setImage，其他用 insertFileNode）。
 *  插入 atom 节点后立即将光标移至节点之后，避免 Tiptap NodeSelection 使下次插入覆盖当前节点。 */
function insertMedia(media: MediaItem) {
  if (!props.editor) return
  if (media.mime_type.startsWith('image/')) {
    props.editor.chain().focus().setImage({
      src: media.file_url, alt: media.file_name, 'data-media-id': media.id,
    } as any).run()
  } else {
    props.editor.chain().focus().insertFileNode({
      href: media.file_url,
      fileName: media.file_name,
      fileSize: media.file_size,
      mimeType: media.mime_type,
      mediaId: media.id,
    }).run()
  }
  props.editor.commands.focus(props.editor.state.selection.to)
}

/** 拖拽文件到编辑区（捕获阶段处理，先于 ProseMirror 内部处理器执行） */
function handleDrop(e: DragEvent) {
  const files = Array.from(e.dataTransfer?.files ?? [])
  if (!files.length) return
  e.preventDefault()
  e.stopPropagation()
  uploadStore.requestUploadWithPreview(files, { onComplete: insertMedia })
}

/**
 * 剪贴板文件/图片粘贴（截图、浏览器复制图片、文件管理器复制文件等）
 * 使用 capture 阶段提前拦截，防止 Tiptap 将图片数据当作文本处理。
 * 仅当剪贴板含有 file 类型时才接管，否则放行让 Tiptap 的插件链处理。
 */
function handlePaste(e: ClipboardEvent) {
  const fileItems = Array.from(e.clipboardData?.items ?? []).filter(item => item.kind === 'file')
  if (!fileItems.length) return

  e.preventDefault()
  e.stopPropagation()

  const files = fileItems.map(item => item.getAsFile()).filter((f): f is File => f !== null)
  if (!files.length) return

  uploadStore.requestUploadWithPreview(files, { onComplete: insertMedia })
}

function handleEditorBackgroundMouseDown(e: MouseEvent) {
  if (!props.editor || uploadStore.hasActiveTasks) return
  const editorElement = props.editor.view?.dom as HTMLElement | null
  const target = e.target as Node | null
  if (!editorElement || !target) return
  if (editorElement.contains(target)) return
  e.preventDefault()
  const { state, view } = props.editor
  let lastTextSelectionPos: number | null = null
  state.doc.descendants((node, pos) => {
    if (!node.isTextblock) return true
    lastTextSelectionPos = pos + node.nodeSize - 1
    return true
  })
  if (lastTextSelectionPos !== null) {
    view.dispatch(state.tr.setSelection(TextSelection.create(state.doc, lastTextSelectionPos)).scrollIntoView())
    view.focus()
    return
  }
  const paragraphType = state.schema.nodes.paragraph
  if (!paragraphType) {
    view.focus()
    return
  }
  const tr = state.tr
  const insertPos = tr.doc.content.size
  tr.insert(insertPos, paragraphType.create())
  tr.setSelection(TextSelection.create(tr.doc, insertPos + 1))
  view.dispatch(tr.scrollIntoView())
  view.focus()
}
</script>

<template>
  <div class="scrollbar-visible relative flex-1 min-h-0 overflow-y-auto bg-background cursor-text" @drop.capture="handleDrop"
    @dragover.prevent @paste.capture="handlePaste">
    <!-- Tiptap 编辑器渲染区，prose 样式与阅读页完全一致 -->
    <div class="visual-editor-wrapper mx-auto w-full max-w-4xl min-h-full cursor-text px-8 py-6" @mousedown="handleEditorBackgroundMouseDown">
      <EditorContent v-if="editor" :editor="editor" :class="PROSE_CLASSES" />
    </div>

    <!-- 上传中蒙层 -->
    <div v-if="uploadStore.hasActiveTasks"
      class="absolute inset-0 flex items-center justify-center bg-background/60 backdrop-blur-sm z-10">
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        上传中 {{ uploadStore.totalProgress }}%
      </div>
    </div>
  </div>
</template>

<style>
/*
 * Tiptap ProseMirror 编辑器基础样式
 * 仅设置编辑态必要样式，不覆盖 prose / markdown.css 中的任何规则
 */
.visual-editor-wrapper .ProseMirror {
  outline: none;
  min-height: 300px;
  cursor: text;
}

.visual-editor-wrapper {
  min-height: 100%;
}

.visual-editor-wrapper > div {
  min-height: 100%;
}

/* 占位文本样式 */
.visual-editor-wrapper .ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: var(--color-muted-foreground);
  pointer-events: none;
  height: 0;
  opacity: 0.5;
}

/* 表格编辑态：选中单元格高亮 */
.visual-editor-wrapper .ProseMirror .selectedCell::after {
  content: '';
  position: absolute;
  inset: 0;
  background: oklch(0.6 0.15 250 / 0.15);
  pointer-events: none;
  z-index: 2;
}

.visual-editor-wrapper .ProseMirror table td,
.visual-editor-wrapper .ProseMirror table th {
  position: relative;
}

/* 任务列表样式 —— 按 Tiptap 官方 demo 写法，使用 .tiptap 根选择器 */
.tiptap ul[data-type='taskList'] {
  list-style: none;
  margin-left: 0;
  padding: 0;
}

.tiptap ul[data-type='taskList'] li {
  display: flex;
  align-items: flex-start;
}

.tiptap ul[data-type='taskList'] li>label {
  flex: 0 0 auto;
  margin-right: 0.5rem;
  user-select: none;
}

.tiptap ul[data-type='taskList'] li>div {
  flex: 1 1 auto;
}

.tiptap ul[data-type='taskList'] li>div> :first-child {
  margin-top: 0;
}

.tiptap ul[data-type='taskList'] li>div> :last-child {
  margin-bottom: 0;
}

.tiptap ul[data-type='taskList'] input[type='checkbox'] {
  cursor: pointer;
}

/* 自定义容器在编辑态的可编辑区域 */
.visual-editor-wrapper .ProseMirror .custom-container .custom-container-content {
  min-height: 1.5rem;
}

/* 代码块在编辑态的基础样式（与 shiki 类一致） */
.visual-editor-wrapper .ProseMirror pre {
  padding: 1rem 1.25rem;
  border-radius: 1rem;
  overflow-x: auto;
  font-size: 0.875rem;
  line-height: 1.7;
  border: 1px solid var(--border);
  background-color: var(--color-muted);
  tab-size: 2;
  -moz-tab-size: 2;
  white-space: pre;
}

/* lowlight 语法高亮主题（GitHub Light / Dark） */
.visual-editor-wrapper .ProseMirror pre code .hljs-comment,
.visual-editor-wrapper .ProseMirror pre code .hljs-quote {
  color: #6a737d;
  font-style: italic;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-keyword,
.visual-editor-wrapper .ProseMirror pre code .hljs-selector-tag {
  color: #d73a49;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-string,
.visual-editor-wrapper .ProseMirror pre code .hljs-addition {
  color: #032f62;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-number,
.visual-editor-wrapper .ProseMirror pre code .hljs-literal {
  color: #005cc5;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-built_in,
.visual-editor-wrapper .ProseMirror pre code .hljs-type {
  color: #6f42c1;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-function .hljs-title,
.visual-editor-wrapper .ProseMirror pre code .hljs-title {
  color: #6f42c1;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-attr,
.visual-editor-wrapper .ProseMirror pre code .hljs-variable,
.visual-editor-wrapper .ProseMirror pre code .hljs-template-variable {
  color: #005cc5;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-tag {
  color: #22863a;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-name {
  color: #22863a;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-deletion {
  color: #b31d28;
  background-color: #ffeef0;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-section {
  color: #005cc5;
  font-weight: bold;
}

.visual-editor-wrapper .ProseMirror pre code .hljs-meta {
  color: #735c0f;
}

/* 暗色模式 */
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-comment,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-quote {
  color: #8b949e;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-keyword,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-selector-tag {
  color: #ff7b72;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-string,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-addition {
  color: #a5d6ff;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-number,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-literal {
  color: #79c0ff;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-built_in,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-type {
  color: #d2a8ff;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-function .hljs-title,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-title {
  color: #d2a8ff;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-attr,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-variable,
.dark .visual-editor-wrapper .ProseMirror pre code .hljs-template-variable {
  color: #79c0ff;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-tag {
  color: #7ee787;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-name {
  color: #7ee787;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-deletion {
  color: #ffa198;
  background-color: rgba(248, 81, 73, 0.1);
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-section {
  color: #79c0ff;
  font-weight: bold;
}

.dark .visual-editor-wrapper .ProseMirror pre code .hljs-meta {
  color: #d29922;
}

/* 图片选中态（NodeSelection 单点选中 / AllSelection 全选均适用） */
.visual-editor-wrapper .ProseMirror img.ProseMirror-selectednode,
.visual-editor-wrapper .ProseMirror img.image-in-selection {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

/* ── 文件节点布局 ── */

/* 抑制 ProseMirror 默认 selectednode 轮廓，由 FileCard 内部处理选中态 */
.visual-editor-wrapper .ProseMirror .file-node-wrapper.ProseMirror-selectednode {
  outline: none;
}

/* 基础：full 布局（默认块级，占满一行） */
.visual-editor-wrapper .ProseMirror .file-node-wrapper {
  margin: 0.25rem 0;
}

/* half 布局：一行两个 */
.visual-editor-wrapper .ProseMirror .file-layout-half {
  display: inline-block;
  width: calc(50% - 3px);
  vertical-align: top;
  margin-right: 4px;
}

/* third 布局：一行三个 */
.visual-editor-wrapper .ProseMirror .file-layout-third {
  display: inline-block;
  width: calc(33.333% - 4px);
  vertical-align: top;
  margin-right: 4px;
}
</style>
