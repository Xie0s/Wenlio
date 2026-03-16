<!--
  src/components/personalization/editor/IntroductionEditor.vue
  职责：Introduction 区块配置编辑器，维护标题、布局、卡片风格与特性条目
  对外暴露：Props(config)、Emits(update)
-->
<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { IntroductionConfig, TitleAlign } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { AlignLeft, AlignCenter, AlignRight } from 'lucide-vue-next'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Plus, Trash2 } from 'lucide-vue-next'

const props = defineProps<{ config: IntroductionConfig }>()
const emit = defineEmits<{ update: [config: IntroductionConfig] }>()

const form = reactive<IntroductionConfig>(JSON.parse(JSON.stringify(props.config)))

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

const alignOptions: { value: TitleAlign; label: string; icon: any }[] = [
  { value: 'left', label: '靠左', icon: AlignLeft },
  { value: 'center', label: '居中', icon: AlignCenter },
  { value: 'right', label: '靠右', icon: AlignRight },
]

function addFeature() {
  form.features.push({ icon: 'zap', title: '', description: '', image_url: '', link_text: '', link_url: '' })
}

function removeFeature(i: number) {
  form.features.splice(i, 1)
}

const layouts = [
  {
    value: 'grid-3',
    label: '三列',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/>
      <rect x="22" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/>
      <rect x="40" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/>
    </svg>`,
  },
  {
    value: 'grid-2',
    label: '两列',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="8" width="24" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/>
      <rect x="32" y="8" width="24" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/>
    </svg>`,
  },
  {
    value: 'list',
    label: '列表',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="4" y="7" width="52" height="6" rx="2" fill="currentColor" fill-opacity="0.18"/>
      <rect x="4" y="15" width="52" height="6" rx="2" fill="currentColor" fill-opacity="0.18"/>
      <rect x="4" y="23" width="52" height="6" rx="2" fill="currentColor" fill-opacity="0.18"/>
    </svg>`,
  },
] as const

const cardStyles = [
  {
    value: 'elevated',
    label: '浮起',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.04"/>
      <rect x="8" y="6" width="44" height="24" rx="4" fill="currentColor" fill-opacity="0.12"/>
      <rect x="10" y="8" width="40" height="20" rx="3" fill="currentColor" fill-opacity="0.06"/>
    </svg>`,
  },
  {
    value: 'flat',
    label: '扁平',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.04"/>
      <rect x="10" y="8" width="40" height="20" rx="3" fill="currentColor" fill-opacity="0.05"/>
    </svg>`,
  },
  {
    value: 'bordered',
    label: '边框',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.04"/>
      <rect x="10" y="8" width="40" height="20" rx="3" stroke="currentColor" stroke-opacity="0.3" fill="none"/>
    </svg>`,
  },
  {
    value: 'glass',
    label: '玻璃',
    preview: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/>
      <rect x="10" y="8" width="40" height="20" rx="3" fill="currentColor" fill-opacity="0.08" stroke="currentColor" stroke-opacity="0.1"/>
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
      <div class="grid grid-cols-3 gap-2">
        <button
          v-for="l in layouts"
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

    <!-- 卡片风格 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">卡片风格</Label>
      <div class="grid grid-cols-4 gap-2">
        <button
          v-for="s in cardStyles"
          :key="s.value"
          type="button"
          :class="[
            'flex flex-col items-center gap-1.5 rounded-xl border-2 p-2 text-center transition-all duration-200',
            form.card_style === s.value
              ? 'border-primary bg-primary/5 shadow-sm'
              : 'border-border hover:border-primary/40 hover:bg-muted/50',
          ]"
          @click="form.card_style = s.value"
        >
          <span
            class="w-full"
            :class="form.card_style === s.value ? 'text-primary' : 'text-muted-foreground'"
            v-html="s.preview"
          />
          <span
            class="text-[10px] font-medium leading-tight"
            :class="form.card_style === s.value ? 'text-primary' : 'text-muted-foreground'"
          >{{ s.label }}</span>
        </button>
      </div>
    </div>

    <!-- 文本内容 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">文本内容</legend>
      <div class="space-y-1.5">
        <Label for="intro-title" class="text-xs text-muted-foreground">标题</Label>
        <Input id="intro-title" v-model="form.title" />
      </div>
      <div class="space-y-1.5">
        <Label for="intro-description" class="text-xs text-muted-foreground">描述（可选）</Label>
        <Textarea id="intro-description" v-model="form.description" placeholder="留空不显示" :rows="2" />
      </div>
    </fieldset>

    <!-- 特性条目 -->
    <fieldset class="rounded-xl border p-3.5 space-y-2">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">特性条目</legend>

      <div v-for="(item, i) in form.features" :key="i" class="rounded-xl border bg-muted/20 p-3 space-y-2">
        <div class="flex items-center justify-between gap-2">
          <span class="text-xs font-medium text-muted-foreground">条目 {{ i + 1 }}</span>
          <Tooltip>
            <TooltipTrigger as-child>
              <Button variant="ghost" size="icon" class="rounded-full h-6 w-6 shrink-0 text-destructive/70 hover:text-destructive" @click="removeFeature(i)">
                <Trash2 class="h-3 w-3" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>删除</TooltipContent>
          </Tooltip>
        </div>
        <div class="grid gap-2 sm:grid-cols-2">
          <div class="space-y-1">
            <Label class="text-[11px] text-muted-foreground">图标名</Label>
            <Input :name="`feature-icon-${i}`" v-model="item.icon" placeholder="如 zap、shield、users" />
          </div>
          <div class="space-y-1">
            <Label class="text-[11px] text-muted-foreground">标题</Label>
            <Input :name="`feature-title-${i}`" v-model="item.title" placeholder="特性标题" />
          </div>
        </div>
        <div class="space-y-1">
          <Label class="text-[11px] text-muted-foreground">描述</Label>
          <Input :name="`feature-desc-${i}`" v-model="item.description" placeholder="一句话说明" />
        </div>
        <div class="space-y-1">
          <Label class="text-[11px] text-muted-foreground">截图 URL（可选）</Label>
          <Input :name="`feature-img-${i}`" v-model="item.image_url" placeholder="https://... 留空不显示" />
        </div>
        <div class="grid gap-2 sm:grid-cols-2">
          <div class="space-y-1">
            <Label class="text-[11px] text-muted-foreground">链接文字（可选）</Label>
            <Input :name="`feature-link-text-${i}`" v-model="item.link_text" placeholder="了解更多" />
          </div>
          <div class="space-y-1">
            <Label class="text-[11px] text-muted-foreground">链接地址（可选）</Label>
            <Input :name="`feature-link-url-${i}`" v-model="item.link_url" placeholder="https://..." />
          </div>
        </div>
      </div>

      <p v-if="!form.features.length" class="py-2 text-xs text-center text-muted-foreground">暂无特性条目</p>

      <button
        class="w-full rounded-lg border border-dashed border-border py-2 text-xs text-muted-foreground hover:border-primary/50 hover:text-primary transition-colors"
        @click="addFeature"
      >
        <Plus class="h-3 w-3 inline mr-1" />
        添加特性
      </button>
    </fieldset>

  </div>
</template>
