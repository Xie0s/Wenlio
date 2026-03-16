# DocPlatform — 多租户文档管理平台

## 项目简介

DocPlatform 是一个多租户文档管理平台，支持主题/版本/章节/页面的层级文档管理、Markdown 编辑与阅读、评论系统等功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端框架 | Vue 3 + TypeScript |
| 构建工具 | Vite |
| UI 组件 | shadcn-vue + TailwindCSS |
| 状态管理 | Pinia |
| 路由 | Vue Router |
| 后端框架 | Go (Gin) |
| 数据库 | MongoDB |
| 认证 | JWT |

## 项目结构

```
doc-mather-1003/
├── vue/                    # 前端源码
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── composables/    # 组合式函数
│   │   ├── views/          # 页面视图
│   │   ├── lib/            # 业务逻辑层
│   │   └── assets/         # 静态资源
│   ├── vite.config.ts      # Vite 配置
│   └── package.json
├── backend/                # 后端源码
│   ├── cmd/server/         # 程序入口
│   ├── internal/           # 内部模块（handler/router/middleware 等）
│   ├── config/             # 配置文件
│   ├── frontend/           # 前端构建产物（由 Vite 构建输出）
│   ├── go.mod
│   └── go.sum
└── docs/                   # 项目文档
```

## 环境要求

- **Node.js** >= 18
- **Go** >= 1.26
- **MongoDB** >= 6.0（需本地运行或可访问的实例）

## 快速开始

### 1. 安装前端依赖

```bash
cd vue
npm install
```

### 2. 开发模式（前后端分离）

前端开发服务器（端口 5173，API 请求自动代理到后端）：

```bash
cd vue
npm run dev
```

后端服务（端口 1501）：

```bash
cd backend
go run cmd/server/main.go
```

开发时访问 `http://localhost:5173`，前端热更新，API 自动代理到 `http://localhost:1501`。

### 3. 构建与部署（单体模式）

本项目采用 **Go 后端内嵌前端静态资源** 的单体部署方式，只需运行一个 Go 服务即可同时提供前后端。

#### 步骤一：构建前端

```bash
cd vue
npm run build
```

构建产物会直接输出到 `backend/frontend/` 目录。

#### 步骤二：生成发布包（可选）

可选：一键生成发布包

```powershell
pwsh -File .\scripts\build-release.ps1
```

按目标平台构建：

 ```powershell
.\scripts\build-release.ps1 linux
.\scripts\build-release.ps1 windows
.\scripts\build-release.ps1 all
 ```

 或使用命名参数：

 ```powershell
pwsh -File .\scripts\build-release.ps1 -Target linux
pwsh -File .\scripts\build-release.ps1 -Target windows
pwsh -File .\scripts\build-release.ps1 -Target all
 ```

脚本会自动完成以下工作：

- 构建前端到 `backend/frontend/`
- 根据 `-Target` 参数编译 Go 后端为 `linux/amd64`、`windows/amd64` 或两者同时编译
- 组装发布目录到 `release/linux-amd64/` 与 `release/windows-amd64/`
- 生成压缩包：`release/docplatform-linux-amd64.zip`、`release/docplatform-windows-amd64.zip`

其中发布目录内已包含二进制文件、`config/`、`frontend/` 和空的 `uploads/` 目录，可直接上传后修改 `config/config.yaml` 使用。

#### 步骤三：手动启动后端

```bash
cd backend
go run cmd/server/main.go
```

#### 步骤四：访问应用

启动后控制台会输出访问地址：

```
HTTP 服务启动  {"addr": ":8081", "url": "http://localhost:8081"}
前端静态资源已挂载  {"path": "frontend"}
```

浏览器打开 **http://localhost:8081** 即可访问完整应用。

## 默认账号

首次启动会自动创建超级管理员（可在 `backend/config/config.yaml` 中修改）：

| 字段 | 默认值 |
|------|--------|
| 用户名 | `admin` |
| 密码 | `admin123` |

## 配置说明

后端配置文件位于 `backend/config/config.yaml`：

```yaml
server:
  port: 8081              # 服务端口
  env: development        # development / production

mongodb:
  uri: mongodb://localhost:27017
  database: docplatform   # 数据库名

jwt:
  secret: "..."           # JWT 密钥（生产环境务必修改）
  expire_hours: 24        # Token 过期时间

storage:
  type: local             # 文件存储方式：local | minio
  local_path: ./uploads   # 本地存储路径
```

## 常用命令

| 命令 | 目录 | 说明 |
|------|------|------|
| `npm run dev` | `vue/` | 启动前端开发服务器 |
| `npm run build` | `vue/` | 构建前端到 `backend/frontend/` |
| `go run cmd/server/main.go` | `backend/` | 启动后端服务 |
| `GOOS=linux GOARCH=amd64 go build -o docplatform cmd/server/main.go` | `backend/` | 编译 Linux 生产二进制 |

## 生产部署

详见 [宝塔面板部署指南](docs/deploy-baota.md)，包含：

- Nginx 反向代理 + 静态资源分流配置
- Go 后端 Supervisor 守护进程
- SSL 证书配置
- 生产环境安全建议
