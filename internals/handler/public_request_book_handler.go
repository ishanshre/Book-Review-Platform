package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) RequestBook(w http.ResponseWriter, r *http.Request) {
	var emptyRequest models.RequestedBook
	data := make(map[string]interface{})
	data["requestBook"] = emptyRequest
	render.Template(w, r, "public_request_book.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostRequestBook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	requested_by := m.App.Session.GetInt(r.Context(), "user_id")
	requestedBook := models.RequestedBook{
		BookTitle:     r.Form.Get("book_title"),
		Author:        r.Form.Get("author"),
		RequestedBy:   requested_by,
		RequestedDate: time.Now(),
	}
	form.Required("book_title", "author")
	form.MaxLength("book_title", 255)
	form.MaxLength("author", 255)
	data := make(map[string]interface{})
	data["requestBook"] = requestedBook
	if !form.Valid() {
		render.Template(w, r, "public_request_book.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertRequestedBook(&requestedBook); err != nil {
		helpers.ServerError(w, err)
		return
	}
	user, err := m.DB.GetGlobalUserByIDAny(requested_by)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	msg := models.MailData{
		To:      m.App.AdminEmail,
		From:    user.Email,
		Subject: fmt.Sprintf("Request for %s", requestedBook.BookTitle),
		Content: fmt.Sprintf("Requesting for book %s by %s", requestedBook.BookTitle, requestedBook.Author),
	}
	m.App.MailChan <- msg
	m.App.Session.Put(r.Context(), "flash", "Book Successfully requested")
	http.Redirect(w, r, "/request-book", http.StatusSeeOther)
}
