package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents the database model of roles
type UserRole struct {
	ID          uuid.UUID `gorm:"primarykey;uniqueIndex;not null;type:uuid;" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:50;" json:"name"`
	Description string    `gorm:"size:255;nullable" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (userRole *UserRole) BeforeCreate(tx *gorm.DB) error {
	userRole.ID = uuid.New()
	return nil
}

func (userRole *UserRole) TableName() string {
	return "user_roles"
}
