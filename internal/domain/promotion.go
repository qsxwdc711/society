package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromotionType string

const (
	PromotionTypeNone     PromotionType = "NONE"
	PromotionTypeReduce   PromotionType = "REDUCE"   // 满减/立减（示例）
	PromotionTypeDiscount PromotionType = "DISCOUNT" // 折扣（示例）
)

type Promotion struct {
	Id         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Type       PromotionType        `json:"type" bson:"type"`
	Desc       string               `json:"desc" bson:"desc"`
	ProductIds []primitive.ObjectID `json:"productIds" bson:"productIds"`
	StartAt    time.Time            `json:"startAt" bson:"startAt"`
	EndAt      time.Time            `json:"endAt" bson:"endAt"`
	CreatedAt  time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time            `json:"updatedAt" bson:"updatedAt"`
}
