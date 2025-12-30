package dao

import (
	"context"
	"errors"
	"society/internal/domain"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleDaoInterface interface {
	Add(ctx *gin.Context, role Role) error
	Get(ctx context.Context) ([]Role, error)
	FindRoleById(ctx context.Context, id primitive.ObjectID) (Role, error)
	Update(ctx *gin.Context, role domain.Role) error
}
type RoleDao struct {
	db *mongo.Collection
}

func NewRoleDao(db *mongo.Client) RoleDaoInterface {
	database := viper.GetString("mongo.database")
	return &RoleDao{
		db: db.Database(database).Collection("role"),
	}
}

func (dao *RoleDao) Add(ctx *gin.Context, role Role) error {
	count, err := dao.db.CountDocuments(ctx, bson.M{"name": role.Name})
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("角色名已存在")
	}
	_, err = dao.db.InsertOne(ctx, role)
	if err != nil {
		return err
	}
	return nil
}

func (dao *RoleDao) FindRoleById(ctx context.Context, id primitive.ObjectID) (Role, error) {
	var role Role
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

func (dao *RoleDao) Get(ctx context.Context) ([]Role, error) {
	cursor, err := dao.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var roles []Role
	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil

}

func (dao *RoleDao) Update(ctx *gin.Context, role domain.Role) error {
	apis := slice.Map[domain.Api, Api](role.Apis, func(idx int, src domain.Api) Api {
		return Api{
			Url:    src.Url,
			Method: src.Method,
		}
	})
	filter := bson.M{"_id": role.Id}
	update := bson.M{
		"$set": bson.M{
			"name":    role.Name,
			"apis":    apis,
			"desc":    role.Desc,
			"centers": role.Centers,
			"codes":   role.Codes,
		},
	}
	result := dao.db.FindOneAndUpdate(ctx, filter, update)
	return result.Err()
}

type Role struct {
	Id      primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string               `json:"name" bson:"name"`
	Apis    []Api                `json:"apis" bson:"apis"`
	Centers []primitive.ObjectID `json:"centers" bson:"centers"`
	Codes   []string             `json:"codes" bson:"codes"`
	Desc    string               `json:"desc" bson:"desc"` //角色描述
}
