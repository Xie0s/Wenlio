# 05 - 认证与权限设计

> **文档说明**：定义系统认证机制、JWT 策略、角色权限模型与多租户数据隔离实现，作为安全层开发的唯一口径。

---

## 1. 认证机制

### 1.1 JWT 认证流程

```
┌──────────┐     POST /api/v1/auth/login      ┌──────────┐
│  客户端   │ ──────────────────────────────▶  │  API     │
│          │  { username, password }           │  Server  │
│          │ ◀──────────────────────────────── │          │
│          │  { access_token, expires_in }     │          │
└──────────┘                                   └──────────┘

后续请求：
Authorization: Bearer <access_token>
```

### 1.2 Token 策略

| 项目     | Access Token          |
| -------- | --------------------- |
| 签名算法 | HS256                 |
| 有效期   | 24 小时               |
| 用途     | API 请求认证          |
| 存储位置 | 客户端 localStorage   |
| 密钥存储 | config.yaml `jwt.secret` |

> **简化设计**：本系统不使用 Refresh Token 机制，Token 过期后重新登录即可。管理后台为低频操作场景，24 小时有效期已足够。

### 1.3 Token Payload 结构

```json
{
  "sub": "用户ID (ObjectID hex)",
  "tid": "租户ID (字符串，如 'acme'，super_admin 为空字符串)",
  "role": "super_admin | tenant_admin",
  "iat": 1740000000,
  "exp": 1740086400
}
```

> **注意**：读者（viewer）无需登录，不产生 Token。Token 仅用于管理后台操作。

### 1.4 退出登录

```
客户端行为：
1. 清除 localStorage 中的 access_token
2. 跳转至登录页

服务端行为：
- 无需服务端操作（无 Token 黑名单机制）
- Token 自然过期后失效
```

> **安全说明**：HS256 + 无黑名单方案适用于本系统的场景（管理员数量有限、无高安全要求）。如未来需要即时吊销能力，可引入 Redis 黑名单。

---

## 2. 角色模型

### 2.1 角色定义

| 角色标识       | 角色名称   | 层级   | 数据范围       |
| -------------- | ---------- | ------ | -------------- |
| `super_admin`  | 超级管理员 | 平台层 | 全部租户       |
| `tenant_admin` | 租户管理员 | 租户层 | 本租户全部数据 |
| `viewer`       | 读者       | 公开层 | 已发布文档     |

### 2.2 角色常量

```go
const (
    RoleSuperAdmin  = "super_admin"
    RoleTenantAdmin = "tenant_admin"
)

// 读者不需要角色常量，因为不需要登录
```

### 2.3 角色说明

- **super_admin**：全平台唯一（或极少数），管理所有租户和账号，不参与文档管理业务
- **tenant_admin**：管理本租户的主题、版本、文档、评论，可邀请其他 tenant_admin
- **viewer**：匿名读者，无需登录即可浏览已发布文档和提交评论

---

## 3. 权限矩阵

### 3.1 接口级权限

| 功能模块           | 接口路径模式                     | super_admin | tenant_admin | viewer（匿名）|
| ------------------ | -------------------------------- | ----------- | ------------ | -------------- |
| 租户管理           | `/api/v1/admin/tenants/**`       | ✅          | ❌           | ❌             |
| 超管用户管理       | `/api/v1/admin/users/**`         | ✅          | ❌           | ❌             |
| 租户用户管理       | `/api/v1/tenant/users/**`        | ❌          | ✅           | ❌             |
| 主题管理           | `/api/v1/tenant/themes/**`       | ❌          | ✅           | ❌             |
| 版本管理           | `/api/v1/tenant/versions/**`     | ❌          | ✅           | ❌             |
| 章节管理           | `/api/v1/tenant/sections/**`     | ❌          | ✅           | ❌             |
| 文档页管理         | `/api/v1/tenant/pages/**`        | ❌          | ✅           | ❌             |
| 评论审核           | `/api/v1/tenant/comments/**`     | ❌          | ✅           | ❌             |
| 媒体上传           | `/api/v1/tenant/media/**`        | ❌          | ✅           | ❌             |
| 阅读已发布文档     | `/api/v1/public/**`              | 无需认证    | 无需认证     | ✅             |
| 提交评论           | `POST /api/v1/public/.../comments` | 无需认证  | 无需认证     | ✅             |
| 全文搜索           | `GET /api/v1/public/search`      | 无需认证    | 无需认证     | ✅             |

---

## 4. 数据隔离中间件实现

### 4.1 TenantIsolation 中间件

```go
// TenantIsolation 租户数据隔离中间件
// 从 JWT 中提取 tenant_id 并注入 Gin Context，
// 后续所有 Repository 查询自动附加 tenant_id 过滤条件
func TenantIsolation() gin.HandlerFunc {
    return func(c *gin.Context) {
        authCtx, exists := c.Get("auth_context")
        if !exists {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        auth := authCtx.(*AuthContext)
        // super_admin 不注入 tenant_id（平台管理接口专属）
        if auth.Role == RoleSuperAdmin {
            c.Next()
            return
        }
        if auth.TenantID == "" {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        c.Set("tenant_id", auth.TenantID)
        c.Next()
    }
}
```

### 4.2 RoleGuard 中间件

```go
// RoleGuard 角色校验中间件工厂
// 传入允许的角色列表，不匹配则返回 403
func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authCtx, exists := c.Get("auth_context")
        if !exists {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        auth := authCtx.(*AuthContext)
        for _, role := range allowedRoles {
            if auth.Role == role {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "code":    403001,
            "message": "权限不足",
        })
    }
}
```

### 4.3 路由分组与中间件挂载

```go
func RegisterRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")

    // 认证路由（部分公开，部分需认证）
    auth := v1.Group("/auth")
    {
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.Refresh)
        auth.Use(JWTAuth())
        auth.POST("/logout", authHandler.Logout)
        auth.GET("/me", authHandler.Me)
        auth.PATCH("/me/password", authHandler.ChangePassword)
    }

    // 超级管理员路由
    admin := v1.Group("/admin")
    admin.Use(JWTAuth(), RoleGuard(RoleSuperAdmin))
    {
        admin.GET("/tenants", tenantHandler.List)
        admin.POST("/tenants", tenantHandler.Create)
        // ...
    }

    // 租户管理员路由
    tenant := v1.Group("/tenant")
    tenant.Use(JWTAuth(), RoleGuard(RoleTenantAdmin), TenantIsolation())
    {
        tenant.GET("/themes", themeHandler.List)
        tenant.POST("/themes", themeHandler.Create)
        // ...
    }

    // 公开路由（无需认证）
    public := v1.Group("/public")
    {
        public.GET("/tenants/:tenant_id", publicHandler.GetTenant)
        public.GET("/tenants/:tenant_id/themes", publicHandler.ListThemes)
        // ...
    }
}
```

---

## 5. 登录安全策略

| 规则             | 说明                                        |
| ---------------- | ------------------------------------------- |
| 密码哈希         | bcrypt cost ≥ 12                            |
| 登录失败锁定     | 连续失败 5 次，锁定 15 分钟                 |
| 锁定存储         | `login_fail_count` + `locked_until` 字段   |
| 密码强度要求     | 最少 8 位，含字母和数字                     |
| 登录响应         | 不区分"用户不存在"和"密码错误"，统一返回"用户名或密码错误" |

### 登录流程伪代码

```go
func (s *AuthService) Login(username, password string) (*dto.LoginResp, error) {
    // 1. 查询用户
    user, err := s.userRepo.FindByUsername(username)
    if err != nil {
        return nil, errcode.ErrLoginFailed
    }

    // 2. 检查账号锁定
    if time.Now().Before(user.LockedUntil) {
        return nil, errcode.ErrAccountLocked
    }

    // 3. 检查账号状态
    if user.Status != "active" {
        return nil, errcode.ErrAccountDisabled
    }

    // 4. 校验密码
    if !crypto.CheckPassword(password, user.Password) {
        // 累加失败次数
        user.LoginFailCount++
        if user.LoginFailCount >= 5 {
            user.LockedUntil = time.Now().Add(15 * time.Minute)
        }
        s.userRepo.UpdateLoginFail(user)
        return nil, errcode.ErrLoginFailed
    }

    // 5. 登录成功，重置失败计数
    s.userRepo.ResetLoginFail(user.ID)

    // 6. 签发 Token
    token, err := jwt.GenerateToken(user.ID, user.TenantID, user.Role)
    if err != nil {
        return nil, errcode.ErrInternalServer
    }

    return &dto.LoginResp{
        AccessToken: token,
        ExpiresIn:   86400,
        User: dto.UserInfo{
            ID:       user.ID,
            Name:     user.Name,
            Role:     user.Role,
            TenantID: user.TenantID,
        },
    }, nil
}
```

---

## 6. 接口安全分级

| 分级     | 说明                                          | 示例                                     |
| -------- | --------------------------------------------- | ---------------------------------------- |
| 公开     | 无需任何认证                                  | `/api/v1/public/pages/{id}`              |
| 已认证   | 需要有效 Access Token，无角色限制             | `/api/v1/auth/me`                        |
| 角色限制 | 需要特定角色                                  | `/api/v1/admin/tenants`（super_admin）   |
| 租户隔离 | 认证后通过 `tenant_id` 强制过滤数据           | `/api/v1/tenant/themes`（tenant_admin）  |

---

## 7. 公开接口安全说明

公开接口（`/api/v1/public/*`）面向匿名读者，安全设计：

| 措施               | 说明                                               |
| ------------------ | -------------------------------------------------- |
| 只读原则           | GET 接口仅返回已发布（published）数据              |
| 租户参数强制       | 搜索接口必须携带 `tenant_id` 参数，服务端校验存在性|
| 评论提交限制       | 评论内容限 1000 字，默认 `pending` 状态需审核      |
| 无敏感数据泄露     | 公开接口不返回 draft 内容、管理配置等内部数据      |

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
