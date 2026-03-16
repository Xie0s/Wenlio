<!-- PopoverContent.vue - Popover 弹出内容层
     职责：渲染弹出气泡的内容区域，通过 Portal 渲染到 body 以避免层叠上下文问题
     对外暴露：作为 Popover 的内容容器，接收 PopoverContentProps -->
<script setup lang="ts">
import type { PopoverContentEmits, PopoverContentProps } from "reka-ui"
import { computed, type HTMLAttributes, unref, useAttrs } from "vue"
import { reactiveOmit } from "@vueuse/core"
import {
  PopoverContent,
  PopoverPortal,
  useForwardPropsEmits,
} from "reka-ui"
import { cn } from "@/utils"

defineOptions({
  inheritAttrs: false,
})

const props = withDefaults(
  defineProps<PopoverContentProps & { class?: HTMLAttributes["class"] }>(),
  {
    align: "center",
    sideOffset: 4,
  },
)
const emits = defineEmits<PopoverContentEmits>()

const delegatedProps = reactiveOmit(props, "class")
const attrs = useAttrs()

const forwarded = useForwardPropsEmits(delegatedProps, emits)

const contentBindings = computed(() => ({
  ...attrs,
  ...unref(forwarded),
}))

const contentClass = computed(() =>
  cn(
    "glass bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=open]:duration-150 data-[state=closed]:duration-100 data-[state=open]:ease-out data-[state=closed]:ease-in z-50 w-72 origin-(--reka-popover-content-transform-origin) transform-gpu will-change-[transform,opacity] rounded-2xl border p-4 outline-hidden",
    props.class,
  ),
)
</script>

<template>
  <PopoverPortal>
    <PopoverContent
      data-slot="popover-content"
      v-bind="contentBindings"
      :class="contentClass"
    >
      <slot />
    </PopoverContent>
  </PopoverPortal>
</template>
