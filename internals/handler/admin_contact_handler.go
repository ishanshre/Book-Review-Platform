package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

	// render the detail page with form and data
	render.Template(w, r, "admin-contactdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// ContactUs renders the contacts us page and contact us form.
// It takes HTTP request and response writer as paramters
func (m *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	data := make(map[string]interface{})
	data["contact"] = contact
	render.Template(w, r, "contact-us.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// ContactUs handles the post method after user submits the contact form.
// It takes HTTP request and response writer as paramters
func (m *Repository) PostContactUs(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// initialize a new form using Postform
	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	contact := models.Contact{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Subject:   r.Form.Get("subject"),
		Message:   r.Form.Get("message"),
	}
	form.Required("first_name", "last_name", "email", "phone", "subject", "message")
	form.ValidatePhone("phone")
	data["contact"] = contact
	if !form.Valid() {
		render.Template(w, r, "contact-us.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	contact.SubmittedAt = time.Now()
	contact.IpAddress = r.RemoteAddr
	contact.BrowserInfo = r.UserAgent()
	contact.ReferringPage = r.Referer()
	if err := m.DB.InsertContact(&contact); err != nil {
		helpers.ServerError(w, err)
		return
	}
	msg := models.MailData{
		To:      contact.Email,
		From:    "admin@bookworm.com",
		Subject: fmt.Sprintf("Contact Notification: %v", contact.Subject),
		Content: fmt.Sprintf("%v %v contacted the company. \n %v", contact.FirstName, contact.LastName, contact.Message),
	}
	m.App.MailChan <- msg

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Message Successfull Sent")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
