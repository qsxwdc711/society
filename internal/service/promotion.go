package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromotionInterface interface {
	Add(ctx *gin.Context, req AddPromotionReq) (primitive.ObjectID, error)
	Update(ctx *gin.Context, req UpdatePromotionReq) error
	Delete(ctx *gin.Context, id string) error
	BindProducts(ctx *gin.Context, req BindPromotionProductsReq) error
	List(ctx context.Context, req ginx.Page) (domain.PageResult[domain.Promotion], error)
	FindById(ctx context.Context, id string) (domain.Promotion, error)
}

type AddPromotionReq struct {
	Name    string               `json:"name" binding:"required"`
	Type    domain.PromotionType `json:"type" binding:"required"`
	Desc    string               `json:"desc"`
	StartAt time.Time            `json:"startAt"`
	EndAt   time.Time            `json:"endAt"`
}

type UpdatePromotionReq struct {
	Id      string                `json:"id"`
	Name    *string               `json:"name"`
	Type    *domain.PromotionType `json:"type"`
	Desc    *string               `json:"desc"`
	StartAt *time.Time            `json:"startAt"`
	EndAt   *time.Time            `json:"endAt"`
}

type BindPromotionProductsReq struct {
	Id         string   `json:"id"`
	ProductIds []string `json:"productIds"`
}

type PromotionService struct {
	repo repository.PromotionRepoInterface
}

func NewPromotionService(repo repository.PromotionRepoInterface) PromotionInterface {
	return &PromotionService{repo: repo}
}

func (svc *PromotionService) Add(ctx *gin.Context, req AddPromotionReq) (primitive.ObjectID, error) {
	if strings.TrimSpace(req.Name) == "" {
		return primitive.NilObjectID, errors.New("促销名不能为空")
	}
	// 简单校验时间
	if !req.StartAt.IsZero() && !req.EndAt.IsZero() && req.EndAt.Before(req.StartAt) {
		return primitive.NilObjectID, errors.New("结束时间不能早于开始时间")
	}
	return svc.repo.Add(ctx, domain.Promotion{
		Name:       strings.TrimSpace(req.Name),
		Type:       req.Type,
		Desc:       strings.TrimSpace(req.Desc),
		ProductIds: nil,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
	})
}

func (svc *PromotionService) Update(ctx *gin.Context, req UpdatePromotionReq) error {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return errors.New("促销ID非法")
	}
	set := bson.M{}
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		if n == "" {
			return errors.New("促销名不能为空")
		}
		set["name"] = n
	}
	if req.Type != nil {
		set["type"] = *req.Type
	}
	if req.Desc != nil {
		set["desc"] = strings.TrimSpace(*req.Desc)
	}
	if req.StartAt != nil {
		set["startAt"] = *req.StartAt
	}
	if req.EndAt != nil {
		set["endAt"] = *req.EndAt
	}
	if len(set) == 0 {
		return errors.New("没有可更新字段")
	}
	err = svc.repo.Update(ctx, oid, set)
	if err == mongo.ErrNoDocuments {
		return errors.New("促销不存在")
	}
	return err
}

func (svc *PromotionService) Delete(ctx *gin.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("促销ID非法")
	}
	err = svc.repo.Delete(ctx, oid)
	if err == mongo.ErrNoDocuments {
		return errors.New("促销不存在")
	}
	return err
}

func (svc *PromotionService) BindProducts(ctx *gin.Context, req BindPromotionProductsReq) error {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return errors.New("促销ID非法")
	}
	pids := make([]primitive.ObjectID, 0, len(req.ProductIds))
	for _, s := range req.ProductIds {
		id, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			return errors.New("商品ID非法")
		}
		pids = append(pids, id)
	}
	err = svc.repo.BindProducts(ctx, oid, pids)
	if err == mongo.ErrNoDocuments {
		return errors.New("促销不存在")
	}
	return err
}

func (svc *PromotionService) List(ctx context.Context, req ginx.Page) (domain.PageResult[domain.Promotion], error) {
	if !req.Check() {
		return domain.PageResult[domain.Promotion]{}, errors.New("分页参数非法")
	}
	return svc.repo.List(ctx, req)
}

func (svc *PromotionService) FindById(ctx context.Context, id string) (domain.Promotion, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Promotion{}, errors.New("促销ID非法")
	}
	return svc.repo.FindById(ctx, oid)
}
