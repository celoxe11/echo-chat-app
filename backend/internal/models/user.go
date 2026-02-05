package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user in the system (stored in MySQL)
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FirebaseUID string       `gorm:"uniqueIndex;not null;size:100" json:"firebase_uid"`
	Email     string         `gorm:"uniqueIndex;not null;size:100" json:"email"`
	FullName  string         `gorm:"size:100" json:"full_name"`
	Username  string         `gorm:"size:100" json:"username"`
	AvatarURL string         `gorm:"size:255" json:"avatar_url"`
	Status    string         `gorm:"type:enum('online','offline','away');default:'offline'" json:"status"`
	LastSeen  time.Time      `json:"last_seen"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Friends     []User  `gorm:"many2many:friendships;" json:"friends,omitempty"`
	Groups      []Group `gorm:"many2many:group_members;" json:"groups,omitempty"`
	OwnedGroups []Group `gorm:"foreignKey:OwnerID" json:"owned_groups,omitempty"`
}


// TableName overrides for custom table names (optional)
func (User) TableName() string {
	return "users"
}