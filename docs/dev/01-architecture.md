# 01 - 系统架构设计

> **文档说明**：定义后端系统整体架构、分层职责、目录结构与部署方案，作为开发实施的技术蓝图。

---

## 1. 技术栈总览

| 层级     | 技术选型                    | 说明                              |
| -------- | --------------------------- | --------------------------------- |
| 语言     | Go 1.22+                   | 高性能、强类型、原生并发          |
| Web 框架 | Gin                         | 轻量高性能 HTTP 框架              |
| 数据库   | MongoDB 6.x+               | 文档型数据库，含内置文本索引      |
| 驱动     | mongo-driver v1.x (official)| MongoDB 官方 Go 驱动              |
| 认证     | JWT (HS256)                 | golang-jwt/jwt v5                 |
| 日志     | Zap                         | 高性能结构化日志                  |
| 配置     | Viper                       | 多格式配置读取                    |
| 校验     | go-playground/validator v10 | 结构体参数校验                    |
| 前端内嵌 | Go embed.FS                 | Vue 3 构建产物打包进二进制        |
| 文件存储 | 本地 / MinIO SDK            | 图片等媒体文件                    |

---

## 2. 系统架构图

```
┌─────────────────────────────────────────────────────────┐
│                      客户端层                            │
│  ┌──────────────────────────┐  ┌─────────────────────┐  │
│  │  管理后台（/admin/*）     │  │  文档阅读（/{tid}/*）│  │
│  │  Vue 3 SPA               │  │  Vue 3 SPA          │  │
│  └────────────┬─────────────┘  └──────────┬──────────┘  │
└───────────────┼───────────────────────────┼─────────────┘
                │ HTTPS                     │ HTTPS
┌───────────────▼───────────────────────────▼─────────────┐
│                      Nginx 反向代理                      │
│                (SSL 终结 / 所有请求转发)                  │
└───────────────────────────┬─────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────┐
│                  Go HTTP 服务 (Gin) :8080                │
│  ┌──────────────────────────────────────────────────┐   │
│  │                  中间件层                          │   │
│  │  CORS │ RequestID │ Logger │ Recovery              │   │
│  ├──────────────────────────────────────────────────┤   │
│  │                  路由层 (Router)                   │   │
│  │  /api/v1/auth │ /api/v1/admin │ /api/v1/tenant    │   │
│  │  /api/v1/public │ /* (SPA 回退)                   │   │
│  ├──────────────────────────────────────────────────┤   │
│  │                  认证与鉴权层                      │   │
│  │  JWT 解析 │ 角色校验 │ 租户隔离                    │   │
│  ├──────────────────────────────────────────────────┤   │
│  │                  Handler 层                        │   │
│  │  参数绑定 │ 参数校验 │ 调用 Service │ 响应封装     │   │
│  ├──────────────────────────────────────────────────┤   │
│  │                  Service 层                        │   │
│  │  业务逻辑 │ 权限校验 │ 发布流程 │ 版本克隆         │   │
│  ├──────────────────────────────────────────────────┤   │
│  │                  Repository 层                     │   │
│  │  数据访问 │ MongoDB CRUD │ 查询构建 │ 索引优化     │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  embed.FS：Vue 3 SPA 静态资源                     │   │
│  │  /assets/*   → 静态文件直返                       │   │
│  │  /* (非 API) → index.html（SPA 回退）             │   │
│  └──────────────────────────────────────────────────┘   │
└───────────────────────────┬─────────────────────────────┘
                            │
                    ┌───────▼───────┐
                    │   MongoDB     │
                    │  (数据持久化)  │
                    └───────────────┘
```

---

## 3. 分层架构详解

### 3.1 Router 层（路由注册）

- **职责**：注册路由、绑定中间件、分组 API 版本
- **规则**：不包含任何业务逻辑，仅做路由映射
- **分组**：
  - `/api/v1/auth/*` — 认证路由（公开 + 已认证混合）
  - `/api/v1/admin/*` — 超级管理员路由
  - `/api/v1/tenant/*` — 租户管理员路由
  - `/api/v1/public/*` — 公开路由（读者端，无需认证）
  - `/*` — SPA 回退（embed.FS）

### 3.2 Middleware 层（中间件）

- **职责**：请求预处理与后处理
- **内置中间件**：
  - `RequestID` — 每个请求分配唯一 ID，贯穿日志链路
  - `Logger` — 请求/响应结构化日志
  - `Recovery` — panic 恢复，返回 500
  - `CORS` — 跨域配置（开发环境 :5173 → :8080）
  - `JWTAuth` — Token 解析与验证
  - `RoleGuard` — 角色权限校验
  - `TenantIsolation` — 租户数据隔离注入

### 3.3 Handler 层（控制器）

- **职责**：HTTP 请求/响应处理
- **规则**：
  - 参数绑定与校验（使用 validator tag）
  - 调用 Service 层方法
  - 统一响应封装（不直接操作数据库）
  - 不包含业务逻辑判断

### 3.4 Service 层（业务逻辑）

- **职责**：核心业务逻辑实现
- **规则**：
  - 业务规则校验与状态管理
  - 版本克隆（深拷贝章节 + 文档页）
  - 发布/下线流程编排
  - 全文搜索查询构建
  - 不直接处理 HTTP 请求/响应

### 3.5 Repository 层（数据访问）

- **职责**：数据库操作封装
- **规则**：
  - 提供 CRUD 基础操作
  - 封装复杂查询（文本搜索、聚合管道）
  - 隔离数据库驱动细节
  - 不包含业务逻辑

### 3.6 Model 层（数据模型）

- **职责**：定义数据结构
- **分类**：
  - `entity` — MongoDB 文档映射结构体（含 bson tag）
  - `dto` — 请求/响应数据传输对象（含 json/binding tag）

---

## 4. Go 内嵌 Vue 3 SPA 机制

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

**关键设计**：

- `/api/*` 路由 → API Handler 处理
- `/assets/*` 路由 → embed.FS 直接返回静态资源
- 其余所有路由 → 返回 `index.html`，由 Vue Router 客户端路由接管
- `/admin/*` 和 `/{tenant_id}/*` 均由同一个 Vue 3 应用处理

---

## 5. 项目目录结构

```
/（项目根目录 - Monorepo）
│
├── frontend/                        # Vue 3 单应用（admin + 文档阅读）
│   ├── src/
│   │   ├── components/
│   │   │   ├── admin/               # 管理后台专用组件
│   │   │   ├── reader/              # 文档阅读专用组件
│   │   │   │   ├── MarkdownRenderer.vue
│   │   │   │   ├── DocSidebar.vue
│   │   │   │   └── CommentSection.vue
│   │   │   └── common/              # 通用组件
│   │   ├── views/
│   │   │   ├── admin/               # 管理后台页面
│   │   │   └── reader/              # 文档阅读页面
│   │   ├── lib/
│   │   │   ├── stores/
│   │   │   │   ├── auth.ts          # 登录状态
│   │   │   │   ├── theme.ts         # 主题/版本数据
│   │   │   │   └── page.ts          # 文档页数据（含缓存）
│   │   │   └── markdown.ts          # markdown-it 配置与初始化
│   │   ├── router/
│   │   │   └── index.ts             # 路由配置（含导航守卫）
│   │   ├── api/                     # Go API 请求封装
│   │   └── types/                   # TypeScript 类型定义
│   ├── vite.config.ts
│   └── package.json
│
└── backend/                         # Go 后端
    ├── cmd/
    │   └── main.go                  # 程序入口
    ├── internal/
    │   ├── config/
    │   │   └── config.go            # 配置结构体与加载
    │   ├── middleware/
    │   │   ├── cors.go              # 跨域中间件
    │   │   ├── jwt.go               # JWT 认证中间件
    │   │   ├── logger.go            # 请求日志中间件
    │   │   ├── recovery.go          # 异常恢复中间件
    │   │   ├── request_id.go        # 请求 ID 中间件
    │   │   ├── role_guard.go        # 角色鉴权中间件
    │   │   └── tenant.go            # 租户隔离中间件
    │   ├── router/
    │   │   ├── router.go            # 路由总入口
    │   │   ├── auth.go              # 认证路由
    │   │   ├── admin.go             # 超级管理员路由（租户管理）
    │   │   ├── tenant.go            # 租户管理员路由
    │   │   └── public.go            # 公开接口路由（读者端）
    │   ├── handler/
    │   │   ├── spa_handler.go       # SPA embed 处理
    │   │   ├── auth_handler.go
    │   │   ├── tenant_handler.go
    │   │   ├── user_handler.go
    │   │   ├── theme_handler.go
    │   │   ├── version_handler.go
    │   │   ├── section_handler.go
    │   │   ├── page_handler.go
    │   │   ├── comment_handler.go
    │   │   ├── search_handler.go
    │   │   ├── upload_handler.go
    │   │   └── public_handler.go
    │   ├── service/
    │   │   ├── auth_service.go
    │   │   ├── tenant_service.go
    │   │   ├── user_service.go
    │   │   ├── theme_service.go
    │   │   ├── version_service.go   # 版本克隆、发布、归档
    │   │   ├── section_service.go
    │   │   ├── page_service.go      # 文档发布/下线、导入 Markdown
    │   │   ├── comment_service.go   # 评论审核流程
    │   │   ├── search_service.go    # 全文搜索
    │   │   └── upload_service.go
    │   ├── repository/
    │   │   ├── base_repo.go         # 基础 CRUD 封装
    │   │   ├── tenant_repo.go
    │   │   ├── user_repo.go
    │   │   ├── theme_repo.go
    │   │   ├── version_repo.go
    │   │   ├── section_repo.go
    │   │   ├── page_repo.go
    │   │   ├── comment_repo.go
    │   │   └── media_repo.go
    │   └── model/
    │       ├── entity/              # MongoDB 文档映射
    │       │   ├── tenant.go
    │       │   ├── user.go
    │       │   ├── theme.go
    │       │   ├── version.go
    │       │   ├── section.go
    │       │   ├── page.go
    │       │   ├── comment.go
    │       │   └── media.go
    │       └── dto/                 # 请求/响应传输对象
    │           ├── auth_dto.go
    │           ├── tenant_dto.go
    │           ├── theme_dto.go
    │           ├── version_dto.go
    │           ├── page_dto.go
    │           ├── comment_dto.go
    │           └── common_dto.go    # 分页、排序等通用 DTO
    ├── pkg/
    │   ├── response/
    │   │   └── response.go          # 统一响应封装
    │   ├── errcode/
    │   │   └── errcode.go           # 错误码定义
    │   ├── jwt/
    │   │   └── jwt.go               # JWT 工具
    │   ├── mongo/
    │   │   └── client.go            # MongoDB 连接管理
    │   ├── logger/
    │   │   └── logger.go            # 日志工具
    │   ├── utils/
    │   │   ├── slug.go              # Slug 生成与校验
    │   │   ├── time.go              # 时间工具
    │   │   ├── crypto.go            # 加密工具（密码哈希）
    │   │   └── pagination.go        # 分页工具
    │   └── constants/
    │       ├── role.go              # 角色常量
    │       ├── status.go            # 各实体状态常量
    │       └── reserved.go          # 租户 ID 保留词
    ├── frontend/                    # ← Vue 3 build 产物复制至此（embed 目标）
    │   └── dist/                    # （gitignore，由构建脚本生成）
    ├── config/
    │   └── config.yaml
    ├── scripts/
    │   ├── migrate/
    │   │   └── init_indexes.go      # 索引初始化
    │   └── seed/
    │       └── super_admin.go       # 超级管理员初始化
    ├── go.mod
    ├── go.sum
    ├── Makefile
    └── README.md
```

---

## 6. 请求处理流程

```
HTTP Request
    │
    ▼
┌──────────────┐
│  Middleware   │  RequestID → Logger → CORS → Recovery
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Router     │  路由匹配 → 分组中间件（JWTAuth / RoleGuard / TenantIsolation）
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handler    │  Bind → Validate → Call Service → Response
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Service    │  业务校验 → 发布/克隆/搜索 → 操作执行
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │  Query Build → MongoDB Driver → Return Result
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   MongoDB    │  数据持久化
└──────────────┘
```

**公开接口特殊流程**（读者端，无需认证）：

```
HTTP Request (GET /api/v1/public/pages/:page_id)
    │
    ▼
┌──────────────┐
│  Middleware   │  RequestID → Logger → CORS → Recovery（无 JWT）
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Handler     │  Bind page_id → Call Service
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Service     │  校验 status = published → 返回 Markdown 原文
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │  FindByID + status=published 过滤
└──────────────┘
```

---

## 7. 多租户数据隔离架构

### 7.1 隔离策略实现

采用 **共享数据库 + 字段隔离** 方案：

- 所有租户数据存储在同一 MongoDB 数据库中
- 通过 `tenant_id` 字段实现租户级隔离
- 中间件自动注入隔离条件，业务层无需手动处理
- 公开接口通过 URL 路径中的 `tenant_id` 参数隔离数据

### 7.2 隔离级别矩阵

| 数据类型         | 隔离字段      | 隔离方式                    |
| ---------------- | ------------- | --------------------------- |
| 租户             | —（平台级）   | 超级管理员可见              |
| 用户             | `tenant_id`   | TenantIsolation 中间件      |
| 主题 / 版本      | `tenant_id`   | TenantIsolation 中间件      |
| 章节 / 文档页    | `tenant_id`   | TenantIsolation 中间件      |
| 评论             | `tenant_id`   | TenantIsolation 中间件      |
| 媒体文件         | `tenant_id`   | TenantIsolation 中间件      |
| 公开接口数据     | `tenant_id`   | URL 参数 + Service 层校验   |

### 7.3 上下文传递

```go
// AuthContext 认证上下文，由 JWTAuth 中间件注入 Gin Context
type AuthContext struct {
    UserID   primitive.ObjectID // 当前用户 ID
    TenantID primitive.ObjectID // 租户 ID（super_admin 为空）
    Role     string             // 角色：super_admin / tenant_admin
}
```

---

## 8. 部署架构

### 8.1 开发环境

```bash
# 终端 1：启动 Go 后端（API 服务）
cd backend && go run ./cmd/main.go      # :8080

# 终端 2：启动 Vue 3（Vite 热更新，代理 API 至 :8080）
cd frontend && npm run dev              # :5173

# 访问管理后台：http://localhost:5173/admin/
# 访问文档阅读：http://localhost:5173/{tenant_id}/
```

Vite 开发代理配置：

```ts
// frontend/vite.config.ts
export default defineConfig({
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
    outDir: "../backend/frontend/dist",
    emptyOutDir: true,
  },
});
```

### 8.2 生产环境

```
                    ┌──────────────┐
                    │   Nginx      │
                    │  (SSL 终止)  │
                    │  :443        │
                    └──────┬───────┘
                           │ 所有请求
                    ┌──────▼───────┐
                    │  Go 服务     │
                    │  :8080       │
                    │              │
                    │  /api/*      │ → API 处理
                    │  /assets/*   │ → embed.FS 静态资源
                    │  /*          │ → embed.FS index.html（SPA 回退）
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │   MongoDB    │
                    │  :27017      │
                    └──────────────┘
```

**生产构建流程**：

```bash
# 步骤 1：构建 Vue 3（产物输出到 Go 项目的 frontend/dist/）
cd frontend && npm run build

# 步骤 2：编译 Go 二进制（embed 自动打包 Vue 3 产物）
cd backend && go build -o ../dist/docplatform ./cmd/main.go

# 步骤 3：上传二进制到服务器并重启
scp docplatform user@server:/opt/docplatform/
ssh user@server "systemctl restart docplatform"
```

**服务器只需运行三个进程**：

```bash
mongod              # MongoDB
docplatform         # Go 服务（内含 Vue 3）
nginx               # SSL + 代理
```

### 8.3 Nginx 生产配置

```nginx
server {
    listen 443 ssl;
    server_name microswift.cn;

    ssl_certificate     /etc/ssl/microswift.cn.pem;
    ssl_certificate_key /etc/ssl/microswift.cn.key;

    location / {
        proxy_pass         http://127.0.0.1:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name microswift.cn;
    return 301 https://$host$request_uri;
}
```

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
