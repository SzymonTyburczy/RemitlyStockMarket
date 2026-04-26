package repository

// redisAuditRepo implements AuditRepository using a Redis list.
// Key: "audit:log"  →  List of JSON-encoded LogEntry (append-only, RPUSH / LRANGE)
type redisAuditRepo struct{}
