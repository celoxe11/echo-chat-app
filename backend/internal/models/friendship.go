package models


import (
	"time"
)

// Friendship represents a friend relationship between two users
type Friendship struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index:idx_user_friend" json:"user_id"`
	FriendID   uint       `gorm:"not null;index:idx_user_friend" json:"friend_id"`
	Status     string     `gorm:"type:enum('pending','accepted','blocked');default:'pending'" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`

	// Relationships
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Friend User `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

func (Friendship) TableName() string {
	return "friendships"
}
