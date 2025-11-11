package main

import (
	"blog-system-server/internal/config"
	"blog-system-server/internal/storage"
	"log"
)

func main() {

	// 初始化数据库连接
	cfg := config.Load()
	if err := storage.InitDB(cfg); err != nil {
		// log的Fatal / Fatalf / Fatalln 	输出日志后 调用 os.Exit(1) 直接退出程序。
		log.Fatal("数据库连接失败:", err)
	}

	// 执行迁移
	if err := storage.AutoMigrate(); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 完成迁移
	log.Println("数据库初始化完成！")

}
