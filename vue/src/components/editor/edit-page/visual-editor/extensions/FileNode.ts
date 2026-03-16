/**
 * extensions/FileNode.ts - 自定义文件附件节点扩展
 * 职责：在 Tiptap 编辑器中表示文件附件，渲染为可操作的文件卡片，支持布局切换
 * 对外暴露：FileNode 扩展、insertFileNode / updateFileLayout commands、FileNodeAttrs 类型
 */

import { Node, mergeAttributes } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import FileNodeView from './FileNodeView.vue'

export type FileLayout = 'full' | 'half' | 'third'

export interface FileNodeAttrs {
  href: string
  fileName: string
  fileSize: number | null
  mimeType: string | null
  mediaId: string | null
  layout: FileLayout
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    fileNode: {
      insertFileNode: (attrs: Omit<FileNodeAttrs, 'layout'> & { layout?: FileLayout }) => ReturnType
      updateFileLayout: (layout: FileLayout) => ReturnType
    }
  }
}

export const FileNode = Node.create({
  name: 'fileNode',

  group: 'block',

  atom: true,

  addAttributes() {
    return {
      href: {
        default: null,
        parseHTML: (el: HTMLElement) =>
          el.getAttribute('data-href') || el.querySelector('a')?.getAttribute('href'),
        renderHTML: (attrs: Record<string, any>) =>
          attrs.href ? { 'data-href': attrs.href } : {},
      },
      fileName: {
        default: '附件',
        parseHTML: (el: HTMLElement) =>
          el.getAttribute('data-file-name')
          || el.querySelector('a')?.getAttribute('data-file-name')
          || el.querySelector('a')?.textContent
          || '附件',
        renderHTML: (attrs: Record<string, any>) =>
          attrs.fileName ? { 'data-file-name': attrs.fileName } : {},
      },
      fileSize: {
        default: null,
        parseHTML: (el: HTMLElement) => {
          const v = el.getAttribute('data-file-size')
          return v ? Number(v) : null
        },
        renderHTML: (attrs: Record<string, any>) =>
          attrs.fileSize != null ? { 'data-file-size': String(attrs.fileSize) } : {},
      },
      mimeType: {
        default: null,
        parseHTML: (el: HTMLElement) => el.getAttribute('data-mime-type'),
        renderHTML: (attrs: Record<string, any>) =>
          attrs.mimeType ? { 'data-mime-type': attrs.mimeType } : {},
      },
      mediaId: {
        default: null,
        parseHTML: (el: HTMLElement) => el.getAttribute('data-media-id'),
        renderHTML: (attrs: Record<string, any>) =>
          attrs.mediaId ? { 'data-media-id': attrs.mediaId } : {},
      },
      layout: {
        default: 'full',
        parseHTML: (el: HTMLElement) => el.getAttribute('data-layout') || 'full',
        renderHTML: (attrs: Record<string, any>) =>
          attrs.layout && attrs.layout !== 'full' ? { 'data-layout': attrs.layout } : {},
      },
    }
  },

  parseHTML() {
    return [{ tag: 'div.file-attachment' }]
  },

  renderHTML({ node, HTMLAttributes }) {
    // 编辑器内部 DOM：<div class="file-attachment" data-href="..." ...><a href="...">name</a></div>
    const name = node.attrs.fileName || '附件'
    return [
      'div',
      mergeAttributes({ class: 'file-attachment' }, HTMLAttributes),
      ['a', { href: node.attrs.href }, name],
    ]
  },

  addStorage() {
    return {
      markdown: {
        serialize(state: any, node: any) {
          const { href, fileName, fileSize, mimeType, mediaId, layout } = node.attrs as FileNodeAttrs
          const name = fileName || '附件'
          // div 上放所有 data-* 属性，内部 <a> 提供阅读端原生链接行为
          const divParts: string[] = [`class="file-attachment"`]
          if (fileSize != null) divParts.push(`data-file-size="${fileSize}"`)
          if (mimeType) divParts.push(`data-mime-type="${mimeType}"`)
          if (mediaId) divParts.push(`data-media-id="${mediaId}"`)
          if (layout && layout !== 'full') divParts.push(`data-layout="${layout}"`)

          state.ensureNewLine()
          state.write('\n')
          state.write(`<div ${divParts.join(' ')}><a href="${href}" data-file-name="${name}">${name}</a></div>`)
          state.ensureNewLine()
          state.write('\n')
        },
      },
    }
  },

  addNodeView() {
    return VueNodeViewRenderer(FileNodeView)
  },

  addCommands() {
    return {
      insertFileNode:
        (attrs) =>
        ({ commands }) => {
          return commands.insertContent({
            type: this.name,
            attrs: { layout: 'full', ...attrs },
          })
        },
      updateFileLayout:
        (layout: FileLayout) =>
        ({ commands }) => {
          return commands.updateAttributes(this.name, { layout })
        },
    }
  },
})
