package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Favorite struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID `json:"uid" bson:"uid"`
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// 收藏列表返回：包含商品信息 + 收藏时间
type FavoriteItem struct {
	Product    Product   `json:"product" bson:"product"`
	FavoriteAt time.Time `json:"favoriteAt" bson:"favoriteAt"`
}

// 是否收藏
type FavoriteCheckRes struct {
	IsFavorite bool `json:"isFavorite"`
}
