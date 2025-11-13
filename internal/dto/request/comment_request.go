package request

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint64 `json:"postID" binding:"required"`
	UserID  uint64 `json:"-"` // 从上下文中拿 不从json中拿
}

type SelectCommentsListByPostIDRequest struct {
	PostID uint64 `form:"post_id" bingding:"required,min=1"`
}
