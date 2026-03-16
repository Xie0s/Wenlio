<!--
  LinkMenu.vue - 链接浮动菜单组件
  职责：当光标在链接内时显示浮动工具栏，提供编辑、打开、复制、取消链接操作
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref, watch } from 'vue'
import { Pencil, ExternalLink, Unlink, Copy } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import { useMenuScheduler, isRectInView, isTopInView } from './useMenuScheduler'
import type { Editor } from '@tiptap/vue-3'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import LinkInputPopover from './LinkInputPopover.vue'

const props = defineProps<{ editor: Editor }>()

const isVisible = ref(false)
const position = ref({ top: 0, left: 0 })
const currentUrl = ref('')
const currentText = ref('')
const showEditPopover = ref(false)
const menuRef = ref<HTMLDivElement | null>(null)

// 缓存上一次状态，避免不必要的响应式更新
const lastState = {
  isVisible: false,
  top: 0,
  left: 0,
  url: '',
  text: '',
}

function applyState(nextVisible: boolean, nextTop: number, nextLeft: number, nextUrl: string, nextText: string) {
  const roundedTop = Math.round(nextTop)
  const roundedLeft = Math.round(nextLeft)

  if (lastState.url !== nextUrl) {
    lastState.url = nextUrl
    currentUrl.value = nextUrl
  }
  if (lastState.text !== nextText) {
    lastState.text = nextText
    currentText.value = nextText
  }
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

function computeMenu() {
  // 编辑对话框打开时保持菜单不动
  if (showEditPopover.value) return

  const editor = props.editor
  if (!editor.isActive('link')) {
    applyState(false, 0, 0, '', '')
    return
  }

  const attrs = editor.getAttributes('link')
  const nextUrl = attrs.href || ''
  const { from, to } = editor.state.selection
  const nextText = editor.state.doc.textBetween(from, to) || ''

  try {
    const { view } = editor
    const { $from } = editor.state.selection
    const pos = $from.pos
    const domAtPos = view.domAtPos(pos)
    let linkEl: HTMLAnchorElement | null = null
    if (domAtPos.node) {
      const node = domAtPos.node.nodeType === Node.TEXT_NODE
        ? domAtPos.node.parentElement
        : domAtPos.node as HTMLElement
      linkEl = node?.closest('a') as HTMLAnchorElement | null
    }

    const menuWidth = menuRef.value?.offsetWidth || 160
    const maxLeft = Math.max(8, window.innerWidth - menuWidth - 8)

    const vp = getViewport()

    if (linkEl) {
      const linkRect = linkEl.getBoundingClientRect()
      if (!isRectInView(linkRect, vp)) {
        applyState(false, 0, 0, nextUrl, nextText)
        return
      }
      const top = linkRect.bottom + 8
      if (!isTopInView(top, vp)) {
        applyState(false, 0, 0, nextUrl, nextText)
        return
      }
      const left = Math.min(maxLeft, Math.max(10, linkRect.left + (linkRect.width - menuWidth) / 2))
      applyState(true, top, left, nextUrl, nextText)
    } else {
      const coords = view.coordsAtPos(pos)
      const coordRect = new DOMRect(coords.left, coords.top, 1, coords.bottom - coords.top)
      if (!isRectInView(coordRect, vp)) {
        applyState(false, 0, 0, nextUrl, nextText)
        return
      }
      const top = coords.bottom + 8
      if (!isTopInView(top, vp)) {
        applyState(false, 0, 0, nextUrl, nextText)
        return
      }
      const left = Math.min(maxLeft, Math.max(10, coords.left - menuWidth / 2))
      applyState(true, top, left, nextUrl, nextText)
    }
  } catch {
    applyState(false, 0, 0, nextUrl, nextText)
  }
}

const { scheduleUpdate, getViewport } = useMenuScheduler(props.editor, computeMenu)

// 编辑对话框状态变化时重新计算
watch(showEditPopover, () => {
  if (!showEditPopover.value) scheduleUpdate()
})

// ── 操作 ──
function handleEditConfirm(url: string) {
  showEditPopover.value = false
  setTimeout(() => {
    props.editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
  }, 0)
}

function handleRemove() {
  showEditPopover.value = false
  setTimeout(() => {
    props.editor.chain().focus().unsetLink().run()
  }, 0)
}

function handleOpen() {
  if (currentUrl.value) {
    window.open(currentUrl.value, '_blank', 'noopener,noreferrer')
  }
}

function handleUnlink() {
  props.editor.chain().focus().unsetLink().run()
}

async function copyToClipboard(text: string): Promise<boolean> {
  if (window.isSecureContext && navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return true
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'true')
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  textarea.style.pointerEvents = 'none'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()

  const copied = document.execCommand('copy')
  document.body.removeChild(textarea)
  return copied
}

async function handleCopy() {
  const url = currentUrl.value.trim()
  if (!url) return

  try {
    const copied = await copyToClipboard(url)
    if (copied) {
      toast.success('链接已复制')
      return
    }
    toast.error('复制失败，请手动复制链接')
  } catch {
    toast.error('复制失败，请检查剪贴板权限')
  }
}

</script>

<template>
  <div
    v-if="isVisible"
    ref="menuRef"
    class="glass fixed z-50 flex items-center gap-0.5 p-1.5 rounded-full antialiased"
    :style="{ top: `${position.top}px`, left: `${position.left}px` }"
  >
    <!-- 编辑链接 -->
    <LinkInputPopover
      v-model:open="showEditPopover"
      :initial-url="currentUrl"
      :is-edit-mode="true"
      @confirm="handleEditConfirm"
      @remove="handleRemove"
    >
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': showEditPopover }" @click="showEditPopover = !showEditPopover">
            <Pencil class="h-4 w-4" />
          </Button>
        </TooltipTrigger>
        <TooltipContent v-if="!showEditPopover">编辑链接</TooltipContent>
      </Tooltip>
    </LinkInputPopover>

    <!-- 打开链接 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!currentUrl" @click="handleOpen">
          <ExternalLink class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>打开链接</TooltipContent>
    </Tooltip>

    <!-- 复制链接 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!currentUrl" @click="handleCopy">
          <Copy class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>复制链接</TooltipContent>
    </Tooltip>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 取消链接 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full text-destructive hover:text-destructive hover:bg-foreground/10" @click="handleUnlink">
          <Unlink class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>取消链接</TooltipContent>
    </Tooltip>
  </div>

</template>
