package handler

import (
	"strconv"

	"blog-service/internal/model"
	"blog-service/internal/util"

	"github.com/gin-gonic/gin"
)

func GetCommentList(c *gin.Context) {
	articleId := c.Query("articleId")

	var comments []model.Comment
	model.DB.Where("article_id = ? AND parent_id IS NULL", articleId).Order("created_at DESC").Find(&comments)

	type CommentWithChildren struct {
		model.Comment
		Children []model.Comment `json:"children"`
	}

	result := make([]CommentWithChildren, len(comments))
	for i, comment := range comments {
		result[i] = CommentWithChildren{Comment: comment}
		model.DB.Where("parent_id = ?", comment.ID).Order("created_at ASC").Find(&result[i].Children)
	}

	util.Success(c, result)
}

func AddComment(c *gin.Context) {
	userId := c.GetUint("userId")
	username := c.GetString("username")

	var req struct {
		ArticleID uint   `json:"articleId" binding:"required"`
		Content   string `json:"content" binding:"required"`
		ParentID  *uint  `json:"parentId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	var user model.User
	model.DB.First(&user, userId)

	comment := model.Comment{
		ArticleID: req.ArticleID,
		UserID:    userId,
		Username:  username,
		Avatar:    user.Avatar,
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	if req.ParentID != nil {
		var parent model.Comment
		if model.DB.First(&parent, req.ParentID).Error == nil {
			comment.ReplyTo = parent.Username
		}
	}

	if err := model.DB.Create(&comment).Error; err != nil {
		util.Error(c, 500, "评论失败")
		return
	}

	model.DB.Model(&model.Article{}).Where("id = ?", req.ArticleID).Update("comment_count", model.DB.Raw("comment_count + 1"))

	util.Success(c, comment)
}

func DeleteComment(c *gin.Context) {
	userId := c.GetUint("userId")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var comment model.Comment
	if err := model.DB.First(&comment, id).Error; err != nil {
		util.Error(c, 404, "评论不存在")
		return
	}

	if comment.UserID != userId {
		util.Error(c, 403, "无权限删除")
		return
	}

	articleId := comment.ArticleID
	model.DB.Delete(&comment)
	model.DB.Model(&model.Article{}).Where("id = ?", articleId).Update("comment_count", model.DB.Raw("comment_count - 1"))

	util.Success(c, nil)
}
