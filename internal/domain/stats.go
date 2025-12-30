package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductSalesRank struct {
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	Amount    int64              `json:"amount" bson:"amount"` // 销售额
	Qty       int64              `json:"qty" bson:"qty"`       // 销量
}

type ProductViewRank struct {
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	Views     int64              `json:"views" bson:"views"`
}
