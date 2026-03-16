/**
 * composables/index.ts - 组合式函数入口
 * 职责：统一导出所有 composable 函数，提供可复用的组合式逻辑
 * 对外暴露：所有 composable 函数的聚合导出
 */
export { useOutsideClick } from './useOutsideClick'
export { useAutoSave } from './useAutoSave'
export { useMarkdownPreview } from './useMarkdownPreview'
export { useThemeEditor, useEditorInject, EDITOR_KEY } from './useThemeEditor'
export { useAdminBreadcrumbs, ADMIN_NAME_MAP, isAdminObjectId } from './useAdminBreadcrumbs'
export type { BreadcrumbItemData } from './useAdminBreadcrumbs'
export { useDocReader } from './useDocReader'
