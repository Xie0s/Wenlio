<!--
  src/components/personalization/editor/CtaEditor.vue
  职责：CTA 区块配置编辑器，维护标题、布局、背景与按钮等配置
  对外暴露：Props(config)、Emits(update)
-->
<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import type { CtaConfig, TitleAlign } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { AlignLeft, AlignCenter, AlignRight } from 'lucide-vue-next'
import LinkInput from './LinkInput.vue'
import { useAuthStore } from '@/stores/auth'

const alignOptions: { value: TitleAlign; label: string; icon: any }[] = [
  { value: 'left', label: '靠左', icon: AlignLeft },
  { value: 'center', label: '居中', icon: AlignCenter },
  { value: 'right', label: '靠右', icon: AlignRight },
]

const props = defineProps<{ config: CtaConfig }>()
const emit = defineEmits<{ update: [config: CtaConfig] }>()

const authStore = useAuthStore()
const tenantId = computed(() => authStore.user?.tenant_id ?? '')

const form = reactive<CtaConfig>(JSON.parse(JSON.stringify(props.config)))

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

const layoutOptions = [
  {
    value: 'simple',
    label: '居中',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="15" y="10" width="30" height="4" rx="2" fill="currentColor" fill-opacity="0.4"/>
      <rect x="18" y="16" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/>
      <rect x="22" y="22" width="16" height="6" rx="3" fill="currentColor" fill-opacity="0.3"/>
    </svg>`,
  },
  {
    value: 'card',
    label: '带卡片',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="8" width="24" height="4" rx="2" fill="currentColor" fill-opacity="0.4"/>
      <rect x="4" y="14" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/>
      <rect x="4" y="22" width="14" height="6" rx="3" fill="currentColor" fill-opacity="0.3"/>
      <rect x="34" y="6" width="22" height="24" rx="3" fill="currentColor" fill-opacity="0.15"/>
    </svg>`,
  },
] as const
</script>

<template>
  <div class="space-y-5">

    <!-- 标题对齐 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">标题对齐</Label>
      <Tabs
        :model-value="form.title_align || 'center'"
        class="items-start"
        @update:model-value="(value) => form.title_align = value as TitleAlign"
      >
        <TabsList>
          <TabsTrigger
            v-for="a in alignOptions"
            :key="a.value"
            :value="a.value"
          >
            <component :is="a.icon" class="h-3.5 w-3.5" />
            {{ a.label }}
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </div>

    <!-- 布局模式 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">布局模式</Label>
      <div class="grid grid-cols-2 gap-2">
        <button
          v-for="l in layoutOptions"
          :key="l.value"
          type="button"
          :class="[
            'flex flex-col items-center gap-1.5 rounded-xl border-2 p-2.5 text-center transition-all duration-200',
            form.layout === l.value
              ? 'border-primary bg-primary/5 shadow-sm'
              : 'border-border hover:border-primary/40 hover:bg-muted/50',
          ]"
          @click="form.layout = l.value"
        >
          <span
            class="w-full"
            :class="form.layout === l.value ? 'text-primary' : 'text-muted-foreground'"
            v-html="l.preview"
          />
          <span
            class="text-[11px] font-medium leading-tight"
            :class="form.layout === l.value ? 'text-primary' : 'text-muted-foreground'"
          >{{ l.label }}</span>
        </button>
      </div>
    </div>

    <!-- 文本内容 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">文本内容</legend>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">标题</Label>
        <Input v-model="form.title" placeholder="准备好了吗？" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">描述</Label>
        <Textarea v-model="form.description" placeholder="立即开始探索" :rows="2" />
      </div>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch v-model="form.show_description" />
        <span class="text-sm">显示描述文字</span>
      </label>
    </fieldset>

    <!-- 背景设置 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">背景设置</legend>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">背景图 URL</Label>
        <Input v-model="form.background_image_url" placeholder="https://... 水彩/艺术画风效果最佳" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">背景色（无背景图时生效）</Label>
        <Input v-model="form.background_color" placeholder="如 #f5f0e8 或留空" />
      </div>
    </fieldset>

    <!-- 按钮 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">按钮</legend>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch v-model="form.show_button" />
        <span class="text-sm">显示按钮</span>
      </label>
      <template v-if="form.show_button !== false">
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">按钮文字</Label>
        <Input v-model="form.button.text" placeholder="立即开始" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">链接</Label>
        <LinkInput v-model="form.button.url" :tenant-id="tenantId" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">按钮样式</Label>
        <Select v-model="form.button.variant">
          <SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="dark">深色药丸</SelectItem>
            <SelectItem value="primary">主色</SelectItem>
            <SelectItem value="outline">描边</SelectItem>
            <SelectItem value="secondary">次要</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch v-model="form.button.show_arrow" />
        <span class="text-sm">显示箭头 →</span>
      </label>
      </template>
    </fieldset>

    <!-- 浮动卡片（card 布局时） -->
    <fieldset v-if="form.layout === 'card'" class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">浮动卡片</legend>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">卡片截图 URL</Label>
        <Input v-model="form.card_image_url" placeholder="https://..." />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">卡片标题</Label>
        <Input v-model="form.card_title" placeholder="产品特性展示" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">卡片描述</Label>
        <Textarea v-model="form.card_description" placeholder="简短说明" :rows="2" />
      </div>
      <div class="grid grid-cols-2 gap-2">
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">按钮文字</Label>
          <Input v-model="form.card_button_text" placeholder="了解更多" />
        </div>
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">按钮图标（Lucide id）</Label>
          <Input v-model="form.card_button_icon" placeholder="arrow-right" />
        </div>
      </div>
    </fieldset>

  </div>
</template>
