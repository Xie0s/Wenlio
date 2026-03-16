<!-- HomepageEditorPage.vue - 租户首页编辑器
     职责：左侧实时预览 + 右侧配置侧边栏的页面构建器布局
     对外接口：无（路由页面组件） -->
<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useHomepageConfig, useHomepageEditor } from '@/composables/personalization'
import { useAuthStore } from '@/stores/auth'
import { useReaderStore } from '@/stores/reader'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { ArrowLeft, Loader2, PanelRightClose, PanelRight, RotateCcw, Save, Upload, Plus, Navigation, ImageIcon, LayoutGrid, Library, Megaphone, PanelBottom, ArrowUpRight } from 'lucide-vue-next'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import { Popover, PopoverAnchor, PopoverContent } from '@/components/ui/popover'
import { SECTION_TYPE_META } from '@/components/personalization/types'
import type { SectionType } from '@/components/personalization/types'
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable'
import HomepageRenderer from '@/components/personalization/renderer/HomepageRenderer.vue'
import HomepageEditorSidebar from '@/components/personalization/editor/HomepageEditorSidebar.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const router = useRouter()
const authStore = useAuthStore()
const readerStore = useReaderStore()
const api = useHomepageConfig()
const editor = useHomepageEditor()
const initializing = ref(true)
const sidebarCollapsed = ref(false)
const needsPublish = ref(false)
const confirmResetOpen = ref(false)

const tenantId = computed(() => authStore.user?.tenant_id ?? '')

onMounted(async () => {
  const data = await api.loadDraft()
  if (data?.draft) {
    editor.initFromLayout(data.draft)
    needsPublish.value = true
  } else if (data?.published) {
    editor.initFromLayout(data.published)
  }
  if (tenantId.value) {
    await readerStore.loadThemes(tenantId.value)
  }
  initializing.value = false
})

const addPopoverOpen = ref(false)

const SECTION_ICONS: Record<string, unknown> = {
  navbar: Navigation,
  hero: ImageIcon,
  introduction: LayoutGrid,
  theme_list: Library,
  cta: Megaphone,
  footer: PanelBottom,
}

const addableSections = computed(() =>
  (Object.keys(SECTION_TYPE_META) as SectionType[]).filter(t => editor.canAddSection(t)),
)

function addSection(type: SectionType) {
  editor.addSection(type)
  addPopoverOpen.value = false
}

async function handleSave() {
  const ok = await api.saveDraft(editor.layout.value, { silentSuccess: true })
  if (ok) {
    editor.markSaved()
    needsPublish.value = true
    toast.success('草稿已保存', {
      description: '读者端不会立即更新，请点击“立即发布”后生效。',
      action: {
        label: '立即发布',
        onClick: () => {
          void handlePublishAfterSave()
        },
      },
    })
  }
}

function viewPage() {
  window.open(`/${tenantId.value}`, '_blank')
}

async function handlePublish() {
  const saved = await api.saveDraft(editor.layout.value, { silentSuccess: true })
  if (!saved) return
  const ok = await api.publish()
  if (ok) {
    editor.markSaved()
    needsPublish.value = false
  }
}

async function handlePublishAfterSave() {
  if (api.publishing.value) return
  const ok = await api.publish()
  if (ok) {
    needsPublish.value = false
  }
}

function handleResetToDefault() {
  editor.resetToDefault()
  confirmResetOpen.value = false
}
</script>

<template>
  <div class="h-screen flex flex-col overflow-hidden bg-background">
    <!-- 顶部工具栏 -->
    <div class="shrink-0 border-b bg-background px-4 py-2.5 flex items-center justify-between gap-4 z-10">
      <div class="flex items-center gap-3">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="rounded-full bg-primary/10 text-primary hover:bg-primary/20"
              @click="router.push('/admin')"
            >
              <ArrowLeft class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>返回管理后台</TooltipContent>
        </Tooltip>
        <h1 class="text-sm font-semibold">首页编辑器</h1>
        <span v-if="editor.dirty.value" class="text-xs text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full">未保存</span>
      </div>
      <div class="flex items-center gap-1">
        <Popover v-model:open="addPopoverOpen">
          <PopoverAnchor as-child>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button
                  variant="ghost"
                  size="icon"
                  class="rounded-full transition-colors"
                  :class="{ 'bg-primary/10 text-primary': addPopoverOpen }"
                  :disabled="!addableSections.length"
                  @click="addPopoverOpen = !addPopoverOpen"
                >
                  <Plus class="h-4 w-4 transition-transform duration-200" :class="{ 'rotate-45': addPopoverOpen }" />
                </Button>
              </TooltipTrigger>
              <TooltipContent v-if="!addPopoverOpen">添加区块</TooltipContent>
            </Tooltip>
          </PopoverAnchor>
          <PopoverContent align="end" :side-offset="8" class="w-44 rounded-2xl p-2">
            <div class="grid grid-cols-2 gap-1.5">
              <button
                v-for="type in addableSections" :key="type"
                class="flex flex-col items-center justify-center gap-1.5 p-2.5 rounded-xl hover:bg-muted transition-colors aspect-square"
                @click="addSection(type)"
              >
                <component :is="SECTION_ICONS[type]" class="h-5 w-5 text-muted-foreground" />
                <span class="text-[11px] leading-tight text-center">{{ SECTION_TYPE_META[type].label }}</span>
              </button>
            </div>
            <p v-if="!addableSections.length" class="py-2 text-xs text-muted-foreground text-center">所有区块已添加</p>
          </PopoverContent>
        </Popover>

        <div class="mx-1 h-5 w-px bg-border" />

        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full" :disabled="api.saving.value" @click="confirmResetOpen = true">
              <RotateCcw class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>重置为默认</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="rounded-full transition-colors"
              :class="{ 'bg-primary/10 text-primary': editor.dirty.value }"
              :disabled="api.saving.value"
              @click="handleSave"
            >
              <Loader2 v-if="api.saving.value" class="h-4 w-4 animate-spin" />
              <Save v-else class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>保存草稿</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="rounded-full transition-colors"
              :class="{ 'bg-primary/10 text-primary': needsPublish }"
              :disabled="api.publishing.value"
              @click="handlePublish"
            >
              <Loader2 v-if="api.publishing.value" class="h-4 w-4 animate-spin" />
              <Upload v-else class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>发布首页</TooltipContent>
        </Tooltip>

        <div class="mx-1 h-5 w-px bg-border" />

        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="rounded-full transition-colors"
              :class="{ 'bg-primary/10 text-primary': sidebarCollapsed }"
              @click="sidebarCollapsed = !sidebarCollapsed"
            >
              <PanelRightClose v-if="!sidebarCollapsed" class="h-4 w-4" />
              <PanelRight v-else class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>{{ sidebarCollapsed ? '展开侧边栏' : '收起侧边栏' }}</TooltipContent>
        </Tooltip>

        <div class="mx-1 h-5 w-px bg-border" />

        <ThemeToggle />

        <div class="mx-1 h-5 w-px bg-border" />

        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="ghost"
              size="icon"
              class="rounded-full bg-primary/10 text-primary hover:bg-primary/20"
              @click="viewPage"
            >
              <ArrowUpRight class="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>查看页面</TooltipContent>
        </Tooltip>
      </div>
    </div>

    <!-- 加载中 -->
    <div v-if="initializing" class="flex-1 flex items-center justify-center">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <!-- 主体：左预览 + 可调宽侧边栏 -->
    <ResizablePanelGroup v-else :key="String(sidebarCollapsed)" direction="horizontal" class="flex-1 min-h-0">
      <!-- 左侧：实时预览区 -->
      <ResizablePanel :default-size="sidebarCollapsed ? 100 : 70" :min-size="40">
        <div class="h-full bg-muted/20">
          <HomepageRenderer
            class="h-full"
            :layout="editor.layout.value"
            :tenant-id="tenantId"
            :tenant-name="authStore.user?.tenant_name ?? '预览'"
            tenant-logo=""
            :themes="readerStore.themes"
            @navigate="() => {}"
          />
        </div>
      </ResizablePanel>

      <!-- 拖拽手柄 + 右侧配置侧边栏 -->
      <template v-if="!sidebarCollapsed">
        <ResizableHandle with-handle class="bg-transparent" />
        <ResizablePanel :default-size="30" :min-size="20" :max-size="50">
          <HomepageEditorSidebar :editor="editor" :themes="readerStore.themes" :tenant-id="tenantId" />
        </ResizablePanel>
      </template>
    </ResizablePanelGroup>

    <ConfirmDialog
      v-model:open="confirmResetOpen"
      type="reset"
      @confirm="handleResetToDefault"
      @cancel="confirmResetOpen = false"
    />
  </div>
</template>
