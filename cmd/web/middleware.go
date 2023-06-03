package main

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	})
}
