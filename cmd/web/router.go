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

	// Login routes
	mux.Get("/user/login", handler.Repo.Login)
	mux.Post("/user/login", handler.Repo.PostLogin)
	mux.Get("/user/logout", handler.Repo.Logout)

	// Register routes
	mux.Get("/user/register", handler.Repo.Register)
	mux.Post("/user/register", handler.Repo.PostRegister)

	// create a file server with golang path implementation
	fileServer := http.FileServer(http.Dir("./static/"))

	// handler for the file server with system file implementation path
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/profile", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handler.Repo.PersonalProfile)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Use(Admin)
		mux.Get("/", handler.Repo.AdminDashboard)
		mux.Get("/users", handler.Repo.AdminAllUsers)
		mux.Get("/users/detail/{id}", handler.Repo.AdminGetUserDetailByID)

		mux.Get("/users/create", handler.Repo.AdminUserAdd)
		mux.Post("/users/create", handler.Repo.PostAdminUserAdd)

		mux.Post("/users/detail/{id}/delete", handler.Repo.PostAdminUserDeleteByID)
	})
	return mux
}
