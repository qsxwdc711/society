package ioc

import (
	"context"
	"society/internal/repository/dao"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitApiColl(engine *gin.Engine, db1 *mongo.Client) {
	databaseName := viper.GetString("mongo.database")
	db := db1.Database(databaseName)

	// 超级管理员 role id
	roleId := viper.GetString("role.id")
	if roleId == "" {
		panic("role id is empty")
	}
	id, err := primitive.ObjectIDFromHex(roleId)
	if err != nil {
		panic("获取管理员失败")
	}

	// 生成 API 列表（强类型）
	routes := engine.Routes()
	apiColl := db.Collection("api")

	apis := make([]dao.Api, 0, len(routes))
	for _, route := range routes {
		apis = append(apis, dao.Api{
			Url:    route.Path, // 模板路径：/xxx/:id
			Method: route.Method,
		})
	}

	// 1) 重建 api 集合
	_, err = apiColl.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	// InsertMany 需要 []interface{}
	docs := make([]interface{}, 0, len(apis))
	for _, a := range apis {
		docs = append(docs, a)
	}
	_, err = apiColl.InsertMany(context.TODO(), docs)
	if err != nil {
		panic("api数据库初始化失败：" + err.Error())
	}

	// 2) 更新超级管理员的 apis（必须用强类型 apis）
	roleColl := db.Collection("role")
	_, err = roleColl.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"apis": apis}},
	)
	if err != nil {
		panic("更新超级管理员数据库api失败：" + err.Error())
	}
}
