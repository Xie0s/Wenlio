<!--
  VisualEditorWorkspace.vue - 可视化编辑器工作区
  职责：集成格式化工具栏 + Tiptap 编辑器内容区，提供所见即所得的文档编辑体验
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { watch } from 'vue'
import { useEditorInject } from '@/composables/useThemeEditor'
import { releaseOverlayFocus } from '@/lib/overlay-focus'
import { useVisualEditor } from './useVisualEditor'
import VisualEditorToolbar from './VisualEditorToolbar.vue'
import VisualEditorContent from './VisualEditorContent.vue'
import { LinkMenu, ImageMenu, TableMenu, CodeBlockMenu, FileCardMenu } from './common'
import EditorEmptyGuide from '../EditorEmptyGuide.vue'

const editorCtx = useEditorInject()

const { editor, isReady } = useVisualEditor({
  content: editorCtx.content,
  placeholder: '开始编写内容...',
})

// 当 Dialog / Sheet 等覆盖层打开时，reka-ui 会在页面根元素设置 aria-hidden="true"，
// 与 tiptap 保持焦点冲突产生浏览器警告。打开前先 blur 编辑器即可规避。
watch(() => editorCtx.showSettingsPanel.value, (open) => {
  if (open && editor.value) {
    releaseOverlayFocus({
      editorElement: editor.value.view?.dom as HTMLElement | null,
      blurEditor: () => editor.value?.commands.blur(),
    })
  }
}, { flush: 'sync' })
</script>

<template>
  <!-- 有页面：可视化编辑区（工具栏 + 内容） -->
  <div v-if="editorCtx.page.value" class="flex flex-col flex-1 min-w-0 min-h-0 h-full">
    <!-- 格式化工具栏 -->
    <VisualEditorToolbar :editor="editor" />

    <!-- 编辑器内容区 -->
    <VisualEditorContent :editor="editor" />

    <!-- 浮动菜单：链接 / 图片 / 表格 -->
    <template v-if="editor">
      <LinkMenu :editor="editor" />
      <ImageMenu :editor="editor" />
      <TableMenu :editor="editor" />
      <CodeBlockMenu :editor="editor" />
      <FileCardMenu :editor="editor" />
    </template>

    <!-- 编辑器加载态 -->
    <div v-if="!isReady"
      class="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm z-20">
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        编辑器加载中...
      </div>
    </div>
  </div>

  <!-- 空态引导 -->
  <EditorEmptyGuide v-else />
</template>
