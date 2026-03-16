<!--
  SearchDialog.vue - 文档搜索对话框组件
  职责：提供完整的文档全文搜索功能，含防抖输入、骨架屏过渡、分页加载、关键词高亮
  对外暴露：
    Props: tenantId (string), open (v-model boolean)
    Emits: navigate(path: string)
-->
<template>
  <Dialog v-model:open="isOpen">
    <DialogContent
      class="rounded-2xl p-0 gap-0 overflow-hidden transition-all duration-200"
      :class="hasResults ? 'sm:max-w-4xl' : 'sm:max-w-2xl'"
      :show-close-button="false"
    >
      <!-- 无障碍标题（视觉隐藏） -->
      <VisuallyHidden>
        <DialogTitle>搜索文档</DialogTitle>
        <DialogDescription>输入关键词搜索文档内容</DialogDescription>
      </VisuallyHidden>

      <!-- 搜索输入区域 -->
      <div class="flex items-center border-b px-4 py-3 gap-3">
        <Search class="w-5 h-5 text-muted-foreground shrink-0" />
        <input
          ref="searchInputRef"
          v-model="keyword"
          type="text"
          placeholder="搜索文档..."
          class="flex-1 bg-transparent text-base outline-none placeholder:text-muted-foreground"
          @keydown.enter="handleSearch"
          @input="handleInputChange"
        />
        <div class="flex items-center gap-1 shrink-0">
          <!-- 加载指示 -->
          <Loader2
            v-if="isLoading"
            class="w-5 h-5 text-muted-foreground animate-spin"
          />
          <!-- 清除 / 关闭按钮 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <button
                class="p-1.5 rounded-full hover:bg-muted text-muted-foreground transition-colors"
                @click="keyword ? clearSearch() : closeDialog()"
              >
                <X class="w-5 h-5" />
              </button>
            </TooltipTrigger>
            <TooltipContent>{{ keyword ? '清除' : '关闭' }}</TooltipContent>
          </Tooltip>
        </div>
      </div>

      <!-- 搜索结果区域 -->
      <div class="max-h-[60vh] overflow-y-auto">
        <Transition name="content-fade" mode="out-in">

          <!-- 骨架屏 -->
          <div v-if="showSkeleton" key="skeleton" class="divide-y">
            <SearchResultSkeleton
              v-for="i in 5"
              :key="`sk-${i}`"
              :index="i - 1"
              class="skeleton-item"
              :class="{ 'skeleton-visible': skeletonVisible }"
              :style="{ '--stagger-delay': `${(i - 1) * 60}ms` }"
            />
          </div>

          <!-- 未搜索 -->
          <div v-else-if="!hasSearched" key="idle" class="py-12 text-center">
            <Search class="w-12 h-12 mx-auto text-muted-foreground/40 mb-3" />
            <p class="text-base text-muted-foreground">输入关键词开始搜索</p>
            <p class="text-sm text-muted-foreground/60 mt-1">按 Enter 或等待自动搜索</p>
          </div>

          <!-- 无结果 -->
          <div v-else-if="searchComplete && results.length === 0" key="empty" class="py-12 text-center">
            <FileQuestion class="w-12 h-12 mx-auto text-muted-foreground/40 mb-3" />
            <p class="text-base text-muted-foreground">未找到匹配的文档</p>
            <p class="text-sm text-muted-foreground/60 mt-1">尝试使用其他关键词</p>
          </div>

          <!-- 结果列表 -->
          <div v-else-if="results.length > 0" key="results" class="divide-y">
            <SearchResultItem
              v-for="(item, index) in results"
              :key="item.page_id"
              :item="item"
              :keyword="keyword"
              class="result-item"
              :style="{ '--result-delay': `${Math.min(index, 10) * 25}ms` }"
              @click="handleResultClick(item)"
            />
          </div>

        </Transition>

        <!-- 加载更多 -->
        <div v-if="hasMore && !isLoading && !showSkeleton" class="px-4 py-3 text-center border-t">
          <button
            class="text-sm text-muted-foreground hover:text-foreground transition-colors"
            @click="loadMore"
          >
            加载更多
          </button>
        </div>
      </div>

      <!-- 底部信息栏 -->
      <div v-if="hasSearched && results.length > 0" class="border-t px-4 py-2 bg-muted/30">
        <div class="flex items-center justify-between text-sm text-muted-foreground">
          <span>找到 {{ total }} 条结果，已加载 {{ results.length }} 条</span>
          <span>按 Enter 打开</span>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { Search, X, FileQuestion, Loader2 } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { VisuallyHidden } from 'reka-ui'
import { useDebounceFn } from '@vueuse/core'
import { http } from '@/utils/http'
import type { SearchResult, SearchResponse } from '@/utils/types'
import SearchResultItem from './SearchResultItem.vue'
import SearchResultSkeleton from './SearchResultSkeleton.vue'

const props = defineProps<{
  tenantId: string
}>()

const emit = defineEmits<{
  navigate: [path: string]
}>()

// 双向绑定 open
const isOpen = defineModel<boolean>('open', { default: false })

// 搜索状态
const keyword = ref('')
const isLoading = ref(false)
const hasSearched = ref(false)
const searchComplete = ref(false)
const showSkeleton = ref(false)
const skeletonVisible = ref(false)
const results = ref<SearchResult[]>([])
const page = ref(1)
const pageSize = 10
const total = ref(0)

const searchInputRef = ref<HTMLInputElement | null>(null)

// 已加载数 < 总数时还有更多
const hasMore = computed(() => !isLoading.value && results.value.length < total.value)
const hasResults = computed(() => hasSearched.value && results.value.length > 0)

// 防抖搜索
const debouncedSearch = useDebounceFn(() => {
  if (keyword.value.trim()) {
    executeSearch(true)
  }
}, 500)

// 输入变化：立即显示骨架屏，防抖触发请求
function handleInputChange() {
  const kw = keyword.value.trim()
  if (kw.length >= 1) {
    if (!showSkeleton.value) {
      showSkeleton.value = true
      searchComplete.value = false
      nextTick(() => { skeletonVisible.value = true })
    }
    debouncedSearch()
  } else {
    showSkeleton.value = false
    skeletonVisible.value = false
    searchComplete.value = false
    results.value = []
    hasSearched.value = false
  }
}

// 执行搜索
async function executeSearch(reset = true) {
  if (!keyword.value.trim()) return

  if (reset) {
    page.value = 1
    results.value = []
  }

  isLoading.value = true
  hasSearched.value = true
  searchComplete.value = false

  try {
    const res = await http.get<SearchResponse>('/public/search', {
      q: keyword.value.trim(),
      tenant_id: props.tenantId,
      page: page.value,
      page_size: pageSize,
    })

    if (res.code === 0 && res.data) {
      const { total: t, items } = res.data
      total.value = t ?? 0
      if (reset) {
        results.value = items ?? []
      } else {
        results.value = [...results.value, ...(items ?? [])]
      }
    } else {
      if (reset) results.value = []
      total.value = 0
    }
  } catch (err) {
    console.error('[SearchDialog] 搜索失败:', err)
    if (reset) results.value = []
    total.value = 0
  } finally {
    isLoading.value = false
    searchComplete.value = true
    showSkeleton.value = false
    skeletonVisible.value = false
  }
}

// 手动搜索（按 Enter）
function handleSearch() {
  if (keyword.value.trim()) executeSearch(true)
}

// 加载更多
function loadMore() {
  page.value++
  executeSearch(false)
}

// 清除搜索
function clearSearch() {
  keyword.value = ''
  results.value = []
  total.value = 0
  hasSearched.value = false
  searchComplete.value = false
  showSkeleton.value = false
  skeletonVisible.value = false
  searchInputRef.value?.focus()
}

// 关闭弹窗
function closeDialog() {
  isOpen.value = false
}

// 点击结果项
function handleResultClick(item: SearchResult) {
  isOpen.value = false
  const kw = keyword.value.trim()
  const path = kw ? `${item.path}?q=${encodeURIComponent(kw)}` : item.path
  emit('navigate', path)
}

// 打开时自动聚焦输入框
watch(isOpen, (open) => {
  if (open) {
    nextTick(() => { searchInputRef.value?.focus() })
  }
})
</script>

<style scoped>
/* 筛选面板滑入 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 内容区切换 */
.content-fade-enter-active {
  transition: opacity 0.2s ease;
}
.content-fade-leave-active {
  transition: opacity 0.15s ease;
}
.content-fade-enter-from,
.content-fade-leave-to {
  opacity: 0;
}

/* 骨架屏错开出现 */
.skeleton-item {
  opacity: 0;
  transform: translateY(6px);
  transition: opacity 0.2s ease, transform 0.2s ease;
  transition-delay: var(--stagger-delay, 0ms);
}
.skeleton-item.skeleton-visible {
  opacity: 1;
  transform: translateY(0);
}

/* 结果项入场 */
.result-item {
  animation: result-in 0.2s ease forwards;
  animation-delay: var(--result-delay, 0ms);
  opacity: 0;
}
@keyframes result-in {
  from { opacity: 0; transform: translateY(4px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
