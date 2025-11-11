package services

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/dto/response"
	"blog-system-server/internal/models"
	"blog-system-server/pkg/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

// 构造
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

// 注册 让authService调
func (authService *AuthService) Register(request request.RegisterRequest) (*response.RegisterResponse, error) {
	// 1. 先检查 用户名是否存在
	var existingUser models.Users
	if err := authService.db.Where("username = ? or email = ?", request.Username, request.Email).First(&existingUser).Error; err == nil {
		log.Printf("用户名或邮箱已存在: %v", err)
		return nil, errors.New("用户名或邮箱已存在")
	}

	// 2. 检查邮箱是否存在

	// 3.1 生成盐值
	salt, err := utils.GenerateSalt()
	if err != nil {
		return nil, errors.New("生成盐值失败")
	}

	// 3.2 密码 + 盐值 进行加密
	hashedPassword, err := utils.HashPasswordWithSalt(request.Password, salt)
	if err != nil {
		return nil, errors.New("盐值加密失败")
	}

	// 4.1 创建用户 赋值
	user := models.Users{
		Username: request.Username,
		Email:    request.Email,
		Salt:     salt,
		Password: hashedPassword,
	}
	// 4.2 sql写入库
	if err := authService.db.Create(&user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 5. 返回用户信息（不包含敏感信息）
	return &response.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
