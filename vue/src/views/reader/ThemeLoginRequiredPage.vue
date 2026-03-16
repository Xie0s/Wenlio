<!-- ThemeLoginRequiredPage.vue - 主题需要登录提示页
     职责：当主题 access_mode 为 'login' 且用户未登录时展示提示，引导登录
     对外接口：路由参数 tenantId, themeSlug -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReaderStore } from '@/stores/reader'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Lock, LogIn, ArrowLeft, Home, Loader2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const store = useReaderStore()

const tenantId = computed(() => route.params.tenantId as string)
const themeSlug = computed(() => route.params.themeSlug as string)
const themeName = computed(() => store.currentTheme?.name || themeSlug.value)
const ready = ref(false)

onMounted(async () => {
  if (!store.currentTheme || store.currentTheme.slug !== themeSlug.value) {
    await store.ensureThemesLoaded(tenantId.value)
    store.findThemeBySlug(themeSlug.value)
  }
  ready.value = true
})

function goLogin() {
  const themePath = `/${tenantId.value}/${themeSlug.value}`
  router.push(`/admin/login?redirect=${encodeURIComponent(themePath)}`)
}
function goGallery() {
  router.push(`/${tenantId.value}/themes`)
}
function goHome() {
  router.push(`/${tenantId.value}`)
}
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-background px-4">
    <div v-if="!ready" class="flex justify-center">
      <Loader2 class="size-6 animate-spin text-muted-foreground" />
    </div>
    <div v-else class="flex flex-col items-center gap-6 max-w-md text-center">
      <div class="rounded-full bg-blue-100 dark:bg-blue-900/30 p-5">
        <Lock class="size-10 text-blue-600 dark:text-blue-400" :stroke-width="1.5" />
      </div>

      <div class="space-y-2">
        <h1 class="text-2xl font-semibold tracking-tight">{{ themeName }}</h1>
        <p class="text-muted-foreground text-base leading-relaxed">
          该主题仅对登录用户可见，请登录后访问。
        </p>
      </div>

      <div class="flex items-center gap-3">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="outline" size="icon" class="rounded-full h-11 w-11" @click="goHome">
              <Home class="size-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>返回首页</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="outline" size="icon" class="rounded-full h-11 w-11" @click="goGallery">
              <ArrowLeft class="size-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>主题画廊</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button size="icon" class="rounded-full h-11 w-11" @click="goLogin">
              <LogIn class="size-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>前往登录</TooltipContent>
        </Tooltip>
      </div>
    </div>
  </div>
</template>
