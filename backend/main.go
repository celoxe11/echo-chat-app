package main

import (
	"echo-chat-app-backend/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize all databases (MySQL, MongoDB, Redis)
	if err := config.InitDatabases(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Setup graceful shutdown
	defer config.CloseDatabases()

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"databases": gin.H{
				"mysql":   "connected",
				"mongodb": "connected",
				"redis":   "connected",
			},
		})
	})

	// TODO: Add your routes here
	// Example:
	// api := router.Group("/api/v1")
	// {
	//     api.POST("/users", controllers.CreateUser)
	//     api.GET("/users/:id", controllers.GetUser)
	//     api.POST("/messages", controllers.SendMessage)
	//     api.GET("/messages", controllers.GetMessages)
	// }

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		config.CloseDatabases()
		os.Exit(0)
	}()

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
