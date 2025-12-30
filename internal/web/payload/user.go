package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserLoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserRegisterReq struct {
	Name     string `json:"name" bson:"name"` //姓名
	Phone    string `json:"phone" bson:"phone"`
	Gender   string `json:"gender" bson:"gender"` // 性别
	Age      int32  `json:"age" bson:"age"`
	Password string `json:"password" bson:"password"` // 密码
	Avatar   string `json:"avatar" bson:"avatar"`     // 头像
	Role     string `json:"role" bson:"role"`
}
type UserUpdateReq struct {
	Id     primitive.ObjectID `json:"id"`               // id
	Name   string             `json:"name" bson:"name"` //姓名
	Phone  string             `json:"phone" bson:"phone"`
	Gender string             `json:"gender" bson:"gender"` // 性别
	Age    int32              `json:"age" bson:"age"`
	Avatar string             `json:"avatar" bson:"avatar"` // 头像
	Role   string             `json:"role" bson:"role"`
}
type UserChangeReq struct {
	OldPass string `json:"oldPass" binding:"required"`
	NewPass string `json:"newPass" binding:"required"`
}
type UserDeleteReq struct {
	Id primitive.ObjectID `json:"id,omitempty"`
}
type UserResetPasswordReq struct {
	Id    primitive.ObjectID `json:"id" binding:"required"`
	Phone string             `json:"phone" binding:"required"`
}
