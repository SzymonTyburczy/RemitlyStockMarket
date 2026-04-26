package handler

import (
	"log/slog"
	"net/http"
	"os"
)

// ChaosHandler handles POST /chaos — kills this instance.
type ChaosHandler struct{}

func NewChaosHandler() *ChaosHandler {
	return &ChaosHandler{}
}

// Kill handles POST /chaos
func (h *ChaosHandler) Kill(w http.ResponseWriter, r *http.Request) {
	slog.Info("chaos endpoint triggered — instance shutting down")
	w.WriteHeader(http.StatusOK)
	// Flush response before exiting so the caller gets a response
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	os.Exit(1)
}
