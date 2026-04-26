package service_test

import (
	"context"
	"testing"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

type mockBankRepo struct{ stocks map[string]int }

func newMockBankRepo() *mockBankRepo { return &mockBankRepo{stocks: make(map[string]int)} }

func (m *mockBankRepo) SetStocks(_ context.Context, stocks []domain.Stock) error {
	m.stocks = make(map[string]int)
	for _, s := range stocks {
		m.stocks[s.Name] = s.Quantity
	}
	return nil
}
func (m *mockBankRepo) GetAllStocks(_ context.Context) ([]domain.Stock, error) {
	out := make([]domain.Stock, 0, len(m.stocks))
	for name, qty := range m.stocks {
		out = append(out, domain.Stock{Name: name, Quantity: qty})
	}
	return out, nil
}
func (m *mockBankRepo) StockExists(_ context.Context, name string) (bool, error) {
	_, ok := m.stocks[name]
	return ok, nil
}
func (m *mockBankRepo) GetQuantity(_ context.Context, name string) (int, error) {
	return m.stocks[name], nil
}
func (m *mockBankRepo) IncrBy(_ context.Context, name string, delta int) error {
	m.stocks[name] += delta
	return nil
}

var _ repository.BankRepository = (*mockBankRepo)(nil)

func TestBankService_SetStocks_PopulatesState(t *testing.T) {
	svc := service.NewBankService(newMockBankRepo())
	ctx := context.Background()
	if err := svc.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 100}, {Name: "GOOG", Quantity: 50}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stocks, _ := svc.GetStocks(ctx)
	if len(stocks) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(stocks))
	}
}

func TestBankService_SetStocks_ReplacesExistingState(t *testing.T) {
	svc := service.NewBankService(newMockBankRepo())
	ctx := context.Background()
	_ = svc.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})
	_ = svc.SetStocks(ctx, []domain.Stock{{Name: "GOOG", Quantity: 5}})
	stocks, _ := svc.GetStocks(ctx)
	if len(stocks) != 1 || stocks[0].Name != "GOOG" {
		t.Errorf("expected only GOOG after replace, got %v", stocks)
	}
}

func TestBankService_GetStocks_EmptyInitially(t *testing.T) {
	svc := service.NewBankService(newMockBankRepo())
	stocks, err := svc.GetStocks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stocks) != 0 {
		t.Errorf("expected empty bank, got %d stocks", len(stocks))
	}
}
