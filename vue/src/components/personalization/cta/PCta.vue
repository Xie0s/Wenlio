<!--
  src/components/personalization/cta/PCta.vue
  职责：个性化首页 CTA 区块展示组件，渲染文案、按钮与卡片布局
  对外暴露：Props(config)、Emits(navigate)
-->
<script setup lang="ts">
import { computed } from 'vue'
import type { CtaConfig } from '../types'
import { ArrowRight } from 'lucide-vue-next'
import DynamicLucideIcon from '@/components/common/DynamicLucideIcon.vue'
import { useThemeStore } from '@/stores/theme'

const props = defineProps<{
  config: CtaConfig
}>()

const emit = defineEmits<{
  navigate: [url: string]
}>()

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.resolvedTheme === 'dark')

const hasBackgroundImage = computed(() => !!props.config.background_image_url)
const isCardLayout = computed(() => props.config.layout === 'card')

const sectionStyle = computed(() => {
  const s: Record<string, string> = {}
  if (!hasBackgroundImage.value && props.config.background_color) {
    s.backgroundColor = props.config.background_color
  }
  s.paddingTop = 'var(--hp-section-py, 96px)'
  s.paddingBottom = 'var(--hp-section-py, 96px)'
  return s
})

const titleAlignClass = computed(() => {
  const map: Record<string, string> = { left: 'text-left', center: 'text-center', right: 'text-right' }
  return map[props.config.title_align] ?? 'text-center'
})

const btnAlignClass = computed(() => {
  const map: Record<string, string> = { left: 'text-left', center: 'text-center', right: 'text-right' }
  return map[props.config.title_align] ?? 'text-center'
})

function handleClick() {
  const url = props.config.button.url
  if (url.startsWith('#')) {
    document.querySelector(url)?.scrollIntoView({ behavior: 'smooth' })
    return
  }
  emit('navigate', url)
}

function btnClass(variant: string): string {
  const base = 'inline-flex items-center justify-center gap-2 rounded-full text-sm font-medium transition-all duration-200 cursor-pointer active:scale-[0.98]'
  const variants: Record<string, string> = {
    dark:      `${base} bg-foreground text-background px-8 py-3.5 hover:opacity-90`,
    primary:   `${base} bg-primary text-primary-foreground px-8 py-3.5 hover:bg-primary/90 shadow-lg shadow-primary/20`,
    secondary: `${base} bg-secondary text-secondary-foreground px-8 py-3.5 hover:bg-secondary/80`,
    outline:   `${base} border-2 border-foreground/20 px-8 py-3.5 hover:border-foreground/40 hover:bg-foreground/5`,
  }
  return variants[variant] ?? variants.dark!
}
</script>

<template>
  <section class="relative overflow-hidden" :style="sectionStyle">

    <!-- 背景图 -->
    <img
      v-if="hasBackgroundImage"
      :src="config.background_image_url"
      alt=""
      class="absolute inset-0 h-full w-full object-cover"
    />
    <!-- 遮罩 -->
    <div
      v-if="hasBackgroundImage"
      class="absolute inset-0"
      :style="{ backgroundColor: isDark ? 'rgba(26,26,26,0.7)' : 'rgba(245,240,232,0.65)' }"
    />

    <div class="relative z-10 mx-auto px-4 sm:px-6 lg:px-8" style="max-width: var(--hp-max-width, 1200px)">
      <!-- Simple 布局 -->
      <div v-if="!isCardLayout" class="mx-auto max-w-2xl" :class="titleAlignClass">
        <h2 class="text-3xl sm:text-4xl font-bold tracking-tight">{{ config.title }}</h2>
        <p v-if="config.show_description !== false && config.description" class="mt-5 text-lg text-muted-foreground leading-relaxed">{{ config.description }}</p>
        <div v-if="config.show_button !== false" class="mt-10" :class="btnAlignClass">
          <button :class="btnClass(config.button.variant)" @click="handleClick">
            {{ config.button.text }}
            <ArrowRight v-if="config.button.show_arrow" class="h-4 w-4" />
          </button>
        </div>
      </div>

      <!-- Card 布局 -->
      <div v-else class="grid items-center gap-12 lg:grid-cols-2 lg:gap-16">
        <div :class="titleAlignClass">
          <h2 class="text-3xl sm:text-4xl font-bold tracking-tight">{{ config.title }}</h2>
          <p v-if="config.show_description !== false && config.description" class="mt-5 text-lg text-muted-foreground leading-relaxed">{{ config.description }}</p>
          <div v-if="config.show_button !== false" class="mt-10" :class="btnAlignClass">
            <button :class="btnClass(config.button.variant)" @click="handleClick">
              {{ config.button.text }}
              <ArrowRight v-if="config.button.show_arrow" class="h-4 w-4" />
            </button>
          </div>
        </div>

        <!-- 浮动卡片 -->
        <div class="relative rounded-3xl border border-border/40 overflow-hidden">
          <!-- 图片 -->
          <img
            v-if="config.card_image_url"
            :src="config.card_image_url"
            :alt="config.card_title"
            class="w-full object-cover block"
          />
          <!-- 空状态（无图片无标题时） -->
          <div v-if="!config.card_image_url && !config.card_title" class="flex flex-col items-center justify-center p-6 py-12 text-muted-foreground/30 bg-card/90">
            <svg class="h-12 w-12 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="0.8">
              <rect x="3" y="3" width="18" height="18" rx="4" />
              <circle cx="8.5" cy="8.5" r="1.5" />
              <path d="m21 15-5-5L5 21" />
            </svg>
            <span class="text-xs">在编辑器中配置卡片内容</span>
          </div>
          <!-- 毛玻璃底栏（HeroUI CardFooter 风格） -->
          <div
            v-if="config.card_title || config.card_description || config.card_button_text"
            class="absolute bottom-4 left-1/2 -translate-x-1/2 w-[calc(100%-16px)] flex items-stretch justify-between gap-3 border border-white/20 bg-white/10 backdrop-blur-md px-4 py-2.5 z-10"
            :class="config.card_description ? 'rounded-2xl' : 'rounded-xl'"
          >
            <!-- 文字区：无描述时居中，有描述时靠左 -->
            <div
              v-if="config.card_title || config.card_description"
              class="flex flex-1 min-w-0"
              :class="config.card_description ? 'flex-col gap-0.5 justify-center' : 'items-center justify-center'"
            >
              <p v-if="config.card_title" class="text-sm font-semibold text-white leading-snug truncate" :class="!config.card_description ? 'text-center' : ''">{{ config.card_title }}</p>
              <p v-if="config.card_description" class="text-xs text-white/70 leading-relaxed truncate">{{ config.card_description }}</p>
            </div>
            <!-- 按钮：随底栏高度自适应拉伸，圆角同步联动 -->
            <button
              v-if="config.card_button_text"
              class="shrink-0 flex items-center justify-center gap-1.5 bg-black/20 hover:bg-black/30 text-white text-xs font-medium px-3.5 transition-colors"
              :class="config.card_description ? 'rounded-xl' : 'rounded-lg'"
            >
              {{ config.card_button_text }}
              <DynamicLucideIcon v-if="config.card_button_icon" :name="config.card_button_icon" class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
