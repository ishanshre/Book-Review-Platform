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

// Genre holds the genres table model
type Genre struct {
	ID    int
	Title string
}

// Publisher holds the publishers table model
type Publisher struct {
	ID              int
	Name            string
	Description     string
	Pic             string
	Address         string
	Phone           string
	Email           string
	Website         string
	EstablishedDate int
	Latitude        string
	Longitude       string
}

// Author struct holds the authors table data
type Author struct {
	ID              int
	FirstName       string
	LastName        string
	Bio             string
	DateOfBirth     int
	Email           string
	CountryOfOrigin string
	Avatar          string
}

// Language struct holds the language table model
type Language struct {
	ID       int
	Language string
}

// Book struct holds the books table data
type Book struct {
	ID            int
	Title         string
	Description   string
	Cover         string
	Isbn          int64
	PublishedDate time.Time
	Paperback     int
	IsActive      bool
	AddedAt       time.Time
	UpdatedAt     time.Time
	PublisherID   int
}

// BookAuthor struct holds the immediate table between book and author
type BookAuthor struct {
	BookID   int
	AuthorID int
}

// BookGenre holds the itermediate table between book and genre
type BookGenre struct {
	BookID  int
	GenreID int
}

// BookLanguage holds the book id and language id
type BookLanguage struct {
	BookID     int
	LanguageID int
}
