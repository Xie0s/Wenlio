<!--
  hero-demo/EditorDemo.vue
  职责：演示文档编辑场景，复刻 ThemeEditorPage 布局
  功能：可视化富文本编辑（Tiptap）+ 源码 Markdown 编辑，支持模式切换
  默认模式：可视化编辑（visual）
  无 store/router 依赖，通过 provide 提供 mock 上下文供子组件使用
-->
<script setup lang="ts">
import { ref, computed, watch, watchEffect, onMounted, provide } from 'vue'
import { renderMarkdown } from '@/lib/markdown'
import { PROSE_CLASSES } from '@/lib/prose'
import { EDITOR_TREE, EDITOR_PAGES, EDITOR_INITIAL_PAGE_ID } from './demo-data'
import { useVisualEditor } from '@/components/editor/edit-page/visual-editor/useVisualEditor'
import { useMarkdownPreview } from '@/composables/useMarkdownPreview'
import { EDITOR_KEY, type EditorMode } from '@/composables/useThemeEditor'
import VisualEditorToolbar from '@/components/editor/edit-page/visual-editor/VisualEditorToolbar.vue'
import VisualEditorContent from '@/components/editor/edit-page/visual-editor/VisualEditorContent.vue'
import {
  ChevronRight, FolderOpen, FileText,
  Check, Upload, Settings, Eye, EyeOff, Rocket,
  Plus, Code, PenLine,
} from 'lucide-vue-next'

// ── 编辑模式（数据驱动，默认可视化） ────────────────────────────
const MODES = [
  { id: 'visual' as EditorMode, label: '可视化', icon: PenLine },
  { id: 'source' as EditorMode, label: '源码',   icon: Code    },
] as const

const editorMode = ref<EditorMode>('visual')

// ── 页面状态 ───────────────────────────────────────────────────
const activePageId = ref(EDITOR_INITIAL_PAGE_ID)

const pageContents = ref<Record<string, string>>(
  Object.fromEntries(
    Object.entries(EDITOR_PAGES).map(([id, p]) => [id, p.content]),
  ),
)

const content = computed({
  get: () => pageContents.value[activePageId.value] ?? '',
  set: (val) => { pageContents.value[activePageId.value] = val },
})

const activeTitle    = computed(() => EDITOR_PAGES[activePageId.value]?.title ?? '')
const totalPageCount = computed(() => EDITOR_TREE.reduce((n, s) => n + s.pages.length, 0))

// ── Mock 上下文（供 VisualEditorToolbar / VisualEditorContent 的 useEditorInject 消费）──
provide(EDITOR_KEY, {
  uploading: ref(false),
  uploadImage: async (_file: File) => '',
  content,
  page: ref(null),
  showSettingsPanel: ref(false),
} as any)

// ── 可视化编辑器（Tiptap） ─────────────────────────────────────
const { editor, isReady } = useVisualEditor({ content, placeholder: '开始编写内容...' })

// ── 源码模式：预览状态 ─────────────────────────────────────────
// useMarkdownPreview 负责内容编辑时的防抖渲染（300ms）
// watch(activePageId) 负责切换页面时立即渲染
const showPreview = ref(true)
const { html: previewHtml } = useMarkdownPreview(content, { delay: 300 })

watch(activePageId, () => {
  renderMarkdown(content.value).then(html => { previewHtml.value = html })
})

// ── 文档树折叠状态 ──────────────────────────────────────────────
const sectionCollapsed = ref<Record<string, boolean>>({})
const isOpen      = (id: string) => !sectionCollapsed.value[id]
const toggleSection = (id: string) => { sectionCollapsed.value[id] = !sectionCollapsed.value[id] }

// ── 保存状态演示（仅视觉） ─────────────────────────────────────
const saving = ref(false)
function demoSave() {
  if (saving.value) return
  saving.value = true
  setTimeout(() => { saving.value = false }, 800)
}

// ── onMounted：初始渲染 + IntersectionObserver ──────────────────
// IntersectionObserver 背景：reka-ui 任何 Overlay 打开时会对 Portal 兄弟节点
// 设置 aria-hidden，当用户滚动到其他 section 后那里的弹层打开，hero 区域可能
// 被标记 aria-hidden，而 Tiptap 若仍持有 focus 就会触发浏览器警告。
// 用 watchEffect + onCleanup 统一管理 IO 生命周期，无需 onBeforeUnmount。
const demoRoot = ref<HTMLElement | null>(null)

onMounted(() => {
  // 初始 source 预览
  renderMarkdown(content.value).then(html => { previewHtml.value = html })

  // 编辑器就绪后挂载 IntersectionObserver
  watchEffect((onCleanup) => {
    if (!isReady.value || !demoRoot.value) return
    const io = new IntersectionObserver(
      ([entry]) => {
        if (entry && !entry.isIntersecting && editor.value?.isFocused)
          editor.value.commands.blur()
      },
      { threshold: 0.1 },
    )
    io.observe(demoRoot.value)
    onCleanup(() => io.disconnect())
  })
})
</script>

<template>
  <div ref="demoRoot" class="flex flex-col h-full bg-background overflow-hidden">

    <!-- 顶部工具栏 -->
    <div class="shrink-0 h-11 flex items-center gap-2 border-b px-3 bg-background z-10">

      <!-- 左侧：页面标题 + 状态徽章 -->
      <div class="flex items-center gap-2 min-w-0 flex-1 overflow-hidden">
        <component
          :is="editorMode === 'visual' ? PenLine : Code"
          class="h-4 w-4 text-muted-foreground/60 shrink-0"
        />
        <span class="text-sm font-medium text-foreground/80 truncate select-none">{{ activeTitle }}</span>
        <span class="shrink-0 inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-[10px] font-medium text-gray-600 dark:bg-gray-800 dark:text-gray-400 select-none">
          草稿
        </span>
      </div>

      <!-- 右侧：操作按钮组 -->
      <div class="flex items-center gap-1 shrink-0">

        <!-- 保存 -->
        <button
          class="h-7 w-7 flex items-center justify-center rounded-full border transition-colors"
          :class="saving ? 'bg-primary text-primary-foreground border-primary' : 'bg-background hover:bg-accent'"
          title="保存"
          @click="demoSave"
        >
          <Check class="h-3.5 w-3.5" />
        </button>

        <!-- 导入 -->
        <button class="h-7 w-7 flex items-center justify-center rounded-full border bg-background hover:bg-accent transition-colors" title="导入 Markdown">
          <Upload class="h-3.5 w-3.5" />
        </button>

        <!-- 设置 -->
        <button class="h-7 w-7 flex items-center justify-center rounded-full border bg-background hover:bg-accent transition-colors" title="页面设置">
          <Settings class="h-3.5 w-3.5" />
        </button>

        <div class="h-4 w-px bg-border mx-0.5" />

        <!-- 编辑模式切换（数据驱动） -->
        <div class="flex items-center rounded-full border bg-muted p-0.5 gap-0.5">
          <button
            v-for="mode in MODES"
            :key="mode.id"
            class="h-6 flex items-center gap-1 px-2 rounded-full text-[11px] font-medium transition-all select-none"
            :class="editorMode === mode.id
              ? 'bg-background shadow-sm text-foreground'
              : 'text-muted-foreground hover:text-foreground'"
            @click="editorMode = mode.id"
          >
            <component :is="mode.icon" class="h-3 w-3" />
            {{ mode.label }}
          </button>
        </div>

        <!-- 预览切换（仅源码模式） -->
        <Transition
          enter-active-class="transition-all duration-150"
          enter-from-class="opacity-0 scale-90"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition-all duration-100"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-90"
        >
          <button
            v-if="editorMode === 'source'"
            class="h-7 flex items-center gap-1 px-2.5 rounded-full border text-xs font-medium transition-colors select-none"
            :class="showPreview
              ? 'bg-primary text-primary-foreground border-primary hover:bg-primary/90'
              : 'bg-background hover:bg-accent'"
            @click="showPreview = !showPreview"
          >
            <component :is="showPreview ? Eye : EyeOff" class="h-3 w-3" />
            预览
          </button>
        </Transition>

        <!-- 发布 -->
        <button
          class="h-7 flex items-center gap-1 px-2.5 rounded-full bg-primary text-primary-foreground text-xs font-medium transition-opacity hover:opacity-90 select-none"
          title="发布文档"
        >
          <Rocket class="h-3 w-3" />
          发布
        </button>

      </div>
    </div>

    <!-- 编辑区主体 -->
    <div class="flex flex-1 min-h-0 overflow-hidden">

      <!-- 编辑器内容区 -->
      <div class="flex flex-col flex-1 min-w-0 min-h-0 overflow-hidden">

        <!-- 可视化编辑模式 -->
        <template v-if="editorMode === 'visual'">
          <VisualEditorToolbar :editor="editor" />
          <div class="relative flex flex-col flex-1 min-h-0">
            <VisualEditorContent :editor="editor" />
            <div
              v-if="!isReady"
              class="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm z-20"
            >
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                编辑器加载中...
              </div>
            </div>
          </div>
        </template>

        <!-- 源码编辑模式 -->
        <template v-else>
          <div class="flex flex-1 min-w-0 min-h-0 overflow-hidden">
            <textarea
              :value="content"
              placeholder="在此编写 Markdown 内容..."
              class="flex-1 min-w-0 h-full resize-none border-0 bg-background px-5 py-4 font-mono text-sm leading-relaxed outline-none text-foreground/90 placeholder:text-muted-foreground/40"
              spellcheck="false"
              @input="content = ($event.target as HTMLTextAreaElement).value"
            />
            <Transition
              enter-active-class="transition-all duration-200 ease-out"
              enter-from-class="opacity-0 translate-x-2"
              enter-to-class="opacity-100 translate-x-0"
              leave-active-class="transition-all duration-150 ease-in"
              leave-from-class="opacity-100"
              leave-to-class="opacity-0"
            >
              <div v-if="showPreview" class="w-72 shrink-0 border-l overflow-y-auto bg-card">
                <div class="px-5 py-4">
                  <div :class="[PROSE_CLASSES, 'prose-sm']" v-html="previewHtml" />
                  <p v-if="!previewHtml" class="text-center text-xs text-muted-foreground py-10 select-none">
                    预览区域
                  </p>
                </div>
              </div>
            </Transition>
          </div>
        </template>

      </div>

      <!-- 分隔线 -->
      <div class="w-px bg-border shrink-0" />

      <!-- 右侧文档树侧边栏 -->
      <aside class="w-44 shrink-0 flex flex-col overflow-hidden bg-background">

        <div class="flex items-center justify-between px-3 pt-3 pb-1.5 shrink-0">
          <span class="text-xs text-muted-foreground select-none">{{ totalPageCount }} 篇文档</span>
          <button class="h-6 w-6 flex items-center justify-center rounded-full hover:bg-accent transition-colors" title="添加章节">
            <Plus class="h-3.5 w-3.5 text-muted-foreground" />
          </button>
        </div>

        <div class="flex-1 overflow-y-auto px-2 pb-3 space-y-0.5">
          <div v-for="section in EDITOR_TREE" :key="section.id">

            <button
              class="w-full flex items-center gap-1.5 rounded-lg px-1.5 py-1.5 hover:bg-accent/50 transition-colors"
              @click="toggleSection(section.id)"
            >
              <ChevronRight
                class="h-3 w-3 text-muted-foreground transition-transform duration-200 shrink-0"
                :class="isOpen(section.id) ? 'rotate-90' : ''"
              />
              <FolderOpen class="h-3.5 w-3.5 text-muted-foreground shrink-0" />
              <span class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground truncate flex-1 text-left select-none">
                {{ section.title }}
              </span>
            </button>

            <ul v-if="isOpen(section.id)" class="ml-4 border-l pl-2 space-y-0.5 py-0.5">
              <li v-for="page in section.pages" :key="page.id">
                <button
                  class="w-full flex items-center gap-1.5 rounded-md px-2 py-1.5 text-left transition-colors hover:bg-accent"
                  :class="page.id === activePageId ? 'bg-accent font-medium text-primary' : 'text-muted-foreground'"
                  @click="activePageId = page.id"
                >
                  <FileText class="h-3 w-3 shrink-0 text-muted-foreground" />
                  <span class="truncate text-xs select-none">{{ page.title }}</span>
                </button>
              </li>
            </ul>

          </div>
        </div>

      </aside>
    </div>

  </div>
</template>
