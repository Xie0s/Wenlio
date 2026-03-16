<!--
  DynamicLucideIcon.vue - 根据字符串名称动态渲染 Lucide 图标
  支持 kebab-case（如 'book-open'）和 PascalCase（如 'BookOpen'）两种格式
  fallback 为 BookOpen
-->
<script setup lang="ts">
import { computed } from 'vue'
import * as LucideIcons from 'lucide-vue-next'
import { BookOpen } from 'lucide-vue-next'

const props = defineProps<{
  name: string
}>()

function toPascalCase(str: string): string {
  return str
    .split(/[-_\s]/)
    .filter(Boolean)
    .map(s => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase())
    .join('')
}

const icon = computed(() => {
  if (!props.name) return BookOpen
  const pascal = toPascalCase(props.name)
  return (LucideIcons as Record<string, any>)[pascal] ?? BookOpen
})
</script>

<template>
  <component :is="icon" v-bind="$attrs" />
</template>
