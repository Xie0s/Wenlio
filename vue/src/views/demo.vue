<!--
  demo.vue - 组件演示页面
  职责：作为组件库的入口页面，展示所有基础 UI 组件的测试及示例
  对外接口：无
-->
<script setup lang="ts">
import ComponentGallery from '@/components/common/demo/ComponentGallery.vue'
</script>

<template>
  <div class="demo-page">
    <!-- 主体组件库展示 -->
    <ComponentGallery />
  </div>
</template>

<style scoped>
.demo-page {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}
</style>
