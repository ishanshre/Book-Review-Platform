package models

type BookWithAverageRating struct {
	Book          Book
	Authors       []Author
	LenAuthors    int
	AverageRating float64
	NumReviews    int
}
