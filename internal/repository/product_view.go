package repository

import (
	"society/internal/domain"
	"society/internal/repository/dao"

	"github.com/gin-gonic/gin"
)

type ProductViewRepoInterface interface {
	Add(ctx *gin.Context, v domain.ProductView) error
}

type ProductViewRepo struct {
	dao dao.ProductViewDaoInterface
}

func NewProductViewRepo(dao dao.ProductViewDaoInterface) ProductViewRepoInterface {
	return &ProductViewRepo{dao: dao}
}

func (repo *ProductViewRepo) Add(ctx *gin.Context, v domain.ProductView) error {
	return repo.dao.Add(ctx, v)
}
