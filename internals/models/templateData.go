package models

import "github.com/ishanshre/Book-Review-Platform/internals/forms"

// TemplateData is a struct type that holds data that is passed to go html templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	Username        string
}
