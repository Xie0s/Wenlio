<!-- HomepageEditorSidebar.vue - 首页编辑器右侧侧边栏
     职责：提供区块列表管理（折叠展开编辑、排序、显隐、增删）和全局样式编辑
     对外接口：
       Props: editor (EditorState) -->
<script setup lang="ts">
import { ref } from 'vue'
import type { Ref, ComputedRef } from 'vue'
import type { Theme } from '@/utils/types'
import type { HomepageSection, SectionType, HomepageLayout, HomepageGlobal } from '../types'
import { SECTION_TYPE_META } from '../types'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  ChevronUp, ChevronDown, Trash2,
  Eye, EyeOff, ChevronRight, Settings2,
} from 'lucide-vue-next'
import NavbarEditor from './NavbarEditor.vue'
import HeroEditor from './HeroEditor.vue'
import IntroductionEditor from './IntroductionEditor.vue'
import ThemeListEditor from './ThemeListEditor.vue'
import CtaEditor from './CtaEditor.vue'
import FooterEditor from './FooterEditor.vue'
import GlobalEditor from './GlobalEditor.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

interface EditorState {
  layout: Ref<HomepageLayout>
  sections: ComputedRef<HomepageSection[]>
  selectedSectionId: Ref<string | null>
  selectedSection: ComputedRef<HomepageSection | null>
  dirty: Ref<boolean>
  canAddSection: (type: SectionType) => boolean
  addSection: (type: SectionType, index?: number) => HomepageSection | null
  removeSection: (id: string) => void
  toggleVisibility: (id: string) => void
  moveUp: (id: string) => void
  moveDown: (id: string) => void
  updateGlobal: (global: HomepageGlobal) => void
  updateSectionConfig: (id: string, config: HomepageSection['config']) => void
}

const props = defineProps<{
  editor: EditorState
  themes?: Theme[]
  tenantId?: string
}>()

// 全局设置面板展开状态
const globalExpanded = ref(false)

// 删除确认弹窗
const confirmOpen = ref(false)
const pendingDeleteId = ref<string | null>(null)

function askRemoveSection(id: string) {
  pendingDeleteId.value = id
  confirmOpen.value = true
}

function confirmRemove() {
  if (pendingDeleteId.value) {
    props.editor.removeSection(pendingDeleteId.value)
    pendingDeleteId.value = null
  }
  confirmOpen.value = false
}

function toggleSection(id: string) {
  props.editor.selectedSectionId.value =
    props.editor.selectedSectionId.value === id ? null : id
}

function handleConfigUpdate(section: HomepageSection, config: HomepageSection['config']) {
  props.editor.updateSectionConfig(section.id, config)
}

function handleGlobalUpdate(config: HomepageGlobal) {
  props.editor.updateGlobal(config)
}

</script>

<template>
  <aside class="h-full flex flex-col bg-background border-l overflow-hidden">
    <!-- 侧边栏头部 -->
    <div class="shrink-0 px-4 py-2.5 border-b flex items-center justify-center">
      <span class="text-sm font-semibold">页面配置</span>
    </div>

    <!-- 可滚动内容区 -->
    <div class="flex-1 overflow-y-auto">
      <!-- 全局样式折叠面板 -->
      <div class="border-b">
        <button
          class="w-full flex items-center gap-2.5 px-4 py-3 text-sm transition-colors"
          :class="globalExpanded ? 'bg-primary/10 text-primary' : 'hover:bg-muted/50'"
          @click="globalExpanded = !globalExpanded"
        >
          <ChevronRight
            class="h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200"
            :class="{ 'rotate-90': globalExpanded }"
          />
          <Settings2 class="h-4 w-4 shrink-0 text-muted-foreground" />
          <span class="font-medium">全局样式</span>
        </button>
        <div v-if="globalExpanded" class="px-4 pb-4">
          <GlobalEditor :config="editor.layout.value.global" @update="handleGlobalUpdate" />
        </div>
      </div>

      <!-- 区块列表 -->
      <div class="divide-y">
        <div v-for="(section, i) in editor.sections.value" :key="section.id">
          <!-- 区块行头 -->
          <div
            class="flex items-center gap-1.5 px-4 py-2.5 transition-colors"
            :class="{
              'bg-primary/10': editor.selectedSectionId.value === section.id,
              'opacity-50': !section.visible,
            }"
          >
            <!-- 展开/收起 -->
            <button class="p-0.5 hover:bg-muted rounded-lg transition-colors" @click="toggleSection(section.id)">
              <ChevronRight
                class="h-4 w-4 text-muted-foreground transition-transform duration-200"
                :class="{ 'rotate-90': editor.selectedSectionId.value === section.id }"
              />
            </button>

            <!-- 区块名称（点击展开） -->
            <button class="flex-1 text-left text-sm font-medium truncate hover:text-primary transition-colors" @click="toggleSection(section.id)">
              {{ SECTION_TYPE_META[section.type].label }}
            </button>

            <!-- 操作按钮组 -->
            <div class="flex items-center shrink-0">
              <Tooltip>
                <TooltipTrigger as-child>
                  <button
                    class="p-1 rounded-lg hover:bg-muted transition-colors disabled:opacity-30"
                    :disabled="i === 0"
                    @click="editor.moveUp(section.id)"
                  >
                    <ChevronUp class="h-3.5 w-3.5 text-muted-foreground" />
                  </button>
                </TooltipTrigger>
                <TooltipContent>上移</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger as-child>
                  <button
                    class="p-1 rounded-lg hover:bg-muted transition-colors disabled:opacity-30"
                    :disabled="i === editor.sections.value.length - 1"
                    @click="editor.moveDown(section.id)"
                  >
                    <ChevronDown class="h-3.5 w-3.5 text-muted-foreground" />
                  </button>
                </TooltipTrigger>
                <TooltipContent>下移</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger as-child>
                  <button class="p-1 rounded-lg hover:bg-muted transition-colors" @click="editor.toggleVisibility(section.id)">
                    <EyeOff v-if="!section.visible" class="h-3.5 w-3.5 text-muted-foreground" />
                    <Eye v-else class="h-3.5 w-3.5 text-muted-foreground" />
                  </button>
                </TooltipTrigger>
                <TooltipContent>{{ section.visible ? '隐藏' : '显示' }}</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger as-child>
                  <button class="p-1 rounded-lg hover:bg-destructive/10 transition-colors" @click="askRemoveSection(section.id)">
                    <Trash2 class="h-3.5 w-3.5 text-destructive/70" />
                  </button>
                </TooltipTrigger>
                <TooltipContent>删除区块</TooltipContent>
              </Tooltip>
            </div>
          </div>

          <!-- 区块编辑表单（折叠内容） -->
          <div v-if="editor.selectedSectionId.value === section.id" class="px-4 pb-4 pt-1 bg-muted/20">
            <NavbarEditor v-if="section.type === 'navbar'" :config="(section.config as any)" :themes="themes" :tenant-id="tenantId" @update="handleConfigUpdate(section, $event)" />
            <HeroEditor v-else-if="section.type === 'hero'" :config="(section.config as any)" @update="handleConfigUpdate(section, $event)" />
            <IntroductionEditor v-else-if="section.type === 'introduction'" :config="(section.config as any)" @update="handleConfigUpdate(section, $event)" />
            <ThemeListEditor v-else-if="section.type === 'theme_list'" :config="(section.config as any)" @update="handleConfigUpdate(section, $event)" />
            <CtaEditor v-else-if="section.type === 'cta'" :config="(section.config as any)" @update="handleConfigUpdate(section, $event)" />
            <FooterEditor v-else-if="section.type === 'footer'" :config="(section.config as any)" @update="handleConfigUpdate(section, $event)" />
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!editor.sections.value.length" class="px-4 py-8 text-center text-sm text-muted-foreground">
        暂无区块，请添加
      </div>
    </div>

  </aside>

  <ConfirmDialog
    v-model:open="confirmOpen"
    type="delete"
    description="删除后该区块配置将丢失，确认继续？"
    @confirm="confirmRemove"
    @cancel="confirmOpen = false"
  />
</template>
