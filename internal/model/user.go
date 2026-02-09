package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50"`
	Password  string    `json:"-" gorm:"size:100"`
	Nickname  string    `json:"nickname" gorm:"size:50"`
	Email     string    `json:"email" gorm:"size:100"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Role      string    `json:"role" gorm:"size:20;default:user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
