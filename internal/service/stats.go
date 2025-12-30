package service

import (
	"context"

	"society/internal/domain"
	"society/internal/repository"
)

type StatsInterface interface {
	ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error)
	ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error)
}

type StatsService struct {
	repo repository.StatsRepoInterface
}

func NewStatsService(repo repository.StatsRepoInterface) StatsInterface {
	return &StatsService{repo: repo}
}

func (svc *StatsService) ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error) {
	return svc.repo.ProductSalesRank(ctx, limit)
}

func (svc *StatsService) ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error) {
	return svc.repo.ProductViewRank(ctx, limit, days)
}
