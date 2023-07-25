package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

// AllBooks return all books in pages
func (m *Repository) AllBooks(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	books, err := m.DB.AllBookData(limit, page)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["books"] = books
	render.Template(w, r, "public_books.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// BookDetailByISBN returns the detail of the book using ISBN
func (m *Repository) BookDetailByISBN(w http.ResponseWriter, r *http.Request) {
	isbn, err := strconv.ParseInt(chi.URLParam(r, "isbn"), 10, 64)
	if err != nil {
		helpers.PageNotFound(w, err)
		return
	}
	book, err := m.DB.GetBookByISBN(isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.PageNotFound(w, err)
			return
		}
		helpers.ServerError(w, err)
		return
	}
	bookAuthors, err := m.DB.GetBookAuthorByBookID(book.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	authors := []*models.Author{}
	for _, bookAuthor := range bookAuthors {
		author, err := m.DB.GetAuthorByID(bookAuthor.AuthorID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		authors = append(authors, author)
	}

	publisher, err := m.DB.GetPublisherByID(book.PublisherID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reviews, err := m.DB.GetReviewsByBookID(book.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var totalRatings float64
	var numReviews int
	var averageRating float64
	reviewDatas := []*models.ReviewUserData{}
	for _, review := range reviews {
		user, err := m.DB.GetGlobalUserByIDAny(review.UserID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		reviewData := &models.ReviewUserData{
			Review: review,
			User:   user,
		}
		reviewDatas = append(reviewDatas, reviewData)
		totalRatings += review.Rating
		numReviews++
	}
	if numReviews > 0 {
		averageRating = totalRatings / float64(numReviews)
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["authors"] = authors
	data["publisher"] = publisher
	data["reviewDatas"] = reviewDatas
	data["averageRating"] = averageRating
	data["lastIndexAuthors"] = len(authors) - 1
	render.Template(w, r, "public_book_detail.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
