package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Start of handler for admin book-genre relationship

// AdminAllBookGenre retrieves all book-genre relationships in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function retrieves all book-genre relationships, all books, and all genres from the database.
// If any errors occur during the retrieval process, a server error is returned.
// The function prepares the necessary data and renders the "admin-allbookgenres.page.tmpl" template,
// displaying the list of book-genre relationships as well as new book genre relationship add form.
func (m *Repository) AdminAllBookGenre(w http.ResponseWriter, r *http.Request) {
	var bookGenre models.BookGenre
	allBooks, allGenres, bookGenres, bookGenreDatas, err := m.genericBookGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["bookGenres"] = bookGenres
	data["bookGenreDatas"] = bookGenreDatas
	data["bookGenre"] = bookGenre
	data["allGenres"] = allGenres
	data["allBooks"] = allBooks
	data["base_path"] = base_bookGenres_path
	render.Template(w, r, "admin-allbookgenres.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeleteBookGenre handles the deletion of a book-genre relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function extracts the book ID and genre ID from the URL path parameters,
// deletes the book-genre relationship from the database, and redirects the user to the "/admin/bookGenres" page.
// If any errors occur during the process, a server error is returned.
func (m *Repository) PostAdminDeleteBookGenre(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre_id, err := strconv.Atoi(chi.URLParam(r, "genre_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteBookGenre(book_id, genre_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Genre Relationship Deleted")

	http.Redirect(w, r, "/admin/bookGenres", http.StatusSeeOther)
}

// AdminGetBookGenreByID handes the detail logic
func (m *Repository) AdminGetBookGenreByID(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre_id, err := strconv.Atoi(chi.URLParam(r, "genre_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookGenre, err := m.DB.GetBookGenreByID(book_id, genre_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id
	genre, err := m.DB.GetGenreByID(genre_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre.ID = genre_id
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allGenres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["genre"] = genre
	data["allGenres"] = allGenres
	data["bookGenre"] = bookGenre
	data["base_path"] = base_bookGenres_path
	render.Template(w, r, "admin-bookgenredetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// AdminGetBookGenreByID retrieves information related to a specific book-genre relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function extracts the book ID and genre ID from the URL path parameters,
// retrieves the book-genre relationship, book title, genre information, all books, and all genres from the database.
// It then prepares the necessary data and renders the "admin-bookgenredetail.page.tmpl" template to display the details.
// If any errors occur during the process, a server error is returned.
func (m *Repository) PostAdminUpdateBookGenre(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre_id, err := strconv.Atoi(chi.URLParam(r, "genre_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	updated_book_id, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	updated_genre_id, err := strconv.Atoi(r.Form.Get("genre_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookGenre := models.BookGenre{
		BookID:  updated_book_id,
		GenreID: updated_genre_id,
	}
	exists, err := m.DB.BookGenreExists(bookGenre.BookID, bookGenre.GenreID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-author relationship already exists")
		form.Errors.Add("genre_id", "book-author relationship already exists")
	}
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id
	genre, err := m.DB.GetGenreByID(genre_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre.ID = genre_id
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allGenres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["genre"] = genre
	data["allGenres"] = allGenres
	data["bookGenre"] = bookGenre
	data["base_path"] = base_bookGenres_path

	form.Required("book_id", "genre_id")
	if !form.Valid() {
		render.Template(w, r, "admin-bookgenredetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateBookGenre(&bookGenre, book_id, genre_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Genre Relation Updated")

	http.Redirect(w, r, fmt.Sprintf("/admin/bookGenres/detail/%d/%d", bookGenre.BookID, bookGenre.GenreID), http.StatusSeeOther)
}

// PostAdminInsertBookGenre handles the insertion of a new book-genre relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request and validates it. If any parsing or validation errors occur,
// a server error is returned. The function retrieves the book ID and genre ID from the form data and creates a new
// BookGenre object with the provided IDs. It then retrieves all book-genre relationships, all books, and all genres
// from the database to prepare the necessary data for rendering the template. The function checks if the book-genre
// relationship already exists and adds an error to the form if it does. If the form is not valid, the template is
// rendered with the form and data. If the form is valid, the new book-genre relationship is inserted into the database
// and the user is redirected to the "/admin/bookGenres" page.
func (m *Repository) PostAdminInsertBookGenre(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	data := make(map[string]interface{})

	book_id, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	genre_id, err := strconv.Atoi(r.Form.Get("genre_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	bookGenre := models.BookGenre{
		BookID:  book_id,
		GenreID: genre_id,
	}

	allBooks, allGenres, bookGenres, bookGenreDatas, err := m.genericBookGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["allBooks"] = allBooks
	data["allGenres"] = allGenres
	data["bookGenre"] = bookGenre
	data["bookGenres"] = bookGenres
	data["bookGenreDatas"] = bookGenreDatas
	form.Required("book_id", "genre_id")

	exists, err := m.DB.BookGenreExists(bookGenre.BookID, bookGenre.GenreID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-genre relationship already exists")
		form.Errors.Add("genre_id", "book-genre relationship already exists")
	}

	data["base_path"] = base_bookGenres_path

	if !form.Valid() {
		render.Template(w, r, "admin-allbookgenres.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err := m.DB.InsertBookGenre(&bookGenre); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Genre Relationship Added")

	http.Redirect(w, r, "/admin/bookGenres", http.StatusSeeOther)
}

func (m *Repository) genericBookGenre() ([]*models.Book, []*models.Genre, []*models.BookGenre, []*models.BookGenreData, error) {
	allBookGenresCh := make(chan []*models.BookGenre)
	allBooksCh := make(chan []*models.Book)
	allGenresCh := make(chan []*models.Genre)
	errorCh := make(chan error)

	go func() {
		books, err := m.DB.AllBook()
		if err != nil {
			errorCh <- err
			return
		}
		allBooksCh <- books
	}()
	go func() {
		allBookGenres, err := m.DB.AllBookGenre()
		if err != nil {
			errorCh <- err
			return
		}
		allBookGenresCh <- allBookGenres
	}()
	go func() {
		allGenres, err := m.DB.AllGenre()
		if err != nil {
			errorCh <- err
			return
		}
		allGenresCh <- allGenres
	}()

	var allBooks []*models.Book
	var allGenres []*models.Genre
	var bookGenres []*models.BookGenre
	var err error
	for i := 0; i < 3; i++ {
		select {
		case allBooks = <-allBooksCh:
		case allGenres = <-allGenresCh:
		case bookGenres = <-allBookGenresCh:
		case err = <-errorCh:
		}
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}

	bookGenreDatas := []*models.BookGenreData{}
	for _, v := range bookGenres {
		book, err := m.DB.GetBookTitleByID(v.BookID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		genre, err := m.DB.GetGenreByID(v.GenreID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		bookGenreData := &models.BookGenreData{
			BookData:  book,
			GenreData: genre,
		}
		bookGenreDatas = append(bookGenreDatas, bookGenreData)
	}
	return allBooks, allGenres, bookGenres, bookGenreDatas, nil
}
