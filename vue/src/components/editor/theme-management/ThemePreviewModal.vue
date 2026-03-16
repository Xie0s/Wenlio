<!-- ThemePreviewModal.vue - 主题内容预览弹窗
     职责：以全屏模态形式展示主题的版本/章节/文章树形目录及简约文章内容预览，供编辑器列表页使用
     数据通过公开 API 加载，局部 pageCache 避免重复请求，不依赖全局 reader store
     对外暴露：
       Props: theme(Theme), open(boolean)
       Emits: update:open(boolean) -->
<script setup lang="ts">
import { ref, watch, computed, onUnmounted } from 'vue'
import { http } from '@/utils/http'
import { renderMarkdown } from '@/lib/markdown'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import DocContent from '@/components/preview/DocContent.vue'
import type { Theme, Version, Section, DocPage } from '@/utils/types'
import type { SectionWithPages } from '@/components/editor/edit-page/types'
import { X, ChevronDown, FileText, BookOpen } from 'lucide-vue-next'

const props = defineProps<{
  theme: Theme
  open: boolean
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

// ── 数据 ──
const versions = ref<Version[]>([])
const activeVersionId = ref('')
const tree = ref<SectionWithPages[]>([])
const activePageId = ref('')
const renderedHtml = ref('')
const activePageTitle = ref('')

// ── 加载状态 ──
const loadingVersions = ref(false)
const loadingTree = ref(false)
const loadingPage = ref(false)

// ── 局部页面缓存，避免重复请求 ──
const pageCache = new Map<string, DocPage>()

// ── 章节折叠状态 ──
const sectionCollapsed = ref<Record<string, boolean>>({})

// ── 已加载的版本 id，避免重复加载树 ──
const loadedVersionId = ref('')

function close() {
  emit('update:open', false)
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

watch(() => props.open, async (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', onKeyDown)
    await loadVersions()
  } else {
    document.body.style.overflow = 'auto'
    window.removeEventListener('keydown', onKeyDown)
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  document.body.style.overflow = 'auto'
})

async function loadVersions() {
  loadingVersions.value = true
  const res = await http.get<Version[]>(`/tenant/themes/${props.theme.id}/versions`)
  loadingVersions.value = false
  if (res.code === 0 && res.data && res.data.length > 0) {
    versions.value = res.data
    const defaultVer = res.data.find(v => v.is_default) || res.data[0]!
    if (activeVersionId.value !== defaultVer.id) {
      activeVersionId.value = defaultVer.id
    }
    if (loadedVersionId.value !== activeVersionId.value) {
      await loadTree(activeVersionId.value)
    }
  }
}

async function loadTree(versionId: string) {
  loadingTree.value = true
  tree.value = []
  sectionCollapsed.value = {}
  activePageId.value = ''
  renderedHtml.value = ''
  activePageTitle.value = ''
  const res = await http.get<Section[]>(`/tenant/versions/${versionId}/sections`)
  loadingTree.value = false
  if (res.code === 0 && res.data) {
    const pageResults = await Promise.all(
      res.data.map(section => http.get<DocPage[]>(`/tenant/sections/${section.id}/pages`)),
    )
    tree.value = res.data.map((section, index) => ({
      ...section,
      pages: pageResults[index]!.code === 0 && pageResults[index]!.data ? pageResults[index]!.data! : [],
    }))
    loadedVersionId.value = versionId
    // 自动选中第一篇文章
    for (const section of tree.value) {
      if (section.pages.length > 0) {
        await selectPage(section.pages[0]!.id, section.pages[0]!.title)
        break
      }
    }
  }
}

async function onVersionChange(versionId: string) {
  if (activeVersionId.value === versionId) return
  activeVersionId.value = versionId
  await loadTree(versionId)
}

async function selectPage(pageId: string, title: string) {
  if (activePageId.value === pageId) return
  activePageId.value = pageId
  activePageTitle.value = title

  for (const section of tree.value) {
    const page = section.pages.find(item => item.id === pageId)
    if (page) {
      pageCache.set(pageId, page)
      renderedHtml.value = await renderMarkdown(page.content || '')
      return
    }
  }

  const cached = pageCache.get(pageId)
  if (cached) {
    renderedHtml.value = await renderMarkdown(cached.content || '')
    return
  }

  loadingPage.value = true
  renderedHtml.value = ''
  const res = await http.get<DocPage>(`/tenant/pages/${pageId}`)
  loadingPage.value = false
  if (res.code === 0 && res.data) {
    pageCache.set(pageId, res.data)
    renderedHtml.value = await renderMarkdown(res.data.content || '')
  }
}

function toggleSection(id: string) {
  sectionCollapsed.value[id] = !sectionCollapsed.value[id]
}

const activeVersion = computed(() => versions.value.find(v => v.id === activeVersionId.value))

const totalPages = computed(() => tree.value.reduce((n, s) => n + s.pages.length, 0))
</script>

<template>
  <Teleport to="body">
    <!-- 遮罩 -->
    <Transition name="preview-fade">
      <div
        v-if="open"
        class="fixed inset-0 bg-black/30 z-[90] backdrop-blur-[2px]"
        @click="close"
      />
    </Transition>

    <!-- 弹窗 -->
    <Transition name="preview-scale">
      <div
        v-if="open"
        class="fixed inset-0 z-[100] flex items-center justify-center p-4 pointer-events-none"
      >
        <div
          class="pointer-events-auto relative flex w-full max-w-5xl overflow-hidden rounded-3xl border border-border/40 bg-background shadow-2xl"
          style="height: min(82vh, 720px)"
          @click.stop
        >

          <!-- 左侧：目录树 -->
          <div class="flex w-56 shrink-0 flex-col border-r border-border/40">

            <!-- 左侧头部：主题名 + 版本选择 -->
            <div class="shrink-0 border-b border-border/30 px-4 py-4 space-y-3">
              <div class="flex items-center gap-2 min-w-0">
                <BookOpen class="h-4 w-4 shrink-0 text-primary/70" :stroke-width="1.5" />
                <span class="truncate text-base font-normal text-foreground">{{ theme.name }}</span>
              </div>
              <Select
                v-if="versions.length > 1"
                :model-value="activeVersionId"
                @update:model-value="(v) => v && onVersionChange(v as string)"
              >
                <SelectTrigger class="h-8 w-full rounded-lg text-sm">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent class="rounded-xl">
                  <SelectItem v-for="v in versions" :key="v.id" :value="v.id" class="text-sm">
                    {{ v.label || v.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <div v-else-if="activeVersion" class="text-sm text-muted-foreground px-0.5">
                {{ activeVersion.label || activeVersion.name }}
              </div>
            </div>

            <!-- 目录树区域 -->
            <div class="flex-1 overflow-y-auto py-3 [scrollbar-width:thin]">

              <!-- 加载中 -->
              <div v-if="loadingVersions || loadingTree" class="flex items-center justify-center pt-10">
                <div class="preview-spinner" />
              </div>

              <!-- 空状态 -->
              <div v-else-if="tree.length === 0" class="px-4 pt-8 text-center">
                <FileText class="mx-auto h-8 w-8 text-muted-foreground/30" :stroke-width="1" />
                <p class="mt-2 text-sm text-muted-foreground">暂无内容</p>
              </div>

              <!-- 章节列表 -->
              <nav v-else class="space-y-1 px-2">
                <div v-for="section in tree" :key="section.id">
                  <!-- 章节标题 -->
                  <button
                    class="flex w-full items-center justify-between rounded-lg px-2 py-1.5 text-left transition-colors hover:bg-accent/60"
                    @click="toggleSection(section.id)"
                  >
                    <span class="truncate text-sm font-light text-foreground/80">
                      {{ section.title }}
                    </span>
                    <ChevronDown
                      class="h-3 w-3 shrink-0 text-muted-foreground transition-transform duration-200"
                      :class="sectionCollapsed[section.id] ? '-rotate-90' : ''"
                    />
                  </button>

                  <!-- 文章列表 -->
                  <ul v-if="!sectionCollapsed[section.id]" class="relative space-y-0.5 pl-3 pb-1">
                    <div class="absolute left-[13px] top-0 bottom-0 w-px bg-border/50" />
                    <li
                      v-for="page in section.pages"
                      :key="page.id"
                      class="relative -ml-[1px] rounded-r-lg transition-colors duration-150"
                      :class="page.id === activePageId ? 'bg-accent/50' : 'hover:bg-accent/40'"
                    >
                      <span
                        v-if="page.id === activePageId"
                        class="absolute left-0 inset-y-0.5 w-0.5 rounded-full bg-primary"
                      />
                      <button
                        class="w-full truncate py-1.5 pl-4 pr-2 text-left text-sm transition-colors duration-150"
                        :class="page.id === activePageId
                          ? 'text-primary font-normal'
                          : 'text-foreground/70 hover:text-foreground font-thin'"
                        @click="selectPage(page.id, page.title)"
                      >
                        {{ page.title }}
                      </button>
                    </li>
                  </ul>
                </div>
              </nav>
            </div>

            <!-- 左侧底部：统计信息 -->
            <div class="shrink-0 border-t border-border/30 px-4 py-2.5 flex items-center gap-3 text-muted-foreground/60">
              <span class="text-sm">{{ tree.length }} 章节</span>
              <span class="w-px h-3 bg-border/50" />
              <span class="text-sm">{{ totalPages }} 篇文章</span>
            </div>
          </div>

          <!-- 右侧：文章内容 -->
          <div class="flex flex-1 min-w-0 flex-col">

            <!-- 右侧头部：文章标题 + 关闭按钮 -->
            <div class="shrink-0 flex items-center gap-3 border-b border-border/30 px-6 py-4">
              <h2 class="flex-1 truncate text-base font-semibold text-foreground">
                {{ activePageTitle || theme.name }}
              </h2>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-8 w-8 shrink-0 rounded-full"
                    @click="close"
                  >
                    <X class="h-4 w-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent class="z-[110]">关闭预览</TooltipContent>
              </Tooltip>
            </div>

            <!-- 内容区 -->
            <div class="flex-1 min-h-0 overflow-y-auto px-6 py-6 [scrollbar-width:thin]">

              <!-- 加载中 -->
              <div v-if="loadingPage" class="flex h-40 items-center justify-center">
                <div class="preview-spinner" />
              </div>

              <!-- 有内容 -->
              <DocContent v-else-if="renderedHtml" :html="renderedHtml" />

              <!-- 空状态 -->
              <div v-else class="flex h-40 flex-col items-center justify-center gap-2">
                <FileText class="h-10 w-10 text-muted-foreground/20" :stroke-width="1" />
                <p class="text-sm text-muted-foreground">
                  {{ loadingTree ? '正在加载目录…' : '从左侧选择文章开始预览' }}
                </p>
              </div>
            </div>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.preview-fade-enter-active,
.preview-fade-leave-active {
  transition: opacity 0.2s ease;
}
.preview-fade-enter-from,
.preview-fade-leave-to {
  opacity: 0;
}

.preview-scale-enter-active,
.preview-scale-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.preview-scale-enter-from,
.preview-scale-leave-to {
  opacity: 0;
  transform: scale(0.96) translateY(8px);
}

.preview-spinner {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid transparent;
  border-top-color: var(--color-primary);
  animation: preview-spin 0.7s linear infinite;
}

@keyframes preview-spin {
  to { transform: rotate(360deg); }
}
</style>
