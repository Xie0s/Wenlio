<!--
  src/components/home/HomeNav.vue
  职责：首页顶部导航栏
  对外暴露：无
-->
<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const scrolled = ref(false)
const activeSection = ref('')
const navRef = ref<HTMLElement | null>(null)
const pill = ref({ left: 0, width: 0, visible: false })

const navItems = [
  { id: 'features', label: '核心能力' },
  { id: 'structure', label: '组织结构' },
  { id: 'workflow', label: '发布流程' },
  { id: 'more-features', label: '功能总览' },
  { id: 'personalization', label: '个性化组件' },
]

function onScroll() {
  scrolled.value = window.scrollY > 20
}

let observer: IntersectionObserver | null = null

function movePillToEl(el: HTMLElement) {
  if (!navRef.value) return
  const navRect = navRef.value.getBoundingClientRect()
  const rect = el.getBoundingClientRect()
  pill.value = { left: rect.left - navRect.left, width: rect.width, visible: true }
}

function onLinkEnter(e: MouseEvent) {
  movePillToEl(e.currentTarget as HTMLElement)
}

function snapToActive() {
  nextTick(() => {
    const el = navRef.value?.querySelector<HTMLElement>(`[data-section="${activeSection.value}"]`)
    if (el) movePillToEl(el)
    else pill.value.visible = false
  })
}

function onNavLeave() {
  snapToActive()
}

watch(activeSection, snapToActive)

onMounted(() => {
  window.addEventListener('scroll', onScroll)

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) activeSection.value = entry.target.id
      })
    },
    { rootMargin: '-64px 0px -55% 0px', threshold: 0 }
  )

  navItems.forEach(({ id }) => {
    const el = document.getElementById(id)
    if (el) observer!.observe(el)
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  observer?.disconnect()
})
</script>

<template>
  <header class="fixed top-0 left-0 right-0 z-50 transition-all duration-300"
    :class="scrolled ? 'bg-[#f5f3ef]/80 dark:bg-[#1a1a1a]/75 backdrop-blur-[16px] backdrop-saturate-[160%]' : 'bg-transparent'">
    <div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
      <!-- Brand -->
      <a href="/" class="flex items-center gap-2">
        <span class="text-xl font-light tracking-tight">Wenlio 文流</span>
      </a>

      <!-- Nav Links -->
      <nav ref="navRef" class="relative hidden items-center gap-1 md:flex" @mouseleave="onNavLeave">
        <!-- 滑动背景胶囊 -->
        <span
          class="pointer-events-none absolute -inset-y-0.5 rounded-full bg-foreground/[0.07] transition-all duration-200 ease-out"
          :style="{ left: `${pill.left}px`, width: `${pill.width}px`, opacity: pill.visible ? 1 : 0 }" />
        <a v-for="item in navItems" :key="item.id" :href="`#${item.id}`" :data-section="item.id"
          class="relative rounded-full px-3 py-1.5 text-base font-light text-foreground transition-colors duration-200"
          @mouseenter="onLinkEnter">{{ item.label }}</a>
      </nav>

      <!-- CTA -->
      <div class="flex items-center gap-3">
        <a href="https://microswift.cn" target="_blank" rel="noopener noreferrer"
          class="text-base font-light text-foreground transition-opacity hover:opacity-80">
          关于我们
        </a>
        <button v-if="authStore.isLoggedIn"
          class="text-base font-light text-foreground transition-opacity hover:opacity-80"
          @click="router.push('/admin')">进入后台</button>
        <router-link v-else to="/admin/login"
          class="text-base font-light text-foreground transition-opacity hover:opacity-80">登录</router-link>
        <ThemeToggle />
      </div>
    </div>
  </header>
</template>
