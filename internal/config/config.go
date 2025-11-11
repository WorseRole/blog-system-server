package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MYSQLHost     string
	MYSQLPort     string
	MYSQLUser     string
	MYSQLPassword string
	MYSQLDatabase string
	ServerPort    string
	Env           string
}

func Load() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件, 使用环境变量")
	}
	return &Config{
		MYSQLHost: getEnv("MYSQL_HOST", "localhost"),
		MYSQLPort: getEnv("MYSQL_PORT", "3306"),
		MYSQLUser: getEnv("MYSQL_USER", "root"),
		// password 和 database 隐藏 一下下
		MYSQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MYSQLDatabase: getEnv("MYSQL_DATABASE", ""),
	}
}

func getEnv(key string, value string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return value
}
