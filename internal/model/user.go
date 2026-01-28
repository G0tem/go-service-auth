package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/G0tem/go-service-auth/internal"
)

// User struct
type User struct {
	ID             uuid.UUID `gorm:"primarykey;uniqueIndex;not null;type:uuid;"`
	Username       string    `gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email          string    `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	EmailConfirmed bool      `gorm:"not null;column:email_confirmed;" json:"email_confirmed"`
	PasswordHash   string    `gorm:"column:password_hash;" validate:"required,min=6,max=50" json:"password"`
	AvatarURL      string    `json:"avatar_url"`
	RoleID         uuid.UUID `gorm:"type:uuid;column:role_id" json:"role_id"`
	Role           UserRole  `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	// UUID version 4
	user.ID = uuid.New()
	return nil
}

func (user *User) GetAvatarUrl(cdnUrl string) string {
	if user.AvatarURL == "" {
		return ""
	}
	return internal.JoinUrl(cdnUrl, user.AvatarURL)
}

func (user *User) TableName() string {
	return "users"
}
