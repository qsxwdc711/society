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

type PromotionRepoInterface interface {
	Add(ctx *gin.Context, p domain.Promotion) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	BindProducts(ctx *gin.Context, id primitive.ObjectID, productIds []primitive.ObjectID) error
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Promotion, error)
	List(ctx context.Context, page ginx.Page) (domain.PageResult[domain.Promotion], error)
}

type PromotionRepo struct {
	dao dao.PromotionDaoInterface
}

func NewPromotionRepo(dao dao.PromotionDaoInterface) PromotionRepoInterface {
	return &PromotionRepo{dao: dao}
}

func (repo *PromotionRepo) Add(ctx *gin.Context, p domain.Promotion) (primitive.ObjectID, error) {
	return repo.dao.Add(ctx, p)
}

func (repo *PromotionRepo) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
	return repo.dao.Update(ctx, id, set)
}

func (repo *PromotionRepo) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *PromotionRepo) BindProducts(ctx *gin.Context, id primitive.ObjectID, productIds []primitive.ObjectID) error {
	return repo.dao.BindProducts(ctx, id, productIds)
}

func (repo *PromotionRepo) FindById(ctx context.Context, id primitive.ObjectID) (domain.Promotion, error) {
	return repo.dao.FindById(ctx, id)
}

func (repo *PromotionRepo) List(ctx context.Context, page ginx.Page) (domain.PageResult[domain.Promotion], error) {
	list, total, err := repo.dao.List(ctx, page)
	if err != nil {
		return domain.PageResult[domain.Promotion]{}, err
	}
	return domain.PageResult[domain.Promotion]{
		List:  list,
		Total: total,
		Page:  int64(page.Page),
		Size:  int64(page.Size),
	}, nil
}
