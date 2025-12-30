package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"society/internal/domain"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mongoWriter struct {
	logColl *mongo.Collection
}

func NewMongoWriter(db *mongo.Client) mongoWriter {
	databaseName := viper.GetString("mongo.database")
	return mongoWriter{
		logColl: db.Database(databaseName).Collection("log"),
	}
}

func InitLogger(db *mongo.Client) domain.Loggers {
	fmt.Println("InitLogger called, db==nil?", db == nil, "GOOS=", runtime.GOOS)

	// 创建 encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	consoleWriter := zapcore.AddSync(os.Stdout)
	mongoWriter := zapcore.AddSync(NewMongoWriter(db))

	// 本地和Linux分别处理
	var core zapcore.Core

	if runtime.GOOS == "linux" {
		// Linux：控制台 + MongoDB
		consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zap.InfoLevel)
		mongoCore := zapcore.NewCore(jsonEncoder, mongoWriter, zap.InfoLevel)

		core = zapcore.NewTee(consoleCore, mongoCore)

	} else {
		// 本地：只输出控制台
		core = zapcore.NewCore(consoleEncoder, consoleWriter, zap.InfoLevel)
	}

	// 创建 logger
	log := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(log)

	// 返回给你的 domain，需要 logInfo 则分开
	return domain.Loggers{
		Logg:    log,
		LogInfo: log, // 统一使用一个即可
	}
}

func (mw mongoWriter) Write(p []byte) (n int, err error) {
	var logMap map[string]interface{}
	if err = json.Unmarshal(p, &logMap); err != nil {
		fmt.Println("mongoWriter unmarshal error:", err, "raw:", string(p))
		return 0, err
	}
	var log logStruct
	if logMap["caller"] != nil {
		log.Caller = logMap["caller"].(string)
	}
	if logMap["error"] != nil {
		log.Error = logMap["error"].(string)
	}
	if logMap["route"] != nil {
		log.Route = logMap["route"].(string)
	}
	if logMap["msg"] != nil {
		log.Msg = logMap["msg"].(string)
	}
	log.Level = logMap["level"].(string)
	log.Time = time.Now().Format("2006-01-02 15:04:05")
	_, err = mw.logColl.InsertOne(context.Background(), log)
	if err != nil {
		fmt.Println("mongoWriter InsertOne error:", err)
		return 0, err
	}
	return len(p), nil
}

type logStruct struct {
	Level  string `bson:"level" json:"level"`
	Time   string `bson:"time" json:"ts"`
	Caller string `bson:"caller" json:"caller"`
	Msg    string `bson:"msg" json:"msg"`
	Route  string `bson:"route" json:"route"`
	Error  string `bson:"error" json:"error"`
}
