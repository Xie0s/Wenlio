<!--
  AlertDialogAction.vue - 警示对话框确认操作组件
  职责：执行模态框的确认/主操作。
  功能说明：
   1. 样式调整为圆形主题色按钮 (primary + rounded-full)。
  2. 默认显示 Check (对勾) 图标，隐藏文字。
-->
<script setup lang="ts">
import type { AlertDialogActionProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { Check } from "lucide-vue-next"
import { AlertDialogAction } from "reka-ui"
import { cn } from "@/utils"
import { buttonVariants } from '@/components/ui/button'

const props = defineProps<AlertDialogActionProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <AlertDialogAction
    v-bind="delegatedProps"
    :class="cn(
      buttonVariants({ variant: 'default', size: 'icon' }),
      'rounded-full',
      props.class,
    )"
  >
    <Check class="size-4" />
    <span class="sr-only"><slot>确认</slot></span>
  </AlertDialogAction>
</template>
