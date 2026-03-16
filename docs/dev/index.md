# 文档管理平台 - 后端开发文档

> **文档说明**：本目录为后端开发文档的总索引，涵盖系统架构、开发规范、数据库设计、API 接口、认证权限、核心业务逻辑与错误码定义。
>
> **技术栈**：Go (Gin) + MongoDB + JWT + Vue 3 (embed.FS 内嵌)
>
> **需求文档**：[需求文档 v2.0](../需求文档v2.md) | [技术栈与开发规范 v2.0](../index-v2.md)

---

## 文档索引

### 1. [系统架构设计](./01-architecture.md)

- 技术栈总览
- 系统架构图（客户端 → Nginx → Go HTTP 服务 → MongoDB）
- 分层架构详解（Router → Middleware → Handler → Service → Repository → Model）
- Go 内嵌 Vue 3 SPA 机制（embed.FS）
- 项目目录结构（Monorepo：frontend/ + backend/）
- 请求处理流程
- 多租户数据隔离架构（共享数据库 + `tenant_id` 字段隔离）
- 部署架构（开发环境 / 生产环境）

### 2. [开发约束与代码规范](./02-conventions.md)

- 技术栈约束（语言版本、框架版本、依赖限制）
- 文件注释规范（文件头部注释格式、函数注释、结构体注释）
- 命名规范（包命名、文件命名、变量与函数命名、常量命名）
- 统一响应结构（成功/失败/分页响应格式）
- HTTP 状态码使用规则
- 分页与排序规范
- 错误处理规范（AppError 结构、错误传递原则）
- MongoDB 操作规范（文档设计、通用字段、软删除、索引）
- 日志规范（级别、格式、必须记录场景）
- 安全规范（密码、JWT、注入防护）
- Git 规范（分支策略、Commit Message 格式）

### 3. [数据库设计](./03-database.md)

- 设计原则（租户隔离、软删除、时间统一、引用完整性）
- 通用基础结构（BaseModel / TenantModel）
- 全部集合设计（8 个集合）：
  - `tenants` — 租户
  - `users` — 用户（含角色）
  - `themes` — 文档主题
  - `versions` — 主题版本
  - `sections` — 章节
  - `pages` — 文档页（含 Markdown 内容）
  - `comments` — 评论
  - `media` — 上传文件记录
- 索引设计
- 数据库初始化策略

### 4. [API 接口设计](./04-api-design.md)

- API 设计规范（URL 规则、HTTP 方法语义、通用请求头/查询参数）
- 完整接口清单（10 个模块，50+ 接口）：
  - 认证模块（登录、刷新、登出、当前用户）
  - 租户管理（CRUD + 启用/停用）— super_admin 专属
  - 用户管理（CRUD + 重置密码）— super_admin / tenant_admin
  - 主题管理（CRUD + 排序）
  - 版本管理（CRUD + 发布/归档 + 设为默认 + 克隆）
  - 章节管理（CRUD + 排序）
  - 文档页管理（CRUD + 发布/下线 + 导入 Markdown + 排序）
  - 评论管理（列表 + 审核/拒绝/删除 + 回复）
  - 媒体文件上传
  - 公开接口（租户信息、主题列表、版本列表、文档树、页面内容、评论、搜索）
- 请求/响应示例

### 5. [认证与权限设计](./05-auth.md)

- JWT 认证流程（登录 → Token 签发 → 请求认证 → 刷新 → 登出）
- Token 策略（HS256、有效期、吊销机制）
- Token Payload 结构
- 角色模型（super_admin / tenant_admin / viewer）
- 完整权限矩阵（功能 × 角色）
- 数据隔离中间件实现（TenantIsolation）
- 登录安全策略（失败锁定、密码策略）
- 接口安全分级（公开 / 认证 / 角色 / 隔离）

### 6. [核心业务逻辑](./06-business-logic.md)

- 文档发布流程（即时发布，无构建步骤）
- 版本管理逻辑（创建、克隆、发布、归档、设为默认）
- 版本克隆深拷贝实现（章节 + 文档页完整复制）
- 文档页状态管理（draft / published）
- 批量发布版本逻辑
- 评论审核流程（pending → approved / rejected）
- 全文搜索实现（MongoDB Text Index）
- 租户 ID 保留词校验
- Markdown 文件导入逻辑（frontmatter 解析）
- 文档自动保存机制
- Slug 唯一性校验
- 文档排序规则

### 7. [错误码与异常处理](./07-error-codes.md)

- 错误码编码规则（`XXYYZZ` = HTTP状态码 + 模块编号 + 错误序号）
- 模块编号分配（10 个模块）
- **完整错误码表**（60+ 错误码）：
  - 通用错误（参数、认证、权限、系统）
  - 认证/租户/用户错误
  - 主题/版本/章节/文档页错误
  - 评论/媒体错误
- Go 实现（AppError 结构体、错误码常量、错误响应处理、参数校验转换）
- 异常处理策略（Recovery 中间件、错误日志分级）

---

## 快速参考

### 角色权限速查

| 角色       | 标识           | 数据范围 | 核心职责                                           |
| ---------- | -------------- | -------- | -------------------------------------------------- |
| 超级管理员 | `super_admin`  | 全平台   | 租户创建与平台治理，不参与业务操作                 |
| 租户管理员 | `tenant_admin` | 本租户   | 管理本租户的主题/版本/文档/评论，邀请其他管理员    |
| 读者       | `viewer`       | 已发布   | 无需登录，浏览已发布文档，可匿名评论              |

### 文档层级速查

```
租户（Tenant）
  └── 主题（Theme）         ← 一个产品或文档分类
        └── 版本（Version） ← 如 v1.0、v2.0、latest
              └── 章节（Section）  ← 侧边栏分组
                    └── 文档页（Page）
```

### 文档状态速查

| 实体    | 状态值                                    | 说明                                  |
| ------- | ----------------------------------------- | ------------------------------------- |
| 租户    | `active` / `suspended`                    | 停用后文档不可访问                    |
| 版本    | `draft` / `published` / `archived`        | draft 不可见，archived 只读但可访问   |
| 文档页  | `draft` / `published`                     | draft 仅管理后台可见                  |
| 评论    | `pending` / `approved` / `rejected`       | 仅 approved 对读者可见                |

### API 前缀

```
认证接口：/api/v1/auth/*
管理接口：/api/v1/admin/*       （super_admin 专属，需认证）
租户接口：/api/v1/tenant/*      （tenant_admin 专属，需认证 + 租户隔离）
公开接口：/api/v1/public/*      （无需认证，面向读者）
```

### 路由结构速查

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

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**对应需求**：需求文档 v2.0
**最后更新**：2026 年 2 月
