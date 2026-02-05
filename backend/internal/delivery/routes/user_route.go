package routes

import (
	"echo-chat-app-backend/internal/delivery/middlewares"

	"echo-chat-app-backend/internal/delivery/controllers"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.RouterGroup, authClient *auth.Client, ctrl *controllers.UserController, mysqlDB *gorm.DB) {
	userGroup := router.Group("/users")
	userGroup.Use(middlewares.AuthMiddleware(mysqlDB, authClient))
	{
		userGroup.GET("/me", ctrl.Me)
		userGroup.GET("/search", ctrl.SearchUserByUsername)
		userGroup.PATCH("/me", ctrl.UpdateProfile)
	}
}
