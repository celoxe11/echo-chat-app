package main

import (
	"echo-chat-app-backend/config"
	"echo-chat-app-backend/internal/models"
	"log"
)

// MigrateMySQL runs auto-migration for MySQL tables
func MigrateMySQL() error {
	log.Println("🔄 Running MySQL migrations...")

	// Auto-migrate all models
	err := config.DB.MySQL.AutoMigrate(
		&models.User{},
		&models.Friendship{},
		&models.Group{},
		&models.GroupMember{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ MySQL migrations completed successfully")
	return nil
}

// CreateMongoIndexes creates indexes for MongoDB collections
func CreateMongoIndexes() error {
	log.Println("🔄 Creating MongoDB indexes...")

	// TODO: Add MongoDB indexes here
	// Example:
	// ctx := context.Background()
	// collection := DB.MongoDB.Collection("chat_messages")
	// indexModel := mongo.IndexModel{
	//     Keys: bson.D{{"sender_id", 1}, {"created_at", -1}},
	// }
	// _, err := collection.Indexes().CreateOne(ctx, indexModel)

	log.Println("✅ MongoDB indexes created successfully")
	return nil
}

func main() {
	// Initialize all databases (MySQL, MongoDB, Redis)
	if err := config.InitDatabases(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Run MySQL migrations
	if err := MigrateMySQL(); err != nil {
		log.Fatalf("Failed to run MySQL migrations: %v", err)
	}

	// Create MongoDB indexes
	if err := CreateMongoIndexes(); err != nil {
		log.Fatalf("Failed to create MongoDB indexes: %v", err)
	}

	// Close databases
	config.CloseDatabases()
}