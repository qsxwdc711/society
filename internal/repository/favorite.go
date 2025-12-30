package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FavoriteRepoInterface interface {
	Add(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error
	Remove(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error
	Check(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) (bool, error)
	List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) (domain.PageResult[domain.FavoriteItem], error)
}

type FavoriteRepo struct {
	dao dao.FavoriteDaoInterface
}

func NewFavoriteRepo(dao dao.FavoriteDaoInterface) FavoriteRepoInterface {
	return &FavoriteRepo{dao: dao}
}

func (repo *FavoriteRepo) Add(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error {
	return repo.dao.Add(ctx, uid, productId)
}

func (repo *FavoriteRepo) Remove(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error {
	return repo.dao.Remove(ctx, uid, productId)
}

func (repo *FavoriteRepo) Check(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) (bool, error) {
	return repo.dao.Check(ctx, uid, productId)
}

func (repo *FavoriteRepo) List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) (domain.PageResult[domain.FavoriteItem], error) {
	list, total, err := repo.dao.List(ctx, uid, page)
	if err != nil {
		return domain.PageResult[domain.FavoriteItem]{}, err
	}
	return domain.PageResult[domain.FavoriteItem]{
		List:  list,
		Total: total,
		Page:  int64(page.Page),
		Size:  int64(page.Size),
	}, nil
}
