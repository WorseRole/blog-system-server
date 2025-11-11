package models

import "time"

type Users struct {
	ID        uint64    `gorm:"primaryKey;column:id" json:"id"`
	UserName  string    `gorm:"column:user_name; type:varchar(100)" json:"user_name"`
	Email     string    `gorm:"column:email; type:varchar(200)"  json:"email"`
	IsDel     string    `gorm:"column:is_del; type:int" json:"is_del"`
	CreatedAt time.Time `gorm:"column:created_at; type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:datetime" json:"updated_at"`
}
