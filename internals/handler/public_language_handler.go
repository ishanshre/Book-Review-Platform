package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) AllBookFilterByLanguage(w http.ResponseWriter, r *http.Request) {
	language := chi.URLParam(r, "language")
	data := make(map[string]interface{})
	data["language"] = language
	render.Template(w, r, "public_books_by_language.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AllBooksFilterByLanguageApi(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
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
	filteredBooks, err := m.DB.GetAllBooksByLanguage(limit, page, searchKey, sort, language)
	log.Println(filteredBooks)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBooks)
}
