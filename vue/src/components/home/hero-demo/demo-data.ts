/**
 * hero-demo/demo-data.ts
 * 演示组件所用的静态模拟数据：文档树、页面内容
 */
import type { SectionTree } from '@/utils/types'

// ── 查看场景文档树 ──────────────────────────────────────────
export const VIEWER_TREE: SectionTree[] = [
  {
    id: 's1', title: '入门指南', sort_order: 1,
    pages: [
      { id: 'p1', slug: 'intro',       title: '介绍',     sort_order: 1 },
      { id: 'p2', slug: 'quickstart',  title: '快速开始', sort_order: 2 },
      { id: 'p3', slug: 'install',     title: '安装配置', sort_order: 3 },
    ],
  },
  {
    id: 's2', title: 'API 参考', sort_order: 2,
    pages: [
      { id: 'p4', slug: 'api-themes',  title: '主题 API',  sort_order: 1 },
      { id: 'p5', slug: 'api-pages',   title: '文档页 API', sort_order: 2 },
    ],
  },
  {
    id: 's3', title: '进阶', sort_order: 3,
    pages: [
      { id: 'p6', slug: 'faq', title: '常见问题', sort_order: 1 },
    ],
  },
]

// ── 编辑场景文档树 ──────────────────────────────────────────
export const EDITOR_TREE: SectionTree[] = [
  {
    id: 's1', title: '入门指南', sort_order: 1,
    pages: [
      { id: 'p1', slug: 'intro',      title: '介绍',     sort_order: 1 },
      { id: 'p2', slug: 'quickstart', title: '快速开始', sort_order: 2 },
      { id: 'p3', slug: 'install',    title: '安装配置', sort_order: 3 },
    ],
  },
  {
    id: 's2', title: 'API 参考', sort_order: 2,
    pages: [
      { id: 'p4', slug: 'api-themes', title: '主题 API', sort_order: 1 },
    ],
  },
]

// ── 查看场景页面内容（Markdown） ────────────────────────────
export const VIEWER_PAGES: Record<string, { title: string; markdown: string }> = {
  p1: {
    title: '介绍',
    markdown: `# 介绍

文档管理平台为团队提供**多租户文档托管**与阅读服务，支持主题划分、版本管理与全文搜索。

## 核心能力

- **多租户隔离**：每个租户独立管理文档，互不干扰
- **主题与版本**：支持多主题、多版本文档管理
- **Markdown 编辑**：源码与可视化双模式，30 秒自动保存
- **即时发布**：发布仅变更状态，无需等待构建

## 适用场景

- 产品帮助文档
- 开发者 API 参考
- 团队内部知识库
`,
  },
  p2: {
    title: '快速开始',
    markdown: `# 快速开始

欢迎使用文档管理平台。本指南将帮助您在 5 分钟内完成首个文档站点的搭建。

## 前置要求

在开始之前，请确认您已具备以下条件：

- 已注册账号并获得管理员权限
- 熟悉基本的 Markdown 语法

## 创建主题

主题是文档站点的容器，您可以为每个产品线创建独立主题。

1. 进入后台管理 → **主题列表**
2. 点击右上角「新建主题」
3. 填写主题名称和 URL Slug

::: tip
Slug 仅支持英文小写字母、数字和连字符，创建后不可修改。
:::

## 编写内容

平台支持 **Markdown 源码** 与可视化富文本双模式编辑，每 30 秒自动保存草稿。

\`\`\`bash
# 发布整个版本下所有草稿
POST /api/versions/:id/publish
\`\`\`

## 发布版本

完成内容编写后，在工具栏点击「发布」即可上线，读者刷新页面即见最新内容。
`,
  },
  p3: {
    title: '安装配置',
    markdown: `# 安装配置

## 环境要求

- Go 1.21+
- PostgreSQL 15+
- Node.js 20+

## 配置文件

\`\`\`yaml
# config.yaml
database:
  host: localhost
  port: 5432
  name: docmather
server:
  port: 8080
  jwt_secret: your-secret-key
\`\`\`

## 启动服务

\`\`\`bash
go run ./cmd/server
\`\`\`
`,
  },
  p4: {
    title: '主题 API',
    markdown: `# 主题 API

## 获取主题列表

\`\`\`bash
GET /api/themes
Authorization: Bearer {token}
\`\`\`

**响应示例：**

\`\`\`json
{
  "items": [
    { "id": "theme-1", "name": "用户手册", "slug": "user-guide" }
  ]
}
\`\`\`

## 创建主题

\`\`\`bash
POST /api/themes
Content-Type: application/json

{
  "name": "API 文档",
  "slug": "api-docs"
}
\`\`\`
`,
  },
  p5: {
    title: '文档页 API',
    markdown: `# 文档页 API

## 获取页面内容

\`\`\`bash
GET /api/pages/:id
Authorization: Bearer {token}
\`\`\`

## 更新页面

\`\`\`bash
PUT /api/pages/:id
Content-Type: application/json

{
  "title": "快速开始",
  "content": "# 快速开始\\n\\n..."
}
\`\`\`
`,
  },
  p6: {
    title: '常见问题',
    markdown: `# 常见问题

## 如何修改主题 Slug？

Slug 创建后**不可修改**，如需更改请重新创建主题并迁移文档内容。

## 发布后读者看不到更新？

请确认版本状态已切换为「已发布」，读者刷新页面即可看到最新内容，无需构建。

## 是否支持私有文档？

当前版本所有已发布文档均公开可访问，权限控制功能在路线图中。
`,
  },
}

// ── 编辑场景页面内容（可编辑 Markdown） ────────────────────
export const EDITOR_PAGES: Record<string, { title: string; content: string }> = {
  p1: {
    title: '介绍',
    content: `# 介绍

文档管理平台为团队提供多租户文档托管与阅读服务。

## 核心能力

- **多租户隔离**：每个租户独立管理文档，互不干扰
- **主题与版本**：支持多主题、多版本文档管理
- **Markdown 编辑**：源码与可视化双模式，30 秒自动保存
- **即时发布**：发布仅变更状态，无需等待构建
`,
  },
  p2: {
    title: '快速开始',
    content: `# 快速开始

本指南将帮助您在 5 分钟内完成首个文档站点的搭建。

## 创建主题

进入后台 → **主题列表** → 点击「新建主题」，填写名称与 Slug。

::: tip
Slug 仅支持英文小写字母、数字和连字符。
:::

## 编写内容

1. 在主题下创建版本
2. 添加章节与文档页
3. 在编辑器中填写 Markdown 内容

## 发布上线

点击顶部「发布」按钮，读者刷新页面即可看到最新内容。
`,
  },
  p3: {
    title: '安装配置',
    content: `# 安装配置

## 环境要求

- Go 1.21+
- PostgreSQL 15+

## 配置文件

\`\`\`yaml
database:
  host: localhost
  port: 5432
\`\`\`
`,
  },
  p4: {
    title: '主题 API',
    content: `# 主题 API

## 获取主题列表

\`\`\`bash
GET /api/themes
Authorization: Bearer {token}
\`\`\`

## 创建主题

\`\`\`bash
POST /api/themes
Content-Type: application/json

{
  "name": "用户手册",
  "slug": "user-guide"
}
\`\`\`
`,
  },
}

export const VIEWER_INITIAL_PAGE_ID = 'p2'
export const EDITOR_INITIAL_PAGE_ID = 'p2'
