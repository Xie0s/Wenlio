<!--
  src/components/home/HomeStructure.vue
  职责：首页文档层级结构展示区，可视化 URL 路径规则与租户-主题-版本-页面层级
  对外暴露：无
-->
<script setup lang="ts">
import { ref } from 'vue'
import { ChevronRight, Globe, BookOpen, GitBranch, FileText, Layers } from 'lucide-vue-next'

const activeLevel = ref(3)

const levels = [
  {
    id: 0,
    icon: Globe,
    label: '租户',
    slug: 'acme',
    color: 'text-blue-500',
    bg: 'bg-blue-50',
    border: 'border-blue-100',
    desc: '先锁定租户边界：每个团队拥有独立工作空间，数据完全隔离。',
  },
  {
    id: 1,
    icon: BookOpen,
    label: '主题',
    slug: 'api-reference',
    color: 'text-purple-500',
    bg: 'bg-purple-50',
    border: 'border-purple-100',
    desc: '再按业务拆分主题：例如 API 文档、用户手册，结构清晰可扩展。',
  },
  {
    id: 2,
    icon: GitBranch,
    label: '版本',
    slug: 'v2.0',
    color: 'text-orange-500',
    bg: 'bg-orange-50',
    border: 'border-orange-100',
    desc: '随后进入版本管理：同一主题可并行维护多个版本，读者可自由切换。',
  },
  {
    id: 3,
    icon: FileText,
    label: '文档页',
    slug: 'quick-start',
    color: 'text-green-500',
    bg: 'bg-green-50',
    border: 'border-green-100',
    desc: '最终落到文档页：具体 Markdown 页面可独立发布与下线。',
  },
]

function urlParts(levelId: number) {
  const parts = ['acme', 'api-reference', 'v2.0', 'quick-start']
  return parts.slice(0, levelId + 1)
}

const urlExamples = [
  { path: '/acme/', desc: '第 1 层：租户文档首页（主题列表）', methodClass: 'bg-blue-50 text-blue-600' },
  { path: '/acme/api-reference/', desc: '第 2 层：主题首页（自动跳转默认版本）', methodClass: 'bg-purple-50 text-purple-600' },
  { path: '/acme/api-reference/v2.0/', desc: '第 3 层：指定版本首页', methodClass: 'bg-orange-50 text-orange-600' },
  { path: '/acme/api-reference/v2.0/quick-start', desc: '第 4 层：具体文档页', methodClass: 'bg-green-50 text-green-600' },
]
</script>

<template>
  <section id="structure" class="py-16">
    <div class="mx-auto max-w-6xl px-6">
      <div class="mb-14 text-left">
        <p class="mb-3 text-sm text-muted-foreground">第三部分 · 组织结构</p>
        <h2 class="text-4xl font-medium tracking-tight">把核心能力落到可管理的层级结构</h2>
        <p class="mt-3 text-muted-foreground">四层路径从租户到文档页逐级展开，URL 即结构，结构同时定义数据与权限边界</p>
      </div>

      <div class="grid gap-6 md:grid-cols-2">
        <!-- 左：层级交互卡片 -->
        <div class="flex flex-col gap-3">
          <div
            v-for="level in levels"
            :key="level.id"
            class="cursor-pointer overflow-hidden rounded-2xl border bg-white dark:bg-[#242424] transition-all duration-200"
            :class="activeLevel === level.id ? 'border-foreground/20' : 'border-black/[0.06]'"
            @click="activeLevel = level.id"
          >
            <div class="flex items-start gap-4 p-5">
              <div class="flex h-9 w-9 flex-none items-center justify-center rounded-xl" :class="level.bg">
                <component :is="level.icon" class="h-4 w-4" :class="level.color" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-sm font-medium">{{ level.label }}</span>
                  <span class="rounded-md border px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground" :class="level.border">
                    /{{ level.slug }}
                  </span>
                </div>
                <p class="text-xs text-muted-foreground">{{ level.desc }}</p>
              </div>
              <ChevronRight class="h-4 w-4 flex-none text-muted-foreground/40 transition-transform" :class="activeLevel === level.id ? 'rotate-90' : ''" />
            </div>

            <!-- 展开详情 -->
            <div v-if="activeLevel === level.id" class="border-t border-black/[0.05] bg-muted/20 px-5 py-3">
              <p class="mb-2 text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60">完整 URL</p>
              <div class="flex flex-wrap items-center gap-0.5 font-mono text-xs">
                <span class="text-muted-foreground">https://docs.example.com/</span>
                <span
                  v-for="(part, idx) in urlParts(level.id)"
                  :key="idx"
                  class="flex items-center gap-0.5"
                >
                  <span class="rounded px-1 py-0.5" :class="idx === level.id ? `${level.bg} ${level.color} font-medium` : 'text-muted-foreground'">{{ part }}</span>
                  <span v-if="idx < urlParts(level.id).length - 1" class="text-muted-foreground/50">/</span>
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 右：完整层级树 + URL 示例 -->
        <div class="flex flex-col gap-4">
          <!-- 层级树 -->
          <div class="flex-1 overflow-hidden rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
            <div class="mb-4 flex items-center gap-2">
              <Layers class="h-4 w-4 text-muted-foreground" />
              <span class="text-sm font-medium">文档层级树</span>
            </div>
            <div class="space-y-1 font-mono text-xs">
              <!-- 租户 -->
              <div class="flex items-center gap-2">
                <Globe class="h-3 w-3 text-blue-400" />
                <span class="text-blue-500 font-medium">acme</span>
                <span class="text-muted-foreground/50">（租户）</span>
              </div>
              <!-- 主题 1 -->
              <div class="ml-5 flex items-center gap-2">
                <span class="text-muted-foreground/30">└─</span>
                <BookOpen class="h-3 w-3 text-purple-400" />
                <span class="text-purple-500">api-reference</span>
                <span class="text-muted-foreground/50">（主题）</span>
              </div>
              <!-- 版本 -->
              <div class="ml-10 flex items-center gap-2">
                <span class="text-muted-foreground/30">├─</span>
                <GitBranch class="h-3 w-3 text-orange-400" />
                <span class="text-orange-500">v1.0</span>
                <span class="text-[10px] rounded bg-muted px-1 text-muted-foreground">已归档</span>
              </div>
              <div class="ml-10 flex items-center gap-2">
                <span class="text-muted-foreground/30">└─</span>
                <GitBranch class="h-3 w-3 text-orange-400" />
                <span class="text-orange-500">v2.0</span>
                <span class="text-[10px] rounded bg-green-50 px-1 text-green-600 border border-green-100">默认版本</span>
              </div>
              <!-- 文档页 -->
              <div class="ml-[60px] flex items-center gap-2">
                <span class="text-muted-foreground/30">├─</span>
                <FileText class="h-3 w-3 text-green-400" />
                <span class="text-green-600">quick-start</span>
              </div>
              <div class="ml-[60px] flex items-center gap-2">
                <span class="text-muted-foreground/30">└─</span>
                <FileText class="h-3 w-3 text-green-400" />
                <span class="text-green-600">configuration</span>
              </div>
              <!-- 主题 2 -->
              <div class="ml-5 mt-2 flex items-center gap-2">
                <span class="text-muted-foreground/30">└─</span>
                <BookOpen class="h-3 w-3 text-purple-400" />
                <span class="text-purple-500">user-guide</span>
                <span class="text-muted-foreground/50">（主题）</span>
              </div>
              <div class="ml-10 flex items-center gap-2">
                <span class="text-muted-foreground/30">└─</span>
                <GitBranch class="h-3 w-3 text-orange-400" />
                <span class="text-orange-500">latest</span>
                <span class="text-[10px] rounded bg-green-50 px-1 text-green-600 border border-green-100">默认版本</span>
              </div>
            </div>
          </div>

          <!-- 访问地址卡片 -->
          <div class="overflow-hidden rounded-3xl border border-black/[0.06] dark:border-white/[0.06] bg-white dark:bg-[#242424] p-6">
            <p class="mb-3 text-sm font-medium">访问地址示例</p>
            <div class="space-y-2">
              <div v-for="url in urlExamples" :key="url.path"
                class="flex items-start gap-3 rounded-xl bg-muted/30 px-3 py-2"
              >
                <span class="mt-0.5 flex-none rounded px-1.5 py-0.5 text-[9px] font-medium" :class="url.methodClass">GET</span>
                <div class="min-w-0 flex-1">
                  <code class="block text-[11px] text-foreground/80 truncate">{{ url.path }}</code>
                  <p class="mt-0.5 text-[10px] text-muted-foreground">{{ url.desc }}</p>
                </div>
              </div>
            </div>
            <a
              href="#workflow"
              class="mt-4 inline-flex h-9 items-center rounded-full border border-black/10 px-5 text-sm font-medium text-foreground transition-colors hover:bg-muted/40"
            >
              继续看第四部分：发布流程
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
