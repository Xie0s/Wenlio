<!-- ThemeVerifyCodePage.vue - 主题验证码输入页面
     职责：当主题 access_mode 为 'code' 时，用户需输入6位验证码才能访问
     对外接口：路由参数 tenantId, themeSlug；验证成功后跳转到主题首页 -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReaderStore } from '@/stores/reader'
import { http, hasToken } from '@/utils/http'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { KeyRound, Check, ArrowLeft, Home, Loader2 } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const route = useRoute()
const router = useRouter()
const store = useReaderStore()

const tenantId = computed(() => route.params.tenantId as string)
const themeSlug = computed(() => route.params.themeSlug as string)
const themeName = computed(() => store.currentTheme?.name || themeSlug.value)
const ready = ref(false)

const code = ref('')
const verifying = ref(false)
const error = ref('')

onMounted(async () => {
  if (!store.currentTheme || store.currentTheme.slug !== themeSlug.value) {
    await store.ensureThemesLoaded(tenantId.value)
    store.findThemeBySlug(themeSlug.value)
  }
  const theme = store.currentTheme
  if (theme) {
    // 已登录用户直接跳过验证码
    if (hasToken()) {
      router.replace(`/${tenantId.value}/${themeSlug.value}`)
      return
    }
    const storedToken = localStorage.getItem(`theme_access_${theme.id}`)
    if (storedToken) {
      router.replace(`/${tenantId.value}/${themeSlug.value}`)
      return
    }
  }
  ready.value = true
})

async function handleSubmit() {
  if (code.value.length !== 6) {
    error.value = '请输入6位验证码'
    return
  }

  const theme = store.currentTheme
  if (!theme) {
    toast.error('主题信息未加载')
    return
  }

  verifying.value = true
  error.value = ''

  const res = await http.post<{ access_token: string }>(`/public/themes/${theme.id}/verify-code`, {
    code: code.value,
  })

  verifying.value = false

  if (res.code === 0 && res.data?.access_token) {
    localStorage.setItem(`theme_access_${theme.id}`, res.data.access_token)
    toast.success('验证成功')
    router.replace(`/${tenantId.value}/${themeSlug.value}`)
  } else {
    error.value = res.message || '验证码错误'
  }
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
    <div v-else class="flex flex-col items-center gap-6 max-w-sm w-full text-center">
      <div class="rounded-full bg-amber-100 dark:bg-amber-900/30 p-5">
        <KeyRound class="size-10 text-amber-600 dark:text-amber-400" :stroke-width="1.5" />
      </div>

      <div class="space-y-2">
        <h1 class="text-2xl font-semibold tracking-tight">{{ themeName }}</h1>
        <p class="text-muted-foreground text-base leading-relaxed">
          该主题需要验证码才能访问，请输入6位验证码。
        </p>
      </div>

      <form class="w-full flex flex-col items-center gap-4" @submit.prevent="handleSubmit">
        <div class="w-full space-y-1.5">
          <Label class="text-sm text-muted-foreground">验证码</Label>
          <Input
            v-model="code"
            placeholder="输入6位验证码"
            maxlength="6"
            class="rounded-xl text-center font-mono text-lg tracking-[0.3em] uppercase h-12"
            autofocus
            @input="error = ''"
          />
          <p v-if="error" class="text-xs text-destructive text-left">{{ error }}</p>
        </div>

        <div class="flex items-center gap-3">
          <Tooltip>
            <TooltipTrigger as-child>
              <Button type="button" variant="outline" size="icon" class="rounded-full h-11 w-11" @click="goHome">
                <Home class="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>返回首页</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger as-child>
              <Button type="button" variant="outline" size="icon" class="rounded-full h-11 w-11" @click="goGallery">
                <ArrowLeft class="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>主题画廊</TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger as-child>
              <Button type="submit" size="icon" class="rounded-full h-11 w-11" :disabled="verifying || code.length !== 6">
                <Loader2 v-if="verifying" class="size-5 animate-spin" />
                <Check v-else class="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>验证</TooltipContent>
          </Tooltip>
        </div>
      </form>
    </div>
  </div>
</template>
