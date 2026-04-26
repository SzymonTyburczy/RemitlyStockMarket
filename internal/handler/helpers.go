package handler

import (
	"encoding/json"
	"net/http"
)

// respondJSON writes v as JSON with the given status code.
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// respondError writes a plain-text error message with the given status code.
func respondError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
