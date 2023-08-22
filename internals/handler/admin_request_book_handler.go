package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) AdminAllRequestBookList(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["base_path"] = base_request_book_path
	render.Template(w, r, "admin-allrequestbooks.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllRequestedBookssApi(w http.ResponseWriter, r *http.Request) {
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
	filteredRequestedBookss, err := m.DB.RequestedBooksListFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredRequestedBookss)
}

func (m *Repository) AdminDeleteRequestedBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.PageNotFound(w, r, err)
		return
	}

	// The function calls DeleteUser interface to delete the user form the database
	if err := m.DB.DeleteRequestBooks(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Requested Book Deleted")
	// Redirect the admin to all users page
	http.Redirect(w, r, "/admin/request-books", http.StatusSeeOther)
}
