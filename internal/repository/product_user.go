package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProductRepoInterface interface {
	List(ctx context.Context, q string, categoryId *primitive.ObjectID, sort string, minPrice *int64, maxPrice *int64, page ginx.Page) (domain.PageResult[domain.Product], error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error)
}

type UserProductRepo struct {
	dao dao.UserProductDaoInterface
}

func NewUserProductRepo(dao dao.UserProductDaoInterface) UserProductRepoInterface {
	return &UserProductRepo{dao: dao}
}

func (repo *UserProductRepo) List(ctx context.Context, q string, categoryId *primitive.ObjectID, sort string, minPrice *int64, maxPrice *int64, page ginx.Page) (domain.PageResult[domain.Product], error) {
	list, total, err := repo.dao.List(ctx, q, categoryId, sort, minPrice, maxPrice, page)
	if err != nil {
		return domain.PageResult[domain.Product]{}, err
	}
	return domain.PageResult[domain.Product]{
		List:  list,
		Total: total,
		Page:  int64(page.Page),
		Size:  int64(page.Size),
	}, nil
}

func (repo *UserProductRepo) FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error) {
	return repo.dao.FindById(ctx, id)
}
