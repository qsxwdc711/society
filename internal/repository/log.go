package repository

import (
	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type LogRepositoryInterface interface {
	Get(ctx *gin.Context, req ginx.Page, level string) ([]domain.Log, int64, error)
}

type LogRepository struct {
	dao dao.LogDaoInterface
}

func NewLogRepository(dao dao.LogDaoInterface) LogRepositoryInterface {
	return &LogRepository{
		dao: dao,
	}
}

func (repo *LogRepository) Get(ctx *gin.Context, req ginx.Page, level string) ([]domain.Log, int64, error) {
	data, count, err := repo.dao.Get(ctx, req, level)
	logs := slice.Map[dao.Log, domain.Log](data, func(idx int, src dao.Log) domain.Log {
		return domain.Log{
			Level:  src.Level,
			Time:   src.Time,
			Caller: src.Caller,
			Msg:    src.Msg,
			Route:  src.Route,
			Error:  src.Error,
		}
	})
	return logs, count, err

}
