package routes

import (
	"echo-chat-app-backend/internal/delivery/controllers"
	"echo-chat-app-backend/internal/repositories"
	"echo-chat-app-backend/internal/usecases"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(mysqlDB *gorm.DB, firebaseAuth *auth.Client) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// declare repositories, usecases, controllers here
	// Dependency Injection
	authRepo := repositories.NewAuthRepository(mysqlDB)
    authUseCase := usecases.NewAuthUseCase(authRepo)
	authController := controllers.NewAuthController(authUseCase)

	api := router.Group("/api")
	{
		SetupAuthRoutes(api, firebaseAuth, authController)
		SetupUserRoutes(api)
	}

	return router
}