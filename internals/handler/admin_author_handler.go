package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// Author handlers

// AdminAllAuthor retrieves all authors from the database and renders the "admin-allauthors.page.tmpl" template.
// It takes the HTTP response writer and request as parameters.
// The function calls the AllAuthor interface to retrieve all authors from the database.
// If an error occurs during retrieval, a server error is returned.
// The authors are stored in the "authors" key of the data map.
// The function renders the template with the authors data.
func (m *Repository) AdminAllAuthor(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["base_path"] = base_authors_path
	render.Template(w, r, "admin-allauthors.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllAuthorApi(w http.ResponseWriter, r *http.Request) {
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
	filteredAuthors, err := m.DB.AllAuthorsFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredAuthors)
}

// PostAdminDeleteAuthor handles the deletion of an author.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function calls the DeleteAuthor interface to delete the author from the database.
// If an error occurs during deletion, a server error is returned.
// The function redirects the user to the "/admin/authors" page.
func (m *Repository) PostAdminDeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteAuthor(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Author Deleted")

	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

// AdminGetAuthorDetailByID retrieves the author details by ID and renders the admin-authordetail page.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function calls the GetAuthorByID interface to retrieve the author from the database.
// If an error occurs during retrieval, a server error is returned.
// The function creates a data map and adds the author to the "author" key in the data map.
// The function renders the "admin-authordetail.page.tmpl" template with the data map.
func (m *Repository) AdminGetAuthorDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author, err := m.DB.GetAuthorByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["author"] = author
	data["base_path"] = base_authors_path
	render.Template(w, r, "admin-authordetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdateAuthor handles the update author logic.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the author ID from the URL parameter and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function parses the form data from the request.
// If an error occurs during parsing, a server error is returned.
// It creates a new form and adds the parsed form data to it.
// The function retrieves the date of birth from the form and converts it to an integer.
// If the conversion fails, a server error is returned.
// The function creates a new Author model with the form data and the converted date of birth.
// The author's ID is set to the retrieved ID from the URL parameter.
// The function creates a data map and adds the author to the "author" key in the data map.
// If the form is not valid, the function renders the "admin-authordetail.page.tmpl" template
// with the form errors and the data map.
// The function calls the UpdateAuthor interface to update the author in the database.
// If an error occurs during the update, a server error is returned.
// The function redirects the user to the author's detail page.
func (m *Repository) PostAdminUpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	dob, err := strconv.Atoi(r.Form.Get("date_of_birth"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	author := models.Author{
		ID:              id,
		FirstName:       r.Form.Get("first_name"),
		LastName:        r.Form.Get("last_name"),
		Bio:             r.Form.Get("bio"),
		DateOfBirth:     dob,
		Email:           r.Form.Get("email"),
		CountryOfOrigin: r.Form.Get("country_of_origin"),
		Avatar:          r.Form.Get("avatar"),
	}
	data := make(map[string]interface{})
	data["author"] = author
	data["base_path"] = base_authors_path

	if !form.Valid() {
		render.Template(w, r, "admin-authordetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateAuthor(&author); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Author Updated")

	http.Redirect(w, r, fmt.Sprintf("/admin/authors/detail/%d", id), http.StatusSeeOther)
}

// AdminInsertAuthor renders the page inserting new Author.
// It hanldes the get method for inserting Author
// It takes HTTP response writer and response as paramters.
// It creates an empty Author model.
// It create a data map that stores the empty Author model.
// Finally, "admin-authorinsert.page.tmpl" is rendered with additional data.
func (m *Repository) AdminInsertAuthor(w http.ResponseWriter, r *http.Request) {
	var emptyAuthor models.Author
	data := make(map[string]interface{})
	data["author"] = emptyAuthor
	data["base_path"] = base_authors_path

	render.Template(w, r, "admin-authorinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertAuthor handles the insertion of a new author.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// It creates a new form and adds the parsed form data to it.
// The function retrieves the date of birth from the form and converts it to an integer.
// If the conversion fails, an error is added to the form errors.
// The function uploads the author's avatar image using the AdminPublicUploadImage helper function.
// If the upload fails, an error is added to the form errors.
// The function creates a new Author model with the form data and the uploaded avatar.
// The "date_of_birth" field is set to the converted date of birth.
// The function adds the required validation for the "date_of_birth" field.
// It creates a data map and adds the author to the "author" key in the data map.
// If the form is not valid, the function renders the "admin-authorinsert.page.tmpl" template
// with the form errors and the data map.
// The function calls the InsertAuthor interface to insert the author into the database.
// If an error occurs during the insertion, a server error is returned.
// The function redirects the user to the "/admin/authors" page.
func (m *Repository) PostAdminInsertAuthor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)

	dob, err := strconv.Atoi(r.Form.Get("date_of_birth"))
	if err != nil {
		form.Errors.Add("date_of_birth", "Invalid date of birth")
	}
	idString := fmt.Sprintf("%d%s", dob, helpers.RandomAlphaNum(8))
	avatar, err := helpers.AdminPublicUploadImage(r, "avatar", "author", idString)
	if err != nil {
		form.Errors.Add("avatar", "No picture was choosen")
	}
	author := &models.Author{
		FirstName:       r.Form.Get("first_name"),
		LastName:        r.Form.Get("last_name"),
		Bio:             r.Form.Get("bio"),
		DateOfBirth:     dob,
		Email:           r.Form.Get("email"),
		CountryOfOrigin: r.Form.Get("country_of_origin"),
		Avatar:          avatar,
	}
	form.Required("date_of_birth")
	data := make(map[string]interface{})
	data["author"] = author
	data["base_path"] = base_authors_path

	if !form.Valid() {
		render.Template(w, r, "admin-authorinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertAuthor(author); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Author Inserted")

	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}
