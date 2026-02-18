package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"url-shortener/internal/api/utils"
	"url-shortener/internal/database"
	"url-shortener/internal/store"

	"github.com/go-chi/chi/v5"
)

func DeleteShortenerURL(db *sql.DB, s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		queries := database.New(db)

		if err := queries.DeleteShortURL(ctx, uint64(id)); err != nil {
			slog.Error("erro ao deletar o link", "err", err)
			utils.SendJSON(w, utils.ResponseJSON{Error: "Erro ao deletar o link"}, http.StatusBadRequest)
			return
		}

		data, err := queries.GetShortURLById(ctx, uint64(id))
		if err == nil {
			s.RemoveShortURLCache(ctx, data.Code)
		}

		utils.SendJSON(w, utils.ResponseJSON{}, http.StatusNoContent)
	}
}
