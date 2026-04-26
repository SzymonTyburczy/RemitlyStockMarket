// Integration tests for TradeService — require a real Redis instance.
// Set REDIS_TEST_URL=localhost:6379 to run; skipped otherwise.
package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

func redisClientForTest(t *testing.T) *redis.Client {
	t.Helper()
	addr := os.Getenv("REDIS_TEST_URL")
	if addr == "" {
		t.Skip("REDIS_TEST_URL not set — skipping integration test")
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("cannot connect to Redis at %s: %v", addr, err)
	}
	t.Cleanup(func() {
		client.FlushAll(context.Background())
		client.Close()
	})
	return client
}

func buildTradeService(t *testing.T, client *redis.Client) service.TradeService {
	t.Helper()
	bankRepo := repository.NewBankRepository(client)
	tradeRepo := repository.NewTradeRepository(client)
	auditRepo := repository.NewAuditRepository(client)
	return service.NewTradeService(bankRepo, tradeRepo, auditRepo)
}

func TestTradeService_Buy_Success(t *testing.T) {
	client := redisClientForTest(t)
	ctx := context.Background()
	bankRepo := repository.NewBankRepository(client)
	_ = bankRepo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})

	tradeSvc := buildTradeService(t, client)
	if err := tradeSvc.ExecuteTrade(ctx, "alice", "AAPL", domain.Buy); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	qty, _ := bankRepo.GetQuantity(ctx, "AAPL")
	if qty != 9 {
		t.Errorf("bank should have 9 AAPL, got %d", qty)
	}
	walletRepo := repository.NewWalletRepository(client)
	wQty, _ := walletRepo.GetQuantity(ctx, "alice", "AAPL")
	if wQty != 1 {
		t.Errorf("wallet should have 1 AAPL, got %d", wQty)
	}
}

func TestTradeService_Buy_StockNotFound_Returns404Error(t *testing.T) {
	client := redisClientForTest(t)
	tradeSvc := buildTradeService(t, client)
	err := tradeSvc.ExecuteTrade(context.Background(), "alice", "UNKNOWN", domain.Buy)
	if !errors.Is(err, service.ErrStockNotFound) {
		t.Errorf("expected ErrStockNotFound, got %v", err)
	}
}

func TestTradeService_Buy_BankEmpty_Returns400Error(t *testing.T) {
	client := redisClientForTest(t)
	ctx := context.Background()
	bankRepo := repository.NewBankRepository(client)
	_ = bankRepo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 0}})

	tradeSvc := buildTradeService(t, client)
	err := tradeSvc.ExecuteTrade(ctx, "alice", "AAPL", domain.Buy)
	if !errors.Is(err, service.ErrInsufficientBank) {
		t.Errorf("expected ErrInsufficientBank, got %v", err)
	}
}

func TestTradeService_Sell_Success(t *testing.T) {
	client := redisClientForTest(t)
	ctx := context.Background()
	bankRepo := repository.NewBankRepository(client)
	walletRepo := repository.NewWalletRepository(client)
	_ = bankRepo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 5}})
	_ = walletRepo.IncrBy(ctx, "alice", "AAPL", 3)

	tradeSvc := buildTradeService(t, client)
	if err := tradeSvc.ExecuteTrade(ctx, "alice", "AAPL", domain.Sell); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	qty, _ := bankRepo.GetQuantity(ctx, "AAPL")
	if qty != 6 {
		t.Errorf("bank should have 6 AAPL, got %d", qty)
	}
}

func TestTradeService_Sell_WalletEmpty_Returns400Error(t *testing.T) {
	client := redisClientForTest(t)
	ctx := context.Background()
	bankRepo := repository.NewBankRepository(client)
	_ = bankRepo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})

	tradeSvc := buildTradeService(t, client)
	err := tradeSvc.ExecuteTrade(ctx, "alice", "AAPL", domain.Sell)
	if !errors.Is(err, service.ErrInsufficientWallet) {
		t.Errorf("expected ErrInsufficientWallet, got %v", err)
	}
}

func TestTradeService_AuditLogAppended_OnSuccess(t *testing.T) {
	client := redisClientForTest(t)
	ctx := context.Background()
	bankRepo := repository.NewBankRepository(client)
	_ = bankRepo.SetStocks(ctx, []domain.Stock{{Name: "AAPL", Quantity: 10}})

	tradeSvc := buildTradeService(t, client)
	_ = tradeSvc.ExecuteTrade(ctx, "alice", "AAPL", domain.Buy)

	auditRepo := repository.NewAuditRepository(client)
	entries, _ := auditRepo.GetAll(ctx)
	if len(entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].Type != domain.Buy || entries[0].WalletID != "alice" || entries[0].StockName != "AAPL" {
		t.Errorf("unexpected audit entry: %+v", entries[0])
	}
}

func TestTradeService_InvalidOperation(t *testing.T) {
	client := redisClientForTest(t)
	tradeSvc := buildTradeService(t, client)
	err := tradeSvc.ExecuteTrade(context.Background(), "alice", "AAPL", "UNKNOWN")
	if !errors.Is(err, service.ErrInvalidOperation) {
		t.Errorf("expected ErrInvalidOperation, got %v", err)
	}
}
