# CDN 缓存与版本号方案

> 本文档记录前端项目部署到腾讯 EdgeOne CDN 后的缓存优化方案，供后续项目参考。

---

## 一、版本号实现

### 1.1 版本号格式

采用可追溯版本号格式：`v{packageVersion}+{buildSerial}.{commitShortSha}`

- **packageVersion**：从 `package.json` 的 `version` 字段读取（如 `0.5.0`）
- **buildSerial**：优先读取 CI 构建号（如 `GITHUB_RUN_NUMBER` / `CI_PIPELINE_IID`）
- **commitShortSha**：优先读取 CI 提交 SHA，兜底读取本地 Git 提交 SHA
- **本地兜底**：无 CI 时使用 `scripts/.build-counter` 生成 `local.{n}`

示例输出：`v0.5.0+1024.a1b2c3d4e5f6`、`v0.5.0+local.12.a1b2c3d4e5f6`

### 1.2 实现方式

#### package.json 构建脚本

```json
{
  "scripts": {
    "build": "vue-tsc -b && vite build && node scripts/postbuild.js"
  }
}
```

#### 后处理脚本 `scripts/postbuild.js`

```javascript
import { readFileSync, writeFileSync, existsSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const distDir = join(__dirname, '..', 'dist')
const swPath = join(distDir, 'sw.js')
const pkgPath = join(__dirname, '..', 'package.json')
const buildCounterPath = join(__dirname, '.build-counter')

// 1) 优先使用 CI 构建号 + Commit SHA
// 2) 无 CI 时使用 scripts/.build-counter 本地递增兜底
// 3) 生成 buildVersion：v{packageVersion}+{buildSerial}.{commitShortSha}
// 4) 注入 dist/sw.js 中 __SW_BUILD_VERSION__ 占位符
// 5) 同时生成 dist/version.json 版本清单

// 替换 SW 版本号占位符
let swContent = readFileSync(swPath, 'utf-8')
swContent = swContent.replace(/__SW_BUILD_VERSION__/g, buildVersion)
writeFileSync(swPath, swContent, 'utf-8')

console.log(`✅ Service Worker 版本号已更新: ${buildVersion}`)
```

#### Service Worker 版本号占位符 `public/sw.js`

```javascript
// 构建时由 postbuild.js 替换为实际版本号
const CACHE_VERSION = '__SW_BUILD_VERSION__' !== '__SW_BUILD_' + 'VERSION__' 
  ? '__SW_BUILD_VERSION__' 
  : 'v' + Date.now()
const FONT_CACHE = `font-cache-${CACHE_VERSION}`
```

---

## 二、Vite 构建配置

### 2.1 文件名 Hash 策略

确保所有静态资源带 hash，CDN 缓存能正确失效：

```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        // 显式配置 hash，确保 CDN 缓存正确失效
        entryFileNames: 'assets/[name]-[hash].js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash][extname]',
        manualChunks: {
          'three': ['three']  // 大型库单独打包
        }
      }
    }
  }
})
```

---

## 三、Nginx 缓存配置

### 3.1 缓存策略分层

| 资源类型 | 缓存时间 | Cache-Control |
|---------|---------|---------------|
| `index.html` | 不缓存 | `no-cache, no-store, must-revalidate` |
| `/assets/*` (带 hash) | 365 天 | `public, immutable` |
| `/api/*` | 不缓存 | `no-cache, no-store, must-revalidate, private` |
| 图片/字体 | 365 天 | `public, immutable` |
| `favicon.ico` | 30 天 | 默认 |

### 3.2 关键 Location 配置

```nginx
# Vite 构建产物（带 hash）- 长期缓存
location ^~ /assets/ {
    expires 365d;
    access_log off;
    add_header Cache-Control "public, immutable";
}

# API 接口 - 禁止缓存
location ^~ /api/v1/ {
    proxy_pass http://127.0.0.1:1340/api/v1/;
    # ... proxy headers ...
    
    add_header Cache-Control "no-cache, no-store, must-revalidate, private" always;
    add_header Pragma "no-cache" always;
    add_header Expires "0" always;
}

# SPA 路由 - HTML 不缓存
location / {
    try_files $uri $uri/ /index.html;
    
    add_header Cache-Control "no-cache, no-store, must-revalidate";
    add_header Pragma "no-cache";
    add_header Expires "0";
}
```

---

## 四、腾讯 EdgeOne CDN 配置建议

### 4.1 缓存规则配置

在 EdgeOne 控制台 → 站点 → 缓存配置：

1. **遵循源站**：开启"遵循源站 Cache-Control 头"
2. **节点缓存 TTL**：
   - `/index.html` → 不缓存
   - `/assets/*` → 缓存 365 天
   - `/api/*` → 不缓存
   - `/sw.js` → 不缓存

### 4.2 部署后刷新策略

每次部署后需刷新的文件：
- `/index.html`
- `/sw.js`
- `/version.json`

带 hash 的资源（`/assets/*`）无需手动刷新。

---

## 五、常见问题排查

### 5.1 用户看到旧版本

1. 检查 `index.html` 是否被 CDN 缓存（响应头 `X-Cache: HIT`）
2. 刷新 CDN 缓存：`/index.html`、`/sw.js`、`/version.json`
3. 用户清除浏览器缓存或硬刷新（Ctrl+Shift+R）

### 5.2 Service Worker 不更新

1. 检查 `sw.js` 版本号是否正确替换
2. 检查 CDN 是否缓存了旧的 `sw.js`
3. 用户需要关闭所有标签页后重新打开

### 5.3 验证构建版本号

```powershell
# 查看 dist/sw.js 中的版本号
Get-Content dist/sw.js | Select-String "CACHE_VERSION"

# 预期输出类似：
# const CACHE_VERSION = 'v0.4.2'
```

---

## 六、文件清单

| 文件 | 用途 |
|------|------|
| `vite.config.ts` | Vite 构建配置，hash 策略 |
| `public/sw.js` | Service Worker，字体缓存 |
| `scripts/postbuild.js` | 构建后处理，注入版本号 |
| `scripts/.build-counter` | 本地构建兜底计数器（自动生成，可加入 .gitignore） |
| `dist/version.json` | 构建版本清单（构建后自动生成） |
| `nginx.survey.conf` | Nginx 配置示例 |

---

## 七、版本号更新规则

- **buildSerial**：优先使用 CI 构建号（`GITHUB_RUN_NUMBER` / `CI_PIPELINE_IID` 等）
- **local buildSerial**：仅在无 CI 时使用 `scripts/.build-counter` 自动 +1（`local.{n}`）
- **minor**：功能迭代时手动修改 `package.json` 的 `version`
- **major**：大版本/破坏性变更时手动修改

示例：
```json
// 当前 v0.4.x → 升级到 v0.5.x
{ "version": "0.5.0" }

// 当前 v0.x.x → 升级到 v1.0.x
{ "version": "1.0.0" }
```
