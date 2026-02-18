package api

import (
	"database/sql"
	"net/http"
	"url-shortener/internal/api/handlers"
	"url-shortener/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MyHandler(db *sql.DB, store store.Store) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/{code:[a-zA-Z0-9]{8}}", handlers.GetShortenerURL(db, store))

	r.Route("/api", func(r chi.Router) {
		r.Route("/url", func(r chi.Router) {
			r.Post("/shorten", handlers.PostShortenerURL(db, store))
			r.Patch("/{id:[0-9]+}", handlers.PatchShortenerURL(db))
			r.Delete("/{id:[0-9]+}", handlers.DeleteShortenerURL(db, store))
		})

	})

	return r
}
