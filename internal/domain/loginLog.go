package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRole string

const (
	LoginRoleAdmin LoginRole = "ADMIN"
	LoginRoleUser  LoginRole = "USER"
)

type LoginLog struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID `json:"uid" bson:"uid"`
	Role      LoginRole          `json:"role" bson:"role"`
	IP        string             `json:"ip" bson:"ip"`
	UA        string             `json:"ua" bson:"ua"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
