<!-- TenantHomePage.vue - 租户文档首页
     职责：路由入口页面，优先渲染个性化首页（HomepageRenderer），
           未配置时回退到主题画廊页（ThemeGalleryPage） -->
<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HomepageRenderer from '@/components/personalization/renderer/HomepageRenderer.vue'
import ThemeGalleryPage from '@/views/reader/ThemeGalleryPage.vue'
import { applyTenantBrowserBranding } from '@/lib/browser-branding'
import { useReaderStore } from '@/stores/reader'
import { Loader2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const store = useReaderStore()

const tenantId = computed(() => route.params.tenantId as string)

// 是否存在个性化首页配置
const hasCustomHomepage = computed(() => !!store.homepageLayout)

async function init() {
  store.clearContext()
  // 并行加载：首页配置 + 主题列表（主题列表在两种模式下都需要）
  await Promise.all([
    store.loadTenant(tenantId.value),
    store.loadHomepage(tenantId.value),
    store.loadThemes(tenantId.value),
  ])
  if (!store.tenant) {
    router.replace({ name: 'NotFound' })
  }
}

onMounted(init)
watch(tenantId, init)
watch(
  () => [
    store.tenant?.name ?? '',
    store.tenant?.browser_title ?? '',
    store.tenant?.browser_icon_url ?? '',
  ],
  () => {
    applyTenantBrowserBranding({
      tenantName: store.tenant?.name,
      browserTitle: store.tenant?.browser_title,
      browserIconUrl: store.tenant?.browser_icon_url,
    })
  },
  { immediate: true },
)

function handleNavigate(path: string) {
  router.push(path)
}
</script>

<template>
  <!-- 加载中 -->
  <div v-if="store.loadingHomepage" class="flex h-screen items-center justify-center bg-background"
    :key="`ld-${tenantId}`">
    <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
  </div>

  <!-- 个性化首页 -->
  <HomepageRenderer v-else-if="hasCustomHomepage" class="h-screen" :key="`hp-${tenantId}`"
    :layout="store.homepageLayout!" :tenant-id="tenantId" :tenant-name="store.tenantName"
    :tenant-logo="store.tenant?.logo_url ?? ''" :themes="store.themes" @navigate="handleNavigate" />

  <!-- 经典首页（回退） -->
  <ThemeGalleryPage v-else :key="`gal-${tenantId}`" />
</template>
