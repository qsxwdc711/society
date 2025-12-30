package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"
)

type StatsRepoInterface interface {
	ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error)
	ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error)
}

type StatsRepo struct {
	dao dao.StatsDaoInterface
}

func NewStatsRepo(dao dao.StatsDaoInterface) StatsRepoInterface {
	return &StatsRepo{dao: dao}
}

func (repo *StatsRepo) ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error) {
	return repo.dao.ProductSalesRank(ctx, limit)
}

func (repo *StatsRepo) ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error) {
	return repo.dao.ProductViewRank(ctx, limit, days)
}
