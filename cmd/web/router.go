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

	mux.Get("/user/reset-password", handler.Repo.ResetPassword)
	mux.Post("/user/reset-password", handler.Repo.PostResetPassword)
	mux.Get("/user/reset", handler.Repo.ResetPasswordChange)
	mux.Post("/user/reset", handler.Repo.PostResetPasswordChange)

	// Register routes
	mux.Get("/user/register", handler.Repo.Register)
	mux.Post("/user/register", handler.Repo.PostRegister)

	// create a file server with golang path implementation
	fileServer := http.FileServer(http.Dir("./static/"))

	// handler for the file server with system file implementation path
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	fileServerPublic := http.FileServer(http.Dir("./public/"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServerPublic))

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
		mux.Post("/users/detail/{id}", handler.Repo.AdminUpdateUser)
		mux.Post("/users/detail/{id}/profile", handler.Repo.PostAdminUserProfileUpdate)

		mux.Get("/users/create", handler.Repo.AdminUserAdd)
		mux.Post("/users/create", handler.Repo.PostAdminUserAdd)

		mux.Post("/users/detail/{id}/delete", handler.Repo.PostAdminUserDeleteByID)

		// admin genre router
		mux.Get("/genres", handler.Repo.AdminAllGenres)
		mux.Post("/genres", handler.Repo.PostAdminAddGenre)
		mux.Get("/genres/detail/{id}", handler.Repo.AdminGetGenreByID)
		mux.Post("/genres/detail/{id}", handler.Repo.PostAdminGetGenreByID)
		mux.Post("/genres/detail/{id}/delete", handler.Repo.AdminDeleteGenre)

		// admin publisher router
		mux.Get("/publishers", handler.Repo.AdminAllPublusher)
		mux.Get("/publishers/detail/{id}", handler.Repo.AdminGetPublisherDetailByID)
		mux.Post("/publishers/detail/{id}/update", handler.Repo.PostAdminUpdatePublisher)
		mux.Post("/publishers/detail/{id}/delete", handler.Repo.PostAdminDeletePublisher)
		mux.Get("/publishers/create", handler.Repo.AdminInsertPublisher)
		mux.Post("/publishers/create", handler.Repo.PostAdminInsertPublisher)

		// admin author router
		mux.Get("/authors", handler.Repo.AdminAllAuthor)
		mux.Post("/authors/detail/{id}/delete", handler.Repo.PostAdminDeleteAuthor)
		mux.Get("/authors/detail/{id}", handler.Repo.AdminGetAuthorDetailByID)
		mux.Post("/authors/detail/{id}/update", handler.Repo.PostAdminUpdateAuthor)
		mux.Get("/authors/create", handler.Repo.AdminInsertAuthor)
		mux.Post("/authors/create", handler.Repo.PostAdminInsertAuthor)

		// admin language router
		mux.Get("/languages", handler.Repo.AdminAllLanguage)
		mux.Post("/languages/detail/{id}/delete", handler.Repo.PostAdminDeleteLanguage)
		mux.Post("/languages/detail/{id}/update", handler.Repo.PostAdminUpdateLanguage)
		mux.Post("/languages/create", handler.Repo.PostAdminInsertLanguage)

		// admin book router
		mux.Get("/books", handler.Repo.AdminAllBook)
		mux.Post("/books/detail/{id}/delete", handler.Repo.PostAdminDeleteBook)
		mux.Get("/books/detail/{id}", handler.Repo.AdminGetBookDetailByID)
		mux.Get("/books/create", handler.Repo.AdminInsertBook)
		mux.Post("/books/create", handler.Repo.PostAdminInsertBook)
		mux.Post("/books/detail/{id}/update", handler.Repo.PostAdminUpdateBook)

		// book-admin router
		mux.Get("/bookAuthors", handler.Repo.AdminAllBookAuthor)
		mux.Post("/bookAuthors/create", handler.Repo.PostAdminInsertBookAuthor)
		mux.Get("/bookAuthors/detail/{book_id}/{author_id}", handler.Repo.AdminGetBookAuthorByID)
		mux.Post("/bookAuthors/detail/{book_id}/{author_id}/delete", handler.Repo.PostAdminDeleteBookAuthor)
		mux.Post("/bookAuthors/detail/{book_id}/{author_id}/update", handler.Repo.PostAdminUpdateBookAuthor)
	})
	return mux
}
