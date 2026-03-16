<!-- ContextMenu.vue - 右键菜单根组件
     职责：提供右键菜单的状态管理和上下文
     modal 默认设为 false：右键菜单是轻量交互，打开时不应锁定背景滚动、隐藏滚动条 -->
<script setup lang="ts">
import type { ContextMenuRootEmits, ContextMenuRootProps } from "reka-ui"
import { ContextMenuRoot, useForwardPropsEmits } from "reka-ui"

const props = withDefaults(defineProps<ContextMenuRootProps>(), {
  modal: false,
})
const emits = defineEmits<ContextMenuRootEmits>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <ContextMenuRoot
    data-slot="context-menu"
    v-bind="forwarded"
  >
    <slot />
  </ContextMenuRoot>
</template>
