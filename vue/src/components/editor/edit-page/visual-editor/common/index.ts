/**
 * common/index.ts - 可视化编辑器公共组件统一导出
 * 职责：导出链接、图片、表格、代码块、文件附件相关的浮动菜单和对话框组件
 * 对外暴露：LinkMenu, LinkInputDialog, LinkInputPopover, ImageMenu, TableMenu, TableInsertDialog, CodeBlockMenu, FileCard, FileCardMenu, FilePickerDialog
 */

export { default as LinkMenu } from './LinkMenu.vue'
export { default as LinkInputDialog } from './LinkInputDialog.vue'
export { default as LinkInputPopover } from './LinkInputPopover.vue'
export { default as ImageMenu } from './ImageMenu.vue'
export { default as TableMenu } from './TableMenu.vue'
export { default as TableInsertDialog } from './TableInsertDialog.vue'
export { default as CodeBlockMenu } from './CodeBlockMenu.vue'
export { default as FileCard } from './FileCard.vue'
export { default as FileCardMenu } from './FileCardMenu.vue'
export { default as FilePickerDialog } from './FilePickerDialog.vue'
export { cellBackgroundColors, canMergeCells, canSplitCell } from './table-helpers'
