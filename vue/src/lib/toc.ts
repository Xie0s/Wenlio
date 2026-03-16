/**
 * lib/toc.ts - TOC（目录）相关纯工具函数
 * 职责：提供 TOC 渲染所需的无副作用工具函数，供编辑器侧边栏与阅读页 TOC 组件共用
 * 对外暴露：levelPaddingClass(level)
 */

/**
 * 根据标题层级返回对应的左缩进 Tailwind class
 * h2 → pl-2, h3 → pl-4, h4+ → pl-8
 */
export function levelPaddingClass(level: number): string {
  if (level >= 4) return 'pl-8'
  if (level === 3) return 'pl-4'
  return 'pl-2'
}
