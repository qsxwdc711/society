package dao

import (
	"context"

	"society/internal/domain"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProductDaoInterface interface {
	List(ctx context.Context, q string, categoryId *primitive.ObjectID, sort string, minPrice *int64, maxPrice *int64, page ginx.Page) ([]domain.Product, int64, error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error)
}

type UserProductDao struct {
	db *mongo.Collection
}

func NewUserProductDao(db *mongo.Client) UserProductDaoInterface {
	database := viper.GetString("mongo.database")
	return &UserProductDao{db: db.Database(database).Collection("product")}
}

func (dao *UserProductDao) List(
	ctx context.Context,
	q string,
	categoryId *primitive.ObjectID,
	sort string,
	minPrice *int64,
	maxPrice *int64,
	page ginx.Page,
) ([]domain.Product, int64, error) {

	filter := bson.M{
		"status": "ON", // 用户端只展示上架商品
	}
	if q != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": q, "$options": "i"}},
			{"intro": bson.M{"$regex": q, "$options": "i"}},
		}
	}
	if categoryId != nil {
		filter["categoryId"] = *categoryId
	}
	if minPrice != nil || maxPrice != nil {
		r := bson.M{}
		if minPrice != nil {
			r["$gte"] = *minPrice
		}
		if maxPrice != nil {
			r["$lte"] = *maxPrice
		}
		filter["price"] = r
	}

	total, err := dao.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opt := page.GetOpts()

	// sort: price_asc / price_desc / created_desc
	switch sort {
	case "price_asc":
		opt.Sort = bson.M{"price": 1}
	case "price_desc":
		opt.Sort = bson.M{"price": -1}
	default:
		opt.Sort = bson.M{"createdAt": -1}
	}

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

func (dao *UserProductDao) FindById(ctx context.Context, id primitive.ObjectID) (domain.Product, error) {
	var p domain.Product
	// 用户端详情也只允许看 ON 商品（防止下架仍可访问）
	err := dao.db.FindOne(ctx, bson.M{"_id": id, "status": "ON"}).Decode(&p)
	return p, err
}
