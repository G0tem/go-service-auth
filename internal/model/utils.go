package model

import (
	"crypto/rand"
	"encoding/base64"
)

func UniqueRandomString(len int) string {
	token := make([]byte, len)
	rand.Reader.Read(token)
	return base64.RawURLEncoding.EncodeToString(token)[:len]
}
