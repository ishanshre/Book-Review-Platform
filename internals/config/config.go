package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Global configurations for the application.
// It must not import any internal packages to avoid import cycle
type AppConfig struct {
	InProduction  bool
	UseRedis      bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	MailChan      chan models.MailData
	AdminEmail    string
}
