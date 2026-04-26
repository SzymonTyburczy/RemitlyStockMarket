package service_test

import (
	"context"
	"testing"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

type mockWalletRepo struct{ wallets map[string]map[string]int }

func newMockWalletRepo() *mockWalletRepo {
	return &mockWalletRepo{wallets: make(map[string]map[string]int)}
}
func (m *mockWalletRepo) GetAllStocks(_ context.Context, walletID string) ([]domain.Stock, error) {
	out := make([]domain.Stock, 0)
	for name, qty := range m.wallets[walletID] {
		out = append(out, domain.Stock{Name: name, Quantity: qty})
	}
	return out, nil
}
func (m *mockWalletRepo) GetQuantity(_ context.Context, walletID, stockName string) (int, error) {
	return m.wallets[walletID][stockName], nil
}
func (m *mockWalletRepo) IncrBy(_ context.Context, walletID, stockName string, delta int) error {
	if m.wallets[walletID] == nil {
		m.wallets[walletID] = make(map[string]int)
	}
	m.wallets[walletID][stockName] += delta
	return nil
}

var _ repository.WalletRepository = (*mockWalletRepo)(nil)

func TestWalletService_GetWallet_EmptyWallet(t *testing.T) {
	svc := service.NewWalletService(newMockWalletRepo())
	wallet, err := svc.GetWallet(context.Background(), "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wallet.ID != "alice" {
		t.Errorf("expected ID alice, got %s", wallet.ID)
	}
	if len(wallet.Stocks) != 0 {
		t.Errorf("expected empty stocks, got %v", wallet.Stocks)
	}
}

func TestWalletService_GetStockQuantity(t *testing.T) {
	repo := newMockWalletRepo()
	_ = repo.IncrBy(context.Background(), "alice", "AAPL", 5)
	svc := service.NewWalletService(repo)
	qty, err := svc.GetStockQuantity(context.Background(), "alice", "AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if qty != 5 {
		t.Errorf("expected 5, got %d", qty)
	}
}

func TestWalletService_GetWallet_WithStocks(t *testing.T) {
	repo := newMockWalletRepo()
	ctx := context.Background()
	_ = repo.IncrBy(ctx, "bob", "AAPL", 3)
	_ = repo.IncrBy(ctx, "bob", "GOOG", 7)
	svc := service.NewWalletService(repo)
	wallet, _ := svc.GetWallet(ctx, "bob")
	if len(wallet.Stocks) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(wallet.Stocks))
	}
}
