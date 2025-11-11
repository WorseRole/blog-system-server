### 先搭建一个Go服务

##### 1. 创建并初始化项目：

~~~bash
mkdir blog-system-server
cd blog-system-server
go mod init blog-system-server
~~~

##### 2. 安装依赖： 

使用go get 添加 Gin 框架

~~~bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

go get github.com/joho/godotenv
~~~

##### 3. 编写测试代码

创建main.go 文件，写入一个简单的测试用例

~~~go
package main

import "github.com/gin-gonic/gin"

func main() {
  r := gin.Default()
  r.Get("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hello World"
    })
  })
  // 默认在127.0.0.1 8080
  r.Run()
}
~~~



### 设计表结构

~~~sql
-- users 表
id int 	// 主键
name varchar	// 姓名
password varchar	// 密码
salt varchar	// 盐
email varchar	// 邮箱
is_del int	// 0: 正常，1:禁用
created_at datetime	// 创建时间
updated_at datetime	// 更新时间


-- posts 表
id int	// 主键
title varchar	// 标题
content text	// 内容
user_id int	// 对应userId 文章的作者
is_del int 	// 0:未删除，1:已删除
created_at datetime	// 创建时间
updated_at datetime	// 更新时间


--comments 表
id int // 主键
content text	// 评论内容
user_id int	// 对应userId 评论的作者
post_id int	// 对应postId 评论在此文章下
is_del int	// 0:未删除，1:已删除
created_at	// 创建时间
updated_at	// 更新时间

~~~





### 规划项目结构

项目结构：

~~~text
blog-system-server/
├── cmd/
│   └── server/
│       └── main.go              # 入口(依赖注入)
├── internal/
│   ├── config/                  # 配置
│   ├── models/                  # GORM模型
│   ├── handlers/                # HTTP处理
│   │   ├── auth.go
│   │   ├── post.go  
│   │   └── comment.go
│   ├── services/                # 业务逻辑
│   │   ├── auth_service.go
│   │   ├── post_service.go
│   │   └── comment_service.go
│   ├── middleware/              # 中间件
│   │   ├── auth.go
│   │   └── cors.go
│   └── storage/                 # 数据存储
│       ├── mysql.go
│       └── redis.go
├── pkg/
│   ├── utils/
│   │   ├── jwt.go
│   │   └── password.go
│   └── response/
│       └── response.go          # 统一响应
└── scripts/
    └── migrate.sql              # 数据库初始化
~~~



