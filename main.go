package main

import (
	"blog-service/internal/handler"
	"blog-service/internal/middleware"
	"blog-service/internal/model"

	"github.com/gin-gonic/gin"
)

func main() {
	model.InitDB()

	r := gin.Default()
	r.Use(middleware.Cors())

	// 认证路由
	auth := r.Group("/api_v2/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	// 用户路由
	user := r.Group("/api_v2/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/info", handler.GetUserInfo)
		user.PUT("/info", handler.UpdateUserInfo)
		user.PUT("/password", handler.UpdatePassword)
	}

	// 文章路由
	article := r.Group("/api_v2/article")
	{
		article.GET("/list", handler.GetArticleList)
		article.GET("/recommend", handler.GetRecommendArticles)
		article.GET("/archive", handler.GetArchiveList)
		article.GET("/:id", handler.GetArticleDetail)
		article.GET("/:id/interaction", middleware.JWTAuth(), handler.GetArticleInteraction)
		article.POST("/:id/like", middleware.JWTAuth(), handler.LikeArticle)
		article.POST("/:id/collect", middleware.JWTAuth(), handler.CollectArticle)
		article.POST("", middleware.JWTAuth(), handler.CreateArticle)
		article.PUT("/:id", middleware.JWTAuth(), handler.UpdateArticle)
		article.DELETE("/:id", middleware.JWTAuth(), handler.DeleteArticle)
	}

	// 分类路由
	r.GET("/api_v2/category/list", handler.GetCategoryList)

	// 标签路由
	r.GET("/api_v2/tag/list", handler.GetTagList)

	// 评论路由
	comment := r.Group("/api_v2/comment")
	{
		comment.GET("/list", handler.GetCommentList)
		comment.POST("/add", middleware.JWTAuth(), handler.AddComment)
		comment.DELETE("/:id", middleware.JWTAuth(), handler.DeleteComment)
	}

	r.Run(":6000")
}
