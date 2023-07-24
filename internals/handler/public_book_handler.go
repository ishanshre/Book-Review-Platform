package handler

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	allGenres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allBooks, err := m.DB.AllBookDataRandom()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	topRatedBooks := []models.BookWithAverageRating{}
	for _, book := range allBooks {
		reviews, err := m.DB.GetReviewsByBookID(book.ID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		var totalRatings float64
		var numReviews int
		for _, review := range reviews {
			totalRatings += review.Rating
			numReviews++
		}
		if numReviews > 0 {
			averageRating := totalRatings / float64(numReviews)
			if averageRating > 4.0 {
				bookAuthors, err := m.DB.GetBookAuthorByBookID(book.ID)
				if err != nil {
					helpers.ServerError(w, err)
					return
				}
				authors := []models.Author{}
				for _, bookAuthor := range bookAuthors {
					author, err := m.DB.GetAuthorByID(bookAuthor.AuthorID)
					if err != nil {
						helpers.ServerError(w, err)
						return
					}
					authors = append(authors, *author)
				}
				topRatedBooks = append(topRatedBooks, models.BookWithAverageRating{
					Book:          *book,
					Authors:       authors,
					LenAuthors:    len(authors) - 1,
					AverageRating: averageRating,
					NumReviews:    numReviews,
				})
			}
		}
		if len(topRatedBooks) == 10 {
			break
		}
	}

	data := make(map[string]interface{})
	data["allGenre"] = allGenres
	data["allBooks"] = allBooks
	data["topRatedBooks"] = topRatedBooks
	render.Template(w, r, "public_home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
