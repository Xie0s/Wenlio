/**
 * lib/prose.ts - Markdown 内容渲染区 prose 样式常量
 * 职责：统一定义 prose Tailwind 类字符串，确保编辑器预览区与阅读页渲染样式完全一致
 * 对外暴露：PROSE_CLASSES
 */

/**
 * 所有 Markdown 内容渲染容器共用的 Tailwind prose 类
 * 编辑器源码预览、可视化编辑器、文档阅读页均应引用此常量
 */
export const PROSE_CLASSES =
  'prose prose-slate dark:prose-invert max-w-none ' +
  'prose-headings:font-semibold prose-headings:tracking-tight ' +
  'prose-a:text-primary prose-a:no-underline hover:prose-a:underline ' +
  'prose-code:text-foreground prose-pre:p-0 prose-pre:bg-transparent'
