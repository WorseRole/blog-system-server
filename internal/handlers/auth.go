package handlers

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/services"
	"blog-system-server/pkg/response"

	"github.com/gin-gonic/gin"
)

/*
这么写的好处: AuthHandler
确保依赖注入（Dependency Injection）规范；
避免全局变量；
保持 handler 逻辑可测试。
*/
type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// 登录接口
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginUserRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误,参数不对")
		return
	}
	// 调用业务服务
	result, err := h.authService.Login(req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	// 返回成功响应
	response.Success(c, "登录成功", result)
}

// 注册接口
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误,参数不对")
		return
	}

	// 调用业务服务
	result, err := h.authService.Register(req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 返回成功响应
	response.Success(c, "注册成功,请登录", result)
}
