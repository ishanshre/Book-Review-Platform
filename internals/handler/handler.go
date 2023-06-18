package handler

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
	"github.com/ishanshre/Book-Review-Platform/internals/repository"
	"github.com/ishanshre/Book-Review-Platform/internals/repository/dbrepo"
)

// Repository used to get global app config and database access
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo is of type Repository and used by handlers to get access to global app config and datbase
var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Assign Repository to Repo for handler to access
func NewHandler(r *Repository) {
	Repo = r
}

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

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

	// Redirect the admin to all users page
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminUserAdd renders page for adding user by admin.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(map[string]interface{})
	data["register"] = emptyUser
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
	form.Required("username", "email", "password")
	form.MinLength("username", 5)
	form.MinLength("password", 8)
	form.HasLowerCase("password")
	form.HasUpperCase("password")
	form.HasNumber("password")
	form.HasSpecialCharacter("password")

	// Store the data from post form to register_user.
	register_user := models.User{
		Username: r.Form.Get("username"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
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

	// Redirect the admin to all users page.
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// ResetPassword render the password reset page.
// It takes HTTP response writer and request as parameters.
// It renders password reset page.
func (m *Repository) ResetPassword(w http.ResponseWriter, r *http.Request) {

	// create a data map that holds empty User model
	var emptyUser models.User
	data := make(map[string]interface{})
	data["reset_user"] = emptyUser

	// Render the "reset-password.page.tmpl" with form and data.
	render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// userStore the token and user data in cache
var userStore = &models.UserTokenStore{
	Users:             make(map[string]models.User),
	PasswordResetRepo: make(map[string]models.PasswordResetToken),
}

// PostResetPassword handles the post method that takes in the email address and send the reset token to user email.
// It takes HTTP response writer and request as parameters.
// It parse the form, validates the email, checks if email exists, then send reset token to email if exists.
func (m *Repository) PostResetPassword(w http.ResponseWriter, r *http.Request) {

	// Parse the form data from the request.
	// If any error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form with post form values.
	form := forms.New(r.PostForm)

	// Add form filed validation/
	form.Required("email")

	// Create a variable that stores the USER model/
	reset_user := models.User{
		Email: r.Form.Get("email"),
	}

	// Create a data map that holds the reset_user User model/
	data := make(map[string]interface{})
	data["reset_user"] = reset_user

	// Check if email exists.
	// If error occurs, a server error is returned.
	exists, err := m.DB.EmailExists(reset_user.Email)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If exists then add error to the form.
	if exists {
		form.Errors.Add("email", "This email does not exist")
	}

	// If form is not valid then render the "reset-password.page.tmpl" go template with form and data.
	if !form.Valid() {
		render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// store the user email in cache.
	_ = userStore.Users[reset_user.Email]

	// create a token.
	token, err := helpers.GenerateRandomToken(32)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// Create a PasswordResetToken model and store token, email and expiry date
	resetToken := models.PasswordResetToken{
		Token:     token,
		Email:     reset_user.Email,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	userStore.PasswordResetRepo[token] = resetToken

	// send the email to email address with reset token
	body := fmt.Sprintf(`
		<h1>Reset Password</h1>
			<strong>The token for password change = </strong> %s<br><hr>
			<button><a href="%s/user/reset">Reset</a></button><br><hr>
			<strong>Ignore it, if you did not apply for reset password</strong>
	`, token, r.Header["Origin"][0])
	msg := models.MailData{
		To:      reset_user.Email,
		From:    "admin@bookworm.com",
		Subject: "Change Password",
		Content: body,
	}
	m.App.MailChan <- msg

	// Redirect user to password reset page.
	http.Redirect(w, r, "/user/reset", http.StatusSeeOther)
}

// ResetPasswordChange renders password change form with new and confirm password.
// It takes HTTP response writer and request as parameters.
func (m *Repository) ResetPasswordChange(w http.ResponseWriter, r *http.Request) {

	// Create a empty ResetPassword model
	var emptyPass models.ResetPassword

	// Create data map that holds emptyPass
	data := make(map[string]interface{})
	data["reset_password"] = emptyPass

	// render the "reset-password-change.page.tmpl" go template
	render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostResetPasswordChange handles post method and changes the password.
// It takes HTTP response writer and request as parameters.
// It parses the form data from the request, validates the form, checks the validity of the reset token, encrypts the new password,
// updates the user's password in the database, sends a notification email, and redirects the user to the login page.
func (m *Repository) PostResetPasswordChange(w http.ResponseWriter, r *http.Request) {

	// Parse the form data from the request.
	// If error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form wth parsed data
	form := forms.New(r.PostForm)
	passReset := models.ResetPassword{
		Token:              r.Form.Get("reset_token"),
		NewPassword:        r.Form.Get("new_password"),
		NewPasswordConfirm: r.Form.Get("confirm_new_password"),
	}

	// Form validation to form fileds
	form.Required("reset_token", "new_password", "confirm_new_password")

	// Create a data map to store teh reset_password
	data := make(map[string]interface{})
	data["reset_password"] = passReset

	// check if new password and confirm password matches
	if passReset.NewPassword != passReset.NewPasswordConfirm {
		form.Errors.Add("new_password", "password mismatch")
		form.Errors.Add("confirm_new_password", "password mismatch")
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// add password validations
	form.MinLength("new_password", 8)
	form.HasUpperCase("new_password")
	form.HasLowerCase("new_password")
	form.HasSpecialCharacter("new_password")
	form.HasNumber("new_password")

	// get the reset token
	resetToken := userStore.PasswordResetRepo[passReset.Token]
	if time.Now().After(resetToken.ExpiresAt) {
		form.Errors.Add("reset_token", "Token is invalid or expired")
	}

	// if form is not valid render the "reset-password-change.page.tmpl" with form and data
	if !form.Valid() {
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Encrypt the new password.
	hashed_password, err := helpers.EncryptPassword(passReset.NewPassword)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// Store the hashed password and email in user model.
	user := userStore.Users[resetToken.Email]
	user.Password = hashed_password
	userStore.Users[resetToken.Email] = user

	// Call ChangePassword interface to change the password.
	// If any error occurs, a server error is returned
	if err := m.DB.ChangePassword(hashed_password, resetToken.Email); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Deletes the token stored in cache.
	delete(userStore.PasswordResetRepo, resetToken.Token)
	// To do add notification for successfull password change
	body := fmt.Sprintf(`
		<h1>Password Reset Successfull<h1>
			<p>You can login from here<p><br>
			<button><a href="%s/user/login">Login</a></button>
	`, r.Header["Origin"][0])
	msg := models.MailData{
		To:      resetToken.Email,
		From:    "admin@bookworm.com",
		Subject: "Password Reset Successfull",
		Content: body,
	}

	// send the msg to email channel
	m.App.MailChan <- msg

	// Redirect user to login page if password reset is successfull
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminAllGenres renders the admin gender page.
// It takes HTTP response writer and request as parameters.
// It fetched all genre records from db and renders the page.
func (m *Repository) AdminAllGenres(w http.ResponseWriter, r *http.Request) {

	// The function calls AllGenre interface to retrive all the records from genre table.
	// If error occurs, a server error is returned
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a map that holds the genres data
	data := make(map[string]interface{})
	data["genres"] = genres

	// Render the template with nill form and data
	render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminAddGenre is a handler that handles add new genre.
// It takes HTTP response writer and request as parameters.
// It parses the data from form, validates it, check for existing genres and only then add new genre.
// Finally, admin is redirected to all genres page if adding genre is successfull.
func (m *Repository) PostAdminAddGenre(w http.ResponseWriter, r *http.Request) {

	// Parse the form to populate post form.
	// If any error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Retrives all genres record from db using AllGenre() interface
	// If any error occurs, a server error is returned.
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Initiate a new form with post form values
	form := forms.New(r.PostForm)

	// Add form field validation
	form.Required("title")

	// Create a Genre model that holds the form data
	add_genre := models.Genre{
		Title: r.Form.Get("title"),
	}

	// Check if genre exists using GenreExists() interface
	// If any error occurs, a server error is returned
	exists, err := m.DB.GenreExists(add_genre.Title)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// If exists then add form error
	if exists {
		form.Errors.Add("title", "Genre already exists")
	}

	// Create a data map that holds genres and add_genres
	data := make(map[string]interface{})
	data["genres"] = genres
	data["add_genre"] = add_genre

	// If form is not valid render "admin-allgenres.page.tmpl" with form and data
	if !form.Valid() {
		render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// If form is valid then call InsertGenre interface to add new genre to db
	// If any error occurs, a server error is returned.
	if err := m.DB.InsertGenre(&add_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Finally, admin is redirected to the all genres pages
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}

// AdminGetGenreByID renders the genre detail and update form page.
// It mainly handle the get request method.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetGenreByID(w http.ResponseWriter, r *http.Request) {

	// Retrive "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Retrive the genre from db using GetGenreByID interface with id as parameter.
	// If any error occurs, a server error is returned
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a  data map and store the genre model
	data := make(map[string]interface{})
	data["genre"] = genre

	// Render the "admin-genre-read-update.page.tmpl" template with empty form and data
	render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminGetGenreByID fetch the specific genre from the database in admin interface as well as handles the update for the genre.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminGetGenreByID(w http.ResponseWriter, r *http.Request) {

	// Retrive "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It parse the form to populate the PostForm.
	// If any error occurs, a server error is returned
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It creates a new form with the post form values.
	form := forms.New(r.PostForm)

	// Then initate a Genre model that stores the updated values.
	update_genre := models.Genre{
		ID:    id,
		Title: r.Form.Get("title"),
	}

	// Before updating the function calls GenreExists interface to check for existing genres.
	exists, err := m.DB.GenreExists(update_genre.Title)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// if exists add a error to the form.
	if exists {
		form.Errors.Add("title", "Genre already exists")
	}

	// Add a validation to the form field.
	form.Required("title")

	// retrive genre using the id
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a data map and store the genre model
	data := make(map[string]interface{})
	data["genre"] = genre

	// If form is invalid render "admin-genre-read-update.page.tmpl" page with form and data
	if !form.Valid() {
		render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// If form is valid, then call UpdateGenre interface to update the genre.
	// If any error occurs, a server error is returned.
	if err := m.DB.UpdateGenre(&update_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If successull, then admin is redirected to genre detail page.
	http.Redirect(w, r, fmt.Sprintf("/admin/genres/detail/%d", id), http.StatusSeeOther)
}

// AdminDeleteGenre deletes the genre from the database in admin context.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminDeleteGenre(w http.ResponseWriter, r *http.Request) {

	// Fetch the parameter "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It calls the DeleteGenre interface with passing id as parameter to delete the record.
	// If any error occurs, a server error is retured.
	if err := m.DB.DeleteGenre(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// If successfull, admin is redirected to all genres page.
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}

// AdminAllPublisher renders admin all publisher page
func (m *Repository) AdminAllPublusher(w http.ResponseWriter, r *http.Request) {

	// Retrive all publishers record from database using AllPublishers database
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		// returns a server error if any error occurs
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores the publishers
	data := make(map[string]interface{})
	data["publishers"] = publishers

	// render the "admin-allpublishers.page.tmpl" with data and new empty form
	render.Template(w, r, "admin-allpublishers.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeletePublisher handles the POST request to delete a publisher.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The publisher is deleted from the database using the ID.
// If there is an error during the deletion process, a server error is returned.
// Otherwise, the user is redirected to the publishers admin page.
func (m *Repository) PostAdminDeletePublisher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeletePublisher(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/publishers", http.StatusSeeOther)
}

// AdminGetPublisherDetailByID handles the GET request to retrieve the details of a publisher by its ID.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The publisher is retrieved from the database using the ID.
// The data map is populated with the publisher object.
// The template is rendered with the form object and data map.
func (m *Repository) AdminGetPublisherDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publisher, err := m.DB.GetPublisherByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["publisher"] = publisher
	render.Template(w, r, "admin-publisherdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdatePublisher handles the post method logic for updating a publisher.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The form is parsed and a new form object is created.
// The established date is converted to an integer.
// The publisher object is populated with the form data.
// The ID of the publisher is set.
// The "name" field is marked as required.
// The data map is populated with the publisher object.
// If the form is not valid, the template is rendered with the form object and data map.
// The publisher is updated in the database.
// The user is redirected to the publisher detail page.
func (m *Repository) PostAdminUpdatePublisher(w http.ResponseWriter, r *http.Request) {
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
	establishedDate, _ := strconv.Atoi(r.Form.Get("established_date"))
	publisher := models.Publisher{
		Name:            r.Form.Get("name"),
		Description:     r.Form.Get("description"),
		Pic:             r.Form.Get("pic"),
		Address:         r.Form.Get("address"),
		Phone:           r.Form.Get("phone"),
		Email:           r.Form.Get("email"),
		Website:         r.Form.Get("website"),
		EstablishedDate: establishedDate,
		Latitude:        r.Form.Get("latitude"),
		Longitude:       r.Form.Get("longitude"),
	}
	publisher.ID = id
	form.Required("name")
	data := make(map[string]interface{})
	data["publisher"] = publisher
	if !form.Valid() {
		render.Template(w, r, "admin-publisherdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdatePublisher(&publisher); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/publishers/detail/%d", id), http.StatusSeeOther)
}

// AdminInsertPublisher handles the logic for rendering the publisher insert form.
// It takes the HTTP response writer and request as parameters.
// An empty publisher object is created.
// The data map is populated with the empty publisher object.
// The template is rendered with the form object and data map.
func (m *Repository) AdminInsertPublisher(w http.ResponseWriter, r *http.Request) {
	var emptyPublisher models.Publisher
	data := make(map[string]interface{})
	data["publisher"] = emptyPublisher
	render.Template(w, r, "admin-publisherinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertPublisher handles the logic for inserting a new publisher using the POST method.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// It creates a new form object and validates the required fields.
// The established date is converted to an integer.
// The picture path is obtained by uploading an image using the helper function.
// A new publisher object is created with the form data and the picture path.
// The publisher data is stored in the "publisher" key of the data map.
// If the form is not valid, the template is rendered with the form errors and data.
// If there is an error during insertion, a server error is returned.
// Otherwise, the publisher is inserted into the database and a redirect response is sent to the publishers list page.
func (m *Repository) PostAdminInsertPublisher(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)

	establishedDate, err := strconv.Atoi(r.Form.Get("established_date"))
	if err != nil {
		form.Errors.Add("established_date", "Invalid established date")
	}
	pic_path, err := helpers.AdminPublicUploadImage(r, "pic", "publisher", establishedDate)
	if err != nil {
		form.Errors.Add("pic", "No picture was choosen")
	}
	publisher := &models.Publisher{
		Name:            r.Form.Get("name"),
		Description:     r.Form.Get("description"),
		Pic:             pic_path,
		Address:         r.Form.Get("address"),
		Phone:           r.Form.Get("phone"),
		Email:           r.Form.Get("email"),
		Website:         r.Form.Get("website"),
		EstablishedDate: establishedDate,
		Latitude:        r.Form.Get("latitude"),
		Longitude:       r.Form.Get("longitude"),
	}
	data := make(map[string]interface{})
	data["publisher"] = publisher
	if !form.Valid() {
		render.Template(w, r, "admin-publisherinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertPublisher(publisher); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/publishers", http.StatusSeeOther)
}

// Author handlers

// AdminAllAuthor retrieves all authors from the database and renders the "admin-allauthors.page.tmpl" template.
// It takes the HTTP response writer and request as parameters.
// The function calls the AllAuthor interface to retrieve all authors from the database.
// If an error occurs during retrieval, a server error is returned.
// The authors are stored in the "authors" key of the data map.
// The function renders the template with the authors data.
func (m *Repository) AdminAllAuthor(w http.ResponseWriter, r *http.Request) {
	authors, err := m.DB.AllAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["authors"] = authors
	render.Template(w, r, "admin-allauthors.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// PostAdminDeleteAuthor handles the deletion of an author.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function calls the DeleteAuthor interface to delete the author from the database.
// If an error occurs during deletion, a server error is returned.
// The function redirects the user to the "/admin/authors" page.
func (m *Repository) PostAdminDeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteAuthor(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

// AdminGetAuthorDetailByID retrieves the author details by ID and renders the admin-authordetail page.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function calls the GetAuthorByID interface to retrieve the author from the database.
// If an error occurs during retrieval, a server error is returned.
// The function creates a data map and adds the author to the "author" key in the data map.
// The function renders the "admin-authordetail.page.tmpl" template with the data map.
func (m *Repository) AdminGetAuthorDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author, err := m.DB.GetAuthorByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["author"] = author
	render.Template(w, r, "admin-authordetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateAuthor handles the update author logic.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function parses the form data from the request.
// If an error occurs during parsing, a server error is returned.
// It creates a new form and adds the parsed form data to it.
// The function retrieves the date of birth from the form and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function creates a new Author model with the form data and the converted date of birth.
// The author's ID is set to the retrieved ID from the URL parameter.
// The function creates a data map and adds the author to the "author" key in the data map.
// If the form is not valid, the function renders the "admin-authordetail.page.tmpl" template
// with the form errors and the data map.
// The function calls the UpdateAuthor interface to update the author in the database.
// If an error occurs during the update, a server error is returned.
// The function redirects the user to the author's detail page.
func (m *Repository) PostAdminUpdateAuthor(w http.ResponseWriter, r *http.Request) {
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
	dob, err := strconv.Atoi(r.Form.Get("date_of_birth"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author := models.Author{
		ID:              id,
		FirstName:       r.Form.Get("first_name"),
		LastName:        r.Form.Get("last_name"),
		Bio:             r.Form.Get("bio"),
		DateOfBirth:     dob,
		Email:           r.Form.Get("email"),
		CountryOfOrigin: r.Form.Get("country_of_origin"),
		Avatar:          r.Form.Get("avatar"),
	}
	data := make(map[string]interface{})
	data["author"] = author
	if !form.Valid() {
		render.Template(w, r, "admin-authordetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateAuthor(&author); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/authors/detail/%d", id), http.StatusSeeOther)
}

// AdminInsertAuthor renders the page inserting new Author.
// It hanldes the get method for inserting Author
// It takes HTTP response writer and response as paramters.
// It creates an empty Author model.
// It create a data map that stores the empty Author model.
// Finally, "admin-authorinsert.page.tmpl" is rendered with additional data.
func (m *Repository) AdminInsertAuthor(w http.ResponseWriter, r *http.Request) {
	var emptyAuthor models.Author
	data := make(map[string]interface{})
	data["author"] = emptyAuthor
	render.Template(w, r, "admin-authorinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertAuthor handles the insertion of a new author.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// It creates a new form and adds the parsed form data to it.
// The function retrieves the date of birth from the form and converts it to an integer.
// If the conversion fails, an error is added to the form errors.
// The function uploads the author's avatar image using the AdminPublicUploadImage helper function.
// If the upload fails, an error is added to the form errors.
// The function creates a new Author model with the form data and the uploaded avatar.
// The "date_of_birth" field is set to the converted date of birth.
// The function adds the required validation for the "date_of_birth" field.
// It creates a data map and adds the author to the "author" key in the data map.
// If the form is not valid, the function renders the "admin-authorinsert.page.tmpl" template
// with the form errors and the data map.
// The function calls the InsertAuthor interface to insert the author into the database.
// If an error occurs during the insertion, a server error is returned.
// The function redirects the user to the "/admin/authors" page.
func (m *Repository) PostAdminInsertAuthor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)

	dob, err := strconv.Atoi(r.Form.Get("date_of_birth"))
	if err != nil {
		form.Errors.Add("date_of_birth", "Invalid date of birth")
	}
	avatar, err := helpers.AdminPublicUploadImage(r, "avatar", "author", dob)
	if err != nil {
		form.Errors.Add("avatar", "No picture was choosen")
	}
	author := &models.Author{
		FirstName:       r.Form.Get("first_name"),
		LastName:        r.Form.Get("last_name"),
		Bio:             r.Form.Get("bio"),
		DateOfBirth:     dob,
		Email:           r.Form.Get("email"),
		CountryOfOrigin: r.Form.Get("country_of_origin"),
		Avatar:          avatar,
	}
	form.Required("date_of_birth")
	data := make(map[string]interface{})
	data["author"] = author
	if !form.Valid() {
		render.Template(w, r, "admin-authorinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertAuthor(author); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

// AdminAllLanguage retrieves all languages from the database and renders the admin-alllanguages page.
// It takes the HTTP response writer and request as parameters.
// The function calls the AllLanguage interface to retrieve all languages from the database.
// If an error occurs during the retrieval, a server error is returned.
// The function creates a data map and adds the retrieved languages to the "languages" key in the data map.
// It creates an empty Language model and adds it to the "language" key in the data map.
// The function renders the "admin-alllanguages.page.tmpl" template with the data map and an empty form.
func (m *Repository) AdminAllLanguage(w http.ResponseWriter, r *http.Request) {
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["languages"] = languages
	var emptyLanguage models.Language
	data["language"] = emptyLanguage
	render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeleteLanguage deletes a language from the database in the admin context.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the "id" parameter from the URL and converts it to an integer.
// If there is an error during the conversion, a server error is returned.
// The function calls the DeleteLanguage interface to delete the language with the specified ID from the database.
// If an error occurs during the deletion, a server error is returned.
// After successful deletion, the function redirects the user to the "/admin/languages" page.
func (m *Repository) PostAdminDeleteLanguage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteLanguage(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)
}

// PostAdminUpdateLanguage handles the update language logic
func (m *Repository) PostAdminUpdateLanguage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["languages"] = languages
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language := models.Language{
		ID:       id,
		Language: r.Form.Get("language"),
	}
	form.Required("language")
	data["language"] = language
	exists, err := m.DB.LanguageExists(language.Language)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("language", "This language already exists")
	}
	if !form.Valid() {
		render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateLanguage(&language); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)
}

// PostAdminInsertLanguage inserts a new language into the database in the admin context.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// If there is an error during form parsing, a server error is returned.
// The function creates a new form instance and initializes a language model with the language value from the form.
// It sets the "language" field as required in the form.
// The function creates a data map and adds the "add_language" key with the language model as its value.
// It calls the AllLanguage interface to retrieve all existing languages from the database.
// If an error occurs during the retrieval, a server error is returned.
// The function adds the retrieved languages to the data map with the "languages" key.
// It checks if the language already exists in the database using the LanguageExists interface.
// If the language exists, an error is added to the form's errors.
// If the form is not valid, the function renders the "admin-alllanguages.page.tmpl" template with the form errors and data.
// If there are no form errors, the function calls the InsertLanguage method to insert the language into the database.
// If an error occurs during the insertion, a server error is returned.
// After successful insertion, the function redirects the user to the "/admin/languages" page.
func (m *Repository) PostAdminInsertLanguage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	language := models.Language{
		Language: r.Form.Get("language"),
	}
	form.Required("language")
	data := make(map[string]interface{})
	data["add_language"] = language
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["languages"] = &languages
	stat, _ := m.DB.LanguageExists(language.Language)
	if stat {
		form.Errors.Add("language", "This language already exists")
	}
	if !form.Valid() {
		render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertLanguage(&language); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)

}

// AdminAllBook handles logic for reteriving all Books in admin page.
// It takes HTTP response writer and request as parameters.
// The function calls AllBook interface to reterive all the books from the database
// If any error occurs, a server error is returned.
// A data map is created that stores the books
// Finally, "admin-allbooks.page.tmpl" go template is rendered with data.
func (m *Repository) AdminAllBook(w http.ResponseWriter, r *http.Request) {
	books, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["books"] = books
	render.Template(w, r, "admin-allbooks.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteBook(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
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
		helpers.ServerError(w, err)
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
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["book"] = book
	data["publishers"] = publishers
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
		helpers.ServerError(w, err)
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
	data["book"] = book
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
		helpers.ServerError(w, err)
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
	http.Redirect(w, r, fmt.Sprintf("/admin/books/detail/%d", book.ID), http.StatusSeeOther)

}

// Start of handler for admin book-author

// AdminAllBookAuthor retrieves all book-author relationships and renders the "admin-allbookauthors.page.tmpl" template.
// It takes the HTTP response writer and request as parameters. The function calls the database's AllBookAuthor method
// to retrieve all book-author relationships. If an error occurs during the retrieval, a server error is returned.
// The function also retrieves all books and authors from the database using the AllBook and AllAuthor methods, respectively.
// If any errors occur during the retrieval, a server error is returned. The retrieved data, including book-authors,
// book-author, all authors, and all books, is stored in a data map. The function then renders the template,
// passing the data map and an empty form to the template for rendering.
func (m *Repository) AdminAllBookAuthor(w http.ResponseWriter, r *http.Request) {
	bookAuthors, err := m.DB.AllBookAuthor()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var bookAuthor models.BookAuthor
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
	data["bookAuthors"] = bookAuthors
	data["bookAuthor"] = bookAuthor
	data["allAuthors"] = allAuthors
	data["allBooks"] = allBooks
	render.Template(w, r, "admin-allbookauthors.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
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

	data["allBooks"] = allBooks
	data["allAuthors"] = allAuthors
	data["bookAuthor"] = bookAuthor
	data["bookAuthors"] = bookAuthors
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

	http.Redirect(w, r, "/admin/bookAuthors", http.StatusSeeOther)
}

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

	http.Redirect(w, r, "/admin/bookGenres", http.StatusSeeOther)
}

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
	data := make(map[string]interface{})
	data["bookLanguages"] = bookLanguages
	data["bookLanguage"] = bookLanguage
	data["allLanguages"] = allLanguages
	data["allBooks"] = allBooks
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

	data["allBooks"] = allBooks
	data["allLanguages"] = allLanguages
	data["bookLanguage"] = bookLanguage
	data["bookLanguages"] = bookLanguages
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

	http.Redirect(w, r, "/admin/bookLanguages", http.StatusSeeOther)
}
