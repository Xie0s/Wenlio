/**
 * extensions/index.ts - Tiptap 可视化编辑器扩展配置集合
 * 职责：统一导出所有 Tiptap 扩展，供 useVisualEditor 消费
 * 对外暴露：createExtensions(options) → Extension[]
 */

import StarterKit from '@tiptap/starter-kit'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import ImageBase from '@tiptap/extension-image'
import { Table } from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableCellBase from '@tiptap/extension-table-cell'
import TableHeaderBase from '@tiptap/extension-table-header'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import Placeholder from '@tiptap/extension-placeholder'
import { Markdown } from 'tiptap-markdown'
import { Extension } from '@tiptap/core'
import { Plugin } from '@tiptap/pm/state'
import { Decoration, DecorationSet } from '@tiptap/pm/view'
import { DOMSerializer } from '@tiptap/pm/model'
import { CustomContainer } from './CustomContainer'
import { FileNode } from './FileNode'

function escapeHtml(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

/**
 * 将表格单元格内容序列化为 HTML 字符串
 * 使用 ProseMirror DOMSerializer 而非 markdown renderContent，
 * 防止 markdown 特殊字符在 HTML 表格内被转义导致保存/加载循环中指数级膨胀
 */
function serializeCellToHtml(cell: any): string {
  const serializer = DOMSerializer.fromSchema(cell.type.schema)
  const fragment = serializer.serializeFragment(cell.content)
  const wrapper = document.createElement('div')
  wrapper.appendChild(fragment)
  if (wrapper.childElementCount === 1 && wrapper.firstElementChild?.tagName === 'P') {
    return wrapper.firstElementChild.innerHTML
  }
  return wrapper.innerHTML
}

/**
 * 自定义 Table 扩展 — 统一序列化为 HTML，保留单元格属性（如 backgroundColor）
 */
const CustomTable = Table.extend({
  addStorage() {
    return {
      markdown: {
        serialize(state: any, node: any) {
          state.ensureNewLine()
          state.write('<table>\n')

          node.forEach((row: any) => {
            state.write('  <tr>\n')

            row.forEach((cell: any) => {
              const isHeader = cell.type.name === 'tableHeader'
              const tag = isHeader ? 'th' : 'td'
              const attrs = cell.attrs || {}
              const attrParts: string[] = []

              if (attrs.colspan && attrs.colspan !== 1) {
                attrParts.push(`colspan="${escapeHtml(String(attrs.colspan))}"`)
              }
              if (attrs.rowspan && attrs.rowspan !== 1) {
                attrParts.push(`rowspan="${escapeHtml(String(attrs.rowspan))}"`)
              }
              if (Array.isArray(attrs.colwidth) && attrs.colwidth.length > 0) {
                attrParts.push(`data-colwidth="${escapeHtml(attrs.colwidth.join(','))}"`)
              }
              if (attrs.backgroundColor) {
                const bg = escapeHtml(String(attrs.backgroundColor))
                attrParts.push(`data-background-color="${bg}"`)
                attrParts.push(`style="background-color: ${bg}"`)
              }

              const attrText = attrParts.length > 0 ? ` ${attrParts.join(' ')}` : ''
              const cellHtml = serializeCellToHtml(cell)
              state.write(`    <${tag}${attrText}>${cellHtml}</${tag}>\n`)
            })

            state.write('  </tr>\n')
          })

          state.write('</table>\n\n')
        },
      },
    }
  },
})

/**
 * 统一粘贴增强插件，注册在 Markdown 扩展之前（ProseMirror 按顺序调用 handlePaste）。
 *
 * 处理以下场景（按优先级）：
 *
 * 1. URL 粘贴 + 选区非空 → 将选中文字变为超链接，不插入 URL 文本
 * 2. Markdown 从代码编辑器粘贴（剪贴板同时含 HTML）
 *    tiptap-markdown 在检测到 HTML 时直接 return false 跳过 markdown 解析，
 *    此处接管：用 tiptap-markdown 内部 parser 将纯文本解析为 ProseMirror 节点插入。
 */
const MarkdownPasteOverride = Extension.create({
  name: 'markdownPasteOverride',
  addProseMirrorPlugins() {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const ext = this
    return [
      new Plugin({
        props: {
          handlePaste(view, event) {
            const { clipboardData } = event
            if (!clipboardData) return false

            const html = clipboardData.getData('text/html')
            const text = clipboardData.getData('text/plain').trim()
            if (!text) return false

            // ── 场景 1：URL 粘贴覆盖选区 → 自动创建链接 ──────────────────────────
            const isUrl = /^https?:\/\/[^\s]+$/.test(text)
            if (isUrl && !view.state.selection.empty) {
              event.preventDefault()
              ext.editor.chain().focus().setLink({ href: text }).run()
              return true
            }

            // ── 场景 2：Markdown 从代码编辑器粘贴（HTML + markdown 特征）──────────
            // 只在同时有 HTML 时介入（纯文本粘贴由 tiptap-markdown 自身处理）
            if (!html) return false

            const looksLikeMarkdown =
              /^#{1,6} /m.test(text) || /^```/m.test(text) || /^\|.+\|/m.test(text)
            if (!looksLikeMarkdown) return false

            const parser = (ext.editor.storage as any)?.markdown?.parser
            if (!parser) return false

            event.preventDefault()
            try {
              const parsed = parser.parse(text)

              // 兼容不同版本 parser 返回：string(HTML) / Slice / Node
              if (typeof parsed === 'string') {
                ext.editor.chain().focus().insertContent(parsed).run()
                return true
              }

              if (parsed && typeof parsed === 'object') {
                const maybeSlice = parsed as any
                if ('content' in maybeSlice && 'openStart' in maybeSlice && 'openEnd' in maybeSlice) {
                  view.dispatch(view.state.tr.replaceSelection(maybeSlice).scrollIntoView())
                  return true
                }

                if (typeof (parsed as any).slice === 'function') {
                  const slice = (parsed as any).slice(0)
                  if (slice && typeof slice === 'object' && 'content' in slice) {
                    view.dispatch(view.state.tr.replaceSelection(slice).scrollIntoView())
                    return true
                  }
                }
              }

              return false
            } catch (e) {
              console.warn('[MarkdownPaste] 解析失败，回退默认处理', e)
              return false
            }
          },
        },
      }),
    ]
  },
})

/**
 * 与 markdown-it-anchor 保持一致的 slugify 实现，确保两套渲染下标题 id 相同
 */
function slugifyHeading(text: string): string {
  return text
    .trim()
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^\w\u4e00-\u9fff-]/g, '')
}

/**
 * 通过 ProseMirror Decoration 给所有标题节点附加 id 属性
 * 不修改文档数据模型，仅影响渲染 DOM，确保 useToc 的 h2[id]/h3[id]/h4[id] 查询可命中
 */
const HeadingIdPlugin = Extension.create({
  name: 'headingId',
  addProseMirrorPlugins() {
    return [
      new Plugin({
        props: {
          decorations(state) {
            const decorations: Decoration[] = []
            state.doc.descendants((node, pos) => {
              if (node.type.name === 'heading') {
                const id = slugifyHeading(node.textContent)
                if (id) {
                  decorations.push(
                    Decoration.node(pos, pos + node.nodeSize, { id }),
                  )
                }
              }
            })
            return DecorationSet.create(state.doc, decorations)
          },
        },
      }),
    ]
  },
})

/**
 * 自定义 TableHeader 扩展 — 支持背景色
 */
const CustomTableHeader = TableHeaderBase.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      backgroundColor: {
        default: null,
        parseHTML: (element: HTMLElement) =>
          element.getAttribute('data-background-color') || element.style.backgroundColor || null,
        renderHTML: (attributes: Record<string, any>) => {
          if (!attributes.backgroundColor) return {}
          return {
            'data-background-color': attributes.backgroundColor,
            style: `background-color: ${attributes.backgroundColor}`,
          }
        },
      },
    }
  },
})

/**
 * 自定义 TableCell 扩展 — 支持背景色
 */
const CustomTableCell = TableCellBase.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      backgroundColor: {
        default: null,
        parseHTML: (element: HTMLElement) =>
          element.getAttribute('data-background-color') || element.style.backgroundColor || null,
        renderHTML: (attributes: Record<string, any>) => {
          if (!attributes.backgroundColor) return {}
          return {
            'data-background-color': attributes.backgroundColor,
            style: `background-color: ${attributes.backgroundColor}`,
          }
        },
      },
    }
  },
})

/**
 * 自定义 Image 扩展 — 支持 width 和 data-align 属性
 */
const CustomImage = ImageBase.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      width: {
        default: null,
        parseHTML: (element: HTMLElement) => element.getAttribute('width') || element.style.width || null,
        renderHTML: (attributes: Record<string, any>) => {
          if (!attributes.width) return {}
          return { style: `width: ${attributes.width};` }
        },
      },
      'data-align': {
        default: null,
        parseHTML: (element: HTMLElement) => element.getAttribute('data-align') || null,
        renderHTML: (attributes: Record<string, any>) => {
          const align = attributes['data-align']
          if (!align) return {}
          return { 'data-align': align }
        },
      },
      'data-rounded': {
        default: null,
        parseHTML: (element: HTMLElement) => element.getAttribute('data-rounded') || null,
        renderHTML: (attributes: Record<string, any>) => {
          const rounded = attributes['data-rounded']
          if (!rounded) return {}
          return { 'data-rounded': rounded }
        },
      },
      'data-media-id': {
        default: null,
        parseHTML: (element: HTMLElement) => element.getAttribute('data-media-id') || null,
        renderHTML: (attributes: Record<string, any>) => {
          const mediaId = attributes['data-media-id']
          if (!mediaId) return {}
          return { 'data-media-id': mediaId }
        },
      },
    }
  },
  addProseMirrorPlugins() {
    return [
      new Plugin({
        props: {
          decorations(state) {
            const { selection, doc } = state
            const decorations: Decoration[] = []
            doc.descendants((node, pos) => {
              if (node.type.name !== 'image') return
              if (selection.from <= pos && selection.to >= pos + node.nodeSize) {
                decorations.push(Decoration.node(pos, pos + node.nodeSize, { class: 'image-in-selection' }))
              }
            })
            return DecorationSet.create(doc, decorations)
          },
        },
      }),
    ]
  },
  addStorage() {
    return {
      markdown: {
        serialize(state: any, node: any) {
          const attrs = node.attrs as { src: string; alt?: string; width?: string; 'data-align'?: string; 'data-rounded'?: string; 'data-media-id'?: string }
          const { src, alt, width } = attrs
          const dataAlign = attrs['data-align']
          const dataRounded = attrs['data-rounded']
          const dataMediaId = attrs['data-media-id']
          // 有自定义属性时序列化为 HTML img 标签，保留 width、data-align、data-rounded、data-media-id
          if (width || dataAlign || dataRounded || dataMediaId) {
            const widthAttr = width ? ` style="width: ${width};"` : ''
            const alignAttr = dataAlign ? ` data-align="${dataAlign}"` : ''
            const roundedAttr = dataRounded ? ` data-rounded="${dataRounded}"` : ''
            const mediaIdAttr = dataMediaId ? ` data-media-id="${dataMediaId}"` : ''
            state.ensureNewLine()
            state.write(`<img src="${src}" alt="${alt || ''}"${widthAttr}${alignAttr}${roundedAttr}${mediaIdAttr} />`)
            state.ensureNewLine()
          } else {
            // 无自定义属性时使用标准 markdown 格式（block 节点同样需要换行隔离）
            state.ensureNewLine()
            state.write(`![${alt || ''}](${src})`)
            state.ensureNewLine()
          }
        },
        parse: {
          // 解析由 markdown-it 处理
        },
      },
    }
  },
})

export interface ExtensionOptions {
  placeholder?: string
}

/**
 * 创建完整的 Tiptap 扩展集合
 * StarterKit v3 已包含：Blockquote, Bold, BulletList, Code, CodeBlock, Document,
 *   Dropcursor, Gapcursor, HardBreak, Heading, History, HorizontalRule,
 *   Italic, Link, ListItem, ListKeymap, OrderedList, Paragraph, Strike, Text, Underline
 *   （v3.0.1+ 新增 Link / Underline / ListKeymap，Link 通过 StarterKit.configure 配置）
 */
export function createExtensions(options: ExtensionOptions = {}) {
  return [
    StarterKit.configure({
      heading: { levels: [1, 2, 3, 4] },
      codeBlock: false,
      // StarterKit v3 已内置 Link，在此配置以避免与 tiptap-markdown 内置 link 重复注册
      link: {
        openOnClick: false,
        autolink: true,
        HTMLAttributes: { target: '_blank', rel: 'noopener noreferrer' },
      },
    }),

    // 代码块（lowlight 语法高亮，替代 StarterKit 默认 codeBlock）
    CodeBlockLowlight.configure({
      lowlight: createLowlight(common),
      HTMLAttributes: { class: 'shiki' },
    }),

    // 图片（自定义扩展，支持 width 和对齐）
    CustomImage.configure({
      inline: false,
      allowBase64: false,
    }),

    // 表格（自定义扩展，支持背景色）
    CustomTable.configure({ resizable: false }),
    TableRow,
    CustomTableCell,
    CustomTableHeader,

    // 任务列表
    TaskList,
    TaskItem.configure({
      nested: true,
      HTMLAttributes: {
        style: 'display: flex; align-items: flex-start;',
      },
    }),

    // 占位文本
    Placeholder.configure({
      placeholder: options.placeholder || '开始编写内容...',
    }),

    // 优先注册粘贴修复插件（必须在 Markdown 之前，确保 handlePaste 先于 tiptap-markdown 运行）
    MarkdownPasteOverride,

    // Markdown 双向转换
    // 过滤 tiptap-markdown 内置扩展：
    // link  → 避免与 StarterKit v3 已内置的 link 重名
    // image → 避免与我们的 CustomImage 重名
    // table → 强制表格序列化为 HTML，保留单元格背景色等自定义属性
    Markdown.extend({
      addExtensions() {
        return (this.parent?.() ?? []).filter(
          (ext: any) => !['link', 'image', 'italic', 'table'].includes(ext.name),
        )
      },
    }).configure({
      html: true,              // 允许 HTML（自定义容器使用 HTML 中转）
      tightLists: true,        // 紧凑列表
      bulletListMarker: '-',   // 无序列表标记
      linkify: true,           // 自动链接识别
      breaks: false,           // 不使用 soft break
      transformPastedText: true,   // 粘贴 markdown 文本自动转换
      transformCopiedText: true,   // 复制时保留 markdown 格式
    }),

    // 标题自动添加 id 属性（供 TOC 导航使用）
    HeadingIdPlugin,

    // 自定义容器（:::tip / :::warning / :::danger / :::info）
    CustomContainer,

    // 文件附件节点（支持布局切换的文件卡片）
    FileNode,
  ]
}
