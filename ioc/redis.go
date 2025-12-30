package ioc

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Address  string
		Password string `yaml:"password"`
		Port     string
	}
	var config Config
	if err := viper.UnmarshalKey("redis", &config); err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Address + ":" + config.Port,
		Password: config.Password,
		DB:       0,
	})
	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	return redisClient
}
