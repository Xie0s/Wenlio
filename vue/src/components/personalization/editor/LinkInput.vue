<!-- LinkInput.vue - 链接类型选择 + 输入
     支持三种模式：锚点 / 本地路径（自动注入 tenantId） / 外部链接
     对外接口：
       Props: modelValue (string), tenantId? (string)
       Emits: update:modelValue -->
<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Input } from '@/components/ui/input'

const props = defineProps<{
  modelValue: string
  /** 当前租户 ID，用于本地链接自动补前缀 */
  tenantId?: string
}>()
const emit = defineEmits<{ 'update:modelValue': [v: string] }>()

type LinkType = 'anchor' | 'local' | 'external'

function detectType(v: string): LinkType {
  if (v.startsWith('#')) return 'anchor'
  if (v.startsWith('http://') || v.startsWith('https://') || v.startsWith('//')) return 'external'
  return 'local'
}

const linkType = ref<LinkType>(detectType(props.modelValue))

// ── 本地路径：拆解 /{tenantId}/ 前缀，只让用户编辑后半段 ──────
function stripTenantPrefix(url: string): string {
  if (!props.tenantId) return url.replace(/^\//, '')
  const prefix = `/${props.tenantId}/`
  if (url.startsWith(prefix)) return url.slice(prefix.length)
  if (url === `/${props.tenantId}`) return ''
  return url.replace(/^\//, '')
}

function buildLocalUrl(path: string): string {
  const p = path.trim()
  if (!props.tenantId) return p ? `/${p}` : '/'
  return p ? `/${props.tenantId}/${p}` : `/${props.tenantId}`
}

// 锚点 / 外部 共用一个 ref，本地模式单独拆分
const rawInput = ref(
  detectType(props.modelValue) === 'local'
    ? stripTenantPrefix(props.modelValue)
    : props.modelValue,
)

/** 供锚点/外部 v-model 使用 */
const rawModel = computed({
  get: () => rawInput.value,
  set: (v: string) => {
    rawInput.value = v
    emit('update:modelValue', v)
  },
})

/** 本地模式：用户只输入 slug 部分 */
const localModel = computed({
  get: () => rawInput.value,
  set: (v: string) => {
    rawInput.value = v
    emit('update:modelValue', buildLocalUrl(v))
  },
})

watch(() => props.modelValue, (v) => {
  const t = detectType(v)
  linkType.value = t
  rawInput.value = t === 'local' ? stripTenantPrefix(v) : v
})

function switchType(type: LinkType) {
  const prev = linkType.value
  linkType.value = type

  // 尽量保留用户已填内容
  let bare = rawInput.value
    .replace(/^#/, '')
    .replace(/^https?:\/\//, '')
    .replace(/^\/\//, '')
    .replace(/^\//, '')

  if (prev === 'local' && props.tenantId) {
    // rawInput 存的是 slug，不需要再处理
    bare = rawInput.value
  }

  if (type === 'anchor') {
    rawInput.value = bare ? `#${bare}` : '#'
    emit('update:modelValue', rawInput.value)
  } else if (type === 'local') {
    rawInput.value = bare
    emit('update:modelValue', buildLocalUrl(bare))
  } else {
    rawInput.value = bare ? `https://${bare}` : 'https://'
    emit('update:modelValue', rawInput.value)
  }
}

function pickAnchor(val: string) {
  rawModel.value = val
}

const localPrefix = computed(() =>
  props.tenantId ? `/${props.tenantId}/` : '/',
)

const placeholder = computed(() => ({
  anchor:   '#themes、#introduction …',
  local:    'theme-slug 或 theme-slug/v1.0',
  external: 'https://example.com',
}[linkType.value]))

const anchorSuggestions = [
  { label: '主题列表',   value: '#themes' },
  { label: '特性介绍',   value: '#introduction' },
  { label: '行动号召',   value: '#cta' },
  { label: '页面顶部',   value: '#hero' },
]

const tabs: { key: LinkType; label: string; path: string }[] = [
  {
    key: 'anchor',
    label: '锚点',
    path: 'M13 10V3L4 14h7v7l9-11h-7z',
  },
  {
    key: 'local',
    label: '本地',
    path: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
  },
  {
    key: 'external',
    label: '外部',
    path: 'M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14',
  },
]
</script>

<template>
  <div class="space-y-2">
    <!-- 类型切换 Tab -->
    <div class="inline-flex rounded-lg border border-border bg-muted/40 p-0.5 gap-0.5">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        @click.prevent="switchType(tab.key)"
        :class="[
          'inline-flex items-center gap-1 rounded-md px-2.5 py-1 text-[11px] font-medium transition-all select-none',
          linkType === tab.key
            ? 'bg-background text-foreground shadow-sm'
            : 'text-muted-foreground hover:text-foreground',
        ]"
      >
        <svg class="h-3 w-3 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" :d="tab.path" />
        </svg>
        {{ tab.label }}
      </button>
    </div>

    <!-- 本地模式：前缀 + 输入 -->
    <div v-if="linkType === 'local'" class="flex items-center overflow-hidden rounded-md border border-border focus-within:ring-1 focus-within:ring-ring">
      <span class="shrink-0 border-r border-border bg-muted px-2.5 py-2 font-mono text-[11px] text-muted-foreground select-none whitespace-nowrap">
        {{ localPrefix }}
      </span>
      <Input
        v-model="localModel"
        :placeholder="placeholder"
        class="rounded-none border-0 shadow-none focus-visible:ring-0 font-mono text-[12px]"
      />
    </div>

    <!-- 锚点 / 外部输入框 -->
    <Input
      v-else
      v-model="rawModel"
      :placeholder="placeholder"
    />

    <!-- 本地模式提示 -->
    <p v-if="linkType === 'local' && tenantId" class="text-[10px] text-muted-foreground">
      完整链接：<span class="font-mono text-foreground/70">{{ buildLocalUrl(rawInput) }}</span>
    </p>

    <!-- 锚点快捷选项 -->
    <div v-if="linkType === 'anchor'" class="flex flex-wrap gap-1.5">
      <button
        v-for="s in anchorSuggestions"
        :key="s.value"
        type="button"
        @click.prevent="pickAnchor(s.value)"
        :class="[
          'rounded-full border px-2 py-0.5 text-[10px] font-medium transition-colors',
          rawInput === s.value
            ? 'border-primary bg-primary/10 text-primary'
            : 'border-border text-muted-foreground hover:border-primary/50 hover:text-foreground',
        ]"
      >
        <span class="font-mono">{{ s.value }}</span>
        <span class="ml-1 opacity-60">{{ s.label }}</span>
      </button>
    </div>
  </div>
</template>
