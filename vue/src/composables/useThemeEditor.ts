/**
 * composables/useThemeEditor.ts - 主题编辑器核心 composable
 * 职责：管理编辑器全部业务状态与 API 操作，通过 provide/inject 供子组件消费
 * 对外暴露：useThemeEditor(themeId)、useEditorInject()、EDITOR_KEY
 */

import { ref, watch, nextTick, type InjectionKey, inject } from 'vue'
import { http } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'
import { useToc } from '@/composables/useToc'
import { useAutoSave } from '@/composables/useAutoSave'
import { useMarkdownPreview } from '@/composables/useMarkdownPreview'
import { useMedia } from '@/lib/media'
import { toast } from 'vue-sonner'
import type { Theme, Version, Section, DocPage, CreateVersionReq } from '@/utils/types'
import type { SectionWithPages } from '@/components/editor/edit-page/types'

export type EditorMode = 'source' | 'visual'

// ── provide/inject key ──
export const EDITOR_KEY: InjectionKey<ThemeEditorContext> = Symbol('theme-editor')

export function useEditorInject() {
  const ctx = inject(EDITOR_KEY)
  if (!ctx) throw new Error('useEditorInject 必须在 ThemeEditorPage 内使用')
  return ctx
}

export type ThemeEditorContext = ReturnType<typeof useThemeEditor>

// ── 主 composable ──
export function useThemeEditor(themeId: string) {
  const authStore = useAuthStore()

  // ═══════════ 主题 ═══════════
  const theme = ref<Theme | null>(null)

  async function loadTheme() {
    const res = await http.get<Theme>(`/tenant/themes/${themeId}`)
    if (res.code === 0 && res.data) theme.value = res.data
  }

  async function updateTheme(form: { name: string; slug: string; description: string; access_mode?: string; access_code?: string }) {
    const res = await http.patch(`/tenant/themes/${themeId}`, form)
    res.code === 0 ? (toast.success('主题已更新'), loadTheme()) : toast.error(res.message)
  }

  // ═══════════ 版本 ═══════════
  const versions = ref<Version[]>([])
  const activeVersionId = ref('')

  async function loadVersions() {
    const res = await http.get<Version[]>(`/tenant/themes/${themeId}/versions`)
    if (res.code === 0 && res.data) {
      versions.value = res.data
      if (!activeVersionId.value && res.data.length > 0) {
        activeVersionId.value = (res.data.find(v => v.is_default) || res.data[0])!.id
      }
    }
  }

  function switchVersion(vId: string) {
    if (!vId || typeof vId !== 'string') return
    activeVersionId.value = vId
    activePageId.value = ''
    page.value = null
    content.value = ''
    previewHtml.value = ''
    sections.value = []
  }

  async function createVersion(form: CreateVersionReq) {
    const res = await http.post<Version>(`/tenant/themes/${themeId}/versions`, form)
    res.code === 0 ? (toast.success('版本创建成功'), loadVersions()) : toast.error(res.message)
  }

  async function updateVersion(id: string, form: { name: string; label: string }) {
    const res = await http.patch(`/tenant/versions/${id}`, form)
    res.code === 0 ? (toast.success('版本已更新'), loadVersions()) : toast.error(res.message)
  }

  async function deleteVersion(id: string) {
    const res = await http.delete(`/tenant/versions/${id}`)
    if (res.code === 0) {
      toast.success('版本已删除')
      // 立即从本地列表移除，避免依赖后续 loadVersions 的时序或缓存问题
      versions.value = versions.value.filter(v => v.id !== id)
      if (activeVersionId.value === id) {
        activeVersionId.value = ''
        sections.value = []
        activePageId.value = ''
        page.value = null
        content.value = ''
        previewHtml.value = ''
      }
      loadVersions()
    } else toast.error(res.message)
  }

  async function publishVersion(id: string): Promise<boolean> {
    const res = await http.post(`/tenant/versions/${id}/publish`)
    if (res.code === 0) {
      toast.success('已发布')
      loadVersions()
      return true
    }
    toast.error(res.message)
    return false
  }

  async function unpublishVersion(id: string) {
    const res = await http.post(`/tenant/versions/${id}/unpublish`)
    res.code === 0 ? (toast.success('已取消发布'), loadVersions()) : toast.error(res.message)
  }

  async function archiveVersion(id: string) {
    const res = await http.post(`/tenant/versions/${id}/archive`)
    res.code === 0 ? (toast.success('已归档'), loadVersions()) : toast.error(res.message)
  }

  async function unarchiveVersion(id: string) {
    const res = await http.post(`/tenant/versions/${id}/unarchive`)
    res.code === 0 ? (toast.success('已取消归档'), loadVersions()) : toast.error(res.message)
  }

  async function setDefaultVersion(id: string) {
    const res = await http.post(`/tenant/versions/${id}/set-default`)
    res.code === 0 ? (toast.success('已设为默认'), loadVersions()) : toast.error(res.message)
  }

  async function cloneVersion(id: string) {
    const name = prompt('新版本名称：')
    if (!name) return
    const res = await http.post(`/tenant/versions/${id}/clone`, { name, label: name })
    res.code === 0 ? (toast.success('克隆成功'), loadVersions()) : toast.error(res.message)
  }

  function openVersionReader(versionName: string) {
    if (theme.value) {
      window.open(`/${authStore.user?.tenant_id}/${theme.value.slug}/${versionName}`, '_blank')
    }
  }

  // ═══════════ 章节 + 文档树 ═══════════
  const sections = ref<SectionWithPages[]>([])
  const loadingTree = ref(false)

  async function loadTree() {
    if (!activeVersionId.value) return
    loadingTree.value = true
    const res = await http.get<Section[]>(`/tenant/versions/${activeVersionId.value}/sections`)
    if (res.code === 0 && res.data) {
      const pageResults = await Promise.all(
        res.data.map(s => http.get<DocPage[]>(`/tenant/sections/${s.id}/pages`)),
      )
      const result: SectionWithPages[] = res.data.map((s, i) => ({
        ...s,
        pages: pageResults[i]!.code === 0 && pageResults[i]!.data ? pageResults[i]!.data! : [],
      }))
      sections.value = result
      if (!activePageId.value && result.length > 0 && result[0]!.pages.length > 0) {
        selectPage(result[0]!.pages[0]!.id)
      }
    }
    loadingTree.value = false
  }

  watch(activeVersionId, () => { if (activeVersionId.value) loadTree() })

  async function createSection(title: string) {
    if (!activeVersionId.value) { toast.error('请先选择或创建一个版本'); return }
    const res = await http.post(`/tenant/versions/${activeVersionId.value}/sections`, { title })
    res.code === 0 ? (toast.success('章节创建成功'), loadTree()) : toast.error(res.message)
  }

  async function updateSection(id: string, title: string) {
    const res = await http.patch(`/tenant/sections/${id}`, { title })
    res.code === 0 ? (toast.success('章节已更新'), loadTree()) : toast.error(res.message)
  }

  async function deleteSection(id: string) {
    const res = await http.delete(`/tenant/sections/${id}`)
    res.code === 0 ? (toast.success('已删除'), loadTree()) : toast.error(res.message)
  }

  async function sortSections(orderedIds: string[]) {
    const items = orderedIds.map((id, i) => ({ id, sort_order: i }))
    const res = await http.put(`/tenant/versions/${activeVersionId.value}/sections/sort`, { items })
    if (res.code !== 0) toast.error(res.message)
  }

  async function createPage(sectionId: string, form: { title: string; slug: string; content: string }) {
    const res = await http.post<DocPage>(`/tenant/sections/${sectionId}/pages`, form)
    if (res.code === 0 && res.data) {
      toast.success('文档页创建成功')
      const newPage = res.data
      const target = sections.value.find(s => s.id === sectionId)
      if (target) {
        target.pages = [...target.pages, newPage]
      } else {
        loadTree()
      }
      activePageId.value = newPage.id
      page.value = newPage
      title.value = newPage.title
      slug.value = newPage.slug
      content.value = newPage.content
      lastSaved.value = ''
    } else {
      toast.error(res.message)
    }
  }

  async function updatePage(id: string, form: { title: string; slug: string }) {
    const res = await http.patch(`/tenant/pages/${id}`, form)
    if (res.code === 0) {
      toast.success('文档页已更新')
      if (activePageId.value === id) {
        title.value = form.title
        slug.value = form.slug
        if (page.value) { page.value.title = form.title; page.value.slug = form.slug }
      }
      loadTree()
    } else toast.error(res.message)
  }

  async function deletePage(id: string) {
    const res = await http.delete(`/tenant/pages/${id}`)
    if (res.code === 0) {
      toast.success('已删除')
      if (activePageId.value === id) {
        activePageId.value = ''
        page.value = null
        content.value = ''
        previewHtml.value = ''
      }
      loadTree()
    } else toast.error(res.message)
  }

  // ═══════════ 文档编辑 ═══════════
  const activePageId = ref('')
  const page = ref<DocPage | null>(null)
  const title = ref('')
  const slug = ref('')
  const content = ref('')
  const saving = ref(false)
  const lastSaved = ref('')

  const { html: previewHtml, enabled: showPreview } = useMarkdownPreview(content)
  const { uploadFile, uploadImage, uploading } = useMedia()

  const { tocItems, activeId: tocActiveId, scrollTo: tocScrollTo, cleanup: cleanupToc, refresh: refreshToc } = useToc(() => previewHtml.value)

  async function selectPage(pageId: string) {
    activePageId.value = pageId
    for (const s of sections.value) {
      const cached = s.pages.find(p => p.id === pageId)
      if (cached) {
        page.value = cached
        title.value = cached.title
        slug.value = cached.slug
        content.value = cached.content
        lastSaved.value = ''
        return
      }
    }
    const res = await http.get<DocPage>(`/tenant/pages/${pageId}`)
    if (res.code === 0 && res.data) {
      page.value = res.data
      title.value = res.data.title
      slug.value = res.data.slug
      content.value = res.data.content
      lastSaved.value = ''
    }
  }

  function syncPageCache(id: string, fields: Partial<Pick<DocPage, 'title' | 'slug' | 'content'>>) {
    for (const s of sections.value) {
      const p = s.pages.find(p => p.id === id)
      if (p) { Object.assign(p, fields); break }
    }
  }

  async function savePage() {
    if (!activePageId.value) return
    saving.value = true
    const res = await http.put(`/tenant/pages/${activePageId.value}`, {
      title: title.value, slug: slug.value, content: content.value,
    })
    saving.value = false
    if (res.code === 0) {
      lastSaved.value = new Date().toLocaleTimeString()
      const fields = { title: title.value, slug: slug.value, content: content.value }
      if (page.value) Object.assign(page.value, fields)
      syncPageCache(activePageId.value, fields)
      toast.success('已保存')
    } else toast.error(res.message)
  }

  async function autoSave() {
    if (!content.value || !page.value) return
    if (content.value === page.value.content && title.value === page.value.title && slug.value === page.value.slug) return
    const res = await http.patch(`/tenant/pages/${activePageId.value}`, { content: content.value, title: title.value, slug: slug.value })
    if (res.code === 0) {
      lastSaved.value = new Date().toLocaleTimeString()
      const fields = { title: title.value, slug: slug.value, content: content.value }
      if (page.value) Object.assign(page.value, fields)
      syncPageCache(activePageId.value, fields)
    }
  }

  async function publishPage() {
    const res = await http.post(`/tenant/pages/${activePageId.value}/publish`)
    if (res.code === 0) { toast.success('已发布'); selectPage(activePageId.value); loadTree() }
    else toast.error(res.message)
  }

  async function unpublishPage() {
    const res = await http.post(`/tenant/pages/${activePageId.value}/unpublish`)
    if (res.code === 0) { toast.success('已下线'); selectPage(activePageId.value); loadTree() }
    else toast.error(res.message)
  }

  // 乐观同步：title 变化时立即更新 sections 中对应 page 的标题
  watch(title, (newTitle) => {
    if (!activePageId.value) return
    for (const s of sections.value) {
      const p = s.pages.find(p => p.id === activePageId.value)
      if (p) { p.title = newTitle; break }
    }
  })

  // ── 自动保存 ──
  const { start: startAutoSave, stop: stopAutoSave } = useAutoSave({
    save: autoSave,
    isDirty: () => !!content.value && !!page.value
      && (content.value !== page.value!.content || title.value !== page.value!.title),
  })

  // ═══════════ 编辑模式 ═══════════
  const editorMode = ref<EditorMode>('visual')

  // 切换模式后重新提取 TOC：新模式的 .prose DOM 已就位，不依赖 previewHtml 就能居然刷新
  watch(editorMode, () => {
    nextTick(() => refreshToc())
  })

  // ═══════════ 设置面板 ═══════════
  const showSettingsPanel = ref(false)

  // ═══════════ 生命周期 ═══════════
  async function init() {
    await Promise.all([loadTheme(), loadVersions()])
    startAutoSave()
  }

  function destroy() {
    stopAutoSave()
    cleanupToc()
  }

  return {
    theme, loadTheme, updateTheme,
    versions, activeVersionId, loadVersions, switchVersion,
    createVersion, updateVersion, deleteVersion,
    publishVersion, unpublishVersion, archiveVersion, unarchiveVersion, setDefaultVersion, cloneVersion, openVersionReader,
    sections, loadingTree, loadTree,
    createSection, updateSection, deleteSection, sortSections,
    createPage, updatePage, deletePage,
    activePageId, page, title, slug, content, saving, lastSaved,
    selectPage, savePage, publishPage, unpublishPage,
    showPreview, previewHtml, tocItems, tocActiveId, tocScrollTo,
    uploading, uploadImage, uploadFile,
    editorMode,
    showSettingsPanel,
    init, destroy,
  }
}
