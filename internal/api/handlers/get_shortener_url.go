package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"url-shortener/internal/database"
	"url-shortener/internal/store"

	"github.com/go-chi/chi/v5"
)

func GetShortenerURL(db *sql.DB, s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		code := chi.URLParam(r, "code")

		log := slog.With(
			"code", code,
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"ua", r.UserAgent(),
		)

		urlCache := s.GetShortURLCache(ctx, code)
		if urlCache != "" {
			log.Info("redirect via cache", "code", code, "to", urlCache)
			http.Redirect(w, r, urlCache, http.StatusPermanentRedirect)
		} else {
			queries := database.New(db)

			url, err := queries.GetShortURL(ctx, code)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Info("url não encontrada")
					http.Error(w, "URL não encontrada", http.StatusNotFound)
					return
				}

				log.Error("erro ao buscar url", "err", err)
				http.Error(w, "Erro ao buscar URL", http.StatusInternalServerError)
				return
			}

			_, err = db.ExecContext(ctx, `UPDATE short_urls SET clicks = clicks + 1 WHERE id = ?`, url.ID)
			if err != nil {
				log.Error("erro ao contabilizar click", "err", err, "id", url.ID)
			}

			s.SetShortURLCache(ctx, code, url.Url)

			log.Info("redirect", "id", url.ID, "to", url.Url)
			http.Redirect(w, r, url.Url, http.StatusPermanentRedirect)
		}
	}
}
