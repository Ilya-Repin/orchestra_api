package handler

import (
	"encoding/json"
	"github.com/Ilya-Repin/orchestra_api/internal/openapi"
	"net/http"
)

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, openapi.ErrorResponse{Message: &message})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
	}
}
