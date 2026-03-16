# 文档管理平台 - 技术栈与开发规范 v2.0

> 版本：v2.0  
> 日期：2026-02-27  
> 变更说明：单应用架构，Go 内嵌 Vue 3，去除 VitePress 静态构建流程

---

## 一、系统架构

### 1.1 架构总览

**一个后端，一个前端，一个数据库。**

```
                        互联网访问者
                             │
              ┌──────────────▼──────────────┐
              │            Nginx             │
              │       microswift.cn          │
              │  SSL 终止 + 反向代理          │
              └──────────────┬──────────────┘
                             │ 所有请求
              ┌──────────────▼──────────────┐
              │        Go HTTP 服务          │
              │          :8080              │
              │                             │
              │  ┌─────────────────────┐    │
              │  │  /api/v1/*          │    │
              │  │  API 路由处理        │    │
              │  └─────────────────────┘    │
              │                             │
              │  ┌─────────────────────┐    │
              │  │  /* （其余所有路径） │    │
              │  │  Vue 3 SPA          │    │
              │  │  embed.FS 内嵌      │    │
              │  └─────────────────────┘    │
              │                             │
              └──────────────┬──────────────┘
                             │
                          MongoDB
```

### 1.2 Vue 3 单应用的两个功能区

同一个 Vue 3 应用，通过路由区分两个完全不同的界面：

| 路由前缀         | 界面     | 用途                                    |
| ---------------- | -------- | --------------------------------------- |
| `/admin/*`       | 管理后台 | 租户/主题/文档/版本的创建与管理，需登录 |
| `/{tenant_id}/*` | 文档阅读 | 面向读者的文档浏览，无需登录            |

**路由设计：**

```
/admin/                          ← 管理后台首页（登录页或仪表盘）
/admin/tenants                   ← 超级管理员：租户列表
/admin/themes                    ← 租户管理员：主题列表
/admin/themes/:themeId/versions  ← 版本管理
/admin/editor/:pageId            ← 文档编辑器

/{tenant_id}/                    ← 租户文档首页（主题列表）
/{tenant_id}/:themeSlug/         ← 主题首页（重定向至默认版本）
/{tenant_id}/:themeSlug/:version/            ← 版本首页
/{tenant_id}/:themeSlug/:version/:pageSlug   ← 具体文档页
```

### 1.3 Go 内嵌 Vue 3

Vue 3 构建产物通过 Go `embed` 包打包进二进制，**部署只需一个文件**。

```go
// backend/internal/handler/spa.go

//go:embed frontend/dist
var frontendDist embed.FS

func RegisterSPAHandler(r *gin.Engine) {
    sub, _ := fs.Sub(frontendDist, "frontend/dist")
    staticHandler := http.FileServer(http.FS(sub))

    // 静态资源（JS/CSS/图片）直接返回
    r.GET("/assets/*filepath", gin.WrapH(staticHandler))

    // SPA 回退：所有非 API 路由返回 index.html，由 Vue Router 接管
    r.NoRoute(func(c *gin.Context) {
        if strings.HasPrefix(c.Request.URL.Path, "/api/") {
            c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "接口不存在"})
            return
        }
        data, _ := frontendDist.ReadFile("frontend/dist/index.html")
        c.Data(http.StatusOK, "text/html; charset=utf-8", data)
    })
}
```

### 1.4 发布极简流程

无任何构建步骤，发布 = 数据库字段更新：

```
管理员点击"发布文档"
        │
        ▼
Go API: UPDATE pages SET status='published', published_at=NOW()
        │
        ▼
读者刷新页面 → Vue 3 调用 /api/v1/public/pages/:id → 返回已发布内容
        │
        ▼
markdown-it 渲染 Markdown → 展示给读者
```

### 1.5 Nginx 配置

```nginx
server {
    listen 443 ssl;
    server_name microswift.cn;

    ssl_certificate     /etc/ssl/microswift.cn.pem;
    ssl_certificate_key /etc/ssl/microswift.cn.key;

    # 所有请求转发至 Go 服务
    location / {
        proxy_pass         http://127.0.0.1:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}

# HTTP 强制跳转 HTTPS
server {
    listen 80;
    server_name microswift.cn;
    return 301 https://$host$request_uri;
}
```

### 1.6 部署步骤

```bash
# 开发完成后，执行一次构建流程：

# 步骤 1：构建 Vue 3（产物输出到 Go 项目的 frontend/dist/）
cd frontend && npm run build

# 步骤 2：编译 Go 二进制（embed 自动打包 Vue 3 产物）
cd backend && go build -o docplatform ./cmd/main.go

# 步骤 3：上传二进制到服务器并重启
scp docplatform user@server:/opt/docplatform/
ssh user@server "systemctl restart docplatform"
```

**服务器只需运行：**

```bash
mongod              # MongoDB
docplatform         # Go 服务（内含 Vue 3）
nginx               # SSL + 代理
```

---

## 二、技术栈

### 2.1 前端（Vue 3）

| 分类              | 选型                          | 说明                                                      |
| ----------------- | ----------------------------- | --------------------------------------------------------- |
| 框架              | Vue 3（Composition API）      | 核心框架                                                  |
| 构建              | Vite                          | 开发/构建工具                                             |
| UI 组件           | shadcn-vue                    | 基础组件库                                                |
| 样式              | TailwindCSS v3                | 原子化 CSS                                                |
| Markdown 排版样式 | @tailwindcss/typography       | 恢复 TailwindCSS Reset 后的 Markdown 渲染样式（prose 类） |
| 状态管理          | Pinia                         | 全局状态                                                  |
| 路由              | vue-router v4                 | 客户端路由，history 模式                                  |
| HTTP              | Axios                         | API 请求，统一拦截器                                      |
| 编辑器            | Tiptap v2                     | 富文本 / Markdown 双模式                                  |
| Markdown 渲染     | markdown-it                   | Markdown → HTML 解析（html: true，配合 DOMPurify 防 XSS） |
| 代码高亮          | shiki（@shikijs/markdown-it） | 200+ 语言，双主题（light/dark）                           |
| 自定义容器        | markdown-it-container         | Tip/Warning/Danger/Info                                   |
| 标题锚点          | markdown-it-anchor            | 生成标题 id，配合 TOC                                     |
| 目录生成          | markdown-it-toc-done-right    | 提取标题生成右侧浮动目录                                  |
| HTML 净化         | DOMPurify                     | 对 markdown-it 输出的 HTML 做 XSS 净化                    |
| 数学公式          | markdown-it-mathjax3          | 数学公式渲染（选配，按需引入）                            |
| 图标              | Lucide Vue                    | 图标库                                                    |

### 2.2 后端（Go）

| 分类       | 选型                        | 说明             |
| ---------- | --------------------------- | ---------------- |
| 语言       | Go 1.22+                    | 核心语言         |
| Web 框架   | Gin                         | HTTP 路由/中间件 |
| 认证       | golang-jwt/jwt v5           | JWT 签发与校验   |
| 数据库驱动 | mongo-driver v1.x           | MongoDB 官方驱动 |
| 配置       | Viper                       | 配置文件读取     |
| 日志       | Zap                         | 结构化日志       |
| 参数校验   | go-playground/validator v10 | 请求参数校验     |
| 文件存储   | 本地 / MinIO SDK            | 图片等媒体文件   |

### 2.3 基础设施

| 分类     | 选型                          |
| -------- | ----------------------------- |
| 数据库   | MongoDB 6.x（含内置文本索引） |
| 反向代理 | Nginx（SSL 终止）             |
| 版本控制 | Git                           |

---

## 三、项目目录结构

### 3.1 仓库结构（Monorepo）

```
/（项目根目录）
│
├── frontend/                    # Vue 3 单应用（admin + 文档阅读）
│   ├── src/
│   │   ├── components/
│   │   │   ├── admin/           # 管理后台专用组件
│   │   │   ├── reader/          # 文档阅读专用组件
│   │   │   │   ├── MarkdownRenderer.vue   # Markdown 渲染核心组件
│   │   │   │   ├── DocSidebar.vue         # 文档侧边栏
│   │   │   │   └── CommentSection.vue     # 评论区组件
│   │   │   └── common/          # 通用组件
│   │   ├── views/
│   │   │   ├── admin/           # 管理后台页面
│   │   │   └── reader/          # 文档阅读页面
│   │   ├── lib/                 # Pinia stores + 业务逻辑
│   │   │   ├── stores/
│   │   │   │   ├── auth.ts      # 登录状态
│   │   │   │   ├── theme.ts     # 主题/版本数据
│   │   │   │   └── page.ts      # 文档页数据（含缓存）
│   │   │   └── markdown.ts      # markdown-it 配置与初始化
│   │   ├── router/
│   │   │   └── index.ts         # 路由配置（含导航守卫）
│   │   ├── api/                 # Go API 请求封装
│   │   └── types/               # TypeScript 类型定义
│   ├── vite.config.ts           # proxy /api/* → :8080（开发阶段）
│   └── package.json
│
└── backend/                     # Go 后端
    ├── cmd/
    │   └── main.go              # 程序入口
    ├── internal/
    │   ├── handler/
    │   │   ├── spa.go           # Vue 3 SPA embed 处理
    │   │   ├── auth.go
    │   │   ├── tenant.go
    │   │   ├── theme.go
    │   │   ├── version.go
    │   │   ├── page.go
    │   │   ├── section.go
    │   │   ├── comment.go
    │   │   ├── search.go
    │   │   └── upload.go
    │   ├── service/             # 业务逻辑层
    │   ├── repository/          # MongoDB 数据访问层
    │   ├── model/               # 数据模型（struct + bson tag）
    │   ├── middleware/          # JWT、CORS、TenantGuard
    │   └── pkg/                 # 工具函数（response、validator 等）
    ├── frontend/                # ← Vue 3 build 产物复制至此（embed 目标）
    │   └── dist/                # （gitignore，由构建脚本生成）
    ├── config/
    │   └── config.yaml
    └── go.mod
```

### 3.2 构建脚本（Makefile）

```makefile
.PHONY: build dev

# 生产构建（前端 → 后端 embed → 编译二进制）
build:
	cd frontend && npm run build
	cp -r frontend/dist backend/frontend/dist
	cd backend && go build -o ../dist/docplatform ./cmd/main.go

# 开发模式（并行启动 Go + Vite dev server）
dev:
	@echo "启动 Go 后端 :8080"
	cd backend && go run ./cmd/main.go &
	@echo "启动 Vue 3 开发服务器 :5173"
	cd frontend && npm run dev
```

---

## 四、前端核心设计

### 4.1 路由守卫

```ts
// frontend/src/router/index.ts

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  // /admin/* 路由需要登录
  if (to.path.startsWith("/admin") && !authStore.isLoggedIn) {
    next("/admin/login");
    return;
  }
  next();
});
```

### 4.2 MarkdownRenderer 组件设计

#### markdown.ts — 异步单例初始化

> **关键设计**：Shiki 初始化是异步的，使用 Promise 单例确保只初始化一次，避免组件挂载时的竞态问题。  
> `html: true` 与 VitePress 行为保持一致（允许 Markdown 中嵌入 HTML），输出时用 DOMPurify 净化防 XSS。

```ts
// frontend/src/lib/markdown.ts
// 职责：markdown-it 单例异步初始化，全局复用
// 对外接口：renderMarkdown(content: string) → Promise<string>

import MarkdownIt from "markdown-it";
import shikiPlugin from "@shikijs/markdown-it";
import container from "markdown-it-container";
import anchor from "markdown-it-anchor";
import toc from "markdown-it-toc-done-right";
import DOMPurify from "dompurify";

const CONTAINERS = ["tip", "warning", "danger", "info"] as const;

// 单例：只初始化一次
let _mdInstance: MarkdownIt | null = null;
let _initPromise: Promise<MarkdownIt> | null = null;

async function createMarkdownIt(): Promise<MarkdownIt> {
  const md = new MarkdownIt({
    html: true, // 与 VitePress 一致，允许 Markdown 中嵌入 HTML
    linkify: true,
    typographer: true,
  });

  // Shiki 双主题：通过 CSS 变量切换，无需重渲染
  md.use(
    await shikiPlugin({
      themes: { light: "github-light", dark: "github-dark" },
      defaultColor: false,
    }),
  );

  // 标题锚点（需在 TOC 插件之前注册）
  md.use(anchor, {
    permalink: anchor.permalink.ariaHidden({ placement: "before" }),
    slugify: (s: string) =>
      encodeURIComponent(s.trim().toLowerCase().replace(/\s+/g, "-")),
  });

  // 右侧浮动目录（[[toc]] 占位符自动替换）
  md.use(toc, { containerClass: "doc-toc", listType: "ul" });

  // 自定义容器
  CONTAINERS.forEach((name) => md.use(container, name));

  return md;
}

function getMd(): Promise<MarkdownIt> {
  if (_mdInstance) return Promise.resolve(_mdInstance);
  if (!_initPromise) {
    _initPromise = createMarkdownIt().then((inst) => {
      _mdInstance = inst;
      return inst;
    });
  }
  return _initPromise;
}

/**
 * 渲染 Markdown 并净化输出 HTML
 * DOMPurify 过滤 <script>、事件属性等危险内容，保留代码高亮所需的 class/style
 */
export async function renderMarkdown(content: string): Promise<string> {
  const md = await getMd();
  const rawHtml = md.render(content);
  return DOMPurify.sanitize(rawHtml, {
    ADD_ATTR: ["class", "style", "tabindex"],
    USE_PROFILES: { html: true },
  });
}
```

#### MarkdownRenderer.vue — 渲染组件

> `prose` 类由 `@tailwindcss/typography` 提供，恢复 TailwindCSS Preflight Reset 清零的排版样式（标题层级、列表缩进、代码块、表格等），是与 VitePress 视觉一致的**前提条件**。

```vue
<!-- frontend/src/components/reader/MarkdownRenderer.vue -->
<template>
  <div
    class="prose prose-slate dark:prose-invert max-w-none"
    v-html="renderedHtml"
  />
</template>

<script setup lang="ts">
import { ref, watchEffect } from "vue";
import { renderMarkdown } from "@/lib/markdown";

const props = defineProps<{ content: string }>();
const renderedHtml = ref("");

// watchEffect 替代 computed：处理 renderMarkdown 的异步特性
watchEffect(async () => {
  renderedHtml.value = await renderMarkdown(props.content);
});
</script>
```

> **安全说明**：`html: true` 允许在 Markdown 中嵌入原始 HTML（与 VitePress 一致），`DOMPurify.sanitize()` 在输出阶段过滤所有危险内容，兼顾功能完整性与 XSS 安全。

### 4.3 文档页数据缓存

```ts
// frontend/src/lib/stores/page.ts

export const usePageStore = defineStore("page", () => {
  const cache = new Map<string, Page>(); // pageId → Page

  async function fetchPage(pageId: string): Promise<Page> {
    if (cache.has(pageId)) return cache.get(pageId)!;
    const page = await api.public.getPage(pageId);
    cache.set(pageId, page);
    return page;
  }

  return { fetchPage };
});
```

### 4.4 Markdown 样式规范

#### 为什么必须使用 @tailwindcss/typography

TailwindCSS 的 Preflight 会将所有 HTML 元素样式归零（`h1`-`h6` 字号相同、`ul/li` 无 bullet、`table` 无边框等）。`v-html` 注入的 Markdown HTML 如果不加 `prose` 类，渲染结果将**完全失去排版**。

**安装与配置：**

```bash
npm install -D @tailwindcss/typography
```

```ts
// tailwind.config.ts
export default {
  plugins: [require("@tailwindcss/typography")],
  darkMode: "class", // 暗色模式由 <html class="dark"> 控制
};
```

#### 暗色模式切换机制

Shiki 配置 `defaultColor: false` 后，代码块颜色通过 CSS 变量输出，无需重渲染即可响应主题切换。配合 TailwindCSS `dark:prose-invert` 实现全局反色。

```ts
// frontend/src/lib/stores/theme.ts（扩展现有 theme store）
// 职责：管理应用暗色/亮色模式，持久化到 localStorage

export const useThemeStore = defineStore("theme", () => {
  const isDark = ref(localStorage.getItem("theme") === "dark");

  function toggleDark() {
    isDark.value = !isDark.value;
    localStorage.setItem("theme", isDark.value ? "dark" : "light");
    document.documentElement.classList.toggle("dark", isDark.value);
  }

  // 初始化时同步 DOM
  document.documentElement.classList.toggle("dark", isDark.value);

  return { isDark, toggleDark };
});
```

Shiki 双主题 CSS 变量（自动生成，无需手写）：

```css
/* Shiki defaultColor: false 时生成的 CSS 变量结构 */
.shiki,
.shiki span {
  color: var(--shiki-light);
  background-color: var(--shiki-light-bg);
}

html.dark .shiki,
html.dark .shiki span {
  color: var(--shiki-dark);
  background-color: var(--shiki-dark-bg);
}
```

#### 自定义容器样式

`markdown-it-container` 生成的 HTML 结构如下，需在全局 CSS 中定义样式：

```html
<!-- 生成结构示例 -->
<div class="custom-block tip">
  <p class="custom-block-title">TIP</p>
  <p>内容</p>
</div>
```

```css
/* frontend/src/assets/markdown.css */
/* 职责：自定义容器（tip/warning/danger/info）的视觉样式，对齐 VitePress 风格 */

.custom-block {
  padding: 1rem 1.25rem;
  border-left: 4px solid;
  border-radius: 0 0.5rem 0.5rem 0;
  margin: 1.25rem 0;
}

.custom-block-title {
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
}

.custom-block.tip {
  border-color: #3b82f6;
  background: #eff6ff;
  color: #1e3a5f;
}
.custom-block.warning {
  border-color: #f59e0b;
  background: #fffbeb;
  color: #5c3d00;
}
.custom-block.danger {
  border-color: #ef4444;
  background: #fff1f2;
  color: #5c0011;
}
.custom-block.info {
  border-color: #8b5cf6;
  background: #f5f3ff;
  color: #2e1065;
}

/* 暗色模式 */
html.dark .custom-block.tip {
  background: #172554;
  color: #bfdbfe;
}
html.dark .custom-block.warning {
  background: #422006;
  color: #fde68a;
}
html.dark .custom-block.danger {
  background: #450a0a;
  color: #fecaca;
}
html.dark .custom-block.info {
  background: #2e1065;
  color: #ddd6fe;
}
```

在 `main.ts` 中全局引入：

```ts
import "@/assets/markdown.css";
```

---

## 五、后端核心设计

### 5.1 分层职责

```
Handler（路由层）
  → 绑定参数、校验、调用 Service、返回统一格式响应
  → 不含业务逻辑

Service（业务层）
  → 执行业务规则，组合 Repository 操作
  → 不直接操作数据库驱动

Repository（数据层）
  → 封装所有 MongoDB 操作，返回 Model struct
  → 所有查询必须携带 tenant_id 条件

Model（数据模型）
  → struct + bson tag，与 HTTP DTO 分离
```

### 5.2 统一响应格式

```go
// backend/internal/pkg/response.go

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
    c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus, code int, msg string) {
    c.JSON(httpStatus, Response{Code: code, Message: msg})
}
```

### 5.3 全文搜索（MongoDB Text Index）

```go
// 在 pages 集合上建立文本索引（启动时自动创建）
indexModel := mongo.IndexModel{
    Keys: bson.D{
        {Key: "title",   Value: "text"},
        {Key: "content", Value: "text"},
    },
    Options: options.Index().SetWeights(bson.D{
        {Key: "title",   Value: 10},  // 标题权重更高
        {Key: "content", Value: 1},
    }),
}
```

```go
// 搜索查询（强制 tenant_id 隔离）
filter := bson.D{
    {Key: "$text",      Value: bson.D{{Key: "$search", Value: keyword}}},
    {Key: "tenant_id",  Value: tenantID},
    {Key: "status",     Value: "published"},
}
```

### 5.4 公开 API 完整列表

读者端（无需鉴权）所需的完整公开接口，路由统一挂载在 `/api/v1/public/` 下：

| 方法 | 路径                                       | 说明                                                |
| ---- | ------------------------------------------ | --------------------------------------------------- |
| GET  | `/api/v1/public/tenants/:tenant_id`        | 租户基础信息（名称、Logo）                          |
| GET  | `/api/v1/public/tenants/:tenant_id/themes` | 租户下所有已发布主题列表                            |
| GET  | `/api/v1/public/themes/:theme_id/versions` | 主题下所有 `published`/`archived` 版本列表          |
| GET  | `/api/v1/public/versions/:version_id/tree` | 版本的章节+文档树（侧边栏数据，一次性返回）         |
| GET  | `/api/v1/public/pages/:page_id`            | 单页内容（Markdown 原文，仅 published 状态）        |
| GET  | `/api/v1/public/pages/:page_id/comments`   | 页面 approved 评论列表                              |
| POST | `/api/v1/public/pages/:page_id/comments`   | 提交评论（默认 pending 状态）                       |
| GET  | `/api/v1/public/search`                    | 全文搜索（**无需鉴权**，必须携带 `tenant_id` 参数） |

> **搜索接口鉴权说明**：读者无需登录即可搜索，因此搜索接口归入 `/public/` 路由组，无需 JWT Token。  
> 搜索接口通过 `tenant_id` 查询参数隔离数据，而非从 Token 中读取，服务端强制校验该参数存在。

**版本文档树接口设计**（侧边栏核心数据）：

```go
// GET /api/v1/public/versions/:version_id/tree
// 返回该版本下所有章节及其文档页（按 sort_order 排序），用于渲染侧边栏

type SectionTree struct {
    ID        string     `json:"id"`
    Title     string     `json:"title"`
    SortOrder int        `json:"sort_order"`
    Pages     []PageMeta `json:"pages"`  // 只含元数据，不含 content
}

type PageMeta struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Slug      string `json:"slug"`
    SortOrder int    `json:"sort_order"`
}
```

### 5.5 租户 ID 保留词校验

路由设计 `/{tenant_id}/*` 中，租户 ID 与系统路由共用同一路径层级，**必须**在创建租户时拒绝保留词：

```go
// backend/internal/service/tenant_service.go

// 系统保留的路径段，禁止用作 tenant_id
var reservedTenantIDs = map[string]struct{}{
    "admin":  {},
    "api":    {},
    "assets": {},
    "static": {},
    "health": {},
    "favicon.ico": {},
}

func validateTenantID(id string) error {
    if _, reserved := reservedTenantIDs[id]; reserved {
        return errors.New("该 ID 为系统保留词，不可使用")
    }
    // 只允许小写字母、数字、连字符，长度 3-32
    matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$`, id)
    if !matched {
        return errors.New("租户 ID 只允许小写字母、数字和连字符，长度 3-32 位")
    }
    return nil
}
```

---

## 六、编码规范

### 6.1 通用约定

- **语言**：注释、文档、错误信息使用中文；变量/函数/文件名使用英文
- **命名语义化**：禁止 `temp`、`data2`、`obj` 等无意义命名
- **单一职责**：函数/组件超过 80 行考虑拆分
- **禁止魔法值**：常量必须有命名定义

### 6.2 Go 规范

```
命名规则
├── 文件名：snake_case（page_service.go）
├── 包名：全小写单词（handler、service、repository）
├── 结构体/接口：PascalCase（PageService）
└── 变量：camelCase（tenantID、pageSlug）

错误处理
├── 所有 error 必须显式处理，禁止 _ 忽略
└── 向上传递时：fmt.Errorf("context: %w", err)
```

### 6.3 Vue 3 规范

```
命名规则
├── 组件文件：PascalCase（MarkdownRenderer.vue）
├── Composable：useXxx（usePageStore）
├── Pinia Store：defineStore('page', ...)
└── CSS 类：TailwindCSS 原子类，禁止自定义 class

组件职责
├── 组件层：只负责 UI 渲染和用户交互
├── lib/stores：状态管理和业务逻辑
└── api/：只做请求封装

样式规范
├── Modal/Dialog/Card：根组件统一 rounded-3xl
├── 按钮：基础 Button 仅图标 + Tooltip，不显示文字
└── 禁止使用 !important
```

### 6.4 API 规范

```
URL 规则
├── 全小写 kebab-case：/api/v1/tenant/doc-pages
├── 资源用名词复数：/themes、/pages、/versions
├── 层级不超过 3 层
└── 操作用 HTTP 方法区分：GET/POST/PUT/PATCH/DELETE

统一分页参数
GET /api/v1/tenant/pages?page=1&page_size=20&sort=created_at&order=desc
```

---

## 七、数据库规范

### 7.1 集合命名

| 集合       | 说明                       |
| ---------- | -------------------------- |
| `tenants`  | 租户                       |
| `users`    | 用户（含所有角色）         |
| `themes`   | 文档主题                   |
| `versions` | 主题版本                   |
| `sections` | 章节                       |
| `pages`    | 文档页（含 Markdown 内容） |
| `comments` | 评论                       |
| `media`    | 上传文件记录               |

### 7.2 字段约定

- 主键：`_id`（ObjectID），JSON 序列化输出为字符串 `id`
- 所有集合含 `tenant_id`（`super_admin` 用户除外）
- 时间：UTC，字段名 `created_at` / `updated_at` / `deleted_at`
- 软删除：`deleted_at != null` 视为已删除，查询时默认过滤

### 7.3 核心索引

```js
// pages 集合
{ tenant_id: 1, version_id: 1, status: 1 }   // 按版本查询文档列表
{ tenant_id: 1, slug: 1 }                     // 按 slug 查找文档页
{ title: "text", content: "text" }            // 全文搜索索引

// versions 集合
{ tenant_id: 1, theme_id: 1, status: 1 }

// comments 集合
{ tenant_id: 1, page_id: 1, status: 1 }
```

---

## 八、环境配置

### 8.1 后端配置（config.yaml）

```yaml
server:
  port: 8080
  env: development # development / production

mongodb:
  uri: mongodb://localhost:27017
  database: docplatform

jwt:
  secret: "your-secret-key-change-in-production"
  expire_hours: 24

storage:
  type: local # local | minio
  local_path: ./uploads
```

### 8.2 前端开发配置（vite.config.ts）

```ts
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "../backend/frontend/dist", // 直接输出到 Go embed 目录
    emptyOutDir: true,
  },
});
```

### 8.3 开发启动

```bash
# 终端 1：启动 Go 后端
cd backend && go run ./cmd/main.go

# 终端 2：启动 Vue 3（热更新）
cd frontend && npm run dev

# 访问管理后台：http://localhost:5173/admin/
# 访问文档阅读：http://localhost:5173/{tenant_id}/
```

---

## 九、Git 工作流

### 9.1 分支策略

```
main       ← 生产环境，只接受 PR 合并
dev        ← 开发主线
feature/*  ← 功能分支，从 dev 切出
fix/*      ← 修复分支
```

### 9.2 Commit 规范

```
格式：<type>(<scope>): <subject>

type:
  feat     新功能
  fix      Bug 修复
  refactor 重构
  style    代码格式
  docs     文档
  chore    构建/依赖/配置

示例：
  feat(reader): 新增文档页评论区组件
  fix(auth): 修复 Token 过期后未跳转登录页的问题
  feat(editor): 支持导入 .md 文件并解析 frontmatter
```
