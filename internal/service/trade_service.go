package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

type tradeService struct {
	bankRepo   repository.BankRepository
	tradeRepo  repository.TradeRepository
	auditRepo  repository.AuditRepository
}

func NewTradeService(
	bankRepo repository.BankRepository,
	tradeRepo repository.TradeRepository,
	auditRepo repository.AuditRepository,
) TradeService {
	return &tradeService{
		bankRepo:  bankRepo,
		tradeRepo: tradeRepo,
		auditRepo: auditRepo,
	}
}

func (s *tradeService) ExecuteTrade(ctx context.Context, walletID, stockName string, op domain.OperationType) error {
	if op != domain.Buy && op != domain.Sell {
		return ErrInvalidOperation
	}

	// Verify stock exists in bank (registered via POST /stocks) — returns 404 if not
	exists, err := s.bankRepo.StockExists(ctx, stockName)
	if err != nil {
		return fmt.Errorf("checking stock existence: %w", err)
	}
	if !exists {
		return ErrStockNotFound
	}

	// Execute atomic trade via repository Lua scripts
	switch op {
	case domain.Buy:
		err = s.tradeRepo.ExecuteBuy(ctx, walletID, stockName)
	case domain.Sell:
		err = s.tradeRepo.ExecuteSell(ctx, walletID, stockName)
	}

	// Map repository-level errors to service-level sentinel errors
	if errors.Is(err, repository.ErrInsufficientBank) {
		return ErrInsufficientBank
	}
	if errors.Is(err, repository.ErrInsufficientWallet) {
		return ErrInsufficientWallet
	}
	if err != nil {
		return fmt.Errorf("executing trade: %w", err)
	}

	// Append to audit log only on success
	return s.auditRepo.Append(ctx, domain.LogEntry{
		Type:      op,
		WalletID:  walletID,
		StockName: stockName,
	})
}
