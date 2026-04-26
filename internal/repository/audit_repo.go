package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
)

// auditLogKey is the Redis list key for the append-only audit log.
const auditLogKey = "audit:log"

type redisAuditRepo struct {
	client *redis.Client
}

func NewAuditRepository(client *redis.Client) AuditRepository {
	return &redisAuditRepo{client: client}
}

// Append adds a log entry to the tail of the audit log list.
func (r *redisAuditRepo) Append(ctx context.Context, entry domain.LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return r.client.RPush(ctx, auditLogKey, data).Err()
}

// GetAll returns all log entries in insertion order.
func (r *redisAuditRepo) GetAll(ctx context.Context) ([]domain.LogEntry, error) {
	result, err := r.client.LRange(ctx, auditLogKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	entries := make([]domain.LogEntry, 0, len(result))
	for _, s := range result {
		var entry domain.LogEntry
		if err := json.Unmarshal([]byte(s), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
