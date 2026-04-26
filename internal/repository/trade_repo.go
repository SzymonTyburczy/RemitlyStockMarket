package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// TradeRepository executes atomic buy/sell operations via Redis Lua scripts.
// It lives in the repository layer so the service layer stays free of Redis imports.
type TradeRepository interface {
	ExecuteBuy(ctx context.Context, walletID, stockName string) error
	ExecuteSell(ctx context.Context, walletID, stockName string) error
}

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

type redisTradeRepo struct {
	client *redis.Client
}

func NewTradeRepository(client *redis.Client) TradeRepository {
	return &redisTradeRepo{client: client}
}

func (r *redisTradeRepo) ExecuteBuy(ctx context.Context, walletID, stockName string) error {
	err := buyScript.Run(ctx, r.client,
		[]string{bankStocksKey, walletStocksKey(walletID)},
		stockName,
	).Err()
	if err != nil && err.Error() == "INSUFFICIENT_BANK" {
		return ErrInsufficientBank
	}
	return err
}

func (r *redisTradeRepo) ExecuteSell(ctx context.Context, walletID, stockName string) error {
	err := sellScript.Run(ctx, r.client,
		[]string{walletStocksKey(walletID), bankStocksKey},
		stockName,
	).Err()
	if err != nil && err.Error() == "INSUFFICIENT_WALLET" {
		return ErrInsufficientWallet
	}
	return err
}

// Sentinel errors for trade operations — checked by the service layer.
var (
	ErrInsufficientBank   = fmt.Errorf("insufficient stock in bank")
	ErrInsufficientWallet = fmt.Errorf("insufficient stock in wallet")
)
