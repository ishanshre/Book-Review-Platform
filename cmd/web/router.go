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

	// create a file server with golang path implementation
	fileServer := http.FileServer(http.Dir("./static/"))

	// Get route for Home page
	mux.Get("/", handler.Repo.Home)

	// handler for the file server with system file implementation path
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
