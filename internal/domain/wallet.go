package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID `json:"uid" bson:"uid"`
	Balance   int64              `json:"balance" bson:"balance"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type BillType string

const (
	BillTypePay    BillType = "PAY"
	BillTypeTopUp  BillType = "TOPUP"
	BillTypeTrans  BillType = "TRANSFER"
	BillTypeRefund BillType = "REFUND"
)

type Bill struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Uid       primitive.ObjectID  `json:"uid" bson:"uid"`
	Type      BillType            `json:"type" bson:"type"`
	Amount    int64               `json:"amount" bson:"amount"` // 正数入账，负数出账
	OrderId   *primitive.ObjectID `json:"orderId,omitempty" bson:"orderId,omitempty"`
	Remark    string              `json:"remark" bson:"remark"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}
