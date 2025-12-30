package repository

import (
	"context"
	"society/internal/domain"
	"society/internal/repository/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepoInterface interface {
	FindOneByPhone(ctx context.Context, phone string) (domain.User, error)
	InsertOne(ctx context.Context, user domain.User) (domain.User, error)
	FindOneById(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	ChangePassword(ctx context.Context, uid primitive.ObjectID, newPassword string) error
}

type UserRepository struct {
	dao dao.UserDaoInterface
}

func NewUserRepo(dao dao.UserDaoInterface) UserRepoInterface {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) FindOneByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(user), nil
}

func (repo *UserRepository) InsertOne(ctx context.Context, user domain.User) (domain.User, error) {
	return repo.dao.InsertOne(ctx, dao.User{
		Name:     user.Name,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Age:      user.Age,
		Password: user.Password,
		Avatar:   user.Avatar,
		Role:     user.Role,
	})
}

func (repo *UserRepository) FindOneById(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	user, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(user), nil
}

func (repo *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	return repo.dao.Update(ctx, user)
}

func (repo *UserRepository) ChangePassword(ctx context.Context, uid primitive.ObjectID, newPassword string) error {
	return repo.dao.ChangePassword(ctx, uid, newPassword)
}

func toUserDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Name:     user.Name,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Age:      user.Age,
		Password: user.Password,
		Avatar:   user.Avatar,
		Role:     user.Role,
	}
}
