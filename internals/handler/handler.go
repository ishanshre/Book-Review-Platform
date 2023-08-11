package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/driver"
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

const base_users_path = "/admin/users"
const base_genres_path = "/admin/genres"
const base_publishers_path = "/admin/publishers"
const base_authors_path = "/admin/authors"
const base_languages_path = "/admin/languages"
const base_books_path = "/admin/books"
const base_bookAuthors_path = "/admin/bookAuthors"
const base_bookGenres_path = "/admin/bookGenres"
const base_bookLanguages_path = "/admin/bookLanguages"
const base_readLists_path = "/admin/readLists"
const base_buyLists_path = "/admin/buyLists"
const base_followers_path = "/admin/followers"
const base_reviews_path = "/admin/reviews"
const base_contacts_path = "/admin/contacts"

// ClearSessionMessage clears the session message like flash, error and warning after being displayed
func (m *Repository) ClearSessionMessage(w http.ResponseWriter, r *http.Request) {
	msgType := chi.URLParam(r, "type")
	m.App.Session.Pop(r.Context(), msgType)
	w.WriteHeader(http.StatusOK)
}
