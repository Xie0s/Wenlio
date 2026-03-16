<!--
  FilePickerDialog.vue - 已有文件选择浮窗
  职责：以可拖拽浮窗形式展示媒体库文件，左右分栏，多选后插入编辑器
  当选中文件中包含图片时，点击确认弹出 Popover 让用户选择「图片插入」或「文件插入」
  对外暴露：Props: open (v-model), Emits: select(items, insertType: 'image' | 'file')
-->
<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useMedia, type MediaItem } from '@/lib/media'
import { useFloatingPanel } from '@/composables/useFloatingPanel'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverAnchor, PopoverContent } from '@/components/ui/popover'
import { Loader2, X, Paperclip, File, Check, Image, FileText } from 'lucide-vue-next'
import MediaPickerPanel from '@/components/editor/media/MediaPickerPanel.vue'
import FileCard from './FileCard.vue'

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  select: [items: MediaItem[], insertType: 'image' | 'file']
}>()

const { loading, items, loadMedia } = useMedia()
const selected = ref<Set<string>>(new Set())

const insertPopoverOpen = ref(false)

const { panelRef, style, onMouseDown, reset } = useFloatingPanel()

watch(open, async (val) => {
  if (val) {
    ;(document.activeElement as HTMLElement)?.blur()
    selected.value = new Set()
    insertPopoverOpen.value = false
    reset()
    loadMedia()
    await nextTick()
    panelRef.value?.getBoundingClientRect()
  }
}, { flush: 'sync' })

function toggle(id: string) {
  const s = new Set(selected.value)
  if (s.has(id)) s.delete(id)
  else s.add(id)
  selected.value = s
}

const pickedItems = computed(() => items.value.filter(m => selected.value.has(m.id)))
const hasImage = computed(() => pickedItems.value.some(m => m.mime_type?.startsWith('image/')))
const firstImageItem = computed(() => pickedItems.value.find(m => m.mime_type?.startsWith('image/')))

function confirm() {
  if (pickedItems.value.length === 0) return
  if (hasImage.value) {
    insertPopoverOpen.value = true
  } else {
    doInsert('file')
  }
}

function doInsert(insertType: 'image' | 'file') {
  insertPopoverOpen.value = false
  emit('select', pickedItems.value, insertType)
  open.value = false
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-all duration-150 ease-out"
      leave-active-class="transition-all duration-100 ease-in"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        ref="panelRef"
        class="fixed z-50 flex flex-col bg-background border rounded-3xl shadow-2xl overflow-hidden"
        style="width: 860px; height: 580px;"
        :style="style"
      >
        <!-- 拖拽标题栏 -->
        <div
          class="flex items-center gap-2.5 px-5 py-3.5 border-b cursor-move select-none shrink-0"
          @mousedown="onMouseDown"
        >
          <Paperclip class="size-4 text-muted-foreground shrink-0" :stroke-width="1.5" />
          <h3 class="text-sm font-medium flex-1">插入已有文件</h3>
          <span v-if="selected.size > 0" class="text-xs text-primary font-medium">
            已选 {{ selected.size }} 个
          </span>
          <button
            class="size-7 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted/60 transition-colors"
            @click="open = false"
          >
            <X class="size-4" :stroke-width="1.5" />
          </button>
        </div>

        <!-- 主体内容 -->
        <div class="flex-1 overflow-hidden">
          <!-- 加载态 -->
          <div v-if="loading" class="flex items-center justify-center h-full">
            <Loader2 class="size-6 animate-spin text-muted-foreground" />
          </div>

          <!-- 空态：媒体库为空 -->
          <div v-else-if="items.length === 0" class="flex flex-col items-center justify-center h-full text-muted-foreground">
            <File class="size-12 mb-3 opacity-20" :stroke-width="1" />
            <p class="text-sm">暂无可用文件</p>
          </div>

          <!-- 左右分栏选择面板 -->
          <MediaPickerPanel
            v-else
            :items="items"
            :selected="selected"
            @toggle="toggle"
          />
        </div>

        <!-- 底部操作栏 -->
        <div class="flex items-center px-5 py-3 border-t shrink-0">
          <span class="flex-1 text-xs text-primary font-medium">
            <template v-if="selected.size > 0">已选 {{ selected.size }} 个</template>
          </span>
          <div class="flex items-center gap-2">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="outline" size="icon" class="size-8 rounded-full" @click="open = false">
                  <X class="size-4" :stroke-width="1.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>取消</TooltipContent>
            </Tooltip>
            <Popover v-model:open="insertPopoverOpen">
              <PopoverAnchor as-child>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button size="icon" class="size-8 rounded-full" :disabled="selected.size === 0" @click="confirm">
                      <Check class="size-4" :stroke-width="2" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>插入</TooltipContent>
                </Tooltip>
              </PopoverAnchor>
              <PopoverContent side="top" align="end" class="w-72 rounded-3xl p-3 z-[200]">
                <p class="text-xs text-muted-foreground mb-3 font-medium">选择插入方式</p>
                <div class="flex flex-col gap-2">
                  <!-- 图片插入 -->
                  <button
                    class="flex items-start gap-3 p-3 rounded-xl border hover:border-primary/50 hover:bg-primary/5 transition-colors text-left w-full"
                    @click="doInsert('image')"
                  >
                    <div class="size-12 rounded-lg bg-muted overflow-hidden shrink-0 flex items-center justify-center">
                      <img v-if="firstImageItem" :src="firstImageItem.file_url" class="size-full object-cover" />
                      <Image v-else class="size-5 text-muted-foreground" :stroke-width="1.5" />
                    </div>
                    <div>
                      <p class="text-sm font-medium">图片插入</p>
                      <p class="text-xs text-muted-foreground">以图片形式直接显示</p>
                    </div>
                  </button>
                  <!-- 文件插入 -->
                  <button
                    class="flex items-start gap-3 p-3 rounded-xl border hover:border-primary/50 hover:bg-primary/5 transition-colors text-left w-full"
                    @click="doInsert('file')"
                  >
                    <div class="size-9 rounded-lg bg-muted flex items-center justify-center shrink-0 self-center">
                      <FileText class="size-4 text-muted-foreground" :stroke-width="1.5" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-medium mb-1.5">文件插入</p>
                      <FileCard
                        v-if="pickedItems[0]"
                        :file-name="pickedItems[0].file_name"
                        :file-size="pickedItems[0].file_size"
                        :mime-type="pickedItems[0].mime_type"
                      />
                    </div>
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
