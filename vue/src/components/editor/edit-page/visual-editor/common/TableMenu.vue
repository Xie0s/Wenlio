<!--
  TableMenu.vue - 表格浮动菜单组件
  职责：当光标在表格内时显示浮动工具栏，提供行列操作、合并拆分、背景色、表格设置
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref } from 'vue'
import {
  Rows3, Columns3, Plus, Minus,
  Trash2, Merge, Split,
  TableProperties, Palette,
} from 'lucide-vue-next'
import { useMenuScheduler, isRectInView, isTopInView } from './useMenuScheduler'
import type { Editor } from '@tiptap/vue-3'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { cellBackgroundColors, canMergeCells, canSplitCell } from './table-helpers'

const props = defineProps<{ editor: Editor }>()

const isVisible = ref(false)
const position = ref({ top: 0, left: 0 })
const canMerge = ref(false)
const canSplit = ref(false)
const menuRef = ref<HTMLDivElement | null>(null)
const dropdownOpenCount = ref(0)
const lastTableAnchorPos = ref<number | null>(null)
const rowDropdownOpen = ref(false)
const colDropdownOpen = ref(false)
const colorDropdownOpen = ref(false)
const settingsDropdownOpen = ref(false)

function onDropdownOpenChange(open: boolean) {
  dropdownOpenCount.value = Math.max(0, dropdownOpenCount.value + (open ? 1 : -1))
}

const lastState = { isVisible: false, top: 0, left: 0, canMerge: false, canSplit: false }

// 缓存菜单实际渲染宽度，用于 JS 层居中计算（替代 CSS translateX(-50%)）
let cachedMenuWidth = 0

function applyState(nextVisible: boolean, nextTop: number, nextLeft: number, nextCanMerge: boolean, nextCanSplit: boolean) {
  const roundedTop = Math.round(nextTop)
  const roundedLeft = Math.round(nextLeft)

  if (lastState.isVisible !== nextVisible) {
    lastState.isVisible = nextVisible
    isVisible.value = nextVisible
  }
  if (lastState.canMerge !== nextCanMerge) {
    lastState.canMerge = nextCanMerge
    canMerge.value = nextCanMerge
  }
  if (lastState.canSplit !== nextCanSplit) {
    lastState.canSplit = nextCanSplit
    canSplit.value = nextCanSplit
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
  const editor = props.editor
  const isInteracting = !!(menuRef.value && document.activeElement && menuRef.value.contains(document.activeElement))

  if (!editor.isActive('table')) {
    // 下拉打开或菜单内交互时不隐藏（焦点可能临时离开编辑器）
    if (dropdownOpenCount.value > 0 || isInteracting) return
    lastTableAnchorPos.value = null
    applyState(false, 0, 0, false, false)
    return
  }

  // 下拉打开时冻结位置，避免重新计算导致位移
  if (dropdownOpenCount.value > 0) return

  // 更新菜单实际宽度缓存（canMerge/canSplit 改变时宽度可能变化）
  if (menuRef.value) {
    const w = menuRef.value.offsetWidth
    if (w > 0) cachedMenuWidth = w
  }

  const nextCanMerge = canMergeCells(editor)
  const nextCanSplit = canSplitCell(editor)

  try {
    const { view, state } = editor
    const anchorPos = state.selection.$anchor.pos
    lastTableAnchorPos.value = anchorPos
    const domAtPos = view.domAtPos(anchorPos)
    const currentNode = domAtPos.node
    const targetElement = currentNode.nodeType === Node.TEXT_NODE
      ? currentNode.parentElement
      : currentNode as HTMLElement
    const tableEl = targetElement?.closest('table')

    if (tableEl) {
      const tableRect = tableEl.getBoundingClientRect()

      const vp = getViewport()

      if (!isRectInView(tableRect, vp)) {
        applyState(false, 0, 0, nextCanMerge, nextCanSplit)
        return
      }

      let top = tableRect.top - 45
      if (top < vp.top + 8) top = tableRect.bottom + 8

      if (!isTopInView(top, vp)) {
        applyState(false, 0, 0, nextCanMerge, nextCanSplit)
        return
      }

      // 在 JS 层计算居中 left，避免 CSS translateX(-50%) 在宽度变化时产生半像素抖动
      const centerLeft = tableRect.left + tableRect.width / 2
      const menuW = cachedMenuWidth
      if (menuW > 0) {
        const left = Math.round(Math.max(12, Math.min(window.innerWidth - menuW - 12, centerLeft - menuW / 2)))
        applyState(true, top, left, nextCanMerge, nextCanSplit)
      } else {
        // 首次显示时 menuRef 尚未挂载，宽度不可测；先按中心点展示，再下一帧重算真实居中位置
        const left = Math.round(Math.min(window.innerWidth - 12, Math.max(12, centerLeft)))
        applyState(true, top, left, nextCanMerge, nextCanSplit)
        window.requestAnimationFrame(computeMenu)
      }
    } else {
      applyState(false, 0, 0, nextCanMerge, nextCanSplit)
    }
  } catch {
    applyState(false, 0, 0, nextCanMerge, nextCanSplit)
  }
}

const { getViewport } = useMenuScheduler(props.editor, computeMenu)

// ── 行操作 ──
function runInTable(command: (chain: any) => any) {
  const anchorPosSnapshot = lastTableAnchorPos.value
  setTimeout(() => {
    const chain = props.editor.chain().focus(null, { scrollIntoView: false })
    if (!props.editor.isActive('table') && anchorPosSnapshot !== null) {
      chain.setTextSelection(anchorPosSnapshot)
    }
    command(chain).run()
  }, 0)
}

function addRowBefore() {
  runInTable((chain) => chain.addRowBefore())
}
function addRowAfter() {
  runInTable((chain) => chain.addRowAfter())
}
function deleteRow() {
  runInTable((chain) => chain.deleteRow())
}

// ── 列操作 ──
function addColumnBefore() {
  runInTable((chain) => chain.addColumnBefore())
}
function addColumnAfter() {
  runInTable((chain) => chain.addColumnAfter())
}
function deleteColumn() {
  runInTable((chain) => chain.deleteColumn())
}

// ── 表格操作 ──
function deleteTable() {
  runInTable((chain) => chain.deleteTable())
}
function toggleHeaderRow() {
  runInTable((chain) => chain.toggleHeaderRow())
}
function toggleHeaderColumn() {
  runInTable((chain) => chain.toggleHeaderColumn())
}

// ── 单元格操作 ──
function handleMergeCells() {
  runInTable((chain) => chain.mergeCells())
}
function handleSplitCell() {
  runInTable((chain) => chain.splitCell())
}

// 设置单元格背景色
function setCellBackground(color: string | null) {
  runInTable((chain) => chain.setCellAttribute('backgroundColor', color))
}


</script>

<template>
  <div
    v-if="isVisible"
    ref="menuRef"
    class="glass fixed z-50 flex items-center gap-0.5 p-1.5 rounded-full antialiased"
    :style="{ top: `${position.top}px`, left: `${position.left}px` }"
  >
    <!-- 行操作下拉 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { rowDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': rowDropdownOpen }">
          <Rows3 class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem @click="addRowBefore">
          <Plus class="h-4 w-4 mr-2" />在上方插入行
        </DropdownMenuItem>
        <DropdownMenuItem @click="addRowAfter">
          <Plus class="h-4 w-4 mr-2" />在下方插入行
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem class="text-destructive" @click="deleteRow">
          <Minus class="h-4 w-4 mr-2" />删除当前行
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <!-- 列操作下拉 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { colDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': colDropdownOpen }">
          <Columns3 class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem @click="addColumnBefore">
          <Plus class="h-4 w-4 mr-2" />在左侧插入列
        </DropdownMenuItem>
        <DropdownMenuItem @click="addColumnAfter">
          <Plus class="h-4 w-4 mr-2" />在右侧插入列
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem class="text-destructive" @click="deleteColumn">
          <Minus class="h-4 w-4 mr-2" />删除当前列
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 合并/拆分单元格 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!canMerge" @click="handleMergeCells">
          <Merge class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>合并单元格</TooltipContent>
    </Tooltip>

    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :disabled="!canSplit" @click="handleSplitCell">
          <Split class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>拆分单元格</TooltipContent>
    </Tooltip>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 单元格背景色 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { colorDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': colorDropdownOpen }">
          <Palette class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" class="w-40">
        <DropdownMenuItem
          v-for="color in cellBackgroundColors"
          :key="color.name"
          class="flex items-center gap-2"
          @click="setCellBackground(color.value)"
        >
          <div
            class="w-4 h-4 rounded border border-border"
            :class="!color.value && 'bg-transparent'"
            :style="color.value ? { backgroundColor: color.value } : undefined"
          />
          {{ color.name }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 表格设置 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { settingsDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': settingsDropdownOpen }">
          <TableProperties class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem @click="toggleHeaderRow">切换表头行</DropdownMenuItem>
        <DropdownMenuItem @click="toggleHeaderColumn">切换表头列</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem class="text-destructive" @click="deleteTable">
          <Trash2 class="h-4 w-4 mr-2" />删除表格
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
