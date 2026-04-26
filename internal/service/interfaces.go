package service

import (
	"context"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
)

// BankService manages the bank's stock inventory.
type BankService interface {
	SetStocks(ctx context.Context, stocks []domain.Stock) error
	GetStocks(ctx context.Context) ([]domain.Stock, error)
}

// WalletService manages wallet state queries.
type WalletService interface {
	GetWallet(ctx context.Context, walletID string) (*domain.Wallet, error)
	GetStockQuantity(ctx context.Context, walletID, stockName string) (int, error)
}

// TradeService orchestrates atomic buy/sell between bank and wallet.
type TradeService interface {
	ExecuteTrade(ctx context.Context, walletID, stockName string, op domain.OperationType) error
}

// AuditService provides read access to the audit log.
type AuditService interface {
	GetAll(ctx context.Context) ([]domain.LogEntry, error)
}
