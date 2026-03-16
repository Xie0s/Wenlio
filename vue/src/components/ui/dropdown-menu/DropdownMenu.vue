<!-- DropdownMenu.vue - 下拉菜单根组件
     职责：提供下拉菜单的状态管理和上下文
     modal 默认设为 false：下拉菜单是轻量交互，打开时不应锁定背景滚动、隐藏滚动条 -->
<script setup lang="ts">
import type { DropdownMenuRootEmits, DropdownMenuRootProps } from "reka-ui"
import { DropdownMenuRoot, useForwardPropsEmits } from "reka-ui"

const props = withDefaults(defineProps<DropdownMenuRootProps>(), {
  modal: false,
})
const emits = defineEmits<DropdownMenuRootEmits>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <DropdownMenuRoot
    v-slot="slotProps"
    data-slot="dropdown-menu"
    v-bind="forwarded"
  >
    <slot v-bind="slotProps" />
  </DropdownMenuRoot>
</template>
