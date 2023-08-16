package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) AllBookFilterByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")
	data := make(map[string]interface{})
	data["genre"] = genre
	render.Template(w, r, "public_books_by_genre.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AllBooksFilterByGenreApi(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
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
	if sort == "" {
		sort = "asc"
	}
	filteredBooks, err := m.DB.GetAllBooksByGenre(limit, page, searchKey, sort, genre)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBooks)
}
