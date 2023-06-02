package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// NoSurf implement csrf token middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next) // creates a new handler
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode, // allows cookies to sent in cross site
	})
	return csrfHandler
}
