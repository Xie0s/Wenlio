<!-- DashboardPage.vue - 管理后台仪表盘
     职责：展示当前用户的角色信息和快捷入口 -->
<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Building2, BookOpen, MessageSquare, Users, ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const isSuperAdmin = computed(() => authStore.user?.role === 'super_admin')

const quickActions = computed(() => {
  if (isSuperAdmin.value) {
    return [
      { icon: Building2, label: '租户管理', desc: '创建和管理租户', to: '/admin/tenants' },
      { icon: Users, label: '用户管理', desc: '管理超级管理员账号', to: '/admin/users' },
    ]
  }
  return [
    { icon: BookOpen, label: '主题管理', desc: '管理文档主题', to: '/admin/themes' },
    { icon: MessageSquare, label: '评论管理', desc: '审核读者评论', to: '/admin/comments' },
    { icon: Users, label: '用户管理', desc: '管理租户管理员', to: '/admin/tenant-users' },
  ]
})
</script>

<template>
  <div class="space-y-6">
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="action in quickActions" :key="action.to" class="rounded-3xl cursor-pointer hover:border-primary/30 transition-colors" @click="router.push(action.to)">
        <CardHeader class="flex flex-row items-center gap-3 pb-2">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <component :is="action.icon" class="h-5 w-5" />
          </div>
          <CardTitle class="text-base font-medium">{{ action.label }}</CardTitle>
        </CardHeader>
        <CardContent class="flex items-center justify-between">
          <p class="text-sm text-muted-foreground">{{ action.desc }}</p>
          <Tooltip>
            <TooltipTrigger as-child>
              <Button variant="ghost" size="icon" class="h-8 w-8 rounded-full">
                <ArrowRight class="h-4 w-4" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>前往</TooltipContent>
          </Tooltip>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
