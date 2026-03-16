/**
 * extensions/CustomContainer.ts - Tiptap 自定义容器节点扩展
 * 职责：实现 :::tip / :::warning / :::danger / :::info 容器在可视化编辑器中的渲染与编辑
 * 渲染 HTML 结构与 markdown-it-container 输出完全一致，确保样式复用
 * 对外暴露：CustomContainer (Tiptap Node Extension)
 */

import { Node, mergeAttributes } from '@tiptap/core'

export type ContainerType = 'tip' | 'warning' | 'danger' | 'info'

const CONTAINER_TITLES: Record<ContainerType, string> = {
  tip: '提示',
  warning: '注意',
  danger: '危险',
  info: '信息',
}

export const CONTAINER_TYPES: ContainerType[] = ['tip', 'warning', 'danger', 'info']

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    customContainer: {
      /**
       * 插入或切换自定义容器
       */
      setContainer: (type: ContainerType) => ReturnType
      /**
       * 移除自定义容器（将内容提升为普通段落）
       */
      unsetContainer: () => ReturnType
    }
  }
}

export const CustomContainer = Node.create({
  name: 'customContainer',

  group: 'block',

  content: 'block+',

  defining: true,

  addAttributes() {
    return {
      type: {
        default: 'tip',
        parseHTML: (element) => {
          for (const t of CONTAINER_TYPES) {
            if (element.classList.contains(t)) return t
          }
          return 'tip'
        },
        renderHTML: (attributes) => ({
          class: `custom-container ${attributes.type}`,
        }),
      },
      title: {
        default: null,
        parseHTML: (element) => {
          const titleEl = element.querySelector('.custom-container-title')
          return titleEl?.textContent || null
        },
        // title 通过子节点渲染，不作为 DOM 属性输出
        renderHTML: () => ({}),
      },
    }
  },

  parseHTML() {
    return [
      {
        tag: 'div.custom-container',
        // contentElement 指定内容解析区域（若存在 .custom-container-content 则从中取子节点）
        contentElement: '.custom-container-content',
      },
      {
        // 兼容 markdown-it-container 原始输出（无 .custom-container-content 包裹）
        tag: 'div.custom-container',
        contentElement: 'div.custom-container',
      },
    ]
  },

  renderHTML({ node, HTMLAttributes }) {
    const type = node.attrs.type as ContainerType
    const title = node.attrs.title || CONTAINER_TITLES[type] || '提示'

    return [
      'div',
      mergeAttributes(HTMLAttributes, { class: `custom-container ${type}` }),
      ['p', { class: 'custom-container-title' }, title],
      ['div', { class: 'custom-container-content' }, 0],
    ]
  },

  addCommands() {
    return {
      setContainer:
        (type: ContainerType) =>
        ({ commands }) => {
          return commands.wrapIn(this.name, {
            type,
            title: CONTAINER_TITLES[type],
          })
        },
      unsetContainer:
        () =>
        ({ commands }) => {
          return commands.lift(this.name)
        },
    }
  },

  // Markdown 序列化/反序列化配置（供 tiptap-markdown 使用）
  addStorage() {
    return {
      markdown: {
        serialize(state: any, node: any) {
          const type = node.attrs.type || 'tip'
          const title = node.attrs.title || ''
          const defaultTitle = CONTAINER_TITLES[type as ContainerType] || ''
          // 如果 title 与默认值相同，不输出 title
          const titlePart = title && title !== defaultTitle ? ` ${title}` : ''
          state.write(`:::${type}${titlePart}\n`)
          state.renderContent(node)
          state.write(':::\n\n')
        },
        parse: {
          // tiptap-markdown 通过 HTML 解析，parseHTML 已处理
        },
      },
    }
  },
})

export { CONTAINER_TITLES }
export default CustomContainer
