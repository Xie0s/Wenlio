<!--
  src/components/home/HomeFooter.vue
  职责：首页页脚，居中布局，参考 PFooter 设计风格
  对外暴露：无
-->
<script setup lang="ts">
import { RouterLink } from 'vue-router'

const navLinks = [
  { label: '核心能力', href: '#features' },
  { label: '组织结构', href: '#structure' },
  { label: '发布流程', href: '#workflow' },
  { label: '功能总览', href: '#more-features' },
  { label: '个性化组件', href: '#personalization' },
]

const metaLinks = [
  { label: '关于我们', href: 'https://microswift.cn', external: true },
  { label: '后台登录', href: '/admin/login', external: false },
  { label: '服务条款', href: '#', external: false },
  { label: '隐私政策', href: '#', external: false },
]

const LOGO_URL = 'https://ziyuan-1302808601.cos.ap-chengdu.myqcloud.com/microswift-logo/%E5%BE%AE%E8%AE%AF%E4%BA%91%E6%96%B0LOGO.gif'
</script>

<template>
  <footer class="relative w-full border-t border-black/[0.06] dark:border-white/[0.06] mt-20 text-muted-foreground">
    <div class="mx-auto max-w-6xl px-6 pt-10 pb-6">

      <!-- 主体：居中 -->
      <div class="flex flex-col items-center space-y-6 text-center">

        <!-- Logo -->
        <div class="flex justify-center">
          <img :src="LOGO_URL" alt="微讯云" class="h-auto w-auto object-contain dark:brightness-0 dark:invert"
            style="max-height: 64px;" />
        </div>

        <!-- 标语 -->
        <p class="text-base font-light text-muted-foreground max-w-lg leading-relaxed">
          支持多租户隔离、版本化管理、即时发布与评论审核的一体化文档平台。
        </p>

        <!-- 导航链接 -->
        <div class="flex flex-wrap justify-center items-center gap-x-6 gap-y-2">
          <a v-for="link in navLinks" :key="link.label" :href="link.href"
            class="text-sm text-muted-foreground transition-colors hover:text-foreground">
            {{ link.label }}
          </a>
        </div>

      </div>

      <!-- 分割线 -->
      <div class="mt-8 mb-4 h-px bg-black/[0.06] dark:bg-white/[0.06]" />

      <!-- 底部版权行 -->
      <div class="flex flex-col items-center justify-between gap-2 text-xs sm:flex-row sm:items-center">
        <!-- 版权 -->
        <span class="shrink-0">Microswift Core™ 版权所有 ©2020-2026</span>

        <!-- 微讯云信息 -->
        <div class="flex flex-wrap justify-center gap-4">
          <template v-for="link in metaLinks" :key="link.label">
            <a v-if="link.external" :href="link.href" target="_blank" rel="noopener noreferrer"
              class="transition-colors hover:text-foreground">{{ link.label }}</a>
            <RouterLink v-else-if="link.href.startsWith('/')" :to="link.href"
              class="transition-colors hover:text-foreground">{{ link.label }}</RouterLink>
            <a v-else :href="link.href" class="transition-colors hover:text-foreground">{{ link.label }}</a>
          </template>
        </div>

        <!-- 备案信息 -->
        <div class="flex items-center gap-2 shrink-0">
          <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer"
            class="transition-colors hover:text-foreground">陇ICP备20002844号-1</a>
          <span class="text-muted-foreground/30">|</span>
          <a href="https://www.beian.gov.cn/" target="_blank" rel="noopener noreferrer"
            class="transition-colors hover:text-foreground">甘公网安备62090202000540号</a>
        </div>
      </div>

    </div>
  </footer>
</template>
