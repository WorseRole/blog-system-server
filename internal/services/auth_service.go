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

// 登录 让authService调
func (authService *AuthService) Login(req request.LoginUserRequest) (*response.LoginUserResponse, error) {
	// 1. 先检查用户是否存在
	var existingUser = models.Users{}

	if err := authService.db.Where("username = ?", req.Username).First(&existingUser).Error; err != nil {
		// 用户不存在
		log.Printf("用户不存在:%s, err:%v", req.Username, err)
		return nil, errors.New("用户不存在")
	}
	// 2. 通过查询拿到密码+盐 进行Hash 校验密码
	bol := utils.CheckPasswordWithSalt(req.Password, existingUser.Password, existingUser.Salt)
	if !bol {
		// 校验未通过
		log.Printf("密码不对.username:%s, password:%s", req.Username, req.Password)
		return nil, errors.New("密码不匹配")
	}
	// 3. 生成JWT Token
	token, err := utils.GenerateJWT(existingUser.ID, existingUser.Username)
	if err != nil {
		log.Printf("生成Token失败.username:%s", req.Username)
		return nil, errors.New("生成Token失败")
	}
	// 返回结果
	return &response.LoginUserResponse{
		Token:    token,
		ID:       existingUser.ID,
		Username: existingUser.Username,
		Email:    existingUser.Email,
	}, nil
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
