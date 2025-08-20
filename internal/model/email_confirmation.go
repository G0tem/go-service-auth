package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailTemplate struct
type EmailConfirmation struct {
	UserID    uuid.UUID `gorm:"primarykey;type:uuid;column:user_id;"`
	User      User      `gorm:"foreignKey:UserID;references:ID;"`
	Token     string    `gorm:"uniqueIndex;not null;"`
	Password  string    `gorm:"-"` // Временное поле, не сохраняется в БД
	Timestamp time.Time
}

func NewEmailConfirmation(userId uuid.UUID, password string) *EmailConfirmation {
	return &EmailConfirmation{
		UserID:    userId,
		Token:     UniqueRandomString(100),
		Password:  password,
		Timestamp: time.Now(),
	}
}

func (emailConfirmation *EmailConfirmation) BeforeCreate(tx *gorm.DB) error {
	// Unique link part (it's security safe)
	emailConfirmation.Token = UniqueRandomString(100)
	emailConfirmation.Timestamp = time.Now()

	return nil
}

func (*EmailConfirmation) TableName() string {
	return "email_confirmation_request"
}

func (*EmailConfirmation) AfterInsertTrigger() string {
	return `
DROP TRIGGER IF EXISTS email_confirmation_request_delete_old_rows_trigger ON email_confirmation_request;
DROP FUNCTION IF EXISTS email_confirmation_request_delete_old_rows;

CREATE FUNCTION email_confirmation_request_delete_old_rows() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  DELETE FROM email_confirmation_request WHERE timestamp < NOW() - INTERVAL '1 day';
  RETURN NEW;
END;
$$;

CREATE TRIGGER email_confirmation_request_delete_old_rows_trigger
    AFTER INSERT ON email_confirmation_request
    EXECUTE PROCEDURE email_confirmation_request_delete_old_rows();
	`
}
