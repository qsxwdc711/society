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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductDaoInterface interface {
	Add(ctx *gin.Context, p domain.Product) (primitive.ObjectID, error)
	Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error
	UpdateStatus(ctx *gin.Context, id primitive.ObjectID, status domain.ProductStatus) error
	Delete(ctx *gin.Context, id primitive.ObjectID) error
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error)
	List(ctx context.Context, q string, categoryId *primitive.ObjectID, status *domain.ProductStatus, page ginx.Page) ([]domain.Product, int64, error)
}

type ProductDao struct {
	db *mongo.Collection
}

func NewProductDao(db *mongo.Client) ProductDaoInterface {
	database := viper.GetString("mongo.database")
	return &ProductDao{db: db.Database(database).Collection("product")}
}

func (dao *ProductDao) Add(ctx *gin.Context, p domain.Product) (primitive.ObjectID, error) {
	count, err := dao.db.CountDocuments(ctx, bson.M{"title": p.Title})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if count != 0 {
		return primitive.NilObjectID, errors.New("商品名称已存在")
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Status == "" {
		p.Status = domain.ProductStatusOff
	}
	res, err := dao.db.InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("商品ID类型转换失败")
	}
	return oid, nil
}

func (dao *ProductDao) Update(ctx *gin.Context, id primitive.ObjectID, set bson.M) error {
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

func (dao *ProductDao) UpdateStatus(ctx *gin.Context, id primitive.ObjectID, status domain.ProductStatus) error {
	r, err := dao.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}})
	if err != nil {
		return err
	}
	if r.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *ProductDao) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	r, err := dao.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (dao *ProductDao) FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error) {
	var p domain.Product
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	return p, err
}

func (dao *ProductDao) List(ctx context.Context, q string, categoryId *primitive.ObjectID, status *domain.ProductStatus, page ginx.Page) ([]domain.Product, int64, error) {
	filter := bson.M{}
	if q != "" {
		filter["title"] = bson.M{"$regex": q, "$options": "i"}
	}
	if categoryId != nil {
		filter["categoryId"] = *categoryId
	}
	if status != nil && *status != "" {
		filter["status"] = *status
	}

	total, err := dao.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opt := page.GetOpts()
	opt.Sort = bson.M{"createdAt": -1}

	cur, err := dao.db.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	res := make([]domain.Product, 0, page.Size)
	for cur.Next(ctx) {
		var p domain.Product
		if err := cur.Decode(&p); err != nil {
			return nil, 0, err
		}
		res = append(res, p)
	}
	return res, total, cur.Err()
}

var _ = options.FindOptions{}
