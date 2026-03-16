/**
 * category-tag.ts - 分类与标签业务逻辑层
 *
 * 职责：封装分类树和标签列表的状态管理与 API 操作（CRUD、排序）
 * 对外暴露：useCategoryList()、useTagList() composable 函数
 */
import { ref, reactive } from 'vue'
import { http } from '@/utils/http'
import type { Category, Tag, CreateCategoryReq, CreateTagReq } from '@/utils/types'
import { toast } from 'vue-sonner'

export type CategoryListLib = ReturnType<typeof useCategoryList>
export type TagListLib = ReturnType<typeof useTagList>

// ============================================================
// 分类
// ============================================================

export function useCategoryList() {
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const manageMode = ref(false)

  // 选中分类（用于主题列表筛选）
  const selectedCategoryId = ref<string>('')
  function selectCategory(id: string) {
    selectedCategoryId.value = selectedCategoryId.value === id ? '' : id
  }

  // 创建弹窗
  const showCreate = ref(false)
  const form = reactive<CreateCategoryReq>({ name: '', slug: '', description: '', parent_id: '' })

  // 编辑弹窗
  const showEdit = ref(false)
  const editId = ref('')
  const editForm = reactive({ name: '', slug: '', description: '' })

  async function loadCategories() {
    loading.value = true
    const res = await http.get<Category[]>('/tenant/categories')
    loading.value = false
    if (res.code === 0 && res.data) {
      categories.value = res.data
    }
  }

  async function createCategory(): Promise<boolean> {
    const payload: CreateCategoryReq = {
      name: form.name,
      slug: form.slug,
      description: form.description,
    }
    if (form.parent_id) payload.parent_id = form.parent_id
    const res = await http.post('/tenant/categories', payload)
    if (res.code === 0) {
      toast.success('分类创建成功')
      showCreate.value = false
      resetForm()
      loadCategories()
      return true
    } else {
      toast.error(res.message)
      return false
    }
  }

  async function createChildCategory(parentId: string, data: { name: string; slug: string }): Promise<boolean> {
    const res = await http.post('/tenant/categories', { name: data.name, slug: data.slug, parent_id: parentId })
    if (res.code === 0) {
      toast.success('子分类创建成功')
      await loadCategories()
      return true
    } else {
      toast.error(res.message)
      return false
    }
  }

  function openEdit(cat: Category) {
    editId.value = cat.id
    Object.assign(editForm, { name: cat.name, slug: cat.slug, description: cat.description || '' })
    showEdit.value = true
  }

  function startInlineEdit(cat: Category) {
    editId.value = cat.id
    Object.assign(editForm, { name: cat.name, slug: cat.slug, description: cat.description || '' })
  }

  async function updateCategory() {
    const res = await http.patch(`/tenant/categories/${editId.value}`, editForm)
    if (res.code === 0) {
      toast.success('分类已更新')
      showEdit.value = false
      editId.value = ''
      loadCategories()
    } else {
      toast.error(res.message)
    }
  }

  async function deleteCategory(id: string) {
    const res = await http.delete(`/tenant/categories/${id}`)
    if (res.code === 0) {
      toast.success('已删除')
      loadCategories()
    } else {
      toast.error(res.message)
    }
  }

  function resetForm() {
    Object.assign(form, { name: '', slug: '', description: '', parent_id: '' })
  }

  function openCreate(parentId?: string) {
    resetForm()
    if (parentId) form.parent_id = parentId
    showCreate.value = true
  }

  return {
    categories,
    loading,
    manageMode,
    selectedCategoryId,
    selectCategory,
    showCreate,
    form,
    showEdit,
    editId,
    editForm,
    loadCategories,
    openCreate,
    createCategory,
    createChildCategory,
    openEdit,
    startInlineEdit,
    updateCategory,
    deleteCategory,
  }
}

// ============================================================
// 标签
// ============================================================

export function useTagList(options?: { afterDeleteTag?: () => void }) {
  const tags = ref<Tag[]>([])
  const loading = ref(false)

  // 选中标签（用于主题列表筛选，多选）
  const selectedTagIds = ref<string[]>([])
  function toggleTag(id: string) {
    const idx = selectedTagIds.value.indexOf(id)
    if (idx >= 0) selectedTagIds.value.splice(idx, 1)
    else selectedTagIds.value.push(id)
  }

  // 创建弹窗
  const showCreate = ref(false)
  const form = reactive<CreateTagReq>({ name: '', slug: '', description: '' })

  // 编辑弹窗
  const showEdit = ref(false)
  const editId = ref('')
  const editForm = reactive({ name: '', slug: '', description: '' })

  async function loadTags() {
    loading.value = true
    const res = await http.get<Tag[]>('/tenant/tags')
    loading.value = false
    if (res.code === 0 && res.data) {
      tags.value = res.data
    }
  }

  function startInlineEdit(tag: Tag) {
    editId.value = tag.id
    Object.assign(editForm, { name: tag.name, slug: tag.slug, description: tag.description || '' })
  }

  async function createTag(): Promise<boolean> {
    const res = await http.post('/tenant/tags', { ...form })
    if (res.code === 0) {
      toast.success('标签创建成功')
      showCreate.value = false
      Object.assign(form, { name: '', slug: '', description: '' })
      loadTags()
      return true
    } else {
      toast.error(res.message)
      return false
    }
  }

  function openEdit(tag: Tag) {
    editId.value = tag.id
    Object.assign(editForm, { name: tag.name, slug: tag.slug, description: tag.description || '' })
    showEdit.value = true
  }

  async function updateTag() {
    const res = await http.patch(`/tenant/tags/${editId.value}`, editForm)
    if (res.code === 0) {
      toast.success('标签已更新')
      showEdit.value = false
      editId.value = ''
      loadTags()
    } else {
      toast.error(res.message)
    }
  }

  async function deleteTag(id: string, force = false) {
    const url = force ? `/tenant/tags/${id}?force=true` : `/tenant/tags/${id}`
    const res = await http.delete(url)
    if (res.code === 0) {
      toast.success('已删除')
      loadTags()
      options?.afterDeleteTag?.()
    } else if (res.code === 409103) {
      await loadTags()
    } else {
      toast.error(res.message)
    }
  }

  return {
    tags,
    loading,
    selectedTagIds,
    toggleTag,
    showCreate,
    form,
    showEdit,
    editId,
    editForm,
    loadTags,
    createTag,
    openEdit,
    startInlineEdit,
    updateTag,
    deleteTag,
  }
}
