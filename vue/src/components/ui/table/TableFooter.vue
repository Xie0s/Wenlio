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
  <tfoot :class="cn('border-t bg-muted/50 font-medium [&>tr]:last:border-b-0', props.class)" v-bind="attrs">
    <slot />
  </tfoot>
</template>
