package service

import (
	"context"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

type bankService struct {
	repo repository.BankRepository
}

func NewBankService(repo repository.BankRepository) BankService {
	return &bankService{repo: repo}
}

func (s *bankService) SetStocks(ctx context.Context, stocks []domain.Stock) error {
	return s.repo.SetStocks(ctx, stocks)
}

func (s *bankService) GetStocks(ctx context.Context) ([]domain.Stock, error) {
	return s.repo.GetAllStocks(ctx)
}
