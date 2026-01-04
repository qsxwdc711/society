package repository

import (
	"context"
	"society/internal/repository/dao"
)

type BillRepoInterface interface {
	Insert(ctx context.Context, doc any) error
	PageByUid(ctx context.Context, uid any, page int, size int, typ string) ([]any, int64, error)
}

type BillRepo struct {
	dao dao.BillDaoInterface
}

func NewBillRepo(dao dao.BillDaoInterface) BillRepoInterface {
	return &BillRepo{dao: dao}
}

func (r *BillRepo) Insert(ctx context.Context, doc any) error {
	return r.dao.Insert(ctx, doc)
}

func (r *BillRepo) PageByUid(ctx context.Context, uid any, page int, size int, typ string) ([]any, int64, error) {
	return r.dao.PageByUid(ctx, uid, page, size, typ)
}
