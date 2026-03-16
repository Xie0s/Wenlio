# 04 - API 接口设计

> **文档说明**：定义所有后端 API 接口的路由、参数、响应结构与权限要求，作为前后端协作的唯一口径。

---

## 1. API 设计规范

### 1.1 URL 规则

| 规则              | 示例                                          |
| ----------------- | --------------------------------------------- |
| RESTful 风格      | `/api/v1/tenant/pages/{id}`                   |
| 使用短横线分隔    | `/api/v1/tenant/doc-pages`，不使用下划线或驼峰 |
| 资源名用复数      | `/themes`、`/pages`、`/versions`              |
| 嵌套资源限制两级  | `/themes/{id}/versions`，不超过两级嵌套       |
| 操作型接口用动词  | `POST /pages/{id}/publish`                    |

### 1.2 HTTP 方法语义

| 方法     | 语义              |
| -------- | ----------------- |
| `GET`    | 查询（列表/详情） |
| `POST`   | 创建              |
| `PUT`    | 全量更新          |
| `PATCH`  | 局部更新          |
| `DELETE` | 删除（软删除）    |

### 1.3 通用请求头

| 请求头          | 说明                                        |
| --------------- | ------------------------------------------- |
| `Authorization` | `Bearer <access_token>`（公开接口不需要）   |
| `Content-Type`  | `application/json`                          |
| `X-Request-ID`  | 客户端可选传入，服务端若未传则自动生成      |

### 1.4 通用查询参数（列表接口）

| 参数          | 说明                                   |
| ------------- | -------------------------------------- |
| `page`        | 页码，默认 1                           |
| `page_size`   | 每页条数，默认 20，最大 100            |
| `sort_by`     | 排序字段（白名单校验）                 |
| `sort_order`  | `asc` / `desc`，默认 `desc`            |
| `keyword`     | 关键字模糊搜索                         |

---

## 2. 接口清单

### 2.1 认证模块（/api/v1/auth）

| 方法   | 路径                   | 说明              | 权限     |
| ------ | ---------------------- | ----------------- | -------- |
| POST   | `/auth/login`          | 登录，获取 Token  | 公开     |
| POST   | `/auth/refresh`        | 刷新 Token        | 公开     |
| POST   | `/auth/logout`         | 登出              | 已认证   |
| GET    | `/auth/me`             | 获取当前用户信息  | 已认证   |
| PATCH  | `/auth/me/password`    | 修改当前用户密码  | 已认证   |

**登录请求**：
```json
{ "username": "string", "password": "string" }
```

**登录响应**：
```json
{
  "access_token": "string",
  "expires_in": 86400,
  "user": {
    "id": "...",
    "name": "...",
    "role": "tenant_admin",
    "tenant_id": "acme"
  }
}
```

---

### 2.2 租户管理（/api/v1/admin/tenants）— super_admin 专属

| 方法   | 路径                            | 说明                    |
| ------ | ------------------------------- | ----------------------- |
| GET    | `/admin/tenants`                | 租户列表                |
| POST   | `/admin/tenants`                | 创建租户（含初始管理员）|
| GET    | `/admin/tenants/{id}`           | 租户详情                |
| PATCH  | `/admin/tenants/{id}`           | 更新租户信息            |
| POST   | `/admin/tenants/{id}/suspend`   | 封禁租户                |
| POST   | `/admin/tenants/{id}/activate`  | 解封租户                |

**创建租户请求**：
```json
{
  "id": "acme",
  "name": "Acme Corp",
  "logo_url": "",
  "admin_username": "admin",
  "admin_password": "password123",
  "admin_name": "管理员"
}
```

> 创建租户时同步创建该租户的初始管理员账号（`tenant_admin` 角色）。

---

### 2.3 用户管理（/api/v1/admin/users）— super_admin 专属

| 方法   | 路径                              | 说明                    |
| ------ | --------------------------------- | ----------------------- |
| GET    | `/admin/users`                    | 全平台用户列表          |
| POST   | `/admin/users`                    | 创建超级管理员账号      |
| POST   | `/admin/users/{id}/reset-password`| 重置密码                |

---

### 2.4 租户用户管理（/api/v1/tenant/users）— tenant_admin

| 方法   | 路径                              | 说明                    | 权限         |
| ------ | --------------------------------- | ----------------------- | ------------ |
| GET    | `/tenant/users`                   | 本租户用户列表          | tenant_admin |
| POST   | `/tenant/users`                   | 邀请/创建管理员         | tenant_admin |
| PATCH  | `/tenant/users/{id}`              | 更新用户信息            | tenant_admin |
| POST   | `/tenant/users/{id}/deactivate`   | 禁用用户                | tenant_admin |
| POST   | `/tenant/users/{id}/activate`     | 启用用户                | tenant_admin |
| POST   | `/tenant/users/{id}/reset-password` | 重置密码              | tenant_admin |

---

### 2.5 主题管理（/api/v1/tenant/themes）

| 方法   | 路径                    | 说明              | 权限         |
| ------ | ----------------------- | ----------------- | ------------ |
| GET    | `/tenant/themes`        | 主题列表          | tenant_admin |
| POST   | `/tenant/themes`        | 创建主题          | tenant_admin |
| GET    | `/tenant/themes/{id}`   | 主题详情          | tenant_admin |
| PATCH  | `/tenant/themes/{id}`   | 更新主题信息      | tenant_admin |
| DELETE | `/tenant/themes/{id}`   | 删除主题          | tenant_admin |
| PUT    | `/tenant/themes/sort`   | 批量更新排序      | tenant_admin |

**创建主题请求**：
```json
{
  "name": "API Reference",
  "slug": "api-reference",
  "description": "API 接口文档",
  "icon": "book"
}
```

**批量排序请求**：
```json
{
  "items": [
    { "id": "...", "sort_order": 1 },
    { "id": "...", "sort_order": 2 }
  ]
}
```

---

### 2.6 版本管理（/api/v1/tenant/themes/{themeId}/versions）

| 方法   | 路径                                       | 说明                 | 权限         |
| ------ | ------------------------------------------ | -------------------- | ------------ |
| GET    | `/tenant/themes/{themeId}/versions`        | 版本列表             | tenant_admin |
| POST   | `/tenant/themes/{themeId}/versions`        | 创建版本             | tenant_admin |
| GET    | `/tenant/versions/{id}`                    | 版本详情             | tenant_admin |
| PATCH  | `/tenant/versions/{id}`                    | 更新版本信息         | tenant_admin |
| POST   | `/tenant/versions/{id}/publish`            | 发布版本             | tenant_admin |
| POST   | `/tenant/versions/{id}/archive`            | 归档版本             | tenant_admin |
| POST   | `/tenant/versions/{id}/set-default`        | 设为默认版本         | tenant_admin |
| POST   | `/tenant/versions/{id}/clone`              | 克隆版本（深拷贝）   | tenant_admin |

**创建版本请求**：
```json
{
  "name": "v2.0",
  "label": "2.0（测试版）"
}
```

**克隆版本请求**：
```json
{
  "name": "v2.1",
  "label": "2.1（基于 v2.0 克隆）"
}
```

> 克隆操作深拷贝源版本下的所有章节和文档页，新版本状态为 `draft`。

---

### 2.7 章节管理（/api/v1/tenant/versions/{versionId}/sections）

| 方法   | 路径                                            | 说明              | 权限         |
| ------ | ----------------------------------------------- | ----------------- | ------------ |
| GET    | `/tenant/versions/{versionId}/sections`         | 章节列表          | tenant_admin |
| POST   | `/tenant/versions/{versionId}/sections`         | 创建章节          | tenant_admin |
| PATCH  | `/tenant/sections/{id}`                         | 更新章节标题      | tenant_admin |
| DELETE | `/tenant/sections/{id}`                         | 删除章节          | tenant_admin |
| PUT    | `/tenant/versions/{versionId}/sections/sort`    | 批量更新排序      | tenant_admin |

**创建章节请求**：
```json
{ "title": "快速开始" }
```

> 删除章节时，该章节下所有文档页也同步删除。

---

### 2.8 文档页管理（/api/v1/tenant/pages）

| 方法   | 路径                          | 说明                     | 权限         |
| ------ | ----------------------------- | ------------------------ | ------------ |
| GET    | `/tenant/sections/{sectionId}/pages` | 章节下文档页列表  | tenant_admin |
| POST   | `/tenant/sections/{sectionId}/pages` | 创建文档页        | tenant_admin |
| GET    | `/tenant/pages/{id}`          | 文档页详情（含内容）     | tenant_admin |
| PUT    | `/tenant/pages/{id}`          | 更新文档页（全量）       | tenant_admin |
| PATCH  | `/tenant/pages/{id}`          | 局部更新（自动保存）     | tenant_admin |
| DELETE | `/tenant/pages/{id}`          | 删除文档页               | tenant_admin |
| POST   | `/tenant/pages/{id}/publish`  | 发布文档页               | tenant_admin |
| POST   | `/tenant/pages/{id}/unpublish`| 下线文档页               | tenant_admin |
| POST   | `/tenant/pages/import`        | 导入 Markdown 文件       | tenant_admin |
| PUT    | `/tenant/sections/{sectionId}/pages/sort` | 批量更新排序 | tenant_admin |

**创建文档页请求**：
```json
{
  "title": "Quick Start",
  "slug": "quick-start",
  "content": "# Quick Start\n\n..."
}
```

**更新文档页请求（PUT 全量更新）**：
```json
{
  "title": "Quick Start",
  "slug": "quick-start",
  "content": "# Quick Start\n\n...",
  "section_id": "..."
}
```

**局部更新请求（PATCH 自动保存）**：
```json
{
  "content": "# Quick Start\n\n..."
}
```

**导入 Markdown 请求**（multipart/form-data）：

| 字段         | 类型   | 说明                        |
| ------------ | ------ | --------------------------- |
| `file`       | file   | `.md` 文件                  |
| `section_id` | string | 目标章节 ID                 |
| `version_id` | string | 目标版本 ID                 |

> 导入时自动解析 frontmatter 中的 `title`、`description`，Slug 根据 title 自动生成。

**批量发布版本接口**（独立于单页发布）：

通过版本管理接口 `POST /tenant/versions/{id}/publish` 触发，将该版本下所有 `draft` 文档页批量更新为 `published`。

---

### 2.9 评论管理（/api/v1/tenant/comments）

| 方法   | 路径                                | 说明              | 权限         |
| ------ | ----------------------------------- | ----------------- | ------------ |
| GET    | `/tenant/comments`                  | 评论列表（含筛选）| tenant_admin |
| POST   | `/tenant/comments/{id}/approve`     | 批准评论          | tenant_admin |
| POST   | `/tenant/comments/{id}/reject`      | 拒绝评论          | tenant_admin |
| DELETE | `/tenant/comments/{id}`             | 删除评论          | tenant_admin |
| POST   | `/tenant/comments/{id}/reply`       | 管理员回复评论    | tenant_admin |

**评论列表筛选参数**：

| 参数      | 类型   | 说明                                    |
| --------- | ------ | --------------------------------------- |
| `status`  | string | `pending` / `approved` / `rejected`     |
| `page_id` | string | 按文档页筛选                            |

**管理员回复请求**：
```json
{ "content": "感谢反馈，已修复。" }
```

> 管理员回复直接以 `approved` 状态创建，不需审核。

---

### 2.10 媒体文件上传（/api/v1/tenant/media）

| 方法   | 路径                   | 说明                    | 权限         |
| ------ | ---------------------- | ----------------------- | ------------ |
| POST   | `/tenant/media/upload` | 上传文件（图片/附件）   | tenant_admin |

**请求**：`multipart/form-data`，字段 `file`

**响应**：
```json
{
  "url": "https://...",
  "file_name": "screenshot.png",
  "file_size": 102400
}
```

---

### 2.12 公开接口（/api/v1/public）— 无需认证

读者端接口，路由统一挂载在 `/api/v1/public/` 下：

| 方法 | 路径                                       | 说明                                                |
| ---- | ------------------------------------------ | --------------------------------------------------- |
| GET  | `/public/tenants/{tenant_id}`              | 租户基础信息（名称、Logo）                          |
| GET  | `/public/tenants/{tenant_id}/themes`       | 租户下所有已发布主题列表                            |
| GET  | `/public/themes/{theme_id}/versions`       | 主题下所有 `published`/`archived` 版本列表          |
| GET  | `/public/versions/{version_id}/tree`       | 版本的章节+文档树（侧边栏数据，一次性返回）         |
| GET  | `/public/pages/{page_id}`                  | 单页内容（Markdown 原文，仅 published 状态）        |
| GET  | `/public/pages/{page_id}/comments`         | 页面 approved 评论列表                              |
| POST | `/public/pages/{page_id}/comments`         | 提交评论（默认 pending 状态）                       |
| GET  | `/public/search`                           | 全文搜索（必须携带 `tenant_id` 参数）               |

**版本文档树响应**（侧边栏核心数据）：

```json
{
  "sections": [
    {
      "id": "...",
      "title": "快速开始",
      "sort_order": 1,
      "pages": [
        {
          "id": "...",
          "title": "安装指南",
          "slug": "installation",
          "sort_order": 1
        }
      ]
    }
  ]
}
```

**提交评论请求**：
```json
{
  "author": {
    "name": "张三",
    "email": "zhang@example.com"
  },
  "content": "文档很清晰，感谢！",
  "parent_id": ""
}
```

**搜索接口参数**：

| 参数         | 类型   | 必填 | 说明               |
| ------------ | ------ | ---- | ------------------ |
| `q`          | string | 是   | 搜索关键词         |
| `tenant_id`  | string | 是   | 租户 ID            |
| `theme_id`   | string | 否   | 按主题筛选         |
| `version_id` | string | 否   | 按版本筛选         |

**搜索响应**：
```json
{
  "list": [
    {
      "page_id": "...",
      "title": "Quick Start",
      "snippet": "...匹配的<mark>关键词</mark>内容片段...",
      "theme_name": "API Reference",
      "version_name": "v1.0",
      "path": "/acme/api-reference/v1.0/quick-start"
    }
  ],
  "pagination": { "page": 1, "page_size": 20, "total": 5, "total_pages": 1 }
}
```

> **搜索接口鉴权说明**：读者无需登录即可搜索，搜索接口归入 `/public/` 路由组。通过 `tenant_id` 查询参数隔离数据，服务端强制校验该参数存在。

---

## 3. 错误响应示例

```json
{
  "code": 422301,
  "message": "该版本已处于归档状态，不可修改",
  "data": null,
  "request_id": "req_abc123"
}
```

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
