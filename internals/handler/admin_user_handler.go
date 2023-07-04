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

// AdminDashboard renders admin page for admin user only
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminAllUsers is a handler that handles the HTTP request for retrieving all users in the admin panel.
// It retrieves the page and limit parameters from the URL query, calculates the offset based on the page and limit,
// retrieves the users from the database with the specified limit and offset, creates a data map containing the users
// for rendering the template, and renders the admin all users page.
func (m *Repository) AdminAllUsers(w http.ResponseWriter, r *http.Request) {

	// set the default limit and offset values
	limit := 10
	offset := 0
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	p := true
	if err != nil {
		p = false
	}
	if p {
		offset = (page - 1) * limit
	}
	filter, err := strconv.Atoi(r.URL.Query().Get("limit"))
	p = true
	if err != nil {
		p = false
	}
	if p {
		limit = filter
	}
	users, err := m.DB.AllUsers(limit, offset)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["base_path"] = base_users_path
	data["users"] = users
	render.Template(w, r, "admin-allusers.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminGetUserDetailByID is a handler that handles the HTTP request for retrieving a user's detail by ID in the admin panel.
// It retrieves the user ID from the URL parameters, retrieves the user's information from the database, creates a data map
// containing the user instance for rendering the template, and renders the user detail page.
func (m *Repository) AdminGetUserDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
	}
	user.ID = id
	data := make(map[string]interface{})
	data["user"] = user
	data["base_path"] = base_users_path
	render.Template(w, r, "admin-userdetail.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// AdminUpdateUser is a handler that handles the HTTP request for updating a user's information in the admin panel.
// It retrieves the user ID from the URL parameters, parses the form data from the request, creates a new user instance
// with the updated information, validates the form, updates the user's information in the database, and redirects
// the user to the user detail page.
func (m *Repository) AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	user := &models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Gender:    r.Form.Get("gender"),
		Phone:     r.Form.Get("phone"),
		Address:   r.Form.Get("address"),
	}
	user.ID = id
	data := make(map[string]interface{})
	data["base_path"] = base_users_path
	data["user"] = user
	if !form.Valid() {
		render.Template(w, r, "admin-userdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateUser(user); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "User Updated")
	http.Redirect(w, r, fmt.Sprintf("/admin/users/detail/%d", id), http.StatusSeeOther)
}

// PostAdminUserProfileUpdate is a handler that handles the HTTP request for updating a user's profile picture in the admin panel.
// It retrieves the user ID from the URL parameters, gets the username from the session, uploads the profile picture file,
// updates the profile picture path in the database, and redirects the user to the user detail page.
func (m *Repository) PostAdminUserProfileUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	username := m.App.Session.Get(r.Context(), "username")
	path, err := helpers.UserRegitserFileUpload(r, "profile_pic", username.(string))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.UpdateProfilePic(path, id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Profile Picture Updated")
	http.Redirect(w, r, fmt.Sprintf("/admin/users/detail/%d", id), http.StatusSeeOther)
}

// PostAdminUserDeleteByID renders a confim page to delete users.
// It takes HTTP response writer and request as parameters.
// It parse the id from url, check if auth user id and id from url mathches or not, delete the user if not match.
func (m *Repository) PostAdminUserDeleteByID(w http.ResponseWriter, r *http.Request) {
	// Parse the id from url into integer.
	// If any error occurs, a server error is returned.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Retrive the auth user id from the session.
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	// if user_id and id from url matches then return a client error.
	// Admin himself cannot delete himself
	if id == user_id {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	// The function calls DeleteUser interface to delete the user form the database
	if err := m.DB.DeleteUser(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "User Deleted")
	// Redirect the admin to all users page
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminUserAdd renders page for adding user by admin.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(map[string]interface{})
	data["register"] = emptyUser
	data["base_path"] = base_users_path
	render.Template(w, r, "admin-usercreate.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUserAdd handles post method for creating user from admin interface.
// It takes HTTP response writer and request as parameters.
// It parses form, store the form data, add validations, check for existing user and then only add user if not exists.
func (m *Repository) PostAdminUserAdd(w http.ResponseWriter, r *http.Request) {
	// Parse the form from the request.
	// If error occurs, a server error is returned
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form with post form
	form := forms.New(r.PostForm)

	// Add form field validations
	form.Required("username", "email", "password", "citizenship_number")
	form.MinLength("username", 5)
	form.MinLength("password", 8)
	form.HasLowerCase("password")
	form.HasUpperCase("password")
	form.HasNumber("password")
	form.HasSpecialCharacter("password")

	// Store the data from post form to register_user.
	register_user := models.User{
		Username:          r.Form.Get("username"),
		Email:             r.Form.Get("email"),
		Password:          r.Form.Get("password"),
		CitizenshipNumber: r.Form.Get("citizenship_number"),
	}

	// UsernameExists interface is called to check if username already exists.
	exists, err := m.DB.UsernameExists(register_user.Username)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// if exists add error to form
	if exists {
		form.Errors.Add("username", "Username already exists")
	}
	exists, err = m.DB.EmailExists(register_user.Email)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// if exists add error to form
	if exists {
		form.Errors.Add("email", "Email already exists")
	}

	// create a data map that holds register_user.
	data := make(map[string]interface{})
	data["register"] = register_user
	data["base_path"] = base_users_path
	// If form is invalid render "admin-usercreate.page.tmpl" with form and data.
	if !form.Valid() {
		render.Template(w, r, "admin-usercreate.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// create a hash for new passowrd and store in register_user
	hashed_password, err := helpers.EncryptPassword(register_user.Password)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	register_user.Password = hashed_password

	// Call AdminInsertUser interface for inserting new user.
	// If any error occurs, a server error is returned.
	if err := m.DB.AdminInsertUser(&register_user); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "User Added")

	// Redirect the admin to all users page.
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
