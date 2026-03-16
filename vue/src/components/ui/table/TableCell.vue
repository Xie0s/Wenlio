<script setup lang="ts">
import { useAttrs } from 'vue'
import { cn } from '@/utils'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  class?: string
}>()

const attrs = useAttrs()
</script>

<template>
  <td :class="cn('p-4 align-middle [&:has([role=checkbox])]:pr-0', props.class)" v-bind="attrs">
    <slot />
  </td>
</template>
