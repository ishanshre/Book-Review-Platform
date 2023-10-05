package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// AdminAllBook handles logic for reteriving all Books in admin page.
// It takes HTTP response writer and request as parameters.
// The function calls AllBook interface to reterive all the books from the database
// If any error occurs, a server error is returned.
// A data map is created that stores the books
// Finally, "admin-allbooks.page.tmpl" go template is rendered with data.
func (m *Repository) AdminAllBook(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["base_path"] = base_books_path
	render.Template(w, r, "admin-allbooks.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllBookApi(w http.ResponseWriter, r *http.Request) {
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
	filteredBooks, err := m.DB.AllBooksFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBooks)
}

// PostAdminDeleteBook handles the delete logic of book for admin.
// It takes HTTP response writer and response as parameters.
// It fetches the book id from the url and parse it into integer.
// If it fails to parse, a server error is returned.
// The function calls the DeleteBook method and pass book id as parameter.
// If any error occurs, then a server error is returned.
// If successfull, then admin is redirected to all books page.
func (m *Repository) PostAdminDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}
	if err := m.DB.DeleteBook(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Deleted")

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

// AdminGetBookDetailByID retrieves and displays the details of a book in the admin.
// It takes the HTTP response writer and request as parameters.
// The function extracts the book ID from the URL parameter and converts it to an integer.
// If there is an error during the conversion, a server error is returned.
// The function retrieves the book from the database using the GetBookByID method.
// If an error occurs during the retrieval, a server error is returned.
// The function retrieves all publishers from the database using the AllPublishers method.
// If an error occurs during the retrieval, a server error is returned.
// The function retrieves the publisher of the book using the GetPublisherByID method.
// If an error occurs during the retrieval, a server error is returned.
// A data map is created to store the book, publishers, and publisher.
// The "admin-bookdetail.page.tmpl" template is rendered, passing a new form instance and the data map.
// The function returns after rendering the template.
func (m *Repository) AdminGetBookDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}
	book, err := m.DB.GetBookByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publisher, err := m.DB.GetPublisherByID(book.PublisherID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["publishers"] = publishers
	data["publisher"] = publisher
	data["base_path"] = base_books_path

	render.Template(w, r, "admin-bookdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// AdminInsertBook Handles the get method of adding books to the database.
// It takes the HTTP response writer and request as a parameters.
// It renders the add book form.
// It fetches all the publisers from the database by calling AllPublishers interface.
// If any error occurs during the retrieval, a server error is returned.
// A data map is created to store the book and publishers.
// The "admin-bookinsert.page.tmpl" go template is rendered, passing a new form and data.
func (m *Repository) AdminInsertBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["publishers"] = publishers
	data["base_path"] = base_books_path

	render.Template(w, r, "admin-bookinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertBook handles the insertion of a new book in the admin context.
// It takes the HTTP response writer and request as parameters. The function first parses the form data from the request.
// If an error occurs during the parsing, a server error is returned.
// A new form instance is created to handle form validation, and a data map is initialized to store the retrieved data.
// The function parses the published date, paperback, publisher ID, ISBN, and isActive from the form fields,
// ensuring their appropriate data types. If any errors occur during the parsing, appropriate form errors are added.
// The function constructs a new book instance of the models.Book struct with the values from the form fields,
// including the title, description, ISBN, published date, paperback, isActive, addedAt, updatedAt, and publisherID.
// The function then attempts to upload the book cover image using the AdminPublicUploadImage2 helper function,
// passing the request, "cover" form field, "book" as the upload type, and the book's ISBN as the filename.
// If an error occurs during the image upload, a form error is added.
// The book cover filename is assigned to the book instance.
// The required and length validations are performed on the form fields using the form instance.
// If the form is not valid, the function renders the "admin-bookinsert.page.tmpl" template,
// passing the form and data map. The function then returns.
// If there are no form validation errors, the function calls the InsertBook method of the database with the book instance.
// If an error occurs during the insertion, a server error is returned.
// Finally, the function redirects the user to the list of all books in the admin context.
func (m *Repository) PostAdminInsertBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	publishedDate, err := time.Parse(time.DateOnly, r.Form.Get("published_date"))
	if err != nil {
		form.Errors.Add("published_date", "Enter the valid date")
	}
	paperback, err := strconv.Atoi(r.Form.Get("paperback"))
	if err != nil {
		form.Errors.Add("paperback", "Paperback must be an integer")
	}
	publisherID, err := strconv.Atoi(r.Form.Get("publisher_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isbn, err := strconv.ParseInt(r.Form.Get("isbn"), 10, 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isActive, err := strconv.ParseBool(r.Form.Get("is_active"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book := models.Book{
		Title:         r.Form.Get("title"),
		Description:   r.Form.Get("description"),
		Isbn:          isbn,
		PublishedDate: publishedDate,
		Paperback:     paperback,
		IsActive:      isActive,
		AddedAt:       time.Now(),
		UpdatedAt:     time.Now(),
		PublisherID:   publisherID,
	}
	cover, err := helpers.AdminPublicUploadImage2(r, "cover", "book", book.Isbn)
	if err != nil {
		form.Errors.Add("cover", "No image uploaded")
	}
	book.Cover = cover
	form.Required("isbn", "title")
	form.MinLength("isbn", 13)
	form.MaxLength("isbn", 13)
	form.MaxLength("title", 100)
	form.MaxLength("description", 10000)
	data["book"] = book
	data["base_path"] = base_books_path
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["publishers"] = publishers
	if !form.Valid() {
		render.Template(w, r, "admin-bookinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertBook(&book); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Added")

	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

// PostAdminUpdateBook handles the update of a book in the admin context.
// It takes the HTTP response writer and request as parameters. The function retrieves the book ID from the URL parameter
// and parses it into an integer. If an error occurs during the parsing, a server error is returned.
// The function then parses the form data from the request. If an error occurs during the parsing, a server error is returned.
// A new form instance is created to handle form validation, and a data map is initialized to store the retrieved data.
// The function retrieves the book from the database using the GetBookByID method based on the parsed book ID.
// If an error occurs during the retrieval, a server error is returned.
// The function also retrieves the publisher of the book using the GetPublisherByID method based on the publisher ID
// stored in the retrieved book. If an error occurs during the retrieval, a server error is returned.
// Additionally, the function retrieves all publishers from the database using the AllPublishers method.
// If an error occurs during the retrieval, a server error is returned.
// The retrieved book, publisher, and publishers are stored in the data map.
// The function parses other form fields such as ISBN, published date, paperback, and isActive, ensuring their appropriate data types.
// If any errors occur during the parsing, a server error is returned.
// The function constructs an updated_book instance of the models.Book struct with the updated values from the form fields.
// The required and length validations are performed on the form fields using the form instance.
// If the form is not valid, the function renders the "admin-bookdetail.page.tmpl" template, passing the form and data map.
// The function then returns.
// If there are no form validation errors, the function calls the UpdateBook method of the database with the updated_book instance.
// If an error occurs during the update, a server error is returned.
// Finally, the function redirects the user to the detailed view of the updated book using the book ID.
func (m *Repository) PostAdminUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	book, err := m.DB.GetBookByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publisher, err := m.DB.GetPublisherByID(book.PublisherID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["book"] = book
	data["publisher"] = publisher
	data["publishers"] = publishers
	isbn, err := strconv.ParseInt(r.Form.Get("isbn"), 10, 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publishedDate, err := time.Parse(time.DateOnly, r.Form.Get("published_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	paperback, err := strconv.Atoi(r.Form.Get("paperback"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isActive, err := strconv.ParseBool(r.Form.Get("is_active"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publishedBy, err := strconv.Atoi(r.Form.Get("publisher_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	updated_book := models.Book{
		ID:            book.ID,
		Title:         r.Form.Get("title"),
		Description:   r.Form.Get("description"),
		Isbn:          isbn,
		PublishedDate: publishedDate,
		Paperback:     paperback,
		IsActive:      isActive,
		PublisherID:   publishedBy,
		UpdatedAt:     time.Now(),
	}
	form.Required("title", "isbn")
	form.MinLength("isbn", 13)
	form.MaxLength("isbn", 13)

	form.MaxLength("title", 100)
	form.MaxLength("description", 10000)
	data["base_path"] = base_books_path

	if !form.Valid() {
		render.Template(w, r, "admin-bookdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateBook(&updated_book); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Book Added")

	http.Redirect(w, r, fmt.Sprintf("/admin/books/detail/%d", book.ID), http.StatusSeeOther)

}
