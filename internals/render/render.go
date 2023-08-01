package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"github.com/justinas/nosurf"
)

// app store the pointer to global app config
var app *config.AppConfig

var pathToTemplate = "./templates"

// This functions assign global app config to app in render package from main package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefualtData returns default template data to every templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)

	// checkes if user_id exists in session then set true to IsAuthenticated
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
		td.Username = app.Session.GetString(r.Context(), "username")
		td.UserID = app.Session.Get(r.Context(), "user_id").(int)
		td.AccessLevel = app.Session.GetInt(r.Context(), "access_level")
	}
	td.Flash = app.Session.GetString(r.Context(), "flash")
	td.Error = app.Session.GetString(r.Context(), "error")
	td.Warning = app.Session.GetString(r.Context(), "warning")
	return td
}

// Template renders the template using http/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	// render template cache from template if UseCache is true in global configuration
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get a request from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println("Could not get the template from the template cache")
		return errors.New("cannot get the template cache from the cache")
	}

	// create a new buffer to store the templates and data to pass to template
	buff := new(bytes.Buffer)

	// add default template to all templates
	td = AddDefaultData(td, r)

	// add the parsed template and data to buffer
	if err := t.Execute(buff, td); err != nil {
		log.Println(err)
		return err
	}

	// redner the template using buffer.WriteTo
	_, err := buff.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache
func CreateTemplateCache() (map[string]*template.Template, error) {

	// Defineing a custom function map for function to run in go templates
	funcMap := template.FuncMap{
		"TimeSince": TimeSince,
		"DateOnly":  DateOnly,
	}

	// path pattern to layout and pages template.
	pathLayoutPattern := filepath.Join(pathToTemplate, "layout", "*.layout.tmpl")
	pathPagePattern := filepath.Join(pathToTemplate, "page", "*.page.tmpl")

	// myCache is an empty cache using map.
	myCache := map[string]*template.Template{}

	// pages is a slice of string.
	// It stores all the name of all files matching the pattern with its relative path.
	// i.e. template/page/home.page.tmpl
	pages, err := filepath.Glob(pathPagePattern)
	if err != nil {
		return myCache, err
	}

	// loop through all the pages and add base template to each pages
	for _, page := range pages {
		name := filepath.Base(page) // filepath.Base name returns file name with its extension

		// create and parse new template
		ts, err := template.New(name).Funcs(funcMap).ParseFiles(page)
		if err != nil {
			return myCache, fmt.Errorf("error in parsing template %s", err)
		}

		// find the base template using filepath.Glob and pattern
		matches, err := filepath.Glob(pathLayoutPattern)
		if err != nil {
			return myCache, err
		}

		// if found
		if len(matches) > 0 {
			// add layout templates to the page templates
			ts, err = ts.ParseGlob(pathLayoutPattern)
			if err != nil {
				return myCache, err
			}
		}

		// assign new template to cache
		myCache[name] = ts
	}
	return myCache, nil
}

func TimeSince(t time.Time) string {
	duration := time.Since(t)
	// return fmt.Sprintf("%s ago", duration.Round(time.Second))

	// more readable format
	years := int(duration.Hours() / 24 / 365)
	months := int(duration.Hours() / 24 / 30)
	weeks := int(duration.Hours() / 24 / 7)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())
	minutes := int(duration.Hours() * 60)
	seconds := int(duration.Hours() * 60 * 60)
	var result string
	switch {
	case years > 0:
		result = fmt.Sprintf("%d years ago", years)
	case months > 0:
		result = fmt.Sprintf("%d months ago", months)
	case weeks > 0:
		result = fmt.Sprintf("%d weeks ago", weeks)
	case days > 0:
		result = fmt.Sprintf("%d days ago", days)
	case hours > 0:
		result = fmt.Sprintf("%d hours ago", hours)
	case minutes > 0:
		result = fmt.Sprintf("%d minutes ago", minutes)
	default:
		result = fmt.Sprintf("%d seconds ago", seconds)
	}
	return result
}

func DateOnly(t time.Time) string {
	return t.Format(time.DateOnly)
}
