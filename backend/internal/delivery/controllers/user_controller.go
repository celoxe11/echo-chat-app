package controllers

import (
	"echo-chat-app-backend/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase *usecases.UserUseCase
}

func NewUserController() *UserController {
	userUseCase := &usecases.UserUseCase{}
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (uc *UserController) Me(c *gin.Context) {
	// TODO: Implementation for handling the "Me" endpoint
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	// TODO: Implementation for handling the "GetUserByID" endpoint
}

func (uc *UserController) SearchUser(c *gin.Context) {
	// TODO: Implementation for handling the "SearchUser" endpoint
}