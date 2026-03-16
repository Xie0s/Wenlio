/**
 * useOutsideClick.ts - 点击元素外部检测 composable
 * 职责：监听 mousedown/touchstart 事件，当点击发生在目标元素外部时触发回调
 * 对外暴露：useOutsideClick(elementRef, callback)
 */
import { onMounted, onUnmounted } from 'vue'
import type { Ref } from 'vue'

export function useOutsideClick(
  elementRef: Ref<HTMLElement | null>,
  callback: (event: MouseEvent | TouchEvent) => void,
) {
  const listener = (event: MouseEvent | TouchEvent) => {
    const el = elementRef.value
    if (!el || el.contains(event.target as Node)) return
    callback(event)
  }

  onMounted(() => {
    document.addEventListener('mousedown', listener)
    document.addEventListener('touchstart', listener)
  })

  onUnmounted(() => {
    document.removeEventListener('mousedown', listener)
    document.removeEventListener('touchstart', listener)
  })
}
