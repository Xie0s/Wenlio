<!-- DocPageViewer.vue - 文档阅读页组件（侧边栏 + 正文 + TOC + 评论区 + 移动端抽屉）
     职责：根据传入的路由参数加载并展示文档内容、侧边栏导航、TOC、评论
     对外接口：
       Props: tenantId, themeSlug, versionName, pageSlug
       Emits: navigate(slug), version-change(version), redirect(path) -->
<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReaderStore } from '@/stores/reader'
import { renderMarkdown } from '@/lib/markdown'
import { useToc } from '@/composables/useToc'
import { X } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import DocSidebar from './DocSidebar.vue'
import DocToc from './DocToc.vue'
import DocContent from './DocContent.vue'
import DocComments from './DocComments.vue'
import DocPageNav from './DocPageNav.vue'
import DocMobileToolbar from './DocMobileToolbar.vue'

const props = defineProps<{
  tenantId: string
  themeSlug: string
  versionName: string
  pageSlug: string
}>()

const emit = defineEmits<{
  navigate: [slug: string]
  'version-change': [version: string]
  redirect: [path: string]
}>()

const store = useReaderStore()
const route = useRoute()
const router = useRouter()

// Markdown 渲染 HTML
const renderedHtml = ref('')
const renderingMd = ref(false)

// 搜索高亮
const contentRef = ref<HTMLElement | null>(null)
const searchKeyword = computed(() => (route.query.q as string) || '')

function applyHighlight(container: HTMLElement, keyword: string) {
  container.querySelectorAll('mark[data-search-hl]').forEach(el => {
    const parent = el.parentNode!
    parent.replaceChild(document.createTextNode(el.textContent ?? ''), el)
    parent.normalize()
  })
  if (!keyword.trim()) return
  const re = new RegExp(keyword.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi')
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT)
  const textNodes: Text[] = []
  let node: Node | null
  while ((node = walker.nextNode())) textNodes.push(node as Text)
  for (const textNode of textNodes) {
    const text = textNode.textContent ?? ''
    re.lastIndex = 0
    if (!re.test(text)) continue
    re.lastIndex = 0
    const frag = document.createDocumentFragment()
    let last = 0
    let m: RegExpExecArray | null
    while ((m = re.exec(text)) !== null) {
      if (m.index > last) frag.appendChild(document.createTextNode(text.slice(last, m.index)))
      const mark = document.createElement('mark')
      mark.setAttribute('data-search-hl', '')
      mark.className = 'bg-yellow-200 dark:bg-yellow-900/50 text-yellow-900 dark:text-yellow-100 px-0.5 rounded'
      mark.textContent = m[0]
      frag.appendChild(mark)
      last = m.index + m[0].length
    }
    if (last < text.length) frag.appendChild(document.createTextNode(text.slice(last)))
    textNode.parentNode!.replaceChild(frag, textNode)
  }
  container.querySelector('mark[data-search-hl]')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
}

function clearHighlight() {
  const query = { ...route.query }
  delete query.q
  router.replace({ query })
}

// TOC
const { tocItems, activeId, scrollTo, cleanup } = useToc(() => renderedHtml.value)

// 评论提交状态
const submittingComment = ref(false)
const commentsRef = ref<{ resetForm: () => void } | null>(null)

// 翻页导航：将所有页面展开为平铺列表
const allPages = computed(() => {
  const pages: { id: string; slug: string; title: string }[] = []
  for (const section of store.tree) {
    for (const page of section.pages) {
      pages.push(page)
    }
  }
  return pages
})

const currentPageIndex = computed(() =>
  allPages.value.findIndex(p => p.slug === props.pageSlug),
)

const prevPage = computed(() => {
  const idx = currentPageIndex.value
  return idx > 0 ? (allPages.value[idx - 1] ?? null) : null
})

const nextPage = computed(() => {
  const idx = currentPageIndex.value
  return idx >= 0 && idx < allPages.value.length - 1 ? (allPages.value[idx + 1] ?? null) : null
})

async function loadData() {
  try {
    // 路由守卫 / beforeEach 可能已准备好部分数据，此处按当前租户补齐所需上下文
    if (!store.currentTheme || store.currentTheme.slug !== props.themeSlug) {
      if (store.hasThemesForTenant(props.tenantId)) {
        store.findThemeBySlug(props.themeSlug)
      }
      if (!store.currentTheme || store.currentTheme.slug !== props.themeSlug) {
        await store.ensureThemesLoaded(props.tenantId)
        store.findThemeBySlug(props.themeSlug)
      }
    }
    const theme = store.currentTheme
    if (!theme) {
      toast.error('文档主题不存在')
      emit('redirect', `/${props.tenantId}`)
      return
    }

    if (!store.currentVersion || store.currentVersion.name !== props.versionName) {
      if (store.versions.length > 0) {
        store.findVersionByName(props.versionName)
      }
      if (!store.currentVersion || store.currentVersion.name !== props.versionName) {
        await store.loadVersions(theme.id)
        store.findVersionByName(props.versionName)
      }
    }
    const version = store.currentVersion
    if (!version) {
      toast.error('版本不存在')
      emit('redirect', `/${props.tenantId}/${props.themeSlug}`)
      return
    }

    if (store.tree.length === 0) {
      await store.loadTree(version.id)
    }

    if (props.pageSlug) {
      await loadPage()
    }
  } catch {
    toast.error('文档加载失败，请刷新重试')
  }
}

async function loadPage() {
  const pageId = store.findPageIdBySlug(props.pageSlug)
  if (!pageId) return

  try {
    await store.loadPage(pageId)
  } catch {
    toast.error('页面加载失败，请刷新重试')
    return
  }

  if (store.currentPage?.content) {
    renderingMd.value = true
    renderedHtml.value = await renderMarkdown(store.currentPage.content)
    renderingMd.value = false
  } else {
    renderedHtml.value = ''
  }

  store.loadComments(pageId)
}

async function handleSubmitComment(data: { name: string; email: string; content: string }) {
  if (!data.name || !data.content || !store.currentPage) {
    toast.error('请填写昵称和评论内容')
    return
  }
  submittingComment.value = true
  const result = await store.submitComment(
    store.currentPage.id,
    { name: data.name, email: data.email || undefined },
    data.content,
  )
  submittingComment.value = false
  if (result.success) {
    toast.success('评论已提交，等待审核')
    commentsRef.value?.resetForm()
  } else {
    toast.error(result.message || '提交失败')
  }
}

function onVersionChange(newVersionName: string) {
  emit('version-change', newVersionName)
}

function navigate(slug: string) {
  emit('navigate', slug)
}

function getDocBasePath() {
  return `/${encodeURIComponent(props.tenantId)}/${encodeURIComponent(props.themeSlug)}/${encodeURIComponent(props.versionName)}/`
}

function handleContentLinkClick(event: MouseEvent) {
  const target = event.target as HTMLElement | null
  if (!target) return

  const link = target.closest('a') as HTMLAnchorElement | null
  if (!link) return

  const rawHref = link.getAttribute('href')?.trim()
  if (!rawHref) return

  // 同页标题锚点：阻止默认 hash 导航，改为现有滚动逻辑，避免触发路由加载态。
  if (rawHref.startsWith('#')) {
    const id = rawHref.slice(1)
    if (!id) return
    event.preventDefault()
    scrollTo(id)
    return
  }

  let url: URL
  try {
    url = new URL(rawHref, window.location.href)
  } catch {
    return
  }

  if (url.origin !== window.location.origin) return

  const basePath = getDocBasePath()
  if (!url.pathname.startsWith(basePath)) return

  const slug = decodeURIComponent(url.pathname.slice(basePath.length))
  if (!slug) return

  event.preventDefault()

  // 同页锚点（完整路径写法）：不切页，仅滚动。
  if (slug === props.pageSlug && url.hash) {
    const id = decodeURIComponent(url.hash.slice(1))
    if (id) scrollTo(id)
    return
  }

  emit('navigate', slug)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && searchKeyword.value) clearHighlight()
}

onMounted(() => {
  loadData()
  window.addEventListener('keydown', onKeydown)
})
onUnmounted(() => {
  cleanup()
  window.removeEventListener('keydown', onKeydown)
})

watch(() => props.pageSlug, async () => {
  if (props.pageSlug) {
    await loadPage()
    // 切换页面后重置滚动至顶部
    nextTick(() => {
      document.getElementById('reader-scroll-container')?.scrollTo({ top: 0, behavior: 'instant' })
    })
  }
})

watch(renderedHtml, async () => {
  await nextTick()
  if (contentRef.value && searchKeyword.value) applyHighlight(contentRef.value, searchKeyword.value)
})

watch(searchKeyword, async (kw) => {
  await nextTick()
  if (contentRef.value) applyHighlight(contentRef.value, kw)
})

watch(() => props.versionName, () => {
  store.clearCache()
  renderedHtml.value = ''
  loadData()
})
</script>

<template>
  <div class="mx-auto flex w-full min-h-full items-start xl:gap-0">

    <!-- 左侧导航（桌面端，lg+） -->
    <DocSidebar :tree="store.tree" :versions="store.versions" :version-name="versionName" :current-slug="pageSlug"
      :tenant-id="tenantId" :theme-slug="themeSlug" @navigate="navigate" @version-change="onVersionChange" />

    <!-- 中间内容区 -->
    <div class="flex-1 min-w-0 flex flex-col self-stretch">

      <!-- 移动端章节/TOC 按钮（通过 Teleport 注入顶部 header，此处不渲染 UI） -->
      <DocMobileToolbar :tree="store.tree" :versions="store.versions" :version-name="versionName"
        :current-slug="pageSlug" :current-page-title="store.currentPage?.title || '文档'" :toc-items="tocItems"
        :active-id="activeId" @navigate="navigate" @version-change="onVersionChange" @scroll-to="scrollTo" />

      <!-- 搜索高亮提示条 -->
      <Transition enter-active-class="transition-all duration-200" enter-from-class="opacity-0 -translate-y-1"
        enter-to-class="opacity-100 translate-y-0" leave-active-class="transition-all duration-150"
        leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 -translate-y-1">
        <div v-if="searchKeyword"
          class="flex items-center gap-2 px-4 py-1.5 bg-yellow-50 dark:bg-yellow-900/20 border-b border-yellow-200 dark:border-yellow-800 text-xs text-yellow-800 dark:text-yellow-200">
          <span class="flex-1">正在高亮关键词：<strong>{{ searchKeyword }}</strong></span>
          <button class="p-0.5 rounded hover:bg-yellow-200 dark:hover:bg-yellow-800/50 transition-colors"
            @click="clearHighlight">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>
      </Transition>

      <!-- 正文滚动区 -->
      <div class="flex-1 flex flex-col px-4 py-8 lg:px-6 xl:px-8 2xl:px-10 lg:py-10">

        <!-- 加载中 -->
        <div v-if="store.loadingPage || renderingMd" class="flex flex-1 items-center justify-center">
          <div class="google-spinner" />
        </div>

        <!-- 有内容 -->
        <template v-else-if="store.currentPage">
          <div class="w-full">
            <div ref="contentRef" @click="handleContentLinkClick">
              <DocContent :html="renderedHtml" />
            </div>

            <DocPageNav :prev-page="prevPage" :next-page="nextPage" @navigate="navigate" />

            <DocComments ref="commentsRef" :comments="store.comments" :submitting="submittingComment"
              @submit-comment="handleSubmitComment" />
          </div>
        </template>

        <!-- 空状态 -->
        <div v-else class="flex flex-1 items-center justify-center">
          <p class="text-sm text-muted-foreground">从左侧选择一篇文档开始阅读</p>
        </div>

      </div>
    </div>

    <!-- 右侧 TOC（桌面端，xl+） -->
    <DocToc v-if="tocItems.length > 0" :items="tocItems" :active-id="activeId" @scroll-to="scrollTo" />

  </div>
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
