package repository

import (
	"context"

	"society/internal/domain"
	"society/internal/repository/dao"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type LoginLogRepoInterface interface {
	Add(ctx *gin.Context, log domain.LoginLog) error
	List(ctx context.Context, role domain.LoginRole, page ginx.Page) (domain.PageResult[domain.LoginLog], error)
}

type LoginLogRepo struct {
	dao dao.LoginLogDaoInterface
}

func NewLoginLogRepo(dao dao.LoginLogDaoInterface) LoginLogRepoInterface {
	return &LoginLogRepo{dao: dao}
}

func (repo *LoginLogRepo) Add(ctx *gin.Context, log domain.LoginLog) error {
	return repo.dao.Add(ctx, log)
}

func (repo *LoginLogRepo) List(ctx context.Context, role domain.LoginRole, page ginx.Page) (domain.PageResult[domain.LoginLog], error) {
	list, total, err := repo.dao.List(ctx, role, page)
	if err != nil {
		return domain.PageResult[domain.LoginLog]{}, err
	}
	return domain.PageResult[domain.LoginLog]{
		List:  list,
		Total: total,
		Page:  int64(page.Page),
		Size:  int64(page.Size),
	}, nil
}
