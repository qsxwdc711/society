package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepoInterface interface {
	Add(ctx context.Context, uid, pid primitive.ObjectID) error
	UpdateQty(ctx context.Context, uid, pid primitive.ObjectID, qty int64) error
	Remove(ctx context.Context, uid, pid primitive.ObjectID) error
	List(ctx context.Context, uid primitive.ObjectID) (domain.CartListRes, error)
}

type CartRepo struct {
	dao dao.CartDaoInterface
}

func NewCartRepo(dao dao.CartDaoInterface) CartRepoInterface {
	return &CartRepo{dao: dao}
}

func (r *CartRepo) Add(ctx context.Context, uid, pid primitive.ObjectID) error {
	return r.dao.Add(ctx, uid, pid)
}

func (r *CartRepo) UpdateQty(ctx context.Context, uid, pid primitive.ObjectID, qty int64) error {
	return r.dao.UpdateQty(ctx, uid, pid, qty)
}

func (r *CartRepo) Remove(ctx context.Context, uid, pid primitive.ObjectID) error {
	return r.dao.Remove(ctx, uid, pid)
}

func (r *CartRepo) List(ctx context.Context, uid primitive.ObjectID) (domain.CartListRes, error) {
	list, totalQty, totalPrice, err := r.dao.List(ctx, uid)
	if err != nil {
		return domain.CartListRes{}, err
	}
	return domain.CartListRes{
		List:       list,
		TotalQty:   totalQty,
		TotalPrice: totalPrice,
	}, nil
}
