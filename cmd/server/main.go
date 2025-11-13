package main

import (
	"blog-system-server/internal/config"
	"blog-system-server/internal/handlers"
	"blog-system-server/internal/middleware"
	"blog-system-server/internal/services"
	"blog-system-server/internal/storage"
	"blog-system-server/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// 加载配置
	cfg := config.Load()

	// 初始化日志系统
	if err := logger.InitLogger(cfg); err != nil {
		log.Fatal("日志初始化失败:", err)
	}
	defer logger.Logger.Sync()

	logger.Info("应用启动",
		zap.String("env", cfg.Env),
		zap.String("port", cfg.ServerPort),
	)

	// 连接数据库
	if err := storage.InitDB(cfg); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化服务
	authService := services.NewAuthService(storage.DB)
	authHandler := handlers.NewAuthHandler(authService)

	postService := services.NewPostService(storage.DB)
	postHandler := handlers.NewPostHandler(postService)

	commentService := services.NewCommentService(storage.DB)
	commentHandler := handlers.NewCommentHandler(commentService)

	// 设置路由
	r := gin.Default()

	// 注册访问日志中间件（所有请求都会记录）
	r.Use(middleware.AccessLog())

	// 公开路由 不需要认证的路由
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	// 探活 hello blog
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello blog",
		})
	})

	// 需要认证的路由组
	authRoutes := r.Group("/api")
	// 应用JWT组件
	authRoutes.Use(middleware.AuthMiddleware())
	{
		// 这里添加需要认证的路由 测试
		// authRoutes.GET/POST

		// 文章
		authRoutes.POST("/post/create", postHandler.CreatePost)
		authRoutes.POST("/post/update", postHandler.UpdatePost)
		authRoutes.POST("/post/delete", postHandler.DeletePost)
		authRoutes.GET("/post/lists", postHandler.SelectPostsLists)
		authRoutes.GET("/post/detail", postHandler.SelectPostDetial)

		// 评论
		authRoutes.POST("/comment/create", commentHandler.CreateComment)
		authRoutes.GET("/comment/lists", commentHandler.SelectCommentsListByPostID)

	}

	// 启动服务
	logger.Logger.Info("服务器启动成功", zap.String("port", cfg.ServerPort))
	r.Run(cfg.ServerPort)
}
