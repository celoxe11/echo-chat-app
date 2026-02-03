package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database holds all database connections
type Database struct {
	MySQL   *gorm.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

var DB *Database

// InitDatabases initializes all database connections
func InitDatabases() error {
	DB = &Database{}

	// Initialize MySQL
	if err := initMySQL(); err != nil {
		return fmt.Errorf("failed to initialize MySQL: %w", err)
	}

	// Initialize MongoDB
	if err := initMongoDB(); err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	// Initialize Redis
	if err := initRedis(); err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	log.Println("✅ All databases initialized successfully")
	return nil
}

// initMySQL initializes MySQL connection using GORM
func initMySQL() error {
	// Load .env file
	godotenv.Load()

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	database := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB.MySQL = db
	log.Println("✅ MySQL connected successfully")
	return nil
}

// initMongoDB initializes MongoDB connection
func initMongoDB() error {
	godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	DB.MongoDB = client.Database(dbName)
	log.Println("✅ MongoDB connected successfully")
	return nil
}

// initRedis initializes Redis connection
func initRedis() error {
	godotenv.Load()

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return fmt.Errorf("invalid REDIS_DB value: %w", err)
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping Redis to verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	DB.Redis = client
	log.Println("✅ Redis connected successfully")
	return nil
}

// CloseDatabases closes all database connections gracefully
func CloseDatabases() {
	// Close MySQL
	if DB.MySQL != nil {
		sqlDB, err := DB.MySQL.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("MySQL connection closed")
		}
	}

	// Close MongoDB
	if DB.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		DB.MongoDB.Client().Disconnect(ctx)
		log.Println("MongoDB connection closed")
	}

	// Close Redis
	if DB.Redis != nil {
		DB.Redis.Close()
		log.Println("Redis connection closed")
	}
}
