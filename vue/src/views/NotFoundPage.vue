<!-- NotFoundPage.vue - 404 页面
     职责：路由未匹配或资源不存在时的友好错误提示页
     对外暴露：无（路由组件） -->
<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { FileSearch, ArrowLeft, Home } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-background p-4 text-center">
    <div class="space-y-8 max-w-md">
      <div class="flex items-center justify-center">
        <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
          <FileSearch class="h-8 w-8" />
        </div>
      </div>

      <div class="space-y-2">
        <h1 class="text-6xl font-medium tracking-tight text-foreground">404</h1>
        <p class="text-lg font-medium">页面不存在</p>
        <p class="text-sm text-muted-foreground">你访问的页面不存在或已被移除</p>
      </div>

      <div class="flex items-center justify-center gap-3">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="outline" size="icon" class="rounded-full" @click="router.back()">
              <ArrowLeft class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>返回上一页</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button size="icon" class="rounded-full" @click="router.push('/')">
              <Home class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>回到首页</TooltipContent>
        </Tooltip>
      </div>
    </div>
  </div>
</template>
