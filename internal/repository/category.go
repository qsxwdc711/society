package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepoInterface interface {
	Add(ctx *gin.Context, c domain.Category) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]domain.Category, error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Category, error)
}

type CategoryRepo struct {
	dao dao.CategoryDaoInterface
}

func NewCategoryRepo(dao dao.CategoryDaoInterface) CategoryRepoInterface {
	return &CategoryRepo{dao: dao}
}

func (repo *CategoryRepo) Add(ctx *gin.Context, c domain.Category) (primitive.ObjectID, error) {
	return repo.dao.Add(ctx, c)
}

func (repo *CategoryRepo) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
	return repo.dao.Update(ctx, id, set)
}

func (repo *CategoryRepo) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	return repo.dao.List(ctx)
}

func (repo *CategoryRepo) FindById(ctx context.Context, id primitive.ObjectID) (domain.Category, error) {
	return repo.dao.FindById(ctx, id)
}
