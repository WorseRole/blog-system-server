package services

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/dto/response"
	"blog-system-server/internal/models"
	"errors"
	"time"

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
		UserID:    post.UserID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

// 更新文章Post
func (postService *PostService) UpdatePost(req *request.UpdatePostRequest) (*response.CreatePostResponse, error) {
	// 校验参数
	if req.Content == "" && req.Title == "" {
		return nil, errors.New("修改文章不能传参为空")
	}
	var oldPost models.Posts
	// 校验只能自己更新自己的
	if err := postService.db.Where("is_del= 0 and id = ?", req.ID).First(&oldPost).Error; err != nil {
		return nil, errors.New("没找到对应的文章")
	}
	if oldPost.UserID != req.UserID {
		return nil, errors.New("不是文章的作者不能修改")
	}
	if req.Content == oldPost.Content && req.Title == oldPost.Title {
		return nil, errors.New("文章没有修改项")
	}

	// 构建更新数据
	updateData := map[string]interface{}{}

	if req.Content != oldPost.Content {
		updateData["content"] = req.Content
	}
	if req.Title != oldPost.Title {
		updateData["title"] = req.Title
	}
	// 设置更新时间
	updateData["updated_at"] = time.Now()

	result := postService.db.Model(&models.Posts{}).Where("is_del= 0 and id = ?", req.ID).Updates(updateData)

	if result.Error != nil {
		return nil, errors.New("更新文章失败")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("文章更新失败")
	}

	// 返回更新后的文章信息
	var updatePost models.Posts
	if err := postService.db.Where("is_del= 0 and id = ?", req.ID).First(&updatePost).Error; err != nil {
		return nil, errors.New("获取更新后的文章失败")
	}

	return &response.CreatePostResponse{
		ID:        updatePost.ID,
		Title:     updatePost.Title,
		Content:   updatePost.Content,
		UserID:    updatePost.UserID,
		CreatedAt: updatePost.CreatedAt,
		UpdatedAt: updatePost.UpdatedAt,
	}, nil
}

// 删除文章
func (postService *PostService) DeletePost(req *request.DeletePostRequest) (*response.CreatePostResponse, error) {
	// 校验参数
	if req.ID == 0 {
		return nil, errors.New("参数不对")
	}
	// 查询
	var oldPost models.Posts
	if err := postService.db.Where("is_del= 0 and id = ?", req.ID).First(&oldPost).Error; err != nil {
		return nil, errors.New("没找到对应的文章")
	}
	// 校验只能自己更新自己的
	if oldPost.UserID != req.UserID {
		return nil, errors.New("不是文章的作者不能删除")
	}

	// 更新is_del=1
	result := postService.db.Model(&models.Posts{}).Where("is_del= 0 and id = ?", req.ID).Update("is_del", 1)
	if result.Error != nil {
		return nil, errors.New("删除文章失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("文章删除失败")
	}
	// 返回已删除的数据
	var deletePost models.Posts
	if err := postService.db.Where("is_del= 1 and id = ?", req.ID).First(&deletePost).Error; err != nil {
		return nil, errors.New("返回删除后的文章失败")
	}

	return &response.CreatePostResponse{
		ID:        deletePost.ID,
		Title:     deletePost.Title,
		Content:   deletePost.Content,
		UserID:    deletePost.UserID,
		CreatedAt: deletePost.CreatedAt,
		UpdatedAt: deletePost.UpdatedAt,
	}, nil
}
