<!-- PThemeList.vue - 文档主题列表区块入口
     职责：根据 ThemeListConfig 以列表形式渲染文档主题，支持 6 种节省空间的列表样式
     对外接口：
       Props: config, themes, tenantId
       Emits: navigate(url: string) -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import type { Theme } from '@/utils/types'
import type { ThemeListConfig } from '../types'
import { useReaderStore } from '@/stores/reader'
import ThemeRowSimple from './ThemeRowSimple.vue'
import ThemeRowCompact from './ThemeRowCompact.vue'
import ThemeRowNumbered from './ThemeRowNumbered.vue'
import ThemeRowTwoCol from './ThemeRowTwoCol.vue'
import ThemeRowTag from './ThemeRowTag.vue'
import ThemeRowHeadline from './ThemeRowHeadline.vue'

const props = defineProps<{
  config: ThemeListConfig
  themes: Theme[]
  tenantId: string
}>()

const emit = defineEmits<{
  navigate: [url: string]
}>()

const store = useReaderStore()
const loadingThemeId = ref<string | null>(null)

const titleAlignClass = computed(() => {
  const map: Record<string, string> = { left: 'text-left', center: 'text-center', right: 'text-right' }
  return map[props.config.title_align] ?? 'text-center'
})

const descAlignClass = computed(() => {
  const a = props.config.title_align
  if (a === 'center') return 'mx-auto'
  if (a === 'right') return 'ml-auto'
  return ''
})

/** 带分割线的统一边框容器模式 */
const isDividedMode = computed(() =>
  props.config.card_style === 'list-simple' ||
  props.config.card_style === 'list-compact' ||
  props.config.card_style === 'list-numbered',
)

const dividedComponent = computed(() => {
  if (props.config.card_style === 'list-compact') return ThemeRowCompact
  if (props.config.card_style === 'list-numbered') return ThemeRowNumbered
  return ThemeRowSimple
})

/** 兼容旧版存储配置：字段不存在时默认 false */
const showSlug = computed(() => props.config.show_slug ?? false)
const showDate = computed(() => props.config.show_date ?? false)

async function goToTheme(theme: Theme) {
  if (loadingThemeId.value) return
  loadingThemeId.value = theme.id
  try {
    await store.loadVersions(theme.id)
    const defaultVersion = store.getDefaultVersion()
    if (!defaultVersion) {
      toast.error('该主题暂无可用版本')
      return
    }
    const targetUrl = `/${props.tenantId}/${theme.slug}/${defaultVersion.name}`
    emit('navigate', targetUrl)
  } catch {
    toast.error('加载版本失败，请重试')
  } finally {
    loadingThemeId.value = null
  }
}
</script>

<template>
  <section id="themes" style="padding-top: var(--hp-section-py, 96px); padding-bottom: var(--hp-section-py, 96px)">
    <div class="mx-auto px-4 sm:px-6 lg:px-8" style="max-width: var(--hp-max-width, 1200px)">

      <!-- 标题区 -->
      <div :class="titleAlignClass" class="mb-10 sm:mb-14">
        <h2 class="text-2xl sm:text-4xl font-bold tracking-tight">{{ config.title }}</h2>
        <p v-if="config.description" class="mt-3 text-sm sm:mt-4 sm:text-lg text-muted-foreground leading-relaxed max-w-2xl" :class="descAlignClass">
          {{ config.description }}
        </p>
      </div>

      <!-- 空状态 -->
      <p v-if="themes.length === 0" class="py-8 text-center text-muted-foreground">暂无文档主题</p>

      <template v-else>

        <!-- 简洁行 / 密集行 / 编号行：统一边框容器 + 分割线 -->
        <div
          v-if="isDividedMode"
          class="overflow-hidden rounded-3xl border border-border/60 bg-card divide-y divide-border/50"
        >
          <component
            :is="dividedComponent"
            v-for="(t, i) in themes"
            :key="t.id"
            :theme="t"
            :index="i + 1"
            :show-description="config.show_description"
            :show-slug="showSlug"
            :show-date="showDate"
            :loading="loadingThemeId === t.id"
            @click="goToTheme(t)"
          />
        </div>

        <!-- 双列网格：每项独立边框卡片 -->
        <div
          v-else-if="config.card_style === 'list-two-col'"
          class="grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-4"
        >
          <ThemeRowTwoCol
            v-for="t in themes"
            :key="t.id"
            :theme="t"
            :show-description="config.show_description"
            :show-slug="showSlug"
            :show-date="showDate"
            :loading="loadingThemeId === t.id"
            @click="goToTheme(t)"
          />
        </div>

        <!-- 标签云：flex-wrap 药丸 -->
        <div
          v-else-if="config.card_style === 'list-tag'"
          class="flex flex-wrap gap-2 sm:gap-2.5"
        >
          <ThemeRowTag
            v-for="t in themes"
            :key="t.id"
            :theme="t"
            :show-description="config.show_description"
            :show-slug="showSlug"
            :show-date="showDate"
            :loading="loadingThemeId === t.id"
            @click="goToTheme(t)"
          />
        </div>

        <!-- 大标题行：无容器边框，底部分割线 -->
        <div
          v-else-if="config.card_style === 'list-headline'"
          class="flex flex-col"
        >
          <ThemeRowHeadline
            v-for="t in themes"
            :key="t.id"
            :theme="t"
            :show-description="config.show_description"
            :show-slug="showSlug"
            :show-date="showDate"
            :loading="loadingThemeId === t.id"
            @click="goToTheme(t)"
          />
        </div>

      </template>
    </div>
  </section>
</template>
