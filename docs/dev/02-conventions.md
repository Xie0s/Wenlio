# 02 - 开发约束与代码规范

> **文档说明**：定义后端开发的技术约束、编码规范、命名规则与统一响应结构，确保团队代码风格一致、可维护。

---

## 1. 技术栈约束

| 类别     | 约束                                         |
| -------- | -------------------------------------------- |
| 语言版本 | Go 1.22+，启用 Go Modules                    |
| Web 框架 | Gin v1.10+，禁止混用其他 HTTP 框架           |
| 数据库   | MongoDB 6.x+，使用官方 mongo-driver v1.x     |
| 认证     | JWT HS256（golang-jwt/jwt v5），禁止明文存储密钥 |
| 日志     | Zap，禁止使用 fmt.Println / log 标准库输出    |
| 配置     | Viper，配置文件使用 YAML 格式                 |
| 校验     | go-playground/validator v10                   |
| 前端内嵌 | Go embed.FS，Vue 3 构建产物打包进二进制       |

---

## 2. 文件注释规范

### 2.1 文件头部注释（强制）

**每个 `.go` 文件顶部必须添加功能说明注释**，格式如下：

```go
// Package handler 文档页管理 HTTP 处理器
//
// 职责：处理文档页相关的 HTTP 请求，包括创建、查询、发布/下线、导入等操作。
// 对外接口：RegisterPageRoutes() 注册文档页路由组
package handler
```

**规则**：
- 第一行：`// Package <包名> <一句话功能概述>`
- 第二行空行后：`// 职责：<该文件的职责与功能边界>`
- 第三行：`// 对外接口：<主要导出的函数/类型/接口>`
- `package` 声明紧跟注释之后

### 2.2 函数/方法注释

导出函数必须添加 GoDoc 注释：

```go
// CloneVersion 克隆版本
// 深拷贝指定版本下的所有章节和文档页，生成新版本（draft 状态）
func (s *VersionService) CloneVersion(ctx context.Context, sourceVersionID primitive.ObjectID, req *dto.CloneVersionReq) (*entity.Version, error) {
```

### 2.3 结构体注释

导出结构体必须添加用途说明：

```go
// Page 文档页模型
// 对应 MongoDB pages 集合，存储 Markdown 原文内容
type Page struct {
```

---

## 3. 命名规范

### 3.1 包命名

| 规则       | 示例                   | 说明               |
| ---------- | ---------------------- | ------------------ |
| 全小写     | `handler`, `service`   | 禁止下划线和驼峰   |
| 简短有意义 | `dto`, `errcode`       | 避免过长包名       |
| 单数形式   | `model`, `entity`      | 不使用复数         |

### 3.2 文件命名

| 规则       | 示例                          | 说明                     |
| ---------- | ----------------------------- | ------------------------ |
| 蛇形命名   | `page_handler.go`             | 下划线分隔               |
| 模块前缀   | `version_service.go`          | 按业务模块前缀归类       |
| 测试文件   | `page_service_test.go`        | `_test.go` 后缀          |

### 3.3 变量与函数命名

| 类型           | 规则       | 正确示例               | 错误示例           |
| -------------- | ---------- | ---------------------- | ------------------ |
| 导出函数/类型  | 大驼峰     | `PublishPage`          | `publish_page`     |
| 私有函数/变量  | 小驼峰     | `buildQuery`           | `build_query`      |
| 常量           | 大驼峰     | `RoleTenantAdmin`      | `ROLE_TENANT_ADMIN`|
| 接口           | 大驼峰+er  | `PageRepository`       | `IPageRepo`        |
| MongoDB 字段   | 蛇形命名   | `tenant_id`            | `tenantId`         |
| JSON 字段      | 蛇形命名   | `theme_slug`           | `themeSlug`        |
| URL 路径       | 短横线     | `/api/v1/doc-pages`    | `/api/v1/docPages` |

### 3.4 常量命名

```go
// 角色常量
const (
    RoleSuperAdmin  = "super_admin"
    RoleTenantAdmin = "tenant_admin"
)

// 版本状态常量
const (
    VersionStatusDraft     = "draft"
    VersionStatusPublished = "published"
    VersionStatusArchived  = "archived"
)

// 文档页状态常量
const (
    PageStatusDraft     = "draft"
    PageStatusPublished = "published"
)

// 评论状态常量
const (
    CommentStatusPending  = "pending"
    CommentStatusApproved = "approved"
    CommentStatusRejected = "rejected"
)

// 租户状态常量
const (
    TenantStatusActive    = "active"
    TenantStatusSuspended = "suspended"
)
```

---

## 4. 统一响应结构

### 4.1 响应格式

所有 API 接口使用统一 JSON 响应结构：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "req_abc123"
}
```

**字段说明**：

| 字段         | 类型     | 必填 | 说明                              |
| ------------ | -------- | ---- | --------------------------------- |
| `code`       | `int`    | 是   | 业务状态码，0 为成功，非 0 为失败 |
| `message`    | `string` | 是   | 状态描述                          |
| `data`       | `any`    | 否   | 响应数据，失败时可为 null          |
| `request_id` | `string` | 是   | 请求唯一标识，用于日志追踪        |

### 4.2 分页响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  },
  "request_id": "req_abc123"
}
```

### 4.3 Go 实现规范

```go
// Response 统一响应结构体
type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    RequestID string      `json:"request_id"`
}

// Pagination 分页信息
type Pagination struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
}

// PageData 分页数据容器
type PageData struct {
    List       interface{} `json:"list"`
    Pagination Pagination  `json:"pagination"`
}

// 响应快捷方法
func Success(c *gin.Context, data interface{})
func SuccessWithPage(c *gin.Context, list interface{}, pagination Pagination)
func Fail(c *gin.Context, err *AppError)
func FailWithCode(c *gin.Context, code int, msg string)
```

### 4.4 HTTP 状态码使用规则

| HTTP 状态码 | 适用场景                         |
| ----------- | -------------------------------- |
| `200`       | 请求成功（查询、更新、删除）     |
| `201`       | 资源创建成功                     |
| `400`       | 参数校验失败                     |
| `401`       | 未认证 / Token 过期              |
| `403`       | 权限不足                         |
| `404`       | 资源不存在                       |
| `409`       | 资源冲突（如 slug 重复）         |
| `422`       | 业务规则校验失败                 |
| `500`       | 服务器内部错误                   |

---

## 5. 分页与排序规范

### 5.1 分页请求参数

| 参数          | 类型     | 默认值 | 范围           | 说明     |
| ------------- | -------- | ------ | -------------- | -------- |
| `page`        | `int`    | 1      | ≥ 1            | 当前页码 |
| `page_size`   | `int`    | 20     | 1 ~ 100        | 每页条数 |
| `sort_by`     | `string` | —      | 允许字段白名单 | 排序字段 |
| `sort_order`  | `string` | `desc` | `asc` / `desc` | 排序方向 |

### 5.2 关键字搜索

| 参数      | 类型     | 说明                   |
| --------- | -------- | ---------------------- |
| `keyword` | `string` | 模糊搜索关键字（可选） |

---

## 6. 错误处理规范

### 6.1 错误传递原则

- Service 层返回业务错误码（`*AppError`），不返回原生 `error`
- Repository 层返回包装后的错误，包含上下文信息
- Handler 层统一通过 `response.Fail()` 输出
- 禁止吞掉错误（忽略 error 返回值）

### 6.2 AppError 结构

```go
// AppError 业务错误
type AppError struct {
    HTTPCode int    // HTTP 状态码
    Code     int    // 业务错误码
    Message  string // 用户可见消息
    Internal error  // 内部错误（仅日志，不暴露）
}
```

### 6.3 错误日志记录

```go
// 所有 500 级别错误必须记录完整堆栈
// 所有 4xx 级别错误记录请求参数与错误原因
// 日志必须包含 request_id 用于链路追踪
```

---

## 7. MongoDB 操作规范

### 7.1 文档设计原则

- **内嵌优先**：一对少量关系使用内嵌文档（如评论的 `author` 内嵌昵称/邮箱）
- **引用分离**：一对多量关系使用 ObjectID 引用（如版本 → 章节 → 文档页）
- **禁止深层嵌套**：内嵌层级不超过 3 层
- **统一时间格式**：使用 `time.Time`，MongoDB 存储为 UTC

### 7.2 通用字段约定

每个集合的文档必须包含以下通用字段：

```go
type BaseModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"    json:"updated_at"`
}

// 租户级数据追加
type TenantModel struct {
    BaseModel `bson:",inline"`
    TenantID  primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
}
```

### 7.3 软删除约定

- 主数据使用 `status` 字段实现逻辑删除（`active` / `suspended`）
- 文档页使用 `status` 字段控制可见性（`draft` / `published`）
- 查询默认过滤活跃/已发布状态
- 评论使用 `status` 字段管理审核状态（`pending` / `approved` / `rejected`）

### 7.4 索引规范

- 所有 `tenant_id` 字段必须建立索引
- 查询频繁的字段（如 `slug`、`version_id`）建立索引
- 复合索引优先：`{ tenant_id: 1, status: 1 }`、`{ tenant_id: 1, version_id: 1 }`
- 全文索引：`{ title: "text", content: "text" }`（pages 集合）
- 禁止在单个集合上创建超过 10 个索引

---

## 8. 日志规范

### 8.1 日志级别

| 级别    | 使用场景                               |
| ------- | -------------------------------------- |
| `DEBUG` | 开发阶段调试信息（生产环境关闭）       |
| `INFO`  | 关键业务操作（文档发布、版本克隆等）   |
| `WARN`  | 非致命异常（超时重试、降级）           |
| `ERROR` | 影响业务的错误（数据库连接失败等）     |
| `FATAL` | 系统无法启动的致命错误                 |

### 8.2 日志格式

```json
{
  "level": "info",
  "ts": "2026-02-27T10:00:00+08:00",
  "caller": "service/page_service.go:42",
  "msg": "page published",
  "request_id": "req_abc123",
  "tenant_id": "65a1b2c3d4e5f6",
  "user_id": "65a1b2c3d4e5f8",
  "page_id": "65a1b2c3d4e5fa",
  "page_title": "Quick Start"
}
```

### 8.3 必须记录日志的场景

- 认证事件：登录成功/失败、Token 刷新
- 文档发布/下线：每次状态变更
- 版本操作：创建、克隆、发布、归档
- 评论审核：批准/拒绝
- 租户管理：创建、停用、解封
- 异常/错误：所有 error 级别事件

---

## 9. 安全规范

| 规则                 | 说明                                       |
| -------------------- | ------------------------------------------ |
| 密码存储             | bcrypt 哈希，cost ≥ 12                     |
| JWT 签名             | HS256，密钥存储于配置文件/环境变量         |
| Token 有效期         | Access Token: 24h                          |
| 输入校验             | 所有用户输入必须经过 validator 校验         |
| NoSQL 注入防护       | 禁止拼接查询语句，统一使用驱动参数化       |
| 敏感字段             | 密码、Token 禁止出现在日志和响应中         |
| 租户 ID 保留词       | 创建租户时校验 `admin/api/assets/static` 等保留词 |
| CORS                 | 仅允许指定域名跨域                         |

---

## 10. Git 规范

### 10.1 分支策略

| 分支        | 用途           | 说明                  |
| ----------- | -------------- | --------------------- |
| `main`      | 生产代码       | 仅通过 PR 合入        |
| `dev`       | 开发集成       | 日常开发合入          |
| `feature/*` | 功能开发       | 从 dev 切出           |
| `fix/*`     | Bug 修复       | 从 dev 切出           |

### 10.2 Commit Message 格式

```
<type>(<scope>): <subject>

<body>
```

**type 类型**：

| 类型       | 说明         |
| ---------- | ------------ |
| `feat`     | 新功能       |
| `fix`      | Bug 修复     |
| `refactor` | 重构         |
| `docs`     | 文档变更     |
| `style`    | 代码格式     |
| `chore`    | 构建/工具    |

**示例**：

```
feat(reader): 新增文档页评论区组件
fix(auth): 修复 Token 过期后未跳转登录页的问题
feat(editor): 支持导入 .md 文件并解析 frontmatter
```

---

## 11. 依赖管理

- 所有第三方依赖通过 `go.mod` 管理
- 禁止使用 `replace` 指令（本地调试除外）
- 依赖更新需在 PR 中说明变更原因
- 核心依赖锁定主版本号，避免破坏性升级

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
