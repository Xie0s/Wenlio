<script lang="ts" setup>
/**
 * DrawerDirectionDemo.vue
 * 多方向抽屉演示组件
 * 展示 Drawer 在不同方向（Top, Left, Right）的布局效果
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

const directions = [
  { value: 'top', label: '上方', icon: Icons.PanelTop },
  { value: 'left', label: '左侧', icon: Icons.PanelLeft },
  { value: 'right', label: '右侧', icon: Icons.PanelRight },
] as const

function handleTriggerClick(event: MouseEvent) {
  (event.currentTarget as HTMLButtonElement | null)?.blur()
}
</script>

<template>
  <div class="flex flex-wrap gap-4">
    <TooltipProvider>
      <template v-for="dir in directions" :key="dir.value">
        <Drawer :direction="dir.value">
          <Tooltip>
            <TooltipTrigger as-child>
              <DrawerTrigger as-child>
                <Button variant="outline" size="icon" class="rounded-full" @click="handleTriggerClick">
                  <component :is="dir.icon" class="size-4" />
                </Button>
              </DrawerTrigger>
            </TooltipTrigger>
            <TooltipContent>从{{ dir.label }}打开</TooltipContent>
          </Tooltip>

          <DrawerContent>
            <div :class="[
              'mx-auto w-full h-full flex flex-col',
              (dir.value === 'left' || dir.value === 'right') ? 'max-w-sm' : 'max-w-md'
            ]">
              <DrawerHeader>
                <DrawerTitle>{{ dir.label }}抽屉</DrawerTitle>
                <DrawerDescription>这是从{{ dir.label }}方向弹出的抽屉示例。</DrawerDescription>
              </DrawerHeader>
              <div class="flex-1 p-4">
                <div class="rounded-2xl border-2 border-dashed border-muted p-8 flex items-center justify-center text-muted-foreground italic">
                  内容区域
                </div>
              </div>
              <DrawerFooter>
                <DrawerClose as-child>
                  <Button variant="outline" class="w-full rounded-full">关闭</Button>
                </DrawerClose>
              </DrawerFooter>
            </div>
          </DrawerContent>
        </Drawer>
      </template>
    </TooltipProvider>
  </div>
</template>
