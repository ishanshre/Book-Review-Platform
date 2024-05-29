package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Login Handles the get method of the login
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
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
	id, access_level, is_validated, err := m.DB.Authenticate(user.Username, user.Password)
	if err != nil {
		log.Println(err)
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
	log.Println(user)
	m.UpdateSession(w, r, id, access_level, user.Username, is_validated)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) UpdateSession(w http.ResponseWriter, r *http.Request, id, access_level int, username string, is_validated bool) {
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "username", username)
	m.App.Session.Put(r.Context(), "access_level", access_level)
	m.App.Session.Put(r.Context(), "is_validated", is_validated)
	m.App.Session.Put(r.Context(), "flash", "Login Successfull")
	if !is_validated {
		m.App.Session.Put(r.Context(), "warning", "Please update your KYC to use its features!")
	}
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
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Creating a new form with form value
	form := forms.New(r.PostForm)
	register := models.User{}
	// storing the form value in user model
	register.Email = r.Form.Get("email")
	register.Username = r.Form.Get("username")
	register.Password = r.Form.Get("password")
	password2 := r.Form.Get("password2")

	// form.Required() for form  field validation
	form.Required(
		"email",
		"username",
		"password",
		"password2",
	)
	form.MinLength("username", 5)
	form.MinLength("password", 8)
	form.HasUpperCase("password")
	form.HasLowerCase("password")
	form.HasNumber("password")
	form.HasSpecialCharacter("password")
	form.IsEmail("email")
	exists, err := m.DB.UsernameExists(register.Username)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("username", "This username already exists")
	}
	exists, err = m.DB.EmailExists(register.Email)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("email", "This email already exists")
	}
	if register.Password != password2 {
		form.Errors.Add("password", "Password mismtach")
		form.Errors.Add("password2", "Password mismtach")
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
	msg := models.MailData{
		To:      register.Email,
		From:    m.App.AdminEmail,
		Subject: fmt.Sprintf("Welcome to BookWorm @%s!", register.Username),
		Content: fmt.Sprintf(`
			<h1>Welcome to BookWorm @%s!</h1>
			<p>Thank you for registering to our book review platform. Please update your kyc to access most of the features</p>
		`, register.Username),
	}
	m.App.MailChan <- msg
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
	id := m.App.Session.GetInt(r.Context(), "user_id")
	userKyc, err := m.DB.GetUserWithKyc(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	following, err := m.DB.FollowerCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	read_list_count, err := m.DB.ReadListCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	buy_list_count, err := m.DB.BuyListCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["user"] = userKyc.User
	data["kyc"] = userKyc.Kyc
	data["following"] = following
	data["read_list_count"] = read_list_count
	data["buy_list_count"] = buy_list_count
	render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) PublicUpdateKYC(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	id := m.App.Session.GetInt(r.Context(), "user_id")
	update_kyc := &models.Kyc{}
	form := forms.New(r.PostForm)

	userKyc, err := m.DB.GetUserWithKyc(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	following, err := m.DB.FollowerCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	read_list_count, err := m.DB.ReadListCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	buy_list_count, err := m.DB.BuyListCount(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	layout := "2006-01-02"
	dob, err := time.Parse(layout, r.Form.Get("date_of_birth"))
	if err != nil {
		form.Errors.Add("date_of_birth", err.Error())
	}
	update_kyc.FirstName = r.Form.Get("first_name")
	update_kyc.LastName = r.Form.Get("last_name")
	update_kyc.Gender = r.Form.Get("gender")
	update_kyc.Phone = r.Form.Get("phone")
	update_kyc.Address = r.Form.Get("address")
	update_kyc.DateOfBirth = dob
	update_kyc.DocumentType = r.Form.Get("document_type")
	update_kyc.DocumentNumber = r.Form.Get("document_number")
	update_kyc.IsValidated = false
	update_kyc.UpdatedAt = time.Now()
	update_kyc.ID = userKyc.Kyc.ID
	document_front, err := helpers.MediaPicUpload(r, "document_front", userKyc.User.Username)
	if err != nil {
		form.Errors.Add("document_front", "Document Required!")
	}
	document_back, err := helpers.MediaPicUpload(r, "document_back", userKyc.User.Username)
	if err != nil {
		form.Errors.Add("document_back", "Document Required!")
	}
	update_kyc.DocumentFront = document_front
	update_kyc.DocumentBack = document_back
	form.Required("first_name", "last_name", "gender", "phone", "address", "date_of_birth", "document_type", "document_number")
	form.MaxLength("phone", 10)
	form.MaxLength("first_name", 50)
	form.MaxLength("last_name", 50)
	form.MaxLength("address", 255)
	form.MaxLength("document_number", 50)
	data := make(map[string]interface{})
	data["base_path"] = base_users_path
	data["user"] = userKyc.User
	data["kyc"] = userKyc.Kyc
	data["following"] = following
	data["read_list_count"] = read_list_count
	data["buy_list_count"] = buy_list_count
	if !form.Valid() {
		log.Println("inside")
		render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.PublicKycUpdate(update_kyc); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "KYC Form Uploaded! Email Notification sent to admin. Please wait for admin to verify")
	msg := models.MailData{
		To:      m.App.AdminEmail,
		From:    userKyc.User.Email,
		Subject: fmt.Sprintf("Please check and validate kyc for username: %s", userKyc.User.Username),
		Content: fmt.Sprintf(`
		     <h1>Update and Validate User's Kyc: %s</h1>
			 <p>Please check, validate the updated KYC of user at <a href="%s/admin/users/detail/%s">@%s</a></p>
		`, userKyc.User.Username, r.Host, userKyc.User.Username, userKyc.User.Username),
	}
	m.App.MailChan <- msg
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (m *Repository) PostUserProfilePicUpdate(w http.ResponseWriter, r *http.Request) {
	username := m.App.Session.GetString(r.Context(), "username")
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	path, err := helpers.MediaPicUpload(r, "profile_pic", username)
	if err != nil {
		// helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Please Add a picture to upload")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	if err := m.DB.UpdateProfilePic(path, user_id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Profile Picture Updated")
	http.Redirect(w, r, "/profile", http.StatusFound)
}

// This function handles the get method of the admin login page.
// It renders the admin login page.
func (m *Repository) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(map[string]interface{})
	data["user"] = emptyUser
	render.Template(w, r, "admin_login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostAdminLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("username", "password")
	form.MinLength("username", 5)
	form.MaxLength("username", 25)
	form.MinLength("password", 8)
	user := models.User{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}
	data := make(map[string]interface{})
	data["user"] = user
	if !form.Valid() {
		render.Template(w, r, "admin_login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	id, access_level, is_validated, err := m.DB.Authenticate(user.Username, user.Password)
	if err != nil {
		form.Errors.Add("username", "Invalid username/password")
		form.Errors.Add("password", "Invalid username/password")
		render.Template(w, r, "admin_login.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if access_level != 1 {
		form.Errors.Add("username", "Invalid admin username/password")
		form.Errors.Add("password", "Invalid admin username/password")
		render.Template(w, r, "admin_login.page.tmpl", &models.TemplateData{
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
	m.UpdateSession(w, r, id, access_level, user.Username, is_validated)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
