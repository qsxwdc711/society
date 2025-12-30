package dao

import (
	"time"

	"society/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductViewDaoInterface interface {
	Add(ctx *gin.Context, v domain.ProductView) error
}

type ProductViewDao struct {
	db *mongo.Collection
}

func NewProductViewDao(db *mongo.Client) ProductViewDaoInterface {
	database := viper.GetString("mongo.database")
	return &ProductViewDao{db: db.Database(database).Collection("product_view")}
}

func (dao *ProductViewDao) Add(ctx *gin.Context, v domain.ProductView) error {
	v.CreatedAt = time.Now()
	_, err := dao.db.InsertOne(ctx, v)
	return err
}
