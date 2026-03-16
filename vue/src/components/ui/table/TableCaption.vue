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
  <caption :class="cn('mt-4 text-sm text-muted-foreground', props.class)" v-bind="attrs">
    <slot />
  </caption>
</template>
