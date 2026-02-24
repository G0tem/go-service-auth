package handler

import (
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/G0tem/go-service-auth/internal/model"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *Handler) getUserByEmail(e string) (*model.User, error) {
	var user model.User
	if err := h.db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (h *Handler) getUserByUsername(u string) (*model.User, error) {
	var user model.User
	if err := h.db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// type SortDirection string

// const (
// 	SortDirectionUnsorted SortDirection = ""
// 	SortDirectionAsc      SortDirection = "ASC"
// 	SortDirectionDesc     SortDirection = "DESC"
// )

// type SortByField struct {
// 	FieldName     string
// 	SortDirection SortDirection
// }
