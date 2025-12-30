package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminOrderRepoInterface interface {
	List(ctx context.Context, status *domain.OrderStatus, page ginx.Page) (domain.PageResult[domain.Order], error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error
}

type AdminOrderRepo struct {
	dao dao.AdminOrderDaoInterface
}

func NewAdminOrderRepo(dao dao.AdminOrderDaoInterface) AdminOrderRepoInterface {
	return &AdminOrderRepo{dao: dao}
}

func (repo *AdminOrderRepo) List(ctx context.Context, status *domain.OrderStatus, page ginx.Page) (domain.PageResult[domain.Order], error) {
	list, total, err := repo.dao.List(ctx, status, page)
	if err != nil {
		return domain.PageResult[domain.Order]{}, err
	}
	return domain.PageResult[domain.Order]{
		List:  list,
		Total: total,
		Page:  int64(page.Page),
		Size:  int64(page.Size),
	}, nil
}

func (repo *AdminOrderRepo) FindById(ctx context.Context, id primitive.ObjectID) (domain.Order, error) {
	return repo.dao.FindById(ctx, id)
}

func (repo *AdminOrderRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error {
	return repo.dao.UpdateStatus(ctx, id, status)
}
