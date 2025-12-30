package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderCreated  OrderStatus = "CREATED" // 已创建（待支付）
	OrderPaid     OrderStatus = "PAID"
	OrderShipped  OrderStatus = "SHIPPED"
	OrderDone     OrderStatus = "DONE"
	OrderCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	ProductId     primitive.ObjectID `json:"productId" bson:"productId"`
	TitleSnapshot string             `json:"titleSnapshot" bson:"titleSnapshot"`
	CoverSnapshot string             `json:"coverSnapshot" bson:"coverSnapshot"`
	PriceSnapshot int64              `json:"priceSnapshot" bson:"priceSnapshot"`
	Qty           int64              `json:"qty" bson:"qty"`
	Amount        int64              `json:"amount" bson:"amount"` // price * qty
}

type Order struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID `json:"uid" bson:"uid"`
	Items     []OrderItem        `json:"items" bson:"items"`
	Amount    int64              `json:"amount" bson:"amount"` // 总金额
	Status    OrderStatus        `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`

	// 支付相关（可选字段）
	PaidAt *time.Time `json:"paidAt,omitempty" bson:"paidAt,omitempty"`
}
