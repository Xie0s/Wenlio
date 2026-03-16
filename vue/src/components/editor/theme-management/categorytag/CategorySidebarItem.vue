<!--
  CategorySidebarItem.vue - 分类侧栏项（递归，重设计版）
  职责：渲染单个分类节点，支持展开/折叠子分类、内联展开编辑、Popover新建子分类、Popover删除确认
  设计：EfficientSidebar 极简风格，容器递进缩进 + VS Code 竖线，hover 浮出操作组
  对外暴露：props(cat: Category, catLib: CategoryListLib)
-->
<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import type { Category } from '@/utils/types'
import type { CategoryListLib } from '@/lib/category-tag'
import { titleToSlug } from '@/lib/slug'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { ChevronRight, Plus, Pencil, Trash2, Check, X } from 'lucide-vue-next'

const props = defineProps<{
  cat: Category
  catLib: CategoryListLib
  selectedId?: string
  onSelect?: (id: string) => void
}>()

const hasChildren = computed(() => (props.cat.children?.length ?? 0) > 0)
const isOpen = ref(false)
const isEditing = computed(() => props.catLib.editId.value === props.cat.id)
const isMainCategory = computed(() => (props.cat.level ?? 0) === 0)

const effectiveSelectedId = computed(() =>
  props.selectedId !== undefined ? props.selectedId : props.catLib.selectedCategoryId.value
)
function doSelect(id: string) {
  if (props.onSelect) props.onSelect(id)
  else props.catLib.selectCategory(id)
}

const totalCount = computed(() => {
  const own = props.cat.theme_count ?? 0
  if (isMainCategory.value && props.cat.children?.length) {
    return own + props.cat.children.reduce((sum, c) => sum + (c.theme_count ?? 0), 0)
  }
  return own
})

const addChildOpen = ref(false)
const addChildName = ref('')
const addChildSlug = ref('')
const deleteOpen = ref(false)
const deleting = ref(false)

const nameEl = ref<HTMLElement | null>(null)
const isNameTruncated = ref(false)

function checkTruncation() {
  const el = nameEl.value
  if (el) isNameTruncated.value = el.scrollWidth > el.clientWidth
}

onMounted(checkTruncation)
watch(() => props.cat.name, () => nextTick(checkTruncation))

function handleCatClick(id: string) {
  doSelect(id)
  if (isMainCategory.value && hasChildren.value) {
    isOpen.value = true
  }
}

function startEdit() { props.catLib.startInlineEdit(props.cat) }
function cancelEdit() { props.catLib.editId.value = '' }
async function confirmEdit() { await props.catLib.updateCategory() }

function onAddChildNameInput() { addChildSlug.value = titleToSlug(addChildName.value) }
async function confirmAddChild() {
  if (!addChildName.value) return
  const ok = await props.catLib.createChildCategory(props.cat.id, {
    name: addChildName.value,
    slug: addChildSlug.value || titleToSlug(addChildName.value),
  })
  if (ok) {
    addChildName.value = ''
    addChildSlug.value = ''
    addChildOpen.value = false
    isOpen.value = true
  }
}
async function confirmDelete() {
  deleting.value = true
  await props.catLib.deleteCategory(props.cat.id)
  deleting.value = false
  deleteOpen.value = false
}
</script>

<template>
  <div>
    <!-- 分类行 -->
    <div class="group/cat relative flex items-center gap-1.5 px-3 py-2.5 transition-colors cursor-pointer hover:z-10"
      :class="isMainCategory ? (effectiveSelectedId === cat.id ? 'bg-accent/40' : (isOpen && hasChildren ? 'bg-accent/50' : '')) : ''"
      @click="!isEditing && hasChildren && (isOpen = !isOpen)">
      <!-- 内联编辑模式 -->
      <template v-if="isEditing">
        <textarea v-model="catLib.editForm.name" rows="1"
          class="flex-1 min-w-0 bg-transparent text-sm outline-none text-foreground border-b border-primary resize-none field-sizing-content"
          @keydown.escape="cancelEdit" @click.stop />
        <button type="button" class="flex-none text-primary hover:text-primary/80 transition-colors"
          @click.stop="confirmEdit">
          <Check class="h-5 w-5" />
        </button>
        <button type="button" class="flex-none text-muted-foreground hover:text-foreground transition-colors"
          @click.stop="cancelEdit">
          <X class="h-5 w-5" />
        </button>
      </template>

      <!-- 普通显示模式 -->
      <template v-else>
        <!-- 非管理模式：名称完整显示，无 tooltip -->
        <button v-if="!catLib.manageMode.value" type="button"
          class="flex-1 min-w-0 flex items-start gap-1.5 text-lg font-thin text-left transition-colors cursor-pointer"
          :class="effectiveSelectedId === cat.id
            ? 'text-primary font-extralight'
            : 'text-foreground hover:text-primary/80'" @click.stop="handleCatClick(cat.id)">
          <span class="whitespace-pre-wrap break-words">{{ cat.name }}</span>
          <span v-if="totalCount > 0"
            class="shrink-0 inline-flex items-center justify-center h-5 min-w-5 px-1.5 rounded-full bg-muted text-xs tabular-nums text-muted-foreground">
            {{ totalCount }}
          </span>
        </button>

        <!-- 管理模式：名称截断 + tooltip（hover 操作组出现时名称可能被截断） -->
        <Tooltip v-else>
          <TooltipTrigger as-child>
            <button type="button"
              class="flex-1 min-w-0 flex items-center gap-1.5 truncate text-lg font-thin text-left transition-colors cursor-pointer"
              :class="effectiveSelectedId === cat.id
                ? 'text-primary font-extralight'
                : 'text-foreground hover:text-primary/80'" @click.stop="handleCatClick(cat.id)">
              <span ref="nameEl" class="truncate">{{ cat.name }}</span>
              <span v-if="totalCount > 0"
                class="shrink-0 inline-flex items-center justify-center h-5 min-w-5 px-1.5 rounded-full bg-muted text-xs tabular-nums text-muted-foreground group-hover/cat:hidden"
                :class="{ '!hidden': addChildOpen || deleteOpen }">
                {{ totalCount }}
              </span>
            </button>
          </TooltipTrigger>
          <TooltipContent v-if="isNameTruncated" side="top" class="max-w-64 break-words [text-wrap:wrap]">{{ cat.name }}
          </TooltipContent>
        </Tooltip>

        <!-- 非管理模式：始终可见的展开/折叠箭头 -->
        <button v-if="!catLib.manageMode.value && hasChildren" type="button"
          class="p-2 rounded-full text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
          @click.stop="isOpen = !isOpen">
          <ChevronRight class="h-4 w-4 transition-transform duration-150" :class="{ 'rotate-90': isOpen }" />
        </button>

        <!-- 管理模式：hover 浮出操作组 -->
        <div v-if="catLib.manageMode.value"
          class="relative flex shrink-0 overflow-hidden max-w-0 opacity-0 group-hover/cat:max-w-[156px] group-hover/cat:opacity-100 transition-all duration-200 ease-out"
          :class="{ '!max-w-[156px] !opacity-100': addChildOpen || deleteOpen }">
          <div class="relative flex items-center gap-1 pl-2 pr-1">
            <span class="h-4 w-px shrink-0 bg-border/60" />
            <!-- 折叠箭头 -->
            <button type="button" class="p-2 text-muted-foreground hover:text-foreground transition-colors"
              :class="{ 'opacity-0 pointer-events-none': !hasChildren }" @click.stop="isOpen = !isOpen">
              <ChevronRight class="h-4 w-4 transition-transform duration-150"
                :class="{ 'rotate-90': isOpen && hasChildren }" />
            </button>

            <!-- 添加子分类 Popover（仅一级分类显示） -->
            <Popover v-if="cat.level < 1" v-model:open="addChildOpen">
              <PopoverTrigger as-child>
                <button type="button" class="p-1 transition-colors" :class="addChildOpen
                  ? 'text-primary'
                  : 'text-muted-foreground hover:text-foreground'" @click.stop>
                  <Plus class="h-4 w-4" />
                </button>
              </PopoverTrigger>
              <PopoverContent align="end" class="w-64 rounded-2xl p-5">
                <form class="space-y-2.5" @submit.prevent="confirmAddChild">
                  <p class="text-sm font-medium text-muted-foreground/70 tracking-wide">新建子分类</p>
                  <Input v-model="addChildName" placeholder="子分类名称" class="h-8 rounded-full text-sm" autofocus
                    @input="onAddChildNameInput" />
                  <Input v-model="addChildSlug" placeholder="slug" class="h-8 rounded-full text-sm font-mono" />
                  <div class="flex justify-end gap-1.5">
                    <Button variant="ghost" size="icon" type="button" class="h-6 w-6 rounded-full"
                      @click="addChildOpen = false">
                      <X class="h-3 w-3" />
                    </Button>
                    <Button type="submit" size="icon" class="h-6 w-6 rounded-full" :disabled="!addChildName">
                      <Check class="h-3 w-3" />
                    </Button>
                  </div>
                </form>
              </PopoverContent>
            </Popover>

            <!-- 内联编辑 -->
            <button type="button" class="p-1 text-muted-foreground hover:text-foreground transition-colors"
              @click.stop="startEdit">
              <Pencil class="h-4 w-4" />
            </button>

            <!-- 删除 Popover -->
            <Popover v-model:open="deleteOpen">
              <PopoverTrigger as-child>
                <button type="button" class="p-1 transition-colors" :class="deleteOpen
                  ? 'text-destructive'
                  : 'text-muted-foreground hover:text-destructive'" @click.stop>
                  <Trash2 class="h-4 w-4" />
                </button>
              </PopoverTrigger>
              <PopoverContent align="end" class="w-64 rounded-2xl p-5">
                <p class="text-base font-medium">确认删除该分类？</p>
                <p class="mt-1.5 text-sm text-muted-foreground/60">包含所有子分类，不可撤销。</p>
                <div class="mt-2.5 flex justify-end gap-1.5">
                  <Button variant="ghost" size="icon"
                    class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                    :disabled="deleting" @click="deleteOpen = false">
                    <X class="h-3.5 w-3.5" />
                  </Button>
                  <Button size="icon"
                    class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                    :disabled="deleting" @click="confirmDelete">
                    <Check class="h-3.5 w-3.5" />
                  </Button>
                </div>
              </PopoverContent>
            </Popover>
          </div>
        </div>
      </template>
    </div>

    <!-- 子分类（CSS Grid 展开/折叠动画） -->
    <div v-if="hasChildren" class="grid transition-[grid-template-rows] duration-200 ease-in-out"
      :class="isOpen ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'">
      <div class="overflow-hidden">
        <div class="relative pl-4 space-y-0.5 transition-opacity duration-200"
          :class="isOpen ? 'opacity-100' : 'opacity-0'">
          <div class="absolute left-[9px] top-0 bottom-0 w-px bg-border/60" />
          <div v-for="child in cat.children" :key="child.id" class="relative -ml-[7px] transition-colors duration-150"
            :class="effectiveSelectedId === child.id ? 'bg-accent/40' : 'hover:bg-accent/60'">
            <span v-if="effectiveSelectedId === child.id"
              class="absolute left-0 inset-y-0 w-px bg-primary pointer-events-none" />
            <CategorySidebarItem :cat="child" :cat-lib="catLib" :selected-id="selectedId" :on-select="onSelect" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
