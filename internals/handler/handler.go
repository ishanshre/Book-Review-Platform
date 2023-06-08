package handler

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
	"github.com/ishanshre/Book-Review-Platform/internals/repository"
	"github.com/ishanshre/Book-Review-Platform/internals/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
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
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(register.Password), 14)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	register.Password = string(hashed_password)
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
	})
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
	if err := m.DB.AdminInsertUser(&register_user); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
