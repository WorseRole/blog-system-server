package response

import (
	"time"
)

type CreateCommentResponse struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	UserID    uint64    `json:"user_id"`
	PostID    uint64    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentsResponse struct {
	ID        uint64       `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
