// Package router 路由注册
//
// 职责：注册所有 API 路由，绑定 Handler 和中间件
// 对外接口：RegisterRoutes()
package router

import (
	"docplatform/internal/handler"
	"docplatform/internal/middleware"
	"docplatform/pkg/constants"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	authHandler := handler.NewAuthHandler()
	tenantHandler := handler.NewTenantHandler()
	userHandler := handler.NewUserHandler()
	themeHandler := handler.NewThemeHandler()
	versionHandler := handler.NewVersionHandler()
	sectionHandler := handler.NewSectionHandler()
	pageHandler := handler.NewPageHandler()
	commentHandler := handler.NewCommentHandler()
	mediaHandler := handler.NewMediaHandler()
	publicHandler := handler.NewPublicHandler()
	homepageHandler := handler.NewHomepageHandler()
	settingsHandler := handler.NewTenantSettingsHandler()
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()

	v1 := r.Group("/api/v1")

	// ============================================================
	// 认证路由（部分公开，部分需认证）
	// ============================================================
	auth := v1.Group("/auth")
	{
		auth.POST("/captcha/challenge", authHandler.CreateCaptchaChallenge)
		auth.POST("/captcha/verify", authHandler.VerifyCaptcha)
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		// 以下接口需要认证
		authed := auth.Group("")
		authed.Use(middleware.JWTAuth())
		authed.POST("/logout", authHandler.Logout)
		authed.GET("/me", authHandler.Me)
		authed.PATCH("/me/password", authHandler.ChangePassword)
		authed.PATCH("/me/profile", authHandler.UpdateProfile)
		authed.POST("/me/avatar", authHandler.UploadAvatar)
		authed.POST("/me/profile-bg", authHandler.UploadProfileBg)
	}

	// ============================================================
	// 超级管理员路由
	// ============================================================
	admin := v1.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.RoleGuard(constants.RoleSuperAdmin))
	{
		// 租户管理
		admin.GET("/tenants", tenantHandler.List)
		admin.POST("/tenants", tenantHandler.Create)
		admin.GET("/tenants/:id", tenantHandler.Get)
		admin.PATCH("/tenants/:id", tenantHandler.Update)
		admin.DELETE("/tenants/:id", tenantHandler.Delete)
		admin.POST("/tenants/:id/suspend", tenantHandler.Suspend)
		admin.POST("/tenants/:id/activate", tenantHandler.Activate)

		// 全平台用户管理
		admin.GET("/users", userHandler.ListAll)
		admin.POST("/users", userHandler.CreateSuperAdmin)
		admin.PATCH("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)
		admin.POST("/users/:id/deactivate", userHandler.Deactivate)
		admin.POST("/users/:id/activate", userHandler.Activate)
		admin.POST("/users/:id/reset-password", userHandler.ResetPassword)
	}

	// ============================================================
	// 租户管理员路由
	// ============================================================
	tenant := v1.Group("/tenant")
	tenant.Use(middleware.JWTAuth(), middleware.RoleGuard(constants.RoleTenantAdmin), middleware.TenantIsolation())
	{
		// 租户用户管理
		tenant.GET("/users", userHandler.ListByTenant)
		tenant.POST("/users", userHandler.CreateTenantAdmin)
		tenant.PATCH("/users/:id", userHandler.Update)
		tenant.DELETE("/users/:id", userHandler.Delete)
		tenant.POST("/users/:id/deactivate", userHandler.Deactivate)
		tenant.POST("/users/:id/activate", userHandler.Activate)
		tenant.POST("/users/:id/reset-password", userHandler.ResetPassword)

		// 主题管理
		tenant.GET("/themes", themeHandler.List)
		tenant.POST("/themes", themeHandler.Create)
		tenant.GET("/themes/:id", themeHandler.Get)
		tenant.PATCH("/themes/:id", themeHandler.Update)
		tenant.DELETE("/themes/:id", themeHandler.Delete)
		tenant.PUT("/themes/sort", themeHandler.Sort)

		// 版本管理
		tenant.GET("/themes/:id/versions", versionHandler.List)
		tenant.POST("/themes/:id/versions", versionHandler.Create)
		tenant.GET("/versions/:id", versionHandler.Get)
		tenant.PATCH("/versions/:id", versionHandler.Update)
		tenant.POST("/versions/:id/publish", versionHandler.Publish)
		tenant.POST("/versions/:id/unpublish", versionHandler.Unpublish)
		tenant.POST("/versions/:id/archive", versionHandler.Archive)
		tenant.POST("/versions/:id/unarchive", versionHandler.Unarchive)
		tenant.POST("/versions/:id/set-default", versionHandler.SetDefault)
		tenant.POST("/versions/:id/clone", versionHandler.Clone)
		tenant.DELETE("/versions/:id", versionHandler.Delete)

		// 章节管理
		tenant.GET("/versions/:id/sections", sectionHandler.ListSections)
		tenant.POST("/versions/:id/sections", sectionHandler.CreateSection)
		tenant.PATCH("/sections/:id", sectionHandler.UpdateSection)
		tenant.DELETE("/sections/:id", sectionHandler.DeleteSection)
		tenant.PUT("/versions/:id/sections/sort", sectionHandler.SortSections)

		// 文档页管理
		tenant.GET("/sections/:id/pages", pageHandler.ListPages)
		tenant.POST("/sections/:id/pages", pageHandler.CreatePage)
		tenant.GET("/pages/:id", pageHandler.GetPage)
		tenant.PUT("/pages/:id", pageHandler.UpdatePage)
		tenant.PATCH("/pages/:id", pageHandler.PatchPage)
		tenant.DELETE("/pages/:id", pageHandler.DeletePage)
		tenant.POST("/pages/:id/publish", pageHandler.PublishPage)
		tenant.POST("/pages/:id/unpublish", pageHandler.UnpublishPage)
		tenant.PUT("/sections/:id/pages/sort", pageHandler.SortPages)

		// 评论管理
		tenant.GET("/comments", commentHandler.List)
		tenant.POST("/comments/:id/approve", commentHandler.Approve)
		tenant.POST("/comments/:id/reject", commentHandler.Reject)
		tenant.DELETE("/comments/:id", commentHandler.Delete)
		tenant.POST("/comments/:id/reply", commentHandler.Reply)

		// 媒体文件管理
		tenant.POST("/media/upload", mediaHandler.Upload)
		tenant.GET("/media", mediaHandler.List)
		tenant.DELETE("/media/:id", mediaHandler.Delete)
		tenant.GET("/media/audit", mediaHandler.Audit)
		tenant.GET("/media/usage", mediaHandler.Usage)
		tenant.POST("/media/cleanup-unused", mediaHandler.CleanupUnused)
		tenant.POST("/media/orphan/delete", mediaHandler.DeleteOrphan)

		// 首页个性化
		tenant.GET("/homepage", homepageHandler.GetDraft)
		tenant.PUT("/homepage", homepageHandler.SaveDraft)
		tenant.POST("/homepage/publish", homepageHandler.Publish)

		// 租户设置
		tenant.GET("/settings", settingsHandler.Get)
		tenant.PATCH("/settings/storage", settingsHandler.UpdateStorage)
		tenant.POST("/settings/storage/test", settingsHandler.TestS3)
		tenant.GET("/settings/storage/usage", settingsHandler.GetStorageUsage)
		tenant.PATCH("/settings/ai", settingsHandler.UpdateAI)
		tenant.POST("/settings/ai/test", settingsHandler.TestAI)
		tenant.PATCH("/settings/access", settingsHandler.UpdateAccess)

		// 分类管理
		tenant.GET("/categories", categoryHandler.List)
		tenant.POST("/categories", categoryHandler.Create)
		tenant.PATCH("/categories/:id", categoryHandler.Update)
		tenant.DELETE("/categories/:id", categoryHandler.Delete)
		tenant.PUT("/categories/sort", categoryHandler.Sort)

		// 标签管理
		tenant.GET("/tags", tagHandler.List)
		tenant.POST("/tags", tagHandler.Create)
		tenant.PATCH("/tags/:id", tagHandler.Update)
		tenant.DELETE("/tags/:id", tagHandler.Delete)
	}

	// ============================================================
	// 公开路由（无需认证，面向读者）
	// ============================================================
	public := v1.Group("/public")
	{
		public.GET("/tenants/:tenant_id", publicHandler.GetTenant)
		public.GET("/tenants/:tenant_id/homepage", publicHandler.GetTenantHomepage)
		public.GET("/tenants/:tenant_id/access", publicHandler.GetTenantAccess)                                         // 站点级访问控制设置
		public.GET("/tenants/:tenant_id/themes", middleware.OptionalJWTAuth(), publicHandler.ListThemes)                // 可选JWT：按登录状态过滤主题
		public.GET("/tenants/:tenant_id/themes/filter", middleware.OptionalJWTAuth(), publicHandler.ListThemesByFilter) // 可选JWT：同上
		public.GET("/tenants/:tenant_id/categories", middleware.OptionalJWTAuth(), publicHandler.ListPublicCategories)
		public.GET("/tenants/:tenant_id/tags", middleware.OptionalJWTAuth(), publicHandler.ListPublicTags)
		public.GET("/themes/:theme_id/versions", middleware.OptionalJWTAuth(), publicHandler.ListVersions)
		public.GET("/versions/:version_id/tree", middleware.OptionalJWTAuth(), publicHandler.GetVersionTree)
		public.GET("/pages/:page_id", middleware.OptionalJWTAuth(), publicHandler.GetPage)
		public.GET("/pages/:page_id/comments", middleware.OptionalJWTAuth(), publicHandler.ListComments)
		public.POST("/pages/:page_id/comments", middleware.OptionalJWTAuth(), publicHandler.SubmitComment)
		public.POST("/themes/:theme_id/verify-code", middleware.IPRateLimit(10, 1*time.Minute), publicHandler.VerifyThemeAccessCode)                               // 主题访问码校验（限速：每分钟10次/IP）
		public.POST("/themes/:theme_id/issue-token", middleware.OptionalJWTAuth(), publicHandler.IssueThemeAccessToken) // 已登录用户签发主题访问 token
		public.GET("/search", publicHandler.Search)
	}

	// Raw Markdown（AI 友好，纯文本响应）
	// GET /raw/:tenant_id/:theme_slug/:version           → 目录
	// GET /raw/:tenant_id/:theme_slug/:version/:page_slug → 页面内容
	r.GET("/raw/:tenant_id/:theme_slug/:version", middleware.OptionalJWTAuth(), publicHandler.GetRawDirectory)
	r.GET("/raw/:tenant_id/:theme_slug/:version/:page_slug", middleware.OptionalJWTAuth(), publicHandler.GetRawMarkdown)
}
