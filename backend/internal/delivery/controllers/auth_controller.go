package controllers

import (
	"echo-chat-app-backend/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthController(authUseCase *usecases.AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

func (ac *AuthController) SyncUser(c *gin.Context) {
	uid := c.GetString("firebase_uid")
	email := c.GetString("email")
	name := c.GetString("name")
	avatarURL := c.GetString("avatar_url")

	if uid == "" || email == "" {
		c.JSON(400, gin.H{"error": "Invalid user data"})
		return
	}

	user, err := ac.authUseCase.SyncUser(uid, email, name, avatarURL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to sync user: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User synced successfully", "user": user})
}
