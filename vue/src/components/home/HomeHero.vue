<!--
  src/components/home/HomeHero.vue
  职责：首页 Hero 区域，包含主标题、副标题、CTA 按钮和交互式演示窗口
  对外暴露：无
-->
<script setup lang="ts">
import HeroDemo from './hero-demo/HeroDemo.vue'
import TenantEntryPill from './welcome/TenantEntryPill.vue'
</script>

<template>
  <section class="pt-32 pb-0">

    <div class="relative mx-auto max-w-[1360px] px-6">
      <!-- 标题区 -->
      <div class="mx-auto max-w-3xl text-center">
        <div class="mb-5 inline-flex items-center rounded-full border border-black/10 bg-white/60 px-3 py-1 text-xs text-muted-foreground backdrop-blur-sm dark:border-white/20 dark:bg-white/10 dark:text-white/80">
          文档管理平台 · v2
        </div>

        <h1 class="animate-fade-up mb-5 text-5xl font-medium leading-[1.12] tracking-tight text-foreground sm:text-6xl">
          多租户文档管理与阅读<br />一个应用即刻上线
        </h1>

        <p class="animate-fade-up mb-8 text-lg text-muted-foreground" style="animation-delay: 0.1s">
          从主题与版本管理、Markdown 编辑与渲染，到全文搜索与评论审核，形成完整闭环。<br class="hidden sm:block" />发布仅变更状态，无需构建等待，读者刷新即见最新内容。
        </p>

        <div class="animate-fade-up flex flex-col items-center gap-4" style="animation-delay: 0.2s">
          <!-- 租户入口胶囊 -->
          <TenantEntryPill />

          <div class="flex items-center gap-3">
            <a
              href="#features"
              class="inline-flex h-9 items-center rounded-full border border-black/10 bg-white/70 px-4 text-sm text-foreground backdrop-blur-sm dark:border-white/20 dark:bg-white/10 dark:text-white transition-colors hover:bg-white dark:hover:bg-white/20"
            >
              查看核心能力与流程
            </a>
          </div>
        </div>
      </div>

      <!-- 交互式演示区 -->
      <div class="animate-fade-in mt-10 px-0 sm:px-2" style="animation-delay: 0.35s">
        <HeroDemo />
      </div>
    </div>

  </section>
</template>
