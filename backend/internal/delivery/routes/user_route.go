package routes

import (
	"echo-chat-app-backend/internal/delivery/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup) {
	userController := controllers.NewUserController()

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/me", userController.Me)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/", userController.SearchUser)
	}
}
