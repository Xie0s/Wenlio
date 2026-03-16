/**
 * 演示文件专用，与其他不兼容，其他组件不要使用
 * composables/useDraggablePanel.ts - 可拖拽浮动面板定位
 * 职责：为浮动面板提供拖拽移动能力和初始定位
 * 对外暴露：useDraggablePanel(options?)
 */
import { ref, computed, onUnmounted } from 'vue'

interface DraggablePanelOptions {
  initialRight?: number
  initialTop?: number
}

export function useDraggablePanel(options?: DraggablePanelOptions) {
  const right = ref(options?.initialRight ?? 16)
  const top = ref(options?.initialTop ?? 64)
  const dragging = ref(false)

  let startX = 0
  let startY = 0
  let startRight = 0
  let startTop = 0

  const positionStyle = computed(() => ({
    right: `${right.value}px`,
    top: `${top.value}px`,
  }))

  function startDrag(e: MouseEvent) {
    dragging.value = true
    startX = e.clientX
    startY = e.clientY
    startRight = right.value
    startTop = top.value
    document.addEventListener('mousemove', onDrag)
    document.addEventListener('mouseup', stopDrag)
    e.preventDefault()
  }

  function onDrag(e: MouseEvent) {
    const dx = e.clientX - startX
    const dy = e.clientY - startY
    right.value = Math.max(0, startRight - dx)
    top.value = Math.max(0, startTop + dy)
  }

  function stopDrag() {
    dragging.value = false
    document.removeEventListener('mousemove', onDrag)
    document.removeEventListener('mouseup', stopDrag)
  }

  function resetPosition() {
    right.value = options?.initialRight ?? 16
    top.value = options?.initialTop ?? 64
  }

  onUnmounted(() => {
    document.removeEventListener('mousemove', onDrag)
    document.removeEventListener('mouseup', stopDrag)
  })

  return { positionStyle, startDrag, resetPosition, dragging }
}
