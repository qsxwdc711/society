package service

import (
	"context"
	"errors"
	"strings"

	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductAdminInterface interface {
	Add(ctx *gin.Context, req AddProductReq) (primitive.ObjectID, error)
	Update(ctx *gin.Context, req UpdateProductReq) error
	UpdateStatus(ctx *gin.Context, id string, status domain.ProductStatus) error
	Delete(ctx *gin.Context, id string) error
	FindById(ctx context.Context, id string) (domain.Product, error)
	List(ctx context.Context, req ListProductReq) (domain.PageResult[domain.Product], error)
}

type AddProductReq struct {
	CategoryId string   `json:"categoryId" binding:"required"`
	Title      string   `json:"title" binding:"required"`
	Intro      string   `json:"intro"`
	Cover      string   `json:"cover"`
	Images     []string `json:"images"`
	Price      int64    `json:"price" binding:"required"`
	Stock      int64    `json:"stock" binding:"required"`
}

type UpdateProductReq struct {
	Id         string   `json:"id"`
	CategoryId *string  `json:"categoryId"`
	Title      *string  `json:"title"`
	Intro      *string  `json:"intro"`
	Cover      *string  `json:"cover"`
	Images     []string `json:"images"` // 可选：传空数组代表清空
	Price      *int64   `json:"price"`
	Stock      *int64   `json:"stock"`
}

type ListProductReq struct {
	ginx.Page
	Q          string `form:"q"`
	Status     string `form:"status"`
	CategoryId string `form:"categoryId"`
}

type ProductAdminService struct {
	repo repository.ProductRepoInterface
	cat  repository.CategoryRepoInterface
}

func NewProductAdminService(repo repository.ProductRepoInterface, cat repository.CategoryRepoInterface) ProductAdminInterface {
	return &ProductAdminService{repo: repo, cat: cat}
}

func (svc *ProductAdminService) Add(ctx *gin.Context, req AddProductReq) (primitive.ObjectID, error) {
	cid, err := primitive.ObjectIDFromHex(req.CategoryId)
	if err != nil {
		return primitive.NilObjectID, errors.New("分类ID非法")
	}
	// 校验分类存在（加分项）
	if _, err := svc.cat.FindById(ctx, cid); err != nil {
		return primitive.NilObjectID, errors.New("分类不存在")
	}
	if strings.TrimSpace(req.Title) == "" || req.Price < 0 || req.Stock < 0 {
		return primitive.NilObjectID, errors.New("参数错误")
	}

	return svc.repo.Add(ctx, domain.Product{
		CategoryId: cid,
		Title:      strings.TrimSpace(req.Title),
		Intro:      strings.TrimSpace(req.Intro),
		Cover:      strings.TrimSpace(req.Cover),
		Images:     req.Images,
		Price:      req.Price,
		Stock:      req.Stock,
		Status:     domain.ProductStatusOff,
	})
}

func (svc *ProductAdminService) Update(ctx *gin.Context, req UpdateProductReq) error {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return errors.New("商品ID非法")
	}
	set := bson.M{}
	if req.CategoryId != nil {
		cid, err := primitive.ObjectIDFromHex(*req.CategoryId)
		if err != nil {
			return errors.New("分类ID非法")
		}
		if _, err := svc.cat.FindById(ctx, cid); err != nil {
			return errors.New("分类不存在")
		}
		set["categoryId"] = cid
	}
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			return errors.New("商品名不能为空")
		}
		set["title"] = t
	}
	if req.Intro != nil {
		set["intro"] = strings.TrimSpace(*req.Intro)
	}
	if req.Cover != nil {
		set["cover"] = strings.TrimSpace(*req.Cover)
	}
	// Images：如果你希望“未传则不更新”，可以改为 *[]string 指针
	if req.Images != nil {
		set["images"] = req.Images
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return errors.New("价格非法")
		}
		set["price"] = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return errors.New("库存非法")
		}
		set["stock"] = *req.Stock
	}
	if len(set) == 0 {
		return errors.New("没有可更新字段")
	}

	err = svc.repo.Update(ctx, oid, set)
	if err == mongo.ErrNoDocuments {
		return errors.New("商品不存在")
	}
	return err
}

func (svc *ProductAdminService) UpdateStatus(ctx *gin.Context, id string, status domain.ProductStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("商品ID非法")
	}
	if status != domain.ProductStatusOn && status != domain.ProductStatusOff {
		return errors.New("状态非法")
	}
	err = svc.repo.UpdateStatus(ctx, oid, status)
	if err == mongo.ErrNoDocuments {
		return errors.New("商品不存在")
	}
	return err
}

func (svc *ProductAdminService) Delete(ctx *gin.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("商品ID非法")
	}
	err = svc.repo.Delete(ctx, oid)
	if err == mongo.ErrNoDocuments {
		return errors.New("商品不存在")
	}
	return err
}

func (svc *ProductAdminService) FindById(ctx context.Context, id string) (domain.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Product{}, errors.New("商品ID非法")
	}
	return svc.repo.FindById(ctx, oid)
}

func (svc *ProductAdminService) List(ctx context.Context, req ListProductReq) (domain.PageResult[domain.Product], error) {
	if !req.Page.Check() {
		return domain.PageResult[domain.Product]{}, errors.New("分页参数非法")
	}
	var st *domain.ProductStatus
	if req.Status != "" {
		tmp := domain.ProductStatus(req.Status)
		if tmp != domain.ProductStatusOn && tmp != domain.ProductStatusOff {
			return domain.PageResult[domain.Product]{}, errors.New("状态参数非法")
		}
		st = &tmp
	}
	var cid *primitive.ObjectID
	if req.CategoryId != "" {
		tmp, err := primitive.ObjectIDFromHex(req.CategoryId)
		if err != nil {
			return domain.PageResult[domain.Product]{}, errors.New("分类ID非法")
		}
		cid = &tmp
	}
	return svc.repo.List(ctx, req.Q, cid, st, req.Page)
}
