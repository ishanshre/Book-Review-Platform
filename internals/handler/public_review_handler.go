package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

func (m *Repository) PublicCreateReview(w http.ResponseWriter, r *http.Request) {
	isbn, err := strconv.ParseInt(chi.URLParam(r, "isbn"), 10, 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book, err := m.DB.GetBookByISBN(isbn)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var review models.Review
	data := make(map[string]interface{})
	data["review"] = review
	data["book"] = book
	render.Template(w, r, "public_review_create.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostPublicCreateReview(w http.ResponseWriter, r *http.Request) {
	isbn, err := strconv.ParseInt(chi.URLParam(r, "isbn"), 10, 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book, err := m.DB.GetBookByISBN(isbn)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.Get(r.Context(), "user_id").(int)
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	rating, err := strconv.ParseFloat(r.Form.Get("rating"), 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	body := r.Form.Get("body")
	form.Required("rating", "body")
	form.MinFloatValue("rating", 1)
	form.MaxFloatValue("rating", 5)
	form.MaxLength("body", 10000)
	data := make(map[string]interface{})
	review := &models.Review{
		Rating:    rating,
		Body:      body,
		BookID:    book.ID,
		UserID:    user_id,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	data["review"] = review
	data["book"] = book
	if !form.Valid() {
		render.Template(w, r, "public_review_create.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertReview(review); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%d", isbn), http.StatusSeeOther)

}

func (m *Repository) PostPublicDeleteReview(w http.ResponseWriter, r *http.Request) {
	isbn, err := strconv.ParseInt(chi.URLParam(r, "isbn"), 10, 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	review_id, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	review, err := m.DB.GetReviewByID(review_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := m.App.Session.Get(r.Context(), "user_id").(int)
	if user_id != int(review.UserID) {
		helpers.PageNotFound(w, errors.New("user not authorized"))
		return
	}

	if err := m.DB.DeleteReview(review_id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%d", isbn), http.StatusSeeOther)
}
