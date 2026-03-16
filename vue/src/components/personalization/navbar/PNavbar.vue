<!--
  src/components/personalization/navbar/PNavbar.vue
  职责：个性化首页导航栏展示组件，处理品牌展示、导航交互与移动端菜单
  对外暴露：Props(config, tenantId, tenantName, tenantLogo)、Emits(navigate)
-->
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { NavbarConfig, NavLink, ButtonConfig } from '../types'
import { Menu, X, ChevronDown, ArrowRight } from 'lucide-vue-next'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { useAuthStore } from '@/stores/auth'
import UserMenu from '@/components/auth/UserMenu.vue'

const props = defineProps<{
  config: NavbarConfig
  tenantId: string
  tenantName: string
  tenantLogo: string
}>()

const emit = defineEmits<{
  navigate: [url: string]
}>()

const mobileOpen = ref(false)
const activeDropdown = ref<number | null>(null)
const scrolled = ref(false)
let blurTimer: ReturnType<typeof setTimeout> | null = null
const authStore = useAuthStore()

const navRef = ref<HTMLElement | null>(null)
const hoveredLink = ref<{ left: number; width: number } | null>(null)
const pillPos = ref<{ left: number; width: number } | null>(null)
const dropdownBtnRefs = ref<(HTMLElement | null)[]>([])

function onNavLinkEnter(e: MouseEvent) {
  const nav = navRef.value
  if (!nav) return
  const el = e.currentTarget as HTMLElement
  const navRect = nav.getBoundingClientRect()
  const elRect = el.getBoundingClientRect()
  const pos = { left: elRect.left - navRect.left, width: elRect.width }
  hoveredLink.value = pos
  pillPos.value = pos
}

function onDropdownPanelEnter(i: number) {
  const nav = navRef.value
  const btn = dropdownBtnRefs.value[i]
  if (!nav || !btn) return
  const navRect = nav.getBoundingClientRect()
  const btnRect = btn.getBoundingClientRect()
  const pos = { left: btnRect.left - navRect.left, width: btnRect.width }
  hoveredLink.value = pos
  pillPos.value = pos
}

function onNavLeave() {
  hoveredLink.value = null
}

function onDropdownBlur() {
  blurTimer = setTimeout(() => { activeDropdown.value = null }, 120)
}

function onDropdownMousedown() {
  if (blurTimer) { clearTimeout(blurTimer); blurTimer = null }
}

const displayLogo = computed(() => props.config.logo_url || props.tenantLogo)
const displayBrand = computed(() => props.config.brand_text || props.tenantName)
const brandMode = computed(() => props.config.brand_mode || 'both')

const showLogo = computed(() => brandMode.value === 'both' || brandMode.value === 'logo_only')
const showText = computed(() => brandMode.value === 'both' || brandMode.value === 'text_only')
const showAdminEntry = computed(() =>
  authStore.isLoggedIn && authStore.user?.tenant_id === props.tenantId,
)

function hpBgMix(pct: string) {
  return { backgroundColor: `color-mix(in srgb, var(--hp-bg, hsl(var(--background))) ${pct}, transparent)` }
}

const panelBgStyle = hpBgMix('95%')

const navbarClass = computed(() => {
  const base = 'w-full z-50 transition-all duration-300'
  const sticky = props.config.sticky ? 'sticky top-0' : 'relative'
  const blur = props.config.style === 'blur' ? 'backdrop-blur-xl' : ''
  const shadow = scrolled.value && props.config.sticky ? 'shadow-[0_1px_0_0_rgba(128,128,128,0.1)]' : ''
  return `${base} ${sticky} ${blur} ${shadow}`
})

const navbarBgStyle = computed(() => {
  const s = props.config.style
  if (s === 'transparent') return {}
  return hpBgMix(s === 'solid' ? '95%' : '70%')
})

function onScroll() {
  scrolled.value = window.scrollY > 4
}

function isExternalUrl(url: string): boolean {
  return /^(https?:)?\/\//i.test(url)
}

onMounted(() => {
  if (props.config.sticky) {
    window.addEventListener('scroll', onScroll, { passive: true })
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})

function handleLink(link: NavLink) {
  mobileOpen.value = false
  activeDropdown.value = null
  if (link.url.startsWith('#')) {
    const el = document.querySelector(link.url)
    el?.scrollIntoView({ behavior: 'smooth' })
    return
  }
  if (link.external || isExternalUrl(link.url)) {
    window.open(link.url, '_blank', 'noopener')
    return
  }
  emit('navigate', link.url)
}

function handleButton(btn: ButtonConfig) {
  if (btn.url.startsWith('#')) {
    document.querySelector(btn.url)?.scrollIntoView({ behavior: 'smooth' })
    return
  }
  if (isExternalUrl(btn.url)) {
    window.open(btn.url, '_blank', 'noopener')
    return
  }
  emit('navigate', btn.url)
}

function goToThemes() {
  emit('navigate', `/${props.tenantId}/themes`)
}

function toggleDropdown(index: number) {
  activeDropdown.value = activeDropdown.value === index ? null : index
}
</script>

<template>
  <nav :class="navbarClass" :style="navbarBgStyle">
    <div class="mx-auto flex items-center justify-between px-4 py-2 sm:px-6 lg:px-8"
      style="max-width: var(--hp-max-width, 1200px)">

      <!-- 品牌区域 -->
      <a class="flex items-center gap-2 shrink-0 cursor-pointer hover:opacity-80 transition-opacity"
        @click.prevent="emit('navigate', '')">
        <img v-if="showLogo && displayLogo" :src="displayLogo" :alt="displayBrand"
          class="h-[52px] w-auto max-w-[260px] rounded object-contain" />
        <span v-if="showText" class="text-xl font-normal tracking-tight">{{ displayBrand }}</span>
      </a>

      <!-- 桌面端链接 -->
      <div ref="navRef" class="hidden md:flex items-center gap-0.5 relative" @mouseleave="onNavLeave">
        <!-- 滑动高亮 pill -->
        <span v-if="pillPos"
          class="absolute inset-y-1 rounded-full bg-black/10 dark:bg-white/10 pointer-events-none transition-[left,width,opacity] duration-200 ease-out"
          :class="hoveredLink ? 'opacity-100' : 'opacity-0'"
          :style="{ left: `${pillPos.left}px`, width: `${pillPos.width}px` }" />
        <template v-for="(link, i) in config.links" :key="i">
          <div v-if="link.children?.length" class="relative">
            <button :ref="(el) => dropdownBtnRefs[i] = el as HTMLElement"
              class="flex items-center gap-1 px-4 py-[12px] text-xl font-normal text-foreground/60 hover:text-foreground rounded-full transition-colors relative z-10"
              @click="toggleDropdown(i)" @blur="onDropdownBlur" @mouseenter="onNavLinkEnter">
              {{ link.label }}
              <ChevronDown class="h-3 w-3 transition-transform" :class="{ 'rotate-180': activeDropdown === i }" />
            </button>
            <Transition enter-active-class="transition duration-150" enter-from-class="opacity-0 -translate-y-1"
              enter-to-class="opacity-100 translate-y-0" leave-active-class="transition duration-100"
              leave-from-class="opacity-100" leave-to-class="opacity-0">
              <div v-if="activeDropdown === i"
                class="absolute top-full left-0 mt-1.5 min-w-[180px] rounded-3xl border backdrop-blur-xl p-1.5 shadow-lg shadow-black/5"
                :style="panelBgStyle" @mousedown="onDropdownMousedown" @mouseenter="onDropdownPanelEnter(i)">
                <a v-for="(child, ci) in link.children" :key="ci"
                  class="block px-3 py-2 text-xl font-normal text-foreground/70 hover:text-foreground hover:bg-black/10 dark:hover:bg-white/10 rounded-3xl cursor-pointer transition-colors"
                  @click="handleLink(child)">{{ child.label }}</a>
              </div>
            </Transition>
          </div>
          <a v-else
            class="px-4 py-[12px] text-xl font-normal text-foreground/60 hover:text-foreground rounded-full cursor-pointer transition-colors relative z-10"
            @click.prevent="handleLink(link)" @mouseenter="onNavLinkEnter">{{ link.label }}</a>
        </template>
      </div>

      <!-- 右侧区域：CTA + 移动端菜单 -->
      <div class="flex items-center gap-2">
        <!-- CTA 按钮（桌面端） -->
        <button v-if="config.cta_button"
          class="hidden md:inline-flex items-center gap-1.5 rounded-full bg-foreground text-background px-4 py-2.5 text-sm font-medium transition-all hover:opacity-90 active:scale-[0.98]"
          @click="handleButton(config.cta_button)">
          {{ config.cta_button.text }}
          <ArrowRight v-if="config.cta_button.show_arrow" class="h-3.5 w-3.5" />
        </button>

        <!-- 文档主题（始终显示） -->
        <button
          class="hidden md:inline-flex items-center rounded-full border border-foreground/25 bg-transparent text-foreground px-4 py-2.5 text-sm font-medium transition-colors hover:bg-foreground/5"
          @click="goToThemes">
          文档主题
        </button>

        <!-- 用户菜单（仅当前租户已登录时显示） -->
        <UserMenu v-if="showAdminEntry" :tenant-id="tenantId" class="hidden md:flex" />

        <!-- 主题切换开关 -->
        <ThemeToggle v-if="config.show_theme_toggle" button-size="size-10"
          class="hidden md:flex border border-foreground/20" />

        <!-- 移动端汉堡按钮 -->
        <button class="md:hidden p-2 rounded-full hover:bg-muted/60 transition-colors"
          @click="mobileOpen = !mobileOpen">
          <X v-if="mobileOpen" class="h-5 w-5" />
          <Menu v-else class="h-5 w-5" />
        </button>
      </div>
    </div>

    <!-- 移动端菜单 -->
    <Transition enter-active-class="transition duration-200" enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0" leave-active-class="transition duration-150"
      leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="mobileOpen" class="md:hidden border-t backdrop-blur-xl px-4 py-3 space-y-1" :style="panelBgStyle">
        <template v-for="(link, i) in config.links" :key="i">
          <a class="block px-3 py-2.5 text-xl font-normal text-foreground/60 hover:text-foreground rounded-xl cursor-pointer transition-colors"
            @click.prevent="handleLink(link)">{{ link.label }}</a>
          <template v-if="link.children?.length">
            <a v-for="(child, ci) in link.children" :key="`${i}-${ci}`"
              class="block pl-7 py-2 text-xl font-normal text-foreground/40 hover:text-foreground rounded-xl cursor-pointer transition-colors"
              @click.prevent="handleLink(child)">{{ child.label }}</a>
          </template>
        </template>

        <!-- 主题切换开关（移动端） -->
        <div v-if="config.show_theme_toggle" class="flex items-center justify-between px-3 py-2.5 rounded-xl border">
          <span class="text-sm text-foreground/60">主题切换</span>
          <ThemeToggle />
        </div>

        <!-- 文档主题（始终显示） -->
        <button
          class="w-full flex items-center justify-center rounded-full border border-foreground/25 bg-transparent text-foreground px-4 py-3 text-sm font-medium transition-colors hover:bg-foreground/5"
          @click="goToThemes">
          文档主题
        </button>


        <!-- CTA 按钮（移动端） -->
        <button v-if="config.cta_button"
          class="mt-2 w-full flex items-center justify-center gap-1.5 rounded-full bg-foreground text-background px-4 py-3 text-sm font-medium"
          @click="handleButton(config.cta_button)">
          {{ config.cta_button.text }}
          <ArrowRight v-if="config.cta_button.show_arrow" class="h-3.5 w-3.5" />
        </button>
      </div>
    </Transition>
  </nav>
</template>
