package models

import (
	"gorm.io/gorm"
	"time"
)

// Group represents a chat group (stored in MySQL)
type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	AvatarURL   string         `gorm:"size:255" json:"avatar_url"`
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"`
	IsPrivate   bool           `gorm:"default:false" json:"is_private"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Owner   User   `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []User `gorm:"many2many:group_members;" json:"members,omitempty"`
}

// GroupMember represents the many-to-many relationship between users and groups
type GroupMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	GroupID  uint      `gorm:"not null;index:idx_group_user" json:"group_id"`
	UserID   uint      `gorm:"not null;index:idx_group_user" json:"user_id"`
	Role     string    `gorm:"type:enum('admin','moderator','member');default:'member'" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Group) TableName() string {
	return "groups"
}

func (GroupMember) TableName() string {
	return "group_members"
}
