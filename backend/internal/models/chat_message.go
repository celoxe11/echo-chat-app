package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatMessage represents a chat message stored in MongoDB
type ChatMessage struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Message details
	Content string `bson:"content" json:"content"`
	Type    string `bson:"type" json:"type"` // "text", "image", "file", "audio", "video"

	// Sender information (reference to MySQL User ID)
	SenderID uint `bson:"sender_id" json:"sender_id"`

	// Conversation context
	// For direct messages: recipient_id is set, group_id is null
	// For group messages: group_id is set, recipient_id is null
	RecipientID *uint `bson:"recipient_id,omitempty" json:"recipient_id,omitempty"`
	GroupID     *uint `bson:"group_id,omitempty" json:"group_id,omitempty"`

	// Message metadata
	IsEdited  bool `bson:"is_edited" json:"is_edited"`
	IsDeleted bool `bson:"is_deleted" json:"is_deleted"`

	// Attachments (for media messages)
	Attachments []Attachment `bson:"attachments,omitempty" json:"attachments,omitempty"`

	// Reply/Thread context
	ReplyToID *primitive.ObjectID `bson:"reply_to_id,omitempty" json:"reply_to_id,omitempty"`

	// Read receipts
	ReadBy []ReadReceipt `bson:"read_by,omitempty" json:"read_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// Attachment represents a file attachment in a message
type Attachment struct {
	URL       string `bson:"url" json:"url"`
	FileName  string `bson:"file_name" json:"file_name"`
	FileSize  int64  `bson:"file_size" json:"file_size"`
	MimeType  string `bson:"mime_type" json:"mime_type"`
	Thumbnail string `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
}
