<!--
  AlertDialogCancel.vue - 警示对话框取消操作组件
  职责：关闭模态框且不执行任何操作。
  功能说明：
  1. 样式调整为圆形轮廓线按钮 (outline + rounded-full)。
  2. 默认显示 X 图标，隐藏文字。
-->
<script setup lang="ts">
import type { AlertDialogCancelProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { X } from "lucide-vue-next"
import { AlertDialogCancel } from "reka-ui"
import { cn } from "@/utils"
import { buttonVariants } from '@/components/ui/button'

const props = defineProps<AlertDialogCancelProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <AlertDialogCancel
    v-bind="delegatedProps"
    :class="cn(
      buttonVariants({ variant: 'outline', size: 'icon' }),
      'rounded-full',
      props.class,
    )"
  >
    <X class="size-4" />
    <span class="sr-only"><slot>取消</slot></span>
  </AlertDialogCancel>
</template>
