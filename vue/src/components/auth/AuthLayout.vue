<!--
  AuthLayout.vue - 认证页面通用布局组件

  职责：统一承载登录/注册页面的卡片布局、标题、副标题与底部扩展内容。
  对外暴露：`title`、`description`、`cardClass` 属性，以及默认插槽与 `footer` 插槽。
-->
<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'
import { ShieldCheck, Undo2 } from 'lucide-vue-next'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useThemeStore } from '@/stores/theme'
import { RouterLink } from 'vue-router'

const Antigravity = defineAsyncComponent(() => import('@/components/ui/a-vue-bits/antigravity/Antigravity.vue'))

const props = withDefaults(defineProps<{
  title: string
  description: string
  cardClass?: string
  mode?: 'login' | 'register'
}>(), {
  cardClass: 'max-w-[480px]',
  mode: 'login',
})

const themeStore = useThemeStore()

const isDark = computed(() => themeStore.resolvedTheme === 'dark')
const antigravityColor = computed(() => (isDark.value ? '#ffffff' : '#2563eb'))
const modeLabel = computed(() => (props.mode === 'register' ? '管理员注册' : '后台登录'))
const panelHeightClass = computed(() => (props.mode === 'register' ? 'min-h-0' : 'min-h-[480px]'))
</script>

<template>
  <div class="relative h-svh overflow-hidden bg-background">
    <div class="absolute inset-0">
      <Antigravity
        :count="300"
        :magnetRadius="13"
        :ringRadius="8"
        :waveSpeed="0.5"
        :waveAmplitude="1"
        :particleSize="1"
        :lerpSpeed="0.15"
        :color="antigravityColor"
        :autoAnimate="true"
        :particleVariance="1"
        :rotationSpeed="0"
        :depthFactor="1"
        :pulseSpeed="3"
        particleShape="capsule"
        :fieldStrength="10"
      />
    </div>

    <div class="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,oklch(1_0_0_/_0.52),transparent_38%),linear-gradient(180deg,transparent,oklch(1_0_0_/_0.18))] dark:bg-[radial-gradient(circle_at_top,oklch(0.22_0_0_/_0.4),transparent_34%),linear-gradient(180deg,transparent,oklch(0.06_0_0_/_0.55))]" />
  
    <div class="relative z-20 h-full overflow-y-auto">
    <div class="mx-auto flex min-h-full max-w-[1200px] flex-col px-4 py-5 md:px-6 md:py-6">
      <div class="flex flex-1 items-center justify-center">
        <div :class="['glass relative w-full overflow-hidden rounded-3xl border border-white/30 transition-[max-width] duration-500 ease-[cubic-bezier(0.22,1,0.36,1)] dark:border-white/10', cardClass]">
          <div :class="['relative transition-[min-height] duration-500 ease-[cubic-bezier(0.22,1,0.36,1)]', panelHeightClass]">

            <div
              class="relative z-10 flex w-full flex-col px-5 py-4 sm:px-6 sm:py-5 lg:px-7 lg:py-6"
            >
              <div class="flex items-center justify-between gap-3">
                <Transition name="auth-chip" mode="out-in">
                  <div :key="modeLabel" class="inline-flex items-center gap-2 rounded-full border border-border/60 bg-background/55 px-3 py-1 text-xs font-medium text-muted-foreground backdrop-blur-md">
                    <ShieldCheck class="h-3.5 w-3.5" />
                    <span>{{ modeLabel }}</span>
                  </div>
                </Transition>
                <TooltipProvider>
                  <div class="flex items-center gap-2">
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <ThemeToggle
                          button-size="size-8"
                          class="border border-border/60 bg-background/55 text-foreground backdrop-blur-md hover:bg-muted/60"
                          icon-size="h-4 w-4"
                        />
                      </TooltipTrigger>
                      <TooltipContent>{{ isDark ? '切换到浅色模式' : '切换到深色模式' }}</TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <RouterLink
                          to="/"
                          class="inline-flex h-8 w-8 items-center justify-center rounded-full border border-border/60 bg-background/55 text-foreground backdrop-blur-md transition-colors hover:bg-muted/60"
                        >
                          <Undo2 class="h-4 w-4" />
                        </RouterLink>
                      </TooltipTrigger>
                      <TooltipContent>返回首页</TooltipContent>
                    </Tooltip>
                  </div>
                </TooltipProvider>
              </div>

              <div class="mt-4 flex flex-1 flex-col">
                <Transition name="auth-copy" mode="out-in">
                  <div :key="props.mode" class="space-y-1.5">
                    <h1 class="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">{{ title }}</h1>
                    <p class="text-sm leading-6 text-muted-foreground">{{ description }}</p>
                  </div>
                </Transition>
                <div class="mt-4">
                  <slot />
                </div>
              </div>

              <div class="mt-5 border-t border-border/40 pt-3">
                <slot name="footer" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <footer class="mx-auto mt-6 flex w-full max-w-[720px] flex-col items-center gap-1 px-3 text-center text-sm font-light leading-6 text-muted-foreground/75">
        <div class="flex flex-wrap items-center justify-center gap-x-2 gap-y-0.5">
          <a
            href="https://beian.miit.gov.cn/"
            target="_blank"
            rel="noopener noreferrer"
            class="transition-colors hover:text-foreground"
          >陇ICP备20002844号-1</a>
          <span class="text-muted-foreground/40">|</span>
          <a
            href="https://www.beian.gov.cn/"
            target="_blank"
            rel="noopener noreferrer"
            class="transition-colors hover:text-foreground"
          >甘公网安备62090202000540号</a>
        </div>
        <div class="flex flex-wrap items-center justify-center gap-x-2 text-base font-light text-foreground/70">
          <span>Microswift Core™</span>
          <span class="text-muted-foreground/35">·</span>
          <span>微讯云信息</span>
        </div>
        <a
          href="https://microswift.cn"
          target="_blank"
          rel="noopener noreferrer"
          class="text-sm font-light text-muted-foreground/75 transition-colors hover:text-foreground"
        >版权所有 ©2020-2026</a>
      </footer>
    </div>
    </div>
  </div>
</template>

<style scoped>
.auth-chip-enter-active,
.auth-chip-leave-active {
  transition: opacity 0.24s ease, transform 0.24s ease;
}

.auth-chip-enter-from,
.auth-chip-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.auth-copy-enter-active,
.auth-copy-leave-active {
  transition: opacity 0.32s ease, transform 0.32s ease;
}

.auth-copy-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.auth-copy-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
