# 宝塔面板部署指南

## 部署架构

```
用户浏览器
    │
    ▼
  Nginx (宝塔)
    ├── /assets/*、/favicon.ico 等静态文件 → 直接返回（Nginx 托管）
    ├── /api/*                             → 反向代理 → Go 服务 (:8081)
    ├── /uploads/*                         → 反向代理 → Go 服务 (:8081)
    └── 其他路径                            → 返回 index.html（SPA 路由）
```

Nginx 负责静态资源，Go 只处理 API 和文件上传，性能最优。

## 前置准备

### 1. 服务器环境

- 宝塔面板已安装 **Nginx**
- 已安装 **MongoDB**（宝塔软件商店 → 数据库 → MongoDB）
- 已安装 **Go 运行环境**（或上传编译好的二进制）

### 2. 项目文件上传

将以下内容上传到服务器（例如 `/www/wwwroot/docplatform/`）：

```
/www/wwwroot/docplatform/
├── backend           # Go 编译后的二进制文件
├── config/
│   └── config.yaml   # 生产配置
├── frontend/         # 前端构建产物（npm run build 的输出）
│   ├── index.html
│   ├── assets/
│   ├── sw.js
│   └── ...
└── uploads/          # 上传文件目录（运行后自动创建）
```

**构建步骤：**

```bash
# 本地构建前端
cd vue
npm run build

# 本地编译 Go 二进制（交叉编译到 Linux）
cd backend
set GOOS=linux
set GOARCH=amd64
go build -o docplatform cmd/server/main.go
```

### 3. 生产配置文件

在服务器上修改 `config/config.yaml`：

```yaml
server:
  port: 8081
  env: production        # 切换为生产模式

mongodb:
  uri: mongodb://localhost:27017
  database: docplatform

jwt:
  secret: "请替换为一个长随机字符串"   # 生产环境必须修改！
  expire_hours: 24

storage:
  type: local
  local_path: ./uploads

seed:
  super_admin_username: admin
  super_admin_password: "请替换为强密码"   # 首次启动后修改
```

## 宝塔面板配置步骤

### 步骤一：创建站点

1. 宝塔面板 → **网站** → **添加站点**
2. 域名填写你的域名（如 `docs.example.com`）
3. 根目录设置为前端产物目录：`/www/wwwroot/docplatform/frontend`
4. PHP 版本选 **纯静态**
5. 点击创建

### 步骤二：配置 Nginx

进入站点 → **设置** → **配置文件**，替换 `server {}` 块内容为：

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name docs.example.com;    # 替换为你的域名

    # ─── 宝塔默认日志（保留） ───
    access_log /www/wwwlogs/docplatform.log;
    error_log  /www/wwwlogs/docplatform.error.log;

    # ─── 前端静态资源根目录 ───
    root /www/wwwroot/docplatform/frontend;
    index index.html;

    # ─── 静态资源：Nginx 直接处理 ───
    # 带 hash 的资源文件，强缓存 1 年
    location /assets/ {
        expires 365d;
        add_header Cache-Control "public, immutable";
        try_files $uri =404;
    }

    # 其他静态文件（favicon、logo、sw.js 等）
    location ~* \.(ico|gif|png|jpg|jpeg|svg|webp|js|css|woff2?|ttf|eot)$ {
        expires 7d;
        add_header Cache-Control "public";
        try_files $uri =404;
    }

    # ─── API 请求：反向代理到 Go ───
    location /api/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 30s;
        proxy_read_timeout 60s;
        proxy_send_timeout 60s;
    }

    # ─── 上传文件访问：反向代理到 Go ───
    location /uploads/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # ─── SPA 路由 fallback ───
    # 非静态文件、非 API 的所有路径，返回 index.html 交给 Vue Router
    location / {
        try_files $uri $uri/ /index.html;
    }

    # ─── 禁止访问隐藏文件 ───
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }
}
```

### 步骤三：配置 SSL（推荐）

1. 站点设置 → **SSL** → 选择 **Let's Encrypt** 免费证书
2. 勾选 **强制 HTTPS**
3. 宝塔会自动在配置中添加 443 端口和证书路径

### 步骤四：启动 Go 后端

使用宝塔的 **Supervisor 管理器**（推荐）来守护 Go 进程：

1. 宝塔软件商店 → 搜索安装 **Supervisor管理器**
2. 添加守护进程：
   - **名称**：`docplatform`
   - **启动命令**：`/www/wwwroot/docplatform/docplatform`
   - **运行目录**：`/www/wwwroot/docplatform`
   - **启动用户**：`www`
3. 保存并启动

或者用命令行手动测试：

```bash
cd /www/wwwroot/docplatform
chmod +x ./docplatform
./docplatform
```

## 验证

1. 访问 `https://docs.example.com` → 应看到前端登录页
2. 访问 `https://docs.example.com/api/v1/health` → 应返回 `{"status":"ok"}`
3. 检查 Go 日志确认前端静态资源挂载（在 Nginx 分流模式下，Go 的 serveFrontend 不会被命中，但不影响功能）

## 注意事项

- **CORS**：Nginx 代理后，前后端同域，不存在跨域问题。Go 中的 CORS 中间件在生产环境可以收紧配置
- **文件上传限制**：宝塔 Nginx 默认限制上传 50MB，如需调整，在 Nginx 配置中添加 `client_max_body_size 100m;`
- **防火墙**：8081 端口无需对外开放，只需 80/443 端口。Go 服务仅监听 127.0.0.1 即可
- **Go 服务端口安全**：建议在 `config.yaml` 中将 Go 绑定地址改为 `127.0.0.1:8081`，避免直接暴露
