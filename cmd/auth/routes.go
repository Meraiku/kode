package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.logRequest)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/users", app.handleGetUsers)
		r.Post("/users", app.handlePostUsers)

		r.Route("/", func(r chi.Router) {

			r.Get("/notes", app.authenticateUser(app.handleGetNotes))
			r.Post("/notes", app.authenticateUser(app.handlePostNotes))
		})

		r.Route("/tokens", func(r chi.Router) {
			r.Post("/", app.handleGetTokens)
			r.Post("/refresh", app.handleRefreshTokens)
		})
	})

	return r
}
