package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

// WalletHandler handles all /wallets/* endpoints.
type WalletHandler struct {
	walletSvc service.WalletService
	tradeSvc  service.TradeService
}

func NewWalletHandler(walletSvc service.WalletService, tradeSvc service.TradeService) *WalletHandler {
	return &WalletHandler{walletSvc: walletSvc, tradeSvc: tradeSvc}
}

// GetWallet handles GET /wallets/{wallet_id}
func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	wallet, err := h.walletSvc.GetWallet(r.Context(), walletID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, wallet)
}

// GetStock handles GET /wallets/{wallet_id}/stocks/{stock_name}
func (h *WalletHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	stockName := chi.URLParam(r, "stock_name")
	qty, err := h.walletSvc.GetStockQuantity(r.Context(), walletID, stockName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, qty)
}

// Trade handles POST /wallets/{wallet_id}/stocks/{stock_name}
func (h *WalletHandler) Trade(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "wallet_id")
	stockName := chi.URLParam(r, "stock_name")

	var req domain.TradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.tradeSvc.ExecuteTrade(r.Context(), walletID, stockName, req.Type)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStockNotFound):
			respondError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrInsufficientBank),
			errors.Is(err, service.ErrInsufficientWallet),
			errors.Is(err, service.ErrInvalidOperation):
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
