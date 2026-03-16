// Package errcode 业务错误码定义
//
// 职责：统一定义所有业务错误码常量及 AppError 构造函数
// 对外接口：New() 创建 AppError，预定义各模块错误常量
package errcode

// AppError 业务错误
type AppError struct {
	HTTPCode int    `json:"-"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New 创建 AppError
func New(httpCode, code int, message string) *AppError {
	return &AppError{HTTPCode: httpCode, Code: code, Message: message}
}

// Wrap 包装内部错误（用于日志记录，不暴露给客户端）
func (e *AppError) Wrap(internal error) *AppError {
	return &AppError{
		HTTPCode: e.HTTPCode,
		Code:     e.Code,
		Message:  e.Message,
		Internal: internal,
	}
}

// ============================================================
// 通用错误（模块 00）
// ============================================================

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

// ============================================================
// 认证错误（模块 01）
// ============================================================

var (
	ErrLoginFailed        = New(401, 401101, "用户名或密码错误")
	ErrAccountLocked      = New(401, 401102, "账号已锁定，请稍后重试")
	ErrAccountDisabled    = New(401, 401103, "账号已被禁用")
	ErrPasswordTooWeak    = New(400, 400104, "密码强度不足，至少8位且包含字母和数字")
	ErrPasswordIncorrect  = New(422, 422105, "当前密码错误")
	ErrCaptchaRequired    = New(400, 400106, "请先完成安全检查")
	ErrCaptchaInvalid     = New(400, 400107, "安全检查未通过，请重试")
	ErrCaptchaExpired     = New(400, 400108, "安全检查已过期，请重新验证")
	ErrCaptchaTooFrequent = New(429, 429109, "安全检查请求过于频繁，请稍后再试")
)

// ============================================================
// 租户错误（模块 02）
// ============================================================

var (
	ErrTenantNotFound      = New(404, 404201, "租户不存在")
	ErrTenantIDExists      = New(409, 409202, "租户 ID 已存在")
	ErrTenantSuspended     = New(422, 422203, "租户已被封禁")
	ErrTenantAlreadyActive = New(422, 422204, "租户当前已是活跃状态")
	ErrTenantIDInvalid     = New(400, 400205, "租户 ID 格式不合法")
	ErrTenantIDReserved    = New(409, 409206, "租户 ID 为系统保留词")
	ErrTenantDeleting      = New(422, 422207, "租户正在删除中，请稍后重试")
)

// ============================================================
// 首页个性化错误（模块 02 续）
// ============================================================

var (
	ErrHomepageNoDraft = New(422, 422208, "没有待发布的首页草稿")
)

// ============================================================
// 版本错误（模块 03）
// ============================================================

var (
	ErrVersionNotFound         = New(404, 404301, "版本不存在")
	ErrVersionNotDraft         = New(422, 422302, "版本不是草稿状态")
	ErrVersionNotPublished     = New(422, 422303, "版本不是发布状态")
	ErrVersionArchived         = New(422, 422304, "版本已归档，不可修改")
	ErrVersionDefaultRequired  = New(422, 422305, "不可取消唯一的默认版本")
	ErrVersionSetDefaultNotPub = New(422, 422306, "仅已发布的版本可设为默认")
	ErrVersionDefaultCannotDel = New(422, 422307, "默认版本不可删除，请先切换默认版本")
	ErrVersionAlreadyDraft     = New(422, 422308, "版本已是草稿状态")
	ErrVersionNotArchived      = New(422, 422309, "版本不是归档状态")
)

// ============================================================
// 主题错误（模块 04）
// ============================================================

var (
	ErrThemeNotFound    = New(404, 404401, "主题不存在")
	ErrThemeSlugExists  = New(409, 409402, "主题 Slug 已存在")
	ErrThemeHasVersions = New(422, 422403, "主题下存在版本，无法删除")
	ErrThemeDeleting    = New(422, 422404, "主题正在删除中，请稍后重试")
)

// ============================================================
// 章节错误（模块 05）
// ============================================================

var (
	ErrSectionNotFound        = New(404, 404501, "章节不存在")
	ErrSectionVersionArchived = New(422, 422502, "所属版本已归档，不可操作章节")
)

// ============================================================
// 文档页错误（模块 06）
// ============================================================

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

// ============================================================
// 评论错误（模块 07）
// ============================================================

var (
	ErrCommentNotFound         = New(404, 404701, "评论不存在")
	ErrCommentPageNotPublished = New(422, 422702, "文档页未发布，不可提交评论")
	ErrCommentAlreadyApproved  = New(422, 422703, "评论已批准")
	ErrCommentAlreadyRejected  = New(422, 422704, "评论已拒绝")
	ErrCommentContentTooLong   = New(400, 400705, "评论内容超过 1000 字限制")
	ErrCommentAuthorRequired   = New(400, 400706, "评论作者昵称为必填项")
	ErrCommentNestedTooDeep    = New(422, 422707, "仅支持一层嵌套回复")
)

// ============================================================
// 用户错误（模块 08）
// ============================================================

var (
	ErrUserNotFound      = New(404, 404801, "用户不存在")
	ErrUsernameExists    = New(409, 409802, "用户名已存在")
	ErrUserDisabled      = New(422, 422803, "用户已被禁用")
	ErrCannotDisableSelf = New(422, 422804, "不可禁用自身账号")
	ErrUserAlreadyActive = New(422, 422805, "用户当前已是活跃状态")
	ErrCannotDeleteSelf  = New(422, 422806, "不可删除自身账号")
)

// ============================================================
// 媒体/上传错误（模块 09）
// ============================================================

var (
	ErrUploadFileTooLarge     = New(400, 400901, "文件大小超出限制")
	ErrUploadFileTypeInvalid  = New(400, 400902, "不支持的文件类型")
	ErrUploadStorageFail      = New(500, 500903, "文件存储失败")
	ErrLocalStorageFull       = New(422, 422904, "本地存储空间已满（100MB 限制），请前往设置开启云存储")
	ErrCloudStorageNotEnabled = New(422, 422905, "云存储未启用，请先在设置中配置云存储")
	ErrS3ConnectionFailed     = New(422, 422906, "云存储连接失败，请检查配置")
	ErrMediaStillInUse        = New(409, 409907, "该媒体文件仍被文档引用，无法删除")
)

// ============================================================
// 分类错误（模块 10）
// ============================================================

var (
	ErrCategoryNotFound      = New(404, 404001, "分类不存在")
	ErrCategorySlugExists    = New(409, 409002, "分类 Slug 已存在")
	ErrCategoryHasChildren   = New(422, 422003, "分类下存在子分类，无法删除")
	ErrCategoryHasThemes     = New(422, 422004, "分类下存在主题，无法删除")
	ErrCategoryMaxDepth      = New(422, 422005, "仅支持二级分类，不可继续添加子分类")
	ErrThemeCategoryRequired = New(422, 422006, "创建主题必须选择分类")
)

// ============================================================
// 标签错误（模块 11）
// ============================================================

var (
	ErrTagNotFound   = New(404, 404101, "标签不存在")
	ErrTagSlugExists = New(409, 409102, "标签 Slug 已存在")
	ErrTagInUse      = New(409, 409103, "标签已被主题使用，请确认后继续")
)

// ============================================================
// 访问控制错误（模块 12）
// ============================================================

var (
	ErrAccessCodeInvalid = New(403, 403201, "验证码错误")
	ErrThemeLoginRequired = New(401, 401202, "该主题需要登录后访问")
	ErrThemeCodeRequired  = New(403, 403203, "该主题需要验证码才能访问")
)
