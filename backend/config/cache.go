package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CacheService provides caching operations using Redis
type CacheService struct{}

var Cache = &CacheService{} 

// Set stores a value in Redis with expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return DB.Redis.Set(ctx, key, jsonData, expiration).Err()
}

// Get retrieves a value from Redis and unmarshals it
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := DB.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	return DB.Redis.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in Redis
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := DB.Redis.Exists(ctx, key).Result()
	return count > 0, err
}

// SetUserOnlineStatus sets a user's online status in Redis
func (c *CacheService) SetUserOnlineStatus(ctx context.Context, userID uint, status string) error {
	key := fmt.Sprintf("user:status:%d", userID)
	return DB.Redis.Set(ctx, key, status, 24*time.Hour).Err()
}

// GetUserOnlineStatus gets a user's online status from Redis
func (c *CacheService) GetUserOnlineStatus(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user:status:%d", userID)
	return DB.Redis.Get(ctx, key).Result()
}

// CacheUserSession caches user session data
func (c *CacheService) CacheUserSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return DB.Redis.Set(ctx, key, userID, expiration).Err()
}

// GetUserSession retrieves user ID from session
func (c *CacheService) GetUserSession(ctx context.Context, sessionID string) (uint, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := DB.Redis.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	var userID uint
	if err := json.Unmarshal([]byte(val), &userID); err != nil {
		return 0, err
	}

	return userID, nil
}

// InvalidateUserSession removes a session from cache
func (c *CacheService) InvalidateUserSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return DB.Redis.Del(ctx, key).Err()
}

// CacheConversationList caches a user's conversation list
func (c *CacheService) CacheConversationList(ctx context.Context, userID uint, conversations interface{}) error {
	key := fmt.Sprintf("user:conversations:%d", userID)
	return c.Set(ctx, key, conversations, 5*time.Minute)
}

// GetCachedConversationList retrieves cached conversation list
func (c *CacheService) GetCachedConversationList(ctx context.Context, userID uint, dest interface{}) error {
	key := fmt.Sprintf("user:conversations:%d", userID)
	return c.Get(ctx, key, dest)
}

// InvalidateConversationCache invalidates conversation cache for users
func (c *CacheService) InvalidateConversationCache(ctx context.Context, userIDs ...uint) error {
	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = fmt.Sprintf("user:conversations:%d", userID)
	}
	return c.Delete(ctx, keys...)
}

// IncrementUnreadCount increments unread message count for a user
func (c *CacheService) IncrementUnreadCount(ctx context.Context, userID uint, conversationID string) error {
	key := fmt.Sprintf("user:%d:unread:%s", userID, conversationID)
	return DB.Redis.Incr(ctx, key).Err()
}

// GetUnreadCount gets unread message count
func (c *CacheService) GetUnreadCount(ctx context.Context, userID uint, conversationID string) (int64, error) {
	key := fmt.Sprintf("user:%d:unread:%s", userID, conversationID)
	return DB.Redis.Get(ctx, key).Int64()
}

// ResetUnreadCount resets unread message count
func (c *CacheService) ResetUnreadCount(ctx context.Context, userID uint, conversationID string) error {
	key := fmt.Sprintf("user:%d:unread:%s", userID, conversationID)
	return DB.Redis.Set(ctx, key, 0, 24*time.Hour).Err()
}
