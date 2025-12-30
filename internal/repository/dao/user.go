package dao

import (
	"context"
	"society/internal/domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserDaoInterface interface {
	InsertOne(ctx context.Context, user User) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id primitive.ObjectID) (User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	ChangePassword(ctx context.Context, uid primitive.ObjectID, newPassword string) error
}

type MongoUserDao struct {
	db *mongo.Collection
}

func NewUserDao(mongodb *mongo.Client) UserDaoInterface {
	database := viper.GetString("mongo.database")
	return &MongoUserDao{
		db: mongodb.Database(database).Collection("user"),
	}
}

func (dao *MongoUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		zap.L().Info("dao.FindByPhone not found", zap.String("phone", user.Phone))
		return User{}, err
	}
	zap.L().Info("dao.FindByPhone found", zap.String("phone", user.Phone))
	return user, nil
}

func (dao *MongoUserDao) InsertOne(ctx context.Context, u User) (domain.User, error) {
	res, err := dao.db.InsertOne(ctx, &u)
	if err != nil {
		return domain.User{}, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return domain.User{Id: id}, nil
}

func (dao *MongoUserDao) FindById(ctx context.Context, id primitive.ObjectID) (User, error) {
	//traceID := ctx.Value(middleware.CtxTraceIDKey)
	//traceStr := ""
	//if s, ok := traceID.(string); ok {
	//	traceStr = s
	//}

	zap.L().Info("dao.FindById enter")
	var user User
	err := dao.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		zap.L().Info("dao.FindById error", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func (dao *MongoUserDao) Update(ctx context.Context, user domain.User) (domain.User, error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": bson.M{"avatar": user.Avatar,
		"phone": user.Phone, "name": user.Name, "sex": user.Gender, "role": user.Role}}
	result := dao.db.FindOneAndUpdate(ctx, filter, update)
	// 检查是否出现错误
	if err := result.Err(); err != nil {
		return domain.User{}, err // 返回空的 User 和错误
	}
	// 返回更新后的用户和 nil 错误
	return user, nil
}

func (dao *MongoUserDao) ChangePassword(ctx context.Context, uid primitive.ObjectID, newPassword string) error {
	//traceStr := middleware.GetTraceID(ctx)
	zap.L().Info("dao.ChangePassword enter")
	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"password": newPassword}}
	_, err := dao.db.UpdateOne(ctx, filter, update)
	if err != nil {
		zap.L().Error("dao.ChangePassword failed", zap.Error(err))
		return err
	}
	return nil
}

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"` // id
	Name     string             `json:"name" bson:"name"`        //姓名
	Phone    string             `json:"phone" bson:"phone"`
	Gender   string             `json:"gender" bson:"gender"` // 性别
	Age      int32              `json:"age" bson:"age"`
	Password string             `json:"password" bson:"password"`             // 密码
	Avatar   string             `json:"avatar" bson:"avatar"`                 // 头像
	Role     primitive.ObjectID `json:"role,omitempty" bson:"role,omitempty"` // 角色外键
}
