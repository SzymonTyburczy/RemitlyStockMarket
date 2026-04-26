package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

// NewRouter wires all handlers and returns the configured router.
func NewRouter(
	walletSvc service.WalletService,
	bankSvc service.BankService,
	tradeSvc service.TradeService,
	auditSvc service.AuditService,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	wh := NewWalletHandler(walletSvc, tradeSvc)
	sh := NewStockHandler(bankSvc)
	lh := NewLogHandler(auditSvc)
	ch := NewChaosHandler()

	// Wallet endpoints
	r.Get("/wallets/{wallet_id}", wh.GetWallet)
	r.Get("/wallets/{wallet_id}/stocks/{stock_name}", wh.GetStock)
	r.Post("/wallets/{wallet_id}/stocks/{stock_name}", wh.Trade)

	// Bank (stock) endpoints
	r.Get("/stocks", sh.GetStocks)
	r.Post("/stocks", sh.SetStocks)

	// Audit log
	r.Get("/log", lh.GetLog)

	// Chaos
	r.Post("/chaos", ch.Kill)

	return r
}
