package routes

import (
	"echo-chat-app-backend/internal/delivery/controllers"
	"echo-chat-app-backend/internal/delivery/middlewares"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *gin.RouterGroup, authClient *auth.Client, ctrl *controllers.AuthController, mysqlDB *gorm.DB) {
	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.AuthMiddleware(mysqlDB, authClient))
    {	
        authGroup.POST("/sync", ctrl.SyncUser)
    }
}
