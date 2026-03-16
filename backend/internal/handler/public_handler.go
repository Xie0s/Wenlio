// Package handler HTTP 请求处理器
//
// 职责：处理公开接口（读者端，无需认证）
// 对外接口：PublicHandler
package handler

import (
	"docplatform/internal/middleware"
	"docplatform/internal/model/dto"
	"docplatform/internal/service"
	"docplatform/pkg/errcode"
	"docplatform/pkg/jwt"
	"docplatform/pkg/response"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PublicHandler 公开接口处理器
type PublicHandler struct {
	tenantService   *service.TenantService
	themeService    *service.ThemeService
	versionService  *service.VersionService
	pageService     *service.PageService
	commentService  *service.CommentService
	searchService   *service.SearchService
	categoryService *service.CategoryService
	tagService      *service.TagService
	settingsService *service.TenantSettingsService
}

// extractThemeAccessToken 从请求中提取主题访问 token
// 优先级：X-Theme-Access 头 > theme_token 查询参数
func extractThemeAccessToken(c *gin.Context) string {
	if t := c.GetHeader("X-Theme-Access"); t != "" {
		return t
	}
	return c.Query("theme_token")
}

// NewPublicHandler 创建 PublicHandler
func NewPublicHandler() *PublicHandler {
	return &PublicHandler{
		tenantService:   service.NewTenantService(),
		themeService:    service.NewThemeService(),
		versionService:  service.NewVersionService(),
		pageService:     service.NewPageService(),
		commentService:  service.NewCommentService(),
		searchService:   service.NewSearchService(),
		categoryService: service.NewCategoryService(),
		tagService:      service.NewTagService(),
		settingsService: service.NewTenantSettingsService(),
	}
}

// GetTenant 获取租户基础信息
func (h *PublicHandler) GetTenant(c *gin.Context) {
	tenant, appErr := h.tenantService.GetPublicByID(c.Request.Context(), c.Param("tenant_id"))
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tenant)
}

// ListThemes 租户已发布主题列表（草稿主题不可见）
// 通过 OptionalJWTAuth 中间件判断登录状态：已登录返回全部主题，未登录隐藏 access_mode="login" 的主题
func (h *PublicHandler) ListThemes(c *gin.Context) {
	isAuth := middleware.GetAuthContext(c) != nil
	themes, appErr := h.themeService.ListPublic(c.Request.Context(), c.Param("tenant_id"), isAuth)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, themes)
}

// ListVersions 主题下已发布/已归档版本列表
func (h *PublicHandler) ListVersions(c *gin.Context) {
	themeID, err := primitive.ObjectIDFromHex(c.Param("theme_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	// 访问控制检查
	theme, appErr := h.themeService.GetByIDUnsafe(c.Request.Context(), themeID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	isAuth := middleware.GetAuthContext(c) != nil
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, extractThemeAccessToken(c)); accessErr != nil {
		response.Fail(c, accessErr)
		return
	}
	versions, appErr := h.versionService.ListPublishedByTheme(c.Request.Context(), themeID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, versions)
}

// GetVersionTree 版本文档树（侧边栏数据）
func (h *PublicHandler) GetVersionTree(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("version_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	// 回溯到主题做访问控制
	version, appErr := h.versionService.GetByID(c.Request.Context(), versionID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	theme, appErr := h.themeService.GetByIDUnsafe(c.Request.Context(), version.ThemeID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	isAuth := middleware.GetAuthContext(c) != nil
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, extractThemeAccessToken(c)); accessErr != nil {
		response.Fail(c, accessErr)
		return
	}
	tree, appErr := h.pageService.GetVersionTree(c.Request.Context(), versionID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, gin.H{"sections": tree})
}

// GetPage 获取单页 Markdown 内容（仅 published）
func (h *PublicHandler) GetPage(c *gin.Context) {
	pageID, err := primitive.ObjectIDFromHex(c.Param("page_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	page, appErr := h.pageService.GetPublished(c.Request.Context(), pageID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	// 回溯到主题做访问控制
	version, vErr := h.versionService.GetByID(c.Request.Context(), page.VersionID)
	if vErr != nil {
		response.Fail(c, vErr)
		return
	}
	theme, tErr := h.themeService.GetByIDUnsafe(c.Request.Context(), version.ThemeID)
	if tErr != nil {
		response.Fail(c, tErr)
		return
	}
	isAuth := middleware.GetAuthContext(c) != nil
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, extractThemeAccessToken(c)); accessErr != nil {
		response.Fail(c, accessErr)
		return
	}
	response.Success(c, page)
}

// ListComments 页面已审核评论列表
func (h *PublicHandler) ListComments(c *gin.Context) {
	pageID, err := primitive.ObjectIDFromHex(c.Param("page_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	// 回溯到主题做访问控制（page → version → theme）
	pageEntity, pErr := h.pageService.GetPublished(c.Request.Context(), pageID)
	if pErr != nil {
		response.Fail(c, pErr)
		return
	}
	version, vErr := h.versionService.GetByID(c.Request.Context(), pageEntity.VersionID)
	if vErr != nil {
		response.Fail(c, vErr)
		return
	}
	theme, tErr := h.themeService.GetByIDUnsafe(c.Request.Context(), version.ThemeID)
	if tErr != nil {
		response.Fail(c, tErr)
		return
	}
	isAuth := middleware.GetAuthContext(c) != nil
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, extractThemeAccessToken(c)); accessErr != nil {
		response.Fail(c, accessErr)
		return
	}
	comments, appErr := h.commentService.ListByPage(c.Request.Context(), pageID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, comments)
}

// SubmitComment 提交评论
func (h *PublicHandler) SubmitComment(c *gin.Context) {
	pageID, err := primitive.ObjectIDFromHex(c.Param("page_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.SubmitCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	comment, appErr := h.commentService.Submit(c.Request.Context(), pageID, &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.SuccessCreate(c, gin.H{"message": "评论已提交，等待审核", "id": comment.ID.Hex()})
}

// GetTenantHomepage 获取租户已发布首页配置
func (h *PublicHandler) GetTenantHomepage(c *gin.Context) {
	published, appErr := h.tenantService.GetHomepagePublished(c.Request.Context(), c.Param("tenant_id"))
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, dto.HomepagePublicResp{Published: published})
}

// Search 全文搜索文档（公开接口，必须携带 tenant_id）
func (h *PublicHandler) Search(c *gin.Context) {
	var req dto.SearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	resp, appErr := h.searchService.Search(c.Request.Context(), &req)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, resp)
}

// GetRawMarkdown 获取页面原始 Markdown（AI 可直接采集）
// URL: /raw/:tenant_id/:theme_slug/:version/:page_slug
// 响应：text/plain; charset=utf-8
func (h *PublicHandler) GetRawMarkdown(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	themeSlug := c.Param("theme_slug")
	versionName := c.Param("version")
	pageSlug := c.Param("page_slug")

	// 1. 查找主题
	theme, err := h.themeService.FindBySlug(c.Request.Context(), tenantID, themeSlug)
	if err != nil {
		c.String(404, "# 404\n主题不存在")
		return
	}

	// 访问控制：复用统一检查（JWT 登录 / theme_token 验证码）
	isAuth := middleware.GetAuthContext(c) != nil
	themeToken := extractThemeAccessToken(c)
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, themeToken); accessErr != nil {
		c.String(403, "# 403\n该主题受访问控制保护，请在链接中附带 theme_token 参数")
		return
	}

	// 2. 查找版本（支持 name 或 label 匹配）
	version, verErr := h.versionService.FindByThemeAndName(c.Request.Context(), theme.ID, versionName)
	if verErr != nil {
		c.String(404, "# 404\n版本不存在")
		return
	}

	// 3. 查找已发布页面
	page, pageErr := h.pageService.GetPublishedBySlug(c.Request.Context(), version.ID, pageSlug)
	if pageErr != nil {
		c.String(404, "# 404\n页面不存在或未发布")
		return
	}

	// 返回纯文本 Markdown
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.String(200, page.Content)
}

// GetRawDirectory 获取主题版本的原始 Markdown 目录（所有章节和页面的 raw 链接）
// URL: /raw/:tenant_id/:theme_slug/:version
// 响应：text/plain; charset=utf-8
func (h *PublicHandler) GetRawDirectory(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	themeSlug := c.Param("theme_slug")
	versionName := c.Param("version")

	// 1. 查找主题
	theme, err := h.themeService.FindBySlug(c.Request.Context(), tenantID, themeSlug)
	if err != nil {
		c.String(404, "# 404\n主题不存在")
		return
	}

	// 访问控制
	isAuth := middleware.GetAuthContext(c) != nil
	themeToken := extractThemeAccessToken(c)
	if accessErr := h.themeService.CheckPublicAccess(theme, isAuth, themeToken); accessErr != nil {
		c.String(403, "# 403\n该主题受访问控制保护，请在链接中附带 theme_token 参数")
		return
	}

	// 2. 查找版本
	version, verErr := h.versionService.FindByThemeAndName(c.Request.Context(), theme.ID, versionName)
	if verErr != nil {
		c.String(404, "# 404\n版本不存在")
		return
	}

	// 3. 获取文档树
	tree, appErr := h.pageService.GetVersionTree(c.Request.Context(), version.ID)
	if appErr != nil {
		c.String(500, "# 500\n获取文档树失败")
		return
	}

	// 4. 构建纯文本目录
	basePath := fmt.Sprintf("/raw/%s/%s/%s", tenantID, themeSlug, versionName)
	// 如果有 theme_token，附加到每个链接
	tokenSuffix := ""
	if qt := c.Query("theme_token"); qt != "" {
		tokenSuffix = "?theme_token=" + qt
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s - %s\n\n", theme.Name, version.Label))
	for _, sec := range tree {
		sb.WriteString(fmt.Sprintf("## %s\n", sec.Title))
		for _, p := range sec.Pages {
			sb.WriteString(fmt.Sprintf("- [%s](%s/%s%s)\n", p.Title, basePath, p.Slug, tokenSuffix))
		}
		sb.WriteString("\n")
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.String(200, sb.String())
}

// ListPublicCategories 获取租户分类树（公开，含已发布主题计数）
func (h *PublicHandler) ListPublicCategories(c *gin.Context) {
	isAuth := middleware.GetAuthContext(c) != nil
	tree, appErr := h.categoryService.ListPublic(c.Request.Context(), c.Param("tenant_id"), isAuth)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tree)
}

// ListPublicTags 获取租户标签列表（公开，含已发布主题计数，自动过滤未使用标签）
func (h *PublicHandler) ListPublicTags(c *gin.Context) {
	isAuth := middleware.GetAuthContext(c) != nil
	tags, appErr := h.tagService.ListPublic(c.Request.Context(), c.Param("tenant_id"), isAuth)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, tags)
}

// ListThemesByFilter 支持分类/标签筛选的主题列表（公开，草稿主题不可见）
// 通过 OptionalJWTAuth 中间件判断登录状态
func (h *PublicHandler) ListThemesByFilter(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	categorySlug := c.Query("category")
	tagSlugs := c.QueryArray("tag")
	isAuth := middleware.GetAuthContext(c) != nil
	themes, appErr := h.themeService.ListByFilterPublic(c.Request.Context(), tenantID, categorySlug, tagSlugs, isAuth)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, themes)
}

// GetTenantAccess 获取租户站点级访问控制设置（公开接口，前端用于维护模式/画廊登录判断）
func (h *PublicHandler) GetTenantAccess(c *gin.Context) {
	access, appErr := h.settingsService.GetAccessSettings(c.Request.Context(), c.Param("tenant_id"))
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	response.Success(c, dto.AccessSettingsResp{
		MaintenanceMode:      access.MaintenanceMode,
		GalleryLoginRequired: access.GalleryLoginRequired,
	})
}

// IssueThemeAccessToken 已登录用户为指定主题签发 theme_access token
// 用途：login 类型主题的 raw markdown 链接分享，避免泄露完整 JWT
func (h *PublicHandler) IssueThemeAccessToken(c *gin.Context) {
	if middleware.GetAuthContext(c) == nil {
		response.Fail(c, errcode.ErrUnauthorized)
		return
	}
	themeID, err := primitive.ObjectIDFromHex(c.Param("theme_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	// 验证主题存在
	_, appErr := h.themeService.GetByIDUnsafe(c.Request.Context(), themeID)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	token, tokenErr := jwt.GenerateThemeAccessToken(themeID.Hex())
	if tokenErr != nil {
		response.Fail(c, errcode.ErrInternalServer)
		return
	}
	response.Success(c, dto.VerifyAccessCodeResp{
		AccessToken: token,
		ExpiresIn:   86400,
	})
}

// VerifyThemeAccessCode 验证主题访问码（公开接口）
// 验证通过后签发短期 JWT（24小时有效），payload 包含 theme_id
func (h *PublicHandler) VerifyThemeAccessCode(c *gin.Context) {
	themeID, err := primitive.ObjectIDFromHex(c.Param("theme_id"))
	if err != nil {
		response.Fail(c, errcode.ErrInvalidParam)
		return
	}
	var req dto.VerifyAccessCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	theme, appErr := h.themeService.VerifyAccessCode(c.Request.Context(), themeID, req.Code)
	if appErr != nil {
		response.Fail(c, appErr)
		return
	}
	// 签发主题访问 token（复用 JWT 基础设施，role="theme_access"，tenantID=themeID）
	token, tokenErr := jwt.GenerateThemeAccessToken(theme.ID.Hex())
	if tokenErr != nil {
		response.Fail(c, errcode.ErrInternalServer)
		return
	}
	response.Success(c, dto.VerifyAccessCodeResp{
		AccessToken: token,
		ExpiresIn:   86400,
	})
}
