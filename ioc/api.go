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
	//这里要删除所有角色的缓存可能有新的api加进来
	//获取配置文件的超级管理员id
	roleId := viper.GetString("role.id")
	if roleId == "" {
		panic("role id is empty")
	}
	id, err := primitive.ObjectIDFromHex(roleId)
	if err != nil {
		panic("获取管理员失败")
	}
	//更新api集合的api
	router := engine.Routes()
	apiColl := db.Collection("api")
	var apis []interface{}
	for _, route := range router {
		api := dao.Api{
			Url:    route.Path,
			Method: route.Method,
		}
		apis = append(apis, api)
	}
	_, err = apiColl.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	_, err = apiColl.InsertMany(context.TODO(), apis)
	if err != nil {
		panic("api数据库初始化失败：" + err.Error())
	}

	//更新指定id超级管理员的api
	//更新超级管理员的api
	roleColl := db.Collection("role")
	_, err = roleColl.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"apis": apis}})
	if err != nil {
		panic("更新超级管理员数据库api失败")
	}
}
