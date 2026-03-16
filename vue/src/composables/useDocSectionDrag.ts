/**
 * composables/useDocSectionDrag.ts - 章节拖拽排序与展开/折叠状态管理
 * 职责：封装章节列表的展开/折叠快照恢复、拖拽排序全流程（DnD 事件处理、插入线计算、排序提交）
 * 对外暴露：useDocSectionDrag(editor) → { dragIndex, isSectionOpen, toggleSection, showInsertionLineAt, onDragStart, onDragOver, onDragEnd, onDrop }
 */

import { ref } from 'vue'
import { toast } from 'vue-sonner'
import type { ThemeEditorContext } from '@/composables/useThemeEditor'

export function useDocSectionDrag(editor: ThemeEditorContext) {
  const expandedSections = ref<Record<string, boolean>>({})
  const expandedSectionsSnapshot = ref<Record<string, boolean> | null>(null)

  function isSectionOpen(sectionId: string) {
    return expandedSections.value[sectionId] !== false
  }

  function toggleSection(sectionId: string) {
    expandedSections.value[sectionId] = !isSectionOpen(sectionId)
  }

  const dragIndex = ref<number | null>(null)
  const insertionIndex = ref<number | null>(null)

  function snapshotExpandedSections() {
    const snapshot: Record<string, boolean> = {}
    for (const section of editor.sections.value) {
      snapshot[section.id] = isSectionOpen(section.id)
    }
    expandedSectionsSnapshot.value = snapshot
  }

  function collapseAllSections() {
    const next: Record<string, boolean> = {}
    for (const section of editor.sections.value) {
      next[section.id] = false
    }
    expandedSections.value = next
  }

  function restoreExpandedSections() {
    if (!expandedSectionsSnapshot.value) return
    expandedSections.value = expandedSectionsSnapshot.value
    expandedSectionsSnapshot.value = null
  }

  function resolveInsertionIndex(index: number, event: DragEvent) {
    const currentTarget = event.currentTarget
    if (!(currentTarget instanceof HTMLElement)) return index
    const rect = currentTarget.getBoundingClientRect()
    return event.clientY < rect.top + rect.height / 2 ? index : index + 1
  }

  function normalizeInsertionIndex(rawIndex: number, fromIndex: number) {
    return rawIndex > fromIndex ? rawIndex - 1 : rawIndex
  }

  function hasEffectiveDrop(rawIndex: number | null) {
    if (dragIndex.value === null || rawIndex === null) return false
    return normalizeInsertionIndex(rawIndex, dragIndex.value) !== dragIndex.value
  }

  function showInsertionLineAt(index: number) {
    return insertionIndex.value === index && hasEffectiveDrop(index)
  }

  function onDragStart(index: number, e: DragEvent) {
    dragIndex.value = index
    insertionIndex.value = null
    snapshotExpandedSections()
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = 'move'
      e.dataTransfer.setData('text/plain', String(index))
    }
    requestAnimationFrame(() => {
      collapseAllSections()
    })
  }

  function onDragOver(index: number, e: DragEvent) {
    e.preventDefault()
    if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
    const nextInsertionIndex = resolveInsertionIndex(index, e)
    if (insertionIndex.value !== nextInsertionIndex) {
      insertionIndex.value = nextInsertionIndex
    }
  }

  function onDragEnd() {
    if (dragIndex.value === null && insertionIndex.value === null && expandedSectionsSnapshot.value === null) return
    restoreExpandedSections()
    dragIndex.value = null
    insertionIndex.value = null
  }

  async function onDrop(index: number, event: DragEvent) {
    event.preventDefault()
    const fromIndex = dragIndex.value
    if (fromIndex === null) {
      onDragEnd()
      return
    }

    const rawInsertionIndex = resolveInsertionIndex(index, event)
    const targetIndex = normalizeInsertionIndex(rawInsertionIndex, fromIndex)
    if (targetIndex === fromIndex) {
      onDragEnd()
      return
    }

    const list = [...editor.sections.value]
    const [moved] = list.splice(fromIndex, 1)
    if (!moved) { onDragEnd(); return }
    list.splice(targetIndex, 0, moved)
    editor.sections.value = list
    onDragEnd()
    await editor.sortSections(list.map(s => s.id))
    toast.success('章节排序已更新')
  }

  return {
    dragIndex,
    isSectionOpen,
    toggleSection,
    showInsertionLineAt,
    onDragStart,
    onDragOver,
    onDragEnd,
    onDrop,
  }
}
