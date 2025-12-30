package dao

import (
	"context"
	"time"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"society/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginLogDaoInterface interface {
	Add(ctx *gin.Context, log domain.LoginLog) error
	List(ctx context.Context, role domain.LoginRole, page ginx.Page) ([]domain.LoginLog, int64, error)
}

type LoginLogDao struct {
	db *mongo.Collection
}

func NewLoginLogDao(db *mongo.Client) LoginLogDaoInterface {
	database := viper.GetString("mongo.database")
	return &LoginLogDao{db: db.Database(database).Collection("login_log")}
}

func (dao *LoginLogDao) Add(ctx *gin.Context, log domain.LoginLog) error {
	log.CreatedAt = time.Now()
	_, err := dao.db.InsertOne(ctx, log)
	return err
}

func (dao *LoginLogDao) List(ctx context.Context, role domain.LoginRole, page ginx.Page) ([]domain.LoginLog, int64, error) {
	filter := bson.M{"role": role}
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

	res := make([]domain.LoginLog, 0, page.Size)
	for cur.Next(ctx) {
		var l domain.LoginLog
		if err := cur.Decode(&l); err != nil {
			return nil, 0, err
		}
		res = append(res, l)
	}
	return res, total, cur.Err()
}
