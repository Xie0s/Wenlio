<!-- DocReaderPage.vue - 文档阅读页面
     职责：路由入口页面，组合 DocReaderLayout 布局与 DocPageViewer 内容展示
     对外接口：无（由路由直接加载） -->
<script setup lang="ts">
import { ref, provide, watchEffect, onMounted, onUnmounted, nextTick } from 'vue'
import { useDocReader } from '@/composables/useDocReader'
import DocReaderLayout from '@/components/preview/DocReaderLayout.vue'
import DocPageViewer from '@/components/preview/DocPageViewer.vue'
import { ArrowUp } from 'lucide-vue-next'

const {
  tenantId,
  themeSlug,
  versionName,
  pageSlug,
  navigateTo,
  navigateToPage,
  navigateToVersion,
  redirectTo,
} = useDocReader()

const eyeCareMode = ref(false)
provide('eyeCareMode', eyeCareMode)
provide('toggleEyeCare', () => { eyeCareMode.value = !eyeCareMode.value })

watchEffect(() => {
  document.body.classList.toggle('eye-care', eyeCareMode.value)
})

const showScrollTop = ref(false)

function getScrollContainer() {
  return document.getElementById('reader-scroll-container')
}

function onScroll() {
  showScrollTop.value = (getScrollContainer()?.scrollTop ?? 0) > 300
}

function scrollToTop() {
  getScrollContainer()?.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => nextTick(() => {
  getScrollContainer()?.addEventListener('scroll', onScroll, { passive: true })
}))
onUnmounted(() => {
  getScrollContainer()?.removeEventListener('scroll', onScroll)
  document.body.classList.remove('eye-care')
})
</script>

<template>
  <DocReaderLayout :tenant-id="tenantId" @navigate="navigateTo">
    <DocPageViewer :tenant-id="tenantId" :theme-slug="themeSlug" :version-name="versionName" :page-slug="pageSlug"
      @navigate="navigateToPage" @version-change="navigateToVersion" @redirect="redirectTo" />
  </DocReaderLayout>

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
</template>
