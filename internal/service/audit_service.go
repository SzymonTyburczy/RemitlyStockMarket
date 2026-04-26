package service

import (
	"context"

	"github.com/szymontyburczy/remitly-stock-market/internal/domain"
	"github.com/szymontyburczy/remitly-stock-market/internal/repository"
)

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) GetAll(ctx context.Context) ([]domain.LogEntry, error) {
	return s.repo.GetAll(ctx)
}
