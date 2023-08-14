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

// AdminAllPublisher renders admin all publisher page
func (m *Repository) AdminAllPublusher(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["base_path"] = base_publishers_path
	// render the "admin-allpublishers.page.tmpl" with data and new empty form
	render.Template(w, r, "admin-allpublishers.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminAllPublisherFilterApi(w http.ResponseWriter, r *http.Request) {
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
	filteredPublisher, err := m.DB.AllPublishersFilter(limit, page, searchKey, sort)
	if err != nil {
		helpers.ServerError(w, err)
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, filteredPublisher)
}

// PostAdminDeletePublisher handles the POST request to delete a publisher.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The publisher is deleted from the database using the ID.
// If there is an error during the deletion process, a server error is returned.
// Otherwise, the user is redirected to the publishers admin page.
func (m *Repository) PostAdminDeletePublisher(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err := m.DB.DeletePublisher(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Publisher Deleted")
	http.Redirect(w, r, "/admin/publishers", http.StatusSeeOther)
}

// AdminGetPublisherDetailByID handles the GET request to retrieve the details of a publisher by its ID.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The publisher is retrieved from the database using the ID.
// The data map is populated with the publisher object.
// The template is rendered with the form object and data map.
func (m *Repository) AdminGetPublisherDetailByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	publisher, err := m.DB.GetPublisherByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["publisher"] = publisher
	data["base_path"] = base_publishers_path
	render.Template(w, r, "admin-publisherdetail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminUpdatePublisher handles the post method logic for updating a publisher.
// It takes the HTTP response writer and request as parameters.
// The ID of the publisher is extracted from the URL parameter.
// The form is parsed and a new form object is created.
// The established date is converted to an integer.
// The publisher object is populated with the form data.
// The ID of the publisher is set.
// The "name" field is marked as required.
// The data map is populated with the publisher object.
// If the form is not valid, the template is rendered with the form object and data map.
// The publisher is updated in the database.
// The user is redirected to the publisher detail page.
func (m *Repository) PostAdminUpdatePublisher(w http.ResponseWriter, r *http.Request) {
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
	establishedDate, _ := strconv.Atoi(r.Form.Get("established_date"))
	publisher := models.Publisher{
		Name:            r.Form.Get("name"),
		Description:     r.Form.Get("description"),
		Pic:             r.Form.Get("pic"),
		Address:         r.Form.Get("address"),
		Phone:           r.Form.Get("phone"),
		Email:           r.Form.Get("email"),
		Website:         r.Form.Get("website"),
		EstablishedDate: establishedDate,
		Latitude:        r.Form.Get("latitude"),
		Longitude:       r.Form.Get("longitude"),
	}
	publisher.ID = id
	form.Required("name")
	data := make(map[string]interface{})
	data["publisher"] = publisher
	data["base_path"] = base_publishers_path
	if !form.Valid() {
		render.Template(w, r, "admin-publisherdetail.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.UpdatePublisher(&publisher); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Publisher Updated")
	http.Redirect(w, r, fmt.Sprintf("/admin/publishers/detail/%d", id), http.StatusSeeOther)
}

// AdminInsertPublisher handles the logic for rendering the publisher insert form.
// It takes the HTTP response writer and request as parameters.
// An empty publisher object is created.
// The data map is populated with the empty publisher object.
// The template is rendered with the form object and data map.
func (m *Repository) AdminInsertPublisher(w http.ResponseWriter, r *http.Request) {
	var emptyPublisher models.Publisher
	data := make(map[string]interface{})
	data["publisher"] = emptyPublisher
	data["base_path"] = base_publishers_path
	render.Template(w, r, "admin-publisherinsert.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostAdminInsertPublisher handles the logic for inserting a new publisher using the POST method.
// It takes the HTTP response writer and request as parameters.
// The function parses the form data from the request.
// It creates a new form object and validates the required fields.
// The established date is converted to an integer.
// The picture path is obtained by uploading an image using the helper function.
// A new publisher object is created with the form data and the picture path.
// The publisher data is stored in the "publisher" key of the data map.
// If the form is not valid, the template is rendered with the form errors and data.
// If there is an error during insertion, a server error is returned.
// Otherwise, the publisher is inserted into the database and a redirect response is sent to the publishers list page.
func (m *Repository) PostAdminInsertPublisher(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)

	establishedDate, err := strconv.Atoi(r.Form.Get("established_date"))
	if err != nil {
		form.Errors.Add("established_date", "Invalid established date")
	}
	idString := fmt.Sprintf("%d%s", establishedDate, helpers.RandomAlphaNum(8))
	pic_path, err := helpers.AdminPublicUploadImage(r, "pic", "publisher", idString)
	if err != nil {
		form.Errors.Add("pic", "No picture was choosen")
	}
	publisher := &models.Publisher{
		Name:            r.Form.Get("name"),
		Description:     r.Form.Get("description"),
		Pic:             pic_path,
		Address:         r.Form.Get("address"),
		Phone:           r.Form.Get("phone"),
		Email:           r.Form.Get("email"),
		Website:         r.Form.Get("website"),
		EstablishedDate: establishedDate,
		Latitude:        r.Form.Get("latitude"),
		Longitude:       r.Form.Get("longitude"),
	}
	data := make(map[string]interface{})
	data["publisher"] = publisher
	data["base_path"] = base_publishers_path
	if !form.Valid() {
		render.Template(w, r, "admin-publisherinsert.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	if err := m.DB.InsertPublisher(publisher); err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add success message
	m.App.Session.Put(r.Context(), "flash", "Publisher Added")

	http.Redirect(w, r, "/admin/publishers", http.StatusSeeOther)
}
