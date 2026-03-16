<!--
  VisualEditorToolbar.vue - 可视化编辑器格式化工具栏
  职责：提供 Word 风格的富文本格式化操作按钮，分组布局
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { type Editor } from '@tiptap/vue-3'
import { openOverlaySafely, releaseOverlayFocus } from '@/lib/overlay-focus'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import {
  Undo2, Redo2,
  Bold, Italic, Strikethrough, Code,
  List, ListOrdered, ListChecks,
  Quote, FileCode, Minus, TableIcon,
  Link2, ImagePlus, Paperclip, FolderOpen,
  Info, AlertTriangle, ShieldAlert, Lightbulb,
} from 'lucide-vue-next'
import { CONTAINER_TYPES, CONTAINER_TITLES, type ContainerType } from './extensions/CustomContainer'
import type { FileNodeAttrs } from './extensions/FileNode'
import { LinkInputPopover, TableInsertDialog, FilePickerDialog } from './common'
import { UploadPreviewDialog } from '@/components/editor/media/update'
import { useUploadStore } from '@/stores/upload'
import type { MediaItem } from '@/lib/media'

const props = defineProps<{ editor: Editor | undefined }>()

const uploadStore = useUploadStore()

function getOverlayFocusOptions() {
  return {
    editorElement: props.editor?.view?.dom as HTMLElement | null,
    blurEditor: props.editor ? () => props.editor?.commands.blur() : null,
    settleFrameCount: 2,
  }
}

function blurEditorBeforeOverlay() {
  releaseOverlayFocus(getOverlayFocusOptions())
}

// ── 链接 Popover ──
const linkPopoverOpen = ref(false)

const selectionEmpty = computed(() => props.editor?.state.selection.empty ?? true)

const selectedText = computed(() => {
  if (!props.editor || selectionEmpty.value) return ''
  const { from, to } = props.editor.state.selection
  return props.editor.state.doc.textBetween(from, to)
})

function handleLinkConfirm(url: string, text?: string) {
  if (!props.editor) return
  if (props.editor.state.selection.empty) {
    const displayText = text || url
    props.editor
      .chain()
      .focus()
      .insertContent({ type: 'text', text: displayText, marks: [{ type: 'link', attrs: { href: url } }] })
      .run()
  } else {
    props.editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  }
}

function handleLinkRemove() {
  if (!props.editor) return
  props.editor.chain().focus().unsetLink().run()
}

// ── 图片上传（通过预览弹窗） ──
const imageUploadOpen = ref(false)

function openImageUpload() {
  void openOverlaySafely(() => {
    imageUploadOpen.value = true
  }, getOverlayFocusOptions())
}

function handleImageUploadConfirm(files: File[]) {
  if (!props.editor || files.length === 0) return
  uploadStore.addTasks(files, {
    onComplete: (media) => {
      if (!props.editor) return
      if (media.mime_type.startsWith('image/')) {
        props.editor.chain().focus().setImage({
          src: media.file_url, alt: media.file_name, 'data-media-id': media.id,
        } as any).run()
      }
    },
  })
}

// ── 文件上传（通过预览弹窗） ──
const fileUploadOpen = ref(false)

function openFileUpload() {
  void openOverlaySafely(() => {
    fileUploadOpen.value = true
  }, getOverlayFocusOptions())
}

function handleFileUploadConfirm(files: File[]) {
  if (!props.editor || files.length === 0) return
  uploadStore.addTasks(files, {
    onComplete: (media) => {
      if (!props.editor) return
      const attrs: Omit<FileNodeAttrs, 'layout'> = {
        href: media.file_url,
        fileName: media.file_name,
        fileSize: media.file_size,
        mimeType: media.mime_type,
        mediaId: media.id,
      }
      nextTick(() => {
        props.editor!.chain().focus().insertFileNode(attrs).run()
      })
    },
  })
}

// ── 插入已有文件 ──
const filePickerOpen = ref(false)

function openFilePicker() {
  void openOverlaySafely(() => {
    filePickerOpen.value = true
  }, getOverlayFocusOptions())
}

function handleFilePick(picked: MediaItem[], insertType: 'image' | 'file') {
  if (!props.editor || picked.length === 0) return
  if (insertType === 'image') {
    const imageItems = picked.filter(m => m.mime_type?.startsWith('image/'))
    const fileItems = picked.filter(m => !m.mime_type?.startsWith('image/'))
    nextTick(() => {
      // image 为 block 节点，必须用 focus('end') 确保每次都追加到末尾，否则 focus() 恢复上次选区会替换已插入的节点
      imageItems.forEach(media => {
        props.editor!.chain().focus('end').insertContent({
          type: 'image',
          attrs: {
            src: media.file_url,
            alt: media.file_name,
            'data-media-id': media.id,
          },
        }).run()
      })
      // 混选中的非图片文件以 fileNode 形式插入，而非强制按图片处理
      if (fileItems.length > 0) {
        props.editor!.chain().focus('end').insertContent(
          fileItems.map(media => ({
            type: 'fileNode',
            attrs: {
              href: media.file_url,
              fileName: media.file_name,
              fileSize: media.file_size,
              mimeType: media.mime_type,
              mediaId: media.id,
              layout: 'full' as const,
            },
          }))
        ).run()
      }
    })
  } else {
    const nodes = picked.map(media => ({
      type: 'fileNode',
      attrs: {
        href: media.file_url,
        fileName: media.file_name,
        fileSize: media.file_size,
        mimeType: media.mime_type,
        mediaId: media.id,
        layout: 'full' as const,
      },
    }))
    // 延迟到下一 tick，确保 Dialog 焦点陷阱释放后编辑器能正确获焦
    nextTick(() => {
      props.editor!.chain().focus().insertContent(nodes).run()
    })
  }
}

// ── 表格插入对话框 ──
const tableDialogOpen = ref(false)

function openTableDialog() {
  void openOverlaySafely(() => {
    tableDialogOpen.value = true
  }, getOverlayFocusOptions())
}

function handleTableInsert(rows: number, cols: number, withHeaderRow: boolean) {
  if (!props.editor) return
  props.editor.chain().focus().insertTable({ rows, cols, withHeaderRow }).run()
}

// ── 代码块 ──
function toggleCodeBlock() {
  if (!props.editor) return
  props.editor.chain().focus().toggleCodeBlock().run()
}

// ── 标题级别 Select ──
const headingSelectOpen = ref(false)

watch(headingSelectOpen, (open) => {
  if (open) blurEditorBeforeOverlay()
})

const headingLevel = computed(() => {
  if (!props.editor) return 'paragraph'
  for (const level of [1, 2, 3, 4] as const) {
    if (props.editor.isActive('heading', { level })) return `h${level}`
  }
  return 'paragraph'
})

function setHeading(value: any) {
  if (!props.editor || !value) return
  const v = String(value)
  if (v === 'paragraph') {
    props.editor.chain().focus().setParagraph().run()
  } else {
    const level = parseInt(v.replace('h', '')) as 1 | 2 | 3 | 4
    props.editor.chain().focus().toggleHeading({ level }).run()
  }
}

// ── 容器图标映射 ──
const containerIcons: Record<ContainerType, typeof Info> = {
  tip: Lightbulb,
  warning: AlertTriangle,
  danger: ShieldAlert,
  info: Info,
}

function insertContainer(type: ContainerType) {
  if (!props.editor) return
  props.editor.chain().focus().setContainer(type).run()
}

// ── 通用工具函数 ──
function isActive(name: string, attrs?: Record<string, any>): boolean {
  return props.editor?.isActive(name, attrs) ?? false
}

function canUndo(): boolean { return props.editor?.can().undo() ?? false }
function canRedo(): boolean { return props.editor?.can().redo() ?? false }

function toggleItalic() {
  if (!props.editor) return
  const ok = props.editor.chain().focus().toggleItalic().run()
  if (!ok) {
    console.debug('[VisualEditorToolbar] toggleItalic failed', {
      isEditable: props.editor.isEditable,
      isFocused: props.editor.isFocused,
      selection: props.editor.state.selection,
    })
  }
}
</script>

<template>
  <div v-if="editor"
    class="flex items-center justify-center gap-0.5 px-3 py-1.5 border-b border-border bg-background overflow-x-auto flex-none">

    <!-- 撤销 / 重做 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!canUndo()"
          @click="editor!.chain().focus().undo().run()">
          <Undo2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>撤销</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!canRedo()"
          @click="editor!.chain().focus().redo().run()">
          <Redo2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>重做</TooltipContent>
    </Tooltip>

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 标题选择 -->
    <Select :model-value="headingLevel" :open="headingSelectOpen" @update:open="headingSelectOpen = $event"
      @update:model-value="setHeading">
      <SelectTrigger class="h-8 w-24 text-xs" @pointerdown="blurEditorBeforeOverlay">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="paragraph">正文</SelectItem>
        <SelectItem value="h1">标题 1</SelectItem>
        <SelectItem value="h2">标题 2</SelectItem>
        <SelectItem value="h3">标题 3</SelectItem>
        <SelectItem value="h4">标题 4</SelectItem>
      </SelectContent>
    </Select>

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 文本格式 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('bold') }" @click="editor!.chain().focus().toggleBold().run()">
          <Bold class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>加粗</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('italic') }" @mousedown.prevent="toggleItalic">
          <Italic class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>斜体</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('strike') }" @click="editor!.chain().focus().toggleStrike().run()">
          <Strikethrough class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>删除线</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('code') }" @click="editor!.chain().focus().toggleCode().run()">
          <Code class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>行内代码</TooltipContent>
    </Tooltip>

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 列表 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('bulletList') }"
          @click="editor!.chain().focus().toggleBulletList().run()">
          <List class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>无序列表</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('orderedList') }"
          @click="editor!.chain().focus().toggleOrderedList().run()">
          <ListOrdered class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>有序列表</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('taskList') }" @click="editor!.chain().focus().toggleTaskList().run()">
          <ListChecks class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>任务列表</TooltipContent>
    </Tooltip>

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 插入 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('blockquote') }"
          @click="editor!.chain().focus().toggleBlockquote().run()">
          <Quote class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>引用块</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('codeBlock') }" @click="toggleCodeBlock">
          <FileCode class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>代码块</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @click="editor!.chain().focus().setHorizontalRule().run()">
          <Minus class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>分隔线</TooltipContent>
    </Tooltip>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('table') }" @pointerdown="blurEditorBeforeOverlay"
          @click="openTableDialog">
          <TableIcon class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>插入表格</TooltipContent>
    </Tooltip>

    <!-- 表格插入对话框 -->
    <TableInsertDialog v-model:open="tableDialogOpen" @insert="handleTableInsert" />

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 链接 -->
    <LinkInputPopover v-model:open="linkPopoverOpen" :initial-url="editor?.getAttributes('link')?.href || ''"
      :initial-text="selectedText" :show-text-input="selectionEmpty && !isActive('link')"
      :is-edit-mode="isActive('link')" @confirm="handleLinkConfirm" @remove="handleLinkRemove">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
            :class="{ 'bg-foreground/10': isActive('link') }" @click="linkPopoverOpen = !linkPopoverOpen">
            <Link2 class="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent v-if="!linkPopoverOpen">链接</TooltipContent>
      </Tooltip>
    </LinkInputPopover>

    <!-- 图片上传 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @pointerdown="blurEditorBeforeOverlay" @click="openImageUpload">
          <ImagePlus class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>插入图片</TooltipContent>
    </Tooltip>
    <UploadPreviewDialog
      v-model:open="imageUploadOpen"
      accept="image/*"
      title="上传图片"
      @confirm="handleImageUploadConfirm"
    />

    <!-- 文件上传 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @pointerdown="blurEditorBeforeOverlay" @click="openFileUpload">
          <Paperclip class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>上传文件</TooltipContent>
    </Tooltip>
    <UploadPreviewDialog
      v-model:open="fileUploadOpen"
      title="上传文件"
      @confirm="handleFileUploadConfirm"
    />

    <!-- 插入已有文件 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @pointerdown="blurEditorBeforeOverlay" @click="openFilePicker">
          <FolderOpen class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>插入已有文件</TooltipContent>
    </Tooltip>
    <FilePickerDialog v-model:open="filePickerOpen" @select="handleFilePick" />

    <div class="h-5 w-px bg-border mx-1" />

    <!-- 自定义容器 -->
    <Tooltip v-for="type in CONTAINER_TYPES" :key="type">
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': isActive('customContainer', { type }) }" @click="insertContainer(type)">
          <component :is="containerIcons[type]" class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>{{ CONTAINER_TITLES[type] }}容器</TooltipContent>
    </Tooltip>
  </div>
</template>
