package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Start of handler for admin book-author

// AdminAllBookAuthor retrieves all book-author relationships and renders the "admin-allbookauthors.page.tmpl" template.
// It takes the HTTP response writer and request as parameters. The function calls the database's AllBookAuthor method
// to retrieve all book-author relationships. If an error occurs during the retrieval, a server error is returned.
// The function also retrieves all books and authors from the database using the AllBook and AllAuthor methods, respectively.
// If any errors occur during the retrieval, a server error is returned. The retrieved data, including book-authors,
// book-author, all authors, and all books, is stored in a data map. The function then renders the template,
// passing the data map and an empty form to the template for rendering.
func (m *Repository) AdminAllBookAuthor(w http.ResponseWriter, r *http.Request) {
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["allAuthors"] = allAuthors
	data["allBooks"] = allBooks
	data["base_path"] = base_bookAuthors_path
	render.Template(w, r, "admin-allbookauthors.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}
func (m *Repository) AdminAllBookAuthorsApi(w http.ResponseWriter, r *http.Request) {
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
	filteredBookAuthors, err := m.DB.BookAuthorListFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBookAuthors)
}

// PostAdminDeleteBookAuthor deletes a book-author relationship based on the provided book ID and author ID.
// It takes the HTTP response writer and request as parameters. The function retrieves the book ID and author ID
// from the URL parameters. If any parsing errors occur, a server error is returned. It calls the database's
// DeleteBookAuthor method to delete the corresponding book-author relationship. If an error occurs during the
// deletion, a server error is returned. Otherwise, the function redirects the user to the "/admin/bookAuthors"
// page with a status code of http.StatusSeeOther.
func (m *Repository) PostAdminDeleteBookAuthor(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteBookAuthor(book_id, author_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Deleted")

	http.Redirect(w, r, "/admin/bookAuthors", http.StatusSeeOther)
}

// AdminGetBookAuthorByID retrieves the details of a book-author relationship by its book ID and author ID.
// It takes the HTTP response writer and request as parameters. The function retrieves the book ID and author ID
// from the URL parameters. If any parsing errors occur, a server error is returned. It retrieves the book-author
// relationship, book, author, all books, and all authors from the database. The function prepares the necessary
// data for rendering the template and renders the "admin-bookauthordetial.page.tmpl" template with the form and data.
func (m *Repository) AdminGetBookAuthorByID(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookAuthor, err := m.DB.GetBookAuthorByID(book_id, author_id)
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
	author, err := m.DB.GetAuthorFullNameByID(author_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author.ID = author_id
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["author"] = author
	data["allAuthors"] = allAuthors
	data["bookAuthor"] = bookAuthor
	data["base_path"] = base_bookAuthors_path

	render.Template(w, r, "admin-bookauthordetial.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateBookAuthor handles the update logic of a book-author relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the book ID and author ID from the URL parameters and parses the form data from the request.
// If any parsing errors occur, a server error is returned. It creates a new form object and retrieves the updated
// book ID and author ID from the form data. The function checks if the updated book-author relationship already exists
// and adds an error to the form if it does. It retrieves the book and author details based on their IDs from the database,
// as well as all books and all authors. The function prepares the necessary data for rendering the template.
// The form data is then validated, and if the form is not valid, the template is rendered with the form and data.
// If the form is valid, the book-author relationship is updated in the database, and the user is redirected to the
// detail page of the updated book-author relationship.
func (m *Repository) PostAdminUpdateBookAuthor(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
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
	updated_author_id, err := strconv.Atoi(r.Form.Get("author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookAuthor := models.BookAuthor{
		BookID:   updated_book_id,
		AuthorID: updated_author_id,
	}
	exists, err := m.DB.BookAuthorExists(bookAuthor.BookID, bookAuthor.AuthorID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-author relationship already exists")
		form.Errors.Add("author_id", "book-author relationship already exists")
	}
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id
	author, err := m.DB.GetAuthorFullNameByID(author_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author.ID = author_id
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["author"] = author
	data["allAuthors"] = allAuthors
	data["bookAuthor"] = bookAuthor
	data["base_path"] = base_bookAuthors_path
	form.Required("book_id", "author_id")
	if !form.Valid() {
		render.Template(w, r, "admin-bookauthordetial.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateBookAuthor(&bookAuthor, book_id, author_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Updated")

	http.Redirect(w, r, fmt.Sprintf("/admin/bookAuthors/detail/%d/%d", bookAuthor.BookID, bookAuthor.AuthorID), http.StatusSeeOther)
}

// PostAdminInsertBookAuthor handles the insertion of a new book-author relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request and validates it. If any parsing or validation errors occur,
// a server error is returned. The function retrieves the book ID and author ID from the form data and creates a new
// BookAuthor object with the provided IDs. It then retrieves all book-author relationships, all books, and all authors
// from the database to prepare the necessary data for rendering the template. The function checks if the book-author
// relationship already exists and adds an error to the form if it does. If the form is not valid, the template is
// rendered with the form and data. If the form is valid, the new book-author relationship is inserted into the database
// and the user is redirected to the "/admin/bookAuthors" page.
func (m *Repository) PostAdminInsertBookAuthor(w http.ResponseWriter, r *http.Request) {
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

	author_id, err := strconv.Atoi(r.Form.Get("author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	bookAuthor := models.BookAuthor{
		BookID:   book_id,
		AuthorID: author_id,
	}

	bookAuthors, err := m.DB.AllBookAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookAuthorDatas := []*models.BookAuthorData{}
	for _, v := range bookAuthors {
		book, err := m.DB.GetBookTitleByID(v.BookID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		author, err := m.DB.GetAuthorByID(v.AuthorID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		bookAuthorData := models.BookAuthorData{
			BookData:   book,
			AuthorData: author,
		}
		bookAuthorDatas = append(bookAuthorDatas, &bookAuthorData)
	}
	data["allBooks"] = allBooks
	data["allAuthors"] = allAuthors
	data["bookAuthor"] = bookAuthor
	data["bookAuthors"] = bookAuthors
	data["bookAuthorDatas"] = bookAuthorDatas
	data["base_path"] = base_bookAuthors_path
	form.Required("book_id", "author_id")

	exists, err := m.DB.BookAuthorExists(bookAuthor.BookID, bookAuthor.AuthorID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-author relationship already exists")
		form.Errors.Add("author_id", "book-author relationship already exists")
	}
	if !form.Valid() {
		render.Template(w, r, "admin-allbookauthors.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err := m.DB.InsertBookAuthor(&bookAuthor); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Added")

	http.Redirect(w, r, "/admin/bookAuthors", http.StatusSeeOther)
}
