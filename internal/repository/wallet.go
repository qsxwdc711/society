package repository

import (
	"context"
	"society/internal/repository/dao"
)

type WalletRepoInterface interface {
	EnsureWallet(ctx context.Context, uid any) error
	GetBalance(ctx context.Context, uid any) (int64, error)
	UpdateBalance(ctx context.Context, uid any, delta int64) (int64, error)
}

type WalletRepo struct {
	dao dao.WalletDaoInterface
}

func NewWalletRepo(dao dao.WalletDaoInterface) WalletRepoInterface {
	return &WalletRepo{dao: dao}
}

func (r *WalletRepo) EnsureWallet(ctx context.Context, uid any) error {
	return r.dao.EnsureWallet(ctx, uid)
}

func (r *WalletRepo) GetBalance(ctx context.Context, uid any) (int64, error) {
	return r.dao.GetBalance(ctx, uid)
}

func (r *WalletRepo) UpdateBalance(ctx context.Context, uid any, delta int64) (int64, error) {
	return r.dao.UpdateBalance(ctx, uid, delta)
}
