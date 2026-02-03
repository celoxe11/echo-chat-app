package main

import (
	"context"
	"echo-chat-app-backend/config"
	"echo-chat-app-backend/internal/models"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/iterator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dropTables   = flag.Bool("drop", false, "Drop all tables before seeding")
	fresh        = flag.Bool("fresh", false, "Drop and recreate database (fresh start)")
	seedOnly     = flag.Bool("seed-only", false, "Only seed data without migrations")
	dummyUsers   = flag.Bool("dummy-users", false, "Create dummy users instead of fetching from Firebase")
	skipFirebase = flag.Bool("skip-firebase", false, "Skip Firebase initialization")
)

func main() {
	flag.Parse()

	// Load environment variables from project root
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	log.Println("🌱 Starting database seeder...")

	// 1. Initialize Firebase
	if !*skipFirebase && !*dummyUsers {
		if err := config.InitFirebase(); err != nil {
			log.Fatalf("❌ Error: Firebase initialization failed: %v. Use --dummy-users or --skip-firebase if you don't want to use Firebase.", err)
		}
	}

	// 2. Handle Fresh Start (Drop and Recreate DB)
	if *fresh {
		if err := recreateDatabase(); err != nil {
			log.Fatalf("❌ Failed to recreate database: %v", err)
		}
	}

	// 3. Initialize Databases
	if err := config.InitDatabases(); err != nil {
		log.Fatalf("❌ Failed to initialize databases: %v", err)
	}
	defer config.CloseDatabases()

	// 4. Handle table drops if fresh was not used but drop was requested
	if !*fresh && *dropTables {
		if err := dropAllTables(); err != nil {
			log.Fatalf("❌ Failed to drop tables: %v", err)
		}
	}

	// 5. Run migrations
	if !*seedOnly {
		if err := runMigrations(); err != nil {
			log.Fatalf("❌ Failed to run migrations: %v", err)
		}
	}

	// 6. Seed data
	if err := seedData(); err != nil {
		log.Fatalf("❌ Seeding failed: %v", err)
	}

	log.Println("✅ Database seeding process finished!")
}

func recreateDatabase() error {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DATABASE")

	// 1. Recreate MySQL Database
	// We connect to MySQL server without specifying a database in the DSN to allow dropping/creating it
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, host)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server for recreation: %w", err)
	}

	log.Printf("🔥 Resetting MySQL database: %s", dbName)
	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", dbName)).Error; err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", dbName)).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 2. Drop MongoDB Database (if configured)
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")
	if mongoURI != "" && mongoDBName != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err == nil {
			log.Printf("🔥 Dropping MongoDB database: %s", mongoDBName)
			_ = client.Database(mongoDBName).Drop(ctx)
			_ = client.Disconnect(ctx)
		}
	}

	return nil
}

func dropAllTables() error {
	log.Println("🗑️  Dropping all tables and collections...")
	err := config.DB.MySQL.Migrator().DropTable(
		&models.GroupMember{},
		&models.Group{},
		&models.Friendship{},
		&models.User{},
	)

	if err != nil {
		return err
	}

	ctx := context.Background()
	config.DB.MongoDB.Collection("chat_messages").Drop(ctx)
	config.DB.MongoDB.Collection("conversations").Drop(ctx)

	log.Println("✅ Tables and collections cleared")
	return nil
}

func runMigrations() error {
	log.Println("🔄 Running migrations...")
	return config.DB.MySQL.AutoMigrate(
		&models.User{},
		&models.Friendship{},
		&models.Group{},
		&models.GroupMember{},
	)
}

func seedData() error {
	log.Println("🌱 Seeding data...")

	users, err := seedUsers()
	if err != nil {
		return err
	}

	// SAFETY GUARD: If no users were created, STOP here.
	if len(users) == 0 {
		log.Println("🛑 No users found or created. Stopping seeder to prevent crashes.")
		return nil
	}

	// These functions now only run if we have users
	seedFriendships(users)
	groups, _ := seedGroups(users)

	if len(groups) > 0 {
		seedGroupMembers(groups, users)
	}

	return seedChatMessages(users, groups)
}

func seedUsers() ([]models.User, error) {
	if *dummyUsers {
		return seedDummyUsers()
	}

	if *skipFirebase {
		log.Println("⏭️  Skipping user seeding (skip-firebase flag set)")
		return nil, nil
	}

	log.Println("🔥 Fetching from Firebase...")
	fbUsers, err := listFirebaseUsers()
	if err != nil {
		log.Printf("❌ Error listing Firebase users: %v", err)
		return nil, nil
	}

	if len(fbUsers) == 0 {
		log.Println("ℹ️  No users found in Firebase. Skipping user seeding.")
		return nil, nil
	}

	var users []models.User
	for _, fb := range fbUsers {
		u := models.User{
			FirebaseUID: fb.UID,
			Email:       fb.Email,
			FullName:    fb.DisplayName,
			AvatarURL:   fb.PhotoURL,
			Status:      "offline",
		}
		config.DB.MySQL.Create(&u)
		users = append(users, u)
	}
	return users, nil
}

func seedDummyUsers() ([]models.User, error) {
	log.Println("🎭 Creating dummy users...")
	users := []models.User{
		{FirebaseUID: "uid_1", Email: "alice@example.com", FullName: "Alice", Status: "online"},
		{FirebaseUID: "uid_2", Email: "bob@example.com", FullName: "Bob", Status: "online"},
	}
	for i := range users {
		config.DB.MySQL.Create(&users[i])
	}
	return users, nil
}

func listFirebaseUsers() ([]*auth.ExportedUserRecord, error) {
	var users []*auth.ExportedUserRecord
	iter := config.FirebaseAuth.Users(context.Background(), "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func seedFriendships(users []models.User) {
	if len(users) < 2 {
		return
	}
	log.Println("🤝 Seeding friendships...")
	now := time.Now()
	config.DB.MySQL.Create(&models.Friendship{
		UserID: users[0].ID, FriendID: users[1].ID, Status: "accepted", AcceptedAt: &now,
	})
}

func seedGroups(users []models.User) ([]models.Group, error) {
	if len(users) == 0 {
		return nil, nil
	}
	log.Println("👥 Seeding groups...")
	group := models.Group{Name: "Alpha Team", OwnerID: users[0].ID}
	config.DB.MySQL.Create(&group)
	return []models.Group{group}, nil
}

func seedGroupMembers(groups []models.Group, users []models.User) {
	log.Println("👨‍👩‍👧‍👦 Seeding members...")
	config.DB.MySQL.Create(&models.GroupMember{
		GroupID: groups[0].ID, UserID: users[0].ID, Role: "admin",
	})
}

func seedChatMessages(users []models.User, groups []models.Group) error {
	// FINAL GUARD: Prevent index out of range [0]
	if len(users) < 2 {
		log.Println("⚠️  Not enough users for chat messages.")
		return nil
	}

	log.Println("💬 Seeding messages...")
	msg := models.ChatMessage{
		ID:          primitive.NewObjectID(),
		Content:     "Hello!",
		SenderID:    users[0].ID,
		RecipientID: &users[1].ID,
		CreatedAt:   time.Now(),
	}
	_, err := config.DB.MongoDB.Collection("chat_messages").InsertOne(context.Background(), msg)
	return err
}
