package main

import (
	"blog-system-server/internal/config"
	"blog-system-server/pkg/utils"
	"fmt"
	"log"
)

func main() {
	// r := gin.Default()
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "hello world!",
	// 	})
	// })
	// // 默认在 localhost:8080
	// r.Run()

	salt, err := utils.GenerateSalt()
	if err != nil {
		fmt.Printf("创建失败:%v", err)
	}
	fmt.Printf("salt:%s", salt)

	// 加载配置
	cfg := config.Load()
	log.Printf("=== 配置验证 ===")
	log.Printf("MYSQLHost: %s", cfg.MYSQLHost)
	log.Printf("MYSQLPort: %s", cfg.MYSQLPort)
	log.Printf("MYSQLUser: %s", cfg.MYSQLUser)
	log.Printf("MYSQLPassword: %s", cfg.MYSQLPassword)
	log.Printf("MYSQLDatabase: %s", cfg.MYSQLDatabase)
	log.Printf("ServerPort: %s", cfg.ServerPort)
	log.Printf("Env: %s", cfg.Env)
}
