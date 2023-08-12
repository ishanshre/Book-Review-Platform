package models

import "time"

// Embed the book and author data for BookAuthor
type BookAuthorData struct {
	BookData   *Book
	AuthorData *Author
}

type BookWithPublisher struct {
	ID            int        `json:"id,omitempty"`
	Title         string     `json:"title,omitempty"`
	Description   string     `json:"description,omitempty"`
	Cover         string     `json:"cover,omitempty"`
	Isbn          int64      `json:"isbn,omitempty"`
	PublishedDate time.Time  `json:"published_date,omitempty"`
	Paperback     int        `json:"paperback,omitempty"`
	IsActive      bool       `json:"is_active,omitempty"`
	AddedAt       time.Time  `json:"added_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at,omitempty"`
	Publisher     *Publisher `json:"publisher,omitempty"`
}

type BookInfoData struct {
	BookWithPublisherData *BookWithPublisher
	AuthorsData           []*Author
}

type AuthorBookData struct {
	Author *Author
	Books  []*Book
}

// Embed the book and genre data for BookGenre
type BookGenreData struct {
	BookData  *Book
	GenreData *Genre
}

// Embed the book and language data for BookLanguage
type BookLanguageData struct {
	BookData     *Book
	LanguageData *Language
}

// Embed the book and user data for ReadList
type ReadListData struct {
	BookData  *Book
	UserData  *User
	CreatedAt time.Time
}

// Embed the book and user data for BuyList
type BuyListData struct {
	BookData  *Book
	UserData  *User
	CreatedAt time.Time
}

// Embed the user and author data for Follower Model
type FollowerData struct {
	UserData   *User
	AuthorData *Author
	FollowedAt time.Time
}

type ReviewUserData struct {
	Review *Review
	User   *User
}

type UserKycData struct {
	User *User
	Kyc  *Kyc
}

type AuthorFollowerData struct {
	Author []*Author
}

type PublisherWithBooksData struct {
	Publisher *Publisher
	Books     []*Book
}
