<!-- MaintenancePage.vue - 站点维护模式页面
     职责：当租户开启维护模式且用户未登录时展示维护提示，引导登录
     对外接口：路由参数 tenantId -->
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReaderStore } from '@/stores/reader'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Construction, LogIn } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const store = useReaderStore()

const tenantId = computed(() => route.params.tenantId as string)
const tenantName = computed(() => store.tenant?.name || tenantId.value)

function goLogin() {
  router.push('/admin/login')
}
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-background px-4">
    <div class="flex flex-col items-center gap-6 max-w-md text-center">
      <div class="rounded-full bg-amber-100 dark:bg-amber-900/30 p-5">
        <Construction class="size-10 text-amber-600 dark:text-amber-400" :stroke-width="1.5" />
      </div>

      <div class="space-y-2">
        <h1 class="text-2xl font-semibold tracking-tight">{{ tenantName }}</h1>
        <p class="text-muted-foreground text-base leading-relaxed">
          站点正在维护中，暂时无法访问。<br />
          如您是管理员，请登录后继续。
        </p>
      </div>

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
</template>
