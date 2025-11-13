package handlers

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/middleware"
	"blog-system-server/internal/services"
	"blog-system-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

// 构造
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// 创建文章接口
func (P *PostHandler) CreatePost(c *gin.Context) {

	// 1. 先获取userID
	UserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误:用户丢失")
		return
	}
	// 2. 请求参数
	req := request.CreatePostRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数不对,参数错误")
		return
	}
	req.UserID = UserID

	// 3. 调用业务
	result, err := P.postService.CreatePost(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "创建文章成功", result)
}

// 更新文章接口
func (P *PostHandler) UpdatePost(c *gin.Context) {
	// 1. 先获取userID
	UserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误:用户丢失")
		return
	}
	// 2. 请求参数
	req := request.UpdatePostRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数不对,参数错误")
		return
	}
	req.UserID = UserID
	// 3. 调用业务
	result, err := P.postService.UpdatePost(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "更新文章成功", result)
}

// 删除文章接口
func (P *PostHandler) DeletePost(c *gin.Context) {
	// 1. 先获取userID
	UserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误:用户丢失")
		return
	}
	// 2. 请求参数
	req := request.DeletePostRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数不对,参数错误")
		return
	}
	req.UserID = UserID
	// 3. 调用业务
	result, err := P.postService.DeletePost(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "删除文章成功", result)
}

// 查询文章列表接口
func (P *PostHandler) SelectPostsLists(c *gin.Context) {
	// 1. 先获取UserID
	_, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误，用户丢失")
		return
	}
	// 2. 业务调用 此接口不用参数
	result, err := P.postService.SelectPostsLists()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "成功获取文章列表", result)
}

// 查询单个文章详情接口（包含文章+文章评论（+每个评论的用户信息））
func (P *PostHandler) SelectPostDetial(c *gin.Context) {
	// 1. 先获取UserID
	_, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误，用户丢失")
		return
	}
	// 构建传参
	req := request.SelectPostRequestOne{}
	// 自动绑定query参数并转换类型
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 业务调用
	result, err := P.postService.SelectPostDetial(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "成功获取文章详细信息", result)
}
