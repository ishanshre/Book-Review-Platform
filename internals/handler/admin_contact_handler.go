package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/forms"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/ishanshre/Book-Review-Platform/internals/render"
)

// AdminAllContacts fetches all the record in Contacts.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminAllContacts(w http.ResponseWriter, r *http.Request) {

	// Get all the contacts from the database
	contacts, err := m.DB.AllContacts()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["contacts"] = contacts
	data["base_path"] = base_contacts_path
	render.Template(w, r, "admin-allcontacts.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostAdminDeleteContact Handles the post method for deleting Contact record.
// It takes HTTP response writer and request as paramters
func (m *Repository) PostAdminDeleteContact(w http.ResponseWriter, r *http.Request) {

	// Parsing the contact id from the url.
	// If any error occurs, a server error is returned.
	contact_id, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// DeleteContact interface is used to deleting the record.
	if err := m.DB.DeleteContact(contact_id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Contact record deleted")

	http.Redirect(w, r, "/admin/contacts", http.StatusSeeOther)
}

// AdminGetContactByID handes the detail logic for Contact table.
// It takes HTTP response writer and request as parameters.
func (m *Repository) AdminGetContactByID(w http.ResponseWriter, r *http.Request) {

	// Retrive contact id from the url.
	// Parse them into integer.
	// Return a server error if any error occurs while parsing them
	contact_id, err := strconv.Atoi(chi.URLParam(r, "contact_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the contact using contact id
	contact, err := m.DB.GetContactByID(contact_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create a data map that stores contact
	data := make(map[string]interface{})
	data["contact"] = contact
	data["base_path"] = base_contacts_path

	// render the detail page with form and data
	render.Template(w, r, "admin-contactdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
