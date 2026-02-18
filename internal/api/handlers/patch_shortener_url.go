package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ernestosjunior/shortener-url-go/internal/database"

	"github.com/ernestosjunior/shortener-url-go/internal/api/utils"

	"github.com/go-chi/chi/v5"
)

type patchBody struct {
	Url string `json:"url"`
}

type patchResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func PatchShortenerURL(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		log := slog.With(
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"ua", r.UserAgent(),
		)

		var body patchBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&body); err != nil {
			log.Info("body inválido", "err", err)
			utils.SendJSON(w, utils.ResponseJSON{Error: "Body inválido"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.Url); err != nil {
			log.Info("url inválida", "url", body.Url, "err", err)
			utils.SendJSON(w, utils.ResponseJSON{Error: "URL inválida"}, http.StatusBadRequest)
			return
		}

		idStr := chi.URLParam(r, "id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		queries := database.New(db)

		err := queries.UpdateShortURL(ctx, database.UpdateShortURLParams{Url: body.Url, ID: uint64(id)})
		if err != nil {
			log.Info("erro ao atualizar URL", "url", body.Url, "err", err)
			utils.SendJSON(w, utils.ResponseJSON{Error: "Erro ao atualizar os dados"}, http.StatusBadRequest)
			return
		}

		utils.SendJSON(w, utils.ResponseJSON{Data: patchResponse{
			ID:  id,
			URL: body.Url,
		}}, http.StatusCreated)
	}
}
