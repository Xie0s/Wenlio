<!-- components/DocPage/DocSidebar.vue
     职责：文档阅读页左侧导航侧边栏（桌面端，lg+ 可见）
     设计风格参考 Nextra：无边框分隔线、左边框激活指示器、分组标题 uppercase
     对外暴露事件：navigate(slug), version-change(versionName) -->

<script setup lang="ts">
import { ref, reactive, computed, inject } from 'vue'
import type { Ref } from 'vue'
import { useRoute } from 'vue-router'
import { PanelLeftClose, PanelLeftOpen, ChevronDown, ChartNoAxesColumn, Search, Leaf, ExternalLink, Copy, FolderOpen } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import type { SectionTree, Version } from '@/utils/types'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useReaderStore } from '@/stores/reader'
import { toast } from 'vue-sonner'

defineProps<{
  tree: SectionTree[]
  versions: Version[]
  versionName: string
  currentSlug: string
  tenantId?: string
  themeSlug?: string
}>()

const emit = defineEmits<{
  navigate: [slug: string]
  'version-change': [versionName: string]
}>()

const store = useReaderStore()
const collapsed = ref(false)

const sectionCollapsed = reactive<Record<string, boolean>>({})

const route = useRoute()
const searchKeyword = computed(() => (route.query.q as string) || '')

const openSearch = inject<(() => void) | undefined>('openSearch', undefined)
const eyeCareMode = inject<Ref<boolean>>('eyeCareMode', ref(false))
const toggleEyeCare = inject<() => void>('toggleEyeCare', () => { })

function toggleSection(id: string) {
  sectionCollapsed[id] = !sectionCollapsed[id]
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function highlight(text: string, keyword: string): string {
  if (!keyword.trim()) return escapeHtml(text)
  const escaped = escapeHtml(text)
  const kw = keyword.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return escaped.replace(
    new RegExp(kw, 'gi'),
    match => `<mark class="bg-yellow-200 dark:bg-yellow-900/50 text-yellow-900 dark:text-yellow-100 px-0.5 rounded">${match}</mark>`,
  )
}

// 解码 JWT payload 获取 exp 字段（不验证签名，仅读取）
function decodeJwtExp(token: string): number | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3 || !parts[1]) return null
    const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')))
    return typeof payload.exp === 'number' ? payload.exp : null
  } catch { return null }
}

function formatExpiry(exp: number): string {
  const d = new Date(exp * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 当前 raw 链接携带的 token 过期时间描述
const rawTokenExpiry = computed(() => {
  const theme = store.currentTheme
  if (theme?.access_mode === 'code' || theme?.access_mode === 'login') {
    const token = localStorage.getItem(`theme_access_${theme.id}`)
    if (token) {
      const exp = decodeJwtExp(token)
      if (exp) return formatExpiry(exp)
    }
  }
  return ''
})

// Raw Markdown 链接构建（受保护主题统一使用 theme_access token，避免泄露完整 JWT）
function getRawTokenSuffix(): string {
  const theme = store.currentTheme
  if (theme?.access_mode === 'code' || theme?.access_mode === 'login') {
    const token = localStorage.getItem(`theme_access_${theme.id}`)
    if (token) return `?theme_token=${token}`
  }
  return ''
}

function buildRawPageUrl(props: { tenantId?: string; themeSlug?: string; versionName: string; currentSlug: string }): string {
  return `${window.location.origin}/raw/${props.tenantId}/${props.themeSlug}/${props.versionName}/${props.currentSlug}${getRawTokenSuffix()}`
}

function buildRawDirUrl(props: { tenantId?: string; themeSlug?: string; versionName: string }): string {
  return `${window.location.origin}/raw/${props.tenantId}/${props.themeSlug}/${props.versionName}${getRawTokenSuffix()}`
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    const expiry = rawTokenExpiry.value
    toast.success(expiry ? `链接已复制（有效期至 ${expiry}）` : '链接已复制')
  } catch {
    toast.error('复制失败')
  }
}
</script>

<template>
  <aside class="hidden shrink-0 lg:flex lg:ml-4 flex-col transition-all duration-200 overflow-hidden"
    :class="collapsed ? 'w-10' : 'w-64'" style="height: calc(100vh - 4rem); position: sticky; top: 0;">
    <!-- 导航内容区（展开时显示） -->
    <div v-if="!collapsed" class="flex-1 min-h-0 overflow-y-auto py-8 pr-1 pl-1">

      <!-- 版本选择器 -->
      <div v-if="versions.length > 1" class="mb-6 px-2">
        <Select :model-value="versionName" @update:model-value="(v) => v && emit('version-change', v as string)">
          <SelectTrigger class="w-full h-8 text-sm rounded-lg">
            <SelectValue />
          </SelectTrigger>
          <SelectContent class="rounded-xl">
            <SelectItem v-for="v in versions" :key="v.id" :value="v.name">
              {{ v.label || v.name }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>

      <!-- 导航分组 -->
      <nav class="space-y-6">
        <div v-for="section in tree" :key="section.id">
          <!-- 分组标题 -->
          <button
            class="group mb-2 w-full flex items-center justify-between px-2 text-xl font-light text-foreground transition-colors select-none"
            @click="toggleSection(section.id)">
            <span>{{ section.title }}</span>
            <span class="flex items-center justify-center h-5 w-5 rounded-full transition-colors"
              :class="sectionCollapsed[section.id] ? 'group-hover:bg-accent' : 'bg-accent'">
              <ChevronDown class="size-3 shrink-0 transition-transform duration-200"
                :class="sectionCollapsed[section.id] ? '-rotate-90' : ''" />
            </span>
          </button>
          <!-- 页面列表（VS Code 风格：左侧竖线 + 滑轨激活） -->
          <ul v-if="!sectionCollapsed[section.id]" class="relative pl-4 space-y-0.5">
            <div class="absolute left-[9px] top-0 bottom-0 w-px bg-border/60" />
            <!-- li 向左 7px 对齐竖线，hover/激活背景整体作用，跟上滑轨 -->
            <li v-for="p in section.pages" :key="p.id" class="relative -ml-[7px] transition-colors duration-150"
              :class="p.slug === currentSlug ? 'bg-accent/40' : 'hover:bg-accent/60'">
              <!-- 激活指示器：left-0 = li 左边 = 竖线位置 -->
              <span v-if="p.slug === currentSlug" class="absolute left-0 inset-y-0 w-px bg-primary" />
              <button class="w-full text-left py-1.5 pl-[16px] pr-3 text-lg font-thin transition-colors duration-150"
                :class="p.slug === currentSlug
                  ? 'text-primary font-extralight'
                  : 'text-foreground hover:text-primary/80'" @click="emit('navigate', p.slug)"
                v-html="highlight(p.title, searchKeyword)" />
            </li>
          </ul>
        </div>
      </nav>
    </div>

    <!-- 折叠态占位弹性区 -->
    <div v-else class="flex-1" />

    <!-- 底部工具栏 -->
    <div class="shrink-0 border-t py-4 flex items-center justify-center gap-3"
      :class="collapsed ? 'flex-col' : 'flex-row'">
      <Tooltip v-if="openSearch">
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon-xl" class="rounded-full"
            @click="openSearch">
            <Search class="size-6" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">搜索文档</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon-xl" class="rounded-full"
            :class="eyeCareMode ? 'text-emerald-600 bg-emerald-50 dark:text-emerald-400 dark:bg-emerald-950/30' : ''"
            @click="toggleEyeCare">
            <Leaf class="size-5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">{{ eyeCareMode ? '关闭护眼' : '护眼模式' }}</TooltipContent>
      </Tooltip>
      <ThemeToggle button-size="size-11" icon-size="size-6" />
      <DropdownMenu v-if="currentSlug && tenantId && themeSlug">
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon-xl" class="rounded-full data-[state=open]:bg-accent">
            <Tooltip>
              <TooltipTrigger as-child>
                <ChartNoAxesColumn class="size-6" :stroke-width="1.5" />
              </TooltipTrigger>
              <TooltipContent side="top">原始 Markdown</TooltipContent>
            </Tooltip>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent side="right" align="end" :side-offset="8" class="rounded-xl w-52">
          <DropdownMenuItem as="a" :href="buildRawPageUrl({ tenantId, themeSlug, versionName, currentSlug })" target="_blank" rel="noopener noreferrer">
            <ExternalLink class="size-4 mr-2" />打开当前页 Raw
          </DropdownMenuItem>
          <DropdownMenuItem @click="copyToClipboard(buildRawPageUrl({ tenantId, themeSlug, versionName, currentSlug }))">
            <Copy class="size-4 mr-2" />复制当前页 Raw 链接
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem as="a" :href="buildRawDirUrl({ tenantId, themeSlug, versionName })" target="_blank" rel="noopener noreferrer">
            <FolderOpen class="size-4 mr-2" />打开目录
          </DropdownMenuItem>
          <DropdownMenuItem @click="copyToClipboard(buildRawDirUrl({ tenantId, themeSlug, versionName }))">
            <Copy class="size-4 mr-2" />复制目录链接
          </DropdownMenuItem>
          <template v-if="rawTokenExpiry">
            <DropdownMenuSeparator />
            <div class="px-2 py-1.5 text-[11px] text-muted-foreground leading-tight">
              链接有效期至 {{ rawTokenExpiry }}
            </div>
          </template>
        </DropdownMenuContent>
      </DropdownMenu>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon-xl" class="rounded-full"
            @click="collapsed = !collapsed">
            <PanelLeftClose v-if="!collapsed" class="size-6" :stroke-width="1.5" />
            <PanelLeftOpen v-else class="size-6" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent side="top">{{ collapsed ? '展开导航' : '折叠导航' }}</TooltipContent>
      </Tooltip>
    </div>
  </aside>
</template>
