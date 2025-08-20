package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User settings for marketplaces
type UserMarketplaceSettings struct {
	ID          uint      `gorm:"primarykey;column:id"`
	Marketplace string    `gorm:"size:255;not null;unique;column:marketplace"`
	ApiKey      string    `gorm:"size:1024;column:api_key"`
	UserID      uuid.UUID `gorm:"not null;type:uuid;column:user_id"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (UserMarketplaceSettings) TableName() string {
	return "usermarketplacesettingss"
}
