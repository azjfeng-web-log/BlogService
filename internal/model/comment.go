package model

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"articleId" gorm:"index"`
	UserID    uint      `json:"userId"`
	Username  string    `json:"username" gorm:"size:50"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	Content   string    `json:"content" gorm:"size:1000"`
	ParentID  *uint     `json:"parentId"`
	ReplyTo   string    `json:"replyTo" gorm:"size:50"`
	CreatedAt time.Time `json:"createdAt"`
}
