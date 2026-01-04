package dao

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BillDaoInterface interface {
	Insert(ctx context.Context, doc any) error
	PageByUid(ctx context.Context, uid any, page int, size int, typ string) ([]any, int64, error)
}

type BillDao struct {
	col *mongo.Collection
}

func NewBillDao(db *mongo.Client) BillDaoInterface {
	database := viper.GetString("mongo.database")
	return &BillDao{
		col: db.Database(database).Collection("bill"),
	}
}

func (d *BillDao) Insert(ctx context.Context, doc any) error {
	_, err := d.col.InsertOne(ctx, doc)
	return err
}

func (d *BillDao) PageByUid(ctx context.Context, uid any, page int, size int, typ string) ([]any, int64, error) {
	filter := bson.M{"uid": uid}
	if typ != "" {
		filter["type"] = typ
	}

	total, err := d.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * size)
	limit := int64(size)

	cur, err := d.col.Find(ctx, filter, options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var res []any
	for cur.Next(ctx) {
		var m bson.M
		if err = cur.Decode(&m); err != nil {
			return nil, 0, err
		}
		res = append(res, m)
	}
	return res, total, cur.Err()
}
