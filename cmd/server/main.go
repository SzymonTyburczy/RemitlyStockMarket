package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/szymontyburczy/remitly-stock-market/internal/config"
	"github.com/szymontyburczy/remitly-stock-market/internal/handler"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
	"github.com/szymontyburczy/remitly-stock-market/internal/service"
)

func main() {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// 2. Connect to Redis
	redisClient, err := repository.NewRedisClient(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to connect to Redis", "addr", cfg.RedisURL, "err", err)
		os.Exit(1)
	}

	// 3. Wire repositories
	bankRepo := repository.NewBankRepository(redisClient)
	walletRepo := repository.NewWalletRepository(redisClient)
	auditRepo := repository.NewAuditRepository(redisClient)

	// 4. Wire services
	bankSvc := service.NewBankService(bankRepo)
	walletSvc := service.NewWalletService(walletRepo)
	auditSvc := service.NewAuditService(auditRepo)
	tradeSvc := service.NewTradeService(bankRepo, walletRepo, auditRepo, redisClient)

	// 5. Wire HTTP router
	router := handler.NewRouter(walletSvc, bankSvc, tradeSvc, auditSvc)

	// 6. Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	slog.Info("starting stock market service", "addr", addr, "redis", cfg.RedisURL)
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
