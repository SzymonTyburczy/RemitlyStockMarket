package service

import (
	"context"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetWallet(ctx context.Context, walletID string) (*domain.Wallet, error) {
	stocks, err := s.repo.GetAllStocks(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return &domain.Wallet{ID: walletID, Stocks: stocks}, nil
}

func (s *walletService) GetStockQuantity(ctx context.Context, walletID, stockName string) (int, error) {
	return s.repo.GetQuantity(ctx, walletID, stockName)
}
