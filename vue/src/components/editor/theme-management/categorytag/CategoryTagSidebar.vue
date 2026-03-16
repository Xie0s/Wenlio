<!--
  CategoryTagSidebar.vue - 分类与标签管理侧栏
  职责：在侧栏中提供分类树（内联展开编辑/Popover新建/删除）和标签（# 格式背景色 chip / 侧栏内管理态 Popover 编辑）管理界面
  设计：标签区块标题行提供加号（新建）和 ListFilterPlus（切换侧栏管理态），管理态下点击标签 chip 直接激活 Popover 编辑
  对外暴露：props(catLib: CategoryListLib, tagLib: TagListLib)
-->
<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import type { CategoryListLib, TagListLib } from '@/lib/category-tag'
import type { Tag } from '@/utils/types'
import { titleToSlug } from '@/lib/slug'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import CategorySidebarItem from '@/components/editor/theme-management/categorytag/CategorySidebarItem.vue'
import { Plus, Check, X, Trash2, Loader2, ChevronRight, ListFilterPlus } from 'lucide-vue-next'

const props = defineProps<{
  catLib: CategoryListLib
  tagLib: TagListLib
}>()

// ─── 分类 Popover 新建 ───
const catCreateOpen = ref(false)
const catCreateName = ref('')
const catCreateSlug = ref('')
const catSlugAutoMode = ref(true)
const catNameInputRef = ref<InstanceType<typeof Input> | null>(null)

function onCatNameInput() {
  if (catSlugAutoMode.value) catCreateSlug.value = titleToSlug(catCreateName.value)
}
function onCatSlugInput() {
  catSlugAutoMode.value = catCreateSlug.value === titleToSlug(catCreateName.value)
}
async function submitCatCreate() {
  if (!catCreateName.value) return
  Object.assign(props.catLib.form, {
    name: catCreateName.value,
    slug: catCreateSlug.value || titleToSlug(catCreateName.value),
    description: '',
    parent_id: '',
  })
  const ok = await props.catLib.createCategory()
  if (ok) {
    catCreateName.value = ''
    catCreateSlug.value = ''
    catSlugAutoMode.value = true
    catCreateOpen.value = false
  }
}
function onCatPopoverOpen(open: boolean) {
  catCreateOpen.value = open
  if (!open) {
    catCreateName.value = ''
    catCreateSlug.value = ''
    catSlugAutoMode.value = true
  } else {
    nextTick(() => catNameInputRef.value?.$el?.focus())
  }
}

// ─── 分类管理模式 ───
async function toggleCatManageMode() {
  const next = !props.catLib.manageMode.value
  props.catLib.manageMode.value = next
  catExpanded.value = true
  if (next) {
    await props.catLib.loadCategories()
  } else {
    props.catLib.editId.value = ''
  }
}

// ─── 标签 Popover 新建 ───
const tagCreateOpen = ref(false)
const tagCreateName = ref('')
const tagCreateSlug = ref('')
const tagSlugAutoMode = ref(true)
const tagNameInputRef = ref<InstanceType<typeof Input> | null>(null)

function onTagNameInput() {
  if (tagSlugAutoMode.value) tagCreateSlug.value = titleToSlug(tagCreateName.value)
}
function onTagSlugInput() {
  tagSlugAutoMode.value = tagCreateSlug.value === titleToSlug(tagCreateName.value)
}
async function submitTagCreate() {
  if (!tagCreateName.value) return
  Object.assign(props.tagLib.form, {
    name: tagCreateName.value,
    slug: tagCreateSlug.value || titleToSlug(tagCreateName.value),
    description: '',
  })
  const ok = await props.tagLib.createTag()
  if (ok) {
    tagCreateName.value = ''
    tagCreateSlug.value = ''
    tagSlugAutoMode.value = true
    tagCreateOpen.value = false
  }
}
function onTagPopoverOpen(open: boolean) {
  tagCreateOpen.value = open
  if (!open) {
    tagCreateName.value = ''
    tagCreateSlug.value = ''
    tagSlugAutoMode.value = true
  } else {
    nextTick(() => tagNameInputRef.value?.$el?.focus())
  }
}

// ─── 标签侧栏管理态 ───
const tagManageMode = ref(false)
const tagEditId = ref('')

function resetTagManageState() {
  tagEditId.value = ''
  tagDeleteConfirmId.value = ''
  props.tagLib.editId.value = ''
}

async function toggleTagManageMode() {
  const nextMode = !tagManageMode.value
  tagManageMode.value = nextMode
  tagExpanded.value = true
  if (nextMode) {
    await props.tagLib.loadTags()
  } else {
    resetTagManageState()
  }
}

function toggleTagEditPopover(tag: Tag) {
  if (!tagManageMode.value) return
  if (tagEditId.value === tag.id) {
    resetTagManageState()
  } else {
    tagDeleteConfirmId.value = ''
    tagEditId.value = tag.id
    props.tagLib.startInlineEdit(tag)
  }
}

async function confirmTagEdit() {
  await props.tagLib.updateTag()
  resetTagManageState()
}

function cancelTagEdit() {
  resetTagManageState()
}

const tagDeleteConfirmId = ref('')

async function deleteTagInManage(id: string) {
  await props.tagLib.deleteTag(id, true)
  tagDeleteConfirmId.value = ''
  if (tagEditId.value === id) resetTagManageState()
}

// ─── 选中标签（侧边栏筛选） ───
function isTagSelected(tagId: string) {
  return props.tagLib.selectedTagIds.value.includes(tagId)
}

// ─── 区块折叠状态 ───
const catExpanded = ref(false)
const tagExpanded = ref(false)

// ─── 标签筛选 ───
const tagFilter = ref('')
const filteredTags = computed(() =>
  tagFilter.value.trim()
    ? props.tagLib.tags.value.filter(t => t.name.toLowerCase().includes(tagFilter.value.toLowerCase()))
    : props.tagLib.tags.value
)
</script>

<template>
  <nav class="h-full select-none overflow-visible">

    <!-- ════════ 分类区块 ════════ -->
    <div>

      <!-- 标题行 -->
      <div
        class="group/header flex items-center justify-between px-3 py-3 sticky top-0 z-10 transition-colors cursor-pointer"
        :class="catExpanded ? 'bg-accent' : 'bg-background hover:bg-accent'" @click="catExpanded = !catExpanded">
        <div class="flex items-center gap-2.5">
          <span class="text-xl font-light text-foreground">全部分类</span>
          <span v-if="catLib.categories.value.length > 0" class="text-sm text-muted-foreground tabular-nums">
            {{ catLib.categories.value.length }}
          </span>
        </div>
        <div class="flex items-center gap-1">
          <!-- Popover 新建分类 -->
          <Tooltip>
            <Popover :open="catCreateOpen" @update:open="onCatPopoverOpen">
              <PopoverTrigger as-child>
                <TooltipTrigger as-child>
                  <button
                    class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background"
                    :class="catLib.manageMode.value ? 'opacity-100 bg-primary/10 text-primary' : 'opacity-0 group-hover/header:opacity-100'"
                    @click.stop>
                    <Plus class="h-4 w-4" />
                  </button>
                </TooltipTrigger>
              </PopoverTrigger>
              <PopoverContent side="right" align="start" class="w-72 rounded-2xl p-5">
                <form class="space-y-2.5" @submit.prevent="submitCatCreate">
                  <p class="text-sm font-medium text-muted-foreground/70 mb-2 tracking-wide">新建分类</p>
                  <Input ref="catNameInputRef" v-model="catCreateName" placeholder="分类名称"
                    class="h-8 rounded-full text-sm" @input="onCatNameInput" />
                  <Input v-model="catCreateSlug" placeholder="slug（自动生成）" class="h-8 rounded-full text-sm font-mono"
                    @input="onCatSlugInput" />
                  <div class="flex justify-end gap-1.5 pt-0.5">
                    <Button variant="ghost" size="icon" type="button" class="h-6 w-6 rounded-full"
                      @click="catCreateOpen = false">
                      <X class="h-3 w-3" />
                    </Button>
                    <Button type="submit" size="icon" class="h-6 w-6 rounded-full" :disabled="!catCreateName">
                      <Check class="h-3 w-3" />
                    </Button>
                  </div>
                </form>
              </PopoverContent>
              <TooltipContent side="right">新建分类</TooltipContent>
            </Popover>
          </Tooltip>
          <!-- 管理分类 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <button
                class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background opacity-0 group-hover/header:opacity-100"
                :class="{ '!opacity-100 bg-primary/10 text-primary': catLib.manageMode.value }"
                @click.stop="toggleCatManageMode">
                <ListFilterPlus class="h-4 w-4" />
              </button>
            </TooltipTrigger>
            <TooltipContent side="right">{{ catLib.manageMode.value ? '退出管理' : '管理分类' }}</TooltipContent>
          </Tooltip>
          <!-- 折叠按钮 -->
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background"
            @click.stop="catExpanded = !catExpanded">
            <ChevronRight class="h-4 w-4 transition-transform duration-200" :class="catExpanded ? 'rotate-90' : ''" />
          </button>
        </div>
      </div>

      <!-- 分类内容（展开时，CSS Grid 高度动画） -->
      <div class="grid transition-[grid-template-rows] duration-250 ease-in-out"
        :class="catExpanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'">
        <div class="overflow-hidden">
          <div class="py-1 pb-3 transition-opacity duration-200" :class="catExpanded ? 'opacity-100' : 'opacity-0'">
            <div v-if="catLib.loading.value" class="flex justify-center py-8">
              <Loader2 class="h-3.5 w-3.5 animate-spin text-muted-foreground/40" />
            </div>
            <p v-else-if="catLib.categories.value.length === 0"
              class="pt-5 text-center text-xs text-muted-foreground/35">
              暂无分类
            </p>
            <template v-else>
              <CategorySidebarItem v-for="cat in catLib.categories.value" :key="cat.id" :cat="cat" :cat-lib="catLib" />
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 分隔线 -->
    <div class="flex items-center mx-4 my-2">
      <div class="flex-1 h-px bg-gradient-to-l from-border/70 to-transparent" />
      <div class="w-1 h-1 rounded-full bg-border/60 mx-2" />
      <div class="flex-1 h-px bg-gradient-to-r from-border/70 to-transparent" />
    </div>

    <!-- ════════ 标签区块 ════════ -->
    <div>

      <!-- 标题行 -->
      <div
        class="group/header flex items-center justify-between px-3 py-3 sticky top-0 z-10 transition-colors cursor-pointer"
        :class="tagExpanded ? 'bg-accent' : 'bg-background hover:bg-accent'" @click="tagExpanded = !tagExpanded">
        <div class="flex items-center gap-2.5">
          <span class="text-xl font-light text-foreground">全部标签</span>
          <span v-if="tagLib.tags.value.length > 0" class="text-sm text-muted-foreground tabular-nums">
            {{ tagLib.tags.value.length }}
          </span>
        </div>
        <div class="flex items-center gap-1">
          <!-- Popover 新建标签 -->
          <Tooltip>
            <Popover :open="tagCreateOpen" @update:open="onTagPopoverOpen">
              <PopoverTrigger as-child>
                <TooltipTrigger as-child>
                  <button
                    class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background opacity-0 group-hover/header:opacity-100"
                    :class="{ '!opacity-100 bg-primary/10 text-primary': tagCreateOpen || tagManageMode }" @click.stop>
                    <Plus class="h-4 w-4" />
                  </button>
                </TooltipTrigger>
              </PopoverTrigger>
              <PopoverContent side="right" align="start" class="w-68 rounded-2xl p-5">
                <form class="space-y-2.5" @submit.prevent="submitTagCreate">
                  <p class="text-sm font-medium text-muted-foreground/70 mb-2 tracking-wide">新建标签</p>
                  <Input ref="tagNameInputRef" v-model="tagCreateName" placeholder="标签名称"
                    class="h-8 rounded-full text-sm" @input="onTagNameInput" />
                  <Input v-model="tagCreateSlug" placeholder="slug（自动生成）" class="h-8 rounded-full text-sm font-mono"
                    @input="onTagSlugInput" />
                  <div class="flex justify-end gap-1.5 pt-0.5">
                    <Button variant="ghost" size="icon" type="button" class="h-6 w-6 rounded-full"
                      @click="tagCreateOpen = false">
                      <X class="h-3 w-3" />
                    </Button>
                    <Button type="submit" size="icon" class="h-6 w-6 rounded-full" :disabled="!tagCreateName">
                      <Check class="h-3 w-3" />
                    </Button>
                  </div>
                </form>
              </PopoverContent>
              <TooltipContent side="right">新建标签</TooltipContent>
            </Popover>
          </Tooltip>

          <!-- 管理标签 -->
          <Tooltip>
            <TooltipTrigger as-child>
              <button
                class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background opacity-0 group-hover/header:opacity-100"
                :class="{ '!opacity-100 bg-primary/10 text-primary': tagManageMode }" @click.stop="toggleTagManageMode">
                <ListFilterPlus class="h-4 w-4" />
              </button>
            </TooltipTrigger>
            <TooltipContent side="right">{{ tagManageMode ? '退出管理' : '管理标签' }}</TooltipContent>
          </Tooltip>

          <!-- 折叠按钮 -->
          <button
            class="flex h-7 w-7 items-center justify-center rounded-full transition-all text-muted-foreground hover:text-foreground hover:bg-background"
            @click.stop="tagExpanded = !tagExpanded">
            <ChevronRight class="h-4 w-4 transition-transform duration-200" :class="tagExpanded ? 'rotate-90' : ''" />
          </button>
        </div>
      </div>

      <!-- 标签内容（展开时，CSS Grid 高度动画） -->
      <div class="grid transition-[grid-template-rows] duration-250 ease-in-out"
        :class="tagExpanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'">
        <div class="overflow-hidden">
          <div class="pl-3 pr-3 pb-3 transition-opacity duration-200"
            :class="tagExpanded ? 'opacity-100' : 'opacity-0'">
            <div v-if="tagLib.loading.value" class="flex justify-center py-6">
              <Loader2 class="h-3.5 w-3.5 animate-spin text-muted-foreground/40" />
            </div>
            <p v-else-if="tagLib.tags.value.length === 0" class="pt-5 text-center text-xs text-muted-foreground/35">
              暂无标签
            </p>
            <div v-else class="pt-1.5">
              <!-- 筛选输入框 -->
              <div v-if="tagLib.tags.value.length > 3" class="mb-2.5">
                <input v-model="tagFilter" placeholder="筛选标签…"
                  class="w-full h-7 rounded-full bg-accent/50 px-3 text-xs text-foreground placeholder:text-muted-foreground/50 outline-none focus:bg-accent transition-colors" />
              </div>
              <p v-if="filteredTags.length === 0" class="py-2 text-center text-xs text-muted-foreground/35">无匹配结果</p>
              <!-- 标签云 - 普通态用于筛选，管理态点击打开编辑 Popover -->
              <div class="flex flex-wrap gap-1.5 pt-1">
                <Popover v-for="tag in filteredTags" :key="tag.id" :open="tagManageMode && tagEditId === tag.id"
                  @update:open="(open) => { if (!open && tagEditId === tag.id) cancelTagEdit() }">
                  <PopoverTrigger as-child>
                    <button
                      class="inline-flex items-center gap-0.5 h-7 px-2.5 rounded-full text-sm font-normal transition-all"
                      :class="tagManageMode
                        ? (tagEditId === tag.id
                          ? 'bg-primary/10 text-primary'
                          : 'bg-accent/70 text-foreground hover:bg-accent')
                        : (isTagSelected(tag.id)
                          ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
                          : 'bg-accent/70 text-foreground hover:bg-accent')"
                      @click.stop="tagManageMode ? toggleTagEditPopover(tag) : tagLib.toggleTag(tag.id)">
                      <span class="opacity-40">#</span>{{ tag.name }}
                    </button>
                  </PopoverTrigger>
                  <PopoverContent v-if="tagManageMode" align="start" class="w-64 rounded-2xl p-4">
                    <form class="space-y-2" @submit.prevent="confirmTagEdit">
                      <Input v-model="tagLib.editForm.name" placeholder="标签名称" class="h-8 rounded-full text-sm" />
                      <Input v-model="tagLib.editForm.slug" placeholder="slug（可选）"
                        class="h-8 rounded-full text-sm font-mono text-muted-foreground" />
                      <div v-if="tagDeleteConfirmId === tag.id" class="pt-1 space-y-2">
                        <p v-if="tag.usage_count > 0" class="text-xs text-muted-foreground/70">
                          已被 {{ tag.usage_count }} 个主题使用，删除后将自动移除关联。
                        </p>
                        <div class="flex items-center justify-between">
                          <span class="text-xs text-destructive font-medium">确认删除？</span>
                          <div class="flex gap-1">
                            <Button variant="ghost" size="icon" type="button"
                              class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                              @click.stop="tagDeleteConfirmId = ''">
                              <X class="h-3.5 w-3.5" />
                            </Button>
                            <Button size="icon" type="button"
                              class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                              @click.stop="deleteTagInManage(tag.id)">
                              <Check class="h-3.5 w-3.5" />
                            </Button>
                          </div>
                        </div>
                      </div>
                      <div v-else class="flex items-center justify-between pt-1">
                        <Tooltip>
                          <TooltipTrigger as-child>
                            <button type="button"
                              class="flex h-7 w-7 items-center justify-center rounded-full text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-colors"
                              @click.stop="tagDeleteConfirmId = tag.id">
                              <Trash2 class="h-3.5 w-3.5" />
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>删除标签</TooltipContent>
                        </Tooltip>
                        <div class="flex gap-1">
                          <Button variant="ghost" size="icon" type="button" class="h-7 w-7 rounded-full"
                            @click="cancelTagEdit">
                            <X class="h-3.5 w-3.5" />
                          </Button>
                          <Button type="submit" size="icon" class="h-7 w-7 rounded-full"
                            :disabled="!tagLib.editForm.name">
                            <Check class="h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </div>
                    </form>
                  </PopoverContent>
                </Popover>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

  </nav>
</template>