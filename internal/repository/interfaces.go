package repository

import (
	"context"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
)

// BankRepository defines Redis operations for bank stock state.
type BankRepository interface {
	SetStocks(ctx context.Context, stocks []domain.Stock) error
	GetAllStocks(ctx context.Context) ([]domain.Stock, error)
	StockExists(ctx context.Context, name string) (bool, error)
	GetQuantity(ctx context.Context, name string) (int, error)
	IncrBy(ctx context.Context, name string, delta int) error
}

// WalletRepository defines Redis operations for wallet state.
type WalletRepository interface {
	GetAllStocks(ctx context.Context, walletID string) ([]domain.Stock, error)
	GetQuantity(ctx context.Context, walletID, stockName string) (int, error)
	IncrBy(ctx context.Context, walletID, stockName string, delta int) error
}

// AuditRepository defines the append-only audit log operations.
type AuditRepository interface {
	Append(ctx context.Context, entry domain.LogEntry) error
	GetAll(ctx context.Context) ([]domain.LogEntry, error)
}
