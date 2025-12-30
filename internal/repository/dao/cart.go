package dao

import (
	"context"
	"errors"
	"time"

	"society/internal/domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrCartNotFound = errors.New("购物车商品不存在")

type CartDaoInterface interface {
	Add(ctx context.Context, uid, pid primitive.ObjectID) error
	UpdateQty(ctx context.Context, uid, pid primitive.ObjectID, qty int64) error
	Remove(ctx context.Context, uid, pid primitive.ObjectID) error
	List(ctx context.Context, uid primitive.ObjectID) ([]domain.CartProductItem, int64, int64, error)
}

type CartDao struct {
	cartColl *mongo.Collection
}

func NewCartDao(db *mongo.Client) CartDaoInterface {
	database := viper.GetString("mongo.database")
	return &CartDao{
		cartColl: db.Database(database).Collection("cart"),
	}
}

func (dao *CartDao) Add(ctx context.Context, uid, pid primitive.ObjectID) error {
	now := time.Now()

	_, err := dao.cartColl.UpdateOne(
		ctx,
		bson.M{"uid": uid, "productId": pid},
		bson.M{
			"$inc": bson.M{"quantity": 1},
			"$set": bson.M{"updatedAt": now},
			"$setOnInsert": bson.M{
				"uid":       uid,
				"productId": pid,
				"createdAt": now,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (dao *CartDao) UpdateQty(ctx context.Context, uid, pid primitive.ObjectID, qty int64) error {
	if qty <= 0 {
		return dao.Remove(ctx, uid, pid)
	}
	res, err := dao.cartColl.UpdateOne(
		ctx,
		bson.M{"uid": uid, "productId": pid},
		bson.M{
			"$set": bson.M{
				"quantity":  qty,
				"updatedAt": time.Now(),
			},
		},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrCartNotFound
	}
	return nil
}

func (dao *CartDao) Remove(ctx context.Context, uid, pid primitive.ObjectID) error {
	_, err := dao.cartColl.DeleteOne(ctx, bson.M{"uid": uid, "productId": pid})
	return err
}

func (dao *CartDao) List(ctx context.Context, uid primitive.ObjectID) ([]domain.CartProductItem, int64, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"uid": uid}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "product",
			"localField":   "productId",
			"foreignField": "_id",
			"as":           "product",
		}}},
		bson.D{{Key: "$unwind", Value: "$product"}},
		bson.D{{Key: "$match", Value: bson.M{"product.status": "ON"}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":      0,
			"product":  "$product",
			"quantity": "$quantity",
			"amount": bson.M{
				"$multiply": []interface{}{"$product.price", "$quantity"},
			},
		}}},
	}

	cur, err := dao.cartColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cur.Close(ctx)

	var list []domain.CartProductItem
	var totalQty, totalPrice int64
	for cur.Next(ctx) {
		var item domain.CartProductItem
		if err := cur.Decode(&item); err != nil {
			return nil, 0, 0, err
		}
		list = append(list, item)
		totalQty += item.Quantity
		totalPrice += item.Amount
	}
	return list, totalQty, totalPrice, cur.Err()
}
