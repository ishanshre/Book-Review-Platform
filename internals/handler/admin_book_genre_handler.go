package handler

import (
	"fmt"
	"log"
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
	bookGenres, err := m.DB.AllBookGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var bookGenre models.BookGenre
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
	data["bookGenres"] = bookGenres
	data["bookGenre"] = bookGenre
	data["allGenres"] = allGenres
	data["allBooks"] = allBooks
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

	bookGenres, err := m.DB.AllBookGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

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

	data["allBooks"] = allBooks
	data["allGenres"] = allGenres
	data["bookGenre"] = bookGenre
	data["bookGenres"] = bookGenres
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

	if !form.Valid() {
		log.Println("invlaiud")
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
