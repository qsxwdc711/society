package dao

import (
	"context"
	"errors"
	"time"

	"society/internal/domain"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrFavoriteExists = errors.New("已收藏")

type FavoriteDaoInterface interface {
	Add(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error
	Remove(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error
	Check(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) (bool, error)
	List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) ([]domain.FavoriteItem, int64, error)
}

type FavoriteDao struct {
	favColl *mongo.Collection
}

func NewFavoriteDao(db *mongo.Client) FavoriteDaoInterface {
	database := viper.GetString("mongo.database")
	return &FavoriteDao{
		favColl: db.Database(database).Collection("favorite"),
	}
}

func (dao *FavoriteDao) Add(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error {
	// 防重复收藏
	cnt, err := dao.favColl.CountDocuments(ctx, bson.M{"uid": uid, "productId": productId})
	if err != nil {
		return err
	}
	if cnt > 0 {
		return ErrFavoriteExists
	}

	_, err = dao.favColl.InsertOne(ctx, domain.Favorite{
		Uid:       uid,
		ProductId: productId,
		CreatedAt: time.Now(),
	})
	return err
}

func (dao *FavoriteDao) Remove(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) error {
	_, err := dao.favColl.DeleteOne(ctx, bson.M{"uid": uid, "productId": productId})
	return err
}

func (dao *FavoriteDao) Check(ctx context.Context, uid primitive.ObjectID, productId primitive.ObjectID) (bool, error) {
	cnt, err := dao.favColl.CountDocuments(ctx, bson.M{"uid": uid, "productId": productId})
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (dao *FavoriteDao) List(ctx context.Context, uid primitive.ObjectID, page ginx.Page) ([]domain.FavoriteItem, int64, error) {
	filter := bson.M{"uid": uid}

	total, err := dao.favColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// favorite -> lookup product，且只展示 status=ON
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		bson.D{{Key: "$skip", Value: int64((page.Page - 1) * page.Size)}},
		bson.D{{Key: "$limit", Value: int64(page.Size)}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "product",
			"localField":   "productId",
			"foreignField": "_id",
			"as":           "product",
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$product",
			"preserveNullAndEmptyArrays": false,
		}}},
		// 只展示上架商品
		bson.D{{Key: "$match", Value: bson.M{"product.status": "ON"}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":        0,
			"product":    "$product",
			"favoriteAt": "$createdAt",
		}}},
	}

	cur, err := dao.favColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	res := make([]domain.FavoriteItem, 0, page.Size)
	for cur.Next(ctx) {
		var item domain.FavoriteItem
		if err := cur.Decode(&item); err != nil {
			return nil, 0, err
		}
		res = append(res, item)
	}
	return res, total, cur.Err()
}
