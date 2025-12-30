package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartInterface interface {
	Add(ctx context.Context, uid primitive.ObjectID, productId string) error
	UpdateQty(ctx context.Context, uid primitive.ObjectID, productId string, qty int64) error
	Remove(ctx context.Context, uid primitive.ObjectID, productId string) error
	List(ctx context.Context, uid primitive.ObjectID) (domain.CartListRes, error)
}

type CartService struct {
	repo repository.CartRepoInterface
}

func NewCartService(repo repository.CartRepoInterface) CartInterface {
	return &CartService{repo: repo}
}

func (svc *CartService) Add(ctx context.Context, uid primitive.ObjectID, productId string) error {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return errors.New("商品ID非法")
	}
	return svc.repo.Add(ctx, uid, pid)
}

func (svc *CartService) UpdateQty(ctx context.Context, uid primitive.ObjectID, productId string, qty int64) error {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return errors.New("商品ID非法")
	}
	return svc.repo.UpdateQty(ctx, uid, pid, qty)
}

func (svc *CartService) Remove(ctx context.Context, uid primitive.ObjectID, productId string) error {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return errors.New("商品ID非法")
	}
	return svc.repo.Remove(ctx, uid, pid)
}

func (svc *CartService) List(ctx context.Context, uid primitive.ObjectID) (domain.CartListRes, error) {
	return svc.repo.List(ctx, uid)
}
