# 03 - 数据库设计

> **文档说明**：定义所有 MongoDB 集合的文档结构、索引策略与设计约定，作为数据层开发的权威参考。

---

## 1. 设计原则

| 原则         | 说明                                                                 |
| ------------ | -------------------------------------------------------------------- |
| 租户隔离     | 所有业务数据包含 `tenant_id`，API 层强制过滤                          |
| 软删除       | 主数据禁止物理删除，使用 `status` 字段标记停用/归档                   |
| 时间统一     | 所有时间字段使用 UTC 存储，API 响应时在前端转换                       |
| 引用完整性   | 使用 ObjectID 引用，Service 层保障引用合法性，不依赖数据库外键        |
| Slug 唯一性  | 主题和文档页的 `slug` 在同一租户/版本内唯一，用于 URL 路径段          |

---

## 2. 通用基础结构

```go
// BaseModel 所有集合的基础字段
type BaseModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    CreatedAt time.Time          `bson:"created_at"    json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"    json:"updated_at"`
}

// TenantModel 租户级数据基础字段
type TenantModel struct {
    BaseModel `bson:",inline"`
    TenantID  primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
}
```

---

## 3. 集合设计

### 3.1 tenants（租户）

```go
type Tenant struct {
    BaseModel `bson:",inline"`
    ID        string    `bson:"_id"`           // tenant_id（URL 路径段）
    Name      string    `bson:"name"`          // 显示名称
    LogoURL   string    `bson:"logo_url"`      // Logo URL
    Status    string    `bson:"status"`         // active / suspended
    CreatedAt time.Time `bson:"created_at"`
    UpdatedAt time.Time `bson:"updated_at"`
}
```

> **特殊说明**：`_id` 使用字符串类型（即 URL 中的 `{tenant_id}`），不使用 ObjectID。创建时校验保留词和格式（小写字母、数字、连字符，3-32 位）。

**索引**：`{ status: 1 }`

---

### 3.2 users（用户）

```go
type User struct {
    BaseModel      `bson:",inline"`
    TenantID       string    `bson:"tenant_id"`       // 关联租户 ID（super_admin 为空字符串）
    Username       string    `bson:"username"`         // 登录名
    Password       string    `bson:"password"`         // bcrypt 哈希
    Name           string    `bson:"name"`             // 显示姓名
    Email          string    `bson:"email"`
    Role           string    `bson:"role"`             // super_admin / tenant_admin
    Status         string    `bson:"status"`           // active / inactive
    LastLoginAt    time.Time `bson:"last_login_at"`
    LoginFailCount int       `bson:"login_fail_count"`
    LockedUntil    time.Time `bson:"locked_until"`
}
```

**索引**：
- `{ username: 1 }` unique
- `{ tenant_id: 1, status: 1 }`
- `{ tenant_id: 1, role: 1 }`

---

### 3.3 themes（文档主题）

```go
type Theme struct {
    TenantModel `bson:",inline"`
    Name        string `bson:"name"`        // 主题名称，如 "API Reference"
    Slug        string `bson:"slug"`        // URL 段，如 "api-reference"
    Description string `bson:"description"` // 主题描述
    Icon        string `bson:"icon"`        // 图标标识
    SortOrder   int    `bson:"sort_order"`  // 排序权重
    CreatedBy   primitive.ObjectID `bson:"created_by"`
}
```

**索引**：
- `{ tenant_id: 1, slug: 1 }` unique
- `{ tenant_id: 1, sort_order: 1 }`

---

### 3.4 versions（主题版本）

```go
type Version struct {
    TenantModel  `bson:",inline"`
    ThemeID      primitive.ObjectID `bson:"theme_id"`      // 关联主题
    Name         string             `bson:"name"`          // 如 "v1.0"、"latest"
    Label        string             `bson:"label"`         // 展示标签，如 "1.0（稳定版）"
    Status       string             `bson:"status"`        // draft / published / archived
    IsDefault    bool               `bson:"is_default"`    // 访问主题时默认展示此版本
    PublishedAt  time.Time          `bson:"published_at"`  // 首次发布时间
    CreatedBy    primitive.ObjectID `bson:"created_by"`
}
```

**索引**：
- `{ tenant_id: 1, theme_id: 1, status: 1 }`
- `{ tenant_id: 1, theme_id: 1, is_default: 1 }`

---

### 3.5 sections（章节）

```go
type Section struct {
    TenantModel `bson:",inline"`
    VersionID   primitive.ObjectID `bson:"version_id"` // 关联版本
    Title       string             `bson:"title"`      // 章节标题
    SortOrder   int                `bson:"sort_order"` // 排序权重
}
```

**索引**：`{ tenant_id: 1, version_id: 1, sort_order: 1 }`

---

### 3.6 pages（文档页）— 核心集合

```go
type Page struct {
    TenantModel  `bson:",inline"`
    VersionID    primitive.ObjectID `bson:"version_id"`   // 关联版本
    SectionID    primitive.ObjectID `bson:"section_id"`   // 关联章节
    Title        string             `bson:"title"`        // 文档标题
    Slug         string             `bson:"slug"`         // URL 段，如 "quick-start"
    Content      string             `bson:"content"`      // Markdown 原文
    Status       string             `bson:"status"`       // draft / published
    SortOrder    int                `bson:"sort_order"`   // 排序权重
    PublishedAt  time.Time          `bson:"published_at"` // 发布时间
    CreatedBy    primitive.ObjectID `bson:"created_by"`
}
```

**索引**：
- `{ tenant_id: 1, version_id: 1, status: 1 }` — 按版本查询文档列表
- `{ tenant_id: 1, version_id: 1, slug: 1 }` unique — 按 slug 查找文档页（同版本内唯一）
- `{ tenant_id: 1, section_id: 1, sort_order: 1 }` — 章节内排序
- `{ title: "text", content: "text" }` — 全文搜索索引

**全文搜索索引权重配置**：

```go
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

---

### 3.7 comments（评论）

```go
type Comment struct {
    TenantModel `bson:",inline"`
    PageID      primitive.ObjectID `bson:"page_id"`     // 关联文档页
    ParentID    primitive.ObjectID `bson:"parent_id"`   // 回复时指向父评论，顶层为空
    Author      CommentAuthor      `bson:"author"`      // 内嵌作者信息
    Content     string             `bson:"content"`     // 纯文本，限 1000 字
    Status      string             `bson:"status"`      // pending / approved / rejected
}

type CommentAuthor struct {
    Name  string `bson:"name"`  // 昵称（必填）
    Email string `bson:"email"` // 邮箱（选填，不公开展示）
}
```

**索引**：
- `{ tenant_id: 1, page_id: 1, status: 1 }` — 按页面查询已审核评论
- `{ tenant_id: 1, status: 1, created_at: -1 }` — 管理后台评论列表

---

### 3.8 media（上传文件记录）

```go
type Media struct {
    TenantModel `bson:",inline"`
    FileName    string `bson:"file_name"`    // 原始文件名
    FileURL     string `bson:"file_url"`     // 存储路径/URL
    FileSize    int64  `bson:"file_size"`    // 文件大小（字节）
    MimeType    string `bson:"mime_type"`    // MIME 类型
    UploadedBy  primitive.ObjectID `bson:"uploaded_by"`
}
```

**索引**：`{ tenant_id: 1, created_at: -1 }`

---

## 4. 索引设计汇总

| 集合       | 关键索引                                                          |
| ---------- | ----------------------------------------------------------------- |
| tenants    | `status`                                                          |
| users      | `username` unique; `tenant_id+status`                             |
| themes     | `tenant_id+slug` unique; `tenant_id+sort_order`                   |
| versions   | `tenant_id+theme_id+status`; `tenant_id+theme_id+is_default`     |
| sections   | `tenant_id+version_id+sort_order`                                 |
| pages      | `tenant_id+version_id+slug` unique; `title+content` text          |
| comments   | `tenant_id+page_id+status`; `tenant_id+status+created_at desc`   |
| media      | `tenant_id+created_at desc`                                       |

---

## 5. 数据库初始化策略

### 5.1 索引初始化

应用启动时自动执行索引创建（幂等），脚本位于 `scripts/migrate/init_indexes.go`。

### 5.2 种子数据

首次部署执行 `scripts/seed/super_admin.go`，创建超级管理员账号（用户名/密码从配置文件读取）。

```go
// 默认超级管理员
// username: 从 config.yaml 读取
// password: 从 config.yaml 读取（bcrypt 哈希后存储）
// role: super_admin
// tenant_id: ""（空，平台级）
```

---

## 6. 租户 ID 设计说明

租户 `_id` 使用字符串类型而非 ObjectID，因为它直接出现在 URL 路径中（`/{tenant_id}/`）。

**格式要求**：
- 只允许小写字母、数字、连字符
- 长度 3-32 位
- 不能以连字符开头或结尾
- 正则：`^[a-z0-9][a-z0-9-]{1,30}[a-z0-9]$`

**其他集合的 `tenant_id` 字段**存储的是租户 `_id` 字符串值（非 ObjectID），类型为 `string`。

> **注意**：由于租户 ID 为字符串类型，其他集合中的 `tenant_id` 字段也使用 `string` 类型，而非 `primitive.ObjectID`。TenantModel 需相应调整：

```go
type TenantModel struct {
    BaseModel `bson:",inline"`
    TenantID  string `bson:"tenant_id" json:"tenant_id"`
}
```

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
