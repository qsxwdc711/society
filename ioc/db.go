package ioc

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongodb() *mongo.Client {
	uri := viper.GetString("mongo.uri")
	if uri == "" {
		panic("mongo.uri is empty")
	}

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		panic(err)
	}

	// 可选：启动时直接 ping，提前暴露配置错误
	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}

	return client
}
