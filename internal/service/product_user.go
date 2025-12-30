package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProductInterface interface {
	List(ctx context.Context, req UserProductListReq) (domain.PageResult[domain.Product], error)
	Detail(ctx *gin.Context, id string) (domain.Product, error)
}

type UserProductListReq struct {
	ginx.Page
	Q          string `form:"q"`
	CategoryId string `form:"categoryId"`
	Sort       string `form:"sort"` // price_asc/price_desc/created_desc
	MinPrice   *int64 `form:"minPrice"`
	MaxPrice   *int64 `form:"maxPrice"`
}

type UserProductService struct {
	repo     repository.UserProductRepoInterface
	viewRepo repository.ProductViewRepoInterface
}

func NewUserProductService(repo repository.UserProductRepoInterface, viewRepo repository.ProductViewRepoInterface) UserProductInterface {
	return &UserProductService{repo: repo, viewRepo: viewRepo}
}

func (svc *UserProductService) List(ctx context.Context, req UserProductListReq) (domain.PageResult[domain.Product], error) {
	if !req.Page.Check() {
		return domain.PageResult[domain.Product]{}, errors.New("分页参数非法")
	}

	var cid *primitive.ObjectID
	if req.CategoryId != "" {
		tmp, err := primitive.ObjectIDFromHex(req.CategoryId)
		if err != nil {
			return domain.PageResult[domain.Product]{}, errors.New("分类ID非法")
		}
		cid = &tmp
	}

	return svc.repo.List(ctx, req.Q, cid, req.Sort, req.MinPrice, req.MaxPrice, req.Page)
}

func (svc *UserProductService) Detail(ctx *gin.Context, id string) (domain.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Product{}, errors.New("商品ID非法")
	}

	p, err := svc.repo.FindById(ctx.Request.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Product{}, errors.New("商品不存在或已下架")
		}
		return domain.Product{}, err
	}

	// ===== 浏览埋点：不影响主流程，失败也不影响返回 =====
	var uidPtr *primitive.ObjectID
	if val, ok := ctx.Get("claims"); ok {
		if c, ok := val.(*domain.UserClaims); ok {
			uid := c.Uid
			uidPtr = &uid
		}
	}

	_ = svc.viewRepo.Add(ctx, domain.ProductView{
		Uid:       uidPtr,
		ProductId: oid,
		IP:        ctx.ClientIP(),
		UA:        ctx.GetHeader("User-Agent"),
	})
	// =====================================================

	return p, nil
}
