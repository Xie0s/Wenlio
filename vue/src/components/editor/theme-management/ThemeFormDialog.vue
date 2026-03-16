<!-- ThemeFormDialog.vue - 主题创建/编辑表单 Dialog（左右分栏设计）
     职责：内部管理表单状态（含 slug 自动生成联动），右侧支持分类/标签 CRUD，emit submit 数据给父层
     对外暴露：
       Props: open(boolean), mode('create'|'edit'), initialValues, catLib, tagLib, submitting(boolean)
       Emits: update:open(boolean), submit(ThemeFormData) -->
<script setup lang="ts">
import { ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '@/components/ui/dialog'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable'
import type { CategoryListLib, TagListLib } from '@/lib/category-tag'
import { titleToSlug } from '@/lib/slug'
import { Check, X, Pencil, Plus } from 'lucide-vue-next'
import CategorySidebarItem from '@/components/editor/theme-management/categorytag/CategorySidebarItem.vue'
import ThemeAccessSettings from '@/components/editor/theme-management/ThemeAccessSettings.vue'
import type { AccessMode } from '@/utils/types'

export interface ThemeFormData {
  name: string
  slug: string
  description: string
  category_id?: string
  tag_ids?: string[]
  access_mode?: AccessMode
  access_code?: string
}

const props = defineProps<{
  open: boolean
  mode: 'create' | 'edit'
  initialValues?: Partial<ThemeFormData>
  catLib: CategoryListLib
  tagLib: TagListLib
  submitting?: boolean
}>()

const emit = defineEmits<{
  'update:open': [boolean]
  submit: [data: ThemeFormData]
}>()

const form = ref<ThemeFormData>({
  name: '',
  slug: '',
  description: '',
  category_id: undefined,
  tag_ids: [],
  access_mode: 'public',
  access_code: '',
})
const slugAuto = ref(true)

// ─── 分类新建 Popover ───
const catCreateOpen = ref(false)
const catCreateName = ref('')
const catCreateSlug = ref('')
const catSlugAutoMode = ref(true)

// ─── 标签新建 Popover ───
const tagCreateOpen = ref(false)
const tagCreateName = ref('')
const tagCreateSlug = ref('')
const tagSlugAutoMode = ref(true)

// ─── 标签删除确认 Popover ───
const tagDeleteId = ref('')

// 打开时重置并填充初始值
watch(() => props.open, (open) => {
  if (!open) return
  const iv = props.initialValues ?? {}
  form.value = {
    name: iv.name ?? '',
    slug: iv.slug ?? '',
    description: iv.description ?? '',
    category_id: iv.category_id ?? undefined,
    tag_ids: iv.tag_ids ? [...iv.tag_ids] : [],
    access_mode: iv.access_mode ?? 'public',
    access_code: iv.access_code ?? '',
  }
  slugAuto.value = props.mode === 'create'
    ? (!form.value.slug || form.value.slug === titleToSlug(form.value.name))
    : false
})

watch(() => form.value.name, () => {
  if (slugAuto.value) form.value.slug = titleToSlug(form.value.name)
})

watch(() => form.value.slug, (slug) => {
  const expected = titleToSlug(form.value.name)
  slugAuto.value = !slug || slug === expected
})

// ── 分类选择 ──
function selectCategory(catId: string) {
  form.value.category_id = form.value.category_id === catId ? undefined : catId
}

// ── 分类新建 ──
function onCatCreatePopoverChange(open: boolean) {
  catCreateOpen.value = open
  if (!open) { catCreateName.value = ''; catCreateSlug.value = ''; catSlugAutoMode.value = true }
}
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
  if (ok) catCreateOpen.value = false
}

// ── 标签新建 ──
function onTagCreatePopoverChange(open: boolean) {
  tagCreateOpen.value = open
  if (!open) { tagCreateName.value = ''; tagCreateSlug.value = ''; tagSlugAutoMode.value = true }
}
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
  if (ok) tagCreateOpen.value = false
}

// ── 标签选择切换 ──
function toggleTag(tagId: string) {
  if (!form.value.tag_ids) form.value.tag_ids = []
  const idx = form.value.tag_ids.indexOf(tagId)
  if (idx >= 0) form.value.tag_ids.splice(idx, 1)
  else form.value.tag_ids.push(tagId)
}

function handleSubmit() {
  emit('submit', { ...form.value })
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent :show-close-button="false" class="sm:max-w-[880px] p-0 gap-0 overflow-hidden"
      @pointer-down-outside="(e) => { e.preventDefault(); emit('update:open', false) }">
      <DialogTitle class="sr-only">{{ mode === 'create' ? '创建主题' : '编辑主题' }}</DialogTitle>
      <DialogDescription class="sr-only">{{ mode === 'create' ? '填写基本信息，在右侧选择分类。' : '修改主题信息、分类与标签。' }}
      </DialogDescription>
      <form class="h-[560px]" @submit.prevent="handleSubmit">
        <ResizablePanelGroup direction="horizontal" class="h-full">

          <!-- ══ 左侧：信息填写区 ══ -->
          <ResizablePanel :default-size="60" :min-size="40"
            class="flex flex-col gap-5 p-7 min-w-0 overflow-hidden [contain:layout_style]">
            <div class="flex gap-3">
              <div class="flex-1 flex flex-col gap-1.5">
                <Label class="text-sm font-medium">主题名称 <span class="text-destructive">*</span></Label>
                <Input v-model="form.name" placeholder="如 API Reference" class="rounded-xl" />
              </div>
              <div class="w-36 flex flex-col gap-1.5">
                <div class="flex items-center gap-1.5">
                  <Label class="text-sm font-medium">Slug</Label>
                  <span v-if="slugAuto"
                    class="rounded-full bg-muted/60 px-1.5 py-0.5 text-[10px] text-muted-foreground font-mono leading-none">自动</span>
                </div>
                <Input v-model="form.slug" placeholder="slug" class="rounded-xl font-mono text-sm" />
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <Label class="text-sm font-medium">描述</Label>
              <Textarea v-model="form.description" placeholder="主题描述（可留空）"
                class="rounded-xl resize-none h-20 overflow-y-auto" />
            </div>

            <!-- ── 访问权限区 ── -->
            <ThemeAccessSettings
              :access-mode="form.access_mode ?? 'public'"
              :access-code="form.access_code ?? ''"
              @update:access-mode="form.access_mode = $event"
              @update:access-code="form.access_code = $event"
            />

            <!-- ── 标签区 ── -->
            <div class="flex flex-col gap-1.5 flex-1 min-h-0">
              <div class="flex items-center justify-between">
                <Label class="text-sm font-medium">标签</Label>
                <Tooltip>
                  <Popover :open="tagCreateOpen" @update:open="onTagCreatePopoverChange">
                    <PopoverTrigger as-child>
                      <TooltipTrigger as-child>
                        <button type="button"
                          class="h-5 w-5 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground hover:bg-accent/80 transition-colors">
                          <Plus class="h-3.5 w-3.5" />
                        </button>
                      </TooltipTrigger>
                    </PopoverTrigger>
                    <PopoverContent align="end" class="w-60 rounded-2xl p-4">
                      <p class="text-xs font-medium text-muted-foreground mb-2.5">新建标签</p>
                      <div class="space-y-2">
                        <Input v-model="tagCreateName" placeholder="标签名称" class="h-8 rounded-full text-sm"
                          @input="onTagNameInput" @keydown.enter.prevent="submitTagCreate" />
                        <Input v-model="tagCreateSlug" placeholder="slug（自动生成）"
                          class="h-8 rounded-full text-sm font-mono" @input="onTagSlugInput" />
                      </div>
                      <div class="flex justify-end gap-1.5 mt-3">
                        <Button variant="ghost" size="icon" type="button" class="h-6 w-6 rounded-full"
                          @click="tagCreateOpen = false">
                          <X class="h-3 w-3" />
                        </Button>
                        <Button type="button" size="icon" class="h-6 w-6 rounded-full" :disabled="!tagCreateName"
                          @click="submitTagCreate">
                          <Check class="h-3 w-3" />
                        </Button>
                      </div>
                    </PopoverContent>
                  </Popover>
                  <TooltipContent>新建标签</TooltipContent>
                </Tooltip>
              </div>
              <div class="flex flex-wrap gap-1.5 content-start flex-1 overflow-y-auto">
                <p v-if="tagLib.tags.value.length === 0" class="text-xs text-muted-foreground/50 self-center">暂无标签</p>
                <template v-else>
                  <div v-for="tag in tagLib.tags.value" :key="tag.id" class="group/tag-chip">
                    <div v-if="tagLib.editId.value === tag.id"
                      class="inline-flex items-center gap-1 rounded-full border border-primary/40 bg-primary/5 px-2.5 py-0.5">
                      <input v-model="tagLib.editForm.name"
                        class="w-20 min-w-0 bg-transparent text-sm outline-none text-foreground"
                        @keydown.enter.prevent="tagLib.updateTag()" @keydown.escape="tagLib.editId.value = ''" />
                      <button type="button"
                        class="p-0.5 rounded-full text-primary hover:bg-primary/10 transition-colors"
                        @click="tagLib.updateTag()">
                        <Check class="h-3 w-3" />
                      </button>
                      <button type="button"
                        class="p-0.5 rounded-full text-muted-foreground hover:bg-accent transition-colors"
                        @click="tagLib.editId.value = ''">
                        <X class="h-3 w-3" />
                      </button>
                    </div>
                    <div v-else class="inline-flex items-center rounded-full border text-sm transition-colors"
                      :class="form.tag_ids?.includes(tag.id) ? 'bg-primary/10 border-primary/50 text-primary' : 'border-border/60 text-muted-foreground'">
                      <button type="button" class="px-2.5 py-0.5 transition-colors"
                        :class="!form.tag_ids?.includes(tag.id) ? 'hover:text-foreground' : ''"
                        @click="toggleTag(tag.id)">
                        {{ tag.name }}
                      </button>
                      <span
                        class="inline-flex items-center gap-0.5 overflow-hidden max-w-0 opacity-0 pr-0 group-hover/tag-chip:max-w-[44px] group-hover/tag-chip:opacity-100 group-hover/tag-chip:pr-1 transition-all duration-200 ease-out">
                        <span class="h-3 w-px flex-none bg-border/50" />
                        <button type="button"
                          class="h-4 w-4 flex items-center justify-center rounded-full hover:bg-accent transition-colors"
                          @click.stop="tagLib.startInlineEdit(tag)">
                          <Pencil class="h-2.5 w-2.5" />
                        </button>
                        <Popover :open="tagDeleteId === tag.id" @update:open="(v) => (tagDeleteId = v ? tag.id : '')">
                          <PopoverTrigger as-child>
                            <button type="button"
                              class="h-4 w-4 flex items-center justify-center rounded-full hover:bg-destructive/10 hover:text-destructive transition-colors"
                              @click.stop>
                              <X class="h-2.5 w-2.5" />
                            </button>
                          </PopoverTrigger>
                          <PopoverContent align="end" class="w-52 rounded-2xl p-4">
                            <p class="text-sm font-medium">删除「{{ tag.name }}」？</p>
                            <p v-if="tag.usage_count > 0" class="mt-1 text-xs text-muted-foreground/70">
                              已被 {{ tag.usage_count }} 个主题使用，删除后将自动移除关联。
                            </p>
                            <div class="mt-3 flex justify-end gap-1.5">
                              <Button variant="ghost" size="icon"
                                class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                                @click="tagDeleteId = ''">
                                <X class="h-3.5 w-3.5" />
                              </Button>
                              <Button size="icon"
                                class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                                @click="tagLib.deleteTag(tag.id, true); tagDeleteId = ''">
                                <Check class="h-3.5 w-3.5" />
                              </Button>
                            </div>
                          </PopoverContent>
                        </Popover>
                      </span>
                    </div>
                  </div>
                </template>
              </div>
            </div>

            <div class="flex items-center justify-between gap-2 pt-3 border-t border-border/50">
              <div class="flex flex-col gap-0.5 min-w-0">
                <span class="text-base font-semibold leading-none">
                  {{ mode === 'create' ? '创建主题' : '编辑主题' }}
                </span>
                <span class="text-sm text-muted-foreground leading-tight">
                  {{ mode === 'create' ? '填写基本信息，在右侧选择分类和标签。' : '修改主题信息、分类与标签。' }}
                </span>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button type="button" variant="outline" size="icon" class="rounded-full" :disabled="submitting"
                      @click="emit('update:open', false)">
                      <X class="size-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>取消</TooltipContent>
                </Tooltip>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button type="submit" size="icon" class="rounded-full"
                      :disabled="submitting || !form.name || !form.category_id">
                      <Check class="size-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>{{ mode === 'create' ? '创建' : '保存' }}</TooltipContent>
                </Tooltip>
              </div>
            </div>
          </ResizablePanel>

          <ResizableHandle class="bg-border/30" />

          <!-- ══ 右侧：分类标签编辑选择区 ══ -->
          <ResizablePanel :default-size="40" :min-size="25"
            class="bg-muted/30 flex flex-col overflow-hidden [contain:layout_style]">

            <!-- ── 分类区 ── -->
            <div class="flex flex-col flex-1 min-h-0">
              <!-- 标题行 -->
              <div class="px-4 pt-5 pb-2 shrink-0 flex items-center justify-between">
                <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  分类 <span class="text-destructive">*</span>
                </span>
                <Tooltip>
                  <Popover :open="catCreateOpen" @update:open="onCatCreatePopoverChange">
                    <PopoverTrigger as-child>
                      <TooltipTrigger as-child>
                        <button type="button"
                          class="h-6 w-6 flex items-center justify-center rounded-full text-muted-foreground hover:text-foreground hover:bg-accent/80 transition-colors">
                          <Plus class="h-4 w-4" />
                        </button>
                      </TooltipTrigger>
                    </PopoverTrigger>
                    <PopoverContent align="end" class="w-60 rounded-2xl p-4">
                      <p class="text-xs font-medium text-muted-foreground mb-2.5">新建分类</p>
                      <div class="space-y-2">
                        <Input v-model="catCreateName" placeholder="分类名称" class="h-8 rounded-full text-sm"
                          @input="onCatNameInput" @keydown.enter.prevent="submitCatCreate" />
                        <Input v-model="catCreateSlug" placeholder="slug（自动生成）"
                          class="h-8 rounded-full text-sm font-mono" @input="onCatSlugInput" />
                      </div>
                      <div class="flex justify-end gap-1.5 mt-3">
                        <Button variant="ghost" size="icon" type="button" class="h-6 w-6 rounded-full"
                          @click="catCreateOpen = false">
                          <X class="h-3 w-3" />
                        </Button>
                        <Button type="button" size="icon" class="h-6 w-6 rounded-full" :disabled="!catCreateName"
                          @click="submitCatCreate">
                          <Check class="h-3 w-3" />
                        </Button>
                      </div>
                    </PopoverContent>
                  </Popover>
                  <TooltipContent>新建分类</TooltipContent>
                </Tooltip>
              </div>

              <!-- 分类列表（复用 CategorySidebarItem，通过 selectedId/onSelect 接入表单状态） -->
              <div class="flex-1 overflow-y-auto pb-2">
                <p v-if="catLib.categories.value.length === 0"
                  class="text-xs text-muted-foreground/50 px-3 py-6 text-center">暂无分类
                </p>
                <template v-else>
                  <CategorySidebarItem v-for="cat in catLib.categories.value" :key="cat.id" :cat="cat" :cat-lib="catLib"
                    :selected-id="form.category_id ?? ''" :on-select="selectCategory" />
                </template>
              </div>
            </div>

          </ResizablePanel>
        </ResizablePanelGroup>
      </form>
    </DialogContent>
  </Dialog>
</template>
