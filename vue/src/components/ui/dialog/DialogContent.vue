<!--
  DialogContent.vue - 对话框内容组件
  职责：展示对话框的核心内容区域。
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
  DialogPortal,
  useForwardPropsEmits,
} from "reka-ui"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { buttonVariants } from "@/components/ui/button"
import { cn } from "@/utils"
import DialogOverlay from "./DialogOverlay.vue"

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
    <DialogOverlay />
    <DialogContent
      data-slot="dialog-content"
      v-bind="{ ...$attrs, ...forwarded }"
      :class="
        cn(
          'bg-background data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-3xl border p-6 duration-200 sm:max-w-lg',
          props.class,
        )"
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
