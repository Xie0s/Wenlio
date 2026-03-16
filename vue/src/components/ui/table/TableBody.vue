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
  <tbody :class="cn('[&_tr:last-child]:border-0', props.class)" v-bind="attrs">
    <slot />
  </tbody>
</template>
