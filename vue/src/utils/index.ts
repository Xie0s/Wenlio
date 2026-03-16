/*
 * index.ts - 通用工具函数入口
 * 职责：统一承载与导出纯工具函数（无业务状态、副作用）
 * 对外暴露：
 * - cn(...inputs): 合并 Tailwind className
 */

import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
