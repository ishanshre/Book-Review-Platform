package handler

import (
	"database/sql"
	"errors"
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
	recentBooks, err := m.DB.AllRecentBooks(8, 1)
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
	data["allGenres"] = allGenres
	data["allBooks"] = allBooks
	data["recentBooks"] = recentBooks
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
	genres, err := m.DB.GetGenresFromBookID(book.BookWithPublisherData.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	languages, err := m.DB.GetLanguagesFromBookID(book.BookWithPublisherData.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book.BookWithPublisherData
	data["authors"] = authors
	data["publisher"] = publisher
	data["genres"] = genres
	data["languages"] = languages
	data["reviewDatas"] = reviewDatas
	data["averageRating"] = averageRating
	data["lastIndexAuthors"] = len(authors) - 1
	data["lastIndexGenres"] = len(genres) - 1
	data["lastIndexLanguages"] = len(languages) - 1
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
	for i := 3; i < 31; i++ {
		publisher := &models.Publisher{
			Name:            faker.Word(),
			Description:     faker.Paragraph(),
			Pic:             "public/publisher/pic-publisher-2015fbr9axdq.png",
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
	for i := 0; i < 15; i++ {
		book := &models.Book{
			Title:         faker.Word(),
			Description:   faker.Paragraph(),
			Cover:         "public/book/cover-book-9874563210321.jpg",
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

func (m *Repository) BookReadListExistsApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.ReadListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !exists {
		helpers.WriteJson(w, http.StatusOK, map[string]bool{
			"exists": false,
		})
		return
	}
	helpers.WriteJson(w, http.StatusOK, map[string]bool{
		"exists": true,
	})
}

func (m *Repository) BookBuyListExistsApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.BuyListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !exists {
		helpers.WriteJson(w, http.StatusOK, map[string]bool{
			"exists": false,
		})
		return
	}
	helpers.WriteJson(w, http.StatusOK, map[string]bool{
		"exists": true,
	})
}

func (m *Repository) AddtoReadListApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.ReadListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		helpers.ServerError(w, errors.New("book already in read list"))
		return
	}
	readList := &models.ReadList{
		BookID:    book_id,
		UserID:    user_id,
		CreatedAt: time.Now(),
	}
	if err := m.DB.InsertReadList(readList); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "add to read list success")
}

func (m *Repository) RemoveFromReadListApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.ReadListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !exists {
		helpers.ServerError(w, errors.New("book does not exists in read list"))
		return
	}
	if err := m.DB.DeleteReadList(user_id, book_id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "Removed from read list")
}

func (m *Repository) AddtoBuyListApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.BuyListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		helpers.ServerError(w, errors.New("book already in buy list"))
		return
	}
	buyList := &models.BuyList{
		BookID:    book_id,
		UserID:    user_id,
		CreatedAt: time.Now(),
	}
	if err := m.DB.InsertBuyList(buyList); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "add to buy list success")
}

func (m *Repository) RemoveFromBuyListApi(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	exists, err := m.DB.BuyListExists(user_id, book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !exists {
		helpers.ServerError(w, errors.New("book does not exists in buy list"))
		return
	}
	if err := m.DB.DeleteBuyList(user_id, book_id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "Removed from buy list")
}
