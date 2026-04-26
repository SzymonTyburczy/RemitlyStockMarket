# Stock Market Service

A simplified, highly available in-memory stock exchange built with **Go**, **Redis**, and **nginx**.

## Architecture

```
Client → localhost:<PORT> (nginx — load balancer)
                 ↓  round-robin with failover
    ┌────────────┬────────────┬────────────┐
    │ service-1  │ service-2  │ service-3  │   ← Go HTTP instances
    └────────────┴────────────┴────────────┘
                        ↓
                   Redis :6379  (shared state)
```

- **3 application instances** behind nginx — killing one leaves two running
- **nginx `proxy_next_upstream`** retries failed requests on healthy instances automatically
- **Redis** holds all state: bank stocks (Hash), wallet stocks (Hash per wallet), audit log (List)
- **`POST /chaos`** calls `os.Exit(1)` on the serving instance — the other two continue unaffected

## Requirements

- Docker (with Docker Compose v2)
- `envsubst` (Linux/macOS — part of `gettext`; on Windows handled by PowerShell)

## Quick Start

### Linux / macOS
```bash
./scripts/start.sh 8080
```

### Windows
```bat
scripts\start.bat 8080
```

The service will be available at `http://localhost:8080` (replace `8080` with your chosen port).

To stop:
```bash
docker compose down
```

## API Reference

### Bank

| Method | Path | Body | Description |
|--------|------|------|-------------|
| `GET`  | `/stocks` | — | Returns current bank inventory |
| `POST` | `/stocks` | `{"stocks":[{"name":"AAPL","quantity":100}]}` | Replaces entire bank state |

### Wallets

| Method | Path | Body | Description |
|--------|------|------|-------------|
| `GET`  | `/wallets/{wallet_id}` | — | Returns wallet with all stocks |
| `GET`  | `/wallets/{wallet_id}/stocks/{stock_name}` | — | Returns single stock quantity |
| `POST` | `/wallets/{wallet_id}/stocks/{stock_name}` | `{"type":"buy"}` or `{"type":"sell"}` | Executes trade |

**Trade rules:**
- Stock must exist in bank (registered via `POST /stocks`) → `404` otherwise
- `buy`: bank must have ≥ 1 unit → `400` if empty
- `sell`: wallet must have ≥ 1 unit → `400` if empty
- Wallet is auto-created on first operation
- Price is fixed at 1 (no funds tracking)

### Audit Log

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/log` | Returns all successful trades in order of occurrence |

### Chaos

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/chaos` | Kills the instance serving this request |

## Example Flow

```bash
# 1. Seed the bank
curl -X POST http://localhost:8080/stocks \
  -H 'Content-Type: application/json' \
  -d '{"stocks":[{"name":"AAPL","quantity":100},{"name":"GOOG","quantity":50}]}'

# 2. Buy 1 AAPL for wallet "alice"
curl -X POST http://localhost:8080/wallets/alice/stocks/AAPL \
  -H 'Content-Type: application/json' \
  -d '{"type":"buy"}'

# 3. Check alice's wallet
curl http://localhost:8080/wallets/alice

# 4. Check how many AAPLs alice holds
curl http://localhost:8080/wallets/alice/stocks/AAPL

# 5. Check bank state
curl http://localhost:8080/stocks

# 6. Check audit log
curl http://localhost:8080/log

# 7. Kill one instance (HA demo) — service continues via the other two
curl -X POST http://localhost:8080/chaos
curl http://localhost:8080/stocks  # still works
```

## Running Tests

### Unit tests (no infrastructure needed)
```bash
go test ./internal/service/... ./internal/handler/...
```

### Integration tests (requires Redis)
```bash
# Start Redis locally
docker run -d -p 6379:6379 redis:7-alpine

# Run all tests including integration
REDIS_TEST_URL=localhost:6379 go test ./...
```

Integration tests cover:
- `internal/repository/` — all three Redis repositories
- `internal/service/trade_service_test.go` — atomic buy/sell, all error cases, audit log

## Project Structure

```
.
├── cmd/server/main.go              # Entry point — DI wiring + HTTP server
├── internal/
│   ├── config/                     # Environment-based config
│   ├── domain/                     # Domain models (Stock, Wallet, LogEntry)
│   ├── handler/                    # HTTP handlers + router (chi)
│   ├── service/                    # Business logic + service interfaces
│   └── repository/                 # Redis repository implementations
├── nginx/nginx.conf.template       # nginx load balancer config (${PORT} placeholder)
├── scripts/
│   ├── start.sh                    # Linux/macOS entry point
│   └── start.bat                   # Windows entry point
├── Dockerfile                      # Multi-stage build (Go binary on Alpine)
└── docker-compose.yml              # nginx + 3x stock-service + Redis
```

## Design Notes

- **Atomic trades**: Buy/sell use Redis Lua scripts — the quantity check and decrement are a single atomic operation, eliminating race conditions under concurrent load
- **Layered architecture**: `handler -> service -> repository`, each layer depends only on interfaces — fully mockable for testing
- **Sentinel errors**: `service/errors.go` defines typed errors mapped cleanly to HTTP status codes in handlers
- **Multi-stage Dockerfile**: Final image is a minimal Alpine container with only the static Go binary (`CGO_ENABLED=0`)