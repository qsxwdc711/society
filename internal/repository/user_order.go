package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOrderRepoInterface interface {
	CreateOrderFromCartTx(ctx context.Context, uid primitive.ObjectID, productIds []primitive.ObjectID) (domain.Order, error)
	PayOrderTx(ctx context.Context, uid primitive.ObjectID, orderId primitive.ObjectID) (domain.Order, error)
}

type UserOrderRepo struct {
	dao dao.UserOrderDaoInterface
}

func NewUserOrderRepo(dao dao.UserOrderDaoInterface) UserOrderRepoInterface {
	return &UserOrderRepo{dao: dao}
}

func (r *UserOrderRepo) CreateOrderFromCartTx(ctx context.Context, uid primitive.ObjectID, productIds []primitive.ObjectID) (domain.Order, error) {
	return r.dao.CreateOrderFromCartTx(ctx, uid, productIds)
}

func (r *UserOrderRepo) PayOrderTx(ctx context.Context, uid primitive.ObjectID, orderId primitive.ObjectID) (domain.Order, error) {
	return r.dao.PayOrderTx(ctx, uid, orderId)
}
