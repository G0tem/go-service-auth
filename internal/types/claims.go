package types

import "time"

// JwtClaims represents minimal set of fields extracted from JWT token
// and propagated through Fiber context.
type JwtClaims struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	Exp         time.Time `json:"exp"`
}

