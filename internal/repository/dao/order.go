package dao

import (
	"time"

	"society/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductId     primitive.ObjectID `bson:"productId"`
	TitleSnapshot string             `bson:"titleSnapshot"`
	PriceSnapshot int64              `bson:"priceSnapshot"`
	Qty           int64              `bson:"qty"`
}

type Order struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Uid       primitive.ObjectID `bson:"uid"`
	Items     []OrderItem        `bson:"items"`
	Amount    int64              `bson:"amount"`
	Status    domain.OrderStatus `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
