<!--
  DocPageInlineEditor.vue - 文档页面内联编辑器
  职责：提供文档标题和 slug 的内联编辑功能，支持确认/取消操作
  对外暴露：props (title, slug, slugError, slugHint)，emits (update:title, update:slug, confirm, cancel)
-->
<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Check, X } from 'lucide-vue-next'

const props = defineProps<{
  title: string
  slug: string
  slugError: string
  slugHint: string
}>()

const emit = defineEmits<{
  'update:title': [value: string]
  'update:slug': [value: string]
  confirm: []
  cancel: []
}>()

const titleValue = computed({
  get: () => props.title,
  set: (next: string) => emit('update:title', next),
})

const slugValue = computed({
  get: () => props.slug,
  set: (next: string) => emit('update:slug', next),
})

function onConfirm() {
  if (!props.slugError) emit('confirm')
}

function onCancel() {
  emit('cancel')
}
</script>

<template>
  <div class="flex flex-col gap-1.5 py-0.5">
    <Input v-model="titleValue" placeholder="标题" class="h-8 rounded-full text-sm" @keyup.esc="onCancel" />
    <div>
      <Input v-model="slugValue" placeholder="slug" class="h-8 rounded-full text-sm font-mono"
        :class="{ 'border-destructive': slugError }" @keyup.enter="onConfirm" @keyup.esc="onCancel" />
      <p v-if="slugError" class="text-xs text-destructive mt-0.5">{{ slugError }}</p>
      <p v-else class="text-xs text-muted-foreground mt-0.5">{{ slugHint }}</p>
    </div>
    <div class="flex gap-1 justify-end">
      <Tooltip>
        <TooltipTrigger as-child>
          <Button size="icon" class="rounded-full h-6 w-6" :disabled="!!slugError" @click.stop="onConfirm">
            <Check class="h-3 w-3" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>确认</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon" class="rounded-full h-6 w-6" @click.stop="onCancel">
            <X class="h-3 w-3" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>取消</TooltipContent>
      </Tooltip>
    </div>
  </div>
</template>
