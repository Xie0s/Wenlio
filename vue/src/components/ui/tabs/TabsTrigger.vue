<!--
  TabsTrigger 组件
  设计特点：
  - 选中时有 muted 背景和阴影
  - 平滑的过渡动画
  - hover 时有淡色背景
-->
<script setup lang="ts">
import type { TabsTriggerProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { onMounted, inject } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsTrigger, useForwardProps } from "reka-ui"
import { cn } from "@/utils"
import { TABS_REGISTERED_VALUES_KEY } from "./keys"

const props = defineProps<TabsTriggerProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
const forwardedProps = useForwardProps(delegatedProps)

// 注册 value 到父级 Tabs
const tabsContext = inject(TABS_REGISTERED_VALUES_KEY, null)
onMounted(() => {
  if (tabsContext && props.value) {
    tabsContext.register(String(props.value))
  }
})
</script>

<template>
  <TabsTrigger
    data-slot="tabs-trigger"
    :class="cn(
      'inline-flex h-full items-center justify-center whitespace-nowrap rounded-full px-4 text-sm font-medium leading-none',
      'ring-offset-background transition-all',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      'disabled:pointer-events-none disabled:opacity-50',
      'hover:bg-muted/40 hover:text-foreground',
      'data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow',
      '[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*=size-])]:size-4',
      props.class,
    )"
    v-bind="forwardedProps"
  >
    <slot />
  </TabsTrigger>
</template>
