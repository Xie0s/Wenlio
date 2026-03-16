<!--
  ThemeGalleryPage.vue - 主题画廊（公开分类/标签筛选）
  职责：读者端按分类/标签浏览文档主题，支持 URL 参数筛选
  对外接口：路由参数 tenantId，Query 参数 category / tag[]
-->
<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePublicGallery } from '@/lib/public-gallery'
import ThemeGalleryBrowser from '@/components/editor/theme-management/ThemeGalleryBrowser.vue'
import ThemeCard from '@/components/editor/theme-management/ThemeCard.vue'
import { BookOpen, Leaf, Loader, SearchX, X } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import UserMenu from '@/components/auth/UserMenu.vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const gallery = usePublicGallery()
const authStore = useAuthStore()

const eyeCareMode = ref(false)
function toggleEyeCare() { eyeCareMode.value = !eyeCareMode.value }

watch(eyeCareMode, (on) => {
  if (on) document.body.classList.add('eye-care')
  else document.body.classList.remove('eye-care')
})

onUnmounted(() => {
  document.body.classList.remove('eye-care')
})

function refreshData() {
  gallery.loadFilters()
  gallery.loadThemes()
}

function clearAllTags() {
  const { tag, ...rest } = route.query
  router.replace({ query: rest })
}

function goToTheme(slug: string) {
  router.push(`/${gallery.tenantId.value}/${slug}`)
}

onMounted(() => {
  gallery.loadFilters()
  gallery.loadThemes()
})
watch([gallery.selectedCategorySlug, gallery.selectedTagSlugs], gallery.loadThemes)
</script>

<template>
  <!-- h-screen flex-col：header + 双栏主体 -->
  <div class="h-screen flex flex-col overflow-hidden bg-background">

    <!-- ════════ 顶部导航栏（对齐 DocReaderLayout header）════════ -->
    <header class="shrink-0 z-40 border-b bg-background">
      <div class="flex h-14 items-center gap-1.5 px-6">

        <span class="text-sm font-medium text-foreground">文档主题</span>

        <!-- 右侧：用户菜单（已登录）-->
        <div class="ml-auto flex items-center gap-2">
          <UserMenu v-if="authStore.isLoggedIn" :tenant-id="gallery.tenantId.value" />
        </div>
      </div>
    </header>

    <!-- ════════ 双栏主体 ════════ -->
    <div class="flex-1 min-h-0 overflow-hidden">
      <div class="h-full px-6 flex gap-4">

        <!-- ════════ 左侧侧边栏（对齐 ThemeListSidebar 样式）════════ -->
        <aside class="w-72 shrink-0 h-full overflow-hidden flex flex-col pt-6">

          <!-- 分类/标签浏览器 -->
          <div class="flex-1 overflow-y-auto">
            <ThemeGalleryBrowser :categories="gallery.categories.value"
              :categories-loading="gallery.categoriesLoading.value"
              :selected-category-id="gallery.selectedCategoryId.value" :tags="gallery.tags.value"
              :tags-loading="gallery.tagsLoading.value" :selected-tag-slugs="gallery.selectedTagSlugs.value"
              @select-category="gallery.selectCategory" @toggle-tag="gallery.toggleTag" />
          </div>

          <!-- 底部工具栏 -->
          <div class="shrink-0 border-t py-4 flex items-center justify-center gap-3">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full h-11 w-11 transition-transform hover:scale-110"
                  @click="refreshData">
                  <Loader class="size-6" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">刷新数据</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full h-11 w-11 transition-transform hover:scale-110"
                  :class="eyeCareMode ? 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-950/30' : ''"
                  @click="toggleEyeCare">
                  <Leaf class="size-6" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="top">{{ eyeCareMode ? '关闭护眼' : '护眼模式' }}</TooltipContent>
            </Tooltip>
            <ThemeToggle button-size="size-11" icon-size="size-6" />
          </div>

        </aside>

        <!-- ══════════ 右侧主题列表（独立滚动，对齐 ThemeListPage）══════════ -->
        <div class="flex-1 min-w-0 overflow-y-auto pt-6 pb-24">

          <!-- 筛选状态栏 -->
          <div v-if="gallery.hasFilters.value" class="relative flex items-center justify-center gap-3 mb-4">
            <!-- 分类筛选状态 -->
            <div v-if="gallery.selectedCategoryName.value"
              class="inline-flex items-center gap-2 rounded-full border border-border/60 bg-muted/30 pl-4 pr-1.5 h-11">
              <span class="text-xs text-muted-foreground shrink-0">分类</span>
              <span class="text-sm text-foreground/80 shrink-0">{{ gallery.selectedCategoryName.value }}</span>
              <button title="清除分类筛选"
                class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-destructive/20 text-destructive hover:bg-destructive/30 transition-colors"
                @click="gallery.selectCategory(gallery.selectedCategoryId.value)">
                <X class="h-3.5 w-3.5" />
              </button>
            </div>

            <!-- 标签筛选状态（多选，合并到单个容器）-->
            <div v-if="gallery.selectedTagSlugs.value.length > 0"
              class="inline-flex items-center gap-2 rounded-full border border-border/60 bg-muted/30 pl-4 pr-1.5 h-11">
              <span class="text-xs text-muted-foreground shrink-0">标签</span>
              <div class="flex items-center gap-1.5">
                <span v-for="slug in gallery.selectedTagSlugs.value.slice(0, 5)" :key="slug"
                  class="inline-flex items-center gap-1 rounded-full bg-foreground/10 pl-2.5 pr-1 py-0.5 text-sm text-foreground/80">
                  #{{gallery.tags.value.find(t => t.slug === slug)?.name ?? slug}}
                  <button
                    class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full hover:bg-destructive/20 hover:text-destructive transition-colors"
                    @click="gallery.toggleTag(slug)">
                    <X class="h-2.5 w-2.5" />
                  </button>
                </span>
                <span v-if="gallery.selectedTagSlugs.value.length > 5"
                  class="inline-flex items-center rounded-full bg-foreground/10 px-2.5 py-0.5 text-sm text-muted-foreground shrink-0">
                  +{{ gallery.selectedTagSlugs.value.length - 5 }}
                </span>
              </div>
              <button title="清除标签筛选"
                class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-destructive/20 text-destructive hover:bg-destructive/30 transition-colors"
                @click="clearAllTags">
                <X class="h-3.5 w-3.5" />
              </button>
            </div>

            <!-- 计数固定右侧 -->
            <span v-if="!gallery.loading.value && (gallery.visibleThemes.value.length > 0 || gallery.hasFilters.value)"
              class="absolute right-0 rounded-full bg-muted px-2.5 py-0.5 text-sm tabular-nums text-muted-foreground">
              共 {{ gallery.visibleThemes.value.length }} 个主题
            </span>
          </div>

          <!-- 加载中 -->
          <div v-if="gallery.loading.value" class="flex items-center justify-center h-full">
            <div class="google-spinner" />
          </div>

          <!-- 空状态 -->
          <div v-else-if="gallery.visibleThemes.value.length === 0" class="flex min-h-full flex-col items-center justify-center px-6">
            <div class="w-full max-w-sm">
              <div class="flex flex-col items-center text-center">
                <div
                  class="mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-muted/50 border border-border/40">
                  <SearchX v-if="gallery.hasFilters.value" class="h-7 w-7 text-muted-foreground/50"
                    :stroke-width="1.5" />
                  <BookOpen v-else class="h-7 w-7 text-muted-foreground/50" :stroke-width="1.5" />
                </div>
                <h3 class="text-xl font-normal tracking-tight text-foreground">
                  {{ gallery.hasFilters.value ? '未找到匹配的主题' : '暂无公开主题' }}
                </h3>
                <p class="mt-2 text-sm font-light text-muted-foreground/70 leading-relaxed">
                  {{ gallery.hasFilters.value
                    ? '当前筛选条件下没有匹配的主题，请尝试调整或清除筛选。'
                    : '该文档空间暂时没有公开的主题可浏览。' }}
                </p>
                <div v-if="gallery.hasFilters.value" class="mt-6 flex flex-wrap justify-center gap-2">
                  <button v-if="gallery.selectedCategoryId.value"
                    class="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/30 px-4 h-9 text-sm text-muted-foreground hover:bg-muted/60 transition-colors"
                    @click="gallery.selectCategory(gallery.selectedCategoryId.value)">
                    <X class="h-3.5 w-3.5" />清除分类
                  </button>
                  <button v-if="gallery.selectedTagSlugs.value.length > 0"
                    class="inline-flex items-center gap-1.5 rounded-full border border-border/60 bg-muted/30 px-4 h-9 text-sm text-muted-foreground hover:bg-muted/60 transition-colors"
                    @click="clearAllTags">
                    <X class="h-3.5 w-3.5" />清除标签
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- 主题列表 -->
          <div v-else>
            <ThemeCard v-for="t in gallery.visibleThemes.value" :key="t.id" :theme="t" :deleting="false" :cascade-open="false"
              :cat-lib="gallery.catLib" :tag-lib="gallery.tagLib" :readonly="true" @enter="goToTheme(t.slug)" />
          </div>

        </div>
      </div>
    </div>

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
