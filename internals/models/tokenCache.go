package models

import "time"

// PasswordResetToken stores the token, email, and expire time
type PasswordResetToken struct {
	Token     string
	Email     string
	ExpiresAt time.Time
}

// UserTokenStore stores the user and token data
type UserTokenStore struct {
	Users             map[string]User
	PasswordResetRepo map[string]PasswordResetToken
}
