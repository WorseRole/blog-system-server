package models

import "time"

type Posts struct {
	ID        uint64    `gorm:"column:id;primaryKey" json:"id"`
	Title     string    `gorm:"column:title;type:varchar(200)" json:"title"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	UserId    uint64    `gorm:"column:user_id; type:int" json:"user_id"`
	IsDel     uint8     `gorm:"column:is_del; type:int" json:"is_del"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}
