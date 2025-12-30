package dao

import (
	"context"
	"errors"
	"time"

	"society/internal/domain"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromotionDaoInterface interface {
	Add(ctx *gin.Context, p domain.Promotion) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	BindProducts(ctx *gin.Context, id primitive.ObjectID, productIds []primitive.ObjectID) error
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Promotion, error)
	List(ctx context.Context, page ginx.Page) ([]domain.Promotion, int64, error)
}

type PromotionDao struct {
	db *mongo.Collection
}

func NewPromotionDao(db *mongo.Client) PromotionDaoInterface {
	database := viper.GetString("mongo.database")
	return &PromotionDao{db: db.Database(database).Collection("promotion")}
}

func (dao *PromotionDao) Add(ctx *gin.Context, p domain.Promotion) (primitive.ObjectID, error) {
	count, err := dao.db.CountDocuments(ctx, bson.M{"name": p.Name})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if count != 0 {
		return primitive.NilObjectID, errors.New("促销名称已存在")
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	res, err := dao.db.InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("促销ID类型转换失败")
	}
	return oid, nil
}

func (dao *PromotionDao) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
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

func (dao *PromotionDao) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	r, err := dao.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *PromotionDao) BindProducts(ctx *gin.Context, id primitive.ObjectID, productIds []primitive.ObjectID) error {
	r, err := dao.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"productIds": productIds,
		"updatedAt":  time.Now(),
	}})
	if err != nil {
		return err
	}
	if r.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *PromotionDao) FindById(ctx context.Context, id primitive.ObjectID) (domain.Promotion, error) {
	var p domain.Promotion
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	return p, err
}

func (dao *PromotionDao) List(ctx context.Context, page ginx.Page) ([]domain.Promotion, int64, error) {
	total, err := dao.db.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	opt := page.GetOpts()
	opt.Sort = bson.M{"createdAt": -1}
	cur, err := dao.db.Find(ctx, bson.M{}, opt)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	res := make([]domain.Promotion, 0, page.Size)
	for cur.Next(ctx) {
		var p domain.Promotion
		if err := cur.Decode(&p); err != nil {
			return nil, 0, err
		}
		res = append(res, p)
	}
	return res, total, cur.Err()
}
