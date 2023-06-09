package handler

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Login Handles the get method of the login
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated or not.
	// If authenticated then redirects to home page
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var emptyLogin models.User
	data := make(map[string]interface{})
	data["user"] = emptyLogin

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostLogin is a handler that handles the HTTP POST request for user login.
// It validates the login form data, authenticates the user, and stores the user ID and access level in the session.
// If the authentication is successful, it redirects the user to the home page.
// If the authentication fails, it renders the login page with appropriate error messages.
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated or not.
	// If authenticated then redirects to home page
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Renew session token for user login
	_ = m.App.Session.RenewToken(r.Context())

	//Parse the PostForm and creates r.PostForm or r.Form
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// add post form data into new form
	form := forms.New(r.PostForm)
	// form.Has("username")
	// form.Has("password")

	// Validate Requried fields
	form.Required("username", "password")
	form.MinLength("username", 5)

	// Store the form data in models
	user := models.User{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}

	data := make(map[string]interface{})
	data["user"] = user
	// Check if the form is valid.
	// If valid renders the form with previous data
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	id, access_level, err := m.DB.Authenticate(user.Username, user.Password)
	if err != nil {
		form.Errors.Add("username", "Invalid username/password")
		form.Errors.Add("password", "Invalid username/password")
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateLastLogin(id); err != nil {
		helpers.ServerError(w, err)
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "username", user.Username)
	m.App.Session.Put(r.Context(), "access_level", access_level)
	m.App.Session.Put(r.Context(), "flash", "Login Successfull")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Register handles the get method of the register.
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	var emptyRegister models.User
	var data = make(map[string]interface{})
	data["register"] = emptyRegister
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostRegister is a handler that handles the HTTP POST request for user registration.
// It creates a new session token, parses a multipart form, validates the form fields,
// and inserts the user data into the database if the form is valid.
// If there are any errors during the registration process, the user is redirected to the registration page with error messages.
// If the registration is successful, the user is redirected to the login page.
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {

	// create a new token
	_ = m.App.Session.RenewToken(r.Context())

	// Initially parse a multipart form to make use of form
	if err := r.ParseMultipartForm(5 << 1); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Creating a new form with form value
	form := forms.New(r.MultipartForm.Value)

	// storing the form value in user model
	register := models.User{
		FirstName:         r.Form.Get("first_name"),
		LastName:          r.Form.Get("last_name"),
		Email:             r.Form.Get("email"),
		Username:          r.Form.Get("username"),
		Password:          r.Form.Get("password"),
		Gender:            r.Form.Get("gender"),
		CitizenshipNumber: r.Form.Get("citizenship_number"),
	}

	// form.Required() for form  field validation
	form.Required(
		"first_name",
		"last_name",
		"email",
		"username",
		"password",
		"gender",
		"citizenship_number",
	)
	form.MinLength("username", 8)
	form.MinLength("password", 8)
	form.HasUpperCase("password")
	form.HasLowerCase("password")
	form.HasNumber("password", "username")
	form.HasSpecialCharacter("password")
	exists, err := m.DB.UsernameExists(register.Username)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("username", "This username already exists")
	}
	if !form.Valid() {
		data := make(map[string]interface{})
		data["register"] = register
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Upload front part of citizenship document
	citizenship_front, err := helpers.UserRegitserFileUpload(r, "citizenship_front", register.Username)
	if err != nil {
		form.Errors.Add("citizenship_front", err.Error())
	}

	// Upload back part of citizenship document
	citizenship_back, err := helpers.UserRegitserFileUpload(r, "citizenship_back", register.Username)
	if err != nil {
		form.Errors.Add("citizenship_back", err.Error())
	}

	// storing the uploaded file path in user model
	register.CitizenshipFront = citizenship_front
	register.CitizenshipBack = citizenship_back
	log.Println(register.CitizenshipBack)
	log.Println(register.CitizenshipFront)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["register"] = register
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	hashed_password, err := helpers.EncryptPassword(register.Password)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	register.Password = hashed_password
	if err := m.DB.InsertUser(&register); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "User Registration Successfull")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// Logout is a handler that handles the HTTP request for logging out the user.
// If the user is not authenticated, they are redirected to the login page.
// If the user is authenticated, their session is destroyed and a new session token is generated.
// Finally, the user is redirected to the login page.
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	// cannot accesss logout url if user is not authenticated
	// unauthenticated users are redirected to login page
	if !helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	m.App.Session.Put(r.Context(), "flash", "Logout Successfull")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// PersonalProfile is a handler that handles the HTTP request for displaying the personal profile of a user.
// It retrieves the user ID from the session, fetches the user's profile information from the database,
// opens and reads the user's citizenship front and back images, encodes them as base64 strings,
// creates a data map containing the encoded images and user profile information for rendering the template,
// and renders the personal profile page.
func (m *Repository) PersonalProfile(w http.ResponseWriter, r *http.Request) {
	id := m.App.Session.Get(r.Context(), "user_id")
	user, err := m.DB.GetProfilePersonal(id.(int))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	cit_front, err := os.Open(user.CitizenshipFront)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer cit_front.Close()
	cit_back, err := os.Open(user.CitizenshipBack)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer cit_back.Close()
	frontStat, err := cit_front.Stat()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	frontSize := frontStat.Size()
	frontData := make([]byte, frontSize)
	_, err = cit_front.Read(frontData)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	backStat, err := cit_back.Stat()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	backSize := backStat.Size()
	backData := make([]byte, backSize)
	_, err = cit_back.Read(backData)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	imgBase64Front := base64.StdEncoding.EncodeToString(frontData)
	imgBase64Back := base64.StdEncoding.EncodeToString(backData)
	data := make(map[string]interface{})
	data["citizenship_front"] = imgBase64Front
	data["citizenship_back"] = imgBase64Back
	data["user_profile"] = user
	render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
