##### 1. 创建并初始化项目：

~~~bash
mkdir blog-system-server
cd blog-system-server
go mod init blog-system-server
~~~

##### 2. 安装依赖： 

使用go get 添加 Gin 框架

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

