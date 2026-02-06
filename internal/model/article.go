package model

import "time"

type Article struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"size:200"`
	Summary      string    `json:"summary" gorm:"size:500"`
	Content      string    `json:"content" gorm:"type:text"`
	Cover        string    `json:"cover" gorm:"size:255"`
	Category     string    `json:"category" gorm:"size:50"`
	Tags         string    `json:"tags" gorm:"size:255"` // JSON数组存储
	ViewCount    int       `json:"viewCount" gorm:"default:0"`
	LikeCount    int       `json:"likeCount" gorm:"default:0"`
	CommentCount int       `json:"commentCount" gorm:"default:0"`
	AuthorID     uint      `json:"authorId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:50;uniqueIndex"`
}

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:50;uniqueIndex"`
}

type ArticleLike struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	ArticleID uint `gorm:"index"`
}

type ArticleCollect struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index"`
	ArticleID uint `gorm:"index"`
}
