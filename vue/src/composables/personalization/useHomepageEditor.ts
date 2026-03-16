/**
 * composables/personalization/useHomepageEditor.ts - 首页编辑器状态管理
 *
 * 职责：管理编辑器中的布局状态，提供区块增删、排序、配置修改等操作
 * 对外暴露：useHomepageEditor() composable 函数
 */

import { ref, computed } from 'vue'
import type {
  HomepageLayout,
  HomepageSection,
  SectionType,
} from '@/components/personalization/types'
import { SECTION_TYPE_META } from '@/components/personalization/types'
import {
  defaultHomepageLayout,
  createSection,
} from '@/components/personalization/defaults'

export function useHomepageEditor(initialLayout?: HomepageLayout | null) {
  const layout = ref<HomepageLayout>(initialLayout ?? defaultHomepageLayout())
  const selectedSectionId = ref<string | null>(null)
  const dirty = ref(false)

  const sections = computed(() => layout.value.sections)

  const selectedSection = computed(() =>
    layout.value.sections.find(s => s.id === selectedSectionId.value) ?? null,
  )

  /** 检查某类型是否可继续添加（singleton 类型只允许一个） */
  function canAddSection(type: SectionType): boolean {
    const meta = SECTION_TYPE_META[type]
    if (!meta.singleton) return true
    return !layout.value.sections.some(s => s.type === type)
  }

  /** 添加区块 */
  function addSection(type: SectionType, index?: number): HomepageSection | null {
    if (!canAddSection(type)) return null
    const section = createSection(type)
    const insertAt = index ?? layout.value.sections.length
    layout.value.sections.splice(insertAt, 0, section)
    selectedSectionId.value = section.id
    dirty.value = true
    return section
  }

  /** 删除区块 */
  function removeSection(id: string) {
    const idx = layout.value.sections.findIndex(s => s.id === id)
    if (idx === -1) return
    layout.value.sections.splice(idx, 1)
    if (selectedSectionId.value === id) {
      selectedSectionId.value = null
    }
    dirty.value = true
  }

  /** 切换区块可见性 */
  function toggleVisibility(id: string) {
    const section = layout.value.sections.find(s => s.id === id)
    if (section) {
      section.visible = !section.visible
      dirty.value = true
    }
  }

  /** 上移区块 */
  function moveUp(id: string) {
    const idx = layout.value.sections.findIndex(s => s.id === id)
    if (idx <= 0) return
    const arr = layout.value.sections
    const moved = arr.splice(idx, 1)
    arr.splice(idx - 1, 0, ...moved)
    dirty.value = true
  }

  /** 下移区块 */
  function moveDown(id: string) {
    const idx = layout.value.sections.findIndex(s => s.id === id)
    const arr = layout.value.sections
    if (idx === -1 || idx >= arr.length - 1) return
    const moved = arr.splice(idx, 1)
    arr.splice(idx + 1, 0, ...moved)
    dirty.value = true
  }

  /** 更新全局样式配置 */
  function updateGlobal(global: HomepageLayout['global']) {
    layout.value.global = global
    dirty.value = true
  }

  /** 更新区块配置 */
  function updateSectionConfig(id: string, config: HomepageSection['config']) {
    const section = layout.value.sections.find(s => s.id === id)
    if (section) {
      section.config = config
      dirty.value = true
    }
  }

  /** 重置为默认布局 */
  function resetToDefault() {
    layout.value = defaultHomepageLayout()
    selectedSectionId.value = null
    dirty.value = true
  }

  /** 从外部数据初始化（如从 API 加载后） */
  function initFromLayout(newLayout: HomepageLayout) {
    layout.value = JSON.parse(JSON.stringify(newLayout))
    selectedSectionId.value = null
    dirty.value = false
  }

  /** 标记为已保存（dirty 清零） */
  function markSaved() {
    dirty.value = false
  }

  return {
    layout,
    sections,
    selectedSectionId,
    selectedSection,
    dirty,
    canAddSection,
    addSection,
    removeSection,
    toggleVisibility,
    moveUp,
    moveDown,
    updateGlobal,
    updateSectionConfig,
    resetToDefault,
    initFromLayout,
    markSaved,
  }
}
