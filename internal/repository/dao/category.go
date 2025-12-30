package dao

import (
	"context"
	"errors"
	"time"

	"society/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryDaoInterface interface {
	Add(ctx *gin.Context, c domain.Category) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]domain.Category, error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Category, error)
}

type CategoryDao struct {
	db *mongo.Collection
}

func NewCategoryDao(db *mongo.Client) CategoryDaoInterface {
	database := viper.GetString("mongo.database")
	return &CategoryDao{db: db.Database(database).Collection("category")}
}

func (dao *CategoryDao) Add(ctx *gin.Context, c domain.Category) (primitive.ObjectID, error) {
	count, err := dao.db.CountDocuments(ctx, bson.M{"name": c.Name})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if count != 0 {
		return primitive.NilObjectID, errors.New("分类名已存在")
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	res, err := dao.db.InsertOne(ctx, c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("分类ID类型转换失败")
	}
	return oid, nil
}

func (dao *CategoryDao) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
	set["updatedAt"] = time.Now()
	r, err := dao.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": set})
	if err != nil {
		return err
	}
	if r.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *CategoryDao) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	r, err := dao.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *CategoryDao) List(ctx context.Context) ([]domain.Category, error) {
	cur, err := dao.db.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	res := make([]domain.Category, 0)
	for cur.Next(ctx) {
		var c domain.Category
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, cur.Err()
}

func (dao *CategoryDao) FindById(ctx context.Context, id primitive.ObjectID) (domain.Category, error) {
	var c domain.Category
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	return c, err
}
