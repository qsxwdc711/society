package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginLogInterface interface {
	Add(ctx *gin.Context, uid primitive.ObjectID, role domain.LoginRole) error
	ListAdmin(ctx context.Context, page ginx.Page) (domain.PageResult[domain.LoginLog], error)
	ListUser(ctx context.Context, page ginx.Page) (domain.PageResult[domain.LoginLog], error)
}

type LoginLogService struct {
	repo repository.LoginLogRepoInterface
}

func NewLoginLogService(repo repository.LoginLogRepoInterface) LoginLogInterface {
	return &LoginLogService{repo: repo}
}

func (svc *LoginLogService) Add(ctx *gin.Context, uid primitive.ObjectID, role domain.LoginRole) error {
	return svc.repo.Add(ctx, domain.LoginLog{
		Uid:  uid,
		Role: role,
		IP:   ctx.ClientIP(),
		UA:   ctx.GetHeader("User-Agent"),
	})
}

func (svc *LoginLogService) ListAdmin(ctx context.Context, page ginx.Page) (domain.PageResult[domain.LoginLog], error) {
	if !page.Check() {
		return domain.PageResult[domain.LoginLog]{}, errors.New("分页参数非法")
	}
	return svc.repo.List(ctx, domain.LoginRoleAdmin, page)
}

func (svc *LoginLogService) ListUser(ctx context.Context, page ginx.Page) (domain.PageResult[domain.LoginLog], error) {
	if !page.Check() {
		return domain.PageResult[domain.LoginLog]{}, errors.New("分页参数非法")
	}
	return svc.repo.List(ctx, domain.LoginRoleUser, page)
}
