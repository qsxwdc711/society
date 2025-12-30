package dao

import (
	"context"
	"time"

	"society/internal/domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsDaoInterface interface {
	ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error)
	ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error)
}

type StatsDao struct {
	order *mongo.Collection
	view  *mongo.Collection
}

func NewStatsDao(db *mongo.Client) StatsDaoInterface {
	database := viper.GetString("mongo.database")
	d := db.Database(database)
	return &StatsDao{
		order: d.Collection("order"),
		view:  d.Collection("product_view"),
	}
}

// 统计销售额：从订单聚合 items，按 productId 汇总 amount/qty
func (dao *StatsDao) ProductSalesRank(ctx context.Context, limit int64) ([]domain.ProductSalesRank, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"status": bson.M{"$in": []string{"PAID", "SHIPPED", "DONE"}},
		}}},
		bson.D{{Key: "$unwind", Value: "$items"}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": "$items.productId",
			"qty": bson.M{"$sum": "$items.qty"},
			"amount": bson.M{"$sum": bson.M{
				"$multiply": []any{"$items.priceSnapshot", "$items.qty"},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"amount": -1}}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id": 0, "productId": "$_id", "qty": 1, "amount": 1,
		}}},
	}

	cur, err := dao.order.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	res := make([]domain.ProductSalesRank, 0, limit)
	for cur.Next(ctx) {
		var r domain.ProductSalesRank
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, cur.Err()
}

// 浏览排行：按 productId 计数
func (dao *StatsDao) ProductViewRank(ctx context.Context, limit int64, days int64) ([]domain.ProductViewRank, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	filter := bson.M{}
	if days > 0 {
		filter["createdAt"] = bson.M{"$gte": time.Now().Add(-time.Duration(days) * 24 * time.Hour)}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$productId",
			"views": bson.M{"$sum": 1},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"views": -1}}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id": 0, "productId": "$_id", "views": 1,
		}}},
	}

	cur, err := dao.view.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	res := make([]domain.ProductViewRank, 0, limit)
	for cur.Next(ctx) {
		var r domain.ProductViewRank
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, cur.Err()
}
