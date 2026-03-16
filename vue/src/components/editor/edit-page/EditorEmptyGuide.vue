<!--
  EditorEmptyGuide.vue - 编辑器空态引导组件
  职责：当未选择文档页时显示操作引导步骤卡片，供源码模式和可视化模式共用
  对外暴露：无 props/emits，纯展示组件
-->
<script setup lang="ts">
import { FolderPlus, FilePlus, PencilLine, Rocket } from 'lucide-vue-next'

const steps = [
  {
    icon: FolderPlus,
    step: '01',
    title: '添加章节',
    desc: '在右侧侧边栏顶部点击「+ 章节」，为文档内容分组创建结构层级。',
    color: 'text-blue-500',
    bg: 'bg-blue-50 dark:bg-blue-950/40',
    border: 'border-blue-100 dark:border-blue-900',
  },
  {
    icon: FilePlus,
    step: '02',
    title: '添加文档',
    desc: '在章节右侧点击「+」，在该章节下创建文档页并填写标题与 URL Slug。',
    color: 'text-violet-500',
    bg: 'bg-violet-50 dark:bg-violet-950/40',
    border: 'border-violet-100 dark:border-violet-900',
  },
  {
    icon: PencilLine,
    step: '03',
    title: '编辑内容',
    desc: '点击文档页进入编辑器，支持富文本与 Markdown 源码双模式，每 30 秒自动保存草稿。',
    color: 'text-amber-500',
    bg: 'bg-amber-50 dark:bg-amber-950/40',
    border: 'border-amber-100 dark:border-amber-900',
  },
  {
    icon: Rocket,
    step: '04',
    title: '发布上线',
    desc: '在顶部工具栏点击「发布版本」，整个版本下所有草稿文档将批量发布，读者刷新即可看到。',
    color: 'text-emerald-500',
    bg: 'bg-emerald-50 dark:bg-emerald-950/40',
    border: 'border-emerald-100 dark:border-emerald-900',
  },
]
</script>

<template>
  <div class="flex flex-1 items-start justify-center overflow-y-auto py-16 px-6">
    <div class="w-full max-w-xl">

      <!-- 标题区 -->
      <div class="mb-10 flex flex-col items-center text-center">
        <h2 class="text-2xl font-normal tracking-tight text-foreground">开始创建文档</h2>
        <p class="mt-2 max-w-sm text-base font-light text-muted-foreground/70 leading-relaxed">
          按照以下步骤完成内容创建与发布，读者只有在版本发布后才能访问文档。
        </p>
      </div>

      <!-- 时间线步骤 -->
      <div class="relative ml-7">
        <!-- 连接线 -->
        <div class="absolute left-0 top-4 bottom-4 w-px bg-border/40" />

        <div v-for="(s, i) in steps" :key="i" class="relative flex gap-6 pb-8 last:pb-0">
          <!-- 左侧圆点 -->
          <div class="relative z-10 flex shrink-0 -ml-7">
            <span :class="['flex h-14 w-14 items-center justify-center rounded-full border bg-background text-base font-normal tabular-nums transition-colors', s.border, s.color]">
              {{ s.step }}
            </span>
          </div>
          <!-- 右侧内容 -->
          <div class="flex-1 min-w-0 pt-2.5">
            <div class="flex items-center gap-3 mb-2">
              <component :is="s.icon" :class="['h-[22px] w-[22px] shrink-0', s.color]" :stroke-width="1.5" />
              <span class="text-lg font-normal text-foreground">{{ s.title }}</span>
            </div>
            <p class="text-base font-light text-muted-foreground/70 leading-relaxed pl-[34px]">{{ s.desc }}</p>
          </div>
        </div>
      </div>

      <!-- 发布说明提示 -->
      <div class="mt-10 flex items-start gap-3 rounded-xl border border-amber-200/60 dark:border-amber-800/40 bg-amber-50/50 dark:bg-amber-950/20 px-5 py-4">
        <Rocket class="h-5 w-5 shrink-0 mt-0.5 text-amber-500/70" :stroke-width="1.5" />
        <p class="text-sm font-light text-amber-700/80 dark:text-amber-400/70 leading-relaxed">
          <span class="font-normal">注意发布规则：</span>
          单篇文档发布需要所在版本已处于「已发布」状态；若版本为「草稿」，请先在工具栏发布整个版本，版本下所有草稿将批量上线，读者立即可见，无需等待构建。
        </p>
      </div>

    </div>
  </div>
</template>
