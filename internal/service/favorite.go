package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteInterface interface {
	Add(ctx context.Context, uid primitive.ObjectID, productId string) error
	Remove(ctx context.Context, uid primitive.ObjectID, productId string) error
	Check(ctx context.Context, uid primitive.ObjectID, productId string) (domain.FavoriteCheckRes, error)
	List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) (domain.PageResult[domain.FavoriteItem], error)
}

type FavoriteService struct {
	repo repository.FavoriteRepoInterface
}

func NewFavoriteService(repo repository.FavoriteRepoInterface) FavoriteInterface {
	return &FavoriteService{repo: repo}
}

func (svc *FavoriteService) Add(ctx context.Context, uid primitive.ObjectID, productId string) error {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return errors.New("商品ID非法")
	}
	err = svc.repo.Add(ctx, uid, pid)
	if errors.Is(err, dao.ErrFavoriteExists) {
		return errors.New("已收藏")
	}
	return err
}

func (svc *FavoriteService) Remove(ctx context.Context, uid primitive.ObjectID, productId string) error {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return errors.New("商品ID非法")
	}
	return svc.repo.Remove(ctx, uid, pid)
}

func (svc *FavoriteService) Check(ctx context.Context, uid primitive.ObjectID, productId string) (domain.FavoriteCheckRes, error) {
	pid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return domain.FavoriteCheckRes{}, errors.New("商品ID非法")
	}
	ok, err := svc.repo.Check(ctx, uid, pid)
	if err != nil {
		return domain.FavoriteCheckRes{}, err
	}
	return domain.FavoriteCheckRes{IsFavorite: ok}, nil
}

func (svc *FavoriteService) List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) (domain.PageResult[domain.FavoriteItem], error) {
	if !page.Check() {
		return domain.PageResult[domain.FavoriteItem]{}, errors.New("分页参数非法")
	}
	return svc.repo.List(ctx, uid, page)
}
