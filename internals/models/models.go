package models

import (
	"time"
)

// User is a type struct which holds users table data
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLogin   time.Time `json:"last_login"`
}
type Kyc struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	ProfilePic     string    `json:"profile_pic"`
	DateOfBirth    time.Time `json:"dob"`
	DocumentType   string    `json:"document_type"`
	DocumentNumber string    `json:"document_number"`
	DocumentFront  string    `json:"document_front"`
	DocumentBack   string    `json:"document_back"`
	IsValidated    bool      `json:"is_validated"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AdminUserList struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
	IsValidated bool      `json:"is_validated"`
}

type AdminUserListApi struct {
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	LastPage int              `json:"last_page"`
	Users    []*AdminUserList `json:"users"`
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

type AdminPublisherList struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	EstablishedDate int    `json:"established_date"`
}

type AdminPublisherListApi struct {
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	LastPage   int                   `json:"last_page"`
	Publishers []*AdminPublisherList `json:"publishers"`
}

// Author struct holds the authors table data
type Author struct {
	ID              int    `json:"id,omitempty"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Bio             string `json:"bio,omitempty"`
	DateOfBirth     int    `json:"dob,omitempty"`
	Email           string `json:"email,omitempty"`
	CountryOfOrigin string `json:"coo,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
}

// Language struct holds the language table model
type Language struct {
	ID       int
	Language string
}

// Book struct holds the books table data
type Book struct {
	ID            int       `json:"id,omitempty"`
	Title         string    `json:"title,omitempty"`
	Description   string    `json:"description,omitempty"`
	Cover         string    `json:"cover,omitempty"`
	Isbn          int64     `json:"isbn,omitempty"`
	PublishedDate time.Time `json:"published_date,omitempty"`
	Paperback     int       `json:"paperback,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	AddedAt       time.Time `json:"added_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	PublisherID   int       `json:"publisher_id,omitempty"`
}

type BookApiFilter struct {
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	LastPage int     `json:"last_page"`
	Books    []*Book `json:"books"`
}

type AuthorApiFilter struct {
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	LastPage int       `json:"last_page"`
	Authors  []*Author `json:"authors"`
}

// BookAuthor struct holds the immediate table between book and author
type BookAuthor struct {
	BookID   int
	AuthorID int
}

type BookAuthorList struct {
	BookID          int    `json:"book_id"`
	BookTitle       string `json:"book_title"`
	AuthorID        int    `json:"author_id"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
}

type BookAuthorListApi struct {
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	LastPage    int               `json:"last_page"`
	BookAuthors []*BookAuthorList `json:"book_authors"`
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

// ReadList holds the the book id, user id and created at of readLists table
type ReadList struct {
	UserID    int
	BookID    int
	CreatedAt time.Time
}

type ReadListFilter struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	BookID    int       `json:"book_id"`
	BookTitle string    `json:"book_title"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadListFilterApi struct {
	Total           int               `json:"total"`
	Page            int               `json:"page"`
	LastPage        int               `json:"last_page"`
	ReadListFilters []*ReadListFilter `json:"read_lists"`
}

// BuyList holds the the book id, user id and created at of BuyLists table
type BuyList struct {
	UserID    int
	BookID    int
	CreatedAt time.Time
}

type BuyListFilter struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	BookID    int       `json:"book_id"`
	BookTitle string    `json:"book_title"`
	CreatedAt time.Time `json:"created_at"`
}

type BuyListFilterApi struct {
	Total          int              `json:"total"`
	Page           int              `json:"page"`
	LastPage       int              `json:"last_page"`
	BuyListFilters []*BuyListFilter `json:"buy_lists"`
}

// Follower hold the book and author id for follower relationship
type Follower struct {
	UserID     int
	AuthorID   int
	FollowedAt time.Time
}

// Review holds the review table data
type Review struct {
	ID        int
	Rating    float64
	Body      string
	BookID    int
	UserID    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Contact struct {
	ID            int
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	Subject       string
	Message       string
	SubmittedAt   time.Time
	IpAddress     string
	BrowserInfo   string
	ReferringPage string
}

type RequestedBook struct {
	ID             int
	BookTitle      string
	Author         string
	RequestedEmail string
	RequestedDate  time.Time
}
