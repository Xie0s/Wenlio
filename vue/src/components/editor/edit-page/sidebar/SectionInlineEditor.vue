<!--
  SectionInlineEditor.vue - 章节标题行内编辑器
  职责：提供章节标题的轻量级行内输入、确认与取消操作
  对外暴露：`modelValue`、`placeholder` props；`update:modelValue`、`confirm`、`cancel` emits
-->
<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Check, X } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
}>(), {
  placeholder: '章节标题',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  confirm: []
  cancel: []
}>()

const value = computed({
  get: () => props.modelValue,
  set: (next: string) => emit('update:modelValue', next),
})

function onConfirm() {
  emit('confirm')
}

function onCancel() {
  emit('cancel')
}
</script>

<template>
  <Input
    v-model="value"
    :placeholder="placeholder"
    class="h-9 rounded-xl text-sm flex-1 min-w-0"
    @keyup.enter="onConfirm"
    @keyup.esc="onCancel"
  />
  <Tooltip>
    <TooltipTrigger as-child>
      <Button size="icon" class="rounded-full h-8 w-8 shrink-0" @click.stop="onConfirm">
        <Check class="h-3.5 w-3.5" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>确认</TooltipContent>
  </Tooltip>
  <Tooltip>
    <TooltipTrigger as-child>
      <Button variant="ghost" size="icon" class="rounded-full h-8 w-8 shrink-0" @click.stop="onCancel">
        <X class="h-3.5 w-3.5" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>取消</TooltipContent>
  </Tooltip>
</template>
