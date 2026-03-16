<!-- MediaPage.vue - 媒体文件管理页面
     职责：纯 UI 渲染与用户交互，业务逻辑由 lib/media.ts 提供
     功能：图片卡网格 + 文件列表侧栏分离展示、上传、存储审计 -->
<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { RouterLink } from 'vue-router'
import type { MediaItem } from '@/lib/media'
import { useMedia, formatBytes } from '@/lib/media'
import { useUploadStore } from '@/stores/upload'
import { UploadDropZone, UploadProgressPanel, UploadPreviewDialog } from '@/components/editor/media/update'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import MediaImageCard from '@/components/editor/media/MediaImageCard.vue'
import MediaSidebar from '@/components/editor/media/MediaSidebar.vue'
import {
  Trash2, Images, RefreshCw,
  ImagePlus, Paperclip, ShieldCheck, Loader2,
  AlertTriangle, CheckCircle2,
  Link2, Unlink, BookOpen, Sparkles,
} from 'lucide-vue-next'

const {
  loading, deleting, items,
  loadMedia, deleteMedia,
  auditing, auditResult, auditStorage,
  deletingOrphan, deleteOrphan,
  usageMap, loadingUsage, usageLoaded, loadUsage,
  cleaningUnused, cleanupUnused,
} = useMedia()

const uploadStore = useUploadStore()

const confirmOpen = ref(false)
const pendingId = ref<string | null>(null)
const imageUploadOpen = ref(false)
const fileUploadOpen = ref(false)
const auditPanelOpen = ref(false)
const auditType = ref<'cloud' | 'local'>('cloud')
const usageDetailOpen = ref(false)
const usageDetailItem = ref<MediaItem | null>(null)
const cleanupConfirmOpen = ref(false)

// ── 图片 / 文件分离 ──
const imageItems = computed(() => items.value.filter(i => i.mime_type.startsWith('image/')))
const fileItems = computed(() => items.value.filter(i => !i.mime_type.startsWith('image/')))

function openUsageDetail(item: MediaItem) {
  usageDetailItem.value = item
  usageDetailOpen.value = true
}

async function handleCleanupUnused() {
  cleanupConfirmOpen.value = false
  await cleanupUnused()
}

const unusedCount = computed(() => {
  if (!usageLoaded.value) return 0
  return items.value.filter(i => !usageMap.value[i.id]?.length).length
})

function requestDelete(id: string) {
  pendingId.value = id
  confirmOpen.value = true
}

async function confirmDelete() {
  if (pendingId.value) {
    await deleteMedia(pendingId.value)
  }
  pendingId.value = null
}

// 上传完成回调：刷新列表
const mediaUploadOptions = {
  onComplete: () => { loadMedia() },
}

function handleImageUploadConfirm(files: File[]) {
  uploadStore.addTasks(files, mediaUploadOptions)
}

function handleFileUploadConfirm(files: File[]) {
  uploadStore.addTasks(files, mediaUploadOptions)
}

function handleDropFilesAdded() {
  // 拖拽上传的文件也走 store，完成后刷新
}

function runAudit() {
  auditStorage(auditType.value)
}

watch(auditPanelOpen, (open) => {
  if (open) runAudit()
})

watch(auditType, () => {
  if (auditPanelOpen.value) runAudit()
})

onMounted(loadMedia)
</script>

<template>
  <!-- AdminLayout sidebarLayout 模式提供 h-screen overflow-hidden -->
  <UploadDropZone class="flex flex-col h-full" :upload-options="mediaUploadOptions" @files-added="handleDropFilesAdded">
    <!-- ════════ 顶部工具栏（shrink-0，不参与滚动） ════════ -->
    <div class="shrink-0 px-6 pt-6 pb-3">
    <div class="flex items-center gap-1.5">
      <!-- 上传图片 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full h-9 w-9" @click="imageUploadOpen = true">
            <ImagePlus class="size-4" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>上传图片</TooltipContent>
      </Tooltip>
      <UploadPreviewDialog
        v-model:open="imageUploadOpen"
        accept="image/*"
        title="上传图片"
        @confirm="handleImageUploadConfirm"
      />

      <!-- 上传文件 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full h-9 w-9" @click="fileUploadOpen = true">
            <Paperclip class="size-4" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>上传文件</TooltipContent>
      </Tooltip>
      <UploadPreviewDialog
        v-model:open="fileUploadOpen"
        title="上传文件"
        @confirm="handleFileUploadConfirm"
      />

      <div class="w-px h-5 bg-border mx-0.5" />

      <!-- 刷新 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full h-9 w-9" :disabled="loading" @click="loadMedia">
            <RefreshCw class="size-4" :class="loading ? 'animate-spin' : ''" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>刷新</TooltipContent>
      </Tooltip>

      <!-- 扫描引用来源 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button
            variant="outline" size="icon" class="rounded-full h-9 w-9"
            :class="usageLoaded ? 'border-primary/40 text-primary' : ''"
            :disabled="loadingUsage"
            @click="loadUsage"
          >
            <Loader2 v-if="loadingUsage" class="size-4 animate-spin" :stroke-width="1.5" />
            <Link2 v-else class="size-4" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>扫描引用来源</TooltipContent>
      </Tooltip>

      <!-- 存储审计 -->
      <Tooltip>
        <TooltipTrigger as-child>
          <Button
            variant="outline" size="icon" class="rounded-full h-9 w-9"
            :class="auditPanelOpen ? 'bg-muted' : ''"
            @click="auditPanelOpen = !auditPanelOpen"
          >
            <ShieldCheck class="size-4" :stroke-width="1.5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>存储审计</TooltipContent>
      </Tooltip>

      <!-- 统计信息（紧跟图标组） -->
      <template v-if="!loading && items.length > 0">
        <div class="w-px h-5 bg-border mx-0.5" />
        <div class="flex items-center gap-2 text-xs text-muted-foreground">
          <span>{{ items.length }} 个文件</span>
          <span class="opacity-40">·</span>
          <span>{{ formatBytes(items.reduce((sum: number, m: MediaItem) => sum + m.file_size, 0)) }}</span>
          <template v-if="usageLoaded && unusedCount > 0">
            <span class="opacity-40">·</span>
            <span class="flex items-center gap-1 text-amber-500">
              <Unlink class="size-3" :stroke-width="1.5" />
              {{ unusedCount }} 个未使用
            </span>
            <button
              class="flex items-center gap-1 text-destructive hover:underline disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="cleaningUnused"
              @click="cleanupConfirmOpen = true"
            >
              <Loader2 v-if="cleaningUnused" class="size-3 animate-spin" :stroke-width="1.5" />
              <Sparkles v-else class="size-3" :stroke-width="1.5" />
              {{ cleaningUnused ? '清理中...' : '一键清理' }}
            </button>
          </template>
        </div>
      </template>
    </div>

    <!-- 存储审计展开区（accordion 风格） -->
    <Transition
      enter-active-class="transition-all duration-200 ease-out overflow-hidden"
      leave-active-class="transition-all duration-150 ease-in overflow-hidden"
      enter-from-class="opacity-0 max-h-0"
      enter-to-class="opacity-100 max-h-[600px]"
      leave-from-class="opacity-100 max-h-[600px]"
      leave-to-class="opacity-0 max-h-0"
    >
      <div v-if="auditPanelOpen" class="mb-4 py-3 border-b">
        <!-- 控件 + 说明 -->
        <div class="flex items-center gap-3 mb-4">
          <Tabs v-model="auditType">
            <div class="flex items-center gap-2">
              <TabsList>
                <TabsTrigger value="cloud">云存储</TabsTrigger>
                <TabsTrigger value="local">本地存储</TabsTrigger>
              </TabsList>
              <Loader2 v-if="auditing" class="size-4 animate-spin text-muted-foreground" :stroke-width="1.5" />
            </div>
          </Tabs>
          <span class="text-sm text-muted-foreground">对比数据库与实际存储，检测孤立文件和缺失文件</span>
        </div>

        <!-- 审计结果 -->
        <div v-if="auditResult" class="space-y-3">
          <!-- 汇总 -->
          <div class="flex items-center gap-6 text-sm">
            <span class="flex items-center gap-1.5 text-green-600">
              <CheckCircle2 class="size-4" :stroke-width="2" />
              匹配 {{ auditResult.match_count }}
            </span>
            <span class="flex items-center gap-1.5" :class="auditResult.orphan_keys.length > 0 ? 'text-amber-500' : 'text-muted-foreground'">
              <AlertTriangle class="size-4" :stroke-width="2" />
              孤立 {{ auditResult.orphan_keys.length }}
            </span>
            <span class="flex items-center gap-1.5" :class="auditResult.missing_keys.length > 0 ? 'text-red-500' : 'text-muted-foreground'">
              <AlertTriangle class="size-4" :stroke-width="2" />
              缺失 {{ auditResult.missing_keys.length }}
            </span>
            <span v-if="auditResult.orphan_keys.length === 0 && auditResult.missing_keys.length === 0" class="flex items-center gap-1.5 text-green-600">
              <CheckCircle2 class="size-4" :stroke-width="2" />
              存储状态正常
            </span>
          </div>

          <!-- 孤立文件 -->
          <div v-if="auditResult.orphan_keys.length > 0">
            <p class="text-sm font-medium text-amber-600 mb-1.5 flex items-center gap-1.5">
              <AlertTriangle class="size-4" :stroke-width="2" />
              孤立文件（存储有但数据库无记录）
            </p>
            <div class="max-h-40 overflow-y-auto space-y-1.5">
              <div
                v-for="key in auditResult.orphan_keys"
                :key="key"
                class="flex items-center gap-2 text-sm"
              >
                <span class="truncate flex-1 font-mono text-muted-foreground" :title="key">{{ key }}</span>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button
                      variant="outline" size="icon"
                      class="rounded-full h-7 w-7 shrink-0 text-destructive border-destructive/30 hover:bg-destructive/10"
                      :disabled="deletingOrphan === key"
                      @click="deleteOrphan(auditType, key)"
                    >
                      <Trash2 class="size-3.5" :stroke-width="1.5" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>删除此孤立文件</TooltipContent>
                </Tooltip>
              </div>
            </div>
          </div>

          <!-- 缺失文件 -->
          <div v-if="auditResult.missing_keys.length > 0">
            <p class="text-sm font-medium text-red-600 mb-1.5 flex items-center gap-1.5">
              <AlertTriangle class="size-4" :stroke-width="2" />
              缺失文件（数据库有但存储不存在）
            </p>
            <div class="max-h-40 overflow-y-auto space-y-1.5">
              <p
                v-for="key in auditResult.missing_keys"
                :key="key"
                class="text-sm font-mono text-muted-foreground truncate"
              >
                {{ key }}
              </p>
            </div>
          </div>
        </div>

        <p v-else-if="!auditing" class="text-sm text-muted-foreground">正在准备审计...</p>
      </div>
    </Transition>
    </div>

    <!-- ════════ 主体区域：图片网格 + 文件侧栏 ════════ -->
    <div class="flex flex-1 min-h-0 px-6 pb-6 gap-6">
      <!-- 左侧：图片区（独立滚动） -->
      <div class="flex-1 min-w-0 overflow-y-auto pb-24 pr-3 scrollbar-visible">
        <!-- 加载骨架屏 -->
        <div v-if="loading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3">
          <div v-for="i in 8" :key="i" class="animate-pulse">
            <div class="aspect-square rounded-lg bg-muted" />
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="imageItems.length === 0 && fileItems.length === 0" class="flex flex-col items-center justify-center py-24 text-muted-foreground">
          <Images class="size-12 mb-4 opacity-30" :stroke-width="1" />
          <p class="text-base">暂无媒体文件</p>
          <p class="text-sm mt-1">点击上方按钮上传图片或文件</p>
        </div>

        <!-- 图片网格 -->
        <div v-else-if="imageItems.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3">
          <MediaImageCard
            v-for="item in imageItems"
            :key="item.id"
            :item="item"
            :deleting="deleting"
            :usage-loaded="usageLoaded"
            :usage-refs="usageMap[item.id]"
            @delete="requestDelete"
            @open-usage="openUsageDetail"
          />
        </div>

        <!-- 仅有文件无图片时的提示 -->
        <div v-else class="flex flex-col items-center justify-center py-24 text-muted-foreground">
          <Images class="size-12 mb-4 opacity-30" :stroke-width="1" />
          <p class="text-base">暂无图片</p>
          <p class="text-sm mt-1">文件已显示在右侧面板中</p>
        </div>
      </div>

      <!-- 右侧：文件侧栏（圆角轮廓线） -->
      <MediaSidebar
        :items="fileItems"
        :deleting="deleting"
        :usage-loaded="usageLoaded"
        :usage-map="usageMap"
        @delete="requestDelete"
        @open-usage="openUsageDetail"
      />
    </div>

  <!-- 批量清理未使用文件确认弹窗 -->
  <ConfirmDialog
    v-model:open="cleanupConfirmOpen"
    title="清理未使用文件"
    :description="`将删除 ${unusedCount} 个未被任何文档页引用的文件，同时清理存储空间。此操作无法恢复，请确认。`"
    @confirm="handleCleanupUnused"
  />

  <!-- 删除确认弹窗 -->
  <ConfirmDialog
    v-model:open="confirmOpen"
    title="删除文件"
    description="确认删除该文件？此操作将同时清理存储空间，无法恢复。"
    @confirm="confirmDelete"
  />

  <!-- 使用来源详情弹窗 -->
  <Dialog v-model:open="usageDetailOpen">
    <DialogContent class="rounded-3xl max-w-lg">
      <DialogHeader>
        <DialogTitle>文件引用来源</DialogTitle>
        <DialogDescription v-if="usageDetailItem" class="truncate text-xs">
          {{ usageDetailItem.file_name }}
        </DialogDescription>
      </DialogHeader>
      <div v-if="usageDetailItem && usageMap[usageDetailItem.id]?.length" class="max-h-80 overflow-y-auto space-y-2 mt-2">
        <div
          v-for="usage in usageMap[usageDetailItem.id]"
          :key="usage.page_id"
          class="flex items-start gap-2.5 p-2.5 rounded-2xl border text-sm"
        >
          <BookOpen class="size-4 text-muted-foreground flex-shrink-0 mt-0.5" :stroke-width="1.5" />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-1.5 flex-wrap mb-0.5">
              <Badge variant="destructive" v-if="usage.theme_deleted" class="text-[10px] px-1.5 py-0 h-4">主题已删除</Badge>
              <RouterLink
                v-else
                :to="`/admin/themes/${usage.theme_id}`"
                class="font-medium text-primary hover:underline truncate"
                @click="usageDetailOpen = false"
              >
                {{ usage.theme_name }}
              </RouterLink>
              <span class="text-muted-foreground text-xs">/ {{ usage.version_name }}</span>
            </div>
            <p class="text-xs text-muted-foreground truncate">{{ usage.page_title }}</p>
          </div>
        </div>
      </div>
      <p v-else class="text-sm text-muted-foreground mt-2">该文件暂无引用记录</p>
    </DialogContent>
  </Dialog>

  <!-- 上传进度浮窗 -->
  <UploadProgressPanel />
  </UploadDropZone>
</template>
