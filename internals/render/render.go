package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// app store the pointer to global app config
var app *config.AppConfig

var pathToTemplate = "templates"

// This functions assign global app config to app in render package from main package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefualtData returns default template data to every templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
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
		ts, err := template.New(name).ParseFiles(page)
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
