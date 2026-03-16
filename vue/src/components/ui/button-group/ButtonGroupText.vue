<script setup lang="ts">
import type { PrimitiveProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import type { ButtonGroupVariants } from "."
import { Primitive } from "reka-ui"
import { cn } from "@/utils"

interface Props extends PrimitiveProps {
  class?: HTMLAttributes["class"]
  orientation?: ButtonGroupVariants["orientation"]
}

const props = withDefaults(defineProps<Props>(), {
  as: "div",
})
</script>

<template>
  <Primitive
    role="group"
    data-slot="button-group-text"
    :data-orientation="props.orientation"
    :as="props.as"
    :as-child="props.asChild"
    :class="cn('bg-muted flex w-fit items-center gap-2 rounded-2xl border px-4 py-2 text-sm font-medium [&_svg]:pointer-events-none [&_svg:not([class*=\'size-\'])]:size-4', props.class)"
  >
    <slot />
  </Primitive>
</template>
