package handler

import (
	"encoding/json"
	"net/http"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

// StockHandler handles GET /stocks and POST /stocks.
type StockHandler struct {
	bankSvc service.BankService
}

func NewStockHandler(bankSvc service.BankService) *StockHandler {
	return &StockHandler{bankSvc: bankSvc}
}

// GetStocks handles GET /stocks
func (h *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.bankSvc.GetStocks(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string][]domain.Stock{"stocks": stocks})
}

// SetStocks handles POST /stocks — replaces the entire bank state.
func (h *StockHandler) SetStocks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Stocks []domain.Stock `json:"stocks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.bankSvc.SetStocks(r.Context(), body.Stocks); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
