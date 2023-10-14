package handler

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) PublicAboutUs(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "public_about.page.tmpl", &models.TemplateData{})
}
