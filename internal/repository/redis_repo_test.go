// Integration tests for Redis repositories.
// Set REDIS_TEST_URL=localhost:6379 to run; skipped otherwise.
package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

func redisClient(t *testing.T) *redis.Client {
	t.Helper()
	addr := os.Getenv("REDIS_TEST_URL")
	if addr == "" {
		t.Skip("REDIS_TEST_URL not set — skipping integration test")
	}
	c := redis.NewClient(&redis.Options{Addr: addr})
	if err := c.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("cannot connect to Redis at %s: %v", addr, err)
	}
	t.Cleanup(func() {
		c.FlushAll(context.Background())
		c.Close()
	})
	return c
}

func TestBankRepo_SetAndGetAllStocks(t *testing.T) {
	repo := repository.NewBankRepository(redisClient(t))
	ctx := context.Background()
	input := []domain.Stock{{Name: "AAPL", Quantity: 100}, {Name: "GOOG", Quantity: 50}}
	if err := repo.SetStocks(ctx, input); err != nil {
		t.Fatalf("SetStocks() error: %v", err)
	}
	stocks, err := repo.GetAllStocks(ctx)
	if err != nil {
		t.Fatalf("GetAllStocks() error: %v", err)
	}
	if len(stocks) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(stocks))
	}
}

func TestBankRepo_StockExists(t *testing.T) {
	repo := repository.NewBankRepository(redisClient(t))
	ctx := context.Background()
	_ = repo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})
	ok, err := repo.StockExists(ctx, "AAPL")
	if err != nil || !ok {
		t.Errorf("expected AAPL to exist: err=%v ok=%v", err, ok)
	}
	ok, _ = repo.StockExists(ctx, "UNKNOWN")
	if ok {
		t.Error("expected UNKNOWN to not exist")
	}
}

func TestBankRepo_IncrBy(t *testing.T) {
	repo := repository.NewBankRepository(redisClient(t))
	ctx := context.Background()
	_ = repo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})
	_ = repo.IncrBy(ctx, "AAPL", -3)
	qty, _ := repo.GetQuantity(ctx, "AAPL")
	if qty != 7 {
		t.Errorf("expected 7, got %d", qty)
	}
}

func TestWalletRepo_IncrByAndGetQuantity(t *testing.T) {
	repo := repository.NewWalletRepository(redisClient(t))
	ctx := context.Background()
	_ = repo.IncrBy(ctx, "alice", "AAPL", 5)
	_ = repo.IncrBy(ctx, "alice", "AAPL", 3)
	qty, err := repo.GetQuantity(ctx, "alice", "AAPL")
	if err != nil {
		t.Fatalf("GetQuantity() error: %v", err)
	}
	if qty != 8 {
		t.Errorf("expected 8, got %d", qty)
	}
}

func TestWalletRepo_GetAllStocks(t *testing.T) {
	repo := repository.NewWalletRepository(redisClient(t))
	ctx := context.Background()
	_ = repo.IncrBy(ctx, "bob", "AAPL", 2)
	_ = repo.IncrBy(ctx, "bob", "GOOG", 4)
	stocks, err := repo.GetAllStocks(ctx, "bob")
	if err != nil {
		t.Fatalf("GetAllStocks() error: %v", err)
	}
	if len(stocks) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(stocks))
	}
}

func TestAuditRepo_AppendAndGetAll(t *testing.T) {
	repo := repository.NewAuditRepository(redisClient(t))
	ctx := context.Background()
	entries := []domain.LogEntry{
		{Type: domain.Buy, WalletID: "alice", StockName: "AAPL"},
		{Type: domain.Sell, WalletID: "bob", StockName: "GOOG"},
	}
	for _, e := range entries {
		if err := repo.Append(ctx, e); err != nil {
			t.Fatalf("Append() error: %v", err)
		}
	}
	got, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Type != domain.Buy || got[1].Type != domain.Sell {
		t.Error("entries out of order or wrong type")
	}
}
