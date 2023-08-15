package handler

import (
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

// AdminAllReviews fetches all the record in Reviews.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminAllReviews(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["base_path"] = base_reviews_path
	render.Template(w, r, "admin-allreviews.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminAllReviewApi(w http.ResponseWriter, r *http.Request) {
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
	filterReviews, err := m.DB.ReviewFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filterReviews)
}

// AdminInsertReview handles the get method and renders the review add page.
// It takes HTTP response writer and request as parameters
func (m *Repository) AdminInsertReview(w http.ResponseWriter, r *http.Request) {

	// Initializing empty review to pass to template as a get request method
	var review models.Review

	// Fetching all the books from db
	// When error occurs, a server error is returned.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching all the users from db.
	// When error occurs, a server error is returned.
	allUsers, err := m.DB.AllUsers(10000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Creating the data map that will holds review, books and users data to pass to go template
	data := make(map[string]interface{})
	data["review"] = review
	data["allBooks"] = allBooks
	data["allUsers"] = allUsers
	data["base_path"] = base_reviews_path
	render.Template(w, r, "admin-reviewinsert.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminInsertReview handles the post method of the creating new review.
// It has HTTP response writer and request as parameters.
func (m *Repository) PostAdminInsertReview(w http.ResponseWriter, r *http.Request) {

	// Parse the form from the post method.
	// If any error occurs, a server error is returned
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching all the books from db
	// When error occurs, a server error is returned.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fetching all the users from db.
	// When error occurs, a server error is returned.
	allUsers, err := m.DB.AllUsers(10000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// initializing a new form
	form := forms.New(r.PostForm)

	// initializing a new data map
	data := make(map[string]interface{})
	rating, err := strconv.ParseFloat(r.Form.Get("rating"), 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	bookID, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	userID, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isActive, err := strconv.ParseBool(r.Form.Get("is_active"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	review := models.Review{
		Rating:    rating,
		Body:      r.Form.Get("body"),
		BookID:    bookID,
		UserID:    userID,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// add form field validation
	form.Required("rating", "book_id", "user_id", "is_active")
	form.MaxLength("body", 10000)
	var min float64 = 1.0
	var max float64 = 5.0
	form.MinFloatValue("rating", min)
	form.MaxFloatValue("rating", max)

	// Add stucts to pass to go templates
	data["review"] = review
	data["allBooks"] = allBooks
	data["allUsers"] = allUsers
	data["base_path"] = base_reviews_path
	if !form.Valid() {
		render.Template(w, r, "admin-reviewinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertReview(&review); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Review record added")

	http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)
}

// PostAdminDeleteReview Handles the post method for deleting Reviewer list record.
// It takes HTTP response writer and request as paramters
func (m *Repository) PostAdminDeleteReview(w http.ResponseWriter, r *http.Request) {

	// Parsing the review id from the url.
	// If any error occurs, a server error is returned.
	review_id, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// DeleteBuyList interface is used to deleting the record.
	if err := m.DB.DeleteReview(review_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Review record deleted")

	http.Redirect(w, r, "/admin/reviews", http.StatusSeeOther)
}

// AdminGetReviewByID handes the detail logic for Review table.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetReviewByID(w http.ResponseWriter, r *http.Request) {

	// Retrive user id and book id from the url.
	// Parse them into integer.
	// Return a server error if any error occurs while parsing them
	review_id, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the review using review id
	review, err := m.DB.GetReviewByID(review_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Get the book title using book_id
	book, err := m.DB.GetBookByID(review.BookID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get the user by using user_id
	user, err := m.DB.GetUserByID(review.UserID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = review.UserID
	// Get all books from the AllBook interface.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Get all user from the AllUsers interface.
	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores book, user, all books, all users and review
	data := make(map[string]interface{})
	data["book"] = book
	data["allBooks"] = allBooks
	data["user"] = user
	data["allUsers"] = allUsers
	data["review"] = review
	data["base_path"] = base_reviews_path

	// render the detail page with form and data
	render.Template(w, r, "admin-reviewdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateReview handles the post method for updating the  Review.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminUpdateReview(w http.ResponseWriter, r *http.Request) {

	// Fetches the author id and user id from url and parse them into integer.
	// If any error occurs, a server error is returned
	review_id, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Parse the form and populate the PostForm
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a new form with PostForm
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})

	rating, err := strconv.ParseFloat(r.Form.Get("rating"), 64)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	book_id, err := strconv.Atoi(r.Form.Get("book_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	isActive, err := strconv.ParseBool(r.Form.Get("is_active"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	getReview, _ := m.DB.GetReviewByID(review_id)
	review := models.Review{
		ID:        review_id,
		Rating:    rating,
		Body:      r.Form.Get("body"),
		BookID:    book_id,
		UserID:    user_id,
		IsActive:  isActive,
		CreatedAt: getReview.CreatedAt,
		UpdatedAt: time.Now(),
	}
	// Get the book title using book_id
	book, err := m.DB.GetBookByID(review.BookID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get the user by using user_id
	user, err := m.DB.GetUserByID(review.UserID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user.ID = review.UserID
	// Get all books from the AllBook interface.
	allBooks, err := m.DB.AllBook()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// Get all user from the AllUsers interface.
	allUsers, err := m.DB.AllUsers(100000, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// add form field validation
	form.Required("rating", "book_id", "user_id", "is_active")
	form.MaxLength("body", 10000)
	var min float64 = 1.0
	var max float64 = 5.0
	form.MinFloatValue("rating", min)
	form.MaxFloatValue("rating", max)
	data["book"] = book
	data["allBooks"] = allBooks
	data["user"] = user
	data["allUsers"] = allUsers
	data["review"] = review
	data["base_path"] = base_reviews_path
	if !form.Valid() {
		render.Template(w, r, "admin-reviewdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
	if err := m.DB.UpdateReview(&review); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Review record updated")

	http.Redirect(w, r, fmt.Sprintf("/admin/reviews/detail/%d", review_id), http.StatusSeeOther)
}
