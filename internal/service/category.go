package service

import (
	"context"
	"errors"
	"strings"

	"society/internal/domain"
	"society/internal/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryInterface interface {
	Add(ctx *gin.Context, req AddCategoryReq) (primitive.ObjectID, error)
	Update(ctx *gin.Context, req UpdateCategoryReq) error
	Delete(ctx *gin.Context, id string) error
	List(ctx context.Context) ([]domain.Category, error)
}

type AddCategoryReq struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type UpdateCategoryReq struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
	Desc *string `json:"desc"`
}

type CategoryService struct {
	repo repository.CategoryRepoInterface
}

func NewCategoryService(repo repository.CategoryRepoInterface) CategoryInterface {
	return &CategoryService{repo: repo}
}

func (svc *CategoryService) Add(ctx *gin.Context, req AddCategoryReq) (primitive.ObjectID, error) {
	if strings.TrimSpace(req.Name) == "" {
		return primitive.NilObjectID, errors.New("分类名不能为空")
	}
	return svc.repo.Add(ctx, domain.Category{
		Name: strings.TrimSpace(req.Name),
		Desc: strings.TrimSpace(req.Desc),
	})
}

func (svc *CategoryService) Update(ctx *gin.Context, req UpdateCategoryReq) error {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return errors.New("分类ID非法")
	}
	set := bson.M{}
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		if n == "" {
			return errors.New("分类名不能为空")
		}
		set["name"] = n
	}
	if req.Desc != nil {
		set["desc"] = strings.TrimSpace(*req.Desc)
	}
	if len(set) == 0 {
		return errors.New("没有可更新字段")
	}
	err = svc.repo.Update(ctx, oid, set)
	if err == mongo.ErrNoDocuments {
		return errors.New("分类不存在")
	}
	return err
}

func (svc *CategoryService) Delete(ctx *gin.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("分类ID非法")
	}
	err = svc.repo.Delete(ctx, oid)
	if err == mongo.ErrNoDocuments {
		return errors.New("分类不存在")
	}
	return err
}

func (svc *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	return svc.repo.List(ctx)
}
