package routes

import (
	"echo-chat-app-backend/internal/delivery/middlewares"
	"echo-chat-app-backend/internal/delivery/controllers"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, authClient *auth.Client, ctrl *controllers.AuthController) {
	authGroup := router.Group("/auth")
	authGroup.Use(middlewares.AuthMiddleware(authClient))
    {
        authGroup.POST("/sync", ctrl.SyncUser)
    }
}
