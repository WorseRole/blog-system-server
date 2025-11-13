package services

import (
	"blog-system-server/internal/dto/request"
	"blog-system-server/internal/dto/response"
	"blog-system-server/internal/models"
	"blog-system-server/pkg/logger"
	"errors"
	"time"

	"go.uber.org/zap"
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
	logger.Logger.Info(
		"创建文章",
		zap.Uint64("user_id", req.UserID),
		zap.String("title", req.Title),
	)
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

	logger.Logger.Info("创建文章成功",
		zap.Uint64("post_id", post.ID),
		zap.Uint64("user_id", req.UserID),
	)

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

// 删除文章 是否需要同时也删除文章对应的评论呢？
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

// 读取所有文章列表
func (postService *PostService) SelectPostsLists() (*[]response.SelectPostResponse, error) {

	// 查询所有的post 构造 posts[]
	var posts []struct {
		ID        uint64    `gorm:"column:id"`
		Title     string    `gorm:"column:title"`
		Content   string    `gorm:"column:content"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		UserID    uint64    `gorm:"column:user_id"`
	}
	postsSql := "select * from posts where is_del = 0 order by created_at desc"
	if err := postService.db.Raw(postsSql).Scan(&posts).Error; err != nil {
		return nil, errors.New("查询文章失败")
	}
	// 构建返回
	var result []response.SelectPostResponse
	for _, post := range posts {
		result = append(result, response.SelectPostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}
	return &result, nil
}

// 读取单个文章所有信息
func (postService *PostService) SelectPostDetial(req *request.SelectPostRequestOne) (*response.SelectPostDetialResponse, error) {
	// 获取单个文章
	var posts models.Posts
	if err := postService.db.Find(&posts, models.Posts{ID: req.ID, IsDel: 0}).Error; err != nil {
		return nil, errors.New("查询单个文章失败")
	}
	// 根据文章ID 查所有的Comment
	var comments []models.Comments
	if err := postService.db.Find(&comments, models.Comments{PostID: req.ID, IsDel: 0}).Error; err != nil {
		return nil, errors.New("查询单个文章对应的所有评论失败")
	}

	// 根据评论中的userID 查询 所有评论的用户信息
	var userIDList []uint64 = make([]uint64, 0)
	for _, comment := range comments {
		userIDList = append(userIDList, comment.UserID)
	}

	var userList []models.Users = make([]models.Users, 0)
	if err := postService.db.Find(&userList, "is_del = 0 and id in ?", userIDList).Error; err != nil {
		return nil, errors.New("查询单个文章对应的所有评论的用户失败")
	}

	// userID : users
	var usersMap map[uint64]response.UserResponse = map[uint64]response.UserResponse{}
	for i := 0; i < len(userList); i++ {
		var user = userList[i]
		usersMap[user.ID] = response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
		}
	}

	// comment 构造返回参数[]Comments
	var commentsResponse []response.CommentsResponse = make([]response.CommentsResponse, 0)
	for _, comment := range comments {
		commentsResponse = append(commentsResponse, response.CommentsResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			User:      usersMap[comment.UserID],
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	return &response.SelectPostDetialResponse{
		ID:               posts.ID,
		Title:            posts.Title,
		Content:          posts.Content,
		CommentsResponse: commentsResponse,
		CreatedAt:        posts.CreatedAt,
		UpdatedAt:        posts.UpdatedAt,
	}, nil

}
