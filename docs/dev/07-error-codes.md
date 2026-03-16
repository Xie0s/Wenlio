# 07 - 错误码与异常处理

> **文档说明**：定义系统所有错误码、编码规则与 Go 实现规范，作为错误处理的唯一口径。

---

## 1. 错误码编码规则

错误码格式：`XXYYZZ`

| 段   | 位数 | 说明                              |
| ---- | ---- | --------------------------------- |
| `XX` | 2位  | HTTP 状态码（40=400, 42=422, 50=500 等） |
| `YY` | 2位  | 模块编号（见下表）                |
| `ZZ` | 2位  | 模块内错误序号，从 01 开始        |

**示例**：`422301` = HTTP 422 + 模块 03（版本） + 错误 01

---

## 2. 模块编号分配

| 模块编号 | 模块名称     |
| -------- | ------------ |
| `00`     | 通用/系统    |
| `01`     | 认证         |
| `02`     | 租户         |
| `03`     | 版本         |
| `04`     | 主题         |
| `05`     | 章节         |
| `06`     | 文档页       |
| `07`     | 评论         |
| `08`     | 用户         |
| `09`     | 媒体/上传    |

---

## 3. 完整错误码表

### 3.1 通用错误（模块 00）

| 错误码   | HTTP | 常量名                    | 说明                   |
| -------- | ---- | ------------------------- | ---------------------- |
| `400001` | 400  | `ErrInvalidParam`         | 请求参数格式错误       |
| `400002` | 400  | `ErrMissingParam`         | 缺少必填参数           |
| `400003` | 400  | `ErrInvalidPageParam`     | 分页参数非法           |
| `401001` | 401  | `ErrUnauthorized`         | 未认证，请先登录       |
| `401002` | 401  | `ErrTokenExpired`         | Token 已过期           |
| `401003` | 401  | `ErrTokenInvalid`         | Token 无效             |
| `403001` | 403  | `ErrForbidden`            | 权限不足               |
| `403002` | 403  | `ErrTenantMismatch`       | 数据不属于当前租户     |
| `404001` | 404  | `ErrResourceNotFound`     | 资源不存在             |
| `500001` | 500  | `ErrInternalServer`       | 服务器内部错误         |
| `500002` | 500  | `ErrDatabase`             | 数据库操作失败         |

---

### 3.2 认证错误（模块 01）

| 错误码   | HTTP | 常量名                       | 说明                     |
| -------- | ---- | ---------------------------- | ------------------------ |
| `401101` | 401  | `ErrLoginFailed`             | 用户名或密码错误         |
| `401102` | 401  | `ErrAccountLocked`           | 账号已锁定，请稍后重试   |
| `401103` | 401  | `ErrAccountDisabled`         | 账号已被禁用             |
| `400104` | 400  | `ErrPasswordTooWeak`         | 密码强度不足（至少8位，含字母和数字） |
| `422105` | 422  | `ErrPasswordIncorrect`       | 当前密码错误（修改密码时）|

---

### 3.3 租户错误（模块 02）

| 错误码   | HTTP | 常量名                      | 说明                               |
| -------- | ---- | --------------------------- | ---------------------------------- |
| `404201` | 404  | `ErrTenantNotFound`         | 租户不存在                         |
| `409202` | 409  | `ErrTenantIDExists`         | 租户 ID 已存在                     |
| `422203` | 422  | `ErrTenantSuspended`        | 租户已被封禁                       |
| `422204` | 422  | `ErrTenantAlreadyActive`    | 租户当前已是活跃状态               |
| `400205` | 400  | `ErrTenantIDInvalid`        | 租户 ID 格式不合法（需小写字母/数字/连字符，3-32位） |
| `409206` | 409  | `ErrTenantIDReserved`       | 租户 ID 为系统保留词，不可使用     |

---

### 3.4 版本错误（模块 03）

| 错误码   | HTTP | 常量名                       | 说明                                  |
| -------- | ---- | ---------------------------- | ------------------------------------- |
| `404301` | 404  | `ErrVersionNotFound`         | 版本不存在                            |
| `422302` | 422  | `ErrVersionNotDraft`         | 版本不是草稿状态，无法执行此操作      |
| `422303` | 422  | `ErrVersionNotPublished`     | 版本不是发布状态，无法归档            |
| `422304` | 422  | `ErrVersionArchived`         | 版本已归档，不可修改                  |
| `422305` | 422  | `ErrVersionDefaultRequired`  | 不可取消唯一的默认版本                |
| `422306` | 422  | `ErrVersionSetDefaultNotPub` | 仅已发布的版本可设为默认              |

---

### 3.5 主题错误（模块 04）

| 错误码   | HTTP | 常量名                    | 说明                   |
| -------- | ---- | ------------------------- | ---------------------- |
| `404401` | 404  | `ErrThemeNotFound`        | 主题不存在             |
| `409402` | 409  | `ErrThemeSlugExists`      | 主题 Slug 已存在       |
| `422403` | 422  | `ErrThemeHasVersions`     | 主题下存在版本，无法删除 |

---

### 3.6 章节错误（模块 05）

| 错误码   | HTTP | 常量名                     | 说明                            |
| -------- | ---- | -------------------------- | ------------------------------- |
| `404501` | 404  | `ErrSectionNotFound`       | 章节不存在                      |
| `422502` | 422  | `ErrSectionVersionArchived`| 所属版本已归档，不可操作章节    |

---

### 3.7 文档页错误（模块 06）

| 错误码   | HTTP | 常量名                         | 说明                                     |
| -------- | ---- | ------------------------------ | ---------------------------------------- |
| `404601` | 404  | `ErrPageNotFound`              | 文档页不存在                             |
| `409602` | 409  | `ErrPageSlugExists`            | 同版本内 Slug 已存在                     |
| `422603` | 422  | `ErrPageAlreadyPublished`      | 文档页已处于发布状态                     |
| `422604` | 422  | `ErrPageAlreadyDraft`          | 文档页已处于草稿状态                     |
| `422605` | 422  | `ErrPageVersionNotPublished`   | 文档页所属版本未发布，不可单独发布文档页 |
| `422606` | 422  | `ErrPageVersionArchived`       | 所属版本已归档，不可操作文档页           |
| `400607` | 400  | `ErrPageContentEmpty`          | 文档内容不可为空                         |
| `400608` | 400  | `ErrPageSlugInvalid`           | Slug 格式不合法                          |
| `400609` | 400  | `ErrPageImportInvalidFile`     | 导入文件格式不支持（仅支持 .md）         |

---

### 3.8 评论错误（模块 07）

| 错误码   | HTTP | 常量名                        | 说明                                 |
| -------- | ---- | ----------------------------- | ------------------------------------ |
| `404701` | 404  | `ErrCommentNotFound`          | 评论不存在                           |
| `422702` | 422  | `ErrCommentPageNotPublished`  | 文档页未发布，不可提交评论           |
| `422703` | 422  | `ErrCommentAlreadyApproved`   | 评论已批准                           |
| `422704` | 422  | `ErrCommentAlreadyRejected`   | 评论已拒绝                           |
| `400705` | 400  | `ErrCommentContentTooLong`    | 评论内容超过 1000 字限制             |
| `400706` | 400  | `ErrCommentAuthorRequired`    | 评论作者昵称为必填项                 |
| `422707` | 422  | `ErrCommentNestedTooDeep`     | 仅支持一层嵌套回复                   |

---

### 3.9 用户错误（模块 08）

| 错误码   | HTTP | 常量名                      | 说明                          |
| -------- | ---- | --------------------------- | ----------------------------- |
| `404801` | 404  | `ErrUserNotFound`           | 用户不存在                    |
| `409802` | 409  | `ErrUsernameExists`         | 用户名已存在                  |
| `422803` | 422  | `ErrUserDisabled`           | 用户已被禁用                  |
| `422804` | 422  | `ErrCannotDisableSelf`      | 不可禁用自身账号              |
| `422805` | 422  | `ErrUserAlreadyActive`      | 用户当前已是活跃状态          |

---

### 3.10 媒体/上传错误（模块 09）

| 错误码   | HTTP | 常量名                     | 说明                          |
| -------- | ---- | -------------------------- | ----------------------------- |
| `400901` | 400  | `ErrUploadFileTooLarge`    | 文件大小超出限制（最大 10MB） |
| `400902` | 400  | `ErrUploadFileTypeInvalid` | 不支持的文件类型              |
| `500903` | 500  | `ErrUploadStorageFail`     | 文件存储失败                  |

---

## 4. Go 实现

### 4.1 错误码常量定义

```go
// Package errcode 业务错误码定义
//
// 职责：统一定义所有业务错误码常量及 AppError 构造函数
// 对外接口：New() 创建 AppError，预定义各模块错误常量
package errcode

import "net/http"

// AppError 业务错误
type AppError struct {
    HTTPCode int
    Code     int
    Message  string
    Internal error
}

func (e *AppError) Error() string {
    return e.Message
}

// New 创建 AppError
func New(httpCode, code int, message string) *AppError {
    return &AppError{HTTPCode: httpCode, Code: code, Message: message}
}

// Wrap 包装内部错误
func (e *AppError) Wrap(internal error) *AppError {
    return &AppError{
        HTTPCode: e.HTTPCode,
        Code:     e.Code,
        Message:  e.Message,
        Internal: internal,
    }
}

// 通用错误
var (
    ErrInvalidParam     = New(400, 400001, "请求参数错误")
    ErrMissingParam     = New(400, 400002, "缺少必填参数")
    ErrInvalidPageParam = New(400, 400003, "分页参数非法")
    ErrUnauthorized     = New(401, 401001, "请先登录")
    ErrTokenExpired     = New(401, 401002, "登录已过期，请重新登录")
    ErrTokenInvalid     = New(401, 401003, "Token 无效")
    ErrForbidden        = New(403, 403001, "权限不足")
    ErrTenantMismatch   = New(403, 403002, "数据不属于当前租户")
    ErrResourceNotFound = New(404, 404001, "资源不存在")
    ErrInternalServer   = New(500, 500001, "服务器内部错误")
    ErrDatabase         = New(500, 500002, "数据库操作失败")
)

// 认证错误
var (
    ErrLoginFailed       = New(401, 401101, "用户名或密码错误")
    ErrAccountLocked     = New(401, 401102, "账号已锁定，请稍后重试")
    ErrAccountDisabled   = New(401, 401103, "账号已被禁用")
    ErrPasswordTooWeak   = New(400, 400104, "密码强度不足，至少8位且包含字母和数字")
    ErrPasswordIncorrect = New(422, 422105, "当前密码错误")
)

// 租户错误
var (
    ErrTenantNotFound      = New(404, 404201, "租户不存在")
    ErrTenantIDExists      = New(409, 409202, "租户 ID 已存在")
    ErrTenantSuspended     = New(422, 422203, "租户已被封禁")
    ErrTenantAlreadyActive = New(422, 422204, "租户当前已是活跃状态")
    ErrTenantIDInvalid     = New(400, 400205, "租户 ID 格式不合法")
    ErrTenantIDReserved    = New(409, 409206, "租户 ID 为系统保留词")
)

// 版本错误
var (
    ErrVersionNotFound         = New(404, 404301, "版本不存在")
    ErrVersionNotDraft         = New(422, 422302, "版本不是草稿状态")
    ErrVersionNotPublished     = New(422, 422303, "版本不是发布状态")
    ErrVersionArchived         = New(422, 422304, "版本已归档，不可修改")
    ErrVersionDefaultRequired  = New(422, 422305, "不可取消唯一的默认版本")
    ErrVersionSetDefaultNotPub = New(422, 422306, "仅已发布的版本可设为默认")
)

// 主题错误
var (
    ErrThemeNotFound    = New(404, 404401, "主题不存在")
    ErrThemeSlugExists  = New(409, 409402, "主题 Slug 已存在")
    ErrThemeHasVersions = New(422, 422403, "主题下存在版本，无法删除")
)

// 章节错误
var (
    ErrSectionNotFound        = New(404, 404501, "章节不存在")
    ErrSectionVersionArchived = New(422, 422502, "所属版本已归档，不可操作章节")
)

// 文档页错误
var (
    ErrPageNotFound            = New(404, 404601, "文档页不存在")
    ErrPageSlugExists          = New(409, 409602, "同版本内 Slug 已存在")
    ErrPageAlreadyPublished    = New(422, 422603, "文档页已处于发布状态")
    ErrPageAlreadyDraft        = New(422, 422604, "文档页已处于草稿状态")
    ErrPageVersionNotPublished = New(422, 422605, "所属版本未发布，不可单独发布文档页")
    ErrPageVersionArchived     = New(422, 422606, "所属版本已归档，不可操作文档页")
    ErrPageContentEmpty        = New(400, 400607, "文档内容不可为空")
    ErrPageSlugInvalid         = New(400, 400608, "Slug 格式不合法")
    ErrPageImportInvalidFile   = New(400, 400609, "仅支持导入 .md 文件")
)

// 评论错误
var (
    ErrCommentNotFound         = New(404, 404701, "评论不存在")
    ErrCommentPageNotPublished = New(422, 422702, "文档页未发布，不可提交评论")
    ErrCommentAlreadyApproved  = New(422, 422703, "评论已批准")
    ErrCommentAlreadyRejected  = New(422, 422704, "评论已拒绝")
    ErrCommentContentTooLong   = New(400, 400705, "评论内容超过 1000 字限制")
    ErrCommentAuthorRequired   = New(400, 400706, "评论作者昵称为必填项")
    ErrCommentNestedTooDeep    = New(422, 422707, "仅支持一层嵌套回复")
)

// 用户错误
var (
    ErrUserNotFound      = New(404, 404801, "用户不存在")
    ErrUsernameExists    = New(409, 409802, "用户名已存在")
    ErrUserDisabled      = New(422, 422803, "用户已被禁用")
    ErrCannotDisableSelf = New(422, 422804, "不可禁用自身账号")
    ErrUserAlreadyActive = New(422, 422805, "用户当前已是活跃状态")
)

// 媒体/上传错误
var (
    ErrUploadFileTooLarge    = New(400, 400901, "文件大小超出限制")
    ErrUploadFileTypeInvalid = New(400, 400902, "不支持的文件类型")
    ErrUploadStorageFail     = New(500, 500903, "文件存储失败")
)
```

### 4.2 错误响应处理

```go
// Fail 统一错误响应
func Fail(c *gin.Context, err *AppError) {
    c.JSON(err.HTTPCode, Response{
        Code:      err.Code,
        Message:   err.Message,
        Data:      nil,
        RequestID: c.GetString("request_id"),
    })
}
```

### 4.3 参数校验错误转换

```go
// HandleValidationError 将 validator 错误转换为统一格式
func HandleValidationError(c *gin.Context, err error) {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        fields := make([]string, 0, len(ve))
        for _, e := range ve {
            fields = append(fields, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
        }
        c.JSON(http.StatusBadRequest, Response{
            Code:      400001,
            Message:   "参数校验失败: " + strings.Join(fields, "; "),
            RequestID: c.GetString("request_id"),
        })
        return
    }
    Fail(c, ErrInvalidParam)
}
```

---

## 5. 异常处理策略

### 5.1 Recovery 中间件

```go
// Recovery panic 恢复中间件，防止单个请求 panic 导致服务崩溃
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                logger.Error("panic recovered",
                    zap.Any("error", r),
                    zap.String("request_id", c.GetString("request_id")),
                    zap.Stack("stack"),
                )
                Fail(c, ErrInternalServer)
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

### 5.2 错误日志分级策略

| 错误类型      | 日志级别 | 说明                           |
| ------------- | -------- | ------------------------------ |
| 400 参数错误  | `WARN`   | 记录请求参数（不含敏感字段）   |
| 401/403 认证  | `WARN`   | 记录 IP、路径、失败原因        |
| 404 不存在    | `INFO`   | 不需要特别关注                 |
| 422 业务错误  | `INFO`   | 记录错误码和业务上下文         |
| 500 服务器错误| `ERROR`  | 记录完整堆栈，立即告警         |

---

**文档版本**：v1.0
**适用项目**：文档管理平台
**最后更新**：2026 年 2 月
