<!--
  src/components/home/HomeWorkflow.vue
  职责：首页工作流程展示区，展示文档从创建到发布的完整流程
  对外暴露：无
-->
<script setup lang="ts">
import { PenLine, Eye, Users, Globe } from 'lucide-vue-next'
import { Card, CardBody } from '@/components/ui/a-heroui/card'

const steps = [
  {
    num: '01',
    icon: PenLine,
    title: '在版本结构下创建内容',
    desc: '在租户-主题-版本层级中新增章节与文档页，使用编辑器编写 Markdown，支持代码块、表格、容器与图片上传。',
    color: 'text-blue-500',
    bg: 'bg-blue-50 dark:bg-blue-500/10',
    numColor: 'text-blue-100 dark:text-blue-500/20',
    preview: 'editor',
  },
  {
    num: '02',
    icon: Eye,
    title: '实时预览与自动保存',
    desc: '编辑区与预览区同步渲染，草稿每 30 秒自动保存到后端，确保发布前内容与读者侧展示一致。',
    color: 'text-purple-500',
    bg: 'bg-purple-50 dark:bg-purple-500/10',
    numColor: 'text-purple-100 dark:text-purple-500/20',
    preview: 'preview',
  },
  {
    num: '03',
    icon: Users,
    title: '角色协作与评论审核',
    desc: '超级管理员与租户管理员按权限协作；读者评论提交后默认为 pending，管理员批准或拒绝后再决定是否公开。',
    color: 'text-orange-500',
    bg: 'bg-orange-50 dark:bg-orange-500/10',
    numColor: 'text-orange-100 dark:text-orange-500/20',
    preview: 'collab',
  },
  {
    num: '04',
    icon: Globe,
    title: '单页或整版发布上线',
    desc: '可发布/下线单篇文档，也可批量发布整个版本；本质仅变更状态，无需构建，读者刷新即可生效。',
    color: 'text-green-500',
    bg: 'bg-green-50 dark:bg-green-500/10',
    numColor: 'text-green-100 dark:text-green-500/20',
    preview: 'publish',
  },
]
</script>

<template>
  <section id="workflow" class="py-16">
    <div class="mx-auto max-w-6xl px-6">

      <div class="mb-16 text-left">
        <p class="mb-3 text-sm text-muted-foreground">第四部分 · 发布流程</p>
        <h2 class="text-4xl font-medium tracking-tight">沿着需求链路，四步完成上线</h2>
      </div>

      <!-- 竖向时间轴 -->
      <div class="relative space-y-4">
        <!-- 左侧竖线 -->
        <div class="absolute left-7 top-10 hidden h-[calc(100%-80px)] w-px bg-gradient-to-b from-border via-border to-transparent md:block" />

        <div v-for="step in steps" :key="step.num" class="relative flex gap-6">
          <!-- 左侧序号锚点 -->
          <div class="relative hidden flex-none md:flex">
            <div
              class="flex h-14 w-14 items-center justify-center rounded-2xl border border-border/60 bg-background text-xs font-bold tracking-widest"
              :class="step.color"
            >
              {{ step.num }}
            </div>
          </div>

          <!-- 卡片主体 -->
          <Card class="w-full overflow-hidden border-border/70 outline outline-1 outline-border/60">
            <CardBody class="p-0">
              <div class="flex flex-col md:flex-row">

                <!-- 左：文字区 -->
                <div class="flex flex-col justify-center p-6 md:w-1/2">
                  <!-- 移动端序号 -->
                  <div class="mb-3 flex items-center gap-2 md:hidden">
                    <span class="text-xs font-bold tracking-widest" :class="step.color">{{ step.num }}</span>
                    <span class="h-px flex-1 bg-border/40" />
                  </div>
                  <div class="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-muted/40">
                    <component :is="step.icon" class="h-4 w-4" :class="step.color" />
                  </div>
                  <h3 class="mb-2 text-lg font-medium leading-snug">{{ step.title }}</h3>
                  <p class="text-sm leading-relaxed text-muted-foreground">{{ step.desc }}</p>
                </div>

                <!-- 右：预览区（彩色背景面板） -->
                <div
                  class="relative flex items-center justify-center overflow-hidden border-t border-border/50 bg-muted/20 p-6 md:w-1/2 md:border-l md:border-t-0"
                >

                  <!-- 各步骤预览内容 -->
                  <div class="relative z-10 w-full max-w-xs">

                    <!-- 01 编辑器 -->
                    <template v-if="step.preview === 'editor'">
                      <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                        <div class="flex h-8 items-center gap-1.5 border-b border-black/[0.06] bg-muted/40 px-3">
                          <span class="h-2 w-2 rounded-full bg-red-400" />
                          <span class="h-2 w-2 rounded-full bg-yellow-400" />
                          <span class="h-2 w-2 rounded-full bg-green-400" />
                          <span class="ml-2 text-[9px] text-muted-foreground">quick-start.md</span>
                        </div>
                        <div class="p-4 font-mono space-y-2">
                          <div class="flex gap-1 items-baseline">
                            <span class="text-[11px] text-blue-400">#</span>
                            <div class="h-2.5 w-28 rounded bg-foreground/15" />
                          </div>
                          <div class="h-2 w-full rounded bg-muted/50" />
                          <div class="h-2 w-4/5 rounded bg-muted/50" />
                          <div class="mt-3 flex gap-1 items-baseline">
                            <span class="text-[11px] text-blue-400">##</span>
                            <div class="h-2 w-20 rounded bg-foreground/10" />
                          </div>
                          <div class="flex items-center gap-1.5">
                            <span class="text-[11px] text-green-500">-</span>
                            <div class="h-2 w-24 rounded bg-muted/50" />
                          </div>
                          <div class="flex items-center gap-1.5">
                            <span class="text-[11px] text-green-500">-</span>
                            <div class="h-2 w-16 rounded bg-muted/50" />
                          </div>
                        </div>
                      </div>
                    </template>

                    <!-- 02 预览分屏 -->
                    <template v-if="step.preview === 'preview'">
                      <div class="space-y-3">
                        <div class="flex gap-2">
                          <div class="flex-1 rounded-xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a] p-3">
                            <div class="mb-2 text-[9px] font-semibold uppercase tracking-wider text-muted-foreground/60">编辑</div>
                            <div class="space-y-1.5">
                              <div class="h-1.5 w-full rounded bg-muted/60" />
                              <div class="h-1.5 w-3/4 rounded bg-muted/60" />
                              <div class="h-1.5 w-5/6 rounded bg-muted/60" />
                            </div>
                          </div>
                          <div class="flex-1 rounded-xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a] p-3">
                            <div class="mb-2 text-[9px] font-semibold uppercase tracking-wider text-muted-foreground/60">预览</div>
                            <div class="space-y-1.5">
                              <div class="h-3 w-20 rounded bg-foreground/12" />
                              <div class="h-1.5 w-full rounded bg-muted/40" />
                              <div class="h-1.5 w-3/4 rounded bg-muted/40" />
                            </div>
                          </div>
                        </div>
                        <div class="flex items-center justify-end gap-1.5 rounded-lg bg-white/60 dark:bg-white/5 px-3 py-2 text-[10px] text-green-600 dark:text-green-400">
                          <span class="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse" />
                          草稿已自动保存 · 刚刚
                        </div>
                      </div>
                    </template>

                    <!-- 03 协作成员 -->
                    <template v-if="step.preview === 'collab'">
                      <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                        <div class="border-b border-black/[0.06] px-4 py-3">
                          <div class="h-2.5 w-20 rounded bg-foreground/15" />
                        </div>
                        <div class="divide-y divide-black/[0.04]">
                          <div v-for="(m, i) in [
                            { name: '超级管理员', role: 'Super Admin', avatar: 'bg-blue-200' },
                            { name: '租户管理员', role: 'Tenant Admin', avatar: 'bg-purple-200' },
                            { name: '租户协作管理员', role: 'Tenant Admin', avatar: 'bg-orange-200' },
                            { name: '文档读者', role: 'Viewer', avatar: 'bg-green-200' },
                          ]" :key="i" class="flex items-center justify-between px-4 py-2.5">
                            <div class="flex items-center gap-2.5">
                              <div class="h-6 w-6 rounded-full flex-none" :class="m.avatar" />
                              <span class="text-[11px]">{{ m.name }}</span>
                            </div>
                            <span class="rounded-md border border-black/[0.06] px-2 py-0.5 text-[9px] text-muted-foreground">
                              {{ m.role }}
                            </span>
                          </div>
                        </div>
                      </div>
                    </template>

                    <!-- 04 发布状态 -->
                    <template v-if="step.preview === 'publish'">
                      <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a] p-4 space-y-3">
                        <div class="flex items-center justify-between">
                          <div>
                            <div class="mb-1 h-2.5 w-28 rounded bg-foreground/15" />
                            <div class="h-2 w-20 rounded bg-muted/50" />
                          </div>
                          <div class="flex items-center gap-1.5 rounded-full bg-green-50 dark:bg-green-500/10 border border-green-100 dark:border-green-500/20 px-2.5 py-1 text-[10px] font-medium text-green-600">
                            <span class="h-1.5 w-1.5 rounded-full bg-green-500" />
                            已发布
                          </div>
                        </div>
                        <div class="space-y-1.5 rounded-lg bg-muted/20 p-3 text-[10px]">
                          <div class="flex justify-between text-muted-foreground/60">
                            <span>发布记录</span><span>时间</span>
                          </div>
                          <div v-for="j in 3" :key="j" class="flex items-center justify-between">
                            <div class="flex items-center gap-1.5">
                              <span class="h-1.5 w-1.5 rounded-full" :class="j === 1 ? 'bg-green-400' : 'bg-muted-foreground/30'" />
                              <div class="h-2 w-20 rounded bg-muted/50" />
                            </div>
                            <div class="h-2 w-10 rounded bg-muted/40" />
                          </div>
                        </div>
                        <div class="flex h-8 w-full items-center justify-center rounded-lg bg-foreground">
                          <span class="text-[10px] font-medium text-background">发布文档</span>
                        </div>
                      </div>
                    </template>

                  </div>
                </div>
              </div>
            </CardBody>
          </Card>
        </div>
      </div>

      <!-- 底部 CTA 行 -->
      <div class="mt-12 flex flex-col items-center justify-center gap-3 text-center sm:flex-row sm:items-center">
        <router-link
          to="/admin/login"
          class="inline-flex h-10 items-center rounded-full bg-foreground px-6 text-sm font-medium text-background transition-opacity hover:opacity-80"
        >
          完成四步：进入后台
        </router-link>
        <a
          href="#more-features"
          class="inline-flex h-10 items-center rounded-full border border-black/10 dark:border-white/10 px-6 text-sm font-medium text-foreground transition-colors hover:bg-muted/40"
        >
          继续看第五部分：功能总览
        </a>
      </div>

    </div>
  </section>
</template>
