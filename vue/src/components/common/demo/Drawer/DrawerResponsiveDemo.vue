<script lang="ts" setup>
/**
 * DrawerResponsiveDemo.vue
 * 响应式/复杂内容抽屉演示组件
 * 展示包含列表、切换开关等复杂交互的抽屉
 */
import { Button } from '@/components/ui/button'
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer'
import * as Icons from 'lucide-vue-next'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

const items = [
  { id: 1, label: '个人资料', icon: Icons.User },
  { id: 2, label: '安全设置', icon: Icons.ShieldCheck },
  { id: 3, label: '通知偏好', icon: Icons.Bell },
  { id: 4, label: '账户绑定', icon: Icons.Link },
  { id: 5, label: '数据隐私', icon: Icons.Lock },
]

function handleTriggerClick(event: MouseEvent) {
  (event.currentTarget as HTMLButtonElement | null)?.blur()
}
</script>

<template>
  <Drawer>
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger as-child>
          <DrawerTrigger as-child>
            <Button variant="outline" size="icon" class="rounded-full" @click="handleTriggerClick">
              <Icons.Settings2 class="size-4" />
            </Button>
          </DrawerTrigger>
        </TooltipTrigger>
        <TooltipContent>打开复杂抽屉</TooltipContent>
      </Tooltip>
    </TooltipProvider>

    <DrawerContent>
      <div class="mx-auto w-full max-w-sm">
        <DrawerHeader>
          <DrawerTitle>系统菜单</DrawerTitle>
          <DrawerDescription>快速访问您的账户与系统配置。</DrawerDescription>
        </DrawerHeader>
        <div class="p-2 space-y-1">
          <button
            v-for="item in items"
            :key="item.id"
            class="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-accent transition-colors text-sm font-medium"
          >
            <component :is="item.icon" class="size-4 text-muted-foreground" />
            {{ item.label }}
            <Icons.ChevronRight class="ml-auto size-4 text-muted-foreground/50" />
          </button>
        </div>
        <DrawerFooter class="pt-2">
          <Button variant="destructive" class="w-full rounded-full">退出登录</Button>
          <DrawerClose as-child>
            <Button variant="ghost" class="w-full rounded-full">取消</Button>
          </DrawerClose>
        </DrawerFooter>
      </div>
    </DrawerContent>
  </Drawer>
</template>
