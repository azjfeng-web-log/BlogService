package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root:wc8294WCY@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&User{}, &Article{}, &Category{}, &Tag{}, &Comment{}, &ArticleLike{}, &ArticleCollect{})

	initDefaultData()
}

func initDefaultData() {
	var count int64
	DB.Model(&Category{}).Count(&count)
	if count == 0 {
		categories := []Category{
			{Name: "技术"},
			{Name: "生活"},
			{Name: "随笔"},
		}
		DB.Create(&categories)
	}

	DB.Model(&Tag{}).Count(&count)
	if count == 0 {
		tags := []Tag{
			{Name: "Go"},
			{Name: "React"},
			{Name: "Vue"},
			{Name: "JavaScript"},
			{Name: "后端"},
			{Name: "前端"},
		}
		DB.Create(&tags)
	}

	// 将 jameinfeng 设为管理员
	DB.Model(&User{}).Where("username = ?", "jameinfeng").Update("role", "admin")
}
