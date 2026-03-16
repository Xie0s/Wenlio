<!--
  auth/UserMenu.vue - 统一用户头像组件，支持两种使用模式
  职责：展示当前登录用户头像，根据 mode 切换交互方式
  对外暴露：
    Props:
      - mode?:        'dropdown' | 'hover-card'  交互模式（默认 'dropdown'）
      - tenantId?:    string   当前租户 ID，提供时下拉中显示"返回首页"跳转项
      - side?:        string   dropdown 弹出方向（默认 'bottom'）
      - align?:       string   dropdown 对齐方式（默认 'end'）
      - sideOffset?:  number   dropdown 偏移量（默认 8）
      - triggerClass?: string  触发按钮额外 class（默认 'h-10 w-10'）
      - avatarClass?:  string  Avatar 额外 class（默认 'h-9 w-9'）
-->
<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import UserHoverCard from '@/components/common/UserHoverCard.vue'
import { User, LogOut, ArrowUpRight, LayoutDashboard } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const props = withDefaults(defineProps<{
  mode?: 'dropdown' | 'hover-card'
  tenantId?: string
  side?: 'top' | 'bottom' | 'left' | 'right'
  align?: 'start' | 'center' | 'end'
  sideOffset?: number
  triggerClass?: string
  avatarClass?: string
}>(), {
  mode: 'dropdown',
  side: 'bottom',
  align: 'end',
  sideOffset: 8,
  triggerClass: 'h-10 w-10',
  avatarClass: 'h-9 w-9',
})

const router = useRouter()
const authStore = useAuthStore()

const displayName = computed(() => authStore.user?.name || authStore.user?.username || '用户')
const roleLabel = computed(() => {
  if (authStore.isSuperAdmin) return '超级管理员'
  if (authStore.isTenantAdmin) return '租户管理员'
  return '用户'
})
const avatarFallback = computed(() => displayName.value.charAt(0).toUpperCase())
const isAdmin = computed(() => authStore.isSuperAdmin || authStore.isTenantAdmin)

const tooltipSide = computed(() => props.side === 'top' ? 'top' : 'bottom')

function goToProfile() {
  router.push('/admin/profile')
}

function goToAdmin() {
  router.push('/admin')
}

function goToTenantHome() {
  if (props.tenantId) {
    router.push({ name: 'TenantHome', params: { tenantId: props.tenantId } })
    return
  }
  const tid = authStore.user?.tenant_id?.trim()
  if (tid) {
    router.push({ name: 'TenantHome', params: { tenantId: tid } })
    return
  }
  router.push('/')
}

async function handleLogout() {
  await authStore.logout()
  toast.success('已登出')
  router.push('/admin/login')
}
</script>

<template>
  <!-- ═══ Mode: dropdown（点击下拉菜单） ═══ -->
  <DropdownMenu v-if="props.mode === 'dropdown'">
    <DropdownMenuTrigger as-child>
      <button
        class="rounded-full flex items-center justify-center p-0.5 transition-colors focus:outline-none data-[state=open]:ring-2 data-[state=open]:ring-primary/70 data-[state=open]:ring-offset-0"
        :class="props.triggerClass" aria-label="用户中心">
        <Tooltip>
          <TooltipTrigger as-child>
            <Avatar :class="props.avatarClass">
              <AvatarImage v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" />
              <AvatarFallback class="text-sm font-medium bg-primary/15 text-primary">
                {{ avatarFallback }}
              </AvatarFallback>
            </Avatar>
          </TooltipTrigger>
          <TooltipContent :side="tooltipSide">用户中心</TooltipContent>
        </Tooltip>
      </button>
    </DropdownMenuTrigger>

    <DropdownMenuContent :side="props.side" :align="props.align" :side-offset="props.sideOffset"
      class="w-52 rounded-3xl z-[90]">
      <DropdownMenuLabel class="flex flex-col gap-0.5">
        <span class="font-medium">{{ displayName }}</span>
        <span class="text-xs font-normal text-muted-foreground">{{ roleLabel }}</span>
      </DropdownMenuLabel>
      <DropdownMenuSeparator />

      <DropdownMenuItem v-if="isAdmin" class="cursor-pointer" @click="goToProfile">
        <User class="size-4 mr-2" :stroke-width="1.5" />
        个人中心
      </DropdownMenuItem>
      <DropdownMenuItem v-if="isAdmin" class="cursor-pointer" @click="goToAdmin">
        <LayoutDashboard class="size-4 mr-2" :stroke-width="1.5" />
        管理后台
      </DropdownMenuItem>
      <DropdownMenuItem v-if="props.tenantId" class="cursor-pointer" @click="goToTenantHome">
        <ArrowUpRight class="size-4 mr-2" :stroke-width="1.5" />
        返回首页
      </DropdownMenuItem>

      <DropdownMenuSeparator />
      <DropdownMenuItem class="text-destructive focus:text-destructive cursor-pointer" @click="handleLogout">
        <LogOut class="size-4 mr-2" :stroke-width="1.5" />
        登出
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <!-- ═══ Mode: hover-card（悬停显示用户信息卡） ═══ -->
  <UserHoverCard v-else :avatar-url="authStore.user?.avatar_url" :name="displayName"
    :username="authStore.user?.username || ''" :bio="authStore.user?.bio"
    :profile-bg-url="authStore.user?.profile_bg_url" :role-label="roleLabel" :tenant-name="authStore.user?.tenant_name">
    <Avatar :class="[props.avatarClass, 'cursor-pointer']" @click="goToProfile">
      <AvatarImage v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" />
      <AvatarFallback class="text-sm font-medium bg-primary/15 text-primary">
        {{ avatarFallback }}
      </AvatarFallback>
    </Avatar>
  </UserHoverCard>
</template>
