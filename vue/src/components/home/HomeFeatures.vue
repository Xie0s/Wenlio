<!--
  src/components/home/HomeFeatures.vue
  职责：首页功能特性展示区，多个 Feature 块交替排列
  对外暴露：无
-->
<script setup lang="ts">
import { FileText, Users, Zap, Lock, Search } from 'lucide-vue-next'

const features = [
  {
    tag: '第二部分 · 文档编辑',
    title: '富文本与 Markdown\n双模式编辑器',
    desc: '编辑器支持 WYSIWYG 与 Markdown 源码双模式切换，内容互转且渲染一致。支持代码块、表格、自定义容器、图片上传与 .md 导入，并每 30 秒自动保存草稿。',
    icon: FileText,
    preview: 'editor',
  },
  {
    tag: '第二部分 · 多租户隔离',
    title: '路径、数据、权限\n四维完全隔离',
    desc: '所有文档路径以 /{tenant_id}/ 开头，数据与搜索都强制按 tenant_id 过滤；JWT 同步携带租户信息，接口层阻断跨租户访问，确保租户之间完全隔离。',
    icon: Users,
    preview: 'team',
  },
  {
    tag: '第二部分 · 发布管理',
    title: '单页或整版发布，\n改状态即生效',
    desc: '发布仅更新状态，无需构建等待。既可对单篇文档执行发布/下线，也可对整版本批量发布，读者端刷新即可获取最新已发布内容。',
    icon: Zap,
    preview: 'publish',
  },
  {
    tag: '第二部分 · 角色与审核',
    title: '三类角色协同，\n评论先审后显',
    desc: '超级管理员、租户管理员、读者权限边界明确。读者可提交评论但默认 pending，管理员在评论管理中批准或拒绝后才公开展示，保障内容质量与可控性。',
    icon: Lock,
    preview: 'permission',
  },
]

const stats = [
  { value: '30s', label: '草稿自动保存周期' },
  { value: '<200ms', label: '目标文档加载时间' },
  { value: 'P95<300ms', label: '接口响应目标' },
  { value: '100', label: '单租户并发读者目标' },
]

const permItems = [
  { label: '管理租户（仅超级管理员）', locked: true, enabled: false },
  { label: '发布/下线文档（租户管理员）', locked: true, enabled: true },
  { label: '审核评论（租户管理员）', locked: true, enabled: true },
  { label: '公开阅读与发表评论（读者）', locked: false, enabled: true },
]
</script>

<template>
  <section id="features" class="py-16">
    <div class="mx-auto max-w-6xl px-6">

      <!-- Section Header -->
      <div class="mb-20 text-left">
        <p class="mb-3 text-sm text-muted-foreground">第二部分 · 核心能力</p>
        <h2 class="text-4xl font-medium tracking-tight">围绕需求闭环的四项核心能力</h2>
        <p class="mt-4 max-w-2xl text-sm leading-relaxed text-muted-foreground">
          从创作、协作、发布到审核，统一在一个界面里完成，降低团队交付与治理成本。
        </p>
      </div>

      <!-- Feature Blocks -->
      <div class="space-y-6">
        <div
          v-for="(feature, index) in features"
          :key="feature.tag"
          class="overflow-hidden rounded-3xl border border-black/[0.06] bg-white dark:bg-[#242424] dark:border-white/[0.06]"
          :class="index % 2 === 0 ? 'md:flex-row' : 'md:flex-row-reverse'"
        >
          <div class="flex flex-col md:flex-row" :class="index % 2 !== 0 ? 'md:flex-row-reverse' : ''">
            <!-- 文字区 -->
            <div class="flex flex-col justify-center p-10 md:w-[42%]">
              <div class="mb-4 inline-flex w-fit items-center gap-1.5 rounded-full border border-black/[0.08] bg-muted/60 px-2.5 py-1 text-xs text-muted-foreground">
                <component :is="feature.icon" class="h-3 w-3" />
                {{ feature.tag }}
              </div>
              <h3 class="mb-4 whitespace-pre-line text-2xl font-medium leading-snug tracking-tight">{{ feature.title }}</h3>
              <p class="text-sm leading-relaxed text-muted-foreground">{{ feature.desc }}</p>
            </div>

            <!-- 预览区 -->
            <div class="flex-1 border-t border-black/[0.06] dark:border-white/[0.06] bg-[#f1eee7] dark:bg-[#181818] p-8 md:border-t-0"
              :class="index % 2 !== 0 ? 'md:border-r md:border-black/[0.06] dark:md:border-white/[0.06]' : 'md:border-l md:border-black/[0.06] dark:md:border-white/[0.06]'"
            >
              <!-- Editor Preview -->
              <template v-if="feature.preview === 'editor'">
                <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                  <div class="flex h-8 items-center gap-2 border-b border-black/[0.06] bg-[#f9f9f9] px-3">
                    <div class="flex gap-1">
                      <span class="h-2 w-2 rounded-full bg-[#ff5f57]" />
                      <span class="h-2 w-2 rounded-full bg-[#febc2e]" />
                      <span class="h-2 w-2 rounded-full bg-[#28c840]" />
                    </div>
                    <div class="flex gap-0.5">
                      <span class="rounded px-2 py-0.5 text-[9px] font-medium bg-black/[0.06] dark:bg-white/[0.08]">编辑</span>
                      <span class="rounded px-2 py-0.5 text-[9px] text-muted-foreground">预览</span>
                    </div>
                  </div>
                  <div class="flex divide-x divide-black/[0.05]">
                    <div class="flex-1 p-4 font-mono">
                      <div class="space-y-1.5">
                        <div class="flex gap-1">
                          <span class="text-[10px] text-blue-400"># </span>
                          <div class="h-2.5 w-28 rounded bg-foreground/15 mt-0.5" />
                        </div>
                        <div class="h-2 w-full rounded bg-muted/50" />
                        <div class="h-2 w-5/6 rounded bg-muted/50" />
                        <div class="mt-3 flex gap-1">
                          <span class="text-[10px] text-blue-400">## </span>
                          <div class="h-2.5 w-20 rounded bg-foreground/10 mt-0.5" />
                        </div>
                        <div class="flex gap-1 items-center">
                          <span class="text-[10px] text-green-500">- </span>
                          <div class="h-2 w-32 rounded bg-muted/50" />
                        </div>
                        <div class="flex gap-1 items-center">
                          <span class="text-[10px] text-green-500">- </span>
                          <div class="h-2 w-24 rounded bg-muted/50" />
                        </div>
                      </div>
                    </div>
                    <div class="flex-1 p-4">
                      <div class="space-y-1.5">
                        <div class="h-3.5 w-28 rounded bg-foreground/15" />
                        <div class="h-2 w-full rounded bg-muted/40" />
                        <div class="h-2 w-5/6 rounded bg-muted/40" />
                        <div class="mt-3 h-2.5 w-20 rounded bg-foreground/10" />
                        <div class="flex items-center gap-1.5">
                          <span class="h-1.5 w-1.5 rounded-full bg-primary/50" />
                          <div class="h-2 w-32 rounded bg-muted/40" />
                        </div>
                        <div class="flex items-center gap-1.5">
                          <span class="h-1.5 w-1.5 rounded-full bg-primary/50" />
                          <div class="h-2 w-24 rounded bg-muted/40" />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <!-- Team Preview -->
              <template v-if="feature.preview === 'team'">
                <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                  <div class="border-b border-black/[0.06] p-4">
                    <div class="mb-3 flex items-center justify-between">
                      <div class="h-3 w-20 rounded bg-foreground/15" />
                      <div class="h-6 w-16 rounded-full bg-primary/15 text-[9px] flex items-center justify-center text-primary font-medium">+ 邀请成员</div>
                    </div>
                    <div class="space-y-2">
                      <div v-for="member in 4" :key="member" class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                          <div class="h-6 w-6 rounded-full" :class="['bg-blue-200','bg-purple-200','bg-green-200','bg-orange-200'][member-1]" />
                          <div class="h-2 w-20 rounded bg-muted/60" />
                        </div>
                        <div class="h-5 w-12 rounded-md border border-black/[0.06] text-[9px] flex items-center justify-center text-muted-foreground">
                          {{ ['管理员','编辑者','阅读者','编辑者'][member-1] }}
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="p-4">
                    <div class="mb-2 h-2.5 w-16 rounded bg-foreground/10" />
                    <div class="grid grid-cols-2 gap-2">
                      <div v-for="space in 4" :key="space" class="rounded-lg border border-black/[0.06] p-2.5">
                        <div class="mb-1.5 h-2 w-12 rounded bg-muted/60" />
                        <div class="h-1.5 w-8 rounded bg-muted/40" />
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <!-- Publish Preview -->
              <template v-if="feature.preview === 'publish'">
                <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                  <div class="p-5">
                    <div class="mb-4 flex items-center justify-between">
                      <div>
                        <div class="mb-1 h-3 w-32 rounded bg-foreground/15" />
                        <div class="h-2 w-24 rounded bg-muted/50" />
                      </div>
                      <div class="flex items-center gap-1.5 rounded-full bg-green-50 px-2.5 py-1 text-[9px] text-green-600 font-medium border border-green-100">
                        <span class="h-1.5 w-1.5 rounded-full bg-green-500" />
                        已发布
                      </div>
                    </div>
                    <div class="space-y-2 rounded-lg bg-muted/20 p-3">
                      <div class="flex items-center justify-between text-[9px] text-muted-foreground">
                        <span>发布记录</span>
                        <span>时间</span>
                      </div>
                      <div v-for="i in 3" :key="i" class="flex items-center justify-between">
                        <div class="flex items-center gap-1.5">
                          <span class="h-1.5 w-1.5 rounded-full" :class="i === 1 ? 'bg-green-400' : 'bg-muted-foreground/30'" />
                          <div class="h-2 w-24 rounded bg-muted/50" />
                        </div>
                        <div class="h-2 w-12 rounded bg-muted/40" />
                      </div>
                    </div>
                    <div class="mt-4 h-8 w-full rounded-lg bg-foreground flex items-center justify-center">
                      <div class="h-2 w-16 rounded bg-background/30" />
                    </div>
                  </div>
                </div>
              </template>

              <!-- Permission Preview -->
              <template v-if="feature.preview === 'permission'">
                <div class="overflow-hidden rounded-2xl border border-black/[0.07] dark:border-white/[0.07] bg-white dark:bg-[#2a2a2a]">
                  <div class="p-5">
                    <div class="mb-4 h-3 w-28 rounded bg-foreground/15" />
                    <div class="space-y-3">
                      <div v-for="perm in permItems" :key="perm.label" class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                          <component :is="perm.locked ? Lock : Search" class="h-3 w-3 text-muted-foreground/60" />
                          <div class="h-2 w-24 rounded bg-muted/60" />
                        </div>
                        <div class="relative h-4 w-7 rounded-full transition-colors"
                          :class="perm.enabled ? 'bg-primary' : 'bg-muted'"
                        >
                          <span class="absolute top-0.5 h-3 w-3 rounded-full bg-white shadow-sm transition-transform"
                            :class="perm.enabled ? 'right-0.5' : 'left-0.5'"
                          />
                        </div>
                      </div>
                    </div>
                    <div class="mt-4 rounded-lg border border-yellow-200 bg-yellow-50 p-3">
                      <div class="mb-1 h-2 w-16 rounded bg-yellow-200" />
                      <div class="h-2 w-full rounded bg-yellow-100" />
                    </div>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Stats -->
      <div class="mt-20 grid grid-cols-2 gap-6 md:grid-cols-4">
        <div v-for="stat in stats" :key="stat.value" class="text-center">
          <div class="mb-1 text-3xl font-medium tracking-tight">{{ stat.value }}</div>
          <div class="text-sm text-muted-foreground">{{ stat.label }}</div>
        </div>
      </div>
    </div>
  </section>
</template>
