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

// AdminAllGenres renders the admin gender page.
// It takes HTTP response writer and request as parameters.
// It fetched all genre records from db and renders the page.
func (m *Repository) AdminAllGenres(w http.ResponseWriter, r *http.Request) {

	// The function calls AllGenre interface to retrive all the records from genre table.
	// If error occurs, a server error is returned
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a map that holds the genres data
	data := make(map[string]interface{})
	data["genres"] = genres
	data["base_path"] = base_genres_path

	// Render the template with nill form and data
	render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminAddGenre is a handler that handles add new genre.
// It takes HTTP response writer and request as parameters.
// It parses the data from form, validates it, check for existing genres and only then add new genre.
// Finally, admin is redirected to all genres page if adding genre is successfull.
func (m *Repository) PostAdminAddGenre(w http.ResponseWriter, r *http.Request) {

	// Parse the form to populate post form.
	// If any error occurs, a server error is returned.
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Retrives all genres record from db using AllGenre() interface
	// If any error occurs, a server error is returned.
	genres, err := m.DB.AllGenre()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Initiate a new form with post form values
	form := forms.New(r.PostForm)

	// Add form field validation
	form.Required("title")

	// Create a Genre model that holds the form data
	add_genre := models.Genre{
		Title: r.Form.Get("title"),
	}

	// Check if genre exists using GenreExists() interface
	// If any error occurs, a server error is returned
	exists, err := m.DB.GenreExists(add_genre.Title)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// If exists then add form error
	if exists {
		form.Errors.Add("title", "Genre already exists")
	}

	// Create a data map that holds genres and add_genres
	data := make(map[string]interface{})
	data["genres"] = genres
	data["add_genre"] = add_genre
	data["base_path"] = base_genres_path

	// If form is not valid render "admin-allgenres.page.tmpl" with form and data
	if !form.Valid() {
		render.Template(w, r, "admin-allgenres.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// If form is valid then call InsertGenre interface to add new genre to db
	// If any error occurs, a server error is returned.
	if err := m.DB.InsertGenre(&add_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Genre Added")

	// Finally, admin is redirected to the all genres pages
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}

// AdminGetGenreByID renders the genre detail and update form page.
// It mainly handle the get request method.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetGenreByID(w http.ResponseWriter, r *http.Request) {

	// Retrive "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Retrive the genre from db using GetGenreByID interface with id as parameter.
	// If any error occurs, a server error is returned
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a  data map and store the genre model
	data := make(map[string]interface{})
	data["genre"] = genre
	data["base_path"] = base_genres_path

	// Render the "admin-genre-read-update.page.tmpl" template with empty form and data
	render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminGetGenreByID fetch the specific genre from the database in admin interface as well as handles the update for the genre.
// It takes HTTP response writer and request as parameters.
func (m *Repository) PostAdminGetGenreByID(w http.ResponseWriter, r *http.Request) {

	// Retrive "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It parse the form to populate the PostForm.
	// If any error occurs, a server error is returned
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It creates a new form with the post form values.
	form := forms.New(r.PostForm)

	// Then initate a Genre model that stores the updated values.
	update_genre := models.Genre{
		ID:    id,
		Title: r.Form.Get("title"),
	}

	// Before updating the function calls GenreExists interface to check for existing genres.
	exists, err := m.DB.GenreExists(update_genre.Title)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// if exists add a error to the form.
	if exists {
		form.Errors.Add("title", "Genre already exists")
	}

	// Add a validation to the form field.
	form.Required("title")

	// retrive genre using the id
	genre, err := m.DB.GetGenreByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a data map and store the genre model
	data := make(map[string]interface{})
	data["genre"] = genre
	data["base_path"] = base_genres_path

	// If form is invalid render "admin-genre-read-update.page.tmpl" page with form and data
	if !form.Valid() {
		render.Template(w, r, "admin-genre-read-update.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// If form is valid, then call UpdateGenre interface to update the genre.
	// If any error occurs, a server error is returned.
	if err := m.DB.UpdateGenre(&update_genre); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Genre Updated")

	// If successull, then admin is redirected to genre detail page.
	http.Redirect(w, r, fmt.Sprintf("/admin/genres/detail/%d", id), http.StatusSeeOther)
}

// AdminDeleteGenre deletes the genre from the database in admin context.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminDeleteGenre(w http.ResponseWriter, r *http.Request) {

	// Fetch the parameter "id" from the url and parse it into integer.
	// If any error occurs, a server error is returned
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// It calls the DeleteGenre interface with passing id as parameter to delete the record.
	// If any error occurs, a server error is retured.
	if err := m.DB.DeleteGenre(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Genre Deleted")
	// If successfull, admin is redirected to all genres page.
	http.Redirect(w, r, "/admin/genres", http.StatusSeeOther)
}
