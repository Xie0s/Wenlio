<script setup lang="ts">
import type { TabsRootEmits, TabsRootProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { ref, watch, provide } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { TabsRoot, useForwardPropsEmits } from "reka-ui"
import { cn } from "@/utils"
import { TABS_SLIDE_DIRECTION_KEY, TABS_REGISTERED_VALUES_KEY } from "./keys"

const props = defineProps<TabsRootProps & { class?: HTMLAttributes["class"] }>()
const emits = defineEmits<TabsRootEmits>()

const delegatedProps = reactiveOmit(props, "class")
const forwarded = useForwardPropsEmits(delegatedProps, emits)

// 存储 tab 顺序
const registeredValues: string[] = []
const slideDirection = ref<'left' | 'right'>('right')
const previousValue = ref<string | number | undefined>(props.modelValue || props.defaultValue)

// 提供给子组件的方法
provide(TABS_REGISTERED_VALUES_KEY, {
  register: (val: string) => {
    if (!registeredValues.includes(val)) {
      registeredValues.push(val)
    }
  },
  getDirection: () => slideDirection.value
})

provide(TABS_SLIDE_DIRECTION_KEY, slideDirection)

// 监听 modelValue 变化，计算滑动方向
watch(() => props.modelValue, (newValue, oldValue) => {
  if (newValue !== undefined && oldValue !== undefined) {
    const newIndex = registeredValues.indexOf(String(newValue))
    const oldIndex = registeredValues.indexOf(String(oldValue))
    
    if (newIndex !== -1 && oldIndex !== -1) {
      slideDirection.value = newIndex > oldIndex ? 'right' : 'left'
    }
  }
  previousValue.value = newValue
})
</script>

<template>
  <TabsRoot
    v-slot="slotProps"
    data-slot="tabs"
    v-bind="forwarded"
    :class="cn('flex flex-col gap-2', props.class)"
  >
    <slot v-bind="slotProps" />
  </TabsRoot>
</template>
