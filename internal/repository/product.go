package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepoInterface interface {
	Add(ctx *gin.Context, p domain.Product) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	UpdateStatus(ctx *gin.Context, id primitive.ObjectID, status domain.ProductStatus) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error)
	List(ctx context.Context, q string, categoryId *primitive.ObjectID, status *domain.ProductStatus, page ginx.Page) (domain.PageResult[domain.Product], error)
}

type ProductRepo struct {
	dao dao.ProductDaoInterface
}

func NewProductRepo(dao dao.ProductDaoInterface) ProductRepoInterface {
	return &ProductRepo{dao: dao}
}

func (repo *ProductRepo) Add(ctx *gin.Context, p domain.Product) (primitive.ObjectID, error) {
	return repo.dao.Add(ctx, p)
}

func (repo *ProductRepo) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
	return repo.dao.Update(ctx, id, set)
}

func (repo *ProductRepo) UpdateStatus(ctx *gin.Context, id primitive.ObjectID, status domain.ProductStatus) error {
	return repo.dao.UpdateStatus(ctx, id, status)
}

func (repo *ProductRepo) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *ProductRepo) FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error) {
	return repo.dao.FindById(ctx, id)
}

func (repo *ProductRepo) List(ctx context.Context, q string, categoryId *primitive.ObjectID, status *domain.ProductStatus, page ginx.Page) (domain.PageResult[domain.Product], error) {
	list, total, err := repo.dao.List(ctx, q, categoryId, status, page)
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
