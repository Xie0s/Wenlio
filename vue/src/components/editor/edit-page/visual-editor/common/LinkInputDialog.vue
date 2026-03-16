<!--
  LinkInputDialog.vue - 链接输入/编辑对话框
  职责：提供链接地址和显示文本的输入界面，支持新建和编辑模式
  对外暴露：Props: open, initialUrl, initialText, isEditMode
            Emits: update:open, confirm, remove
-->
<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import { Link2, X, Check, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'

const props = withDefaults(defineProps<{
  open: boolean
  initialUrl?: string
  initialText?: string
  isEditMode?: boolean
  title?: string
}>(), {
  initialUrl: '',
  initialText: '',
  isEditMode: false,
  title: '',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  confirm: [url: string, text?: string]
  remove: []
}>()

const url = ref('')
const text = ref('')
const urlError = ref('')
const urlInputRef = ref<{ $el: HTMLInputElement } | null>(null)

function blurActiveElement() {
  const activeEl = document.activeElement
  if (activeEl instanceof HTMLElement) activeEl.blur()
}

// 对话框打开时重置状态
watch(() => props.open, (open) => {
  if (open) {
    blurActiveElement()
    url.value = props.initialUrl
    text.value = props.initialText
    urlError.value = ''
    nextTick(() => {
      setTimeout(() => {
        const el = urlInputRef.value?.$el
        el?.focus()
        el?.select()
      }, 100)
    })
  }
})

// URL 格式校验
function isValidUrl(val: string): boolean {
  if (!val.trim()) return false
  return /^(https?:\/\/|mailto:|tel:|\/|#)/i.test(val) || val.startsWith('www.')
}

// 标准化 URL（添加协议前缀）
function normalizeUrl(raw: string): string {
  const trimmed = raw.trim()
  if (!trimmed) return ''
  if (/^(https?:\/\/|mailto:|tel:|\/|#)/i.test(trimmed)) return trimmed
  return `https://${trimmed}`
}

function handleConfirm() {
  const trimmedUrl = url.value.trim()
  if (!trimmedUrl) {
    urlError.value = '请输入链接地址'
    urlInputRef.value?.$el?.focus()
    return
  }
  const normalizedUrl = normalizeUrl(trimmedUrl)
  if (!isValidUrl(normalizedUrl)) {
    urlError.value = '请输入有效的链接地址'
    urlInputRef.value?.$el?.focus()
    return
  }
  emit('confirm', normalizedUrl, text.value.trim() || undefined)
  emit('update:open', false)
}

function handleRemove() {
  emit('remove')
  emit('update:open', false)
}

function handleCancel() {
  emit('update:open', false)
}

function handleUrlInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  url.value = val
  if (urlError.value && val.trim()) urlError.value = ''
}

function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleConfirm()
  }
}

const dialogTitle = computed(() => props.title || (props.isEditMode ? '编辑链接' : '插入链接'))
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent
      class="sm:max-w-md rounded-3xl"
      @keydown="handleKeyDown"
    >
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2">
          <Link2 class="w-5 h-5" />
          {{ dialogTitle }}
        </DialogTitle>
        <DialogDescription>
          {{ isEditMode ? '修改链接地址或文本' : '输入链接地址，可选择性地修改显示文本' }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4 py-4">
        <!-- URL 输入 -->
        <div class="space-y-2">
          <Label for="link-url" class="text-sm font-medium">
            链接地址 <span class="text-destructive">*</span>
          </Label>
          <Input
            id="link-url"
            ref="urlInputRef"
            type="text"
            placeholder="https://example.com"
            :model-value="url"
            :class="urlError ? 'border-destructive focus-visible:ring-destructive' : ''"
            autocomplete="off"
            @input="handleUrlInput"
          />
          <p v-if="urlError" class="text-xs text-destructive">{{ urlError }}</p>
        </div>

        <!-- 链接文本输入 -->
        <div class="space-y-2">
          <Label for="link-text" class="text-sm font-medium">
            显示文本 <span class="text-muted-foreground text-xs">（可选）</span>
          </Label>
          <Input
            id="link-text"
            v-model="text"
            type="text"
            :placeholder="isEditMode ? '保持原有文本' : '留空则使用链接地址'"
            autocomplete="off"
          />
        </div>
      </div>

      <DialogFooter class="flex items-center justify-between gap-3 pt-2">
        <!-- 移除按钮（仅编辑模式） -->
        <Tooltip v-if="isEditMode">
          <TooltipTrigger as-child>
            <Button
              type="button"
              variant="destructive"
              size="icon"
              class="rounded-full"
              @click="handleRemove"
            >
              <Trash2 class="w-4 h-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>移除链接</TooltipContent>
        </Tooltip>

        <div class="flex items-center justify-end gap-3 ml-auto">
          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                type="button"
                variant="outline"
                size="icon"
                class="rounded-full"
                @click="handleCancel"
              >
                <X class="w-5 h-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>取消</TooltipContent>
          </Tooltip>

          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                type="button"
                size="icon"
                class="rounded-full"
                @click="handleConfirm"
              >
                <Check class="w-5 h-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>{{ isEditMode ? '更新' : '插入' }}</TooltipContent>
          </Tooltip>
        </div>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
