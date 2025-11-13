package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

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
	LogLevel      string
}

func Load() *Config {

	// 获取项目根目录的绝对路径
	projectRoot := getProjectRoot()
	envPath := filepath.Join(projectRoot, ".env")

	// 加载.env文件
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("未找到.env文件, 使用环境变量")
	// }
	// 加载指定路径的.env
	if err := godotenv.Load(envPath); err != nil {
		log.Println("未找到.env文件, 使用环境变量")
	} else {
		log.Println("成功找到.env文件:", envPath)
	}

	return &Config{
		MYSQLHost: getEnv("MYSQL_HOST", "localhost"),
		MYSQLPort: getEnv("MYSQL_PORT", "3306"),
		MYSQLUser: getEnv("MYSQL_USER", "root"),
		// password 和 database 隐藏 一下下
		MYSQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MYSQLDatabase: getEnv("MYSQL_DATABASE", "blog_system"),
		ServerPort:    getEnv("SERVER_PORT", ":8080"),
		Env:           getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

/**
默认确实基于程序的当前工作目录（Current Working Directory）来查找 .env 文件。
这个行为在你的项目里得到了验证：当你在项目根目录运行 main.go 时，工作目录就是项目根目录，所以能找到 .env；
而运行 scripts/migrate.go 时，如果工作目录不是项目根目录，就找不到。
**/
// 获取项目根目录  因为.env 在根目录
func getProjectRoot() string {

	// 获取当前文件的路径(config.go的路径)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 根据项目结构向上找到根目录
	// internal/config/config.go -> internal/config -> internal -> 项目根目录
	return filepath.Join(dir, "..", "..")

}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
