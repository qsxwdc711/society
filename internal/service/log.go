package service

import (
	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type LogServiceInterface interface {
	Get(ctx *gin.Context, req ginx.Page, level string) ([]domain.Log, int64, error)
}
type LogService struct {
	repo repository.LogRepositoryInterface
}

func NewLogService(repo repository.LogRepositoryInterface) LogServiceInterface {
	return &LogService{
		repo: repo,
	}
}

func (svc *LogService) Get(ctx *gin.Context, req ginx.Page, level string) ([]domain.Log, int64, error) {
	return svc.repo.Get(ctx, req, level)
}
