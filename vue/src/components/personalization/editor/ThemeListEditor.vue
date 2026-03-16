<!-- ThemeListEditor.vue - 文档主题列表区块配置编辑器
     职责：编辑 ThemeListConfig 的标题、描述、列表风格和显示选项
     对外接口：
       Props: config
       Emits: update(config: ThemeListConfig) -->
<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import type { ThemeListConfig, TitleAlign } from '../types'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AlignLeft, AlignCenter, AlignRight } from 'lucide-vue-next'

const alignOptions: { value: TitleAlign; label: string; icon: any }[] = [
  { value: 'left', label: '靠左', icon: AlignLeft },
  { value: 'center', label: '居中', icon: AlignCenter },
  { value: 'right', label: '靠右', icon: AlignRight },
]

const props = defineProps<{ config: ThemeListConfig }>()
const emit = defineEmits<{ update: [config: ThemeListConfig] }>()

const form = reactive<ThemeListConfig>(JSON.parse(JSON.stringify(props.config)))

watch(form, () => emit('update', JSON.parse(JSON.stringify(form))), { deep: true })

/** 6 种列表风格选项 */
const listStyles = [
  {
    value: 'list-simple',
    label: '简洁行',
    desc: '图标+名称+描述，带分割线',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="3" width="54" height="34" rx="3" stroke="currentColor" stroke-opacity="0.2" stroke-width="1" fill="none"/>
      <rect x="8" y="9" width="6" height="6" rx="1.5" fill="currentColor" fill-opacity="0.22"/>
      <rect x="18" y="10" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/>
      <rect x="18" y="14" width="14" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <rect x="50" y="11" width="4" height="4" rx="1" fill="currentColor" fill-opacity="0.15"/>
      <line x1="3" y1="21" x2="57" y2="21" stroke="currentColor" stroke-opacity="0.12" stroke-width="1"/>
      <rect x="8" y="25" width="6" height="6" rx="1.5" fill="currentColor" fill-opacity="0.22"/>
      <rect x="18" y="26" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/>
      <rect x="18" y="30" width="16" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <rect x="50" y="27" width="4" height="4" rx="1" fill="currentColor" fill-opacity="0.15"/>
    </svg>`,
  },
  {
    value: 'list-compact',
    label: '密集行',
    desc: '名称与描述同行，行高最小',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="3" width="54" height="34" rx="3" stroke="currentColor" stroke-opacity="0.2" stroke-width="1" fill="none"/>
      <rect x="7" y="9" width="14" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="24" y="9.5" width="18" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <line x1="3" y1="15" x2="57" y2="15" stroke="currentColor" stroke-opacity="0.1" stroke-width="1"/>
      <rect x="7" y="19" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="28" y="19.5" width="12" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <line x1="3" y1="25" x2="57" y2="25" stroke="currentColor" stroke-opacity="0.1" stroke-width="1"/>
      <rect x="7" y="29" width="12" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="22" y="29.5" width="20" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
    </svg>`,
  },
  {
    value: 'list-numbered',
    label: '编号行',
    desc: '序号 + 名称 + 描述',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="3" width="54" height="34" rx="3" stroke="currentColor" stroke-opacity="0.2" stroke-width="1" fill="none"/>
      <rect x="7" y="9" width="6" height="6" rx="1" fill="currentColor" fill-opacity="0.18"/>
      <rect x="17" y="10" width="18" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/>
      <rect x="17" y="14" width="12" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <line x1="3" y1="21" x2="57" y2="21" stroke="currentColor" stroke-opacity="0.12" stroke-width="1"/>
      <rect x="7" y="25" width="6" height="6" rx="1" fill="currentColor" fill-opacity="0.18"/>
      <rect x="17" y="26" width="22" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/>
      <rect x="17" y="30" width="14" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
    </svg>`,
  },
  {
    value: 'list-two-col',
    label: '双列网格',
    desc: '两列布局，每项独立边框',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="4" width="26" height="14" rx="2.5" stroke="currentColor" stroke-opacity="0.25" stroke-width="1" fill="currentColor" fill-opacity="0.04"/>
      <rect x="7" y="8" width="14" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="7" y="12" width="10" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <rect x="31" y="4" width="26" height="14" rx="2.5" stroke="currentColor" stroke-opacity="0.25" stroke-width="1" fill="currentColor" fill-opacity="0.04"/>
      <rect x="35" y="8" width="16" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="35" y="12" width="10" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <rect x="3" y="22" width="26" height="14" rx="2.5" stroke="currentColor" stroke-opacity="0.25" stroke-width="1" fill="currentColor" fill-opacity="0.04"/>
      <rect x="7" y="26" width="10" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="7" y="30" width="14" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
      <rect x="31" y="22" width="26" height="14" rx="2.5" stroke="currentColor" stroke-opacity="0.25" stroke-width="1" fill="currentColor" fill-opacity="0.04"/>
      <rect x="35" y="26" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.55"/>
      <rect x="35" y="30" width="10" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/>
    </svg>`,
  },
  {
    value: 'list-tag',
    label: '标签云',
    desc: '药丸标签，空间利用率最高',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="8" width="18" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="6" y="11" width="12" height="2" rx="1" fill="currentColor" fill-opacity="0.45"/>
      <rect x="25" y="8" width="22" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="28" y="11" width="16" height="2" rx="1" fill="currentColor" fill-opacity="0.45"/>
      <rect x="51" y="8" width="6" height="8" rx="3" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="20" width="14" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="6" y="23" width="8" height="2" rx="1" fill="currentColor" fill-opacity="0.45"/>
      <rect x="21" y="20" width="24" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="24" y="23" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.45"/>
      <rect x="49" y="20" width="8" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="3" y="32" width="20" height="6" rx="3" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/>
      <rect x="6" y="34.5" width="14" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.45"/>
    </svg>`,
  },
  {
    value: 'list-headline',
    label: '大标题行',
    desc: '无容器边框，字体突出',
    preview: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/>
      <rect x="4" y="7" width="30" height="3.5" rx="1.75" fill="currentColor" fill-opacity="0.65"/>
      <rect x="4" y="12.5" width="20" height="2" rx="1" fill="currentColor" fill-opacity="0.2"/>
      <line x1="4" y1="18" x2="56" y2="18" stroke="currentColor" stroke-opacity="0.15" stroke-width="1"/>
      <rect x="4" y="22" width="24" height="3.5" rx="1.75" fill="currentColor" fill-opacity="0.65"/>
      <rect x="4" y="27.5" width="16" height="2" rx="1" fill="currentColor" fill-opacity="0.2"/>
      <line x1="4" y1="33" x2="56" y2="33" stroke="currentColor" stroke-opacity="0.15" stroke-width="1"/>
    </svg>`,
  },
] as const

/** 标签云模式下隐藏描述开关 */
const isTagMode = computed(() => form.card_style === 'list-tag')
</script>

<template>
  <div class="space-y-5">

    <!-- 文本内容 -->
    <fieldset class="rounded-xl border p-3.5 space-y-3">
      <legend class="text-xs font-semibold px-1.5 text-muted-foreground uppercase tracking-wider">文本内容</legend>
      <div class="space-y-1.5">
        <Label for="tl-title" class="text-xs text-muted-foreground">标题</Label>
        <Input id="tl-title" v-model="form.title" />
      </div>
      <div class="space-y-1.5">
        <Label for="tl-description" class="text-xs text-muted-foreground">描述（可选）</Label>
        <Input id="tl-description" v-model="form.description" placeholder="留空不显示" />
      </div>
      <!-- 标题对齐 -->
      <div class="space-y-1.5">
        <Label class="text-xs text-muted-foreground">标题对齐</Label>
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
    </fieldset>

    <!-- 列表风格 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">列表风格</Label>
      <div class="grid grid-cols-2 gap-2">
        <button
          v-for="s in listStyles"
          :key="s.value"
          type="button"
          :class="[
            'flex flex-col items-start gap-2 rounded-xl border-2 p-2.5 text-left transition-all duration-200',
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
          <span class="w-full">
            <span
              class="block text-[11px] font-semibold leading-tight"
              :class="form.card_style === s.value ? 'text-primary' : 'text-foreground/80'"
            >{{ s.label }}</span>
            <span
              class="block text-[10px] leading-snug mt-0.5"
              :class="form.card_style === s.value ? 'text-primary/70' : 'text-muted-foreground/70'"
            >{{ s.desc }}</span>
          </span>
        </button>
      </div>
    </div>

    <!-- 显示选项 -->
    <div class="space-y-2">
      <Label class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">显示选项</Label>
      <label
        v-if="!isTagMode"
        class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors"
      >
        <Switch id="tl-show-description" v-model="form.show_description" />
        <span class="text-sm">显示主题描述</span>
      </label>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch id="tl-show-slug" v-model="form.show_slug" />
        <div>
          <span class="text-sm block">显示 Slug</span>
          <span class="text-[11px] text-muted-foreground/60 leading-tight">主题的 URL 路径标识，如 api-reference</span>
        </div>
      </label>
      <label class="flex items-center gap-2.5 rounded-xl border px-3 py-2.5 cursor-pointer hover:bg-muted/40 transition-colors">
        <Switch id="tl-show-date" v-model="form.show_date" />
        <span class="text-sm">显示创建时间</span>
      </label>
    </div>

  </div>
</template>
