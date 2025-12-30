package repository

import (
	"context"
	"society/internal/domain"
	"society/internal/repository/dao"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleRepoInterface interface {
	Add(ctx *gin.Context, role domain.Role) error
	Get(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx *gin.Context, role domain.Role) error
	FindRoleById(c context.Context, roleId primitive.ObjectID) (domain.Role, error)
}

type RoleRepo struct {
	dao dao.RoleDaoInterface
}

func NewRoleRepo(dao dao.RoleDaoInterface) RoleRepoInterface {
	return &RoleRepo{
		dao: dao,
	}
}

func (repo *RoleRepo) FindRoleById(ctx context.Context, id primitive.ObjectID) (domain.Role, error) {
	user, err := repo.dao.FindRoleById(ctx, id)
	return toRoleDomain(user), err
}

func (repo *RoleRepo) Add(ctx *gin.Context, role domain.Role) error {
	return repo.dao.Add(ctx, dao.Role{
		Id:   role.Id,
		Name: role.Name,
		Apis: slice.Map[domain.Api, dao.Api](role.Apis, func(idx int, src domain.Api) dao.Api {
			return dao.Api{
				Url:    src.Url,
				Method: src.Method,
			}
		}),
		Centers: role.Centers,
		Codes:   role.Codes,
		Desc:    role.Desc,
	})
}

func (repo *RoleRepo) Get(ctx context.Context) ([]domain.Role, error) {
	data, err := repo.dao.Get(ctx)
	if err != nil {
		return nil, err
	}
	var roles []domain.Role
	for _, role := range data {
		roles = append(roles, toRoleDomain(role))
	}
	return roles, nil
}

func (repo *RoleRepo) UpdateRole(ctx *gin.Context, role domain.Role) error {
	return repo.dao.Update(ctx, role)
}

//func toRoleDomain(role dao.Role) domain.Role {
//	apis := slice.Map[dao.Api, domain.Api](role.Apis, func(idx int, src dao.Api) domain.Api {
//		return domain.Api{
//			Url:    src.Url,
//			Method: src.Method,
//		}
//	})
//
//	return domain.Role{
//		Id:      role.Id,
//		Name:    role.Name,
//		Apis:    apis,
//		Centers: role.Centers,
//		Codes:   role.Codes,
//		Desc:    role.Desc,
//	}
//}

// toRoleDomain 将 DAO 层的 RoleEntity 转换为 domain.Role
func toRoleDomain(role dao.Role) domain.Role {
	apis := slice.Map[dao.Api, domain.Api](role.Apis, func(idx int, src dao.Api) domain.Api {
		return domain.Api{
			Url:    src.Url,
			Method: src.Method,
		}
	})

	return domain.Role{
		Id:      role.Id,
		Name:    role.Name,
		Apis:    apis,
		Centers: role.Centers,
		Codes:   role.Codes,
		Desc:    role.Desc,
	}
}
