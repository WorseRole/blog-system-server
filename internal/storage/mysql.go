package storage

import (
	"blog-system-server/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 数据库初始化连接
func InitDB() error {
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败: %v", err)
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	DB = db
	return nil
}

// 迁移表结构
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 自动迁移表结构 初始化建表
	err := DB.AutoMigrate(
		&models.Users{},
		&models.Posts{},
		&models.Comments{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表结构迁移完成")
	return nil
}
