<!--
  src/components/home/HomeMoreFeatures.vue
  职责：首页功能特性补充区，以小卡片网格展示版本管理、全文搜索、评论、自动保存等细节功能
  对外暴露：无
-->
<script setup lang="ts">
import {
  GitBranch,
  Search,
  MessageSquare,
  Save,
  Upload,
  Code2,
  Box,
  Palette,
  BookMarked,
} from 'lucide-vue-next'

const versionItems = [
  { name: 'v2.0', status: '已发布', default: true, statusClass: 'bg-green-50 text-green-600 border border-green-100' },
  { name: 'v1.5', status: '已发布', default: false, statusClass: 'bg-green-50 text-green-600 border border-green-100' },
  { name: 'v1.0', status: '已归档', default: false, statusClass: 'bg-muted text-muted-foreground' },
  { name: 'v3.0-draft', status: '草稿', default: false, statusClass: 'bg-yellow-50 text-yellow-600 border border-yellow-100' },
]

const searchResults = [
  { title: '快速开始指南', snippet: '发布后读者访问 <mark class="bg-yellow-100 rounded px-0.5">/{tenant_id}/{theme}/{version}/{page}</mark> 即可阅读...' },
  { title: 'API 认证', snippet: '管理接口使用 <mark class="bg-yellow-100 rounded px-0.5">JWT Token</mark> 鉴权，公开搜索接口无需登录...' },
]

const commentItems = [
  { name: '张三', content: '文档写得很清楚，解决了我的问题！', status: '已批准', statusClass: 'bg-green-50 text-green-600 border border-green-100' },
  { name: '李四', content: '第三步的配置有些不明确，能补充一下吗？', status: '待审核', statusClass: 'bg-yellow-50 text-yellow-600 border border-yellow-100' },
]

const tocPreviewItems = [
  { label: '介绍', indent: '', active: false },
  { label: '快速开始', indent: '', active: false },
  { label: '安装依赖', indent: 'pl-3', active: false },
  { label: '配置文件', indent: 'pl-3', active: true },
  { label: '环境变量', indent: 'pl-6', active: false },
  { label: 'API 参考', indent: '', active: false },
  { label: '常见问题', indent: '', active: false },
]

const cards = [
  {
    icon: GitBranch,
    title: '多版本管理',
    desc: '同一主题维护 v1.0、v2.0 等多个版本，支持克隆、发布、归档，读者可自由切换版本阅读。',
    color: 'text-orange-500',
    bg: 'bg-orange-50',
    preview: 'version',
  },
  {
    icon: Search,
    title: '全文搜索',
    desc: '基于 MongoDB 文本索引，标题权重高于内容。搜索接口按 tenant_id 强制隔离，并支持按主题与版本过滤。',
    color: 'text-blue-500',
    bg: 'bg-blue-50',
    preview: 'search',
  },
  {
    icon: MessageSquare,
    title: '文档评论',
    desc: '读者可在文档页提交评论，默认 pending。管理员审核通过后才公开，支持一层嵌套回复。',
    color: 'text-purple-500',
    bg: 'bg-purple-50',
    preview: 'comment',
  },
  {
    icon: Save,
    title: '自动保存草稿',
    desc: '编辑器每 30 秒自动保存草稿到后端，防止意外丢失。草稿与已发布版本独立存储。',
    color: 'text-green-500',
    bg: 'bg-green-50',
    preview: 'autosave',
  },
  {
    icon: Upload,
    title: '导入 Markdown',
    desc: '支持上传本地 .md 文件，自动解析 frontmatter 中的 title、description 字段，一键导入。',
    color: 'text-teal-500',
    bg: 'bg-teal-50',
    preview: 'import',
  },
  {
    icon: Code2,
    title: '200+ 语言代码高亮',
    desc: '集成 Shiki 代码高亮引擎，支持 200+ 编程语言。浅色/深色主题随系统自动切换，无需重渲染。',
    color: 'text-indigo-500',
    bg: 'bg-indigo-50',
    preview: 'code',
  },
  {
    icon: Box,
    title: '自定义容器块',
    desc: '支持 ::: tip / warning / danger / info 容器语法，前端配套样式后与阅读界面渲染保持一致。',
    color: 'text-yellow-600',
    bg: 'bg-yellow-50',
    preview: 'container',
  },
  {
    icon: Palette,
    title: '租户主题定制',
    desc: '每个租户独立配置品牌色、Logo 和自定义 CSS，文档阅读界面完全呈现租户品牌风格。',
    color: 'text-pink-500',
    bg: 'bg-pink-50',
    preview: 'theme',
  },
  {
    icon: BookMarked,
    title: '目录自动生成',
    desc: '文档右侧浮动目录由 markdown-it-toc-done-right 自动生成，随滚动高亮当前章节，支持深层嵌套。',
    color: 'text-cyan-500',
    bg: 'bg-cyan-50',
    preview: 'toc',
  },
]
</script>

<template>
  <section id="more-features" class="py-16">
    <div class="mx-auto max-w-6xl px-6">
      <div class="mb-14 text-left">
        <p class="mb-3 text-sm text-muted-foreground">第五部分 · 功能总览</p>
        <h2 class="text-4xl font-medium tracking-tight">补齐细节能力，形成完整体验</h2>
        <p class="mt-3 text-muted-foreground">在核心流程之外，版本、搜索、评论审核与主题定制共同支撑完整文档体系</p>
      </div>

      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="card in cards"
          :key="card.title"
          class="rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6"
        >
          <!-- 图标 + 标题 -->
          <div class="mb-4 flex items-start gap-3">
            <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl" :class="card.bg">
              <component :is="card.icon" class="h-4 w-4" :class="card.color" />
            </div>
            <h3 class="mt-1.5 text-sm font-medium leading-snug">{{ card.title }}</h3>
          </div>

          <p class="text-xs leading-relaxed text-muted-foreground">{{ card.desc }}</p>

          <!-- 版本管理预览 -->
          <template v-if="card.preview === 'version'">
            <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20">
              <div class="flex items-center justify-between border-b border-black/[0.05] dark:border-white/[0.05] px-3 py-2">
                <span class="text-[10px] font-medium text-muted-foreground">版本列表</span>
              </div>
              <div class="divide-y divide-black/[0.04]">
                <div v-for="v in versionItems" :key="v.name" class="flex items-center justify-between px-3 py-2">
                  <div class="flex items-center gap-2">
                    <GitBranch class="h-3 w-3 text-muted-foreground/40" />
                    <span class="text-[11px] font-mono">{{ v.name }}</span>
                    <span v-if="v.default" class="rounded bg-green-50 px-1 text-[9px] text-green-600 border border-green-100">默认</span>
                  </div>
                  <span class="rounded px-1.5 text-[9px]" :class="v.statusClass">{{ v.status }}</span>
                </div>
              </div>
            </div>
          </template>

          <!-- 搜索预览 -->
          <template v-if="card.preview === 'search'">
            <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
              <div class="mb-2 flex items-center gap-2 rounded-lg border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a] px-2.5 py-1.5">
                <Search class="h-3 w-3 text-muted-foreground/50" />
                <span class="text-[11px] text-muted-foreground/60">搜索文档...</span>
              </div>
              <div class="space-y-1.5">
                <div v-for="r in searchResults" :key="r.title" class="rounded-lg bg-white dark:bg-[#2a2a2a] p-2 text-[10px]">
                  <div class="mb-0.5 font-medium">{{ r.title }}</div>
                  <div class="text-muted-foreground" v-html="r.snippet" />
                </div>
              </div>
            </div>
          </template>

          <!-- 评论预览 -->
          <template v-if="card.preview === 'comment'">
            <div class="mt-4 space-y-2">
              <div v-for="c in commentItems" :key="c.name" class="rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
                <div class="mb-1 flex items-center justify-between">
                  <span class="text-[11px] font-medium">{{ c.name }}</span>
                  <span class="rounded px-1.5 text-[9px]" :class="c.statusClass">{{ c.status }}</span>
                </div>
                <p class="text-[10px] text-muted-foreground">{{ c.content }}</p>
              </div>
            </div>
          </template>

          <!-- 自动保存预览 -->
          <template v-if="card.preview === 'autosave'">
            <div class="mt-4 rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-muted/20 p-3">
              <div class="mb-2 flex items-center justify-between">
                <span class="text-[10px] text-muted-foreground">编辑中</span>
                <div class="flex items-center gap-1 text-[10px] text-green-600">
                  <span class="h-1.5 w-1.5 rounded-full bg-green-400" />
                  已自动保存
                </div>
              </div>
              <div class="space-y-1">
                <div class="h-2 w-full rounded bg-muted/50" />
                <div class="h-2 w-5/6 rounded bg-muted/50" />
                <div class="h-2 w-4/6 rounded bg-muted/50" />
              </div>
              <div class="mt-2 text-right text-[9px] text-muted-foreground/50">上次保存：刚刚</div>
            </div>
          </template>

          <!-- 导入预览 -->
          <template v-if="card.preview === 'import'">
            <div class="mt-4 rounded-xl border border-dashed border-black/[0.12] bg-muted/10 p-4 text-center">
              <Upload class="mx-auto mb-2 h-5 w-5 text-muted-foreground/40" />
              <p class="text-[11px] text-muted-foreground">拖放 .md 文件到此处</p>
              <p class="mt-0.5 text-[9px] text-muted-foreground/50">自动解析 frontmatter</p>
            </div>
          </template>

          <!-- 代码高亮预览 -->
          <template v-if="card.preview === 'code'">
            <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] bg-[#1e1e2e] p-3 font-mono text-[10px]">
              <div class="mb-1.5 text-[9px] text-white/30">typescript</div>
              <div class="space-y-0.5">
                <div><span class="text-blue-400">const</span> <span class="text-green-300">greet</span> <span class="text-white/60">=</span> <span class="text-yellow-300">(name: string)</span> <span class="text-blue-400">=&gt;</span> <span class="text-white/80">{</span></div>
                <div class="pl-4"><span class="text-purple-400">return</span> <span class="text-orange-300">`Hello, <span class="text-green-300">${name}</span>!`</span></div>
                <div><span class="text-white/80">}</span></div>
              </div>
            </div>
          </template>

          <!-- 容器块预览 -->
          <template v-if="card.preview === 'container'">
            <div class="mt-4 space-y-2">
              <div class="rounded-lg border-l-2 border-blue-400 bg-blue-50 px-3 py-2">
                <p class="text-[10px] font-medium text-blue-700">TIP</p>
                <p class="text-[10px] text-blue-600">这是一个 tip 提示框</p>
              </div>
              <div class="rounded-lg border-l-2 border-yellow-400 bg-yellow-50 px-3 py-2">
                <p class="text-[10px] font-medium text-yellow-700">WARNING</p>
                <p class="text-[10px] text-yellow-600">这是一个 warning 警告框</p>
              </div>
            </div>
          </template>
          <!-- 租户主题定制预览 -->
          <template v-if="card.preview === 'theme'">
            <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] bg-muted/20">
              <div class="border-b border-black/[0.05] p-3">
                <div class="mb-2 text-[10px] font-medium text-muted-foreground">主题配置</div>
                <div class="flex items-center gap-2">
                  <div class="h-5 w-5 rounded-full bg-[#6366f1]" />
                  <div class="h-5 w-5 rounded-full bg-[#0ea5e9]" />
                  <div class="h-5 w-5 rounded-full bg-[#10b981]" />
                  <div class="h-5 w-5 rounded-full bg-[#f59e0b]" />
                  <div class="ml-1 h-5 w-5 rounded-full ring-2 ring-offset-1 ring-[#ec4899] bg-[#ec4899]" />
                  <span class="text-[9px] text-muted-foreground">当前</span>
                </div>
              </div>
              <div class="p-3">
                <div class="mb-1.5 flex items-center gap-2 rounded-lg bg-white px-2 py-1.5">
                  <div class="h-4 w-4 rounded bg-[#ec4899]/20" />
                  <div class="h-2 w-16 rounded bg-muted/60" />
                </div>
                <div class="flex items-center gap-1.5 px-2">
                  <div class="h-1.5 w-1.5 rounded-full bg-[#ec4899]/60" />
                  <div class="h-2 w-20 rounded bg-muted/40" />
                </div>
                <div class="mt-1 flex items-center gap-1.5 px-2">
                  <div class="h-1.5 w-1.5 rounded-full bg-[#ec4899]/60" />
                  <div class="h-2 w-14 rounded bg-muted/40" />
                </div>
              </div>
            </div>
          </template>

          <!-- 目录自动生成预览 -->
          <template v-if="card.preview === 'toc'">
            <div class="mt-4 overflow-hidden rounded-xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#2a2a2a] p-3">
              <div class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/50">目录</div>
              <div class="space-y-1">
                <div v-for="item in tocPreviewItems" :key="item.label"
                  class="flex items-center rounded-md px-1.5 py-1 text-[10px]"
                  :class="[item.indent, item.active ? 'bg-cyan-50 text-cyan-600 font-medium' : 'text-muted-foreground']"
                >
                  <span v-if="item.active" class="mr-1.5 h-1 w-1 rounded-full bg-cyan-500 flex-none" />
                  <span v-else class="mr-1.5 h-1 w-1 rounded-full bg-transparent flex-none" />
                  {{ item.label }}
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </section>
</template>
