<!--
  src/views/home.vue
  职责：文档管理平台首页，组合 HomeNav / HomeHero / HomeFeatures / HomeWorkflow / HomeCTA / HomeFooter
  对外暴露：无（作为路由组件页面）
-->
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ArrowUp } from 'lucide-vue-next'
import HomeNav from '@/components/home/HomeNav.vue'
import HomeHero from '@/components/home/HomeHero.vue'
import HomeFeatures from '@/components/home/HomeFeatures.vue'
import HomeStructure from '@/components/home/HomeStructure.vue'
import HomeWorkflow from '@/components/home/HomeWorkflow.vue'
import HomeMoreFeatures from '@/components/home/HomeMoreFeatures.vue'
import HomePersonalization from '@/components/home/HomePersonalization.vue'
import HomeCTA from '@/components/home/HomeCTA.vue'
import HomeFooter from '@/components/home/HomeFooter.vue'

const showScrollTop = ref(false)

function onScroll() {
  showScrollTop.value = window.scrollY > 300
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <div class="min-h-screen bg-[#f5f3ef] dark:bg-[#1a1a1a]">
    <HomeNav />
    <main class="space-y-10">
      <HomeHero />
      <HomeFeatures />
      <HomeStructure />
      <HomeWorkflow />
      <HomeMoreFeatures />
      <HomePersonalization />
      <HomeCTA />
    </main>
    <HomeFooter />

    <!-- 回到顶部按钮 -->
    <Transition enter-active-class="transition-all duration-200" enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0" leave-active-class="transition-all duration-200"
      leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-2">
      <button v-if="showScrollTop"
        class="fixed bottom-6 right-6 z-50 flex items-center justify-center h-11 w-11 rounded-full border border-foreground/20 bg-background/80 backdrop-blur-md text-foreground transition-colors hover:bg-foreground/5"
        @click="scrollToTop">
        <ArrowUp class="h-5 w-5" />
      </button>
    </Transition>
  </div>
</template>
