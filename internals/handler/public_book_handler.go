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
	data := make(map[string]interface{})
	data["allGenre"] = allGenres
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
