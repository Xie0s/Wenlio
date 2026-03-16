/**
 * composables/useDocTreeCrud.ts - 章节/文档树 CRUD UI 状态管理
 * 职责：提取 EditorSidebar 和 EditorSettingsPanel 中重复的内联编辑、创建、删除确认逻辑
 * 对外暴露：useDocTreeCrud(editor) → 章节编辑/创建/删除 + 文档页编辑/创建/删除 的完整状态和操作
 */

import { ref, computed, watch } from 'vue'
import type { Section, DocPage } from '@/utils/types'
import { validateSlug } from '@/lib/validation'
import { titleToSlug } from '@/lib/slug'
import type { ThemeEditorContext } from '@/composables/useThemeEditor'

export type DocTreeCrudContext = ReturnType<typeof useDocTreeCrud>

export function useDocTreeCrud(editor: ThemeEditorContext) {
  // ═══════════ 章节内联编辑 ═══════════
  const editingSectionId = ref('')
  const editSectionTitle = ref('')

  function startEditSection(s: Section) {
    editingSectionId.value = s.id
    editSectionTitle.value = s.title
  }

  function confirmEditSection() {
    if (editSectionTitle.value.trim()) {
      editor.updateSection(editingSectionId.value, editSectionTitle.value.trim())
    }
    editingSectionId.value = ''
  }

  // ═══════════ 文档页内联编辑 ═══════════
  const editingPageId = ref('')
  const editPageForm = ref({ title: '', slug: '' })
  const editPageSlugError = computed(() =>
    editPageForm.value.slug ? validateSlug(editPageForm.value.slug) : '',
  )

  function startEditPage(p: DocPage) {
    editingPageId.value = p.id
    editPageForm.value = { title: p.title, slug: p.slug }
  }

  function confirmEditPage() {
    if (editPageForm.value.title.trim()) {
      editor.updatePage(editingPageId.value, { ...editPageForm.value })
    }
    editingPageId.value = ''
  }

  // ═══════════ 创建章节 ═══════════
  const showCreateSection = ref(false)
  const sectionTitle = ref('')

  function toggleCreateSection() {
    showCreateSection.value = !showCreateSection.value
    if (!showCreateSection.value) sectionTitle.value = ''
  }

  async function submitCreateSection() {
    await editor.createSection(sectionTitle.value)
    showCreateSection.value = false
    sectionTitle.value = ''
  }

  // ═══════════ 删除确认 Popover ═══════════
  const deleteSectionPopoverId = ref('')
  const deletePagePopoverId = ref('')

  function setSectionDeletePopover(id: string, open: boolean) {
    deleteSectionPopoverId.value = open ? id : ''
  }

  function setPageDeletePopover(id: string, open: boolean) {
    deletePagePopoverId.value = open ? id : ''
  }

  async function confirmDeleteSection(id: string) {
    await editor.deleteSection(id)
    if (deleteSectionPopoverId.value === id) deleteSectionPopoverId.value = ''
  }

  async function confirmDeletePage(id: string) {
    await editor.deletePage(id)
    if (deletePagePopoverId.value === id) deletePagePopoverId.value = ''
  }

  // ═══════════ 创建文档页 Popover ═══════════
  const createPagePopoverSectionId = ref('')
  const createPageSectionId = ref('')
  const pageForm = ref({ title: '', slug: '', content: '' })
  const createPageSlugError = computed(() =>
    pageForm.value.slug ? validateSlug(pageForm.value.slug) : '',
  )
  const createPageSlugAuto = ref(true)

  watch(() => pageForm.value.title, () => {
    if (createPageSlugAuto.value) {
      pageForm.value.slug = titleToSlug(pageForm.value.title)
    }
  })

  watch(() => pageForm.value.slug, (slug) => {
    const expectedSlug = titleToSlug(pageForm.value.title)
    createPageSlugAuto.value = !slug || slug === expectedSlug
  })

  function setCreatePagePopover(sectionId: string, open: boolean) {
    if (open) {
      createPageSectionId.value = sectionId
      createPagePopoverSectionId.value = sectionId
      createPageSlugAuto.value = true
      pageForm.value = { title: '', slug: '', content: '' }
      return
    }

    if (createPagePopoverSectionId.value === sectionId) {
      createPagePopoverSectionId.value = ''
      createPageSlugAuto.value = true
      pageForm.value = { title: '', slug: '', content: '' }
    }
  }

  async function submitCreatePage() {
    await editor.createPage(createPageSectionId.value, pageForm.value)
    createPagePopoverSectionId.value = ''
    createPageSlugAuto.value = true
    pageForm.value = { title: '', slug: '', content: '' }
  }

  return {
    // 章节编辑
    editingSectionId,
    editSectionTitle,
    startEditSection,
    confirmEditSection,
    // 文档页编辑
    editingPageId,
    editPageForm,
    editPageSlugError,
    startEditPage,
    confirmEditPage,
    // 创建章节
    showCreateSection,
    sectionTitle,
    toggleCreateSection,
    submitCreateSection,
    // 删除确认
    deleteSectionPopoverId,
    deletePagePopoverId,
    setSectionDeletePopover,
    setPageDeletePopover,
    confirmDeleteSection,
    confirmDeletePage,
    // 创建文档页
    createPagePopoverSectionId,
    pageForm,
    createPageSlugError,
    setCreatePagePopover,
    submitCreatePage,
  }
}
