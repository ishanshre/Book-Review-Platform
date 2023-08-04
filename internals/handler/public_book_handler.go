package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-faker/faker/v4"
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
	allBooks, err := m.DB.AllBookRandomPage(1000, 1)
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
	// limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	// if err != nil {
	// 	limit = 10
	// }
	// page, err := strconv.Atoi(r.URL.Query().Get("page"))
	// if err != nil {
	// 	page = 1
	// }
	// books, err := m.DB.AllBookData(limit, page)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }
	// data := make(map[string]interface{})
	// data["books"] = books
	render.Template(w, r, "public_books.page.tmpl", &models.TemplateData{})
}

// BookDetailByISBN returns the detail of the book using ISBN
func (m *Repository) BookDetailByISBN(w http.ResponseWriter, r *http.Request) {
	isbn, err := strconv.ParseInt(chi.URLParam(r, "isbn"), 10, 64)
	if err != nil {
		helpers.PageNotFound(w, err)
		return
	}
	book, err := m.DB.BookDetailWithAuthorPublisherWithIsbn(isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.PageNotFound(w, err)
			return
		}
		helpers.ServerError(w, err)
		return
	}
	authors := book.AuthorsData
	publisher := book.BookWithPublisherData.Publisher

	reviews, err := m.DB.GetReviewsByBookID(book.BookWithPublisherData.ID)
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
	data["book"] = book.BookWithPublisherData
	data["authors"] = authors
	data["publisher"] = publisher
	data["reviewDatas"] = reviewDatas
	data["averageRating"] = averageRating
	data["lastIndexAuthors"] = len(authors) - 1
	render.Template(w, r, "public_book_detail.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AllBooksFilterApi(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	searchKey := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "asc"
	}
	filteredBooks, err := m.DB.AllBooksFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBooks)
}

func (m *Repository) PopulateFakeData(w http.ResponseWriter, r *http.Request) {
	for i := 3; i < 50; i++ {
		publisher := &models.Publisher{
			Name:            faker.Word(),
			Description:     faker.Paragraph(),
			Pic:             "public/publisher/pic-publisher-2222.jpg",
			Address:         faker.Word(),
			Phone:           helpers.RandomPhone(10),
			Email:           faker.Email(),
			Website:         fmt.Sprintf("www.%s.com", faker.Word()),
			EstablishedDate: int(helpers.RandomInt(int64(1800), int64(2023))),
			Latitude:        "71.121212",
			Longitude:       "98.12121",
		}
		if err := m.DB.InsertPublisher(publisher); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			helpers.ServerError(w, err)
			return
		}
	}
	for i := 0; i < 150; i++ {
		book := &models.Book{
			Title:         faker.Word(),
			Description:   faker.Paragraph(),
			Cover:         "public/book/cover-book-1234565432123.jpeg",
			Isbn:          helpers.RandomInt(int64(1000000000000), int64(9999999999999)),
			PublishedDate: time.Now(),
			Paperback:     int(helpers.RandomInt(int64(100), int64(100000))),
			IsActive:      true,
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			PublisherID:   30,
		}
		if err := m.DB.InsertBook(book); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			helpers.ServerError(w, err)
			return
		}
	}
	helpers.ApiStatusOk(w, "success in populating data")
}
