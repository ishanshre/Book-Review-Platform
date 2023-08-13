package middleware

import (
	"net/http"

	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewMiddlewareApp(a *config.AppConfig) {
	app = a
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
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

// Auth is a middleware function that checks if the user is authenticated.
// If the user is not authenticated, it redirects to the login page.
// It takes a next http.Handler as an argument and returns an http.Handler.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuthRedirect restrict authenticated user to access some page such as login, signup, forget passsword, etc.
func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func KycValidated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsValidated(r) {
			app.Session.Put(r.Context(), "warning", "Please Update Kyc")
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Admin is a middleware function that checks if the user is an admin.
// If the user is not an admin, it redirects to the home page.
// It takes a next http.Handler as an argument and returns an http.Handler.
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAdmin(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("| %s | %s | %s ", r.Method, r.URL.Path, r.Proto)
		next.ServeHTTP(w, r)
	})
}
