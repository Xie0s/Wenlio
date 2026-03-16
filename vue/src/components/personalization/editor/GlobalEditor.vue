<!--
  src/components/personalization/editor/GlobalEditor.vue
  职责：首页全局配置编辑器，维护背景、字体与区块间距等全局样式
  对外暴露：Props(config)、Emits(update)
-->
<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import type { HomepageGlobal } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ColorPickerInput } from '@/components/ui/color-picker'
import { normalizeBrowserIconUrl, validateBrowserIconUrl } from '@/lib/validation'

const props = defineProps<{ config: HomepageGlobal }>()
const emit = defineEmits<{ update: [config: HomepageGlobal] }>()

const form = reactive<HomepageGlobal>({
  dark_background_color: '#1a1a1a',
  browser_title: '',
  browser_icon_url: '',
  section_spacing: 96,
  ...JSON.parse(JSON.stringify(props.config)),
})

const browserIconError = computed(() => validateBrowserIconUrl(form.browser_icon_url))

watch(form, () => {
  const next = JSON.parse(JSON.stringify(form)) as HomepageGlobal
  next.browser_icon_url = browserIconError.value ? '' : normalizeBrowserIconUrl(form.browser_icon_url)
  emit('update', next)
}, { deep: true })

const presets = [
  { label: 'Cursor 暖调', bg: '#f5f0e8', dark: '#1a1a1a' },
  { label: '纯白', bg: '#ffffff', dark: '#181818' },
  { label: '象牙', bg: '#faf8f5', dark: '#1c1c1e' },
  { label: '浅灰', bg: '#f8f9fa', dark: '#171717' },
]

function applyPreset(preset: typeof presets[0]) {
  form.background_color = preset.bg
  form.dark_background_color = preset.dark
}

const spacingPresets = [
  { label: '紧凑', value: 48 },
  { label: '适中', value: 72 },
  { label: '宽松', value: 96 },
  { label: '超宽', value: 128 },
]
</script>

<template>
  <div class="space-y-4">
    <!-- 预设色板 -->
    <div class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">配色预设</Label>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="p in presets" :key="p.label"
          type="button"
          class="flex items-center gap-2 rounded-lg border px-2.5 py-1.5 text-xs transition-all hover:border-primary/40"
          :class="form.background_color === p.bg ? 'border-primary bg-primary/5' : 'border-border'"
          @click="applyPreset(p)"
        >
          <span class="h-4 w-4 rounded-full border border-border/60" :style="{ backgroundColor: p.bg }" />
          {{ p.label }}
        </button>
      </div>
    </div>

    <div class="grid gap-3 sm:grid-cols-2">
      <div class="space-y-1.5">
        <Label>亮色背景色</Label>
        <ColorPickerInput v-model="form.background_color" placeholder="#f5f0e8" />
      </div>
      <div class="space-y-1.5">
        <Label>暗色背景色</Label>
        <ColorPickerInput v-model="form.dark_background_color" placeholder="#1a1a1a" />
      </div>
      <div class="space-y-1.5">
        <Label for="global-font-family">字体族</Label>
        <Input id="global-font-family" v-model="form.font_family" placeholder="留空继承默认" />
      </div>
      <div class="space-y-1.5">
        <Label for="global-max-width">内容最大宽度</Label>
        <Input id="global-max-width" v-model="form.max_width" placeholder="1200px" />
      </div>
      <div class="space-y-1.5">
        <Label for="global-browser-title">浏览器标签文字</Label>
        <Input id="global-browser-title" v-model="form.browser_title" placeholder="留空默认使用租户名称" />
      </div>
      <div class="space-y-1.5">
        <Label for="global-browser-icon">浏览器标签图标</Label>
        <Input id="global-browser-icon" v-model="form.browser_icon_url" placeholder="https://example.com/favicon.ico" />
        <p v-if="browserIconError" class="text-xs text-destructive">{{ browserIconError }}</p>
      </div>
    </div>

    <!-- 区块间距 -->
    <div class="space-y-2">
      <Label class="text-xs text-muted-foreground">区块间距（上下内边距）</Label>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="sp in spacingPresets" :key="sp.value"
          type="button"
          class="rounded-lg border px-2.5 py-1.5 text-xs transition-all hover:border-primary/40"
          :class="form.section_spacing === sp.value ? 'border-primary bg-primary/5 font-medium' : 'border-border'"
          @click="form.section_spacing = sp.value"
        >
          {{ sp.label }}
          <span class="ml-1 text-muted-foreground">{{ sp.value }}px</span>
        </button>
      </div>
      <div class="flex items-center gap-2">
        <Input
          id="global-section-spacing"
          type="number"
          :model-value="form.section_spacing"
          @update:model-value="form.section_spacing = Number($event) || 96"
          placeholder="96"
          class="w-24"
        />
        <span class="text-xs text-muted-foreground">px</span>
      </div>
    </div>
  </div>
</template>
