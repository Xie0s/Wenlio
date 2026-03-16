<!--
  TabsList 组件
  设计特点：
  - 圆润的胶囊形状容器（rounded-full）
  - 边框环绕样式（ring-1 ring-border）
  - 紧凑的内边距
-->
<script setup lang="ts">
import type { TabsListProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsList } from "reka-ui"
import { cn } from "@/utils"

const props = defineProps<TabsListProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <TabsList
    data-slot="tabs-list"
    v-bind="delegatedProps"
    :class="cn(
      'inline-flex h-10 items-stretch justify-center rounded-full text-muted-foreground ring-1 ring-border',
      props.class,
    )"
  >
    <slot />
  </TabsList>
</template>
