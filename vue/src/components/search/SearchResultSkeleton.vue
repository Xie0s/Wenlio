<!--
  SearchResultSkeleton.vue - 搜索结果骨架屏组件
  职责：搜索加载时渐进式显示骨架占位，模拟 SearchResultItem 布局
  对外暴露：Props: index (0-based, 用于错开宽度增加视觉层次感)
-->
<template>
  <div class="w-full flex items-center gap-3 px-4 py-3 animate-pulse">
    <!-- 图标占位 -->
    <div class="shrink-0 w-9 h-9 rounded-lg bg-muted/70" />

    <!-- 文字信息占位 -->
    <div class="flex-1 min-w-0 space-y-2">
      <div class="h-4 rounded bg-muted/70" :style="{ width: titleWidth }" />
      <div class="h-3 rounded bg-muted/50" :style="{ width: snippetWidth }" />
    </div>

    <!-- 右侧元信息占位 -->
    <div class="shrink-0 space-y-1.5">
      <div class="h-3 w-14 rounded bg-muted/70 ml-auto" />
      <div class="h-3 w-10 rounded bg-muted/50 ml-auto" />
    </div>

    <!-- 箭头占位 -->
    <div class="shrink-0 w-4 h-4 rounded bg-muted/50" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  index?: number
}>(), {
  index: 0,
})

const titleWidth = computed(() => {
  const widths = ['55%', '70%', '45%', '80%', '65%', '50%', '75%', '60%']
  return widths[props.index % widths.length]
})

const snippetWidth = computed(() => {
  const widths = ['85%', '70%', '90%', '75%', '80%', '65%', '95%', '72%']
  return widths[props.index % widths.length]
})
</script>
