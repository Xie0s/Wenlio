/**
 * composables/useFloatingPanel.ts - 居中浮动面板（可拖拽）
 * 职责：为 ThemeEditorPage 等场景的浮窗提供初始居中定位和鼠标拖拽能力
 * 设计：首次打开时居中显示，拖拽后记住位置，关闭后重置回居中
 * 对外暴露：useFloatingPanel()
 */
import { ref, computed, onUnmounted } from 'vue'

export function useFloatingPanel() {
  const panelRef = ref<HTMLElement | null>(null)
  const left = ref(0)
  const top = ref(0)
  const positioned = ref(false)

  let startX = 0
  let startY = 0
  let startLeft = 0
  let startTop = 0

  const style = computed(() =>
    positioned.value
      ? { left: `${left.value}px`, top: `${top.value}px`, transform: 'none' }
      : { left: '50%', top: '50%', transform: 'translate(-50%, -50%)' },
  )

  function capturePosition() {
    const el = panelRef.value
    if (!el || positioned.value) return
    const rect = el.getBoundingClientRect()
    left.value = rect.left
    top.value = rect.top
    positioned.value = true
  }

  function onMouseDown(e: MouseEvent) {
    capturePosition()
    startX = e.clientX
    startY = e.clientY
    startLeft = left.value
    startTop = top.value
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    e.preventDefault()
  }

  function onMouseMove(e: MouseEvent) {
    left.value = startLeft + (e.clientX - startX)
    top.value = startTop + (e.clientY - startY)
  }

  function onMouseUp() {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  function reset() {
    positioned.value = false
  }

  onUnmounted(() => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  })

  return { panelRef, style, onMouseDown, reset }
}
