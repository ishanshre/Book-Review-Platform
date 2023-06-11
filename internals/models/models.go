package models

import (
	"time"
)

// User is a type struct which holds users table data
type User struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Gender            string    `json:"gender"`
	Address           string    `json:"address"`
	Phone             string    `json:"phone"`
	ProfilePic        string    `json:"profile_pic"`
	CitizenshipNumber string    `json:"citizenship_number"`
	CitizenshipFront  string    `json:"citizenship_front"`
	CitizenshipBack   string    `json:"citizenship_back"`
	AccessLevel       int       `json:"access_level"`
	IsValidated       bool      `json:"is_validated"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastLogin         time.Time `json:"last_login"`
}

// MailData holds the email message
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}

// ResetPassword stores the new and confirm password
type ResetPassword struct {
	Token              string
	NewPassword        string
	NewPasswordConfirm string
}

type Genre struct {
	ID    int
	Title string
}
