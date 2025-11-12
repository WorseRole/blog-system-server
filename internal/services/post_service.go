package services

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/dto/response"
	"blog-system-server/internal/models"
	"errors"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

// 构造
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		db: db,
	}
}

// 创建文章Post
func (postService *PostService) CreatePost(req *request.CreatePostRequest) (*response.CreatePostResponse, error) {
	// 验参
	if req.Content == "" || req.Title == "" || req.UserID == 0 {
		return nil, errors.New("参数不能为空")
	}

	post := models.Posts{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}
	if err := postService.db.Create(&post).Error; err != nil {
		return nil, errors.New("创建文章失败")
	}

	return &response.CreatePostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}
