<!-- NavbarEditor.vue - 导航栏区块配置编辑器
     职责：编辑 NavbarConfig 的所有字段，支持二级链接和内部链接选择
     对外接口：
       Props: config, themes?, tenantId?
       Emits: update(config: NavbarConfig) -->
<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { NavbarConfig, NavLink } from '../types'
import type { Theme } from '@/utils/types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Plus, Trash2, ChevronRight, Link2 } from 'lucide-vue-next'
import LinkInput from './LinkInput.vue'

const props = defineProps<{
  config: NavbarConfig
  themes?: Theme[]
  tenantId?: string
}>()
const emit = defineEmits<{ update: [config: NavbarConfig] }>()

const form = reactive<NavbarConfig>(JSON.parse(JSON.stringify(props.config)))

// 记录每个一级链接是否展开二级链接
const expandedLinks = ref<Set<number>>(new Set())

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

// ── 内部链接选项 ─────────────────────────────────────────────

interface InternalLink { label: string; url: string }

function internalLinks(): InternalLink[] {
  const base: InternalLink[] = [{ label: '首页', url: '/' }]
  if (props.themes?.length && props.tenantId) {
    for (const t of props.themes) {
      base.push({ label: `主题 · ${t.name}`, url: `/${props.tenantId}/${t.slug}` })
    }
  }
  return base
}

// ── 一级链接操作 ─────────────────────────────────────────────

function addLink() {
  form.links.push({ label: '', url: '' })
}

function removeLink(i: number) {
  expandedLinks.value.delete(i)
  form.links.splice(i, 1)
  const updated = new Set<number>()
  for (const idx of expandedLinks.value) {
    if (idx > i) updated.add(idx - 1)
    else updated.add(idx)
  }
  expandedLinks.value = updated
}

function toggleExpand(i: number) {
  if (expandedLinks.value.has(i)) {
    expandedLinks.value.delete(i)
  } else {
    expandedLinks.value.add(i)
    const link = form.links[i]
    if (link && !link.children) link.children = []
  }
}

// ── 二级链接操作 ─────────────────────────────────────────────

function addChild(i: number) {
  const link = form.links[i]
  if (!link) return
  if (!link.children) link.children = []
  link.children.push({ label: '', url: '' })
}

function removeChild(i: number, ci: number) {
  form.links[i]?.children?.splice(ci, 1)
}

function applyInternal(target: NavLink, url: string, label: string) {
  target.url = url
  if (!target.label) target.label = label
}

function needsExternalProtocolHint(url: string): boolean {
  const value = url.trim()
  if (!value) return false
  if (value.startsWith('#') || value.startsWith('/')) return false
  if (/^(https?:)?\/\//i.test(value)) return false
  return value.includes('.')
}

function externalProtocolHint(url: string): string {
  const value = url.trim()
  if (!value) return ''
  return `检测到这更像外部域名，请改为 https://${value.replace(/^\/\//, '')}`
}

// ── 选项元数据 ────────────────────────────────────────────────

const brandModes = [
  {
    value: 'both',
    label: '图文',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="6" y="14" width="10" height="10" rx="2" fill="currentColor" fill-opacity="0.3"/>
      <rect x="19" y="16" width="18" height="3" rx="1.5" fill="currentColor" fill-opacity="0.5"/>
      <rect x="19" y="21" width="12" height="2" rx="1" fill="currentColor" fill-opacity="0.25"/>
    </svg>`,
  },
  {
    value: 'logo_only',
    label: '仅 Logo',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="20" y="11" width="20" height="14" rx="3" fill="currentColor" fill-opacity="0.3"/>
    </svg>`,
  },
  {
    value: 'text_only',
    label: '仅文字',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="12" y="15" width="36" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/>
      <rect x="18" y="21" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/>
    </svg>`,
  },
] as const

const navbarStyles = [
  {
    value: 'solid',
    label: '纯色',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect width="60" height="11" rx="3" fill="currentColor" fill-opacity="0.18"/>
      <rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/>
      <rect x="36" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
      <rect x="47" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
    </svg>`,
  },
  {
    value: 'blur',
    label: '毛玻璃',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="5" y="16" width="50" height="4" rx="2" fill="currentColor" fill-opacity="0.1"/>
      <rect x="5" y="22" width="38" height="3" rx="1.5" fill="currentColor" fill-opacity="0.08"/>
      <rect width="60" height="11" rx="3" fill="currentColor" fill-opacity="0.1"/>
      <rect width="60" height="11" rx="3" fill="url(#blur-grad)" fill-opacity="0.5"/>
      <defs>
        <linearGradient id="blur-grad" x1="0" y1="0" x2="0" y2="11" gradientUnits="userSpaceOnUse">
          <stop stop-color="currentColor" stop-opacity="0.25"/>
          <stop offset="1" stop-color="currentColor" stop-opacity="0.05"/>
        </linearGradient>
      </defs>
      <rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/>
      <rect x="36" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
      <rect x="47" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
    </svg>`,
  },
  {
    value: 'transparent',
    label: '透明',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="5" y="16" width="50" height="4" rx="2" fill="currentColor" fill-opacity="0.1"/>
      <rect x="5" y="22" width="38" height="3" rx="1.5" fill="currentColor" fill-opacity="0.08"/>
      <rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/>
      <rect x="36" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
      <rect x="47" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/>
    </svg>`,
  },
] as const
</script>

<template>
  <div class="space-y-5">

    <!-- 品牌设置 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">品牌设置</legend>
      <div class="grid gap-2 sm:grid-cols-2">
        <div class="space-y-1.5">
          <Label for="navbar-brand-text" class="text-xs text-muted-foreground">品牌文字</Label>
          <Input id="navbar-brand-text" v-model="form.brand_text" placeholder="留空使用租户名称" />
        </div>
        <div class="space-y-1.5">
          <Label for="navbar-logo-url" class="text-xs text-muted-foreground">Logo URL</Label>
          <Input id="navbar-logo-url" v-model="form.logo_url" placeholder="留空使用租户 Logo" />
        </div>
      </div>

      <!-- 品牌区域显示模式 -->
      <div class="space-y-2">
        <Label class="text-xs text-muted-foreground">品牌区域显示</Label>
        <div class="grid grid-cols-3 gap-2">
          <button
            v-for="m in brandModes"
            :key="m.value"
            type="button"
            :class="[
              'group flex flex-col items-center gap-1.5 rounded-xl border-2 p-2.5 text-center transition-all duration-200',
              form.brand_mode === m.value
                ? 'border-primary bg-primary/5 shadow-sm'
                : 'border-border hover:border-primary/40 hover:bg-muted/50',
            ]"
            @click="form.brand_mode = m.value"
          >
            <span
              class="w-full"
              :class="form.brand_mode === m.value ? 'text-primary' : 'text-muted-foreground'"
              v-html="m.preview"
            />
            <span
              class="text-[11px] font-medium leading-tight"
              :class="form.brand_mode === m.value ? 'text-primary' : 'text-muted-foreground'"
            >{{ m.label }}</span>
          </button>
        </div>
      </div>
    </fieldset>

    <!-- 导航栏风格 & 行为 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">导航栏风格</Label>
      <div class="grid grid-cols-3 gap-2">
        <button
          v-for="s in navbarStyles"
          :key="s.value"
          type="button"
          :class="[
            'group flex flex-col items-center gap-1.5 rounded-xl border-2 p-2.5 text-center transition-all duration-200',
            form.style === s.value
              ? 'border-primary bg-primary/5 shadow-sm'
              : 'border-border hover:border-primary/40 hover:bg-muted/50',
          ]"
          @click="form.style = s.value"
        >
          <span
            class="w-full"
            :class="form.style === s.value ? 'text-primary' : 'text-muted-foreground'"
            v-html="s.preview"
          />
          <span
            class="text-[11px] font-medium leading-tight"
            :class="form.style === s.value ? 'text-primary' : 'text-muted-foreground'"
          >{{ s.label }}</span>
        </button>
      </div>

      <!-- 吸顶开关 -->
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch id="navbar-sticky" v-model="form.sticky" />
        <span class="text-sm">吸顶固定</span>
        <span class="ml-auto text-xs text-muted-foreground">滚动时保持在页面顶部</span>
      </label>
    </div>

    <!-- 导航链接 -->
    <fieldset class="rounded-xl border p-3.5 space-y-2">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider flex items-center gap-2">
        <span>导航链接</span>
      </legend>
      <div class="rounded-3xl border border-amber-200/70 bg-amber-50/80 px-3 py-2 text-[11px] leading-5 text-amber-700 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-200">
        自定义外部链接必须填写完整协议，例如 <span class="font-mono">https://survey.microswift.cn</span> 或 <span class="font-mono">http://example.com</span>；如果只填域名，浏览器会把它当作站内相对路径并自动拼接当前站点前缀。
      </div>

      <!-- 每条一级链接 -->
      <div v-for="(link, i) in form.links" :key="i" class="space-y-1">
        <!-- 一级链接行 -->
        <div class="flex items-center gap-1.5">
          <Tooltip>
            <TooltipTrigger as-child>
              <button
                class="p-0.5 rounded hover:bg-muted transition-colors shrink-0"
                @click="toggleExpand(i)"
              >
                <ChevronRight
                  class="h-3.5 w-3.5 text-muted-foreground transition-transform duration-150"
                  :class="{ 'rotate-90': expandedLinks.has(i) }"
                />
              </button>
            </TooltipTrigger>
            <TooltipContent>{{ expandedLinks.has(i) ? '收起子链接' : '展开子链接' }}</TooltipContent>
          </Tooltip>

          <Input v-model="link.label" placeholder="标签" class="flex-1 min-w-0" />

          <div class="flex items-center gap-0.5 flex-1 min-w-0">
            <Input v-model="link.url" placeholder="站内路径 /xxx、锚点 #xxx、外链 https://xxx" class="flex-1 min-w-0" />
            <Popover>
              <PopoverTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full h-7 w-7 shrink-0 data-[state=open]:bg-muted" title="选择内部链接">
                  <Link2 class="h-3.5 w-3.5" />
                </Button>
              </PopoverTrigger>
              <PopoverContent align="end" :side-offset="4" class="w-52 rounded-3xl p-1.5 shadow-none">
                <p class="px-2 py-1 text-[11px] text-muted-foreground font-medium">内部页面</p>
                <div class="max-h-[168px] overflow-y-auto">
                  <button
                    v-for="il in internalLinks()" :key="il.url"
                    class="w-full text-left px-2.5 py-1.5 text-sm rounded-2xl transition-colors active:bg-muted/80"
                    :class="link.url === il.url ? 'bg-muted text-foreground' : 'hover:bg-muted text-foreground/80'"
                    @click="applyInternal(link, il.url, il.label)"
                  >
                    {{ il.label }}
                    <span class="block text-[11px] text-muted-foreground truncate">{{ il.url }}</span>
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          </div>

          <Tooltip>
            <TooltipTrigger as-child>
              <Button variant="ghost" size="icon" class="rounded-full h-7 w-7 shrink-0 text-destructive/70 hover:text-destructive" @click="removeLink(i)">
                <Trash2 class="h-3.5 w-3.5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>删除</TooltipContent>
          </Tooltip>
        </div>
        <p v-if="needsExternalProtocolHint(link.url)" class="ml-6 text-[11px] text-amber-600 dark:text-amber-300">
          {{ externalProtocolHint(link.url) }}
        </p>

        <!-- 二级子链接区域 -->
        <div v-if="expandedLinks.has(i)" class="ml-5 space-y-1 border-l border-border/60 pl-3">
          <div v-for="(child, ci) in (link.children ?? [])" :key="ci" class="flex items-center gap-1.5">
            <Input v-model="child.label" placeholder="子标签" class="flex-1 min-w-0" />
            <div class="flex items-center gap-0.5 flex-1 min-w-0">
              <Input v-model="child.url" placeholder="站内路径 /xxx、锚点 #xxx、外链 https://xxx" class="flex-1 min-w-0" />
              <Popover>
                <PopoverTrigger as-child>
                  <Button variant="ghost" size="icon" class="rounded-full h-7 w-7 shrink-0 data-[state=open]:bg-muted">
                    <Link2 class="h-3.5 w-3.5" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent align="end" :side-offset="4" class="w-52 rounded-3xl p-1.5 shadow-none">
                  <p class="px-2 py-1 text-[11px] text-muted-foreground font-medium">内部页面</p>
                  <div class="max-h-[168px] overflow-y-auto">
                    <button
                      v-for="il in internalLinks()" :key="il.url"
                      class="w-full text-left px-2.5 py-1.5 text-sm rounded-2xl transition-colors active:bg-muted/80"
                      :class="child.url === il.url ? 'bg-muted text-foreground' : 'hover:bg-muted text-foreground/80'"
                      @click="applyInternal(child, il.url, il.label)"
                    >
                      {{ il.label }}
                      <span class="block text-[11px] text-muted-foreground truncate">{{ il.url }}</span>
                    </button>
                  </div>
                </PopoverContent>
              </Popover>
            </div>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="ghost" size="icon" class="rounded-full h-7 w-7 shrink-0 text-destructive/70 hover:text-destructive" @click="removeChild(i, ci)">
                  <Trash2 class="h-3.5 w-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>删除子链接</TooltipContent>
            </Tooltip>
          </div>
          <template v-for="(child, ci) in (link.children ?? [])" :key="`child-hint-${i}-${ci}`">
            <p v-if="needsExternalProtocolHint(child.url)" class="text-[11px] text-amber-600 dark:text-amber-300">
              {{ child.label ? `子链接“${child.label}”` : '子链接' }}：{{ externalProtocolHint(child.url) }}
            </p>
          </template>

          <button
            class="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors py-0.5"
            @click="addChild(i)"
          >
            <Plus class="h-3 w-3" />
            添加子链接
          </button>
        </div>
      </div>

      <p v-if="!form.links.length" class="py-2 text-xs text-center text-muted-foreground">暂无导航链接</p>

      <!-- 添加链接 -->
      <button
        class="w-full rounded-lg border border-dashed border-border py-2 text-xs text-muted-foreground hover:border-primary/50 hover:text-primary transition-colors"
        @click="addLink"
      >
        <Plus class="h-3 w-3 inline mr-1" />
        添加链接
      </button>
    </fieldset>

    <!-- 主题切换开关 -->
    <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
      <Switch id="navbar-theme-toggle" v-model="form.show_theme_toggle" />
      <span class="text-sm">显示主题切换开关</span>
      <span class="ml-auto text-xs text-muted-foreground">显示在 CTA 按钮前</span>
    </label>

    <!-- CTA 按钮（可选） -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">CTA 按钮（可选）</legend>
      <div v-if="form.cta_button" class="space-y-2">
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">按钮文字</Label>
          <Input v-model="form.cta_button.text" placeholder="立即体验" />
        </div>
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">链接</Label>
          <LinkInput v-model="form.cta_button.url" :tenant-id="tenantId" />
        </div>
        <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
          <Switch v-model="form.cta_button.show_arrow" />
          <span class="text-sm">显示箭头 →</span>
        </label>
        <button
          type="button"
          class="text-xs text-destructive hover:underline"
          @click="form.cta_button = null"
        >移除 CTA 按钮</button>
      </div>
      <button
        v-else
        type="button"
        class="w-full rounded-lg border border-dashed border-border py-2 text-xs text-muted-foreground hover:border-primary/50 hover:text-primary transition-colors"
        @click="form.cta_button = { text: '立即体验', url: '#themes', variant: 'dark', show_arrow: true }"
      >+ 添加 CTA 按钮</button>
    </fieldset>

  </div>
</template>
