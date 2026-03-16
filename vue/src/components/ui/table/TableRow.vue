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
  <tr :class="cn('border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted', props.class)" v-bind="attrs">
    <slot />
  </tr>
</template>
