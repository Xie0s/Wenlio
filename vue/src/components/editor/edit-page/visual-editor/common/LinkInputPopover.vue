<!--
  LinkInputPopover.vue - 链接输入 Popover 组件
  职责：以 Popover 形式提供链接地址输入界面，支持新建和编辑模式
  使用 PopoverAnchor 定位（而非 PopoverTrigger），避免与外部 Tooltip 嵌套时的事件冲突
  对外暴露：Props: open, initialUrl, initialText, showTextInput, isEditMode
            Emits: update:open, confirm(url, text?), remove
            Slot: default（接收 trigger 元素，需在外部通过 @click 控制 open）
-->
<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { Check, Trash2, X } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Popover, PopoverAnchor, PopoverContent } from '@/components/ui/popover'

const props = withDefaults(defineProps<{
  open: boolean
  initialUrl?: string
  initialText?: string
  showTextInput?: boolean
  isEditMode?: boolean
}>(), {
  initialUrl: '',
  initialText: '',
  showTextInput: false,
  isEditMode: false,
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  confirm: [url: string, text?: string]
  remove: []
}>()

const url = ref('')
const text = ref('')
const urlError = ref('')
const urlInputRef = ref<InstanceType<typeof Input> | null>(null)

watch(() => props.open, (open) => {
  if (open) {
    url.value = props.initialUrl
    text.value = props.initialText
    urlError.value = ''
    nextTick(() => {
      setTimeout(() => {
        const el = urlInputRef.value?.$el as HTMLInputElement | undefined
        el?.focus()
        el?.select()
      }, 50)
    })
  }
})

function isValidUrl(val: string): boolean {
  if (!val.trim()) return false
  return /^(https?:\/\/|mailto:|tel:|\/|#)/i.test(val) || val.startsWith('www.')
}

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
    ;(urlInputRef.value?.$el as HTMLInputElement | undefined)?.focus()
    return
  }
  const normalizedUrl = normalizeUrl(trimmedUrl)
  if (!isValidUrl(normalizedUrl)) {
    urlError.value = '请输入有效的链接地址'
    ;(urlInputRef.value?.$el as HTMLInputElement | undefined)?.focus()
    return
  }
  emit('confirm', normalizedUrl, text.value.trim() || undefined)
  emit('update:open', false)
}

function handleRemove() {
  emit('remove')
  emit('update:open', false)
}

function handleUrlInput() {
  if (urlError.value && url.value.trim()) urlError.value = ''
}

function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleConfirm()
  }
}
</script>

<template>
  <Popover :open="open" @update:open="emit('update:open', $event)">
    <PopoverAnchor as-child>
      <slot />
    </PopoverAnchor>
    <PopoverContent
      class="w-72 p-3 rounded-2xl"
      :side-offset="8"
      @keydown="handleKeyDown"
      @open-auto-focus.prevent
    >
      <div class="space-y-3">
        <div class="space-y-1.5">
          <Label class="text-xs font-medium">链接地址</Label>
          <Input
            ref="urlInputRef"
            v-model="url"
            type="text"
            placeholder="https://example.com"
            :class="urlError ? 'border-destructive' : ''"
            autocomplete="off"
            class="h-8 text-sm"
            @input="handleUrlInput"
          />
          <p v-if="urlError" class="text-xs text-destructive">{{ urlError }}</p>
        </div>
        <div v-if="showTextInput" class="space-y-1.5">
          <Label class="text-xs font-medium">显示文本 <span class="text-muted-foreground text-[10px]">(可选)</span></Label>
          <Input
            v-model="text"
            type="text"
            :placeholder="isEditMode ? '保持原有文本' : '留空则使用链接地址'"
            autocomplete="off"
            class="h-8 text-sm"
          />
        </div>
        <div class="flex items-center justify-between">
          <div>
            <Tooltip v-if="isEditMode">
              <TooltipTrigger as-child>
                <Button variant="destructive" size="icon" class="h-7 w-7 rounded-full" @click="handleRemove">
                  <Trash2 class="w-3.5 h-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>移除链接</TooltipContent>
            </Tooltip>
          </div>
          <div class="flex gap-1.5">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="outline" size="icon" class="h-7 w-7 rounded-full" @click="emit('update:open', false)">
                  <X class="w-3.5 h-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>取消</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger as-child>
                <Button size="icon" class="h-7 w-7 rounded-full" @click="handleConfirm">
                  <Check class="w-3.5 h-3.5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>{{ isEditMode ? '更新' : '插入' }}</TooltipContent>
            </Tooltip>
          </div>
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
