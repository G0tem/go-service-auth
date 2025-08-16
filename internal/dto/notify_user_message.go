package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

/**

{"type": "@notify.OnLiveBroadcast", "user_info": {"user_id": "a08bf07f-2486-40bf-be9e-ca804b3dd841"}, "params": {"broadcast": "7ea02223-e3cf-4826-8229-5009c6d6f9d0"}}

**/

type NotifyUserInfo struct {
	UserId         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	EmailConfirmed bool      `json:"email_confirmed"`
	Description    string    `json:"description"`
	AvatarURL      string    `json:"avatar_url"`
	CoverURL       string    `json:"cover_url"`
	Roles          []string  `json:"roles"`
}

func (userInfo NotifyUserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(userInfo)
}

type NotifyUserMessage struct {
	Type          string          `json:"type"`
	UserInfo      *NotifyUserInfo `json:"user_info"`
	MessageParams interface{}     `json:"params"`
}

func (notifyMessage NotifyUserMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(notifyMessage)
}
