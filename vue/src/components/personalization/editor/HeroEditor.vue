<!--
  src/components/personalization/editor/HeroEditor.vue
  职责：Hero 区块配置编辑器，维护布局、文案、背景、按钮与动画配置
  对外暴露：Props(config)、Emits(update)
-->
<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import type { HeroConfig } from '../types'
import type { TitleAlign } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import LinkInput from './LinkInput.vue'
import { useAuthStore } from '@/stores/auth'
import { AlignLeft, AlignCenter, AlignRight } from 'lucide-vue-next'

const alignOptions: { value: TitleAlign; label: string; icon: any }[] = [
  { value: 'left', label: '靠左', icon: AlignLeft },
  { value: 'center', label: '居中', icon: AlignCenter },
  { value: 'right', label: '靠右', icon: AlignRight },
]

const props = defineProps<{ config: HeroConfig }>()
const emit = defineEmits<{ update: [config: HeroConfig] }>()

const authStore = useAuthStore()
const tenantId = computed(() => authStore.user?.tenant_id ?? '')

const form = reactive<HeroConfig>(JSON.parse(JSON.stringify(props.config)))

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

const layouts = [
  {
    value: 'centered',
    label: '居中',
    description: '标题居中，全屏展开',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="15" y="10" width="30" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/>
      <rect x="20" y="16" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.3"/>
      <rect x="22" y="21" width="16" height="7" rx="3.5" fill="currentColor" fill-opacity="0.25"/>
    </svg>`,
  },
  {
    value: 'left-right',
    label: '左文右图',
    description: '文字居左，配图居右',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="8" width="22" height="3.5" rx="1.75" fill="currentColor" fill-opacity="0.5"/>
      <rect x="4" y="14" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.3"/>
      <rect x="4" y="18" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.2"/>
      <rect x="4" y="24" width="13" height="6" rx="3" fill="currentColor" fill-opacity="0.25"/>
      <rect x="34" y="6" width="22" height="24" rx="3" fill="currentColor" fill-opacity="0.15"/>
    </svg>`,
  },
  {
    value: 'right-left',
    label: '左图右文',
    description: '配图居左，文字居右',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="6" width="22" height="24" rx="3" fill="currentColor" fill-opacity="0.15"/>
      <rect x="34" y="8" width="22" height="3.5" rx="1.75" fill="currentColor" fill-opacity="0.5"/>
      <rect x="34" y="14" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.3"/>
      <rect x="34" y="18" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.2"/>
      <rect x="34" y="24" width="13" height="6" rx="3" fill="currentColor" fill-opacity="0.25"/>
    </svg>`,
  },
] as const
</script>

<template>
  <div class="space-y-5">

    <!-- 布局模式 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">布局模式</Label>
      <div class="grid grid-cols-3 gap-2">
        <button
          v-for="l in layouts"
          :key="l.value"
          type="button"
          @click="form.layout = l.value as HeroConfig['layout']"
          :class="[
            'group flex flex-col items-center gap-1.5 rounded-xl border-2 p-2.5 text-center transition-all duration-200',
            form.layout === l.value
              ? 'border-primary bg-primary/5 shadow-sm'
              : 'border-border hover:border-primary/40 hover:bg-muted/50',
          ]"
        >
          <span
            class="w-full text-foreground"
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

    <!-- 入场动画 -->
    <div class="space-y-1.5">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">入场动画</Label>
      <Select v-model="form.animation">
        <SelectTrigger class="w-full">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="none">无动画</SelectItem>
          <SelectItem value="fade-up">上滑渐入</SelectItem>
          <SelectItem value="fade-in">渐入</SelectItem>
          <SelectItem value="typewriter">打字机效果</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <!-- 文本内容 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">文本内容</legend>
      <div class="space-y-1.5">
        <Label for="hero-subtitle" class="text-xs text-muted-foreground">副标题（上方标签）</Label>
        <Input id="hero-subtitle" v-model="form.subtitle" placeholder="如：文档中心 · BETA" />
      </div>
      <div class="space-y-1.5">
        <Label for="hero-title" class="text-xs text-muted-foreground">
          主标题
          <span v-if="form.animation === 'typewriter'" class="ml-1.5 inline-flex items-center gap-1 rounded-full bg-amber-100 px-1.5 py-0.5 text-[10px] font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-400">
            <svg class="h-2.5 w-2.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
            打字机将逐字显示
          </span>
        </Label>
        <Input id="hero-title" v-model="form.title" placeholder="欢迎来到文档中心" />
      </div>
      <div class="space-y-1.5">
        <Label for="hero-highlight" class="text-xs text-muted-foreground">灰显文字（标题中要淡化的部分）</Label>
        <Input id="hero-highlight" v-model="form.highlight_text" placeholder="留空不灰显，输入标题中的一段文字" />
        <p class="text-[11px] text-muted-foreground/60">Cursor 风格：标题中部分文字以灰色显示，形成视觉层次</p>
      </div>
      <div class="space-y-1.5">
        <Label for="hero-description" class="text-xs text-muted-foreground">描述文本</Label>
        <Textarea id="hero-description" v-model="form.description" placeholder="在这里查阅我们的技术文档" :rows="2" />
      </div>
    </fieldset>

    <!-- 背景图（居中布局） -->
    <fieldset v-if="form.layout === 'centered'" class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">背景图</legend>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">背景图 URL</Label>
        <Input v-model="form.background_image_url" placeholder="https://... 水彩/艺术画风效果最佳" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">遮罩颜色（亮色模式）</Label>
        <Input v-model="form.background_overlay" placeholder="rgba(245,240,232,0.82)" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">遮罩颜色（暗色模式）</Label>
        <Input v-model="form.background_overlay_dark" placeholder="rgba(26,26,26,0.85)" />
      </div>
    </fieldset>

    <!-- 配图 URL（分栏布局） -->
    <div v-if="form.layout !== 'centered'" class="space-y-1.5">
      <Label for="hero-image-url" class="text-xs text-muted-foreground">配图 URL</Label>
      <Input id="hero-image-url" v-model="form.image_url" placeholder="https://..." />
    </div>

    <!-- 主按钮 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">主按钮</legend>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">按钮文字</Label>
        <Input v-model="form.primary_button!.text" placeholder="开始阅读" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">链接</Label>
        <LinkInput v-model="form.primary_button!.url" :tenant-id="tenantId" />
      </div>
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">按钮样式</Label>
        <div class="flex flex-wrap gap-2">
          <!-- dark -->
          <span
            class="inline-flex cursor-pointer items-center rounded-full bg-foreground px-3 py-1 text-[11px] font-semibold text-background ring-2 ring-offset-1 transition-all"
            :class="form.primary_button!.variant === 'dark' ? 'ring-primary' : 'ring-transparent hover:ring-primary/40'"
            @click="form.primary_button!.variant = 'dark'"
          >深色药丸</span>

          <!-- noise -->
          <span
            class="relative inline-flex cursor-pointer items-center rounded-full p-[3px] ring-2 ring-offset-1 transition-all"
            :style="{ background: 'linear-gradient(to right, rgb(255,100,150), rgb(100,150,255), rgb(255,200,100))' }"
            :class="form.primary_button!.variant === 'noise' ? 'ring-primary' : 'ring-transparent hover:ring-primary/40'"
            @click="form.primary_button!.variant = 'noise'"
          >
            <span class="rounded-full bg-neutral-100 px-3 py-1 text-[11px] font-semibold text-black dark:bg-neutral-900 dark:text-white">噪点渐变</span>
          </span>

          <!-- primary -->
          <span
            class="inline-flex cursor-pointer items-center rounded-full bg-primary px-3 py-1 text-[11px] font-semibold text-primary-foreground ring-2 ring-offset-1 transition-all"
            :class="form.primary_button!.variant === 'primary' ? 'ring-primary' : 'ring-transparent hover:ring-primary/40'"
            @click="form.primary_button!.variant = 'primary'"
          >主色填充</span>

          <!-- outline -->
          <span
            class="inline-flex cursor-pointer items-center rounded-full border border-foreground/30 px-3 py-1 text-[11px] font-semibold text-foreground ring-2 ring-offset-1 transition-all"
            :class="form.primary_button!.variant === 'outline' ? 'ring-primary' : 'ring-transparent hover:ring-primary/40'"
            @click="form.primary_button!.variant = 'outline'"
          >描边</span>
        </div>
      </div>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch v-model="form.primary_button!.show_arrow" />
        <span class="text-sm">显示箭头 →</span>
      </label>
    </fieldset>

    <!-- 次按钮（可选） -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">次按钮（可选）</legend>
      <div v-if="form.secondary_button" class="space-y-2">
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">按钮文字</Label>
          <Input v-model="form.secondary_button.text" placeholder="了解更多" />
        </div>
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">链接</Label>
          <LinkInput v-model="form.secondary_button.url" :tenant-id="tenantId" />
        </div>
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">按钮样式</Label>
          <Select v-model="form.secondary_button.variant">
            <SelectTrigger class="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="plain">纯文字</SelectItem>
              <SelectItem value="outline">描边</SelectItem>
              <SelectItem value="secondary">次要</SelectItem>
              <SelectItem value="dark">深色药丸</SelectItem>
              <SelectItem value="primary">主色</SelectItem>
              <SelectItem value="noise">噪点渐变</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
          <Switch v-model="form.secondary_button.show_arrow" />
          <span class="text-sm">显示箭头 →</span>
        </label>
        <button
          type="button"
          class="text-xs text-destructive hover:underline"
          @click="form.secondary_button = null"
        >移除次按钮</button>
      </div>
      <button
        v-else
        type="button"
        class="w-full rounded-lg border border-dashed border-border py-2 text-xs text-muted-foreground hover:border-primary/50 hover:text-primary transition-colors"
        @click="form.secondary_button = { text: '了解更多', url: '/docs', variant: 'outline', show_arrow: false }"
      >+ 添加次按钮</button>
    </fieldset>

  </div>
</template>
