package handler

import (
	"net/http"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

// LogHandler handles GET /log.
type LogHandler struct {
	auditSvc service.AuditService
}

func NewLogHandler(auditSvc service.AuditService) *LogHandler {
	return &LogHandler{auditSvc: auditSvc}
}

// GetLog handles GET /log
func (h *LogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	entries, err := h.auditSvc.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string][]domain.LogEntry{"log": entries})
}
