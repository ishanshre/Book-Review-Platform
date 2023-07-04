package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// AdminAllLanguage retrieves all languages from the database and renders the admin-alllanguages page.
// It takes the HTTP response writer and request as parameters.
// The function calls the AllLanguage interface to retrieve all languages from the database.
// If an error occurs during the retrieval, a server error is returned.
// The function creates a data map and adds the retrieved languages to the "languages" key in the data map.
// It creates an empty Language model and adds it to the "language" key in the data map.
// The function renders the "admin-alllanguages.page.tmpl" template with the data map and an empty form.
func (m *Repository) AdminAllLanguage(w http.ResponseWriter, r *http.Request) {
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["languages"] = languages
	var emptyLanguage models.Language
	data["language"] = emptyLanguage
	data["base_path"] = base_languages_path
	render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeleteLanguage deletes a language from the database in the admin context.
// It takes the HTTP response writer and request as parameters.
// The function retrieves the "id" parameter from the URL and converts it to an integer.
// If there is an error during the conversion, a server error is returned.
// The function calls the DeleteLanguage interface to delete the language with the specified ID from the database.
// If an error occurs during the deletion, a server error is returned.
// After successful deletion, the function redirects the user to the "/admin/languages" page.
func (m *Repository) PostAdminDeleteLanguage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeleteLanguage(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Language Deleted")

	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)
}

// PostAdminUpdateLanguage handles the update language logic
func (m *Repository) PostAdminUpdateLanguage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["languages"] = languages
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	language := models.Language{
		ID:       id,
		Language: r.Form.Get("language"),
	}
	form.Required("language")
	data["language"] = language
	exists, err := m.DB.LanguageExists(language.Language)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		form.Errors.Add("language", "This language already exists")
	}

	data["base_path"] = base_languages_path

	if !form.Valid() {
		render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdateLanguage(&language); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Language Updated")

	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)
}

// PostAdminInsertLanguage inserts a new language into the database in the admin context.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// If there is an error during form parsing, a server error is returned.
// The function creates a new form instance and initializes a language model with the language value from the form.
// It sets the "language" field as required in the form.
// The function creates a data map and adds the "add_language" key with the language model as its value.
// It calls the AllLanguage interface to retrieve all existing languages from the database.
// If an error occurs during the retrieval, a server error is returned.
// The function adds the retrieved languages to the data map with the "languages" key.
// It checks if the language already exists in the database using the LanguageExists interface.
// If the language exists, an error is added to the form's errors.
// If the form is not valid, the function renders the "admin-alllanguages.page.tmpl" template with the form errors and data.
// If there are no form errors, the function calls the InsertLanguage method to insert the language into the database.
// If an error occurs during the insertion, a server error is returned.
// After successful insertion, the function redirects the user to the "/admin/languages" page.
func (m *Repository) PostAdminInsertLanguage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)
	language := models.Language{
		Language: r.Form.Get("language"),
	}
	form.Required("language")
	data := make(map[string]interface{})
	data["add_language"] = language
	languages, err := m.DB.AllLanguage()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["languages"] = &languages
	stat, _ := m.DB.LanguageExists(language.Language)
	if stat {
		form.Errors.Add("language", "This language already exists")
	}

	data["base_path"] = base_languages_path

	if !form.Valid() {
		render.Template(w, r, "admin-alllanguages.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertLanguage(&language); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Language Added")

	http.Redirect(w, r, "/admin/languages", http.StatusSeeOther)

}
