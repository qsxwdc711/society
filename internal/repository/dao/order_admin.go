package dao

import (
	"context"
	"time"

	"society/internal/domain"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminOrderDaoInterface interface {
	List(ctx context.Context, status *domain.OrderStatus, page ginx.Page) ([]domain.Order, int64, error)
	FindById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error
}

type AdminOrderDao struct {
	db *mongo.Collection
}

func NewAdminOrderDao(db *mongo.Client) AdminOrderDaoInterface {
	database := viper.GetString("mongo.database")
	return &AdminOrderDao{db: db.Database(database).Collection("order")}
}

func (dao *AdminOrderDao) List(ctx context.Context, status *domain.OrderStatus, page ginx.Page) ([]domain.Order, int64, error) {
	filter := bson.M{}
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

	res := make([]domain.Order, 0, page.Size)
	for cur.Next(ctx) {
		var o domain.Order
		if err := cur.Decode(&o); err != nil {
			return nil, 0, err
		}
		res = append(res, o)
	}
	return res, total, cur.Err()
}

func (dao *AdminOrderDao) FindById(ctx context.Context, id primitive.ObjectID) (domain.Order, error) {
	var o domain.Order
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	return o, err
}

func (dao *AdminOrderDao) UpdateStatus(ctx context.Context, id primitive.ObjectID, status domain.OrderStatus) error {
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
