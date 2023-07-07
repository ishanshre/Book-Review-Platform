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

// Start of handler for admin book-language relationship

// AdminAllBookLanguage retrieves all book-language relationships in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function retrieves all book-language relationships, all books, and all languages from the database.
// If any errors occur during the retrieval process, a server error is returned.
// The function prepares the necessary data and renders the "admin-allbooklanguages.page.tmpl" template,
// displaying the list of book-language relationships as well as new book language relationship add form.
func (m *Repository) AdminAllBookLanguage(w http.ResponseWriter, r *http.Request) {
	bookLanguages, err := m.DB.AllBookLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var bookLanguage models.BookLanguage
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allLanguages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookLanguageDatas := []*models.BookLanguageData{}
	for _, v := range bookLanguages {
		book, err := m.DB.GetBookTitleByID(v.BookID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		language, err := m.DB.GetLanguageByID(v.LanguageID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		bookLanguageData := &models.BookLanguageData{
			BookData:     book,
			LanguageData: language,
		}
		bookLanguageDatas = append(bookLanguageDatas, bookLanguageData)
	}
	data := make(map[string]interface{})
	data["bookLanguages"] = bookLanguages
	data["bookLanguageDatas"] = bookLanguageDatas
	data["bookLanguage"] = bookLanguage
	data["allLanguages"] = allLanguages
	data["allBooks"] = allBooks
	data["base_path"] = base_bookLanguages_path
	render.Template(w, r, "admin-allbooklanguages.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeleteBookLanguage handles the deletion of a book-language relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function extracts the book ID and language ID from the URL path parameters,
// deletes the book-language relationship from the database, and redirects the user to the "/admin/bookLanguages" page.
// If any errors occur during the process, a server error is returned.
func (m *Repository) PostAdminDeleteBookLanguage(w http.ResponseWriter, r *http.Request) {
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language_id, err := strconv.Atoi(chi.URLParam(r, "language_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteBookLanguage(book_id, language_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Language Deleted")

	http.Redirect(w, r, "/admin/bookLanguages", http.StatusSeeOther)
}

// AdminGetBookLanguageByID handes the detail logic for book language.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetBookLanguageByID(w http.ResponseWriter, r *http.Request) {

	// Retrive book id and language id from the url.
	// Parse them into integer.
	// Return a server error if any error occurs while parsing them
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language_id, err := strconv.Atoi(chi.URLParam(r, "language_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching the Book Language detail by GetBookLanguageByID interface.
	// If any error occurs, a server error is returned.
	bookLanguage, err := m.DB.GetBookLanguageByID(book_id, language_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the book title using book_id
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id

	// get the language by using language_id
	language, err := m.DB.GetLanguageByID(language_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language.ID = language_id

	// Get all books from the AllBook interface.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Get all languages from the AllLanguage interface.
	allLanguages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores book, language, all books, all languages and book-language
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["language"] = language
	data["allLanguages"] = allLanguages
	data["bookLanguage"] = bookLanguage
	data["base_path"] = base_bookLanguages_path

	// render the detail page with form and data
	render.Template(w, r, "admin-booklanguagedetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateBookLanguage handles the post method for updating the book-language relationship.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminUpdateBookLanguage(w http.ResponseWriter, r *http.Request) {

	// Fetches the book id and language id from url and parse them into integer.
	// If any error occurs, a server error is returned
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language_id, err := strconv.Atoi(chi.URLParam(r, "language_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Parse the form and populate the PostForm
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form with PostForm
	form := forms.New(r.PostForm)

	// Get the updated book id and language id from the post form.
	// If any error occurs, a server error is returned
	updated_book_id, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	updated_language_id, err := strconv.Atoi(r.Form.Get("language_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Populate new BookLanguage instance with update book id and language id.
	bookLanguage := models.BookLanguage{
		BookID:     updated_book_id,
		LanguageID: updated_language_id,
	}

	// Check for existing relationship between book and author.
	// A server error is retrned if any error occurs
	exists, err := m.DB.BookLanguageExists(bookLanguage.BookID, bookLanguage.LanguageID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If exists then add error with message telling the relationship exists
	if exists {
		form.Errors.Add("book_id", "book-author relationship already exists")
		form.Errors.Add("language_id", "book-author relationship already exists")
	}

	// get book title with book_id
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id

	// get the language using langugage id
	language, err := m.DB.GetLanguageByID(language_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language.ID = language_id

	// Get all books
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get all languages
	allLanguages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores book, language, all books, all language
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["language"] = language
	data["allLanguages"] = allLanguages
	data["bookLanguage"] = bookLanguage
	data["base_path"] = base_bookLanguages_path

	// Add required form validation for language id and book id
	form.Required("book_id", "language_id")
	if !form.Valid() {
		render.Template(w, r, "admin-booklanguagedetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Update the book language relationship using UpdateBookLanguage interface.
	// Returns a server error if any error occurs.
	if err := m.DB.UpdateBookLanguage(&bookLanguage, book_id, language_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Language Updated")

	// Redirect to book language detail page if update successfull.
	http.Redirect(w, r, fmt.Sprintf("/admin/bookLanguages/detail/%d/%d", bookLanguage.BookID, bookLanguage.LanguageID), http.StatusSeeOther)
}

// PostAdminInsertBookLanguage handles the insertion of a new book-language relationship in an admin context.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request and validates it. If any parsing or validation errors occur,
// a server error is returned. The function retrieves the book ID and language ID from the form data and creates a new
// BookLanguage object with the provided IDs. It then retrieves all book-Language relationships, all books, and all Languages
// from the database to prepare the necessary data for rendering the template. The function checks if the book-language
// relationship already exists and adds an error to the form if it does. If the form is not valid, the template is
// rendered with the form and data. If the form is valid, the new book-Language relationship is inserted into the database
// and the user is redirected to the "/admin/booklanguages" page.
func (m *Repository) PostAdminInsertBookLanguage(w http.ResponseWriter, r *http.Request) {

	// Parse the form. Returns server error if unable to parse the form
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a new form using the post form
	form := forms.New(r.PostForm)

	// create a data map that stores the values to pass to template
	data := make(map[string]interface{})

	book_id, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language_id, err := strconv.Atoi(r.Form.Get("language_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	bookLanguage := models.BookLanguage{
		BookID:     book_id,
		LanguageID: language_id,
	}

	bookLanguages, err := m.DB.AllBookLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allLanguages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookLanguageDatas := []*models.BookLanguageData{}
	for _, v := range bookLanguages {
		book, err := m.DB.GetBookTitleByID(v.BookID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		language, err := m.DB.GetLanguageByID(v.LanguageID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		bookLanguageData := &models.BookLanguageData{
			BookData:     book,
			LanguageData: language,
		}
		bookLanguageDatas = append(bookLanguageDatas, bookLanguageData)
	}

	data["allBooks"] = allBooks
	data["allLanguages"] = allLanguages
	data["bookLanguage"] = bookLanguage
	data["bookLanguages"] = bookLanguages
	data["bookLanguageDatas"] = bookLanguageDatas
	data["base_path"] = base_bookLanguages_path

	form.Required("book_id", "language_id")

	exists, err := m.DB.BookLanguageExists(bookLanguage.BookID, bookLanguage.LanguageID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-language relationship already exists")
		form.Errors.Add("language_id", "book-language relationship already exists")
	}

	if !form.Valid() {
		render.Template(w, r, "admin-allbooklanguages.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err := m.DB.InsertBookLanguage(&bookLanguage); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Language Relationship Added")
	http.Redirect(w, r, "/admin/bookLanguages", http.StatusSeeOther)
}
