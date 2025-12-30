package dao

import (
	"context"
	"errors"
	"time"

	"society/internal/domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrCartEmpty          = errors.New("购物车为空")
	ErrStockNotEnough     = errors.New("库存不足")
	ErrOrderNotFound      = errors.New("订单不存在")
	ErrOrderStatusInvalid = errors.New("订单状态不允许该操作")
	ErrBalanceNotEnough   = errors.New("余额不足")
)

type UserOrderDaoInterface interface {
	CreateOrderFromCartTx(ctx context.Context, uid primitive.ObjectID, productIds []primitive.ObjectID) (domain.Order, error)
	PayOrderTx(ctx context.Context, uid primitive.ObjectID, orderId primitive.ObjectID) (domain.Order, error)
}

type UserOrderDao struct {
	client    *mongo.Client
	db        *mongo.Database
	cartCol   *mongo.Collection
	prodCol   *mongo.Collection
	orderCol  *mongo.Collection
	walletCol *mongo.Collection
	billCol   *mongo.Collection
}

func NewUserOrderDao(client *mongo.Client) UserOrderDaoInterface {
	database := viper.GetString("mongo.database")
	db := client.Database(database)
	return &UserOrderDao{
		client:    client,
		db:        db,
		cartCol:   db.Collection("cart"),
		prodCol:   db.Collection("product"),
		orderCol:  db.Collection("order"),
		walletCol: db.Collection("wallet"),
		billCol:   db.Collection("bill"),
	}
}

type cartItem struct {
	Uid       primitive.ObjectID `bson:"uid"`
	ProductId primitive.ObjectID `bson:"productId"`
	Quantity  int64              `bson:"quantity"`
}

type product struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title"`
	Cover  string             `bson:"cover"`
	Price  int64              `bson:"price"`
	Stock  int64              `bson:"stock"`
	Status string             `bson:"status"`
}

func (dao *UserOrderDao) CreateOrderFromCartTx(ctx context.Context, uid primitive.ObjectID, productIds []primitive.ObjectID) (domain.Order, error) {
	sess, err := dao.client.StartSession()
	if err != nil {
		return domain.Order{}, err
	}
	defer sess.EndSession(ctx)

	var out domain.Order

	_, err = sess.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// 1) 读取购物车（全部 or 指定商品）
		filter := bson.M{"uid": uid}
		if len(productIds) > 0 {
			filter["productId"] = bson.M{"$in": productIds}
		}

		cur, err := dao.cartCol.Find(sc, filter)
		if err != nil {
			return nil, err
		}
		var carts []cartItem
		if err := cur.All(sc, &carts); err != nil {
			return nil, err
		}
		if len(carts) == 0 {
			return nil, ErrCartEmpty
		}

		// 2) 校验商品 + 扣库存（带条件，防超卖）+ 构建快照 items
		items := make([]domain.OrderItem, 0, len(carts))
		var total int64

		for _, ci := range carts {
			if ci.Quantity <= 0 {
				continue
			}

			var p product
			// 只允许购买上架商品
			if err := dao.prodCol.FindOne(sc, bson.M{"_id": ci.ProductId, "status": "ON"}).Decode(&p); err != nil {
				return nil, err
			}

			// 扣库存：stock >= qty
			res, err := dao.prodCol.UpdateOne(
				sc,
				bson.M{"_id": p.Id, "stock": bson.M{"$gte": ci.Quantity}},
				bson.M{"$inc": bson.M{"stock": -ci.Quantity}},
			)
			if err != nil {
				return nil, err
			}
			if res.MatchedCount == 0 {
				return nil, ErrStockNotEnough
			}

			amount := p.Price * ci.Quantity
			total += amount

			items = append(items, domain.OrderItem{
				ProductId:     p.Id,
				TitleSnapshot: p.Title,
				CoverSnapshot: p.Cover,
				PriceSnapshot: p.Price,
				Qty:           ci.Quantity,
				Amount:        amount,
			})
		}

		if len(items) == 0 {
			return nil, ErrCartEmpty
		}

		// 3) 创建订单（你的状态：CREATED）
		now := time.Now()
		order := domain.Order{
			Id:        primitive.NewObjectID(),
			Uid:       uid,
			Items:     items,
			Amount:    total,
			Status:    domain.OrderCreated,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = dao.orderCol.InsertOne(sc, order)
		if err != nil {
			return nil, err
		}

		// 4) 清理已结算购物车项
		_, err = dao.cartCol.DeleteMany(sc, filter)
		if err != nil {
			return nil, err
		}

		out = order
		return nil, nil
	}, options.Transaction())

	return out, err
}

func (dao *UserOrderDao) PayOrderTx(ctx context.Context, uid primitive.ObjectID, orderId primitive.ObjectID) (domain.Order, error) {
	sess, err := dao.client.StartSession()
	if err != nil {
		return domain.Order{}, err
	}
	defer sess.EndSession(ctx)

	var out domain.Order

	_, err = sess.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// 1) 读取订单（必须本人 + CREATED 才能支付）
		var order domain.Order
		if err := dao.orderCol.FindOne(sc, bson.M{"_id": orderId, "uid": uid}).Decode(&order); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrOrderNotFound
			}
			return nil, err
		}
		if order.Status != domain.OrderCreated {
			return nil, ErrOrderStatusInvalid
		}

		now := time.Now()

		// 2) 确保钱包存在（upsert balance=0）
		_, err := dao.walletCol.UpdateOne(
			sc,
			bson.M{"uid": uid},
			bson.M{
				"$setOnInsert": bson.M{"uid": uid, "balance": int64(0)},
				"$set":         bson.M{"updatedAt": now},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return nil, err
		}

		// 3) 扣余额（balance >= order.Amount）
		res, err := dao.walletCol.UpdateOne(
			sc,
			bson.M{"uid": uid, "balance": bson.M{"$gte": order.Amount}},
			bson.M{
				"$inc": bson.M{"balance": -order.Amount},
				"$set": bson.M{"updatedAt": now},
			},
		)
		if err != nil {
			return nil, err
		}
		if res.MatchedCount == 0 {
			return nil, ErrBalanceNotEnough
		}

		// 4) 写账单（支付：负数）
		oid := order.Id
		_, err = dao.billCol.InsertOne(sc, domain.Bill{
			Uid:       uid,
			Type:      domain.BillTypePay,
			Amount:    -order.Amount,
			OrderId:   &oid,
			Remark:    "订单支付",
			CreatedAt: now,
		})
		if err != nil {
			return nil, err
		}

		// 5) 更新订单状态为 PAID + paidAt + updatedAt
		paidAt := time.Now()
		_, err = dao.orderCol.UpdateOne(
			sc,
			bson.M{"_id": order.Id},
			bson.M{"$set": bson.M{
				"status":    domain.OrderPaid,
				"paidAt":    paidAt,
				"updatedAt": paidAt,
			}},
		)
		if err != nil {
			return nil, err
		}

		order.Status = domain.OrderPaid
		order.PaidAt = &paidAt
		order.UpdatedAt = paidAt
		out = order
		return nil, nil
	}, options.Transaction())

	return out, err
}
