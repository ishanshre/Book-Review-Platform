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

// AdminAllFollowers fetches all the relation record between user and books in Followers
func (m *Repository) AdminAllFollowers(w http.ResponseWriter, r *http.Request) {
	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allUsers, err := m.DB.AllUsers(100000000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["allUsers"] = allUsers
	data["allAuthors"] = allAuthors
	data["base_path"] = base_followers_path
	render.Template(w, r, "admin-allfollowers.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminAllFollowerApi(w http.ResponseWriter, r *http.Request) {
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
	filterFollowers, err := m.DB.FollowerFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filterFollowers)
}

// PostAdminInsertFollower handles post method logic for user following author by admin.
// It takes HTTP response writer and request as paramaters
func (m *Repository) PostAdminInsertFollower(w http.ResponseWriter, r *http.Request) {

	// Parse the form. Returns server error if unable to parse the form
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a new form using the post form
	form := forms.New(r.PostForm)

	// create a data map that stores the values to pass to template
	data := make(map[string]interface{})

	author_id, err := strconv.Atoi(r.Form.Get("author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	follower := models.Follower{
		UserID:     user_id,
		AuthorID:   author_id,
		FollowedAt: time.Now(),
	}

	allAuthors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["allAuthors"] = allAuthors
	data["allUsers"] = allUsers
	data["follower"] = follower
	data["base_path"] = base_followers_path
	form.Required("author_id", "user_id")

	exists, err := m.DB.FollowerExists(&follower)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("author", "book-author relationship already exists")
		form.Errors.Add("user_id", "book-author relationship already exists")
	}

	if !form.Valid() {
		render.Template(w, r, "admin-allfollowers.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if err := m.DB.InsertFollower(&follower); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Follower record added")

	http.Redirect(w, r, "/admin/followers", http.StatusSeeOther)
}

// PostAdminDeleteFollow Handles the post method for deleting follower list record.
// It takes HTTP response writer and request as paramters
func (m *Repository) PostAdminDeleteFollow(w http.ResponseWriter, r *http.Request) {

	// Parsing the book id and user id from the url.
	// If any error occurs, a server error is returned.
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// DeleteBuyList interface is used to deleting the record.
	if err := m.DB.DeleteFollower(user_id, author_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Follower record deleted")

	http.Redirect(w, r, "/admin/followers", http.StatusSeeOther)
}

// AdminGetFollowerByID handes the detail logic for Follower table.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetFollowerByID(w http.ResponseWriter, r *http.Request) {

	// Retrive user id and book id from the url.
	// Parse them into integer.
	// Return a server error if any error occurs while parsing them
	user_id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching the buy list detail by GetFolowerByID interface.
	// If any error occurs, a server error is returned.
	follower, err := m.DB.GetFollowerByID(user_id, author_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the book title using book_id
	author, err := m.DB.GetAuthorByID(author_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author.ID = author_id

	// get the user by using user_id
	user, err := m.DB.GetUserByID(user_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = user_id

	// Get all authors from the AllAuthor interface.
	allAuthors, err := m.DB.AllAuthor()
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

	// create a data map that stores author, user, all authors, all users and followers
	data := make(map[string]interface{})
	data["author"] = author
	data["allAuthors"] = allAuthors
	data["user"] = user
	data["allUsers"] = allUsers
	data["follower"] = follower
	data["base_path"] = base_followers_path

	// render the detail page with form and data
	render.Template(w, r, "admin-followersdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateFollower handles the post method for updating the author-user follower relationship.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminUpdateFollower(w http.ResponseWriter, r *http.Request) {

	// Fetches the author id and user id from url and parse them into integer.
	// If any error occurs, a server error is returned
	author_id, err := strconv.Atoi(chi.URLParam(r, "author_id"))
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
	updated_author_id, err := strconv.Atoi(r.Form.Get("author_id"))
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
	follower := models.Follower{
		UserID:   updated_user_id,
		AuthorID: updated_author_id,
	}

	// Check for existing relationship between book and user in read list.
	// A server error is retrned if any error occurs
	exists, err := m.DB.FollowerExists(&follower)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If exists then add error with message telling the relationship exists
	if exists {
		form.Errors.Add("author_id", "book-user relationship already exists")
		form.Errors.Add("user_id", "book-user relationship already exists")
	}

	// get author detail
	author, err := m.DB.GetAuthorByID(author_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author.ID = author_id

	// get the user using langugage id
	user, err := m.DB.GetUserByID(user_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = user_id

	// Get all authors
	allAuthors, err := m.DB.AllAuthor()
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

	// create a data map that stores Author, user, all Authors, all users
	data := make(map[string]interface{})
	data["author"] = author
	data["allAuthors"] = allAuthors
	data["user"] = user
	data["allUsers"] = allUsers
	data["follower"] = follower
	data["base_path"] = base_followers_path
	// Add required form validation for language id and book id
	form.Required("author_id", "user_id")
	if !form.Valid() {
		render.Template(w, r, "admin-followersdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Update the book language relationship using UpdateBookLanguage interface.
	// Returns a server error if any error occurs.
	if err := m.DB.UpdateFollower(&follower, user_id, author_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Follower record updated")

	// Redirect to book language detail page if update successfull.
	http.Redirect(w, r, fmt.Sprintf("/admin/followers/detail/%d/%d", follower.AuthorID, follower.UserID), http.StatusSeeOther)
}
