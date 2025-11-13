package response

import (
	"time"
)

type CreatePostResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    uint64    `json:"userID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SelectPostResponse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SelectPostDetialResponse struct {
	ID               uint64             `json:"id"`
	UserID           uint64             `json:"user_id"`
	Title            string             `json:"title"`
	Content          string             `json:"content"`
	CommentsResponse []CommentsResponse `json:"comments"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
