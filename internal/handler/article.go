package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"blog-service/internal/model"
	"blog-service/internal/util"

	"github.com/gin-gonic/gin"
)

type ArticleResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Cover        string    `json:"cover"`
	Category     string    `json:"category"`
	Tags         []string  `json:"tags"`
	AuthorID     uint      `json:"authorId"`
	ViewCount    int       `json:"viewCount"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func toArticleResponse(a model.Article) ArticleResponse {
	var tags []string
	json.Unmarshal([]byte(a.Tags), &tags)
	return ArticleResponse{
		ID:           a.ID,
		Title:        a.Title,
		Summary:      a.Summary,
		Content:      a.Content,
		Cover:        a.Cover,
		Category:     a.Category,
		Tags:         tags,
		AuthorID:     a.AuthorID,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		CommentCount: a.CommentCount,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	category := c.Query("category")
	tag := c.Query("tag")
	keyword := c.Query("keyword")

	query := model.DB.Model(&model.Article{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var articles []model.Article
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)

	list := make([]ArticleResponse, len(articles))
	for i, a := range articles {
		list[i] = toArticleResponse(a)
	}

	util.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func GetArticleDetail(c *gin.Context) {
	id := c.Param("id")

	var article model.Article
	if err := model.DB.First(&article, id).Error; err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	model.DB.Model(&article).Update("view_count", article.ViewCount+1)

	util.Success(c, toArticleResponse(article))
}

func GetRecommendArticles(c *gin.Context) {
	var articles []model.Article
	model.DB.Order("view_count DESC").Limit(5).Find(&articles)

	list := make([]ArticleResponse, len(articles))
	for i, a := range articles {
		list[i] = toArticleResponse(a)
	}

	util.Success(c, list)
}

func GetArchiveList(c *gin.Context) {
	var articles []model.Article
	model.DB.Order("created_at DESC").Find(&articles)

	archives := make(map[string][]ArticleResponse)
	for _, a := range articles {
		yearMonth := a.CreatedAt.Format("2006-01")
		archives[yearMonth] = append(archives[yearMonth], toArticleResponse(a))
	}

	util.Success(c, archives)
}

func CreateArticle(c *gin.Context) {
	userId := c.GetUint("userId")

	var req struct {
		Title    string   `json:"title" binding:"required"`
		Summary  string   `json:"summary"`
		Content  string   `json:"content" binding:"required"`
		Cover    string   `json:"cover"`
		Category string   `json:"category" binding:"required"`
		Tags     []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	article := model.Article{
		Title:    req.Title,
		Summary:  req.Summary,
		Content:  req.Content,
		Cover:    req.Cover,
		Category: req.Category,
		Tags:     string(tagsJSON),
		AuthorID: userId,
	}

	if err := model.DB.Create(&article).Error; err != nil {
		util.Error(c, 500, "创建失败")
		return
	}

	util.Success(c, toArticleResponse(article))
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var article model.Article
	if err := model.DB.First(&article, id).Error; err != nil {
		util.Error(c, 404, "文章不存在")
		return
	}

	var req struct {
		Title    string   `json:"title"`
		Summary  string   `json:"summary"`
		Content  string   `json:"content"`
		Cover    string   `json:"cover"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Cover != "" {
		updates["cover"] = req.Cover
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		updates["tags"] = string(tagsJSON)
	}

	model.DB.Model(&article).Updates(updates)

	model.DB.First(&article, id)
	util.Success(c, toArticleResponse(article))
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	if err := model.DB.Delete(&model.Article{}, id).Error; err != nil {
		util.Error(c, 500, "删除失败")
		return
	}

	util.Success(c, nil)
}

func LikeArticle(c *gin.Context) {
	userId := c.GetUint("userId")
	articleId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var like model.ArticleLike
	result := model.DB.Where("user_id = ? AND article_id = ?", userId, articleId).First(&like)

	if result.Error != nil {
		model.DB.Create(&model.ArticleLike{UserID: userId, ArticleID: uint(articleId)})
		model.DB.Model(&model.Article{}).Where("id = ?", articleId).Update("like_count", model.DB.Raw("like_count + 1"))
		util.Success(c, gin.H{"liked": true})
	} else {
		model.DB.Delete(&like)
		model.DB.Model(&model.Article{}).Where("id = ?", articleId).Update("like_count", model.DB.Raw("like_count - 1"))
		util.Success(c, gin.H{"liked": false})
	}
}

func CollectArticle(c *gin.Context) {
	userId := c.GetUint("userId")
	articleId, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var collect model.ArticleCollect
	result := model.DB.Where("user_id = ? AND article_id = ?", userId, articleId).First(&collect)

	if result.Error != nil {
		model.DB.Create(&model.ArticleCollect{UserID: userId, ArticleID: uint(articleId)})
		util.Success(c, gin.H{"collected": true})
	} else {
		model.DB.Delete(&collect)
		util.Success(c, gin.H{"collected": false})
	}
}

func GetCategoryList(c *gin.Context) {
	var categories []model.Category
	model.DB.Find(&categories)
	util.Success(c, categories)
}

func GetTagList(c *gin.Context) {
	var tags []model.Tag
	model.DB.Find(&tags)
	util.Success(c, tags)
}
