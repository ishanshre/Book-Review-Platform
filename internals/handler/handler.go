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

// PostLogin handles the post method of login
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

// PostRegister handles the post method for registering new user.
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
	userStatus, _ := m.DB.UsernameExists(register.Username)
	if userStatus {
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

// Logout log the user out
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

// PersonalProfile returns profile to authenticated user.
// Authenticated user can only view this personal profile
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

// AdminAllUser redners users page. It is available for admin user only
func (m *Repository) AdminAllUsers(w http.ResponseWriter, r *http.Request) {
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

// AdminGetUserDetailByID renders user detail page and for admin user only
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

// AdminUpdateUser updates user detail
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

// PostAdminUserDeleteByID renders a confim page to delete users
func (m *Repository) PostAdminUserDeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	if id == user_id {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	if err := m.DB.DeleteUser(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminUserAdd renders page for adding user by admin
func (m *Repository) AdminUserAdd(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(map[string]interface{})
	data["register"] = emptyUser
	render.Template(w, r, "admin-usercreate.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUserAdd handles post method for creating user by admin
func (m *Repository) PostAdminUserAdd(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("username", "email", "password")
	form.MinLength("username", 5)
	form.MinLength("password", 8)
	form.HasLowerCase("password")
	form.HasUpperCase("password")
	form.HasNumber("password")
	form.HasSpecialCharacter("password")

	register_user := models.User{
		Username: r.Form.Get("username"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	data := make(map[string]interface{})
	data["register"] = register_user

	if !form.Valid() {
		render.Template(w, r, "admin-usercreate.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
	hashed_password, err := helpers.EncryptPassword(register_user.Password)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	register_user.Password = hashed_password
	if err := m.DB.AdminInsertUser(&register_user); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// ResetPassword render the password reset page
func (m *Repository) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var emptyUser models.User
	data := make(map[string]interface{})
	data["reset_user"] = emptyUser
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

// PostResetPassword handles the post method that takes in the email address
func (m *Repository) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
	}
	form := forms.New(r.PostForm)
	form.Required("email")
	reset_user := models.User{
		Email: r.Form.Get("email"),
	}
	data := make(map[string]interface{})
	data["reset_user"] = reset_user
	_, err := m.DB.EmailExists(reset_user.Email)
	if err != nil {
		form.Errors.Add("email", "This email does not exist")
	}
	if !form.Valid() {
		render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	_ = userStore.Users[reset_user.Email]

	token, err := helpers.GenerateRandomToken(32)
	if err != nil {
		helpers.ServerError(w, err)
	}
	resetToken := models.PasswordResetToken{
		Token:     token,
		Email:     reset_user.Email,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	userStore.PasswordResetRepo[token] = resetToken
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
	http.Redirect(w, r, "/user/reset", http.StatusSeeOther)
}

// ResetPasswordChange renders password change form with new and confirm password
func (m *Repository) ResetPasswordChange(w http.ResponseWriter, r *http.Request) {
	var emptyPass models.ResetPassword
	data := make(map[string]interface{})
	data["reset_password"] = emptyPass
	render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostResetPasswordChange handles post method and changes the password
func (m *Repository) PostResetPasswordChange(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	passReset := models.ResetPassword{
		Token:              r.Form.Get("reset_token"),
		NewPassword:        r.Form.Get("new_password"),
		NewPasswordConfirm: r.Form.Get("confirm_new_password"),
	}
	form.Required("reset_token", "new_password", "confirm_new_password")
	data := make(map[string]interface{})
	data["reset_password"] = passReset
	if passReset.NewPassword != passReset.NewPasswordConfirm {
		form.Errors.Add("new_password", "password mismatch")
		form.Errors.Add("confirm_new_password", "password mismatch")
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	form.MinLength("new_password", 8)
	form.HasUpperCase("new_password")
	form.HasLowerCase("new_password")
	form.HasSpecialCharacter("new_password")
	form.HasNumber("new_password")

	resetToken := userStore.PasswordResetRepo[passReset.Token]
	if time.Now().After(resetToken.ExpiresAt) {
		form.Errors.Add("reset_token", "Token is invalid or expired")
	}

	if !form.Valid() {
		render.Template(w, r, "reset-password-change.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	hashed_password, err := helpers.EncryptPassword(passReset.NewPassword)
	if err != nil {
		helpers.ServerError(w, err)
	}
	user := userStore.Users[resetToken.Email]
	user.Password = hashed_password
	userStore.Users[resetToken.Email] = user
	if err := m.DB.ChangePassword(hashed_password, resetToken.Email); err != nil {
		helpers.ServerError(w, err)
		return
	}
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
	m.App.MailChan <- msg
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminAllGenres renders the admin gender page
func (m *Repository) AdminAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["genres"] = genres
	render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminAddGenre is a handler that handles admin add genre
func (m *Repository) PostAdminAddGenre(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title")
	add_genre := models.Genre{
		Title: r.Form.Get("title"),
	}
	data := make(map[string]interface{})
	data["genres"] = genres
	data["add_genre"] = add_genre
	if !form.Valid() {
		render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertGenre(&add_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}

// AdminGetGenreByID renders the genre detail and update page
func (m *Repository) AdminGetGenreByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["genre"] = genre
	render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminGetGenreByID updates the existing genre
func (m *Repository) PostAdminGetGenreByID(w http.ResponseWriter, r *http.Request) {
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
	update_genre := models.Genre{
		ID:    id,
		Title: r.Form.Get("title"),
	}
	form.Required("title")
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	res, err := m.DB.GenreExists(update_genre.Title)
	if err != nil && !res {
		form.Errors.Add("title", "Genre already exists")
	}
	data["genre"] = genre
	if !form.Valid() {
		render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateGenre(&update_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/genres/detail/%d", id), http.StatusSeeOther)
}

// AdminDeleteGenre delete the genre
func (m *Repository) AdminDeleteGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteGenre(id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}

// AdminAllPublisher renders admin all publisher page
func (m *Repository) AdminAllPublusher(w http.ResponseWriter, r *http.Request) {
	publishers, err := m.DB.AllPublishers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["publishers"] = publishers
	render.Template(w, r, "admin-allpublishers.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// AdminDeletePublisher handles the delete logic
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

// AdminGetPublisherDetailByID handles the logic of displaying publisher detail by ID
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

// PostAdminUpdatePublisher handles the update login for publisher
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
		Latitude:        r.Form.Get("name"),
		Longitude:       r.Form.Get("name"),
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

// AdminInsertPublisher renders the insert publisher page
func (m *Repository) AdminInsertPublisher(w http.ResponseWriter, r *http.Request) {
	var emptyPublisher models.Publisher
	data := make(map[string]interface{})
	data["publisher"] = emptyPublisher
	render.Template(w, r, "admin-publisherinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertPublisher handles the post method logic for publisher
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

// AdminAllAuthor handles logic for all authors in admin page
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

// PostAdminDeleteAuthor handles author delete logic
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

// AdminGetAuthorDetailByID handles the logic of displaying Author detail by ID
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

// PostAdminUpdateAuthor handles update logic of author
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

// AdminInsertAuthor renders the insert Author page
func (m *Repository) AdminInsertAuthor(w http.ResponseWriter, r *http.Request) {
	var emptyAuthor models.Author
	data := make(map[string]interface{})
	data["author"] = emptyAuthor
	render.Template(w, r, "admin-authorinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertAuthor handles the post method logic for publisher
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

// AdminAllLanguage renders admin all Language page
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

// PostAdminDeleteLanguage handles the delete logic
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
	// stat, _ := m.DB.LanguageExists(language.Language)
	// if stat {
	// 	form.Errors.Add("language", "This language already exists")
	// }
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

// PostAdminInsertLanguage handles the insert language logic
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

// AdminAllBook handles logic for all Books in admin page
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

// PostAdminDeleteBook handles the delete logic
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

// AdminGetBookDetailByID handles the logic of displaying Book detail by ID
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

// AdminInsertBook handles the get method
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

// PostAdminInsertBook handles insert book logic
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
	log.Println(book)
	if err := m.DB.InsertBook(&book); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/books", http.StatusSeeOther)
}

// PostAdminUpdateBook handles Update book logic
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
