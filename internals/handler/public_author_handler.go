package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (h *Repository) AuthorFiltersApi(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "ASC"
	}
	authors, err := h.DB.AllAuthorsFilter(limit, page, search, sort)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, "error in fetching authors")
		return
	}
	helpers.ApiStatusOkData(w, authors)
}

func (h *Repository) AllAuthors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "public_authors.page.tmpl", &models.TemplateData{})
}

func (h *Repository) PublicGetAuthorByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ClientError(w, http.StatusNotFound)
		return
	}
	authorWithBooks, err := h.DB.GetAuthorWithBooks(id)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["author"] = authorWithBooks.Author
	data["books"] = authorWithBooks.Books
	render.Template(w, r, "public_author_detail.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
