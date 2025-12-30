package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	Id      primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name    string               `bson:"name" json:"name"` // unique e.g. admin, user
	Apis    []Api                `json:"apis"`
	Centers []primitive.ObjectID `json:"centers"`
	Codes   []string             `bson:"code" json:"code"` // e.g. user:read
	Desc    string               `bson:"desc" json:"desc"`
}
