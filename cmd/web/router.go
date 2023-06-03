package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
	"github.com/ishanshre/Book-Review-Platform/internals/handler"
)

func router(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(SessionLoad) // load the session middleware
	mux.Use(NoSurf)      // csrf middleware

	// Get route for Home page
	mux.Get("/", handler.Repo.Home)

	// User routes
	mux.Get("/user/login", handler.Repo.Login)
	mux.Post("/user/login", handler.Repo.PostLogin)

	// create a file server with golang path implementation
	fileServer := http.FileServer(http.Dir("./static/"))

	// handler for the file server with system file implementation path
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
