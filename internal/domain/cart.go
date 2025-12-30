package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID `json:"uid" bson:"uid"`
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	Quantity  int64              `json:"quantity" bson:"quantity"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// 返回给前端的购物车商品
type CartProductItem struct {
	Product  Product `json:"product" bson:"product"`
	Quantity int64   `json:"quantity" bson:"quantity"`
	Amount   int64   `json:"amount" bson:"amount"` // price * quantity
}

type CartListRes struct {
	List       []CartProductItem `json:"list"`
	TotalQty   int64             `json:"totalQty"`
	TotalPrice int64             `json:"totalPrice"`
}
