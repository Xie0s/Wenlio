<!-- ThemeCard.vue - 主题列表卡片
     职责：展示单个主题信息，提供 hover 操作入口（进入编辑器、编辑信息、查看页面、删除）
     对外暴露：
       Props: theme(Theme), deleting(boolean), cascadeOpen(boolean), catLib, tagLib
       Emits: enter, open-edit, view, request-delete, cascade-confirm, cascade-cancel -->
<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { VERSION_STATUS_COLOR, VERSION_STATUS_LABEL, ACCESS_MODE_LABEL, ACCESS_MODE_COLOR, type Theme, type Category, type Tag, type VersionStatus, type AccessMode } from '@/utils/types'
import {
  FileText, Layers, LayoutList, SquarePen, ListFilterPlus, CircleArrowOutUpRight, Trash2, Check, X, Calendar, BookOpen,
} from 'lucide-vue-next'
import ThemePreviewModal from './ThemePreviewModal.vue'

const props = defineProps<{
  theme: Theme
  deleting: boolean
  cascadeOpen: boolean
  catLib: { categories: Ref<Category[]> }
  tagLib: { tags: Ref<Tag[]> }
  readonly?: boolean
}>()

const emit = defineEmits<{
  enter: []
  'open-edit': []
  view: []
  'request-delete': []
  'cascade-confirm': []
  'cascade-cancel': []
}>()

function flatCategories(
  list: Category[],
  depth = 0,
): Array<Category & { depth: number }> {
  const result: Array<Category & { depth: number }> = []
  for (const cat of list) {
    result.push({ ...cat, depth })
    if (cat.children?.length) result.push(...flatCategories(cat.children, depth + 1))
  }
  return result
}

const categoryName = computed(() => {
  if (!props.theme.category_id) return ''
  const flat = flatCategories(props.catLib.categories.value)
  return flat.find((c) => c.id === props.theme.category_id)?.name ?? ''
})

const isReadonlyClickable = computed(() => !!props.readonly)

const previewOpen = ref(false)

function getVersionStatusLabel(status?: string) {
  if (!status) return '未设置'
  return VERSION_STATUS_LABEL[status as VersionStatus] ?? status
}

function getVersionStatusClass(status?: string) {
  if (!status) return 'bg-muted text-muted-foreground'
  return VERSION_STATUS_COLOR[status as VersionStatus] ?? 'bg-muted text-muted-foreground'
}

function formatDate(dateStr: string) {
  if (!dateStr || dateStr.startsWith('0001')) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function handleReadonlyEnter() {
  if (!isReadonlyClickable.value) return
  emit('enter')
}
</script>

<template>
  <!-- 列表行风格：无卡片边框 + 左侧指示线 + 操作按钮位移淡入 + 渐变分隔线 -->
  <div
    class="group relative flex min-h-[76px] items-center gap-5 px-4 py-3 transition-all duration-200 hover:bg-accent/50"
    :class="isReadonlyClickable ? 'cursor-pointer' : ''"
    :role="isReadonlyClickable ? 'button' : undefined"
    :tabindex="isReadonlyClickable ? 0 : undefined"
    @click="handleReadonlyEnter"
    @keydown.enter="handleReadonlyEnter"
    @keydown.space.prevent="handleReadonlyEnter"
  >

    <!-- 主信息区 -->
    <div class="flex-1 min-w-0">
      <!-- 第一行：名称 + 分类 + 标签 -->
      <div class="flex items-center gap-2 flex-wrap" @click.stop="false">
        <span class="font-light text-xl leading-snug text-foreground group-hover:text-primary transition-colors duration-200">{{ theme.name }}</span>
        <span
          v-if="theme.access_mode && theme.access_mode !== 'public'"
          class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
          :class="ACCESS_MODE_COLOR[theme.access_mode as AccessMode]"
          @click.stop
        >
          {{ ACCESS_MODE_LABEL[theme.access_mode as AccessMode] }}
        </span>
        <span
          v-if="categoryName"
          class="inline-flex items-center rounded-full border border-primary/30 bg-primary/5 px-2 py-0.5 text-xs text-primary/80"
          @click.stop
        >
          {{ categoryName }}
        </span>
        <span
          v-for="tagId in theme.tag_ids"
          :key="tagId"
          class="text-xs text-muted-foreground/70"
          @click.stop
        >
          #{{ tagLib.tags.value.find((tg: Tag) => tg.id === tagId)?.name ?? tagId }}
        </span>
      </div>

      <!-- 第二行：描述（两行截断） -->
      <p v-if="theme.description" class="mt-1 text-base font-light text-foreground/60 line-clamp-1 leading-relaxed">
        {{ theme.description }}
      </p>
    </div>

    <!-- 右侧：统计 + hover 操作浮层 -->
    <div class="relative flex-none flex items-center self-center">
      <!-- 统计（非 readonly 时 hover 淡出让位给操作浮层） -->
      <div
        class="flex items-center gap-3 text-sm text-muted-foreground transition-all duration-300 group-hover:opacity-0 group-hover:pointer-events-none"
      >
        <span v-if="theme.current_version?.name" class="hidden md:inline">
          {{ theme.current_version.name }}
        </span>
        <span class="flex items-center gap-1">
          <Layers class="h-3.5 w-3.5" :stroke-width="1.5" />
          {{ theme.version_count ?? 0 }}
        </span>
        <span class="flex items-center gap-1">
          <LayoutList class="h-3.5 w-3.5" :stroke-width="1.5" />
          {{ theme.section_count ?? 0 }}
        </span>
        <span class="flex items-center gap-1">
          <FileText class="h-3.5 w-3.5" :stroke-width="1.5" />
          {{ theme.page_count ?? 0 }}
        </span>
        <span v-if="formatDate(theme.created_at)" class="flex items-center gap-1 hidden md:flex">
          <Calendar class="h-3.5 w-3.5" :stroke-width="1.5" />
          {{ formatDate(theme.created_at) }}
        </span>
        <span
          v-if="theme.current_version"
          class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
          :class="getVersionStatusClass(theme.current_version.status)"
        >
          {{ getVersionStatusLabel(theme.current_version.status) }}
        </span>
      </div>

      <div
        v-if="readonly"
        class="pointer-events-none absolute right-0 top-1/2 flex -translate-y-1/2 items-center gap-2 rounded-full border border-border/50 bg-background/95 px-4 py-2 text-sm text-foreground backdrop-blur-md transition-all duration-300 opacity-0 translate-x-4 group-hover:translate-x-0 group-hover:opacity-100"
      >
        <CircleArrowOutUpRight class="h-4 w-4" />
        <span>点击打开</span>
      </div>

      <!-- 操作按钮浮层（readonly 模式隐藏） -->
      <div
        v-if="!readonly"
        class="absolute right-0 top-1/2 flex -translate-y-1/2 items-center gap-0.5 rounded-full border border-border/40 bg-background/95 px-1 py-1 backdrop-blur-md transition-all duration-300"
        :class="cascadeOpen 
          ? 'opacity-100 translate-x-0 pointer-events-auto' 
          : 'opacity-0 translate-x-4 group-hover:translate-x-0 group-hover:opacity-100 pointer-events-none group-hover:pointer-events-auto'"
        @click.stop
      >
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full h-10 w-10 hover:bg-accent" @click="emit('enter')">
              <SquarePen class="h-5 w-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>进入编辑器</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full h-10 w-10 hover:bg-accent" @click="emit('open-edit')">
              <ListFilterPlus class="h-5 w-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>编辑信息</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full h-10 w-10 hover:bg-accent" @click.stop="previewOpen = true">
              <BookOpen class="h-5 w-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>预览内容</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="icon" class="rounded-full h-10 w-10 hover:bg-accent" @click="emit('view')">
              <CircleArrowOutUpRight class="h-5 w-5" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>查看页面</TooltipContent>
        </Tooltip>

        <Popover :open="cascadeOpen" @update:open="(open) => open ? undefined : emit('cascade-cancel')">
          <PopoverTrigger as-child>
            <Button
              variant="ghost" size="icon"
              :class="['rounded-full h-10 w-10 text-destructive hover:text-destructive hover:bg-destructive/10', cascadeOpen ? 'bg-destructive/10' : '']"
              :disabled="deleting"
              @click="emit('request-delete')"
            >
              <Trash2 class="h-5 w-5" />
            </Button>
          </PopoverTrigger>
          <PopoverContent align="end" class="w-88 rounded-2xl p-6">
            <p class="text-base font-semibold">该主题下存在版本，是否级联删除？</p>
            <p class="mt-2 text-sm text-muted-foreground leading-relaxed">
              将依次删除该主题下的所有版本、章节及文档页，最后删除主题本身。已发布的内容将立即从读者端下线，此操作不可撤销。
            </p>
            <div class="mt-3 flex items-center justify-end gap-2">
              <Button
                variant="ghost"
                size="icon"
                class="h-7 w-7 rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                :disabled="deleting"
                @click="emit('cascade-cancel')">
                <X class="h-3.5 w-3.5" />
              </Button>
              <Button
                size="icon"
                class="h-7 w-7 rounded-full bg-destructive text-white transition-colors hover:bg-destructive/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/90"
                :disabled="deleting"
                @click="emit('cascade-confirm')">
                <Check class="h-3.5 w-3.5" />
              </Button>
            </div>
          </PopoverContent>
        </Popover>
      </div>
    </div>

    <!-- 底部渐变分隔线 -->
    <div class="absolute bottom-0 left-1/2 -translate-x-1/2 w-[90%] flex items-center">
      <div class="flex-1 h-px bg-gradient-to-l from-border/50 to-transparent" />
      <div class="w-0.5 h-0.5 rounded-full bg-border/40 mx-1.5" />
      <div class="flex-1 h-px bg-gradient-to-r from-border/50 to-transparent" />
    </div>
  </div>

  <ThemePreviewModal :theme="theme" :open="previewOpen" @update:open="previewOpen = $event" />
</template>
