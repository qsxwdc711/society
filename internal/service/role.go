package service

import (
	"context"
	"society/internal/domain"
	"society/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserRoleInterface interface {
	Add(ctx *gin.Context, role domain.Role) error
	Get(ctx context.Context) ([]domain.Role, error)
	Update(ctx *gin.Context, req domain.Role) error
}
type RoleService struct {
	repo repository.RoleRepoInterface
}

func NewRoleService(repo repository.RoleRepoInterface) UserRoleInterface {
	return &RoleService{
		repo: repo,
	}
}

func (svc *RoleService) Add(ctx *gin.Context, role domain.Role) error {
	return svc.repo.Add(ctx, role)
}
func (svc *RoleService) Get(ctx context.Context) ([]domain.Role, error) {
	return svc.repo.Get(ctx)
}

func (svc *RoleService) Update(ctx *gin.Context, role domain.Role) error {
	return svc.repo.UpdateRole(ctx, role)
}
