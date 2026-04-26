package repository

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
)

// bankStocksKey is the Redis hash key for the bank's stock inventory.
const bankStocksKey = "bank:stocks"

type redisBankRepo struct {
	client *redis.Client
}

func NewBankRepository(client *redis.Client) BankRepository {
	return &redisBankRepo{client: client}
}

// SetStocks atomically replaces the entire bank stock state.
func (r *redisBankRepo) SetStocks(ctx context.Context, stocks []domain.Stock) error {
	pipe := r.client.TxPipeline()
	pipe.Del(ctx, bankStocksKey)
	for _, s := range stocks {
		pipe.HSet(ctx, bankStocksKey, s.Name, s.Quantity)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisBankRepo) GetAllStocks(ctx context.Context) ([]domain.Stock, error) {
	result, err := r.client.HGetAll(ctx, bankStocksKey).Result()
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

func (r *redisBankRepo) StockExists(ctx context.Context, name string) (bool, error) {
	return r.client.HExists(ctx, bankStocksKey, name).Result()
}

func (r *redisBankRepo) GetQuantity(ctx context.Context, name string) (int, error) {
	val, err := r.client.HGet(ctx, bankStocksKey, name).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *redisBankRepo) IncrBy(ctx context.Context, name string, delta int) error {
	return r.client.HIncrBy(ctx, bankStocksKey, name, int64(delta)).Err()
}
