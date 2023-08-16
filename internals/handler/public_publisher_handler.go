package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) PublisherWithBooksDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.PageNotFound(w, err)
		return
	}

	publisherWithBooks, err := m.DB.GetPublisherWithBookByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["publisher"] = publisherWithBooks.Publisher
	data["books"] = publisherWithBooks.Books
	render.Template(w, r, "public_publisher_detail.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
