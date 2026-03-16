# 访问控制功能变更说明

> 本文档记录「站点级访问控制」和「主题级访问控制」功能的所有变更点，供后续调试和维护定位。

---

## 一、后端变更

### 1. 模型层（Entity）

| 文件 | 变更 |
|------|------|
| `backend/internal/model/entity/tenant_settings.go` | 新增 `Access *AccessSettings` 字段；新增 `AccessSettings` 结构体（`MaintenanceMode bool`, `GalleryLoginRequired bool`） |
| `backend/internal/model/entity/theme.go` | 新增 `AccessMode string`（空/public/login/code）和 `AccessCode string` 字段 |

### 2. DTO 层

| 文件 | 变更 |
|------|------|
| `backend/internal/model/dto/tenant_settings_dto.go` | 新增 `UpdateAccessReq`、`AccessSettingsResp` DTO；`TenantSettingsResp` 中增加 `Access *AccessSettingsResp` |
| `backend/internal/model/dto/theme_dto.go` | `ThemeListItem` 新增 `AccessMode string`；`UpdateThemeReq` 新增 `AccessMode *string`、`AccessCode *string` |

### 3. Service 层

| 文件 | 变更 |
|------|------|
| `backend/internal/service/tenant_settings_service.go` | 新增 `toAccessResp()` 辅助函数；新增 `UpdateAccess()` 方法；新增 `GetAccessSettings()` 方法；`GetSettings()` 中增加 `resp.Access = toAccessResp(settings.Access)` |
| `backend/internal/service/theme_service.go` | `enrichThemes()` 中填充 `AccessMode`；`enrichThemesPublic()` 新增 `isAuthenticated` 参数，按 access_mode 过滤 login 类型主题；`ListPublic()`/`ListByFilterPublic()` 增加 `isAuthenticated` 参数；`Update()` 支持 `access_mode`/`access_code` 字段更新；新增 `VerifyAccessCode()` 方法 |

### 4. Handler 层

| 文件 | 变更 |
|------|------|
| `backend/internal/handler/tenant_settings_handler.go` | 新增 `UpdateAccess()` handler |
| `backend/internal/handler/public_handler.go` | 新增 `settingsService` 字段；`ListThemes()`/`ListThemesByFilter()` 从 OptionalJWTAuth 获取认证状态并传递给 service；新增 `GetTenantAccess()` handler；新增 `VerifyThemeAccessCode()` handler |

### 5. 中间件

| 文件 | 变更 |
|------|------|
| `backend/internal/middleware/jwt.go` | 新增 `OptionalJWTAuth()` 中间件：尝试提取 JWT，成功注入 auth_context，失败静默放行 |

### 6. JWT 扩展

| 文件 | 变更 |
|------|------|
| `backend/pkg/jwt/jwt.go` | 新增 `ThemeAccessClaims` 结构体；新增 `GenerateThemeAccessToken()` 函数（24h 有效期）；新增 `ParseThemeAccessToken()` 函数 |

### 7. 错误码

| 文件 | 变更 |
|------|------|
| `backend/pkg/errcode/errcode.go` | 新增模块 12（访问控制）：`ErrAccessCodeInvalid`、`ErrThemeLoginRequired`、`ErrThemeCodeRequired` |

### 8. 路由

| 文件 | 变更 |
|------|------|
| `backend/internal/router/router.go` | 租户设置路由新增 `PATCH /settings/access`；公开主题列表路由加 `OptionalJWTAuth` 中间件；新增公开路由 `GET /public/tenants/:tenant_id/access`；新增公开路由 `POST /public/themes/:id/verify-code` |

---

## 二、前端变更

### 1. 类型定义

| 文件 | 变更 |
|------|------|
| `vue/src/utils/types.ts` | 新增 `AccessMode` 类型；`Theme` 接口新增 `access_mode`；`UpdateThemeReq` 新增 `access_mode`/`access_code`；新增 `AccessSettings` 接口；`TenantSettings` 新增 `access` 字段；新增 `ACCESS_MODE_LABEL`/`ACCESS_MODE_COLOR` 常量映射 |

### 2. 业务逻辑层（lib / composable / store）

| 文件 | 变更 |
|------|------|
| `vue/src/lib/tenant-settings.ts` | 新增 `accessForm` reactive 状态；`loadSettings()` 中同步 access 数据；新增 `saveAccess()` 方法；return 中暴露新字段 |
| `vue/src/composables/useThemeEditor.ts` | `updateTheme()` 参数类型扩展支持 `access_mode`/`access_code` |
| `vue/src/stores/reader.ts` | 新增 `accessSettings` state 和 `loadAccessSettings()` 方法 |

### 3. 组件

| 文件 | 变更 |
|------|------|
| `vue/src/components/editor/theme-management/ThemeAccessSettings.vue` | **新建**。主题访问权限设置基础组件，提供 access_mode 选择和 access_code 输入 |
| `vue/src/components/editor/theme-management/ThemeFormDialog.vue` | `ThemeFormData` 新增 `access_mode`/`access_code`；表单初始化和重置新增字段；模板中集成 ThemeAccessSettings 组件 |
| `vue/src/components/editor/theme-management/ThemeCard.vue` | 导入 ACCESS_MODE 常量；主题名称后显示访问模式 Badge（仅非 public 时） |
| `vue/src/components/editor/edit-page/tool/EditorSettingsPanel.vue` | 主题表单新增 access_mode/access_code；watch 填充新字段；模板中集成 ThemeAccessSettings |
| `vue/src/components/editor/edit-page/tool/EditorToolbar.vue` | 导入 ACCESS_MODE 常量；设置按钮后显示访问模式指示 Badge |

### 4. 页面

| 文件 | 变更 |
|------|------|
| `vue/src/views/admin/TenantSettingsPage.vue` | 新增「访问控制」Accordion 区块，包含维护模式和画廊登录可见两个 Switch |
| `vue/src/views/reader/MaintenancePage.vue` | **新建**。维护模式展示页，含登录引导按钮 |
| `vue/src/views/reader/ThemeVerifyCodePage.vue` | **新建**。主题验证码输入页，6位验证码表单，成功后存储 theme_access token 并跳转 |

### 5. 路由

| 文件 | 变更 |
|------|------|
| `vue/src/router/index.ts` | 新增 `Maintenance` 和 `ThemeVerifyCode` 路由；`beforeEach` 新增第②步读者端访问控制：维护模式重定向、画廊登录拦截、主题级 code/login 拦截 |

---

## 三、Bug 修复（第二轮）

### 3.1 后端：内容接口缺乏主题访问控制

**问题**：`ListVersions`、`GetVersionTree`、`GetPage`、`ListComments` 四个公开接口未检查主题的 `access_mode`，导致 code/login 类型主题的内容可被直接获取。

**修复**：
| 文件 | 变更 |
|------|------|
| `backend/internal/service/theme_service.go` | 新增 `CheckPublicAccess()` 方法（检查 access_mode + 验证 theme_access token / JWT）；新增 `GetByIDUnsafe()` 方法（公开接口回溯用） |
| `backend/internal/handler/public_handler.go` | 新增 `extractThemeAccessToken()` 辅助函数（从 `X-Theme-Access` 头或 `theme_token` 查询参数提取 token）；`ListVersions` 直接检查 theme → access_mode；`GetVersionTree` 回溯 version → theme 检查；`GetPage` 回溯 page → version → theme 检查；`ListComments` 回溯 page → version → theme 检查；`GetRawMarkdown` 非 public 主题直接返回 403 |
| `backend/internal/router/router.go` | 为 versions/tree/pages/comments 路由添加 `OptionalJWTAuth` 中间件 |

### 3.2 后端：分类/标签公开计数包含 login 类型主题

**问题**：未登录用户查看分类和标签列表时，计数中包含了 `access_mode="login"` 的主题。

**修复**：
| 文件 | 变更 |
|------|------|
| `backend/internal/repository/theme_repo.go` | `ListMinimalByTenant` 投影新增 `access_mode` 字段 |
| `backend/internal/service/category_service.go` | `ListPublic()` 签名新增 `isAuthenticated bool`；计数时排除未登录用户不可见的 login 类型主题 |
| `backend/internal/service/tag_service.go` | `ListPublic()` 签名新增 `isAuthenticated bool`；计数时排除未登录用户不可见的 login 类型主题 |
| `backend/internal/handler/public_handler.go` | `ListPublicCategories`/`ListPublicTags` 传递 `isAuth` 参数 |
| `backend/internal/router/router.go` | 为 categories/tags 公开路由添加 `OptionalJWTAuth` 中间件 |

### 3.3 前端：code 类型主题缺乏前端路由拦截

**问题**：直接通过 URL 访问 code 类型主题的文档页，前端未拦截，内容直接加载。

**修复**：
| 文件 | 变更 |
|------|------|
| `vue/src/router/index.ts` | `resolveFirstPagePath` 中新增 access_mode 检查（login → 重定向登录，code → 重定向验证码页）；`beforeEach` 新增主题级访问控制（针对 DocPageView 等直接访问场景） |
| `vue/src/stores/reader.ts` | 新增 `getThemeAccessParams()` 辅助函数；`loadVersions`/`loadTree`/`loadPage`/`loadComments` 请求附带 `theme_token` 查询参数 |
| `vue/src/views/reader/ThemeVerifyCodePage.vue` | `onMounted` 中加载主题数据，已有有效 token 时直接跳转；修正 DTO 字段名 `access_token` |

---

## 四、接口清单

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| `GET` | `/public/tenants/:tenant_id/access` | 获取站点级访问设置 | 无 |
| `POST` | `/public/themes/:id/verify-code` | 验证主题访问码，返回短期 JWT | 无 |
| `PATCH` | `/tenant/settings/access` | 更新站点级访问设置 | JWTAuth |
| `GET` | `/public/tenants/:tenant_id/themes` | 主题列表 | OptionalJWT |
| `GET` | `/public/tenants/:tenant_id/themes/filter` | 主题筛选列表 | OptionalJWT |
| `GET` | `/public/tenants/:tenant_id/categories` | 分类树（含已发布主题计数） | OptionalJWT |
| `GET` | `/public/tenants/:tenant_id/tags` | 标签列表（含已发布主题计数） | OptionalJWT |
| `GET` | `/public/themes/:theme_id/versions` | 版本列表 | OptionalJWT |
| `GET` | `/public/versions/:version_id/tree` | 文档树 | OptionalJWT |
| `GET` | `/public/pages/:page_id` | 文档页内容 | OptionalJWT |
| `GET` | `/public/pages/:page_id/comments` | 评论列表 | OptionalJWT |

---

## 五、访问控制逻辑总结

### 站点级
- **维护模式**：`maintenance_mode = true` 时，未登录用户访问任何读者端页面均重定向到维护页
- **画廊登录可见**：`gallery_login_required = true` 时，未登录用户访问主题画廊页重定向到登录页

### 主题级
- **public**（默认）：所有人可见，所有公开接口正常返回
- **login**：登录用户可见，未登录用户在主题列表中不显示该主题，内容接口返回 `ErrThemeLoginRequired`
- **code**：主题列表中可见但无法直接访问，前端路由拦截至验证码页；验证成功后获取短期 token，后续请求通过 `theme_token` 参数携带；内容接口校验 token 有效性，无效返回 `ErrThemeCodeRequired`

### 前端拦截链路
1. `beforeEach` → 站点级检查（维护模式/画廊登录）→ 主题级检查（login/code）
2. `resolveFirstPagePath` → 主题级检查（ThemeHome/VersionHome beforeEnter）
3. `readerStore` → API 请求自动附带 `theme_token`（code/login 类型主题）
4. 后端接口 → 最终防线，即使前端绕过也能拦截

---

## 六、第三轮变更：Raw Markdown 访问控制 & 安全加固 & 性能优化

### 6.1 后端：Raw Markdown 访问控制

**问题**：`GetRawMarkdown` 对非 public 主题一刀切返回 403，导致有权限用户也无法使用 Raw 功能。

**修复**：
| 文件 | 变更 |
|------|------|
| `backend/internal/handler/public_handler.go` | `GetRawMarkdown` 改用 `CheckPublicAccess()` 做访问控制，支持 JWT 登录和 `theme_token` 查询参数 |
| `backend/internal/handler/public_handler.go` | **新增** `GetRawDirectory()` handler：返回主题版本下所有章节/页面的 raw 链接纯文本目录 |
| `backend/internal/handler/public_handler.go` | **新增** `IssueThemeAccessToken()` handler：已登录用户为指定主题签发独立的 `theme_access` token（24h 有效） |
| `backend/internal/router/router.go` | 新增路由 `GET /raw/:tenant_id/:theme_slug/:version`（目录）；raw 路由添加 `OptionalJWTAuth` 中间件；新增路由 `POST /public/themes/:theme_id/issue-token` |

### 6.2 后端：login 模式支持 theme_access token

**问题**：`CheckPublicAccess` 的 login 模式仅接受 JWT 认证，raw 链接无法通过独立 token 访问。

**修复**：
| 文件 | 变更 |
|------|------|
| `backend/internal/service/theme_service.go` | `CheckPublicAccess()` 的 `login` 分支新增 `themeAccessToken` 验证作为补充（与 code 模式对齐） |

### 6.3 后端：code 模式已登录用户免验证码

**问题**：`access_mode="code"` 的主题对已登录用户仍需验证码。

**修复**：
| 文件 | 变更 |
|------|------|
| `backend/internal/service/theme_service.go` | `CheckPublicAccess()` 的 `code` 分支：`isAuthenticated` 为 true 时直接放行 |

### 6.4 前端：Raw Markdown 二级下拉菜单

| 文件 | 变更 |
|------|------|
| `vue/src/components/preview/DocSidebar.vue` | 「查看原始 Markdown」按钮改为 `DropdownMenu`，包含：打开当前页 Raw / 复制链接 / 打开目录 / 复制目录链接；受保护主题自动附带 `theme_token`；底部显示链接有效期 |

### 6.5 前端：安全加固（避免完整 JWT 泄露）

**问题**：login 类型主题的 raw 链接若直接携带完整 JWT，泄露后等于泄露管理员权限。

**修复**：
| 文件 | 变更 |
|------|------|
| `vue/src/stores/reader.ts` | 新增 `ensureThemeAccessToken()` 方法：已登录用户进入受保护主题时，调用 `issue-token` 接口签发独立 token 并缓存到 localStorage；`getThemeAccessParams()` 扩展覆盖 login 类型主题 |
| `vue/src/router/index.ts` | `resolveFirstPagePath` 中 login/code 主题已登录时调用 `ensureThemeAccessToken()` |
| `vue/src/components/preview/DocSidebar.vue` | `getRawTokenSuffix()` 统一使用 localStorage 中的 `theme_access` token，不再传递完整 JWT |
| `vue/src/utils/http.ts` | `getAccessToken()` 保持为内部函数，不对外导出 |

### 6.6 前端：login 主题登录后重定向 & code 主题登录免验证

| 文件 | 变更 |
|------|------|
| `vue/src/components/auth/LoginForm.vue` | 登录成功后检查 `redirect` 查询参数，优先跳转到指定地址 |
| `vue/src/views/reader/ThemeLoginRequiredPage.vue` | 登录按钮携带 `redirect` 参数指向当前主题页 |
| `vue/src/router/index.ts` | `resolveFirstPagePath` 和 `beforeEach` 中 code 主题：已登录用户跳过 token 检查 |
| `vue/src/views/reader/ThemeVerifyCodePage.vue` | `onMounted` 中已登录用户直接跳转主题页 |

### 6.7 前端：DocReaderLayout 增强

| 文件 | 变更 |
|------|------|
| `vue/src/components/preview/DocReaderLayout.vue` | 导航栏右侧添加 `UserMenu`（已登录时显示）；主题名称后添加访问模式 `Badge`（非 public 时显示） |

### 6.8 性能优化：消除 API 重复调用

**问题**：每次导航到文档页，`themes`/`versions`/`tree` 各被请求 2 次（路由守卫 + DocPageViewer 组件各调用一次）。

**修复**：
| 文件 | 变更 |
|------|------|
| `vue/src/router/index.ts` | `resolveFirstPagePath` 中 `loadThemes` 改为 `if (store.themes.length === 0)` |
| `vue/src/components/preview/DocPageViewer.vue` | `loadData()` 中 themes/versions/tree 全部改为先检查 store 缓存，命中则跳过请求 |

---

## 七、新增接口清单（第三轮）

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| `GET` | `/raw/:tenant_id/:theme_slug/:version` | Raw 目录（纯文本，章节+页面链接） | OptionalJWT + theme_token |
| `GET` | `/raw/:tenant_id/:theme_slug/:version/:page_slug` | Raw 页面内容（纯文本 Markdown） | OptionalJWT + theme_token |
| `POST` | `/public/themes/:theme_id/issue-token` | 已登录用户签发主题访问 token | OptionalJWT（需已登录） |

---

## 八、安全设计

### Raw 链接 Token 策略
- **完整 JWT 永远不出现在 URL 中**：raw 链接统一使用独立的 `theme_access` token
- **theme_access token 特性**：24h 有效、仅授权特定主题只读访问、payload 仅含 `theme_id`
- **签发方式**：
  - `code` 主题：验证码通过后由 `verify-code` 接口签发
  - `login` 主题：已登录用户进入主题时由 `issue-token` 接口签发
- **存储**：`localStorage` key 为 `theme_access_{themeId}`，code/login 类型统一
