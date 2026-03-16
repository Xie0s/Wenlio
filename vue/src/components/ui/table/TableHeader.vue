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
  <thead :class="cn('[&_tr]:border-b', props.class)" v-bind="attrs">
    <slot />
  </thead>
</template>
