package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"` // id
	Phone    string             `json:"phone"`
	Name     string             `json:"name"`
	Password string             `json:"password"`
	Age      int32              `json:"age"`
	Gender   string             `json:"gender"`
	Avatar   string             `json:"avatar"`         // 头像
	Role     primitive.ObjectID `json:"role,omitempty"` // 角色外键
	Token    string             `json:"token,omitempty" bson:"-"`
}
