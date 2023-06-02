package handler

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}
