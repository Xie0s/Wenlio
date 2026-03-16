<!--
  TableInsertDialog.vue - 表格插入对话框
  职责：提供可视化网格选择器，让用户选择表格行列数和是否包含表头行
  对外暴露：Props: open
            Emits: update:open, insert(rows, cols, withHeaderRow)
-->
<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { X, Check } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'

const props = defineProps<{ open: boolean }>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  insert: [rows: number, cols: number, withHeaderRow: boolean]
}>()

const MAX_ROWS = 10
const MAX_COLS = 10

const hoveredCell = ref<{ row: number; col: number } | null>(null)
const selectedSize = ref({ rows: 3, cols: 3 })
const withHeaderRow = ref(true)

function blurActiveElement() {
  const activeEl = document.activeElement
  if (activeEl instanceof HTMLElement) activeEl.blur()
}

// Dialog 打开时确保当前焦点脱离编辑器，避免 aria-hidden 警告
watch(() => props.open, (open) => {
  if (open) blurActiveElement()
})

function handleCellHover(row: number, col: number) {
  hoveredCell.value = { row, col }
}

function handleCellClick(row: number, col: number) {
  selectedSize.value = { rows: row + 1, cols: col + 1 }
}

function handleGridLeave() {
  hoveredCell.value = null
}

const sizeText = computed(() => {
  if (hoveredCell.value) {
    return `${hoveredCell.value.row + 1} × ${hoveredCell.value.col + 1}`
  }
  return `${selectedSize.value.rows} × ${selectedSize.value.cols}`
})

// 计算每个单元格是否激活
function isCellActive(rowIndex: number, colIndex: number): boolean {
  const activeRows = hoveredCell.value ? hoveredCell.value.row + 1 : selectedSize.value.rows
  const activeCols = hoveredCell.value ? hoveredCell.value.col + 1 : selectedSize.value.cols
  return rowIndex < activeRows && colIndex < activeCols
}

function handleInsert() {
  emit('insert', selectedSize.value.rows, selectedSize.value.cols, withHeaderRow.value)
  emit('update:open', false)
  // 重置状态
  selectedSize.value = { rows: 3, cols: 3 }
  hoveredCell.value = null
}

function handleCancel() {
  emit('update:open', false)
  selectedSize.value = { rows: 3, cols: 3 }
  hoveredCell.value = null
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-[360px] rounded-3xl">
      <DialogHeader>
        <DialogTitle>插入表格</DialogTitle>
        <DialogDescription>
          选择表格的行数和列数，点击网格进行选择
        </DialogDescription>
      </DialogHeader>

      <!-- 网格选择器 -->
      <div class="py-2">
        <div class="flex justify-center">
          <div
            class="inline-grid gap-0.5 p-1.5 rounded-lg bg-muted/50"
            :style="{ gridTemplateColumns: `repeat(${MAX_COLS}, 1fr)` }"
            @mouseleave="handleGridLeave"
          >
            <template v-for="rowIndex in MAX_ROWS" :key="`row-${rowIndex}`">
              <button
                v-for="colIndex in MAX_COLS"
                :key="`${rowIndex}-${colIndex}`"
                type="button"
                class="w-7 h-7 rounded transition-colors duration-75"
                :class="isCellActive(rowIndex - 1, colIndex - 1) ? 'bg-primary' : 'bg-background hover:bg-accent'"
                :tabindex="-1"
                @mouseenter="handleCellHover(rowIndex - 1, colIndex - 1)"
                @click="handleCellClick(rowIndex - 1, colIndex - 1)"
              />
            </template>
          </div>
        </div>

        <div class="text-center mt-2 text-sm font-medium text-foreground">
          {{ sizeText }}
        </div>
      </div>

      <!-- 底部：表头选项 + 按钮 -->
      <div class="flex items-center justify-between pt-2">
        <!-- 左侧：表头选项 -->
        <div class="flex items-center gap-2">
          <Switch
            id="header-row"
            :checked="withHeaderRow"
            @update:checked="withHeaderRow = $event"
          />
          <Label for="header-row" class="text-sm cursor-pointer">
            包含表头行
          </Label>
        </div>

        <!-- 右侧：操作按钮 -->
        <div class="flex items-center gap-3">
          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                type="button"
                size="icon"
                variant="outline"
                class="rounded-full"
                @click="handleCancel"
              >
                <X class="h-5 w-5" />
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
                @click="handleInsert"
              >
                <Check class="h-5 w-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>插入</TooltipContent>
          </Tooltip>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
