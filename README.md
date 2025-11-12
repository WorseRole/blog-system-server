### 先搭建一个Go服务

##### 1. 创建并初始化项目：

~~~bash
mkdir blog-system-server
cd blog-system-server
go mod init blog-system-server
~~~

##### 2. 安装依赖： 

使用go get 添加 Gin 框架等包

~~~bash
// gin框架包
go get -u github.com/gin-gonic/gin
// gorm包
go get -u gorm.io/gorm
// mysql驱动包
go get -u gorm.io/driver/mysql
// 这个是读配置文件的包
go get github.com/joho/godotenv
// 这个是jwt的包
go get -u github.com/golang-jwt/jwt/v4
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

- **Service层**：处理纯业务逻辑，不知道HTTP细节
- **Handler层**：处理HTTP协议，不知道业务细节

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



在目录下创建 `.env` 配置文件，在`config.go` 中读取并配置

~~~.env
# mysql 配置
MYSQL_HOST = localhost
MYSQL_PORT = 3306
MYSQL_USER = root
MYSQL_PASSWORD = xxxyyyzzz
MYSQL_DATABASE = xxxyyyzzz

# 服务器配置
SERVER_PORT = 8080
ENV=development

# JWT配置

# Redis配置
~~~





### 接口文档

#### 1. 用户相关

##### 1.1. 注册接口	✅	已完成

>  接口说明：用户进行注册账号，创建成功后需要重新进行登录，需要账号，密码，邮箱

>  是否需要Token：否

URL：

~~~http
POST  http://localhost:8080/api/auth/register
~~~

传参：

~~~json
{
    "username":"leoYan",
    "password":"xxx***yyy",
    "email":"leoyan9527@outlook.com"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "注册成功,请登录",
    "data": {
        "id": 1,
        "username": "leoYan",
        "email": "leoyan9527@outlook.com"
    }
}
~~~

##### 1.2. 登录接口	✅	已完成

>  接口说明：用户输入账号、密码进行登录，得到JWT Token 认证

>  是否需要Token：否

URL:

~~~http
POST	http://localhost:8080/api/auth/login
~~~

传参：

~~~json
{
    "username":"leoYan",
    "password":"xxx***yyy"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzA2MTA0MCwibmJmIjoxNzYyOTc0NjQwLCJpYXQiOjE3NjI5NzQ2NDB9.7hjKcZ-qExWIbIrCWuFpXyB_bxSHe5OmQd5UaQ9wNHY",
        "id": 1,
        "username": "leoYan",
        "email": "leoyan9527@outlook.com"
    }
}
~~~

#### 2. 文章相关

##### 2.1. 创建post接口	✅	已完成

>  接口说明：创建文章，需要已经登录认证的用户才可以创建文章，创建文章时需要 title和content都不能为空

>  是否需要Token：是

URL:

~~~http
POST  http://localhost:8080/api/post/create
~~~

传参：

~~~json
{
    "title":"GoLang基础学习",
    "content":"1.需要认证学习打牢基础; 2.需要进行实践用起来"
}
~~~

返回结果

~~~json
{
    "code": 200,
    "message": "创建文章成功",
    "data": {
        "id": 1,
        "title": "GoLang基础学习",
        "content": "1.需要认证学习打牢基础; 2.需要进行实践用起来",
        "user_id": 1,
        "created_at": "2025-11-13T03:11:36.581+08:00",
        "updated_at": "2025-11-13T03:11:36.581+08:00"
    }
}
~~~



更新post接口（文章作者才可更新）

删除post接口（文章作者自己才可删）

查询post接口（all 所有文章列表 + 单个文章的详细信息）



// 评论相关

评论创建功能（已认证用户可对文章进行评论）

查询评论功能（读取某文章的所有评论）



// 错误处理和日志记录

response中ERROR 

log打印



