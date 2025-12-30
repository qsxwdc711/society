package repository

import (
	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type ApiRepositoryInterface interface {
	GetAllApis(ctx *gin.Context) ([]domain.Api, int64, error)
	PageGet(ctx *gin.Context, page ginx.Page, url string, method string) ([]domain.Api, int64, error)
}
type ApiRepository struct {
	dao dao.ApiDaoInterface
}

func NewApiRepo(dao dao.ApiDaoInterface) ApiRepositoryInterface {
	return &ApiRepository{dao: dao}
}

func (a *ApiRepository) GetAllApis(ctx *gin.Context) ([]domain.Api, int64, error) {
	apis, count, err := a.dao.FindAllApis(ctx)
	if err != nil {
		return nil, count, err
	}
	return toDomainApis(apis), count, nil
}

func (a *ApiRepository) PageGet(c *gin.Context, page ginx.Page, url string, method string) ([]domain.Api, int64, error) {
	apis, count, err := a.dao.FindByFilter(c, page, url, method)
	if err != nil {
		return nil, count, err
	}
	return toDomainApis(apis), count, nil
}
func toDomainApis(data []dao.Api) []domain.Api {
	return slice.Map[dao.Api, domain.Api](data, func(idx int, src dao.Api) domain.Api {
		return domain.Api{
			Url:    src.Url,
			Method: src.Method,
		}
	})
}
