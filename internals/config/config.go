package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// Global configurations for the application.
// It must not import any internal packages to avoid import cycle
type AppConfig struct {
	InProduction  bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
}
