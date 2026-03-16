<!--
  FileNodeView.vue - 文件附件节点的 Tiptap NodeView 渲染组件
  职责：包裹 FileCard 基础组件，根据 layout 属性控制 NodeViewWrapper 的 CSS 布局类
  对外暴露：Tiptap NodeViewRenderer 组件（由 FileNode 扩展注册）
-->
<script setup lang="ts">
import { NodeViewWrapper } from '@tiptap/vue-3'
import type { NodeViewProps } from '@tiptap/vue-3'
import FileCard from '../common/FileCard.vue'

const props = defineProps<NodeViewProps>()
</script>

<template>
  <NodeViewWrapper
    :class="['file-node-wrapper', `file-layout-${node.attrs.layout || 'full'}`]"
    :contenteditable="false"
  >
    <FileCard
      :file-name="node.attrs.fileName || '附件'"
      :file-size="node.attrs.fileSize"
      :mime-type="node.attrs.mimeType"
      :selected="selected"
    />
  </NodeViewWrapper>
</template>
