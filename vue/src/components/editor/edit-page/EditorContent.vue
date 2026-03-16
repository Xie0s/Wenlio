<!--
  EditorContent.vue - 编辑器中间内容区
  职责：Markdown 原文编辑区 + 实时预览区（分栏显示），自含图片粘贴/拖拽上传
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useEditorInject } from '@/composables/useThemeEditor'
import { useUploadStore } from '@/stores/upload'
import { PROSE_CLASSES } from '@/lib/prose'

const editor = useEditorInject()
const uploadStore = useUploadStore()
const textareaRef = ref<HTMLTextAreaElement | null>(null)

function insertAtCursor(text: string) {
  const el = textareaRef.value
  if (!el) { editor.content.value += text; return }
  const start = el.selectionStart
  const end = el.selectionEnd
  editor.content.value = editor.content.value.slice(0, start) + text + editor.content.value.slice(end)
  nextTick(() => { const pos = start + text.length; el.setSelectionRange(pos, pos); el.focus() })
}

function handlePaste(e: ClipboardEvent) {
  const fileItems = Array.from(e.clipboardData?.items ?? []).filter(item => item.kind === 'file')
  if (!fileItems.length) return
  const files = fileItems.map(item => item.getAsFile()).filter((f): f is File => f !== null)
  if (!files.length) return
  e.preventDefault()
  uploadStore.requestUploadWithPreview(files, {
    onComplete: (media) => {
      if (media.mime_type.startsWith('image/')) {
        insertAtCursor(`![${media.file_name}](${media.file_url})\n`)
      } else {
        insertAtCursor(`[${media.file_name}](${media.file_url})\n`)
      }
    },
  })
}

function handleDrop(e: DragEvent) {
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return
  e.stopPropagation()
  uploadStore.requestUploadWithPreview(Array.from(files), {
    onComplete: (media) => {
      if (media.mime_type.startsWith('image/')) {
        insertAtCursor(`![${media.file_name}](${media.file_url})\n`)
      } else {
        insertAtCursor(`[${media.file_name}](${media.file_url})\n`)
      }
    },
  })
}
</script>

<template>
  <div class="flex flex-1 min-w-0 min-h-0 h-full overflow-hidden">
    <!-- Markdown 编辑区 -->
    <div class="flex-1 min-w-0 relative">
      <textarea ref="textareaRef" :value="editor.content.value" placeholder="在此编写 Markdown 内容..."
        class="scrollbar-visible h-full w-full resize-none border-0 bg-background p-4 font-mono text-sm leading-relaxed outline-none"
        @input="editor.content.value = ($event.target as HTMLTextAreaElement).value" @paste="handlePaste"
        @drop.prevent="handleDrop" @dragover.prevent />
      <!-- 上传中蒙层 -->
      <div v-if="uploadStore.hasActiveTasks"
        class="absolute inset-0 flex items-center justify-center bg-background/60 backdrop-blur-sm">
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          上传中 {{ uploadStore.totalProgress }}%
        </div>
      </div>
    </div>

    <!-- 实时预览区 -->
    <div v-if="editor.showPreview.value" class="scrollbar-visible flex-1 min-w-0 overflow-y-auto border-l bg-card p-6">
      <div :class="PROSE_CLASSES" v-html="editor.previewHtml.value" />
      <p v-if="!editor.previewHtml.value" class="text-center text-sm text-muted-foreground py-12">预览区域</p>
    </div>
  </div>
</template>
