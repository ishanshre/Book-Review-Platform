package handler

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Home handles the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	allGenres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	allBooks, err := m.DB.AllBookPage(10, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["allGenre"] = allGenres
	data["allBooks"] = allBooks
	render.Template(w, r, "public_home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
