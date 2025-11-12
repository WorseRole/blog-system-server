package main

import (
	"blog-system-server/internal/config"
	"blog-system-server/internal/handlers"
	"blog-system-server/internal/middleware"
	"blog-system-server/internal/services"
	"blog-system-server/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// 加载配置
	cfg := config.Load()

	// 连接数据库
	if err := storage.InitDB(cfg); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化服务
	authService := services.NewAuthService(storage.DB)
	authHandler := handlers.NewAuthHandler(authService)

	postService := services.NewPostService(storage.DB)
	postHandler := handlers.NewPostHandler(postService)

	// 设置路由
	r := gin.Default()

	// 公开路由 不需要认证的路由
	// 注册路由
	r.POST("/api/auth/register", authHandler.Register)

	// 登录路由
	r.POST("/api/auth/login", authHandler.Login)

	// 探活 hello blog
	// r := gin.Default()
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
		authRoutes.POST("/post/create", postHandler.CreatePost)
		authRoutes.POST("/post/update", postHandler.UpdatePost)
		authRoutes.POST("/post/delete", postHandler.DeletePost)

	}

	// 启动服务
	log.Printf("服务器启动在 %s 这个端口", cfg.ServerPort)
	r.Run(cfg.ServerPort)

	// r.Run()
}
