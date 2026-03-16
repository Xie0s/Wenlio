<!-- ThemeAccessSettings.vue - 主题访问权限设置组件
     职责：提供 access_mode 选择（公开/登录可见/需验证码）和验证码输入，复用于 ThemeFormDialog 和 EditorSettingsPanel
     对外暴露：
       Props: accessMode(AccessMode), accessCode(string)
       Emits: update:accessMode(AccessMode), update:accessCode(string) -->
<script setup lang="ts">
import { computed } from 'vue'
import type { AccessMode } from '@/utils/types'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Globe, Lock, KeyRound, CircleHelp } from 'lucide-vue-next'

const props = defineProps<{
  accessMode: AccessMode
  accessCode: string
}>()

const emit = defineEmits<{
  'update:accessMode': [AccessMode]
  'update:accessCode': [string]
}>()

const options: { value: AccessMode; label: string; icon: typeof Globe; desc: string }[] = [
  { value: 'public', label: '公开', icon: Globe, desc: '所有人可见' },
  { value: 'login', label: '登录可见', icon: Lock, desc: '未登录时从列表隐藏' },
  { value: 'code', label: '需验证码', icon: KeyRound, desc: '列表可见，访问需输入验证码' },
]

const showCodeInput = computed(() => props.accessMode === 'code')
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="flex items-center gap-1.5">
      <Label class="text-sm font-medium">访问权限</Label>
      <Tooltip>
        <TooltipTrigger as-child>
          <CircleHelp class="h-3.5 w-3.5 text-muted-foreground/50 cursor-help" />
        </TooltipTrigger>
        <TooltipContent side="top" class="max-w-[240px]">
          <p class="text-xs leading-relaxed">
            <strong>公开</strong>：所有人可浏览<br />
            <strong>登录可见</strong>：未登录时主题从列表隐藏<br />
            <strong>需验证码</strong>：列表可见但访问需输入6位验证码
          </p>
        </TooltipContent>
      </Tooltip>
    </div>

    <div class="flex gap-1.5">
      <button
        v-for="opt in options"
        :key="opt.value"
        type="button"
        class="flex-1 flex flex-col items-center gap-1 rounded-xl border px-2 py-2 text-xs transition-colors"
        :class="accessMode === opt.value
          ? 'border-primary/50 bg-primary/5 text-primary'
          : 'border-border/60 text-muted-foreground hover:border-border hover:text-foreground'"
        @click="emit('update:accessMode', opt.value)"
      >
        <component :is="opt.icon" class="h-4 w-4" />
        <span class="font-medium">{{ opt.label }}</span>
      </button>
    </div>

    <div v-if="showCodeInput" class="flex flex-col gap-1">
      <Label class="text-xs text-muted-foreground">验证码（6位字母数字组合）</Label>
      <Input
        :model-value="accessCode"
        placeholder="如 Ab3xK9"
        maxlength="6"
        class="rounded-xl font-mono text-sm tracking-widest uppercase"
        @update:model-value="emit('update:accessCode', ($event as string).toUpperCase())"
      />
      <p v-if="accessCode && accessCode.length < 6" class="text-xs text-destructive">
        验证码需要6位
      </p>
    </div>
  </div>
</template>
