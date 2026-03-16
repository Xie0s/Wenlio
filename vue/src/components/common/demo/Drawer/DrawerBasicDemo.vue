<script lang="ts" setup>
/**
 * DrawerBasicDemo.vue
 * 基础底部抽屉演示组件
 * 展示 Drawer 的标准用法，包含 Header, Footer 和基础内容
 */
import { ref } from 'vue'
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
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import * as Icons from 'lucide-vue-next'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

const isOpen = ref(false)

function handleTriggerClick(event: MouseEvent) {
  (event.currentTarget as HTMLButtonElement | null)?.blur()
}
</script>

<template>
  <div class="flex items-center gap-4">
    <Drawer v-model:open="isOpen">
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger as-child>
            <DrawerTrigger as-child>
              <Button variant="outline" size="icon" class="rounded-full" @click="handleTriggerClick">
                <Icons.PanelBottom class="size-4" />
              </Button>
            </DrawerTrigger>
          </TooltipTrigger>
          <TooltipContent>打开基础抽屉</TooltipContent>
        </Tooltip>
      </TooltipProvider>

      <DrawerContent>
        <div class="mx-auto w-full max-w-sm">
          <DrawerHeader>
            <DrawerTitle>设置选项</DrawerTitle>
            <DrawerDescription>在此调整您的偏好设置。</DrawerDescription>
          </DrawerHeader>
          <div class="p-4 pb-0 space-y-4">
            <div class="space-y-2">
              <label class="text-sm font-medium leading-none">项目名称</label>
              <Input placeholder="输入项目名称..." />
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium leading-none">描述</label>
              <Textarea placeholder="补充说明信息..." />
            </div>
          </div>
          <DrawerFooter>
            <Button class="w-full rounded-full">提交更改</Button>
            <DrawerClose as-child>
              <Button variant="outline" class="w-full rounded-full">取消</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  </div>
</template>
