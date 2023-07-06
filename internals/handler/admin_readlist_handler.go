package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// AdminAllReadList fetches all the relation record between user and books in readLists
func (m *Repository) AdminAllReadList(w http.ResponseWriter, r *http.Request) {
	readLists, err := m.DB.AllReadList()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var readList models.ReadList
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allUsers, err := m.DB.AllUsers(1000000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	readListDatas := []*models.ReadListData{}
	for _, v := range readLists {
		book, err := m.DB.GetBookTitleByID(v.BookID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		user, err := m.DB.GetUserByID(v.UserID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		user.ID = v.UserID
		readListData := &models.ReadListData{
			BookData:  book,
			UserData:  user,
			CreatedAt: v.CreatedAt,
		}
		readListDatas = append(readListDatas, readListData)
	}
	data := make(map[string]interface{})
	data["readLists"] = readLists
	data["readListDatas"] = readListDatas
	data["readList"] = readList
	data["allUsers"] = allUsers
	data["allBooks"] = allBooks
	data["base_path"] = base_readLists_path

	render.Template(w, r, "admin-allreadlists.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})

}

// PostAdminInsertReadList handles post method logic for adding books to read list by admin.
// It takes HTTP response writer and request as paramaters
func (m *Repository) PostAdminInsertReadList(w http.ResponseWriter, r *http.Request) {

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
	user_id, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	readList := models.ReadList{
		UserID:    user_id,
		BookID:    book_id,
		CreatedAt: time.Now(),
	}

	readLists, err := m.DB.AllReadList()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["allBooks"] = allBooks
	data["allUsers"] = allUsers
	data["readList"] = readList
	data["readLists"] = readLists
	data["base_path"] = base_readLists_path
	form.Required("book_id", "user_id")

	exists, err := m.DB.ReadListExists(readList.UserID, readList.BookID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("book_id", "book-language relationship already exists")
		form.Errors.Add("user_id", "book-language relationship already exists")
	}

	if !form.Valid() {
		render.Template(w, r, "admin-allreadlists.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err := m.DB.InsertReadList(&readList); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Read List Added")

	http.Redirect(w, r, "/admin/readLists", http.StatusSeeOther)
}

// AdminGetReadListByID handes the detail logic for readList table.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetReadListByID(w http.ResponseWriter, r *http.Request) {

	// Retrive user id and book id from the url.
	// Parse them into integer.
	// Return a server error if any error occurs while parsing them
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching the Book Language detail by GetReadListByID interface.
	// If any error occurs, a server error is returned.
	readList, err := m.DB.GetReadListByID(user_id, book_id)
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

	// get the user by using user_id
	user, err := m.DB.GetUserByID(user_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = user_id

	// Get all books from the AllBook interface.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Get all user from the AllUsers interface.
	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores book, user, all books, all users and read list
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["user"] = user
	data["allUsers"] = allUsers
	data["readList"] = readList
	data["base_path"] = base_readLists_path

	// render the detail page with form and data
	render.Template(w, r, "admin-readlistsdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminDeleteReadList Handles the post method for deleting read list record.
// It takes HTTP response writer and request as paramters
func (m *Repository) PostAdminDeleteReadList(w http.ResponseWriter, r *http.Request) {

	// Parsing the book id and user id from the url.
	// If any error occurs, a server error is returned.
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// DeleteReadList interface is used to deleting the record.
	if err := m.DB.DeleteReadList(user_id, book_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Read List Record Deleted")

	http.Redirect(w, r, "/admin/readLists", http.StatusSeeOther)
}

// PostAdminUpdateReadList handles the post method for updating the book-user add read list relationship.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminUpdateReadList(w http.ResponseWriter, r *http.Request) {

	// Fetches the book id and user id from url and parse them into integer.
	// If any error occurs, a server error is returned
	book_id, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
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
	updated_user_id, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Populate new BookLanguage instance with update book id and language id.
	readList := models.ReadList{
		UserID: updated_user_id,
		BookID: updated_book_id,
	}

	// Check for existing relationship between book and user in read list.
	// A server error is retrned if any error occurs
	exists, err := m.DB.ReadListExists(readList.UserID, readList.BookID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If exists then add error with message telling the relationship exists
	if exists {
		form.Errors.Add("book_id", "book-user relationship already exists")
		form.Errors.Add("user_id", "book-user relationship already exists")
	}

	// get book title with book_id
	book, err := m.DB.GetBookTitleByID(book_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book.ID = book_id

	// get the user using langugage id
	user, err := m.DB.GetUserByID(user_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = user_id

	// Get all books
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get all languages
	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores book, language, all books, all language
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["user"] = user
	data["allUsers"] = allUsers
	data["readList"] = readList
	data["base_path"] = base_readLists_path

	// Add required form validation for language id and book id
	form.Required("book_id", "user_id")
	if !form.Valid() {
		render.Template(w, r, "admin-readlistsdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Update the book language relationship using UpdateBookLanguage interface.
	// Returns a server error if any error occurs.
	if err := m.DB.UpdateReadList(&readList, book_id, user_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Read List record updated")

	// Redirect to book language detail page if update successfull.
	http.Redirect(w, r, fmt.Sprintf("/admin/readLists/detail/%d/%d", readList.BookID, readList.UserID), http.StatusSeeOther)
}
