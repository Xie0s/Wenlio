<!-- components/DocPage/DocMobileToolbar.vue
     职责：移动端统一汉堡菜单（通过 Teleport 插入顶部 header）
     设计：单个 Menu 按钮 → 底部 Sheet，包含快捷操作行（搜索/主题/首页）+ 章节导航/本页目录 Tabs
     对外暴露事件：navigate(slug), version-change(versionName), scroll-to(id) -->

<script setup lang="ts">
import { ref, reactive, inject } from 'vue'
import { Menu, Search, Home, ChevronDown } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetDescription, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import type { SectionTree, Version } from '@/utils/types'
import type { TocItem } from '@/composables/useToc'

defineProps<{
  tree: SectionTree[]
  versions: Version[]
  versionName: string
  currentSlug?: string
  currentPageTitle: string
  tocItems: TocItem[]
  activeId: string
}>()

const emit = defineEmits<{
  navigate: [slug: string]
  'version-change': [versionName: string]
  'scroll-to': [id: string]
}>()

const openSearch = inject<() => void>('openSearch', () => {})
const goTenantHome = inject<() => void>('goTenantHome', () => {})

const menuOpen = ref(false)
const activeTab = ref<'chapters' | 'toc'>('chapters')
const sectionCollapsed = reactive<Record<string, boolean>>({})

function toggleSection(id: string) {
  sectionCollapsed[id] = !sectionCollapsed[id]
}

function navigate(slug: string) {
  menuOpen.value = false
  emit('navigate', slug)
}

function scrollTo(id: string) {
  menuOpen.value = false
  emit('scroll-to', id)
}

function handleSearch() {
  menuOpen.value = false
  openSearch()
}

function handleHome() {
  menuOpen.value = false
  goTenantHome()
}
</script>

<template>
  <Teleport defer to="#reader-mobile-actions">
    <Sheet v-model:open="menuOpen">
      <SheetTrigger as-child>
        <Button variant="ghost" size="icon" class="rounded-full h-9 w-9">
          <Menu class="h-5 w-5" />
        </Button>
      </SheetTrigger>

      <SheetContent side="bottom" class="rounded-t-3xl p-0" style="max-height: 82vh">
        <SheetTitle class="sr-only">导航菜单</SheetTitle>
        <SheetDescription class="sr-only">文档导航、目录与快捷操作</SheetDescription>

        <!-- 拖拽指示条 -->
        <div class="flex justify-center pt-3 pb-2">
          <div class="h-1.5 w-10 rounded-full bg-muted-foreground/30" />
        </div>

        <!-- 快捷操作行 -->
        <div class="flex items-center justify-around border-b border-border px-6 pb-4 pt-1">
          <button
            class="flex flex-col items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors"
            @click="handleSearch"
          >
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-accent/60">
              <Search class="h-5 w-5" />
            </div>
            <span class="text-xs">搜索</span>
          </button>

          <div class="flex flex-col items-center gap-1.5">
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-accent/60">
              <ThemeToggle />
            </div>
            <span class="text-xs text-muted-foreground">主题</span>
          </div>

          <button
            class="flex flex-col items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors"
            @click="handleHome"
          >
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-accent/60">
              <Home class="h-5 w-5" />
            </div>
            <span class="text-xs">首页</span>
          </button>
        </div>

        <!-- Tab 切换（有 TOC 时才显示） -->
        <div v-if="tocItems.length > 0" class="flex border-b border-border">
          <button
            class="flex-1 py-2.5 text-sm font-medium transition-colors"
            :class="activeTab === 'chapters'
              ? 'border-b-2 border-primary text-primary'
              : 'text-muted-foreground hover:text-foreground'"
            @click="activeTab = 'chapters'"
          >
            章节导航
          </button>
          <button
            class="flex-1 py-2.5 text-sm font-medium transition-colors"
            :class="activeTab === 'toc'
              ? 'border-b-2 border-primary text-primary'
              : 'text-muted-foreground hover:text-foreground'"
            @click="activeTab = 'toc'"
          >
            本页目录
          </button>
        </div>

        <!-- 内容区 -->
        <div class="overflow-y-auto px-5 pb-8 pt-3" style="max-height: calc(82vh - 10rem)">

          <!-- 章节导航 -->
          <div v-show="activeTab === 'chapters'">
            <!-- 版本选择器（原生 select，触发设备自带选择器） -->
            <div v-if="versions.length > 1" class="mb-5">
              <select
                class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm text-foreground cursor-pointer outline-none"
                :value="versionName"
                @change="(e: any) => { menuOpen = false; emit('version-change', e.target.value) }"
              >
                <option v-for="v in versions" :key="v.id" :value="v.name">
                  {{ v.label || v.name }}
                </option>
              </select>
            </div>

            <!-- 导航分组（可折叠） -->
            <nav class="space-y-5">
              <div v-for="section in tree" :key="section.id">
                <button
                  class="mb-2 w-full flex items-center justify-between text-sm font-semibold uppercase tracking-wider text-muted-foreground/80 hover:text-foreground transition-colors select-none"
                  @click="toggleSection(section.id)"
                >
                  <span>{{ section.title }}</span>
                  <ChevronDown
                    class="size-3.5 shrink-0 transition-transform duration-200"
                    :class="sectionCollapsed[section.id] ? '-rotate-90' : ''"
                  />
                </button>
                <ul v-if="!sectionCollapsed[section.id]" class="space-y-0.5">
                  <li v-for="p in section.pages" :key="p.id" class="relative">
                    <span
                      v-if="p.slug === currentSlug"
                      class="absolute left-0 inset-y-0 my-1 w-0.5 rounded-full bg-primary"
                    />
                    <button
                      class="w-full text-left py-1.5 pl-4 pr-3 text-base rounded-md transition-colors"
                      :class="p.slug === currentSlug
                        ? 'text-primary font-medium'
                        : 'text-muted-foreground hover:text-foreground hover:bg-accent/60'"
                      @click="navigate(p.slug)"
                    >
                      {{ p.title }}
                    </button>
                  </li>
                </ul>
              </div>
            </nav>
          </div>

          <!-- 本页目录 -->
          <div v-show="activeTab === 'toc'">
            <nav class="space-y-0.5">
              <button
                v-for="item in tocItems"
                :key="item.id"
                class="block w-full text-left text-sm py-1.5 rounded transition-colors leading-relaxed"
                :class="[
                  item.level === 3 ? 'pl-4' : item.level >= 4 ? 'pl-8' : 'pl-0',
                  activeId === item.id
                    ? 'text-primary font-medium'
                    : 'text-muted-foreground hover:text-foreground',
                ]"
                @click="scrollTo(item.id)"
              >
                {{ item.text }}
              </button>
            </nav>
          </div>

        </div>
      </SheetContent>
    </Sheet>
  </Teleport>
</template>
