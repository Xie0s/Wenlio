<!--
  src/components/personalization/hero/HeroSplit.vue
  职责：分栏布局 Hero 展示组件，支持左右图文切换、标题动画与按钮交互
  对外暴露：Props(config)、Emits(navigate)
-->
<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch, useTemplateRef } from 'vue'
import type { HeroConfig, ButtonConfig } from '../types'
import { ArrowRight } from 'lucide-vue-next'
import NoiseBackground from '@/components/ui/a-aceternity/NoiseBackground.vue'

const props = defineProps<{
  config: HeroConfig
}>()

const emit = defineEmits<{
  navigate: [url: string]
}>()

const reversed = computed(() => props.config.layout === 'right-left')

// ─── 打字机动画 ────────────────────────────────────────────────
const displayedTitle = ref('')
const isVisible = ref(false)
let typeTimer: ReturnType<typeof setTimeout> | null = null
let observer: IntersectionObserver | null = null
const sectionRef = useTemplateRef<HTMLElement>('sectionRef')

function stopTypewriter() {
  if (typeTimer) { clearTimeout(typeTimer); typeTimer = null }
}

function startTypewriter(text: string) {
  stopTypewriter()
  displayedTitle.value = ''
  let i = 0
  let deleting = false

  function tick() {
    if (!deleting) {
      if (i < text.length) {
        displayedTitle.value += text[i++]
        typeTimer = setTimeout(tick, 70 + Math.random() * 40)
      } else {
        typeTimer = setTimeout(() => { deleting = true; tick() }, 2000)
      }
    } else {
      if (displayedTitle.value.length > 0) {
        displayedTitle.value = displayedTitle.value.slice(0, -1)
        typeTimer = setTimeout(tick, 35 + Math.random() * 20)
      } else {
        i = 0; deleting = false
        typeTimer = setTimeout(tick, 800)
      }
    }
  }
  typeTimer = setTimeout(tick, 300)
}

const isTypewriter = computed(() => props.config.animation === 'typewriter')

const titleText = computed(() =>
  isTypewriter.value ? displayedTitle.value : props.config.title
)

const animationClass = computed(() => {
  if (isTypewriter.value) return ''
  const map: Record<string, string> = { 'fade-up': 'animate-fade-up', 'fade-in': 'animate-fade-in' }
  return map[props.config.animation] ?? ''
})

// ─── 标题混合色渲染 ──────────────────────────────────────────
const titleParts = computed(() => {
  const full = titleText.value
  const hl = props.config.highlight_text
  if (!hl || !full.includes(hl)) return [{ text: full, muted: false }]
  const idx = full.indexOf(hl)
  const parts: { text: string; muted: boolean }[] = []
  if (idx > 0) parts.push({ text: full.slice(0, idx), muted: false })
  parts.push({ text: hl, muted: true })
  if (idx + hl.length < full.length) parts.push({ text: full.slice(idx + hl.length), muted: false })
  return parts
})

onMounted(() => {
  if (sectionRef.value) {
    observer = new IntersectionObserver((entries) => {
      if (entries[0]) isVisible.value = entries[0].isIntersecting
    }, { threshold: 0.1 })
    observer.observe(sectionRef.value)
  }
})

onUnmounted(() => {
  stopTypewriter()
  if (observer) { observer.disconnect(); observer = null }
})

watch(isVisible, (visible) => {
  if (visible && isTypewriter.value) startTypewriter(props.config.title)
  else if (!visible) stopTypewriter()
})

watch(
  () => props.config.title,
  (val) => { if (isTypewriter.value && isVisible.value) startTypewriter(val) },
)
watch(
  () => props.config.animation,
  (val) => {
    if (val === 'typewriter' && isVisible.value) startTypewriter(props.config.title)
    else { stopTypewriter(); displayedTitle.value = '' }
  },
)

// ─── 按钮处理 ──────────────────────────────────────────────────
function handleButton(btn: ButtonConfig) {
  if (btn.url.startsWith('#')) {
    document.querySelector(btn.url)?.scrollIntoView({ behavior: 'smooth' })
    return
  }
  emit('navigate', btn.url)
}

const alignClass = computed(() => {
  const map: Record<string, string> = { left: 'text-left', center: 'text-center', right: 'text-right' }
  return map[props.config.title_align] ?? 'text-left'
})

const btnAlignClass = computed(() => {
  const map: Record<string, string> = { left: 'justify-start', center: 'justify-center', right: 'justify-end' }
  return map[props.config.title_align] ?? 'justify-start'
})

function btnClass(variant: string): string {
  const base = 'inline-flex items-center justify-center gap-2 rounded-full text-sm font-medium transition-all duration-200 cursor-pointer active:scale-[0.98]'
  const map: Record<string, string> = {
    dark:      `${base} bg-foreground text-background px-6 py-3 hover:opacity-90`,
    primary:   `${base} bg-primary text-primary-foreground px-6 py-3 hover:bg-primary/90 shadow-lg shadow-primary/20`,
    secondary: `${base} bg-secondary text-secondary-foreground px-6 py-3 hover:bg-secondary/80`,
    outline:   `${base} border border-foreground/20 px-6 py-3 hover:border-foreground/40 hover:bg-foreground/5`,
    plain:     `${base} text-foreground/70 hover:text-foreground underline-offset-4 hover:underline`,
  }
  return map[variant] ?? map.dark!
}
</script>

<template>
  <section ref="sectionRef" class="relative overflow-hidden">
    <div
      class="mx-auto grid items-center gap-12 px-6 sm:px-8 lg:grid-cols-2 lg:gap-20 lg:px-10"
      :class="animationClass"
      style="max-width: var(--hp-max-width, 1200px); padding-top: var(--hp-section-py, 96px); padding-bottom: var(--hp-section-py, 96px)"
    >

      <!-- 文字区 -->
      <div :class="[{ 'lg:order-2': reversed }, alignClass]">
        <!-- 副标题 -->
        <div v-if="config.subtitle" class="mb-5 inline-flex items-center">
          <span class="rounded-full border border-foreground/10 bg-foreground/5 px-3.5 py-1 text-xs font-medium tracking-wide text-foreground/60">
            {{ config.subtitle }}
          </span>
        </div>

        <!-- 主标题 -->
        <h1 class="text-3xl sm:text-4xl lg:text-5xl font-bold tracking-tight leading-[1.1]">
          <template v-for="(part, pi) in titleParts" :key="pi">
            <span :class="part.muted ? 'text-muted-foreground/50' : ''">{{ part.text }}</span>
          </template>
          <span
            v-if="isTypewriter"
            class="cursor-google ml-0.5 inline-block w-[3px] rounded-full align-middle"
            :style="{ height: '0.85em', animationPlayState: isVisible ? 'running' : 'paused' }"
          />
        </h1>

        <!-- 描述 -->
        <p
          v-if="config.description"
          class="mt-6 text-lg text-muted-foreground leading-relaxed"
          :class="isTypewriter ? 'animate-fade-in' : ''"
          :style="isTypewriter ? 'animation-delay: 0.8s; animation-fill-mode: both' : ''"
        >
          {{ config.description }}
        </p>

        <!-- 按钮组 -->
        <div
          v-if="config.primary_button || config.secondary_button"
          class="mt-10 flex flex-wrap items-center gap-4"
          :class="[isTypewriter ? 'animate-fade-in' : '', btnAlignClass]"
          :style="isTypewriter ? 'animation-delay: 1.2s; animation-fill-mode: both' : ''"
        >
          <!-- 主按钮 - noise -->
          <NoiseBackground
            v-if="config.primary_button?.variant === 'noise'"
            :gradient-colors="['rgb(255,100,150)', 'rgb(100,150,255)', 'rgb(255,200,100)']"
            :noise-intensity="0.18"
            :speed="0.08"
            container-class="w-fit !rounded-full !p-1.5"
          >
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-full bg-linear-to-r from-neutral-100 via-neutral-100 to-white px-6 py-3 text-sm font-medium text-black shadow-[0px_2px_0px_0px_var(--color-neutral-50)_inset,0px_0.5px_1px_0px_var(--color-neutral-400)] transition-all duration-100 active:scale-[0.98] dark:from-neutral-900 dark:via-neutral-900 dark:to-black dark:text-white dark:shadow-[0px_1px_0px_0px_var(--color-neutral-950)_inset,0px_1px_0px_0px_var(--color-neutral-800)]"
              @click="handleButton(config.primary_button!)"
            >
              {{ config.primary_button!.text }}
              <ArrowRight v-if="config.primary_button!.show_arrow" class="h-4 w-4" />
            </button>
          </NoiseBackground>

          <button
            v-else-if="config.primary_button"
            :class="btnClass(config.primary_button.variant)"
            @click="handleButton(config.primary_button)"
          >
            {{ config.primary_button.text }}
            <ArrowRight v-if="config.primary_button.show_arrow" class="h-4 w-4" />
          </button>

          <!-- 次按钮 - noise -->
          <NoiseBackground
            v-if="config.secondary_button?.variant === 'noise'"
            :gradient-colors="['rgb(255,100,150)', 'rgb(100,150,255)', 'rgb(255,200,100)']"
            :noise-intensity="0.18"
            :speed="0.08"
            container-class="w-fit !rounded-full !p-1.5"
          >
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-full bg-linear-to-r from-neutral-100 via-neutral-100 to-white px-6 py-3 text-sm font-medium text-black shadow-[0px_2px_0px_0px_var(--color-neutral-50)_inset,0px_0.5px_1px_0px_var(--color-neutral-400)] transition-all duration-100 active:scale-[0.98] dark:from-neutral-900 dark:via-neutral-900 dark:to-black dark:text-white dark:shadow-[0px_1px_0px_0px_var(--color-neutral-950)_inset,0px_1px_0px_0px_var(--color-neutral-800)]"
              @click="handleButton(config.secondary_button!)"
            >
              {{ config.secondary_button!.text }}
              <ArrowRight v-if="config.secondary_button!.show_arrow" class="h-4 w-4" />
            </button>
          </NoiseBackground>

          <button
            v-else-if="config.secondary_button"
            :class="btnClass(config.secondary_button.variant)"
            @click="handleButton(config.secondary_button)"
          >
            {{ config.secondary_button.text }}
            <ArrowRight v-if="config.secondary_button.show_arrow" class="h-4 w-4" />
          </button>
        </div>
      </div>

      <!-- 配图区 -->
      <div :class="{ 'lg:order-1': reversed }">
        <div class="relative aspect-[4/3] overflow-hidden rounded-2xl bg-muted/30">
          <img
            v-if="config.image_url"
            :src="config.image_url"
            :alt="config.title"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full w-full flex-col items-center justify-center gap-3 text-muted-foreground/30">
            <svg class="h-14 w-14" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="0.8">
              <rect x="3" y="3" width="18" height="18" rx="4" />
              <circle cx="8.5" cy="8.5" r="1.5" />
              <path d="m21 15-5-5L5 21" />
            </svg>
            <span class="text-xs">在编辑器中填写配图 URL</span>
          </div>
        </div>
      </div>

    </div>
  </section>
</template>

<style scoped>
.cursor-google {
  animation:
    google-colors 2s linear infinite,
    cursor-blink 1.06s step-end infinite;
}

@keyframes google-colors {
  0%   { background-color: #4285F4; }
  25%  { background-color: #EA4335; }
  50%  { background-color: #FBBC04; }
  75%  { background-color: #34A853; }
  100% { background-color: #4285F4; }
}

@keyframes cursor-blink {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0; }
}
</style>
