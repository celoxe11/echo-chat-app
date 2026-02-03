package main

import (
	"echo-chat-app-backend/config"
	"echo-chat-app-backend/internal/delivery/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

	// Initialize all databases (MySQL, MongoDB, Redis)
	if err := config.InitDatabases(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer config.CloseDatabases()

	// Initialize Firebase
	if err := config.InitFirebase(); err != nil {
		log.Fatal(err)
	}

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
	r := routes.SetupRouter(config.DB.MySQL, config.FirebaseAuth)

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
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
