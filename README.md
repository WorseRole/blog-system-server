## 需求文档

使用 Go 语言结合 Gin 框架和 GORM 库开发一个个人博客系统的后端，实现博客文章的基本管理功能，包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能。 

#### 1. 项目初始化

创建一个新的 Go 项目，使用 go mod init 初始化项目依赖管理。
安装必要的库，如 Gin 框架、GORM 以及数据库驱动（如 MySQL 或 SQLite）。

#### 2.数据库设计与模型定义

设计数据库表结构，至少包含以下几个表：
users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
使用 GORM 定义对应的 Go 模型结构体。

#### 3.用户认证与授权

实现用户注册和登录功能，用户注册时需要对密码进行加密存储，登录时验证用户输入的用户名和密码。
使用 JWT（JSON Web Token）实现用户认证和授权，用户登录成功后返回一个 JWT，后续的需要认证的接口需要验证该 JWT 的有效性。

#### 4.文章管理功能

实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
实现文章的更新功能，只有文章的作者才能更新自己的文章。
实现文章的删除功能，只有文章的作者才能删除自己的文章。

#### 5.评论功能

实现评论的创建功能，已认证的用户可以对文章发表评论。
实现评论的读取功能，支持获取某篇文章的所有评论列表。

#### 6.错误处理与日志记录

对可能出现的错误进行统一处理，如数据库连接错误、用户认证失败、文章或评论不存在等，返回合适的 HTTP 状态码和错误信息。
使用日志库记录系统的运行信息和错误信息，方便后续的调试和维护。



## 先从搭建一个Go服务开始

#### 1. 创建并初始化项目：

~~~bash
mkdir blog-system-server
cd blog-system-server
go mod init blog-system-server
~~~

#### 2. 安装依赖： 

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
// 最后再添加Zap 日志库
go get -u go.uber.org/zap
~~~

#### 3. 编写测试代码

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



### 配置文件

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



### 接口文档 (同样为Postman 测试用例)

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

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

传参body：

~~~json
{
    "title":"GoLang基础学习",
    "content":"1.需要认证学习打牢基础; 2.需要进行实践用起来"
}
~~~

返回结果：

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



##### 2.2 更新post接口	✅	已完成

> 接口说明: 更新文章 标题或内容（文章作者才可更新）

> 是否需要Token：是

URL:

~~~http
POST  http://localhost:8080/api/post/update
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

传参body：

~~~json
{
    "id":1,
    "title":"GoLang基础学习01",
    "content":"1.需要认证学习打牢基础; 2.需要进行实践用起来; 3.做项目运用； 4.看八股面经以及刷算法"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "更新文章成功",
    "data": {
        "id": 1,
        "title": "GoLang基础学习01",
        "content": "1.需要认证学习打牢基础; 2.需要进行实践用起来; 3.做项目运用； 4.看八股面经以及刷算法",
        "user_id": 1,
        "created_at": "2025-11-13T03:11:37+08:00",
        "updated_at": "2025-11-14T02:52:39+08:00"
    }
}
~~~



##### 2.3 删除post接口	✅	已完成

> 接口说明: 删除文章 标题或内容（文章作者自己才可删）

> 是否需要Token：是

URL:

~~~http
POST  http://localhost:8080/api/post/delete
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

传参body：

~~~json
{
    "id":1
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "删除文章成功",
    "data": {
        "id": 1,
        "title": "GoLang基础学习01",
        "content": "1.需要认证学习打牢基础; 2.需要进行实践用起来; 3.做项目运用； 4.看八股面经以及刷算法",
        "user_id": 1,
        "created_at": "2025-11-13T03:11:37+08:00",
        "updated_at": "2025-11-14T02:56:34+08:00"
    }
}
~~~



##### 2.4 查询post列表接口	✅	已完成

> 接口说明: 删除文章 标题或内容（文章作者自己才可删）

> 是否需要Token：是

URL:

~~~http
GET  http://localhost:8080/api/post/lists
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "成功获取文章列表",
    "data": [
        {
            "id": 4,
            "user_id": 1,
            "title": "GoLang高级学习",
            "content": "1.需要高级学习; 2.DApp后端开发学习",
            "created_at": "2025-11-14T03:06:07+08:00",
            "updated_at": "2025-11-14T03:06:07+08:00"
        },
        {
            "id": 3,
            "user_id": 1,
            "title": "GoLang进阶学习",
            "content": "1.需要进阶学习; 2.需要做清楚项目",
            "created_at": "2025-11-14T03:05:35+08:00",
            "updated_at": "2025-11-14T03:05:35+08:00"
        },
        {
            "id": 2,
            "user_id": 1,
            "title": "GoLang基础学习",
            "content": "1.需要认证学习打牢基础; 2.需要进行实践用起来",
            "created_at": "2025-11-14T03:03:05+08:00",
            "updated_at": "2025-11-14T03:03:05+08:00"
        }
    ]
}
~~~



##### 2.5 查询post单个文章详情接口	✅	已完成

> 接口详情：已认证登录的用户可查询单个文章的详情（可以查看文章对应的作者以及评论以及评论的作者信息）

> 是否需要Token：是

URL:

~~~http
GET  http://localhost:8080/api/post/detail?post_id=2
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "成功获取文章详细信息",
    "data": {
        "id": 2,
        "user_id": 0,
        "title": "GoLang基础学习",
        "content": "1.需要认证学习打牢基础; 2.需要进行实践用起来",
        "comments": [
            {
                "id": 1,
                "content": "一本非常优秀的书，基础教学扎实。赞！",
                "user": {
                    "id": 1,
                    "username": "leoYan"
                },
                "created_at": "2025-11-14T03:19:05+08:00",
                "updated_at": "2025-11-14T03:19:05+08:00"
            }
        ],
        "created_at": "2025-11-14T03:03:05+08:00",
        "updated_at": "2025-11-14T03:03:05+08:00"
    }
}
~~~



#### 3. 评论相关

##### 3.1评论创建功能	✅	已完成

> 接口说明：已认证用户可对文章进行评论 创建评论

> 是否需要Token：是

URL:

~~~http
POST  http://localhost:8080/api/comment/create
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

传参body:

~~~json
{
    "post_id":2,
    "content":"一本非常优秀的书，基础教学扎实。赞！"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "创建评论成功",
    "data": {
        "id": 1,
        "content": "一本非常优秀的书，基础教学扎实。赞！",
        "user_id": 1,
        "post_id": 2,
        "created_at": "2025-11-14T03:19:04.568+08:00",
        "updated_at": "2025-11-14T03:19:04.568+08:00"
    }
}
~~~



##### 2.5 查询评论功能	✅	已完成

> 接口说明：已认证用户可对查看某文章的所有评论列表

> 是否需要Token：是

URL:

~~~http
GET  http://localhost:8080/api/comment/lists?post_id=2
~~~

传参header：

~~~json
{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6Imxlb1lhbiIsImlzcyI6ImJsb2ctc3lzdGVtIiwic3ViIjoidXNlci10b2tlbiIsImV4cCI6MTc2MzE0NjE5NiwibmJmIjoxNzYzMDU5Nzk2LCJpYXQiOjE3NjMwNTk3OTZ9.-BpBhTSb7SgR-syS3p77bXeAjKFUw69EUvEyxlo76Tc"
}
~~~

返回结果：

~~~json
{
    "code": 200,
    "message": "通过post_id读取评论列表成功",
    "data": [
        {
            "id": 1,
            "content": "一本非常优秀的书，基础教学扎实。赞！",
            "user_id": 1,
            "post_id": 2,
            "created_at": "2025-11-14T03:19:05+08:00",
            "updated_at": "2025-11-14T03:19:05+08:00"
        }
    ]
}
~~~





## 错误处理和日志记录

#### response

> 用pkg.response.response.go 来进行输出 状态码和 信息

#### log打印

> 用pkg.logger.logger.go 以及 access_log.log 中间件进行日志输出打印
>
> 增加.env的日志配置



