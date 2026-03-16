<!-- ThemeListPage.vue - 主题管理页面
     职责：主题列表展示与管理，分类/标签侧边栏管理
     对外暴露：路由页面组件（sidebarLayout 模式，AdminLayout 提供 h-screen overflow-hidden 固定高度容器，侧边栏固定 + 右侧列表独立滚动） -->
<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useThemeList } from '@/lib/theme-list'
import { useCategoryList, useTagList } from '@/lib/category-tag'
import type { Theme } from '@/utils/types'
import { toast } from 'vue-sonner'
import ThemeCard from '@/components/editor/theme-management/ThemeCard.vue'
import ThemeFormDialog from '@/components/editor/theme-management/ThemeFormDialog.vue'
import type { ThemeFormData } from '@/components/editor/theme-management/ThemeFormDialog.vue'
import ThemeListSidebar from '@/components/editor/theme-management/ThemeListSidebar.vue'
import { Button } from '@/components/ui/button'
import { BookOpen, SearchX, X } from 'lucide-vue-next'

const router = useRouter()
const { themes, loading, loadThemes, createTheme, deleteTheme, deleteThemeCascade, updateTheme } = useThemeList()
const catLib = useCategoryList()
const tagLib = useTagList({
  afterDeleteTag: () => loadThemes(catLib.selectedCategoryId.value || undefined, tagLib.selectedTagIds.value.slice()),
})

const eyeCareMode = ref(false)

watch(eyeCareMode, (on) => {
  if (on) document.body.classList.add('eye-care')
  else document.body.classList.remove('eye-care')
})

onUnmounted(() => {
  document.body.classList.remove('eye-care')
})

function handleRefresh() {
  loadThemes(catLib.selectedCategoryId.value || undefined, tagLib.selectedTagIds.value.slice())
  catLib.loadCategories()
  tagLib.loadTags()
}

function findTagName(id: string): string {
  return tagLib.tags.value.find(t => t.id === id)?.name ?? ''
}

function findCategoryName(id: string): string {
  const find = (list: typeof catLib.categories.value): string => {
    for (const cat of list) {
      if (cat.id === id) return cat.name
      if (cat.children?.length) {
        const name = find(cat.children)
        if (name) return name
      }
    }
    return ''
  }
  return find(catLib.categories.value)
}

// ── Dialog 状态 ──────────────────────────────────────────────
const showCreate = ref(false)
const showEdit = ref(false)
const submitting = ref(false)
const editingTheme = ref<Theme | null>(null)

function openEdit(t: Theme) {
  editingTheme.value = t
  showEdit.value = true
}

async function handleCreate(data: ThemeFormData) {
  submitting.value = true
  const ok = await createTheme(data)
  submitting.value = false
  if (ok) showCreate.value = false
}

async function handleUpdate(data: ThemeFormData) {
  if (!editingTheme.value) return
  submitting.value = true
  const ok = await updateTheme(editingTheme.value.id, data)
  submitting.value = false
  if (ok) showEdit.value = false
}

// ── 删除状态 ────────────────────────────────────────────────
const THEME_HAS_VERSIONS_CODE = 422403
const cascadeId = ref('')
const deletingId = ref('')

async function requestDelete(theme: Theme) {
  if ((theme.version_count ?? 0) > 0) {
    cascadeId.value = theme.id
    return
  }
  if (deletingId.value) return
  deletingId.value = theme.id
  let res: { code: number; message: string } | null = null
  try {
    res = await deleteTheme(theme.id, { silentError: true })
  } catch {
    toast.error('删除失败，请稍后重试')
    return
  } finally {
    deletingId.value = ''
  }
  if (!res) return
  if (res.code === 0) {
    if (cascadeId.value === theme.id) cascadeId.value = ''
    return
  }
  if (res.code === THEME_HAS_VERSIONS_CODE) {
    cascadeId.value = theme.id
    return
  }
  toast.error(res.message)
}

function openReader(tenantId: string, slug: string) {
  window.open(`/${tenantId}/${slug}`, '_blank')
}

async function confirmCascadeDelete(id: string) {
  if (deletingId.value) return
  deletingId.value = id
  try {
    const res = await deleteThemeCascade(id)
    if (res.code === 0 && cascadeId.value === id) cascadeId.value = ''
  } catch {
    toast.error('删除失败，请稍后重试')
  } finally {
    deletingId.value = ''
  }
}

onMounted(() => {
  loadThemes()
  catLib.loadCategories()
  tagLib.loadTags()
})

watch(
  [() => catLib.selectedCategoryId.value, () => tagLib.selectedTagIds.value.slice()],
  ([categoryId, tagIds]) => {
    loadThemes(categoryId || undefined, tagIds as string[])
  }
)
</script>

<template>
  <!-- AdminLayout sidebarLayout 模式提供 h-screen overflow-hidden，此处 flex h-full，左侧栏全高 -->
  <div class="relative flex h-full">
    <div v-if="loading" class="absolute inset-0 z-10 flex items-center justify-center">
      <div class="google-spinner" />
    </div>

    <!-- 左侧固定侧边栏（全高） -->
    <ThemeListSidebar class="pl-6 pt-6" :cat-lib="catLib" :tag-lib="tagLib" :eye-care-mode="eyeCareMode"
      @toggle-eye-care="eyeCareMode = !eyeCareMode" @refresh="handleRefresh" />

    <!-- 右侧区域：操作栏 + 列表 -->
    <div class="flex flex-col flex-1 min-w-0 min-h-0">

      <!-- ════════ 操作栏（新建主题居中，计数固定右侧）════════ -->
      <div class="shrink-0 relative flex items-center justify-center gap-3 pt-6 pb-5 px-6">
        <Button variant="outline" class="rounded-full h-11 px-6 text-base font-light" @click="showCreate = true">
          新建主题
        </Button>

        <!-- 分类筛选状态 -->
        <div v-if="catLib.selectedCategoryId.value"
          class="inline-flex items-center gap-2 rounded-full border border-border/60 bg-muted/30 pl-4 pr-1.5 h-11">
          <span class="text-xs text-muted-foreground shrink-0">分类</span>
          <span class="text-sm text-foreground/80 shrink-0">{{ findCategoryName(catLib.selectedCategoryId.value)
          }}</span>
          <button title="清除分类筛选"
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-destructive/20 text-destructive hover:bg-destructive/30 transition-colors"
            @click="catLib.selectCategory(catLib.selectedCategoryId.value)">
            <X class="h-3.5 w-3.5" />
          </button>
        </div>

        <!-- 标签筛选状态（多选，合并到单个容器）-->
        <div v-if="tagLib.selectedTagIds.value.length > 0"
          class="inline-flex items-center gap-2 rounded-full border border-border/60 bg-muted/30 pl-4 pr-1.5 h-11">
          <span class="text-xs text-muted-foreground shrink-0">标签</span>
          <div class="flex items-center gap-1.5">
            <span v-for="tagId in tagLib.selectedTagIds.value.slice(0, 7)" :key="tagId"
              class="inline-flex items-center gap-1 rounded-full bg-foreground/10 pl-2.5 pr-1 py-0.5 text-sm text-foreground/80">
              #{{ findTagName(tagId) }}
              <button
                class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full hover:bg-destructive/20 hover:text-destructive transition-colors"
                @click="tagLib.toggleTag(tagId)">
                <X class="h-2.5 w-2.5" />
              </button>
            </span>
            <span v-if="tagLib.selectedTagIds.value.length > 7"
              class="inline-flex items-center rounded-full bg-foreground/10 px-2.5 py-0.5 text-sm text-muted-foreground shrink-0">
              +{{ tagLib.selectedTagIds.value.length - 7 }}
            </span>
          </div>
          <button title="清除标签筛选"
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-destructive/20 text-destructive hover:bg-destructive/30 transition-colors"
            @click="tagLib.selectedTagIds.value.splice(0)">
            <X class="h-3.5 w-3.5" />
          </button>
        </div>

        <!-- 计数固定右侧：有筛选时始终显示（含0结果），无筛选时 themes > 0 才显示 -->
        <span
          v-if="!loading && (themes.length > 0 || catLib.selectedCategoryId.value || tagLib.selectedTagIds.value.length > 0)"
          class="absolute right-6 rounded-full bg-muted px-2.5 py-0.5 text-sm tabular-nums text-muted-foreground">
          共 {{ themes.length }} 个主题
        </span>
      </div>

      <!-- 主题列表（独立滚动） -->
      <div class="flex-1 min-h-0 overflow-y-auto pb-24 px-6">
        <!-- 空状态 -->
        <div v-if="!loading && themes.length === 0" class="flex min-h-full flex-col items-center justify-center px-6">
          <div class="w-full max-w-sm">
            <div class="flex flex-col items-center text-center">
              <div
                class="mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-muted/50 border border-border/40">
                <SearchX v-if="catLib.selectedCategoryId.value || tagLib.selectedTagIds.value.length > 0"
                  class="h-7 w-7 text-muted-foreground/50" :stroke-width="1.5" />
                <BookOpen v-else class="h-7 w-7 text-muted-foreground/50" :stroke-width="1.5" />
              </div>
              <h3 class="text-xl font-normal tracking-tight text-foreground">
                {{ (catLib.selectedCategoryId.value || tagLib.selectedTagIds.value.length > 0) ? '未找到匹配的主题' : '暂无主题' }}
              </h3>
              <p class="mt-2 text-sm font-light text-muted-foreground/70 leading-relaxed">
                {{ (catLib.selectedCategoryId.value || tagLib.selectedTagIds.value.length > 0)
                  ? '当前筛选条件下没有匹配的主题，请尝试调整或清除筛选。'
                  : '点击上方「新建主题」按钮，创建第一个文档主题。' }}
              </p>
              <div v-if="catLib.selectedCategoryId.value || tagLib.selectedTagIds.value.length > 0"
                class="mt-6 flex flex-wrap justify-center gap-2">
                <button v-if="catLib.selectedCategoryId.value"
                  class="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/30 px-4 h-9 text-sm text-muted-foreground hover:bg-muted/60 transition-colors"
                  @click="catLib.selectCategory(catLib.selectedCategoryId.value)">
                  <X class="h-3.5 w-3.5" />清除分类
                </button>
                <button v-if="tagLib.selectedTagIds.value.length > 0"
                  class="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/30 px-4 h-9 text-sm text-muted-foreground hover:bg-muted/60 transition-colors"
                  @click="tagLib.selectedTagIds.value.splice(0)">
                  <X class="h-3.5 w-3.5" />清除标签
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 主题列表（分隔线内置于 ThemeCard） -->
        <div v-else-if="themes.length > 0">
          <ThemeCard v-for="t in themes" :key="t.id" :theme="t" :deleting="deletingId === t.id"
            :cascade-open="cascadeId === t.id" :cat-lib="catLib" :tag-lib="tagLib"
            @enter="router.push(`/admin/themes/${t.id}`)" @open-edit="openEdit(t)"
            @view="openReader(t.tenant_id, t.slug)" @request-delete="requestDelete(t)"
            @cascade-confirm="confirmCascadeDelete(t.id)" @cascade-cancel="cascadeId = ''" />
        </div>
      </div>

    </div>
  </div>

  <!-- ══════════ 创建主题 Dialog ══════════ -->
  <ThemeFormDialog :open="showCreate" mode="create" :cat-lib="catLib" :tag-lib="tagLib" :submitting="submitting"
    @update:open="showCreate = $event" @submit="handleCreate" />

  <!-- ══════════ 编辑主题 Dialog ══════════ -->
  <ThemeFormDialog :open="showEdit" mode="edit" :initial-values="editingTheme ? {
    name: editingTheme.name,
    slug: editingTheme.slug,
    description: editingTheme.description,
    category_id: editingTheme.category_id,
    tag_ids: editingTheme.tag_ids,
    access_mode: (editingTheme.access_mode || 'public') as any,
  } : undefined" :cat-lib="catLib" :tag-lib="tagLib" :submitting="submitting" @update:open="showEdit = $event"
    @submit="handleUpdate" />
</template>

<style scoped>
.google-spinner {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: conic-gradient(#4285F4 0deg 90deg,
      #EA4335 90deg 180deg,
      #FBBC05 180deg 270deg,
      #34A853 270deg 360deg);
  -webkit-mask: radial-gradient(farthest-side, transparent calc(100% - 3.5px), #000 calc(100% - 3.5px));
  mask: radial-gradient(farthest-side, transparent calc(100% - 3.5px), #000 calc(100% - 3.5px));
  animation: google-spin 0.8s linear infinite;
}

@keyframes google-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
