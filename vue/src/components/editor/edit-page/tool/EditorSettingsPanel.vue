<!--
  EditorSettingsPanel.vue - 编辑器设置模态框
  职责：左右分栏一页显示主题设置、版本管理、章节/文档管理
  对外暴露：无 props/emits，通过 inject 消费编辑器上下文
-->
<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Version } from '@/utils/types'
import { VERSION_STATUS_LABEL, VERSION_STATUS_COLOR } from '@/utils/types'
import { useEditorInject } from '@/composables/useThemeEditor'
import { useDocTreeCrud } from '@/composables/useDocTreeCrud'
import { validateSlug, validateVersionName, slugHint, versionHint } from '@/lib/validation'
import { titleToSlug } from '@/lib/slug'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Dialog, DialogContent, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Collapsible, CollapsibleContent } from '@/components/ui/collapsible'
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import DocSectionTree from '@/components/editor/edit-page/sidebar/DocSectionTree.vue'
import ThemeAccessSettings from '@/components/editor/theme-management/ThemeAccessSettings.vue'
import type { AccessMode } from '@/utils/types'
import {
  X, Plus, Trash2, Rocket, Archive, ArchiveRestore, Star, Copy, ExternalLink,
  Pencil, Check, EyeOff, CircleHelp,
} from 'lucide-vue-next'

const editor = useEditorInject()
const crud = useDocTreeCrud(editor)

// ── 主题表单（从 theme 派生本地副本） ──
const themeForm = ref({ name: '', slug: '', description: '', access_mode: 'public' as AccessMode, access_code: '' })
const themeSlugAuto = ref(true)

const isOpen = computed({
  get: () => editor.showSettingsPanel.value,
  set: (v) => { editor.showSettingsPanel.value = v },
})

watch(isOpen, (open) => {
  if (open && editor.theme.value) {
    const t = editor.theme.value
    const nextForm = {
      name: t.name,
      slug: t.slug,
      description: t.description,
      access_mode: (t.access_mode || 'public') as AccessMode,
      access_code: '',
    }
    themeSlugAuto.value = !nextForm.slug || nextForm.slug === titleToSlug(nextForm.name)
    themeForm.value = nextForm
  }
})

watch(() => themeForm.value.name, () => {
  if (themeSlugAuto.value) {
    themeForm.value.slug = titleToSlug(themeForm.value.name)
  }
})

watch(() => themeForm.value.slug, (slug) => {
  const expectedSlug = titleToSlug(themeForm.value.name)
  themeSlugAuto.value = !slug || slug === expectedSlug
})

// ── 主题 slug 校验 ──
const themeSlugError = computed(() => themeForm.value.slug ? validateSlug(themeForm.value.slug) : '')

// ── 版本内联编辑 ──
const editingVersionId = ref('')
const editVersionForm = ref({ name: '', label: '' })

function startEditVersion(v: Version) {
  editingVersionId.value = v.id
  editVersionForm.value = { name: v.name, label: v.label }
}

function confirmEditVersion() {
  if (editVersionForm.value.name.trim()) {
    editor.updateVersion(editingVersionId.value, { ...editVersionForm.value })
  }
  editingVersionId.value = ''
}

// ── 删除/取消发布确认（版本独有部分） ──
const confirmUnpublishVersion = ref({ open: false, id: '' })
const deleteVersionPopoverId = ref('')

function blurActiveElement() {
  const activeEl = document.activeElement
  if (activeEl instanceof HTMLElement) activeEl.blur()
}

function requestUnpublishVersion(id: string) {
  blurActiveElement()
  confirmUnpublishVersion.value = { open: true, id }
}

function setVersionDeletePopover(id: string, open: boolean) {
  deleteVersionPopoverId.value = open ? id : ''
}

async function confirmDeleteVersionFromPopover(id: string) {
  await editor.deleteVersion(id)
  if (deleteVersionPopoverId.value === id) {
    deleteVersionPopoverId.value = ''
  }
}

// ── 版本名校验 ──
const createVersionNameError = computed(() => versionForm.value.name ? validateVersionName(versionForm.value.name) : '')
const editVersionNameError = computed(() => editVersionForm.value.name ? validateVersionName(editVersionForm.value.name) : '')

// ── 创建内联表单 ──
const showCreateVersion = ref(false)
const versionForm = ref({ name: '', label: '' })

function toggleCreateVersion() {
  showCreateVersion.value = !showCreateVersion.value
  if (!showCreateVersion.value) versionForm.value = { name: '', label: '' }
}

async function submitCreateVersion() {
  await editor.createVersion(versionForm.value)
  showCreateVersion.value = false
  versionForm.value = { name: '', label: '' }
}

</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent :show-close-button="false" class="sm:max-w-[900px] p-0 gap-0 overflow-hidden"
      @pointer-down-outside="(e) => { e.preventDefault(); isOpen = false }">
      <DialogTitle class="sr-only">编辑器设置</DialogTitle>
      <DialogDescription class="sr-only">管理当前主题的基本信息、版本和文档结构</DialogDescription>

      <ResizablePanelGroup direction="horizontal" class="h-[600px] max-h-[600px] overflow-hidden">

        <!-- ══ 左侧：主题设置 + 版本管理 ══ -->
        <ResizablePanel :default-size="46" :min-size="35" class="flex flex-col overflow-hidden [contain:layout_style]">

          <!-- 主题设置区 -->
          <form id="theme-settings-form" class="flex flex-col gap-3 px-6 pt-6 pb-5 border-b border-border/50 shrink-0"
            @submit.prevent="editor.updateTheme(themeForm)">
            <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">主题设置</span>
            <div class="flex gap-2">
              <div class="flex-1 min-w-0">
                <Label class="text-xs text-muted-foreground mb-1 block">名称 <span class="text-destructive">*</span></Label>
                <Input v-model="themeForm.name" placeholder="如 API Reference" class="rounded-xl h-8 text-sm" required />
              </div>
              <div class="w-28 shrink-0">
                <div class="flex items-center gap-1 mb-1">
                  <Label class="text-xs text-muted-foreground">Slug</Label>
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <CircleHelp class="h-3 w-3 text-muted-foreground/50 cursor-help" />
                    </TooltipTrigger>
                    <TooltipContent side="top">{{ slugHint }}</TooltipContent>
                  </Tooltip>
                  <span v-if="themeSlugAuto"
                    class="text-[10px] text-muted-foreground bg-muted/60 rounded-full px-1.5 py-0.5 leading-none font-mono">自动</span>
                </div>
                <Input v-model="themeForm.slug" placeholder="slug" required
                  class="rounded-xl font-mono text-sm h-8"
                  :class="{ 'border-destructive': themeSlugError }" />
              </div>
            </div>
            <div>
              <Label class="text-xs text-muted-foreground mb-1 block">描述</Label>
              <Textarea v-model="themeForm.description" placeholder="主题描述（可留空）"
                class="h-24 resize-none overflow-y-auto rounded-xl text-sm" />
            </div>
            <p v-if="themeSlugError" class="text-xs text-destructive -mt-1">{{ themeSlugError }}</p>

            <!-- 访问权限设置 -->
            <ThemeAccessSettings
              v-model:access-mode="themeForm.access_mode"
              :access-code="themeForm.access_code"
              @update:access-code="themeForm.access_code = $event"
            />
          </form>

          <!-- 版本管理区 -->
          <div class="flex flex-col flex-1 min-h-0 overflow-hidden">
            <div class="flex items-center justify-between px-6 pt-4 pb-2 shrink-0">
              <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                版本管理
                <span class="normal-case font-normal ml-1">({{ editor.versions.value.length }})</span>
              </span>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button size="icon" class="h-6 w-6 rounded-full"
                    :variant="showCreateVersion ? 'default' : 'ghost'" @click="toggleCreateVersion">
                    <Plus class="h-3.5 w-3.5 transition-transform duration-200"
                      :class="{ 'rotate-45': showCreateVersion }" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{{ showCreateVersion ? '取消' : '创建版本' }}</TooltipContent>
              </Tooltip>
            </div>

            <!-- 内联创建版本表单 -->
            <Collapsible :open="showCreateVersion" class="shrink-0 px-6">
              <CollapsibleContent>
                <form class="rounded-2xl border bg-muted/30 px-3 py-2 mb-2 space-y-1.5"
                  @submit.prevent="submitCreateVersion">
                  <div class="flex items-start gap-2">
                    <div class="flex-1 min-w-0">
                      <Input v-model="versionForm.name" placeholder="版本名称，如 v1.0" class="h-8 text-sm"
                        :class="{ 'border-destructive': createVersionNameError }" required />
                      <p v-if="createVersionNameError" class="text-xs text-destructive mt-0.5">{{ createVersionNameError }}</p>
                      <p v-else class="text-xs text-muted-foreground mt-0.5">{{ versionHint }}</p>
                    </div>
                    <Input v-model="versionForm.label" placeholder="展示标签" class="h-8 text-sm w-24 shrink-0" />
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button type="submit" size="icon" class="rounded-full h-7 w-7 shrink-0"
                          :disabled="!!createVersionNameError">
                          <Check class="h-3.5 w-3.5" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>确认创建</TooltipContent>
                    </Tooltip>
                  </div>
                </form>
              </CollapsibleContent>
            </Collapsible>

            <!-- 版本列表 -->
            <div class="flex-1 overflow-y-auto px-4 pb-4">
              <p v-if="editor.versions.value.length === 0"
                class="text-sm text-center text-muted-foreground py-8">暂无版本</p>
              <div v-else class="space-y-0.5">
                <div v-for="v in editor.versions.value" :key="v.id">

                  <template v-if="editingVersionId === v.id">
                    <div class="rounded-xl border bg-muted/30 px-3 py-2 space-y-1.5">
                      <div>
                        <Input v-model="editVersionForm.name" placeholder="版本名称" class="h-8 text-sm"
                          :class="{ 'border-destructive': editVersionNameError }"
                          @keyup.enter="!editVersionNameError && confirmEditVersion()"
                          @keyup.esc="editingVersionId = ''" />
                        <p v-if="editVersionNameError" class="text-xs text-destructive mt-0.5">{{ editVersionNameError }}</p>
                      </div>
                      <div class="flex items-center gap-2">
                        <Input v-model="editVersionForm.label" placeholder="展示标签" class="h-8 text-sm flex-1"
                          @keyup.enter="confirmEditVersion" @keyup.esc="editingVersionId = ''" />
                        <div class="flex gap-1 shrink-0">
                          <Button size="icon" class="rounded-full h-7 w-7" :disabled="!!editVersionNameError"
                            @click="confirmEditVersion">
                            <Check class="h-3.5 w-3.5" />
                          </Button>
                          <Button variant="ghost" size="icon" class="rounded-full h-7 w-7"
                            @click="editingVersionId = ''">
                            <X class="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  </template>

                  <template v-else>
                    <div class="relative flex items-center gap-1.5 group px-3 py-2.5 hover:bg-accent/50">
                      <span class="text-base font-normal truncate flex-1">{{ v.name }}</span>
                      <span v-if="v.label" :class="[
                        'text-[11px] text-muted-foreground truncate transition-opacity',
                        deleteVersionPopoverId === v.id ? 'opacity-0' : 'group-hover:opacity-0',
                      ]">({{ v.label }})</span>
                      <div :class="[
                        'flex items-center gap-1 shrink-0 transition-opacity',
                        deleteVersionPopoverId === v.id ? 'opacity-0' : 'group-hover:opacity-0',
                      ]">
                        <Badge :class="VERSION_STATUS_COLOR[v.status]" class="text-[10px] px-1.5 py-0">
                          {{ VERSION_STATUS_LABEL[v.status] }}
                        </Badge>
                        <Badge v-if="v.is_default" variant="outline" class="text-[10px] px-1.5 py-0">默认</Badge>
                      </div>
                      <div :class="[
                        'absolute right-1 top-1/2 -translate-y-1/2 flex gap-1 opacity-0 transition-opacity group-hover:opacity-100',
                        deleteVersionPopoverId === v.id ? 'opacity-100' : '',
                      ]">
                        <Tooltip v-if="v.status !== 'archived'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="startEditVersion(v)">
                              <Pencil class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>编辑</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'draft'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.publishVersion(v.id)">
                              <Rocket class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>发布</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'published'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="requestUnpublishVersion(v.id)">
                              <EyeOff class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>取消发布</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'published'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.archiveVersion(v.id)">
                              <Archive class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>归档</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'published' && !v.is_default">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.setDefaultVersion(v.id)">
                              <Star class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>设为默认</TooltipContent>
                        </Tooltip>
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.cloneVersion(v.id)">
                              <Copy class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>克隆</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'archived'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.unarchiveVersion(v.id)">
                              <ArchiveRestore class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>取消归档</TooltipContent>
                        </Tooltip>
                        <Tooltip v-if="v.status === 'published'">
                          <TooltipTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                              @click.stop="editor.openVersionReader(v.name)">
                              <ExternalLink class="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>查看文档</TooltipContent>
                        </Tooltip>
                        <Popover v-if="!v.is_default" :open="deleteVersionPopoverId === v.id"
                          @update:open="(open) => setVersionDeletePopover(v.id, open)">
                          <PopoverTrigger as-child>
                            <Button variant="ghost" size="icon" :class="[
                              'rounded-full h-6 w-6 text-destructive hover:text-destructive',
                              deleteVersionPopoverId === v.id ? 'bg-accent' : '',
                            ]" title="删除版本" @click.stop>
                              <Trash2 class="h-3 w-3" />
                            </Button>
                          </PopoverTrigger>
                          <PopoverContent align="end" class="w-64 rounded-2xl p-3">
                            <p class="text-sm font-medium">确认删除该版本？</p>
                            <p class="mt-1 text-xs text-muted-foreground">该版本将被永久删除，此操作无法撤销。</p>
                            <div class="mt-3 flex items-center justify-end gap-2">
                              <Button
                                variant="ghost"
                                size="icon"
                                class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                                @click="deleteVersionPopoverId = ''">
                                <X class="h-3.5 w-3.5" />
                              </Button>
                              <Button size="icon"
                                class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                                @click="confirmDeleteVersionFromPopover(v.id)">
                                <Check class="h-3.5 w-3.5" />
                              </Button>
                            </div>
                          </PopoverContent>
                        </Popover>
                      </div>
                    </div>
                  </template>

                </div>
              </div>
            </div>
          </div>

          <!-- 底部操作栏 -->
          <div class="flex items-center justify-between gap-2 px-6 py-4 border-t border-border/50 shrink-0">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-base font-normal leading-none">编辑器设置</span>
              <span class="text-sm text-muted-foreground leading-tight">管理主题基本信息、版本和文档结构</span>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button type="submit" form="theme-settings-form" size="icon" class="rounded-full"
                    :disabled="!!themeSlugError">
                    <Check class="size-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>保存主题设置</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Button type="button" variant="outline" size="icon" class="rounded-full"
                    @click="isOpen = false">
                    <X class="size-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>关闭</TooltipContent>
              </Tooltip>
            </div>
          </div>

        </ResizablePanel>

        <ResizableHandle class="bg-border/30" />

        <!-- ══ 右侧：章节/文档树 ══ -->
        <ResizablePanel :default-size="54" :min-size="30"
          class="bg-muted/30 flex flex-col overflow-hidden min-h-0 [contain:layout_style]">
          <DocSectionTree :sections="editor.sections.value"
            :active-page-id="editor.activePageId.value" :slug-hint="slugHint"
            :editor="editor" :crud="crud" disable-page-select />
        </ResizablePanel>

      </ResizablePanelGroup>
    </DialogContent>
  </Dialog>

  <!-- 子弹窗全部置于主 Dialog 组件树之外，避免 reka-ui hideOthers 嵌套 aria-hidden 冲突 -->

  <!-- 取消发布版本确认 -->
  <ConfirmDialog v-model:open="confirmUnpublishVersion.open" type="ban" title="确认取消发布"
    description="取消发布后该版本及其所有页面将对读者不可见，可随时重新发布。"
    @confirm="editor.unpublishVersion(confirmUnpublishVersion.id)"
    @cancel="confirmUnpublishVersion.open = false" />

</template>
