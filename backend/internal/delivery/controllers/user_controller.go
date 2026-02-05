package controllers

import (
	"echo-chat-app-backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase *usecases.UserUseCase
}

func NewUserController(userUseCase *usecases.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (uc *UserController) Me(c *gin.Context) {
	uid := c.GetString("firebase_uid") //ambil uid dari middleware
	user, err := uc.userUseCase.Me(uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Fetch profile successfully", "user": user})
}

func (uc *UserController) SearchUserByUsername(c *gin.Context) {
	username := c.Query("username")
	user, err := uc.userUseCase.SearchUserByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to search user: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User searched successfully", "user": user})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	uid := c.GetString("firebase_uid")

	// ambil data dari body
	name := c.PostForm("name")
	username := c.PostForm("username")
	avatar_url := c.PostForm("avatar_url")

	user, err := uc.userUseCase.UpdateProfile(uid, name, username, avatar_url)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User updated successfully", "user": user})
}	
