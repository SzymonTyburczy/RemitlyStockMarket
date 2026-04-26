package service

import "github.com/szymontyburczy/remitly-stock-market/internal/domain"

// BankService defines the contract for bank (stock supply) operations.
type BankService interface {
	SetStocks(stocks []domain.Stock) error
	GetStocks() ([]domain.Stock, error)
	StockExists(name string) (bool, error)
	// Buy decrements bank supply by 1; returns error if stock unavailable.
	Buy(stockName string) error
	// Sell increments bank supply by 1.
	Sell(stockName string) error
}

// WalletService defines the contract for wallet operations.
type WalletService interface {
	GetWallet(walletID string) (*domain.Wallet, error)
	GetStockQuantity(walletID, stockName string) (int, error)
	// AddStock increments stock in wallet by 1.
	AddStock(walletID, stockName string) error
	// RemoveStock decrements stock in wallet by 1; returns error if none held.
	RemoveStock(walletID, stockName string) error
}

// AuditService defines the contract for audit log operations.
type AuditService interface {
	Append(entry domain.LogEntry) error
	GetAll() ([]domain.LogEntry, error)
}
