package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ishanshre/Book-Review-Platform/internals/config"
)

func router(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// create a file server with golang path implementation
	fileServer := http.FileServer(http.Dir("./static/"))

	// handler for the file server with system file implementation path
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
