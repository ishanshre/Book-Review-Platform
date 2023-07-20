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

// ContactUs renders the contacts us page and contact us form.
// It takes HTTP request and response writer as paramters
func (m *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	data := make(map[string]interface{})
	data["contact"] = contact
	data["base_path"] = base_contacts_path
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
	data["base_path"] = base_contacts_path

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
