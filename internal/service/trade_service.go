package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

// buyScript atomically:
//  1. Checks bank[stock] > 0  → error INSUFFICIENT_BANK
//  2. Decrements bank[stock] by 1
//  3. Increments wallet[stock] by 1
//
// KEYS[1]=bankKey  KEYS[2]=walletKey  ARGV[1]=stockName
var buyScript = redis.NewScript(`
local qty = tonumber(redis.call('HGET', KEYS[1], ARGV[1]))
if qty == nil or qty <= 0 then
  return redis.error_reply('INSUFFICIENT_BANK')
end
redis.call('HINCRBY', KEYS[1], ARGV[1], -1)
redis.call('HINCRBY', KEYS[2], ARGV[1], 1)
return 1
`)

// sellScript atomically:
//  1. Checks wallet[stock] > 0  → error INSUFFICIENT_WALLET
//  2. Decrements wallet[stock] by 1
//  3. Increments bank[stock] by 1
//
// KEYS[1]=walletKey  KEYS[2]=bankKey  ARGV[1]=stockName
var sellScript = redis.NewScript(`
local qty = tonumber(redis.call('HGET', KEYS[1], ARGV[1]))
if qty == nil or qty <= 0 then
  return redis.error_reply('INSUFFICIENT_WALLET')
end
redis.call('HINCRBY', KEYS[1], ARGV[1], -1)
redis.call('HINCRBY', KEYS[2], ARGV[1], 1)
return 1
`)

type tradeService struct {
	bankRepo   repository.BankRepository
	walletRepo repository.WalletRepository
	auditRepo  repository.AuditRepository
	redisClient *redis.Client
}

func NewTradeService(
	bankRepo repository.BankRepository,
	walletRepo repository.WalletRepository,
	auditRepo repository.AuditRepository,
	redisClient *redis.Client,
) TradeService {
	return &tradeService{
		bankRepo:    bankRepo,
		walletRepo:  walletRepo,
		auditRepo:   auditRepo,
		redisClient: redisClient,
	}
}

func (s *tradeService) ExecuteTrade(ctx context.Context, walletID, stockName string, op domain.OperationType) error {
	// Validate operation type first
	if op != domain.Buy && op != domain.Sell {
		return ErrInvalidOperation
	}

	// 1. Verify stock exists in bank (registered via POST /stocks)
	exists, err := s.bankRepo.StockExists(ctx, stockName)
	if err != nil {
		return fmt.Errorf("checking stock existence: %w", err)
	}
	if !exists {
		return ErrStockNotFound
	}

	bankKey := "bank:stocks"
	walletKey := fmt.Sprintf("wallet:%s:stocks", walletID)

	switch op {
	case domain.Buy:
		err = buyScript.Run(ctx, s.redisClient, []string{bankKey, walletKey}, stockName).Err()
		if err != nil && err.Error() == "INSUFFICIENT_BANK" {
			return ErrInsufficientBank
		}
	case domain.Sell:
		err = sellScript.Run(ctx, s.redisClient, []string{walletKey, bankKey}, stockName).Err()
		if err != nil && err.Error() == "INSUFFICIENT_WALLET" {
			return ErrInsufficientWallet
		}
	}
	if err != nil {
		return fmt.Errorf("executing trade script: %w", err)
	}

	// 2. Append to audit log (only on success)
	return s.auditRepo.Append(ctx, domain.LogEntry{
		Type:      op,
		WalletID:  walletID,
		StockName: stockName,
	})
}
