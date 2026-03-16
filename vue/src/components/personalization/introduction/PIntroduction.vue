<!--
  src/components/personalization/introduction/PIntroduction.vue
  职责：个性化首页 Introduction 区块展示组件，渲染标题与特性卡片列表
  对外暴露：Props(config)
-->
<script setup lang="ts">
import { computed } from 'vue'
import type { IntroductionConfig, IntroductionCardStyle } from '../types'
import { ArrowUpRight } from 'lucide-vue-next'
import DynamicLucideIcon from '@/components/common/DynamicLucideIcon.vue'

const props = defineProps<{
  config: IntroductionConfig
}>()

const gridClass = computed(() => {
  const map: Record<string, string> = {
    'grid-3': 'grid gap-5 sm:grid-cols-2 lg:grid-cols-3',
    'grid-2': 'grid gap-5 sm:grid-cols-2',
    'list': 'space-y-5 max-w-2xl mx-auto',
  }
  return map[props.config.layout] || map['grid-3']
})

const cardStyle = computed(() => props.config.card_style || 'elevated')

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

function outerCardClass(style: IntroductionCardStyle): string {
  const base = 'rounded-2xl overflow-hidden flex flex-col'
  const map: Record<IntroductionCardStyle, string> = {
    elevated: `${base} bg-card shadow-sm`,
    flat:     `${base} bg-muted/30`,
    bordered: `${base} border border-border/60`,
    glass:    `${base} glass`,
  }
  return map[style] ?? map.elevated
}
</script>

<template>
  <section style="padding-top: var(--hp-section-py, 96px); padding-bottom: var(--hp-section-py, 96px)">
    <div class="mx-auto px-4 sm:px-6 lg:px-8" style="max-width: var(--hp-max-width, 1200px)">
      <!-- 标题区 -->
      <div :class="titleAlignClass" class="mb-16">
        <h2 class="text-3xl sm:text-4xl font-bold tracking-tight">{{ config.title }}</h2>
        <p v-if="config.description" class="mt-4 text-lg text-muted-foreground max-w-2xl leading-relaxed" :class="descAlignClass">
          {{ config.description }}
        </p>
      </div>

      <!-- 特性网格 -->
      <div :class="gridClass">
        <div
          v-for="(item, i) in config.features" :key="i"
          :class="outerCardClass(cardStyle)"
        >
          <!-- 文字区 -->
          <div class="px-6 pt-6" :class="item.image_url ? 'pb-5' : 'pb-6'">
            <div v-if="item.icon" class="mb-4 flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <DynamicLucideIcon :name="item.icon" class="h-5 w-5" />
            </div>
            <h3 class="text-[15px] font-semibold leading-snug mb-2">{{ item.title }}</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">{{ item.description }}</p>
            <a
              v-if="item.link_text && item.link_url"
              :href="item.link_url"
              class="mt-3 inline-flex items-center gap-0.5 text-sm font-medium text-primary hover:underline underline-offset-2 transition-colors"
            >
              {{ item.link_text }}
              <ArrowUpRight class="h-3.5 w-3.5 shrink-0" />
            </a>
          </div>

          <!-- 内嵌预览区 -->
          <div v-if="item.image_url" class="px-4 pb-4 mt-auto">
            <div class="rounded-xl overflow-hidden bg-muted/50 dark:bg-muted/20 ring-1 ring-border/30">
              <img
                :src="item.image_url"
                :alt="item.title"
                class="w-full h-auto object-cover block"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
