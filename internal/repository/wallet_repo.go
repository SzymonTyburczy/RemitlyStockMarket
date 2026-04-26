package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
)

type redisWalletRepo struct {
	client *redis.Client
}

func NewWalletRepository(client *redis.Client) WalletRepository {
	return &redisWalletRepo{client: client}
}

func walletStocksKey(walletID string) string {
	return fmt.Sprintf("wallet:%s:stocks", walletID)
}

func (r *redisWalletRepo) GetAllStocks(ctx context.Context, walletID string) ([]domain.Stock, error) {
	result, err := r.client.HGetAll(ctx, walletStocksKey(walletID)).Result()
	if err != nil {
		return nil, err
	}
	stocks := make([]domain.Stock, 0, len(result))
	for name, qtyStr := range result {
		qty, _ := strconv.Atoi(qtyStr)
		stocks = append(stocks, domain.Stock{Name: name, Quantity: qty})
	}
	return stocks, nil
}

func (r *redisWalletRepo) GetQuantity(ctx context.Context, walletID, stockName string) (int, error) {
	val, err := r.client.HGet(ctx, walletStocksKey(walletID), stockName).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *redisWalletRepo) IncrBy(ctx context.Context, walletID, stockName string, delta int) error {
	return r.client.HIncrBy(ctx, walletStocksKey(walletID), stockName, int64(delta)).Err()
}
