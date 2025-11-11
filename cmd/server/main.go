package main

import (
	"blog-system-server/internal/config"
	"blog-system-server/internal/handlers"
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

	// 设置路由
	r := gin.Default()

	// 注册路由
	r.POST("/api/auth/register", authHandler.Register)

	// 测试服务 hello blog
	// r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello blog",
		})
	})

	// 启动服务
	log.Printf("服务器启动在 %s 这个端口", cfg.ServerPort)
	r.Run(cfg.ServerPort)

	// r.Run()
}
