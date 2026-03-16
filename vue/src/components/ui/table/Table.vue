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
  <div class="relative w-full overflow-auto rounded-xl">
    <table :class="cn('w-full caption-bottom text-sm', props.class)" v-bind="attrs">
      <slot />
    </table>
  </div>
</template>
