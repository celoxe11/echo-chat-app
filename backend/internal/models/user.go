package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user in the system (stored in MySQL)
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // "-" means don't include in JSON
	FullName  string         `gorm:"size:100" json:"full_name"`
	AvatarURL string         `gorm:"size:255" json:"avatar_url"`
	Status    string         `gorm:"type:enum('online','offline','away');default:'offline'" json:"status"`
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