package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderCreated  OrderStatus = "CREATED"
	OrderPaid     OrderStatus = "PAID"
	OrderShipped  OrderStatus = "SHIPPED"  //上架
	OrderDone     OrderStatus = "DONE"     //下架
	OrderCanceled OrderStatus = "CANCELED" //取消
)

type OrderItem struct {
	ProductId     primitive.ObjectID `json:"productId"`
	TitleSnapshot string             `json:"titleSnapshot"`
	PriceSnapshot int64              `json:"priceSnapshot"`
	Qty           int64              `json:"qty"` //数量
}

type Order struct {
	Id        primitive.ObjectID `json:"id"`
	Uid       primitive.ObjectID `json:"uid"`
	Items     []OrderItem        `json:"items"`
	Amount    int64              `json:"amount"`
	Status    OrderStatus        `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
