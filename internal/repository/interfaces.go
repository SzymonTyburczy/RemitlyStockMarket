package repository

import "github.com/szymontyburczy/remitly-stock-market/internal/domain"

// BankRepository defines Redis operations for bank stock state.
type BankRepository interface {
	SetStocks(stocks []domain.Stock) error
	GetAllStocks() ([]domain.Stock, error)
	StockExists(name string) (bool, error)
	// IncrBy atomically changes stock quantity by delta (+1 or -1).
	IncrBy(name string, delta int) error
	GetQuantity(name string) (int, error)
}

// WalletRepository defines Redis operations for wallet state.
type WalletRepository interface {
	GetAllStocks(walletID string) ([]domain.Stock, error)
	GetQuantity(walletID, stockName string) (int, error)
	// IncrBy atomically changes wallet stock quantity by delta.
	IncrBy(walletID, stockName string, delta int) error
}

// AuditRepository defines Redis operations for the append-only audit log.
type AuditRepository interface {
	Append(entry domain.LogEntry) error
	GetAll() ([]domain.LogEntry, error)
}
