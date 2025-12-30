package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"
	"society/internal/repository/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOrderInterface interface {
	CreateFromCart(ctx context.Context, uid primitive.ObjectID, req CreateOrderReq) (domain.Order, error)
	Pay(ctx context.Context, uid primitive.ObjectID, orderId string) (domain.Order, error)
}

type CreateOrderReq struct {
	// 为空表示结算全部购物车；不为空表示结算指定商品
	ProductIds []string `json:"productIds"`
}

type UserOrderService struct {
	repo repository.UserOrderRepoInterface
}

func NewUserOrderService(repo repository.UserOrderRepoInterface) UserOrderInterface {
	return &UserOrderService{repo: repo}
}

func (svc *UserOrderService) CreateFromCart(ctx context.Context, uid primitive.ObjectID, req CreateOrderReq) (domain.Order, error) {
	var pids []primitive.ObjectID
	if len(req.ProductIds) > 0 {
		pids = make([]primitive.ObjectID, 0, len(req.ProductIds))
		for _, s := range req.ProductIds {
			id, err := primitive.ObjectIDFromHex(s)
			if err != nil {
				return domain.Order{}, errors.New("商品ID非法")
			}
			pids = append(pids, id)
		}
	}
	order, err := svc.repo.CreateOrderFromCartTx(ctx, uid, pids)
	if errors.Is(err, dao.ErrCartEmpty) {
		return domain.Order{}, errors.New("购物车为空")
	}
	if errors.Is(err, dao.ErrStockNotEnough) {
		return domain.Order{}, errors.New("库存不足")
	}
	return order, err
}

func (svc *UserOrderService) Pay(ctx context.Context, uid primitive.ObjectID, orderId string) (domain.Order, error) {
	oid, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return domain.Order{}, errors.New("订单ID非法")
	}

	order, err := svc.repo.PayOrderTx(ctx, uid, oid)
	if errors.Is(err, dao.ErrOrderNotFound) {
		return domain.Order{}, errors.New("订单不存在")
	}
	if errors.Is(err, dao.ErrOrderStatusInvalid) {
		return domain.Order{}, errors.New("订单状态不允许支付")
	}
	if errors.Is(err, dao.ErrBalanceNotEnough) {
		return domain.Order{}, errors.New("余额不足")
	}
	return order, err
}
