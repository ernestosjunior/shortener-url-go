package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ResponseJSON struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, resp ResponseJSON, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)

	if err != nil {
		slog.Error("erro ao fazer marshal da resposta", "error", err)
		SendJSON(w, ResponseJSON{Error: "Something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		slog.Error("erro ao escrever resposta em JSON", "error", err)
		return
	}
}
