package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/handler"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

// ── Mock services ────────────────────────────────────────────────────────────

type mockBankSvc struct {
	stocks []domain.Stock
	err    error
}

func (m *mockBankSvc) SetStocks(_ context.Context, stocks []domain.Stock) error {
	m.stocks = stocks
	return m.err
}
func (m *mockBankSvc) GetStocks(_ context.Context) ([]domain.Stock, error) {
	return m.stocks, m.err
}

var _ service.BankService = (*mockBankSvc)(nil)

type mockWalletSvc struct {
	wallet *domain.Wallet
	qty    int
	err    error
}

func (m *mockWalletSvc) GetWallet(_ context.Context, walletID string) (*domain.Wallet, error) {
	if m.wallet == nil {
		return &domain.Wallet{ID: walletID, Stocks: []domain.Stock{}}, m.err
	}
	return m.wallet, m.err
}
func (m *mockWalletSvc) GetStockQuantity(_ context.Context, _, _ string) (int, error) {
	return m.qty, m.err
}

var _ service.WalletService = (*mockWalletSvc)(nil)

type mockTradeSvc struct{ err error }

func (m *mockTradeSvc) ExecuteTrade(_ context.Context, _, _ string, _ domain.OperationType) error {
	return m.err
}

var _ service.TradeService = (*mockTradeSvc)(nil)

type mockAuditSvc struct{ entries []domain.LogEntry }

func (m *mockAuditSvc) GetAll(_ context.Context) ([]domain.LogEntry, error) {
	return m.entries, nil
}

var _ service.AuditService = (*mockAuditSvc)(nil)

// ── Helpers ──────────────────────────────────────────────────────────────────

func newTestRouter(bankSvc service.BankService, walletSvc service.WalletService, tradeSvc service.TradeService, auditSvc service.AuditService) http.Handler {
	return handler.NewRouter(walletSvc, bankSvc, tradeSvc, auditSvc)
}

func doRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// ── GET /stocks ───────────────────────────────────────────────────────────────

func TestGetStocks_Returns200WithStocks(t *testing.T) {
	bankSvc := &mockBankSvc{stocks: []domain.Stock{{Name: "AAPL", Quantity: 100}}}
	router := newTestRouter(bankSvc, &mockWalletSvc{}, &mockTradeSvc{}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodGet, "/stocks", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var resp map[string][]domain.Stock
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if len(resp["stocks"]) != 1 {
		t.Errorf("expected 1 stock, got %d", len(resp["stocks"]))
	}
}

// ── POST /stocks ──────────────────────────────────────────────────────────────

func TestSetStocks_Returns200(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodPost, "/stocks", map[string]any{
		"stocks": []map[string]any{{"name": "AAPL", "quantity": 100}},
	})
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

// ── GET /wallets/{id} ─────────────────────────────────────────────────────────

func TestGetWallet_Returns200WithWallet(t *testing.T) {
	walletSvc := &mockWalletSvc{wallet: &domain.Wallet{ID: "alice", Stocks: []domain.Stock{{Name: "AAPL", Quantity: 5}}}}
	router := newTestRouter(&mockBankSvc{}, walletSvc, &mockTradeSvc{}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodGet, "/wallets/alice", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var wallet domain.Wallet
	_ = json.NewDecoder(rr.Body).Decode(&wallet)
	if wallet.ID != "alice" {
		t.Errorf("expected wallet alice, got %s", wallet.ID)
	}
}

// ── GET /wallets/{id}/stocks/{name} ──────────────────────────────────────────

func TestGetStock_Returns200WithQuantity(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{qty: 42}, &mockTradeSvc{}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodGet, "/wallets/alice/stocks/AAPL", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var qty int
	_ = json.NewDecoder(rr.Body).Decode(&qty)
	if qty != 42 {
		t.Errorf("expected 42, got %d", qty)
	}
}

// ── POST /wallets/{id}/stocks/{name} ─────────────────────────────────────────

func TestTrade_Buy_Returns200(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodPost, "/wallets/alice/stocks/AAPL", map[string]string{"type": "buy"})
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestTrade_StockNotFound_Returns404(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{err: service.ErrStockNotFound}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodPost, "/wallets/alice/stocks/UNKNOWN", map[string]string{"type": "buy"})
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestTrade_InsufficientBank_Returns400(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{err: service.ErrInsufficientBank}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodPost, "/wallets/alice/stocks/AAPL", map[string]string{"type": "buy"})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestTrade_InsufficientWallet_Returns400(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{err: service.ErrInsufficientWallet}, &mockAuditSvc{})
	rr := doRequest(t, router, http.MethodPost, "/wallets/alice/stocks/AAPL", map[string]string{"type": "sell"})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestTrade_InvalidBody_Returns400(t *testing.T) {
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{}, &mockAuditSvc{})
	req := httptest.NewRequest(http.MethodPost, "/wallets/alice/stocks/AAPL", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

// ── GET /log ──────────────────────────────────────────────────────────────────

func TestGetLog_Returns200WithEntries(t *testing.T) {
	auditSvc := &mockAuditSvc{entries: []domain.LogEntry{
		{Type: domain.Buy, WalletID: "alice", StockName: "AAPL"},
	}}
	router := newTestRouter(&mockBankSvc{}, &mockWalletSvc{}, &mockTradeSvc{}, auditSvc)
	rr := doRequest(t, router, http.MethodGet, "/log", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var resp map[string][]domain.LogEntry
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if len(resp["log"]) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(resp["log"]))
	}
}
