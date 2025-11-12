package services

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/dto/response"
	"blog-system-server/internal/models"
	"errors"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

// 构造
func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{
		db: db,
	}
}

// 创建评论
func (commentService *CommentService) CreateComment(req *request.CreateCommentRequest) (*response.CreateCommentResponse, error) {
	// 校验参数
	if req.Content == "" || req.PostID == 0 || req.UserID == 0 {
		return nil, errors.New("创建评论参数错误")
	}
	var post models.Posts
	// 查看文章是否存在
	if err := commentService.db.Where("is_del = 0 and post_id = ?").First(&post).Error; err != nil {
		return nil, errors.New("评论的该文章不存在")
	}

	comment := models.Comments{
		PostID:  req.PostID,
		UserID:  req.UserID,
		Content: req.Content,
	}
	// 创建评论
	if err := commentService.db.Create(&comment).Error; err != nil {
		return nil, errors.New("创建评论失败")
	}
	// 返回评论
	return &response.CreateCommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}, nil
}
