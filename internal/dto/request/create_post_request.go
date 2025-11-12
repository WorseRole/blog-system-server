package request

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  uint64 `json:"-"` // 不从json获取，是从上下文中获取的
}

// 更新Post请求传参
type UpdatePostRequest struct {
	ID      uint64 `json:"id" binding:"required"`
	Title   string `json:"title" `
	Content string `json:"content"`
	UserID  uint64 `json:"-"` // 不从json获取，是从上下文中获取的
}

// 删除Post请求传参
type DeletePostRequest struct {
	ID     uint64 `json:"id" binding:"required"`
	UserID uint64 `json:"-"` // 不从json获取，是从上下文中获取的
}
