package models

import "time"

// Embed the book and author data for BookAuthor
type BookAuthorData struct {
	BookData   *Book
	AuthorData *Author
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
