package handler

import (
	"net/http"
	"strconv"

	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) AllBooksFilterFromBuyList(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "public_books_from_buyList.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AllBooksFilterFromBuyListApi(w http.ResponseWriter, r *http.Request) {
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
	user_id := m.App.Session.GetInt(r.Context(), "user_id")
	filteredBooks, err := m.DB.GetAllBooksFromBuyListByUserId(limit, page, user_id, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredBooks)
}
