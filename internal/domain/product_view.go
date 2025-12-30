package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductView struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Uid       *primitive.ObjectID `json:"uid" bson:"uid,omitempty"`
	ProductId primitive.ObjectID  `json:"productId" bson:"productId"`
	IP        string              `json:"ip" bson:"ip"`
	UA        string              `json:"ua" bson:"ua"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}
