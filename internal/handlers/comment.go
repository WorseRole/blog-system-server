package handlers

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/middleware"
	"blog-system-server/internal/services"
	"blog-system-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// 创建评论接口
func (commentHandler *CommentHandler) CreateComment(c *gin.Context) {
	// 1. 先获取userID
	UserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误:用户丢失")
		return
	}
	var req request.CreateCommentRequest
	// 2. 获取参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	req.UserID = UserID
	// 3. 调业务
	result, err := commentHandler.commentService.CreateComment(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "创建评论成功", result)
}

// 通过 post_id 读取评论列表接口
func (commentHandler *CommentHandler) SelectCommentsListByPostID(c *gin.Context) {
	// 1. 先获取UserID
	_, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Error(c, "系统错误:用户丢失")
		return
	}
	var req request.SelectCommentsListByPostIDRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
	}
	// 业务调用
	result, err := commentHandler.commentService.SelectCommentsListByPostID(req.PostID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "通过post_id读取评论列表成功", result)
}
