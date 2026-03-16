<!--
  DocSectionTree.vue - 文档章节树共享组件
  职责：集成章节树渲染、展开折叠、拖拽排序、顶部操作栏及创建章节表单；承载章节/页面 CRUD 交互；不负责业务状态创建与排序规则计算。
  对外暴露：props(sections, activePageId, slugHint, editor, crud, disablePageSelect)
-->
<script setup lang="ts">
import type { ThemeEditorContext } from '@/composables/useThemeEditor'
import type { DocTreeCrudContext } from '@/composables/useDocTreeCrud'
import type { SectionWithPages } from '@/components/editor/edit-page/types'
import { useDocSectionDrag } from '@/composables/useDocSectionDrag'
import { Collapsible, CollapsibleContent } from '@/components/ui/collapsible'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import SectionInlineEditor from '@/components/editor/edit-page/sidebar/SectionInlineEditor.vue'
import DocPageInlineEditor from '@/components/editor/edit-page/sidebar/DocPageInlineEditor.vue'
import { ChevronRight, Plus, Trash2, Pencil, Check, X } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  sections: SectionWithPages[]
  activePageId: string
  slugHint: string
  editor: ThemeEditorContext
  crud: DocTreeCrudContext
  disablePageSelect?: boolean
}>(), {
  disablePageSelect: false,
})

const {
  dragIndex,
  isSectionOpen,
  toggleSection,
  showInsertionLineAt,
  onDragStart,
  onDragOver,
  onDragEnd,
  onDrop,
} = useDocSectionDrag(props.editor)
</script>

<template>
  <div class="h-full flex flex-col overflow-hidden">

    <!-- 顶部操作栏 -->
    <div class="relative flex items-center justify-center shrink-0 px-3 pt-3 pb-1.5">
      <span class="text-xs text-muted-foreground">{{ sections.length }} 个章节</span>
      <Tooltip>
        <Popover :open="crud.showCreateSection.value"
          @update:open="(open: boolean) => { if (!open) crud.showCreateSection.value = false }">
          <PopoverTrigger as-child>
            <TooltipTrigger as-child>
              <Button size="icon" class="absolute right-3 h-7 w-7 rounded-full"
                :variant="crud.showCreateSection.value ? 'default' : 'outline'"
                @click="crud.toggleCreateSection">
                <Plus class="h-3.5 w-3.5 transition-transform duration-200"
                  :class="{ 'rotate-45': crud.showCreateSection.value }" />
              </Button>
            </TooltipTrigger>
          </PopoverTrigger>
          <PopoverContent align="end" :side-offset="8" class="w-64 rounded-2xl p-2">
            <form class="flex items-center gap-1.5" @submit.prevent="crud.submitCreateSection">
              <Input v-model="crud.sectionTitle.value" placeholder="章节标题"
                class="h-9 rounded-xl text-sm flex-1 min-w-0 border-none bg-transparent focus-visible:ring-0 shadow-none" autofocus required />
              <Button variant="ghost" size="icon" type="button" class="h-8 w-8 shrink-0 rounded-full"
                @click="crud.showCreateSection.value = false">
                <X class="h-3.5 w-3.5" />
              </Button>
              <Button type="submit" size="icon" class="h-8 w-8 shrink-0 rounded-full">
                <Check class="h-3.5 w-3.5" />
              </Button>
            </form>
          </PopoverContent>
        </Popover>
        <TooltipContent>{{ crud.showCreateSection.value ? '取消' : '添加章节' }}</TooltipContent>
      </Tooltip>
    </div>

    <!-- 章节树列表区域 -->
    <div class="flex-1 overflow-y-auto px-2 pb-3">
      <p v-if="sections.length === 0" class="text-center text-sm text-muted-foreground py-6">
        暂无章节，点击右上角 + 添加
      </p>

      <div v-else class="space-y-0.5">
        <div v-for="(section, index) in sections" :key="section.id" class="relative">
          <div v-if="showInsertionLineAt(index)"
            class="pointer-events-none absolute inset-x-2 top-0 z-10 h-0.5 -translate-y-1/2 rounded-full bg-primary" />
          <div draggable="true" class="transition-transform duration-100 select-none"
            :class="{ 'opacity-40': dragIndex === index }"
            @dragstart="onDragStart(index, $event)"
            @dragover="onDragOver(index, $event)"
            @dragend="onDragEnd()"
            @drop="onDrop(index, $event)">
            <Collapsible :open="isSectionOpen(section.id)" @update:open="toggleSection(section.id)">
              <div class="relative flex items-center gap-1.5 group px-3 py-2.5 hover:bg-accent/50"
                :class="{ 'cursor-pointer': crud.editingSectionId.value !== section.id }"
                @click="crud.editingSectionId.value !== section.id && toggleSection(section.id)">

                <template v-if="crud.editingSectionId.value === section.id">
                  <SectionInlineEditor v-model="crud.editSectionTitle.value" @confirm="crud.confirmEditSection"
                    @cancel="crud.editingSectionId.value = ''" />
                </template>

                <template v-else>
                  <span class="truncate flex-1 text-lg font-thin text-foreground/90">{{ section.title }}</span>
                  <div :class="[
                    'absolute right-1 top-1/2 -translate-y-1/2 flex gap-1 opacity-0 transition-opacity group-hover:opacity-100',
                    crud.createPagePopoverSectionId.value === section.id || crud.deleteSectionPopoverId.value === section.id ? 'opacity-100' : '',
                  ]">
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                          @click.stop="toggleSection(section.id)">
                          <ChevronRight class="h-3 w-3 transition-transform duration-200"
                            :class="{ 'rotate-90': isSectionOpen(section.id) }" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{{ isSectionOpen(section.id) ? '折叠' : '展开' }}</TooltipContent>
                    </Tooltip>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                          @click.stop="crud.startEditSection(section)">
                          <Pencil class="h-3 w-3" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>编辑章节</TooltipContent>
                    </Tooltip>
                    <Popover :open="crud.createPagePopoverSectionId.value === section.id"
                      @update:open="(open: boolean) => crud.setCreatePagePopover(section.id, open)">
                      <PopoverTrigger as-child>
                        <Button variant="ghost" size="icon" :class="[
                          'rounded-full h-6 w-6',
                          crud.createPagePopoverSectionId.value === section.id ? 'bg-background text-foreground shadow-sm' : '',
                        ]" title="添加文档页" @click.stop>
                          <Plus class="h-3 w-3" />
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent align="end" class="w-80 rounded-2xl p-3">
                        <form class="space-y-2" @submit.prevent="crud.submitCreatePage">
                          <Input v-model="crud.pageForm.value.title" placeholder="标题，如 快速开始"
                            class="h-8 rounded-full text-sm" required />
                          <div>
                            <Input v-model="crud.pageForm.value.slug" placeholder="slug，如 kuaisukaishi" required
                              class="h-8 rounded-full text-sm font-mono"
                              :class="{ 'border-destructive': crud.createPageSlugError.value }" />
                            <p v-if="crud.createPageSlugError.value" class="text-xs text-destructive mt-0.5">{{
                              crud.createPageSlugError.value }}</p>
                            <p v-else class="text-xs text-muted-foreground mt-0.5">{{ slugHint }}</p>
                          </div>
                          <div class="flex justify-end">
                            <Button type="submit" size="icon" class="h-7 w-7 rounded-full"
                              :disabled="!!crud.createPageSlugError.value">
                              <Check class="h-3.5 w-3.5" />
                            </Button>
                          </div>
                        </form>
                      </PopoverContent>
                    </Popover>
                    <Popover :open="crud.deleteSectionPopoverId.value === section.id"
                      @update:open="(open: boolean) => crud.setSectionDeletePopover(section.id, open)">
                      <PopoverTrigger as-child>
                        <Button variant="ghost" size="icon" :class="[
                          'rounded-full h-6 w-6 text-destructive hover:text-destructive',
                          crud.deleteSectionPopoverId.value === section.id ? 'bg-background shadow-sm' : '',
                        ]" title="删除章节" @click.stop>
                          <Trash2 class="h-3 w-3" />
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent align="end" class="w-64 rounded-2xl p-3">
                        <p class="text-sm font-medium">确认删除该章节？</p>
                        <p class="mt-1 text-xs text-muted-foreground">章节及其下所有文档页将被永久删除。</p>
                        <div class="mt-3 flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                            @click="crud.deleteSectionPopoverId.value = ''">
                            <X class="h-3.5 w-3.5" />
                          </Button>
                          <Button size="icon"
                            class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                            @click="crud.confirmDeleteSection(section.id)">
                            <Check class="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </PopoverContent>
                    </Popover>
                  </div>
                </template>
              </div>

              <CollapsibleContent>
                <div class="relative pl-6 py-0.5">
                  <div class="pointer-events-none absolute top-0 bottom-0 left-[15px] w-px bg-border/50" />
                  <template v-for="page in section.pages" :key="page.id">
                    <div class="relative group/page -ml-[9px] transition-colors duration-150"
                      :class="activePageId === page.id ? 'bg-accent/40' : 'hover:bg-accent/60'">
                      <span v-if="activePageId === page.id"
                        class="pointer-events-none absolute left-0 inset-y-0 w-px bg-primary" />

                      <template v-if="crud.editingPageId.value === page.id">
                        <div class="pl-[16px] pr-2 py-2">
                          <DocPageInlineEditor v-model:title="crud.editPageForm.value.title"
                            v-model:slug="crud.editPageForm.value.slug" :slug-error="crud.editPageSlugError.value"
                            :slug-hint="slugHint" @confirm="crud.confirmEditPage" @cancel="crud.editingPageId.value = ''" />
                        </div>
                      </template>

                      <template v-else>
                        <button class="w-full text-left py-1.5 pl-[16px] pr-8 text-base font-light transition-colors duration-150 truncate"
                          :class="activePageId === page.id ? 'text-primary' : 'text-foreground/80 hover:text-primary/80'"
                          @click="!disablePageSelect && editor.selectPage(page.id)">
                          {{ page.title }}
                        </button>

                        <div :class="[
                          'absolute right-1 top-1/2 -translate-y-1/2 flex gap-1 opacity-0 group-hover/page:opacity-100 transition-opacity',
                          crud.deletePagePopoverId.value === page.id ? 'opacity-100' : '',
                        ]">
                          <Tooltip>
                            <TooltipTrigger as-child>
                              <Button variant="ghost" size="icon" class="rounded-full h-6 w-6"
                                @click.stop="crud.startEditPage(page)">
                                <Pencil class="h-3 w-3" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>编辑</TooltipContent>
                          </Tooltip>
                          <Popover :open="crud.deletePagePopoverId.value === page.id"
                            @update:open="(open: boolean) => crud.setPageDeletePopover(page.id, open)">
                            <PopoverTrigger as-child>
                              <Button variant="ghost" size="icon" :class="[
                                'rounded-full h-6 w-6 text-destructive hover:text-destructive',
                                crud.deletePagePopoverId.value === page.id ? 'bg-background shadow-sm' : '',
                              ]" title="删除文档页" @click.stop>
                                <Trash2 class="h-3 w-3" />
                              </Button>
                            </PopoverTrigger>
                            <PopoverContent align="end" class="w-64 rounded-2xl p-3">
                              <p class="text-sm font-medium">确认删除该文档页？</p>
                              <p class="mt-1 text-xs text-muted-foreground">该文档页将被永久删除。</p>
                              <div class="mt-3 flex items-center justify-end gap-2">
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                                  @click="crud.deletePagePopoverId.value = ''">
                                  <X class="h-3.5 w-3.5" />
                                </Button>
                                <Button size="icon"
                                  class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                                  @click="crud.confirmDeletePage(page.id)">
                                  <Check class="h-3.5 w-3.5" />
                                </Button>
                              </div>
                            </PopoverContent>
                          </Popover>
                        </div>
                      </template>
                    </div>
                  </template>
                  <p v-if="section.pages.length === 0" class="text-xs text-center text-muted-foreground py-2">暂无文档页</p>
                </div>
              </CollapsibleContent>
            </Collapsible>
          </div>
        </div>
        <div v-if="showInsertionLineAt(sections.length)" class="pointer-events-none relative mx-2 h-0">
          <div class="absolute inset-x-0 top-0 h-0.5 -translate-y-1/2 rounded-full bg-primary" />
        </div>
      </div>
    </div>

  </div>
</template>
