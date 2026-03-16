<!--
  DialogScrollContent.vue - 可滚动对话框内容组件
  职责：当对话框内容过长时，提供独立滚动的容器。
  功能说明：
  1. 升级圆角为 rounded-3xl (1.5rem)。
  2. 右上角关闭按钮更新为圆形轮廓线按钮样式 (outline + rounded-full)。
  3. 集成 Tooltip 功能，悬浮显示“关闭”提示。
  参数：
  - showCloseButton: boolean - 是否显示右上角关闭按钮，默认开启。
-->
<script setup lang="ts">
import type { DialogContentEmits, DialogContentProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { X } from "lucide-vue-next"
import {
  DialogClose,
  DialogContent,
  DialogOverlay,
  DialogPortal,
  useForwardPropsEmits,
} from "reka-ui"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/utils"

defineOptions({
  inheritAttrs: false,
})

const props = withDefaults(defineProps<DialogContentProps & { class?: HTMLAttributes["class"], showCloseButton?: boolean }>(), {
  showCloseButton: true,
})
const emits = defineEmits<DialogContentEmits>()

const delegatedProps = reactiveOmit(props, "class")

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <DialogPortal>
    <DialogOverlay
      class="fixed inset-0 z-50 bg-black/80 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
    />
    <DialogContent
      :class="
        cn(
          'bg-background data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 border border-border rounded-3xl p-6 duration-200 sm:max-w-lg max-h-[calc(100vh-4rem)] overflow-y-auto',
          props.class,
        )
      "
      v-bind="{ ...$attrs, ...forwarded }"
      @pointer-down-outside="(event) => {
        const originalEvent = event.detail.originalEvent;
        const target = originalEvent.target as HTMLElement;
        if (originalEvent.offsetX > target.clientWidth || originalEvent.offsetY > target.clientHeight) {
          event.preventDefault();
        }
      }"
    >
      <slot />

      <Tooltip v-if="showCloseButton">
        <TooltipTrigger as-child>
          <DialogClose
            data-slot="dialog-close"
            :class="cn(
              buttonVariants({ variant: 'outline', size: 'icon' }),
              'absolute top-4 right-4 rounded-full',
            )"
          >
            <X class="size-4" />
            <span class="sr-only">Close</span>
          </DialogClose>
        </TooltipTrigger>
        <TooltipContent>关闭</TooltipContent>
      </Tooltip>
    </DialogContent>
  </DialogPortal>
</template>
