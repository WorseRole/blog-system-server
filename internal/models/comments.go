package models

import "time"

type Comments struct {
	ID        uint64    `gorm:"column:id;primaryKey" json:"id"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	UserID    uint64    `gorm:"column:user_id;type:int" json:"user_id"`
	PostID    uint64    `gorm:"column:post_id;type:int" json:"post_id"`
	IsDel     int       `gorm:"column:is_del;type:int;default:0" json:"is_del"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}
