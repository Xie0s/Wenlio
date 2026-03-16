<!-- ColorPickerInput.vue - 带可视化预览的颜色输入控件
     支持 hex (#ffffff) 和 rgba(r,g,b,a) 两种格式
     Props:
       modelValue  - 颜色字符串
       placeholder - 文本框占位符
       supportAlpha - 是否显示透明度滑条（用于 rgba 场景） -->
<script setup lang="ts">
import { computed } from 'vue'
import { Input } from '@/components/ui/input'

const props = defineProps<{
  modelValue?: string
  placeholder?: string
  supportAlpha?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const hasValue = computed(() => !!props.modelValue?.trim())

const isRgba = computed(() =>
  !!props.modelValue?.trim().match(/^rgba?\s*\(/i),
)

/** 从当前值中提取 hex（供原生 color picker 使用） */
const hexFromValue = computed(() => {
  const val = props.modelValue?.trim() || ''
  if (!val) return '#000000'
  if (val.startsWith('#')) {
    const hex = val.replace('#', '')
    if (hex.length === 3) return '#' + hex.split('').map((c) => c + c).join('')
    return val.length >= 7 ? val.slice(0, 7) : '#000000'
  }
  const m = val.match(/rgba?\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)/)
  if (m) {
    const r = parseInt(m[1] ?? '0', 10).toString(16).padStart(2, '0')
    const g = parseInt(m[2] ?? '0', 10).toString(16).padStart(2, '0')
    const b = parseInt(m[3] ?? '0', 10).toString(16).padStart(2, '0')
    return `#${r}${g}${b}`
  }
  return '#000000'
})

/** 当前透明度 0~1 */
const alphaFromValue = computed(() => {
  const m = props.modelValue?.match(/rgba?\s*\(\s*\d+\s*,\s*\d+\s*,\s*\d+\s*,\s*([\d.]+)/)
  return m ? Math.min(1, Math.max(0, parseFloat(m[1] ?? '1'))) : 1
})

const swatchStyle = computed(() => {
  if (!hasValue.value) return {}
  return { backgroundColor: props.modelValue }
})

function hexToRgba(hex: string, alpha: number) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

function onColorPickerInput(e: Event) {
  const hex = (e.target as HTMLInputElement).value
  if (isRgba.value || props.supportAlpha) {
    emit('update:modelValue', hexToRgba(hex, alphaFromValue.value))
  } else {
    emit('update:modelValue', hex)
  }
}

function onAlphaInput(e: Event) {
  const alpha = parseFloat((e.target as HTMLInputElement).value)
  emit('update:modelValue', hexToRgba(hexFromValue.value, alpha))
}

function onTextInput(val: string | number) {
  emit('update:modelValue', String(val))
}
</script>

<template>
  <div class="space-y-2">
    <!-- 色块 + 文本输入 -->
    <div class="flex items-center gap-2">
      <!-- 颜色预览色块（同时作为原生 color picker 触发器） -->
      <div
        class="relative h-9 w-9 flex-shrink-0 cursor-pointer overflow-hidden rounded-xl border border-input shadow-sm"
        title="点击选色"
      >
        <!-- 透明棋盘格背景（无值时显示） -->
        <div
          v-if="!hasValue"
          class="absolute inset-0"
          style="background: repeating-conic-gradient(#d0d0d0 0% 25%, #f8f8f8 0% 50%) 0 0 / 8px 8px"
        />
        <!-- 颜色填充 -->
        <div v-else class="absolute inset-0" :style="swatchStyle" />
        <!-- 原生 color input，透明覆盖在色块上 -->
        <input
          type="color"
          :value="hexFromValue"
          class="absolute inset-0 h-full w-full cursor-pointer opacity-0"
          @input="onColorPickerInput"
        />
      </div>

      <!-- 文本输入框（支持手动输入 rgba） -->
      <Input
        :model-value="modelValue"
        :placeholder="placeholder"
        class="flex-1 font-mono text-sm"
        @update:model-value="onTextInput"
      />
    </div>

    <!-- 透明度滑条（supportAlpha 或已有 rgba 值时显示） -->
    <div v-if="supportAlpha || isRgba" class="flex items-center gap-2 px-0.5">
      <span class="w-10 flex-shrink-0 text-right text-xs text-muted-foreground">透明</span>
      <div class="relative flex-1 h-4 flex items-center">
        <!-- 渐变轨道 -->
        <div
          class="absolute h-2 w-full rounded-full border border-input"
          :style="{
            background: `linear-gradient(to right, transparent 0%, ${hexFromValue} 100%)`,
          }"
        />
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          :value="alphaFromValue"
          class="relative w-full cursor-pointer appearance-none bg-transparent"
          style="accent-color: var(--primary)"
          @input="onAlphaInput"
        />
      </div>
      <span class="w-10 flex-shrink-0 text-left text-xs tabular-nums text-muted-foreground">
        {{ Math.round(alphaFromValue * 100) }}%
      </span>
    </div>
  </div>
</template>
