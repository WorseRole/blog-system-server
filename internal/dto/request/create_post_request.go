package request

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  uint64 `json:"-"` // 不从json获取，是从上下文中获取的
}
