package handler

import (
	"blog-service/internal/model"
	"blog-service/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	var user model.User
	if err := model.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		util.Error(c, 400, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		util.Error(c, 400, "用户名或密码错误")
		return
	}

	token, _ := util.GenerateToken(user.ID, user.Username)

	util.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	var count int64
	model.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		util.Error(c, 400, "用户名已存在")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Username,
	}

	if err := model.DB.Create(&user).Error; err != nil {
		util.Error(c, 500, "注册失败")
		return
	}

	token, _ := util.GenerateToken(user.ID, user.Username)

	util.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func GetUserInfo(c *gin.Context) {
	userId := c.GetUint("userId")

	var user model.User
	if err := model.DB.First(&user, userId).Error; err != nil {
		util.Error(c, 404, "用户不存在")
		return
	}

	util.Success(c, user)
}

func UpdateUserInfo(c *gin.Context) {
	userId := c.GetUint("userId")

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	model.DB.Model(&model.User{}).Where("id = ?", userId).Updates(updates)

	var user model.User
	model.DB.First(&user, userId)

	util.Success(c, user)
}

func UpdatePassword(c *gin.Context) {
	userId := c.GetUint("userId")

	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, 400, "参数错误")
		return
	}

	var user model.User
	model.DB.First(&user, userId)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		util.Error(c, 400, "原密码错误")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	model.DB.Model(&user).Update("password", string(hashedPassword))

	util.Success(c, nil)
}
