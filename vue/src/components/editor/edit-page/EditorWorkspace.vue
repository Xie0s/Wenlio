<!--
  EditorWorkspace.vue - 源码模式工作区
  职责：渲染 Markdown 编辑内容区（编辑 + 预览），TOC 与章节导航已移至 EditorSidebar
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { useEditorInject } from '@/composables/useThemeEditor'
import EditorContent from './EditorContent.vue'
import EditorEmptyGuide from './EditorEmptyGuide.vue'

const editor = useEditorInject()
</script>

<template>
  <EditorContent v-if="editor.page.value" />
  <EditorEmptyGuide v-else />
</template>
