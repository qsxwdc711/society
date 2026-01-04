package dao

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("余额不足")
)

type WalletDaoInterface interface {
	// 若不存在则创建钱包（balance=0）
	EnsureWallet(ctx context.Context, uid any) error

	// 获取余额
	GetBalance(ctx context.Context, uid any) (int64, error)

	// 变更余额（delta 可正可负）
	// 若 delta<0，会校验余额是否足够
	UpdateBalance(ctx context.Context, uid any, delta int64) (int64, error)
}

type WalletDao struct {
	col *mongo.Collection
}

func NewWalletDao(db *mongo.Client) WalletDaoInterface {
	database := viper.GetString("mongo.database")
	return &WalletDao{
		col: db.Database(database).Collection("wallet"),
	}
}

func (d *WalletDao) EnsureWallet(ctx context.Context, uid any) error {
	// upsert: 如果没有钱包就创建
	_, err := d.col.UpdateOne(
		ctx,
		bson.M{"uid": uid},
		bson.M{
			"$setOnInsert": bson.M{
				"uid":       uid,
				"balance":   int64(0),
				"updatedAt": time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (d *WalletDao) GetBalance(ctx context.Context, uid any) (int64, error) {
	var w struct {
		Balance int64 `bson:"balance"`
	}
	err := d.col.FindOne(ctx, bson.M{"uid": uid}, options.FindOne().SetProjection(bson.M{"balance": 1})).Decode(&w)
	if err == mongo.ErrNoDocuments {
		return 0, ErrWalletNotFound
	}
	if err != nil {
		return 0, err
	}
	return w.Balance, nil
}

func (d *WalletDao) UpdateBalance(ctx context.Context, uid any, delta int64) (int64, error) {
	// 先确保钱包存在（避免 NoDocuments）
	if err := d.EnsureWallet(ctx, uid); err != nil {
		return 0, err
	}

	filter := bson.M{"uid": uid}
	// delta<0 扣款时要校验余额
	if delta < 0 {
		filter["balance"] = bson.M{"$gte": -delta}
	}

	var out struct {
		Balance int64 `bson:"balance"`
	}
	err := d.col.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{
			"$inc": bson.M{"balance": delta},
			"$set": bson.M{"updatedAt": time.Now()},
		},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetProjection(bson.M{"balance": 1}),
	).Decode(&out)

	if err == mongo.ErrNoDocuments && delta < 0 {
		return 0, ErrInsufficientBalance
	}
	if err != nil {
		return 0, err
	}
	return out.Balance, nil
}
