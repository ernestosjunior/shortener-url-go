package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ernestosjunior/shortener-url-go/internal/store"

	"github.com/ernestosjunior/shortener-url-go/internal/database"

	"github.com/ernestosjunior/shortener-url-go/internal/api/utils"

	"github.com/go-sql-driver/mysql"
)

type postBody struct {
	Url string `json:"url"`
}

type postResponse struct {
	ID   int64  `json:"id,omitempty"`
	Code string `json:"code"`
}

func PostShortenerURL(db *sql.DB, s store.Store) http.HandlerFunc {
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

		var body postBody
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

		const maxAttempts = 3
		code := utils.GenCode()

		queries := database.New(db)

		for range maxAttempts {
			res, err := queries.CreateShortURL(ctx, database.CreateShortURLParams{Code: code, Url: body.Url})

			if err == nil {
				log.Info("url encurtada criada", "code", code, "url", body.Url)

				s.SetShortURLCache(ctx, code, body.Url)

				id, err := res.LastInsertId()
				if err != nil {
					utils.SendJSON(w, utils.ResponseJSON{Data: postResponse{
						Code: code,
					}}, http.StatusCreated)

				} else {
					utils.SendJSON(w, utils.ResponseJSON{Data: postResponse{
						ID:   id,
						Code: code,
					}}, http.StatusCreated)
				}

				return
			}

			var me *mysql.MySQLError
			if errors.As(err, &me) && me.Number == 1062 {
				log.Warn("colisão de code, retry", "code", code)
				continue
			}

			log.Error("erro ao salvar", "err", err)
			utils.SendJSON(w, utils.ResponseJSON{Error: "Erro ao salvar"}, http.StatusInternalServerError)
			return
		}

		log.Error("falha ao gerar code único", "url", body.Url)
		utils.SendJSON(w, utils.ResponseJSON{Error: "Não foi possível gerar o código. Tente novamente."}, http.StatusInternalServerError)
	}
}
