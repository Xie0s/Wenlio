/**
 * lib/markdown.ts - Markdown 渲染引擎
 * 职责：
 * 1) 初始化 markdown-it + shiki 代码高亮 + 自定义容器插件
 * 2) 使用 DOMPurify 对渲染输出进行 XSS 净化
 * 3) 对外暴露异步渲染函数，供阅读页和编辑器预览复用
 * 对外暴露：
 * - renderMarkdown(source: string): Promise<string>
 * - renderMarkdownSync(source: string): string（初始化完成后可用）
 */

import MarkdownIt from 'markdown-it'
import markdownItAnchor from 'markdown-it-anchor'
import markdownItContainer from 'markdown-it-container'
import markdownItTaskLists from 'markdown-it-task-lists'
import DOMPurify from 'dompurify'

let md: MarkdownIt | null = null
let initPromise: Promise<void> | null = null

// DOMPurify 净化配置（renderMarkdown / renderMarkdownSync 共用）
// 安全策略：
//   - 不允许未受限的 iframe（移除 ADD_TAGS 中的 iframe）
//   - input 仅允许 checkbox 类型的 task list
//   - 禁止 form / script / style / object / embed 等危险标签（DOMPurify 默认行为）
const PURIFY_OPTIONS = {
  ADD_TAGS: ['input'],
  ADD_ATTR: [
    'target', 'rel', 'class', 'style', 'type', 'checked', 'disabled', 'width',
    // 图片节点自定义属性
    'data-align', 'data-rounded', 'data-media-id',
    // 表格单元格背景色
    'data-background-color',
    // 文件附件节点属性
    'data-file-size', 'data-mime-type', 'data-layout', 'data-file-name',
  ],
  FORBID_TAGS: ['iframe', 'form', 'object', 'embed'],
}

// 自定义容器类型：tip / warning / danger / info
const CONTAINER_TYPES = ['tip', 'warning', 'danger', 'info'] as const
type ContainerType = (typeof CONTAINER_TYPES)[number]

const CONTAINER_TITLES: Record<ContainerType, string> = {
  tip: '提示',
  warning: '注意',
  danger: '危险',
  info: '信息',
}

function createContainerPlugin(md: MarkdownIt, type: ContainerType) {
  markdownItContainer(md, type, {
    validate(params: string) {
      return params.trim().split(/\s+/, 2)[0] === type
    },
    render(tokens: any[], idx: number) {
      if (tokens[idx].nesting === 1) {
        const customTitle = tokens[idx].info.trim().slice(type.length).trim()
        const title = customTitle || CONTAINER_TITLES[type]
        return `<div class="custom-container ${type}"><p class="custom-container-title">${md.utils.escapeHtml(title)}</p>\n`
      }
      return '</div>\n'
    },
  })
}

async function initMarkdownIt(): Promise<MarkdownIt> {
  if (md) return md

  md = new MarkdownIt({
    html: true,
    linkify: true,
    typographer: true,
  })

  // 注册锚点插件
  md.use(markdownItAnchor, {
    permalink: markdownItAnchor.permalink.linkInsideHeader({
      placement: 'after',
      symbol: '#',
      class: 'header-anchor',
    }),
    slugify: (s: string) =>
      s
        .trim()
        .toLowerCase()
        .replace(/\s+/g, '-')
        .replace(/[^\w\u4e00-\u9fff-]/g, ''),
  })

  // 给表格包裹可横向滚动的容器，防止宽表格撑破页面布局
  const defaultTableOpen = md.renderer.rules.table_open
    ?? ((tokens: any[], idx: number, options: any, _env: any, self: any) => self.renderToken(tokens, idx, options))
  const defaultTableClose = md.renderer.rules.table_close
    ?? ((tokens: any[], idx: number, options: any, _env: any, self: any) => self.renderToken(tokens, idx, options))
  md.renderer.rules.table_open = (tokens, idx, options, env, self) =>
    `<div class="table-wrapper">${defaultTableOpen(tokens, idx, options, env, self)}`
  md.renderer.rules.table_close = (tokens, idx, options, env, self) =>
    `${defaultTableClose(tokens, idx, options, env, self)}</div>`

  // 注册任务列表（- [ ] / - [x] 语法）
  md.use(markdownItTaskLists, { enabled: true, label: true })

  // 注册自定义容器
  for (const type of CONTAINER_TYPES) {
    createContainerPlugin(md, type)
  }

  // 注册 Shiki 代码高亮（异步加载）
  try {
    const { fromHighlighter } = await import('@shikijs/markdown-it')
    const { createHighlighter } = await import('shiki')
    const highlighter = await createHighlighter({
      themes: ['github-light', 'github-dark'],
      langs: [
        'javascript', 'typescript', 'vue', 'html', 'css', 'json',
        'bash', 'shell', 'yaml', 'markdown', 'go', 'python',
        'java', 'sql', 'docker', 'nginx', 'xml', 'diff',
        'makefile', 'toml', 'rust', 'cpp', 'c', 'csharp',
        'kotlin', 'swift', 'ruby', 'php', 'powershell',
        'scss', 'less', 'graphql', 'proto', 'ini',
      ],
    })

    // 保存原始 highlight 函数，对未知语言做降级处理
    const shikiPlugin = fromHighlighter(highlighter, {
      themes: {
        light: 'github-light',
        dark: 'github-dark',
      },
      defaultColor: false,
    })
    md.use(shikiPlugin)

    // 包装 highlight，捕获未加载语言的错误，降级为纯文本展示
    const originalHighlight = md.options.highlight
    if (originalHighlight) {
      md.options.highlight = (code, lang, attrs) => {
        try {
          return originalHighlight(code, lang, attrs)
        } catch (e) {
          console.warn(`[markdown] Shiki 不支持语言 "${lang}"，降级为纯文本展示`)
          return `<pre><code>${md!.utils.escapeHtml(code)}</code></pre>`
        }
      }
    }
  } catch (e) {
    console.warn('[markdown] Shiki 加载失败，使用默认代码展示:', e)
  }

  return md
}

/**
 * 异步渲染 Markdown 为安全 HTML
 * 首次调用会初始化 markdown-it + shiki，后续调用直接渲染
 */
export async function renderMarkdown(source: string): Promise<string> {
  if (!initPromise) {
    initPromise = initMarkdownIt().then(() => {})
  }
  await initPromise
  const rawHtml = md!.render(source)
  return DOMPurify.sanitize(rawHtml, PURIFY_OPTIONS)
}

/**
 * 同步渲染（仅在 initPromise 已 resolved 后可用，否则降级为纯文本）
 * 适用于编辑器预览等高频场景
 */
export function renderMarkdownSync(source: string): string {
  if (!md) return source
  const rawHtml = md.render(source)
  return DOMPurify.sanitize(rawHtml, PURIFY_OPTIONS)
}
