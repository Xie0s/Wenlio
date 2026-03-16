<!--
  ConfirmDialog.vue - 通用确认对话框组件
  职责：封装 AlertDialog，提供删除、封禁、保存、重置四种预设确认场景，支持自定义标题和描述。
  对外接口：
    Props:
      - open: boolean        — v-model 控制显隐
      - type: DialogType     — 'delete' | 'ban' | 'save' | 'reset'，默认 'delete'
      - title?: string       — 覆盖预设标题
      - description?: string — 覆盖预设描述
    Emits:
      - update:open(value: boolean)
      - confirm()
      - cancel()
-->
<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import { Trash2, ShieldBan, Save, RotateCcw } from 'lucide-vue-next'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogAction,
  AlertDialogCancel,
} from '@/components/ui/alert-dialog'

type DialogType = 'delete' | 'ban' | 'save' | 'reset'

interface Preset {
  icon: Component
  title: string
  description: string
  iconBg: string
  iconColor: string
}

const props = withDefaults(defineProps<{
  open: boolean
  type?: DialogType
  title?: string
  description?: string
}>(), {
  type: 'delete',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  confirm: []
  cancel: []
}>()

const presets: Record<DialogType, Preset> = {
  delete: {
    icon: Trash2,
    title: '确认删除',
    description: '此操作无法撤销，数据将被永久删除。',
    iconBg: 'bg-destructive/10 dark:bg-destructive/20',
    iconColor: 'text-destructive',
  },
  ban: {
    icon: ShieldBan,
    title: '确认封禁',
    description: '封禁后该用户将无法登录系统，可随时在管理页面解除。',
    iconBg: 'bg-orange-500/10 dark:bg-orange-500/20',
    iconColor: 'text-orange-500',
  },
  save: {
    icon: Save,
    title: '确认保存',
    description: '确认保存当前更改？',
    iconBg: 'bg-primary/10 dark:bg-primary/20',
    iconColor: 'text-primary',
  },
  reset: {
    icon: RotateCcw,
    title: '确认重置',
    description: '将当前编辑内容重置为默认配置，未保存的更改会丢失。',
    iconBg: 'bg-amber-500/10 dark:bg-amber-500/20',
    iconColor: 'text-amber-500',
  },
}

const current = computed(() => presets[props.type])
const dialogTitle = computed(() => props.title ?? current.value.title)
const dialogDescription = computed(() => props.description ?? current.value.description)

function onConfirm() {
  emit('update:open', false)
  emit('confirm')
}

function onCancel() {
  emit('update:open', false)
  emit('cancel')
}
</script>

<template>
  <AlertDialog :open="open" @update:open="(v) => emit('update:open', v)">
    <AlertDialogContent>
      <AlertDialogHeader>
        <div class="flex items-center gap-3">
          <div
            :class="[
              'flex h-10 w-10 shrink-0 items-center justify-center rounded-full',
              current.iconBg,
              current.iconColor,
            ]"
          >
            <component :is="current.icon" class="size-5" />
          </div>
          <AlertDialogTitle>{{ dialogTitle }}</AlertDialogTitle>
        </div>
      </AlertDialogHeader>
      <AlertDialogDescription class="pl-13">
        {{ dialogDescription }}
      </AlertDialogDescription>
      <AlertDialogFooter>
        <AlertDialogCancel @click="onCancel" />
        <AlertDialogAction @click="onConfirm" />
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
