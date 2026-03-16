/**
 * table-helpers.ts - 表格操作辅助函数
 * 职责：提供表格浮动菜单所需的预设数据和工具函数（背景色、合并/拆分检测）
 * 对外暴露：cellBackgroundColors, canMergeCells, canSplitCell
 */

import type { Editor } from '@tiptap/vue-3'

/**
 * 单元格背景色预设
 * 使用 HSL 格式，亮暗主题下均有良好可读性
 */
export const cellBackgroundColors = [
  { name: '无', value: null },
  { name: '灰色', value: 'hsl(220 9% 46% / 0.15)' },
  { name: '红色', value: 'hsl(0 84% 60% / 0.15)' },
  { name: '橙色', value: 'hsl(25 95% 53% / 0.15)' },
  { name: '黄色', value: 'hsl(48 96% 53% / 0.15)' },
  { name: '绿色', value: 'hsl(142 71% 45% / 0.15)' },
  { name: '蓝色', value: 'hsl(217 91% 60% / 0.15)' },
  { name: '紫色', value: 'hsl(271 81% 56% / 0.15)' },
  { name: '粉色', value: 'hsl(330 81% 60% / 0.15)' },
]

/**
 * 检查当前选区是否可以合并单元格
 */
export function canMergeCells(editor: Editor): boolean {
  return editor.can().mergeCells()
}

/**
 * 检查当前选区是否可以拆分单元格
 */
export function canSplitCell(editor: Editor): boolean {
  return editor.can().splitCell()
}
