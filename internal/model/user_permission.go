package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents the database model of permissions
type UserPermission struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	Model  string    `gorm:"not null"`
	Action string    `gorm:"not null"`
}

func (userPermission *UserPermission) BeforeCreate(tx *gorm.DB) error {
	// UUID version 4
	userPermission.ID = uuid.New()

	return nil
}
