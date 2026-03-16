<!--
  EditorToolbar.vue - 编辑器顶部工具栏
  职责：版本选择、保存/发布/导入等操作按钮，注入 AdminHeader #actions slot
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { ref, computed, inject } from 'vue'
import type { Ref } from 'vue'
import { useRouter } from 'vue-router'
import { useEditorInject } from '@/composables/useThemeEditor'
import { openOverlaySafely } from '@/lib/overlay-focus'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ACCESS_MODE_LABEL, ACCESS_MODE_COLOR, type AccessMode } from '@/utils/types'
import {
  ArrowLeft, Check, Navigation, NavigationOff, Upload, Code, Blocks, Settings, X, Leaf,
} from 'lucide-vue-next'

const eyeCareMode = inject<Ref<boolean>>('eyeCareMode', ref(false))
const toggleEyeCare = inject<() => void>('toggleEyeCare', () => { })

const modePopoverOpen = ref(false)
const publishVersionPopoverOpen = ref(false)
const publishingVersion = ref(false)

function switchMode(mode: 'source' | 'visual') {
  editor.editorMode.value = mode
  modePopoverOpen.value = false
}

const router = useRouter()
const editor = useEditorInject()

const fileInput = ref<HTMLInputElement | null>(null)
const unpublishDialogOpen = ref(false)

const activeVersion = computed(() =>
  editor.versions.value.find(v => v.id === editor.activeVersionId.value)
)

function handlePublishClick() {
  if (!activeVersion.value || activeVersion.value.status === 'published') {
    editor.publishPage()
    return
  }
  publishVersionPopoverOpen.value = true
}

async function confirmPublishVersion() {
  if (publishingVersion.value) return
  publishingVersion.value = true
  const ok = await editor.publishVersion(editor.activeVersionId.value)
  if (ok) {
    publishVersionPopoverOpen.value = false
    if (editor.activePageId.value) {
      await editor.selectPage(editor.activePageId.value)
    }
    await editor.loadTree()
  }
  publishingVersion.value = false
}

function openSettingsPanel() {
  void openOverlaySafely(() => {
    editor.showSettingsPanel.value = true
  }, {
    editorElement: document.querySelector('.tiptap.ProseMirror') as HTMLElement | null,
    settleFrameCount: 2,
  })
}

function openUnpublishDialog() {
  void openOverlaySafely(() => {
    unpublishDialogOpen.value = true
  }, {
    editorElement: document.querySelector('.tiptap.ProseMirror') as HTMLElement | null,
    settleFrameCount: 2,
  })
}

// ── 按钮激活状态 ──
const isDirty = computed(() => {
  const p = editor.page.value
  if (!p) return false
  return editor.content.value !== p.content
})

function triggerImport() { fileInput.value?.click() }

async function handleFileImport(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!file.name.endsWith('.md')) { toast.error('请选择 .md 文件'); return }
  const text = await file.text()
  const fmMatch = text.match(/^---\s*\n([\s\S]*?)\n---\s*\n/)
  if (fmMatch?.[1]) {
    const titleMatch = fmMatch[1].match(/^title:\s*["']?(.+?)["']?\s*$/m)
    if (titleMatch?.[1]) editor.title.value = titleMatch[1]
    editor.content.value = text.slice(fmMatch[0]!.length)
  } else { editor.content.value = text }
  toast.success(`已导入 ${file.name}`)
  input.value = ''
}
</script>

<template>
  <div class="flex items-center gap-2">
    <Tooltip>
      <TooltipTrigger as-child>
        <Button size="icon" class="rounded-full h-9 w-9 bg-primary text-primary-foreground hover:bg-primary/90"
          @click="router.push('/admin/themes')">
          <ArrowLeft class="size-[18px]" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>返回主题列表</TooltipContent>
    </Tooltip>

    <!-- 护眼模式 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="outline" size="icon" class="rounded-full h-9 w-9"
          :class="eyeCareMode
            ? 'text-emerald-600 border-emerald-300 bg-emerald-50 dark:text-emerald-400 dark:border-emerald-700 dark:bg-emerald-950/30'
            : 'text-emerald-600/70 border-emerald-200 bg-emerald-50/50 dark:text-emerald-400/60 dark:border-emerald-800/40 dark:bg-emerald-950/15'"
          @click="toggleEyeCare">
          <Leaf class="size-[18px]" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>{{ eyeCareMode ? '关闭护眼' : '护眼模式' }}</TooltipContent>
    </Tooltip>

    <!-- 主题切换 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <ThemeToggle button-size="size-9"
          class="border border-input bg-background hover:bg-accent dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
          icon-size="size-[18px]" />
      </TooltipTrigger>
      <TooltipContent>切换主题</TooltipContent>
    </Tooltip>

    <div class="h-5 w-px bg-border" />

    <!-- 保存 -->
    <Tooltip v-if="editor.page.value">
      <TooltipTrigger as-child>
        <Button size="icon" :variant="isDirty ? 'default' : 'outline'" class="rounded-full h-9 w-9 transition-colors"
          :class="isDirty ? 'bg-primary text-primary-foreground hover:bg-primary/90' : ''"
          :disabled="editor.saving.value" @click="editor.savePage">
          <Check class="size-[18px]" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>保存</TooltipContent>
    </Tooltip>

    <!-- 下线 / 发布 -->
    <template v-if="editor.page.value">
      <Tooltip v-if="editor.page.value.status === 'published'">
        <TooltipTrigger as-child>
          <Button size="icon" variant="outline"
            class="rounded-full h-9 w-9 text-red-600 border-red-300 bg-red-50 hover:bg-red-100 dark:text-red-400 dark:border-red-700 dark:bg-red-950/30 dark:hover:bg-red-950/50"
            @click="openUnpublishDialog">
            <NavigationOff class="size-[18px]" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>下线</TooltipContent>
      </Tooltip>
      <Popover v-if="editor.page.value.status === 'draft'" :open="publishVersionPopoverOpen"
        @update:open="(open: boolean) => { if (!open) publishVersionPopoverOpen = false }">
        <PopoverTrigger as-child>
          <Button size="icon" :class="[
            'rounded-full h-9 w-9 bg-primary text-primary-foreground hover:bg-primary/90',
            publishVersionPopoverOpen ? 'ring-2 ring-primary/30' : '',
          ]" @click="handlePublishClick">
            <Navigation class="size-[18px]" />
          </Button>
        </PopoverTrigger>
        <PopoverContent align="end" class="w-72 rounded-2xl p-3">
          <p class="text-base font-medium">所属版本尚未发布，是否同时发布版本？</p>
          <p class="mt-1 text-sm text-muted-foreground">确认后将发布当前版本，该版本下所有草稿文档页将一并上线，所有读者均可立即访问。如需仅发布单篇文档，请先前往设置面板发布版本。
          </p>
          <div class="mt-3 flex items-center justify-end gap-2">
            <Button variant="outline" size="icon" class="h-7 w-7 rounded-full" :disabled="publishingVersion"
              @click="publishVersionPopoverOpen = false">
              <X class="h-3.5 w-3.5" />
            </Button>
            <Button size="icon" class="h-7 w-7 rounded-full bg-primary text-primary-foreground hover:bg-primary/90"
              :disabled="publishingVersion" @click="confirmPublishVersion">
              <Check class="h-3.5 w-3.5" />
            </Button>
          </div>
        </PopoverContent>
      </Popover>
    </template>

    <!-- 上传 -->
    <template v-if="editor.page.value">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="outline" size="icon" class="rounded-full h-9 w-9" @click="triggerImport">
            <Upload class="size-[18px]" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>导入 Markdown</TooltipContent>
      </Tooltip>
      <input ref="fileInput" type="file" accept=".md" class="hidden" @change="handleFileImport" />
    </template>

    <!-- 设置 -->
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="outline" size="icon" class="rounded-full h-9 w-9" @click="openSettingsPanel">
          <Settings class="size-[18px]" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>设置</TooltipContent>
    </Tooltip>

    <!-- 访问模式指示（仅非 public 时显示） -->
    <Tooltip v-if="editor.theme.value?.access_mode && editor.theme.value.access_mode !== 'public'">
      <TooltipTrigger as-child>
        <Badge :class="ACCESS_MODE_COLOR[editor.theme.value.access_mode as AccessMode]" class="text-[11px] px-2 py-0.5 cursor-default">
          {{ ACCESS_MODE_LABEL[editor.theme.value.access_mode as AccessMode] }}
        </Badge>
      </TooltipTrigger>
      <TooltipContent>当前主题访问权限</TooltipContent>
    </Tooltip>

    <!-- 切换编辑模式 -->
    <template v-if="editor.page.value">
      <Popover v-model:open="modePopoverOpen">
        <PopoverTrigger as-child>
          <Button variant="outline" class="rounded-full h-9 px-3 gap-1.5 text-sm transition-colors"
            :class="modePopoverOpen ? 'bg-primary text-primary-foreground border-primary hover:bg-primary/90 hover:text-primary-foreground' : ''">
            <Code v-if="editor.editorMode.value === 'source'" class="size-[18px]" />
            <Blocks v-else class="size-[18px]" />
            <span>{{ editor.editorMode.value === 'source' ? '源码' : '可视化' }}</span>
          </Button>
        </PopoverTrigger>
        <PopoverContent class="w-auto p-2 rounded-full" :side-offset="8">
          <Tabs :model-value="editor.editorMode.value" @update:model-value="switchMode($event as 'source' | 'visual')">
            <TabsList>
              <TabsTrigger value="source" class="gap-1.5">
                <Code class="h-[15px] w-[15px]" />
                源码编辑
              </TabsTrigger>
              <TabsTrigger value="visual" class="gap-1.5">
                <Blocks class="h-[15px] w-[15px]" />
                可视化
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </PopoverContent>
      </Popover>
    </template>

    <!-- 版本选择 -->
    <Select v-if="editor.versions.value.length > 0" :model-value="editor.activeVersionId.value"
      @update:model-value="editor.switchVersion($event as string)">
      <SelectTrigger class="h-9 w-auto min-w-[110px] text-sm rounded-full px-4">
        <SelectValue placeholder="选择版本" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem v-for="v in editor.versions.value" :key="v.id" :value="v.id">
          {{ v.name }}
          <Badge v-if="v.is_default" variant="outline" class="ml-1 text-[10px] px-1 py-0">默认</Badge>
        </SelectItem>
      </SelectContent>
    </Select>
  </div>

  <!-- 下线确认弹窗 -->
  <ConfirmDialog v-model:open="unpublishDialogOpen" type="ban" title="确认下线" description="下线后该文档将对用户不可见，可随时重新发布。"
    @confirm="editor.unpublishPage" />
</template>
