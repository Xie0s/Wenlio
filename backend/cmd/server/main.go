// Package main 程序入口
//
// 职责：加载配置、初始化日志/数据库、注册路由、启动 HTTP 服务
// 对外接口：main()
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"docplatform/internal/config"
	"docplatform/internal/middleware"
	"docplatform/internal/repository"
	"docplatform/internal/router"
	"docplatform/internal/service"
	"docplatform/pkg/constants"
	"docplatform/pkg/logger"
	mongopkg "docplatform/pkg/mongo"
	"docplatform/pkg/utils"

	"docplatform/internal/model/entity"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化日志
	logger.Init(cfg.Server.Env)
	defer logger.Sync()

	// 3. 连接 MongoDB
	if err := mongopkg.Connect(
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
		cfg.MongoDB.Database,
		cfg.MongoDB.Username,
		cfg.MongoDB.Password,
		cfg.MongoDB.AuthDatabase,
	); err != nil {
		logger.L().Fatal("MongoDB 连接失败", zap.Error(err))
	}
	defer mongopkg.Close()

	// 4. 初始化索引
	ensureIndexes()

	// 5. 初始化种子数据
	seedSuperAdmin(cfg)
	recoverDeletingTenants()

	// 6. 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 7. 创建 Gin 引擎并注册全局中间件
	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(),
		middleware.Recovery(),
		cacheControl(),
	)

	// 8. 静态文件服务（上传文件访问）
	r.Static("/uploads", cfg.Storage.LocalPath)

	// 9. 注册路由
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.RegisterRoutes(r)

	// 10. 前端 SPA 静态资源服务
	serveFrontend(r)

	// 11. 启动 HTTP 服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		logger.L().Info("HTTP 服务启动",
			zap.String("addr", addr),
			zap.String("url", fmt.Sprintf("http://localhost:%d", cfg.Server.Port)),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L().Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 10. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.L().Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Error("服务关闭异常", zap.Error(err))
	}
	logger.L().Info("服务已停止")
}

func cacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		switch {
		case strings.HasPrefix(path, "/api/"):
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		case path == "/sw.js", path == "/version.json", path == "/index.html":
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// ensureIndexes 初始化所有集合索引
func ensureIndexes() {
	ctx := context.Background()
	repos := []interface{ EnsureIndexes(context.Context) error }{
		repository.NewUserRepo(),
		repository.NewTenantRepo(),
		repository.NewThemeRepo(),
		repository.NewVersionRepo(),
		repository.NewSectionRepo(),
		repository.NewPageRepo(),
		repository.NewCommentRepo(),
		repository.NewMediaRepo(),
		repository.NewCategoryRepo(),
		repository.NewTagRepo(),
		repository.NewTenantSettingsRepo(),
	}
	for _, repo := range repos {
		if err := repo.EnsureIndexes(ctx); err != nil {
			logger.L().Warn("索引创建失败", zap.Error(err))
		}
	}
	logger.L().Info("数据库索引初始化完成")
}

func recoverDeletingTenants() {
	tenantService := service.NewTenantService()
	if appErr := tenantService.RecoverDeleting(context.Background()); appErr != nil {
		logger.L().Warn("恢复删除中租户任务失败",
			zap.Int("code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Error(appErr.Internal),
		)
		return
	}
	logger.L().Info("删除中租户任务恢复检查完成")
}

// serveFrontend 挂载前端 SPA 静态资源
// 从 frontend 目录提供静态文件，非 API/uploads 路径回退到 index.html 支持客户端路由
func serveFrontend(r *gin.Engine) {
	distPath := filepath.Join(".", "dist")

	// 检查 dist 目录是否存在
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		logger.L().Warn("前端构建产物不存在，跳过静态资源挂载", zap.String("path", distPath))
		return
	}

	staticFS := http.Dir(distPath)
	fileServer := http.FileServer(staticFS)

	// 读取 index.html 内容用于 SPA fallback
	indexPath := filepath.Join(distPath, "index.html")
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		logger.L().Warn("读取 index.html 失败，跳过 SPA 回退", zap.Error(err))
		return
	}

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 和 uploads 路径不处理
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 尝试提供静态文件
		if f, err := fs.Stat(os.DirFS(distPath), strings.TrimPrefix(path, "/")); err == nil && !f.IsDir() {
			if path == "/sw.js" || path == "/version.json" || path == "/index.html" {
				c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Header("Pragma", "no-cache")
				c.Header("Expires", "0")
			}
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// SPA fallback：返回 index.html
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexBytes)
	})

	logger.L().Info("前端静态资源已挂载", zap.String("path", distPath))
}

// seedSuperAdmin 初始化超级管理员（幂等，已存在则跳过）
func seedSuperAdmin(cfg *config.Config) {
	ctx := context.Background()
	userRepo := repository.NewUserRepo()

	existing, _ := userRepo.FindByUsername(ctx, cfg.Seed.SuperAdminUsername)
	if existing != nil {
		return
	}

	hashed, err := utils.HashPassword(cfg.Seed.SuperAdminPassword)
	if err != nil {
		logger.L().Error("密码哈希失败", zap.Error(err))
		return
	}

	user := &entity.User{
		TenantID: "",
		Username: cfg.Seed.SuperAdminUsername,
		Password: hashed,
		Name:     "超级管理员",
		Role:     constants.RoleSuperAdmin,
		Status:   constants.UserStatusActive,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		logger.L().Error("创建超级管理员失败", zap.Error(err))
		return
	}
	logger.L().Info("超级管理员创建成功", zap.String("username", cfg.Seed.SuperAdminUsername))
}
