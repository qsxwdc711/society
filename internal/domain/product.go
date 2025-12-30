package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryId primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	Title      string             `json:"title" bson:"title"`
	Intro      string             `json:"intro" bson:"intro"`
	Cover      string             `json:"cover" bson:"cover"`
	Images     []string           `json:"images" bson:"images"`
	Price      int64              `json:"price" bson:"price"`
	Stock      int64              `json:"stock" bson:"stock"`
	Status     ProductStatus      `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}
type ProductStatus string

const (
	ProductStatusOn  ProductStatus = "ON"
	ProductStatusOff ProductStatus = "OFF"
)
