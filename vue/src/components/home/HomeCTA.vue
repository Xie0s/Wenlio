<!--
  src/components/home/HomeCTA.vue
  职责：首页底部行动召唤区域（Call to Action）
  对外暴露：无
-->
<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'
</script>

<template>
  <section class="pt-8 pb-2">
    <div class="mx-auto max-w-6xl px-6 text-center">
      <p class="mb-4 text-sm text-muted-foreground">第七部分 · 立即开始</p>
      <h2 class="mb-4 text-4xl font-medium tracking-tight md:text-5xl">
        按需求链路完成导览，立即开始管理文档。
      </h2>
      <p class="mx-auto mb-10 max-w-xl text-base text-muted-foreground">
        从主题与版本管理，到文档发布、搜索与评论审核，后台能力已完整就绪。登录后即可按租户开始配置与发布。
      </p>
      <div class="flex flex-col items-center justify-center gap-3 sm:flex-row">
        <router-link
          to="/admin/login"
          class="inline-flex h-11 items-center gap-2 rounded-full bg-foreground px-6 text-sm font-medium text-background transition-opacity hover:opacity-90"
        >
          进入后台开始使用
          <ArrowRight class="h-4 w-4" />
        </router-link>
        <a
          href="#personalization"
          class="inline-flex h-11 items-center rounded-full border border-black/10 dark:border-white/10 px-6 text-sm font-medium text-foreground transition-colors hover:bg-muted/40"
        >
          回看第六部分：个性化组件说明
        </a>
      </div>

      <!-- 底部小字 -->
      <p class="mt-8 text-xs text-muted-foreground/60">
        支持超级管理员、租户管理员与读者三类角色协作
      </p>
    </div>
  </section>
</template>
