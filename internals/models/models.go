package models

import "time"

type User struct {
	ID          int       `json:"id"`
	FullName    string    `json:"full_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLogin   time.Time `json:"last_login"`
}
