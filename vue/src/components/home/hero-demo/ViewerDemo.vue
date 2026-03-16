<!--
  hero-demo/ViewerDemo.vue
  职责：演示文档阅读场景，自包含模拟数据，无 store/router 依赖
  视觉复刻 DocReaderLayout + DocPageViewer，使用真实 DocContent 渲染 Markdown
-->
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { renderMarkdown } from '@/lib/markdown'
import DocContent from '@/components/preview/DocContent.vue'
import { VIEWER_TREE, VIEWER_PAGES, VIEWER_INITIAL_PAGE_ID } from './demo-data'
import { Search, ChevronRight, ChevronDown, BookOpen } from 'lucide-vue-next'

// ── 页面切换状态 ──────────────────────────────────────────────
const activePageId = ref(VIEWER_INITIAL_PAGE_ID)
const htmlCache = ref<Record<string, string>>({})
const renderedHtml = ref('')
const loading = ref(true)

const activePage = computed(() => VIEWER_PAGES[activePageId.value])
const activeTitle = computed(() => activePage.value?.title ?? '')

// 根据激活页面找到对应的 section title
const activeSectionTitle = computed(() => {
  for (const s of VIEWER_TREE) {
    if (s.pages.some(p => p.id === activePageId.value)) return s.title
  }
  return ''
})

// ── Markdown 渲染（带缓存） ────────────────────────────────────
async function loadPage(id: string) {
  loading.value = true
  if (htmlCache.value[id]) {
    renderedHtml.value = htmlCache.value[id]
    loading.value = false
    return
  }
  const page = VIEWER_PAGES[id]
  if (page) {
    const html = await renderMarkdown(page.markdown)
    htmlCache.value[id] = html
    renderedHtml.value = html
  }
  loading.value = false
}

onMounted(() => loadPage(activePageId.value))
watch(activePageId, loadPage)

// ── TOC（从 markdown 提取标题） ────────────────────────────────
function extractToc(markdown: string) {
  return markdown
    .split('\n')
    .flatMap(line => {
      const m = line.match(/^(#{1,3})\s+(.+)$/)
      if (!m) return []
      const level = m[1]!.length
      const text = m[2]!.replace(/\*\*/g, '').trim()
      const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fff-]/g, '')
      return [{ id, text, level }]
    })
}

const tocItems = computed(() =>
  activePage.value ? extractToc(activePage.value.markdown) : [],
)
const activeTocId = ref('')
watch(tocItems, items => { activeTocId.value = items[1]?.id ?? items[0]?.id ?? '' }, { immediate: true })

// ── 章节折叠 ─────────────────────────────────────────────────
const collapsed = ref<Record<string, boolean>>({})
function toggleSection(id: string) {
  collapsed.value[id] = !collapsed.value[id]
}
function isOpen(id: string) {
  return collapsed.value[id] !== true
}
</script>

<template>
  <div class="flex flex-col h-full bg-background overflow-hidden">

    <!-- 顶部导航栏 -->
    <header class="shrink-0 border-b bg-background z-10">
      <div class="flex h-12 items-center gap-2 px-4">

        <!-- 面包屑 -->
        <div class="flex items-center gap-1 min-w-0 flex-1 overflow-hidden">
          <BookOpen class="h-4 w-4 text-primary/70 shrink-0" />
          <span class="text-sm font-medium text-foreground/80 shrink-0 select-none">Wenlio 文流</span>
          <ChevronRight class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
          <span class="text-sm text-muted-foreground shrink-0 select-none">用户手册</span>
          <ChevronRight class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
          <span class="text-sm text-muted-foreground shrink-0 select-none truncate">{{ activeSectionTitle }}</span>
          <ChevronRight class="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
          <span class="text-sm font-medium text-foreground shrink-0 select-none truncate">{{ activeTitle }}</span>
        </div>

        <!-- 版本标签 -->
        <span class="shrink-0 text-[11px] text-muted-foreground/70 select-none hidden sm:block">v1.0</span>

        <!-- 搜索框（装饰性） -->
        <button
          class="shrink-0 flex items-center gap-1.5 h-7 px-2.5 rounded-full border bg-muted/40 hover:bg-muted transition-colors text-xs text-muted-foreground w-28"
        >
          <Search class="h-3 w-3 shrink-0" />
          <span>搜索文档...</span>
        </button>

      </div>
    </header>

    <!-- 主体区域 -->
    <div class="flex flex-1 min-h-0 overflow-hidden">

      <!-- 左侧导航侧边栏 -->
      <aside class="w-44 shrink-0 border-r overflow-y-auto bg-muted/20">
        <nav class="py-5 px-2 space-y-5">
          <div v-for="section in VIEWER_TREE" :key="section.id">

            <!-- 章节标题（可折叠） -->
            <button
              class="w-full flex items-center justify-between px-2 mb-1.5 text-left transition-colors"
              @click="toggleSection(section.id)"
            >
              <span class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground/60 select-none">
                {{ section.title }}
              </span>
              <ChevronDown
                class="h-3 w-3 text-muted-foreground/50 transition-transform duration-200 shrink-0"
                :class="!isOpen(section.id) ? '-rotate-90' : ''"
              />
            </button>

            <!-- 页面列表 -->
            <ul v-if="isOpen(section.id)" class="space-y-0.5">
              <li
                v-for="page in section.pages"
                :key="page.id"
                class="relative"
              >
                <!-- 激活指示左边框 -->
                <span
                  v-if="page.id === activePageId"
                  class="absolute left-0 inset-y-0 my-0.5 w-0.5 rounded-full bg-primary"
                />
                <button
                  class="w-full text-left py-1.5 pl-4 pr-2 text-sm rounded-md transition-colors select-none"
                  :class="page.id === activePageId
                    ? 'text-primary font-medium'
                    : 'text-muted-foreground hover:text-foreground hover:bg-accent/60'"
                  @click="activePageId = page.id"
                >
                  {{ page.title }}
                </button>
              </li>
            </ul>

          </div>
        </nav>
      </aside>

      <!-- 中间正文区 -->
      <main class="flex-1 min-w-0 overflow-y-auto">
        <div class="px-8 py-7 max-w-2xl">

          <!-- 加载中 -->
          <div v-if="loading" class="flex justify-center py-14">
            <div class="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          </div>

          <!-- 文档内容 -->
          <DocContent v-else :html="renderedHtml" />

        </div>
      </main>

      <!-- 右侧 TOC -->
      <aside class="w-40 shrink-0 border-l py-5 px-3 overflow-y-auto bg-muted/10">
        <p class="mb-3 text-[10px] font-semibold uppercase tracking-widest text-muted-foreground/50 select-none">
          本页目录
        </p>
        <nav class="space-y-0.5">
          <button
            v-for="item in tocItems"
            :key="item.id"
            class="block w-full text-left text-xs py-1.5 rounded-md transition-colors select-none"
            :class="[
              item.level === 1 ? 'px-2 font-medium' : 'pl-4 pr-2',
              item.id === activeTocId
                ? 'bg-accent/60 text-primary'
                : 'text-muted-foreground hover:bg-accent/40 hover:text-foreground',
            ]"
            @click="activeTocId = item.id"
          >
            {{ item.text }}
          </button>
        </nav>
      </aside>

    </div>
  </div>
</template>
