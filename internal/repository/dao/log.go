package dao

import (
	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogDaoInterface interface {
	Get(ctx *gin.Context, req ginx.Page, level string) ([]Log, int64, error)
}

type LogDao struct {
	db *mongo.Collection
}

func NewLogDao(db *mongo.Client) LogDaoInterface {
	database := viper.GetString("mongo.database")
	return &LogDao{
		db: db.Database(database).Collection("log"),
	}
}

func (l *LogDao) Get(ctx *gin.Context, req ginx.Page, level string) ([]Log, int64, error) {
	opts := req.GetOpts().SetSort(bson.M{"level": -1})
	filter := bson.M{"level": level}
	cur, err := l.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var logs []Log
	err = cur.All(ctx, &logs)
	if err != nil {
		return nil, 0, err
	}
	count, err := l.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return logs, count, nil
}

type Log struct {
	Level  string `bson:"level" json:"level"`
	Time   string `bson:"time" json:"ts"`
	Caller string `bson:"caller" json:"caller"`
	Msg    string `bson:"msg" json:"msg"`
	Route  string `bson:"route" json:"route"`
	Error  string `bson:"error" json:"error"`
}
