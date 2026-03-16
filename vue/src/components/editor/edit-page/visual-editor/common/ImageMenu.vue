<!--
  ImageMenu.vue - 图片浮动菜单组件
  职责：当选中图片时显示浮动工具栏，提供对齐、尺寸调整、删除等操作
  对外暴露：Props: editor (Tiptap Editor 实例)
-->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { NodeSelection } from '@tiptap/pm/state'
import {
  AlignLeft, AlignCenter, AlignRight,
  Trash2,
  Image as ImageIcon, RectangleHorizontal,
} from 'lucide-vue-next'
import { useMenuScheduler, isRectInView, isTopInView } from './useMenuScheduler'
import type { Editor } from '@tiptap/vue-3'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
const props = defineProps<{ editor: Editor }>()

// 尺寸预设
const SIZE_PRESETS = [
  { value: 'small', label: '小 (300px)', width: '300px' },
  { value: 'medium', label: '中 (600px)', width: '600px' },
  { value: 'large', label: '大 (900px)', width: '900px' },
  { value: 'full', label: '全宽 (100%)', width: '100%' },
]

// 圆角预设
const ROUNDED_PRESETS = [
  { value: 'none', label: '无圆角' },
  { value: 'xl', label: '小圆角 (xl)' },
  { value: '2xl', label: '中圆角 (2xl)' },
  { value: '3xl', label: '大圆角 (3xl)' },
]

const isVisible = ref(false)
const position = ref({ top: 0, left: 0 })
const menuRef = ref<HTMLDivElement | null>(null)
const dropdownOpenCount = ref(0)
const lastImagePos = ref<number | null>(null)
const sizeDropdownOpen = ref(false)
const roundedDropdownOpen = ref(false)

function onDropdownOpenChange(open: boolean) {
  dropdownOpenCount.value = Math.max(0, dropdownOpenCount.value + (open ? 1 : -1))
}

const lastState = { isVisible: false, top: 0, left: 0 }

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

// 获取当前选中的图片节点位置
function getImagePos(): number | null {
  const { selection } = props.editor.state
  if (selection instanceof NodeSelection && selection.node.type.name === 'image') {
    return selection.from
  }
  return null
}

function computeMenu() {
  const imagePos = getImagePos()
  const isInteracting = !!(menuRef.value && document.activeElement && menuRef.value.contains(document.activeElement))
  if (imagePos === null) {
    if (dropdownOpenCount.value > 0 || isInteracting) return
    lastImagePos.value = null
    applyState(false, 0, 0)
    return
  }

  lastImagePos.value = imagePos
  if (dropdownOpenCount.value > 0) return

  try {
    const { view } = props.editor
    const imageWrapperDom = view.nodeDOM(imagePos)
    const imageEl = imageWrapperDom instanceof HTMLElement
      ? imageWrapperDom
      : (imageWrapperDom?.parentElement ?? null)

    const imgEl = imageEl?.tagName === 'IMG'
      ? imageEl
      : imageEl?.querySelector('img') as HTMLElement | null

    if (imgEl) {
      const imgRect = imgEl.getBoundingClientRect()

      if (imgRect.width < 2 || imgRect.height < 2) {
        if (lastState.isVisible) {
          window.requestAnimationFrame(computeMenu)
          return
        }
        applyState(false, 0, 0)
        return
      }

      const vp = getViewport()

      if (!isRectInView(imgRect, vp)) {
        applyState(false, 0, 0)
        return
      }

      let top = imgRect.top - 45
      if (top < vp.top + 8) top = imgRect.bottom + 8

      if (!isTopInView(top, vp)) {
        applyState(false, 0, 0)
        return
      }

      // 使用中心点定位，避免 menu 宽度变化造成点击时位移
      const centerLeft = imgRect.left + imgRect.width / 2
      const left = Math.min(window.innerWidth - 12, Math.max(12, centerLeft))
      applyState(true, top, left)
    } else {
      applyState(false, 0, 0)
    }
  } catch {
    applyState(false, 0, 0)
  }
}

const { getViewport } = useMenuScheduler(props.editor, computeMenu)

// 获取当前图片属性
const imageAttrs = computed(() => {
  const pos = getImagePos()
  if (pos === null) return null
  const { selection } = props.editor.state
  if (selection instanceof NodeSelection) return selection.node.attrs
  return null
})

const currentAlign = computed(() => imageAttrs.value?.['data-align'] || 'left')
const currentRounded = computed(() => imageAttrs.value?.['data-rounded'] || 'xl')

// ── 操作 ──
// pos 在 setTimeout 外部同步捕获（闭包保存），防止 lastImagePos 在延迟执行前被清空
// 直接操作 tr 而非 updateAttributes，避免 updateAttributes 依赖 state.selection 快照
// 同一 tr 内完成属性修改 + NodeSelection 恢复，防止 setNodeMarkup 副作用导致菜单消失
function runOnImage(attrs: Record<string, unknown>) {
  const posSnapshot = lastImagePos.value
  if (posSnapshot === null) return
  setTimeout(() => {
    const node = props.editor.state.doc.nodeAt(posSnapshot)
    if (!node || node.type.name !== 'image') return
    const chain = props.editor.chain()
    if (!props.editor.isFocused) {
      chain.focus(null, { scrollIntoView: false })
    }
    chain
      .command(({ tr, dispatch }) => {
        if (!dispatch) return true
        tr.setNodeMarkup(posSnapshot, undefined, { ...node.attrs, ...attrs })
        tr.setSelection(NodeSelection.create(tr.doc, posSnapshot))
        return true
      })
      .run()
  }, 0)
}

function setSize(width: string) {
  runOnImage({ width })
}

function setAlign(align: string) {
  runOnImage({ 'data-align': align })
}

function setRounded(rounded: string) {
  // 'xl' 是默认值，设置为 xl 时移除属性以保持 markdown 简洁
  const value = rounded === 'xl' ? null : rounded
  runOnImage({ 'data-rounded': value })
}

function deleteImage() {
  const posSnapshot = lastImagePos.value
  if (posSnapshot === null) return
  setTimeout(() => {
    const node = props.editor.state.doc.nodeAt(posSnapshot)
    if (!node) return
    const chain = props.editor.chain()
    if (!props.editor.isFocused) {
      chain.focus(null, { scrollIntoView: false })
    }
    chain
      .command(({ tr, dispatch }) => {
        if (!dispatch) return true
        tr.delete(posSnapshot, posSnapshot + node.nodeSize)
        return true
      })
      .run()
  }, 0)
}

</script>

<template>
  <div
    v-if="isVisible"
    ref="menuRef"
    class="glass fixed z-50 flex items-center gap-0.5 p-1.5 rounded-full antialiased"
    :style="{ top: `${position.top}px`, left: `${position.left}px`, transform: 'translateX(-50%)' }"
  >
    <!-- 尺寸选择 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { sizeDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': sizeDropdownOpen }" @mousedown.prevent>
          <ImageIcon class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem
          v-for="preset in SIZE_PRESETS"
          :key="preset.value"
          @click="setSize(preset.width)"
        >
          {{ preset.label }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 对齐方式 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="currentAlign === 'left' && 'bg-foreground/10'"
          @mousedown.prevent
          @click="setAlign('left')"
        >
          <AlignLeft class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>左对齐</TooltipContent>
    </Tooltip>

    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="currentAlign === 'center' && 'bg-foreground/10'"
          @mousedown.prevent
          @click="setAlign('center')"
        >
          <AlignCenter class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>居中</TooltipContent>
    </Tooltip>

    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full hover:bg-foreground/10"
          :class="currentAlign === 'right' && 'bg-foreground/10'"
          @mousedown.prevent
          @click="setAlign('right')"
        >
          <AlignRight class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>右对齐</TooltipContent>
    </Tooltip>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 圆角选择 -->
    <DropdownMenu :modal="false" @update:open="(o: boolean) => { roundedDropdownOpen = o; onDropdownOpenChange(o) }">
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full hover:bg-foreground/10" :class="{ 'bg-foreground/10': roundedDropdownOpen }" @mousedown.prevent>
          <RectangleHorizontal class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem
          v-for="preset in ROUNDED_PRESETS"
          :key="preset.value"
          :class="currentRounded === preset.value && 'bg-accent'"
          @click="setRounded(preset.value)"
        >
          {{ preset.label }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <div class="w-px h-5 bg-border mx-1" />

    <!-- 删除 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost" size="icon"
          class="h-8 w-8 rounded-full text-destructive hover:text-destructive hover:bg-destructive/10"
          @mousedown.prevent
          @click="deleteImage"
        >
          <Trash2 class="h-4 w-4" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>从文档移除</TooltipContent>
    </Tooltip>
  </div>
</template>
