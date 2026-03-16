/**
 * components/doc-editor/types.ts - 文档编辑器共享类型
 * 职责：定义编辑器组件树内部共享的类型，消除各组件重复定义
 * 对外暴露：SectionWithPages
 */

import type { Section, DocPage } from '@/utils/types'

export type SectionWithPages = Section & { pages: DocPage[] }
