<!--
  FileCardMenu.vue - 文件节点浮动菜单组件
  职责：选中文件节点时显示浮动工具栏，提供布局切换、复制链接、新标签页打开、从文档移除节点
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { NodeSelection } from '@tiptap/pm/state'
import {
  ExternalLink, Trash2, Copy,
  RectangleHorizontal, Columns2, Columns3,
} from 'lucide-vue-next'
import { useMenuScheduler, isRectInView, isTopInView } from './useMenuScheduler'
import type { Editor } from '@tiptap/vue-3'
import type { FileLayout } from '../extensions/FileNode'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { toast } from 'vue-sonner'

const props = defineProps<{ editor: Editor }>()

const isVisible = ref(false)
const position = ref({ top: 0, left: 0 })
const menuRef = ref<HTMLDivElement | null>(null)
const lastFilePos = ref<number | null>(null)

const lastState = { isVisible: false, top: 0, left: 0 }

let cachedMenuWidth = 0

function applyState(nextVisible: boolean, nextTop: number, nextLeft: number) {
  const roundedTop = Math.round(nextTop)
  const roundedLeft = Math.round(nextLeft)

  if (lastState.isVisible !== nextVisible) {
    lastState.isVisible = nextVisible
    isVisible.value = nextVisible
  }
  if (!nextVisible) return
  const posChanged = lastState.top !== roundedTop || lastState.left !== roundedLeft
  if (posChanged) {
    lastState.top = roundedTop
    lastState.left = roundedLeft
    position.value = { top: roundedTop, left: roundedLeft }
  }
}

function getFilePos(): number | null {
  const { selection } = props.editor.state
  if (selection instanceof NodeSelection && selection.node.type.name === 'fileNode') {
    return selection.from
  }
  return null
}

function computeMenu() {
  const filePos = getFilePos()
  const isInteracting = !!(menuRef.value && document.activeElement && menuRef.value.contains(document.activeElement))

  if (filePos === null) {
    if (isInteracting) return
    lastFilePos.value = null
    applyState(false, 0, 0)
    return
  }

  lastFilePos.value = filePos

  if (menuRef.value) {
    const w = menuRef.value.offsetWidth
    if (w > 0) cachedMenuWidth = w
  }

  try {
    const { view } = props.editor
    const nodeDom = view.nodeDOM(filePos)
    const el = nodeDom instanceof HTMLElement
      ? nodeDom
      : (nodeDom?.parentElement ?? null)

    if (el) {
      const rect = el.getBoundingClientRect()
      const vp = getViewport()

      if (!isRectInView(rect, vp)) {
        applyState(false, 0, 0)
        return
      }

      let top = rect.top - 45
      if (top < vp.top + 8) top = rect.bottom + 8

      if (!isTopInView(top, vp)) {
        applyState(false, 0, 0)
        return
      }

      const centerLeft = rect.left + rect.width / 2
      const menuW = cachedMenuWidth
      if (menuW > 0) {
        const left = Math.round(Math.max(12, Math.min(window.innerWidth - menuW - 12, centerLeft - menuW / 2)))
        applyState(true, top, left)
      } else {
        // 首帧宽度未知：先显示在屏幕外以测量宽度，下一帧再修正位置
        applyState(true, top, -9999)
        window.requestAnimationFrame(computeMenu)
      }
    } else {
      applyState(false, 0, 0)
    }
  } catch {
    applyState(false, 0, 0)
  }
}

const { getViewport } = useMenuScheduler(props.editor, computeMenu)

const fileAttrs = computed(() => {
  const pos = getFilePos()
  if (pos === null) return null
  const { selection } = props.editor.state
  if (selection instanceof NodeSelection) return selection.node.attrs
  return null
})

const currentLayout = computed<FileLayout>(() => {
  return (fileAttrs.value?.layout as FileLayout) || 'full'
})

function runOnFile(command: (chain: any) => any) {
  const posSnapshot = lastFilePos.value
  setTimeout(() => {
    const chain = props.editor.chain().focus(null, { scrollIntoView: false })
    if (!(props.editor.state.selection instanceof NodeSelection) && posSnapshot !== null) {
      chain.setNodeSelection(posSnapshot)
    }
    command(chain).run()
  }, 0)
}

function setLayout(layout: FileLayout) {
  runOnFile((chain) => chain.updateAttributes('fileNode', { layout }))
}

function copyLink() {
  const href = fileAttrs.value?.href
  if (!href) return
  navigator.clipboard.writeText(href).then(() => {
    toast.success('链接已复制')
  }).catch(() => {
    toast.error('复制失败')
  })
}

function openInNewTab() {
  const href = fileAttrs.value?.href
  if (!href) return
  window.open(href, '_blank', 'noopener,noreferrer')
}

function deleteFile() {
  runOnFile((chain) => chain.deleteSelection())
}
</script>

<template>
  <div
    v-if="isVisible"
    ref="menuRef"
    class="glass fixed z-50 flex items-center gap-0.5 p-1.5 rounded-full antialiased"
    :style="{ top: `${position.top}px`, left: `${position.left}px` }"
  >
    <!-- 布局切换 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': currentLayout === 'full' }"
          @click="setLayout('full')"
        >
          <RectangleHorizontal class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>整行</TooltipContent>
    </Tooltip>

    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': currentLayout === 'half' }"
          @click="setLayout('half')"
        >
          <Columns2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>一行两个</TooltipContent>
    </Tooltip>

    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="{ 'bg-foreground/10': currentLayout === 'third' }"
          @click="setLayout('third')"
        >
          <Columns3 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>一行三个</TooltipContent>
    </Tooltip>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 复制链接 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @click="copyLink"
        >
          <Copy class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>复制链接</TooltipContent>
    </Tooltip>

    <!-- 新标签页打开 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          @click="openInNewTab"
        >
          <ExternalLink class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>新标签页打开</TooltipContent>
    </Tooltip>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 从文档移除（不删除媒体库文件，物理删除请在媒体管理页操作） -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full text-destructive hover:text-destructive hover:bg-destructive/10"
          @click="deleteFile"
        >
          <Trash2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>从文档移除</TooltipContent>
    </Tooltip>
  </div>
</template>
