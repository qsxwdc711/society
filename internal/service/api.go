package service

import (
	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type ApiServiceInterface interface {
	GetAllApis(ctx *gin.Context) ([]domain.Api, int64, error)
	PageGet(ctx *gin.Context, page ginx.Page, url string, method string) ([]domain.Api, int64, error)
}
type ApiService struct {
	repo repository.ApiRepositoryInterface
}

func NewApiService(repo repository.ApiRepositoryInterface) ApiServiceInterface {
	return &ApiService{repo: repo}
}

func (a *ApiService) GetAllApis(ctx *gin.Context) ([]domain.Api, int64, error) {
	return a.repo.GetAllApis(ctx)
}

func (a *ApiService) PageGet(c *gin.Context, page ginx.Page, url string, method string) ([]domain.Api, int64, error) {
	return a.repo.PageGet(c, page, url, method)
}
