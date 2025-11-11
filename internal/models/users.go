package models

import "time"

type Users struct {
	ID        uint64    `gorm:"primaryKey;column:id" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(100)" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(200)" json:"password"`
	Salt      string    `gorm:"column:salt;type:varchar(20)" json:"salt"`
	Email     string    `gorm:"column:email; type:varchar(200)"  json:"email"`
	IsDel     int       `gorm:"column:is_del;type:int;default:0" json:"is_del"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}
