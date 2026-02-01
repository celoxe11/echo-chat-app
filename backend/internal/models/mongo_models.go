package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


// ReadReceipt tracks who has read a message
type ReadReceipt struct {
	UserID uint      `bson:"user_id" json:"user_id"`
	ReadAt time.Time `bson:"read_at" json:"read_at"`
}

// Conversation represents a conversation summary (for quick lookups)
// This is useful for displaying conversation lists
type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Participants (for direct messages, 2 users; for groups, reference to group_id)
	Participants []uint `bson:"participants" json:"participants"`
	GroupID      *uint  `bson:"group_id,omitempty" json:"group_id,omitempty"`

	// Last message info (denormalized for performance)
	LastMessageID   primitive.ObjectID `bson:"last_message_id" json:"last_message_id"`
	LastMessageText string             `bson:"last_message_text" json:"last_message_text"`
	LastMessageAt   time.Time          `bson:"last_message_at" json:"last_message_at"`
	LastSenderID    uint               `bson:"last_sender_id" json:"last_sender_id"`

	// Unread counts per user
	UnreadCounts map[string]int `bson:"unread_counts" json:"unread_counts"` // key: user_id as string

	// Timestamps
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// CollectionName returns the MongoDB collection name for ChatMessage
func (ChatMessage) CollectionName() string {
	return "chat_messages"
}

// CollectionName returns the MongoDB collection name for Conversation
func (Conversation) CollectionName() string {
	return "conversations"
}
