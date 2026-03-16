<!--
  src/components/home/HomePersonalization.vue
  职责：首页个性化能力说明区，介绍租户首页个性化组件体系
  设计：参考 HomeMoreFeatures.vue 的简洁卡片风格，卡片内嵌入编辑器设计预览
  对外暴露：无
-->
<script setup lang="ts">
import {
  Navigation,
  Layout,
  Grid3X3,
  List,
  Megaphone,
  PanelBottom,
  Palette,
  ArrowRight,
  Check,
  Plus,
  Trash2,
  Eye,
  GitBranch,
} from 'lucide-vue-next'

// 核心特性小卡片
const coreFeatures = [
  {
    icon: Palette,
    title: '全局样式',
    desc: '背景色、字体、间距统一管控',
  },
  {
    icon: Eye,
    title: '实时预览',
    desc: '配置变更即时生效',
  },
  {
    icon: GitBranch,
    title: '草稿发布',
    desc: '草稿与已发布配置独立存储',
  },
]

// 布局选项预览 SVG
const layoutPreviews = {
  centered: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="15" y="10" width="30" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/><rect x="20" y="16" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.3"/><rect x="22" y="21" width="16" height="7" rx="3.5" fill="currentColor" fill-opacity="0.25"/></svg>`,
  split: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="4" y="8" width="22" height="3.5" rx="1.75" fill="currentColor" fill-opacity="0.5"/><rect x="4" y="14" width="18" height="2" rx="1" fill="currentColor" fill-opacity="0.3"/><rect x="4" y="24" width="13" height="6" rx="3" fill="currentColor" fill-opacity="0.25"/><rect x="34" y="6" width="22" height="24" rx="3" fill="currentColor" fill-opacity="0.15"/></svg>`,
  grid3: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="4" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/><rect x="22" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/><rect x="40" y="8" width="16" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/></svg>`,
  listSimple: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/><rect x="3" y="3" width="54" height="34" rx="3" stroke="currentColor" stroke-opacity="0.2" stroke-width="1" fill="none"/><rect x="8" y="9" width="6" height="6" rx="1.5" fill="currentColor" fill-opacity="0.22"/><rect x="18" y="10" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/><rect x="18" y="14" width="14" height="1.5" rx="0.75" fill="currentColor" fill-opacity="0.2"/><line x1="3" y1="21" x2="57" y2="21" stroke="currentColor" stroke-opacity="0.12" stroke-width="1"/><rect x="8" y="25" width="6" height="6" rx="1.5" fill="currentColor" fill-opacity="0.22"/><rect x="18" y="26" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.55"/></svg>`,
  listTag: `<svg viewBox="0 0 60 40" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="40" rx="3" fill="currentColor" fill-opacity="0.05"/><rect x="3" y="8" width="18" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/><rect x="6" y="11" width="12" height="2" rx="1" fill="currentColor" fill-opacity="0.45"/><rect x="25" y="8" width="22" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/><rect x="3" y="20" width="14" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/><rect x="21" y="20" width="24" height="8" rx="4" stroke="currentColor" stroke-opacity="0.3" stroke-width="1" fill="currentColor" fill-opacity="0.05"/></svg>`,
}

// 品牌模式预览
const brandModePreviews = {
  both: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="6" y="14" width="10" height="10" rx="2" fill="currentColor" fill-opacity="0.3"/><rect x="19" y="16" width="18" height="3" rx="1.5" fill="currentColor" fill-opacity="0.5"/><rect x="19" y="21" width="12" height="2" rx="1" fill="currentColor" fill-opacity="0.25"/></svg>`,
  logoOnly: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="20" y="11" width="20" height="14" rx="3" fill="currentColor" fill-opacity="0.3"/></svg>`,
  textOnly: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="12" y="15" width="36" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/><rect x="18" y="21" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/></svg>`,
}

// 导航栏风格预览
const navbarStylePreviews = {
  solid: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect width="60" height="11" rx="3" fill="currentColor" fill-opacity="0.18"/><rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/><rect x="36" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/></svg>`,
  blur: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect width="60" height="11" rx="3" fill="currentColor" fill-opacity="0.1"/><rect width="60" height="11" rx="3" fill="url(#blur)" fill-opacity="0.5"/><defs><linearGradient id="blur" x1="0" y1="0" x2="0" y2="11"><stop stop-color="currentColor" stop-opacity="0.25"/><stop offset="1" stop-color="currentColor" stop-opacity="0.05"/></linearGradient></defs><rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/></svg>`,
  transparent: `<svg viewBox="0 0 60 36" fill="none" xmlns="http://www.w3.org/2000/svg"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="5" y="3.5" width="14" height="4" rx="2" fill="currentColor" fill-opacity="0.5"/><rect x="36" y="3.5" width="8" height="4" rx="2" fill="currentColor" fill-opacity="0.3"/></svg>`,
}
</script>

<template>
  <section id="personalization" class="py-16">
    <div class="mx-auto max-w-6xl px-6">

      <!-- Section Header -->
      <div class="mb-14 text-left">
        <p class="mb-3 text-sm text-muted-foreground">第六部分 · 租户首页个性化</p>
        <h2 class="text-4xl font-medium tracking-tight">把首页拆成可配置、可复用的组件体系</h2>
        <p class="mt-3 text-muted-foreground">
          六大区块灵活组合，每种区块提供多种布局与风格选项，在后台编辑器中实时预览配置效果
        </p>

        <!-- 核心特性小卡片 -->
        <div class="mt-3 flex flex-wrap gap-3">
          <div
            v-for="feature in coreFeatures"
            :key="feature.title"
            class="inline-flex items-center gap-2 rounded-full border border-black/[0.08] dark:border-white/[0.08] bg-white dark:bg-[#242424] px-3 py-1.5"
          >
            <div class="flex h-5 w-5 items-center justify-center rounded-full bg-muted/60">
              <component :is="feature.icon" class="h-2.5 w-2.5 text-muted-foreground" />
            </div>
            <span class="text-xs font-medium">{{ feature.title }}</span>
            <span class="text-[10px] text-muted-foreground">{{ feature.desc }}</span>
          </div>

          <!-- 进入后台按钮 -->
          <router-link
            to="/admin/login"
            class="inline-flex items-center gap-1 rounded-full bg-foreground px-3 py-1.5 text-xs font-medium text-background transition-opacity hover:opacity-80"
          >
            进入后台体验
            <ArrowRight class="h-3 w-3" />
          </router-link>
        </div>
      </div>

      <!-- 六大区块卡片 -->
      <div class="mt-14 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <!-- Navbar 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-blue-50">
              <Navigation class="h-4 w-4 text-blue-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">导航栏</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">品牌 Logo、导航链接、主题切换、CTA 按钮</p>

          <!-- 编辑器风格预览：品牌模式 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">品牌区域显示</div>
            <div class="grid grid-cols-3 gap-1.5">
              <div class="rounded-lg border-2 border-primary bg-primary/5 p-1.5 text-center">
                <span class="block w-full text-primary" v-html="brandModePreviews.both" />
                <span class="text-[9px] text-primary">图文</span>
              </div>
              <div class="rounded-lg border border-border p-1.5 text-center opacity-60">
                <span class="block w-full text-muted-foreground" v-html="brandModePreviews.logoOnly" />
                <span class="text-[9px] text-muted-foreground">仅 Logo</span>
              </div>
              <div class="rounded-lg border border-border p-1.5 text-center opacity-60">
                <span class="block w-full text-muted-foreground" v-html="brandModePreviews.textOnly" />
                <span class="text-[9px] text-muted-foreground">仅文字</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：导航栏风格 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">导航栏风格</div>
            <div class="grid grid-cols-3 gap-1.5">
              <div class="rounded-lg border border-border p-1.5 text-center opacity-60">
                <span class="block w-full text-muted-foreground" v-html="navbarStylePreviews.solid" />
                <span class="text-[9px] text-muted-foreground">纯色</span>
              </div>
              <div class="rounded-lg border-2 border-primary bg-primary/5 p-1.5 text-center">
                <span class="block w-full text-primary" v-html="navbarStylePreviews.blur" />
                <span class="text-[9px] text-primary">毛玻璃</span>
              </div>
              <div class="rounded-lg border border-border p-1.5 text-center opacity-60">
                <span class="block w-full text-muted-foreground" v-html="navbarStylePreviews.transparent" />
                <span class="text-[9px] text-muted-foreground">透明</span>
              </div>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">吸顶固定</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">二级导航</span>
          </div>
        </div>

        <!-- Hero 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-purple-50">
              <Layout class="h-4 w-4 text-purple-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">首屏主视觉</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">页面首屏，支持多种布局与入场动画</p>

          <!-- 编辑器风格预览：布局选择 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">布局模式</div>
            <div class="grid grid-cols-2 gap-1.5">
              <div class="rounded-lg border-2 border-primary bg-primary/5 p-2 text-center">
                <span class="block w-full text-primary" v-html="layoutPreviews.centered" />
                <span class="text-[9px] text-primary">居中</span>
              </div>
              <div class="rounded-lg border border-border p-2 text-center opacity-60">
                <span class="block w-full text-muted-foreground" v-html="layoutPreviews.split" />
                <span class="text-[9px] text-muted-foreground">左文右图</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：动画选择 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">入场动画</div>
            <div class="flex gap-1.5">
              <span class="flex-1 rounded-md border border-primary bg-primary/5 px-2 py-1 text-center text-[10px] text-primary">fade-up</span>
              <span class="flex-1 rounded-md border border-border px-2 py-1 text-center text-[10px] text-muted-foreground opacity-60">fade-in</span>
              <span class="flex-1 rounded-md border border-border px-2 py-1 text-center text-[10px] text-muted-foreground opacity-60">打字机</span>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">艺术背景图</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">灰显文字</span>
          </div>
        </div>

        <!-- Introduction 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-green-50">
              <Grid3X3 class="h-4 w-4 text-green-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">特性介绍</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">产品特性展示，支持图标卡片网格</p>

          <!-- 编辑器风格预览：布局 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">布局模式</div>
            <div class="flex gap-1.5">
              <div class="flex-1 rounded-lg border-2 border-primary bg-primary/5 p-1.5 text-center">
                <span class="block w-full text-primary" v-html="layoutPreviews.grid3" />
                <span class="text-[9px] text-primary">3列</span>
              </div>
              <div class="flex-1 rounded-lg border border-border p-1.5 text-center opacity-60">
                <svg viewBox="0 0 60 36" fill="none" class="w-full text-muted-foreground"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="4" y="8" width="24" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/><rect x="32" y="8" width="24" height="20" rx="2" fill="currentColor" fill-opacity="0.18"/></svg>
                <span class="text-[9px] text-muted-foreground">2列</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：卡片风格 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">卡片风格</div>
            <div class="flex gap-1">
              <span class="flex-1 rounded-md border border-primary bg-primary/5 px-1 py-1 text-center text-[9px] text-primary">elevated</span>
              <span class="flex-1 rounded-md border border-border px-1 py-1 text-center text-[9px] text-muted-foreground opacity-60">flat</span>
              <span class="flex-1 rounded-md border border-border px-1 py-1 text-center text-[9px] text-muted-foreground opacity-60">bordered</span>
              <span class="flex-1 rounded-md border border-border px-1 py-1 text-center text-[9px] text-muted-foreground opacity-60">glass</span>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">Lucide 图标</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">图片支持</span>
          </div>
        </div>

        <!-- Theme List 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-orange-50">
              <List class="h-4 w-4 text-orange-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">文档主题</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">主题列表展示，多种列表风格</p>

          <!-- 编辑器风格预览：列表风格 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">列表风格</div>
            <div class="grid grid-cols-2 gap-1.5">
              <div class="rounded-lg border-2 border-primary bg-primary/5 p-1.5">
                <span class="block w-full text-primary" v-html="layoutPreviews.listSimple" />
                <span class="text-[9px] text-primary">简洁行</span>
              </div>
              <div class="rounded-lg border border-border p-1.5 opacity-60">
                <span class="block w-full text-muted-foreground" v-html="layoutPreviews.listTag" />
                <span class="text-[9px] text-muted-foreground">标签云</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：显示选项 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">显示选项</div>
            <div class="space-y-1.5">
              <label class="flex items-center gap-2">
                <span class="h-3.5 w-3.5 rounded border border-primary bg-primary/20 flex items-center justify-center">
                  <Check class="h-2.5 w-2.5 text-primary" />
                </span>
                <span class="text-[10px] text-muted-foreground">显示主题描述</span>
              </label>
              <label class="flex items-center gap-2">
                <span class="h-3.5 w-3.5 rounded border border-border bg-muted/40"></span>
                <span class="text-[10px] text-muted-foreground">显示 Slug</span>
              </label>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">6种风格</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">编号行</span>
          </div>
        </div>

        <!-- CTA 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-pink-50">
              <Megaphone class="h-4 w-4 text-pink-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">行动号召</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">转化引导区块，支持简单或卡片布局</p>

          <!-- 编辑器风格预览：布局 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">布局模式</div>
            <div class="flex gap-1.5">
              <div class="flex-1 rounded-lg border-2 border-primary bg-primary/5 p-1.5 text-center">
                <svg viewBox="0 0 60 36" fill="none" class="w-full text-primary"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="15" y="10" width="30" height="4" rx="2" fill="currentColor" fill-opacity="0.4"/><rect x="18" y="16" width="24" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/><rect x="22" y="22" width="16" height="6" rx="3" fill="currentColor" fill-opacity="0.3"/></svg>
                <span class="text-[9px] text-primary">simple</span>
              </div>
              <div class="flex-1 rounded-lg border border-border p-1.5 text-center opacity-60">
                <svg viewBox="0 0 60 36" fill="none" class="w-full text-muted-foreground"><rect width="60" height="36" rx="3" fill="currentColor" fill-opacity="0.06"/><rect x="4" y="8" width="24" height="4" rx="2" fill="currentColor" fill-opacity="0.4"/><rect x="4" y="14" width="20" height="2.5" rx="1.25" fill="currentColor" fill-opacity="0.25"/><rect x="4" y="22" width="14" height="6" rx="3" fill="currentColor" fill-opacity="0.3"/><rect x="34" y="6" width="22" height="24" rx="3" fill="currentColor" fill-opacity="0.15"/></svg>
                <span class="text-[9px] text-muted-foreground">card</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：按钮样式 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">按钮样式</div>
            <div class="flex gap-1">
              <span class="flex-1 rounded-full bg-foreground px-2 py-1 text-center text-[9px] font-medium text-background">dark</span>
              <span class="flex-1 rounded-full border border-foreground/30 px-2 py-1 text-center text-[9px] text-muted-foreground opacity-60">outline</span>
              <span class="flex-1 rounded-full bg-primary px-2 py-1 text-center text-[9px] font-medium text-primary-foreground opacity-60">primary</span>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">背景图</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">浮动卡片</span>
          </div>
        </div>

        <!-- Footer 卡片 -->
        <div class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl bg-cyan-50">
              <PanelBottom class="h-4 w-4 text-cyan-500" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">页脚</h3>
          </div>
          <p class="text-xs leading-relaxed text-muted-foreground">页面底部，版权与链接信息</p>

          <!-- 编辑器风格预览：输入框 -->
          <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3 space-y-2">
            <div class="space-y-1">
              <span class="text-[10px] text-muted-foreground">宣传标语</span>
              <div class="h-6 rounded-md border border-black/[0.08] bg-white dark:bg-[#2a2a2a] px-2 flex items-center">
                <span class="text-[10px] text-muted-foreground">构建更好的文档体验</span>
              </div>
            </div>
            <div class="space-y-1">
              <span class="text-[10px] text-muted-foreground">版权年份</span>
              <div class="h-6 rounded-md border border-black/[0.08] bg-white dark:bg-[#2a2a2a] px-2 flex items-center">
                <span class="text-[10px] text-muted-foreground">2020-2025</span>
              </div>
            </div>
          </div>

          <!-- 编辑器风格预览：链接列表 -->
          <div class="mt-3 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
            <div class="mb-2 flex items-center justify-between">
              <span class="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">自定义链接</span>
              <Plus class="h-3 w-3 text-muted-foreground/40" />
            </div>
            <div class="space-y-1.5">
              <div class="flex items-center justify-between rounded-md bg-white dark:bg-[#2a2a2a] px-2 py-1">
                <span class="text-[10px] text-muted-foreground">关于我们</span>
                <Trash2 class="h-3 w-3 text-muted-foreground/40" />
              </div>
              <div class="flex items-center justify-between rounded-md bg-white dark:bg-[#2a2a2a] px-2 py-1">
                <span class="text-[10px] text-muted-foreground">联系我们</span>
                <Trash2 class="h-3 w-3 text-muted-foreground/40" />
              </div>
            </div>
          </div>

          <!-- 功能标签 -->
          <div class="mt-3 flex flex-wrap gap-1.5">
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">ICP备案</span>
            <span class="rounded-full border border-black/[0.08] bg-muted/40 px-2 py-0.5 text-[10px] text-muted-foreground">公安备案</span>
          </div>
        </div>
      </div>

    </div>
  </section>
</template>
