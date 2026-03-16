/**
 * env.d.ts - 全局模块类型声明
 * 职责：为没有自带 TypeScript 类型的第三方包补充类型声明
 */

/// <reference types="vite/client" />

declare module 'markdown-it-task-lists' {
  import type MarkdownIt from 'markdown-it'
  interface TaskListOptions {
    /** 是否允许用户点击 checkbox（默认 false = disabled） */
    enabled?: boolean
    /** 是否将 checkbox 包裹在 label 标签中 */
    label?: boolean
    /** label 在 checkbox 之后 */
    labelAfter?: boolean
  }
  function markdownItTaskLists(md: MarkdownIt, options?: TaskListOptions): void
  export = markdownItTaskLists
}
