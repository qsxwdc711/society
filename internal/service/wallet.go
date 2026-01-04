package service

import (
	"context"
	"errors"
	"fmt"
	"society/internal/domain"
	"society/internal/repository"
	"society/internal/repository/dao"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidAmount = errors.New("金额必须大于0")
)

type WalletServiceInterface interface {
	GetWallet(ctx context.Context, uid primitive.ObjectID) (domain.Wallet, error)
	TopUp(ctx context.Context, uid primitive.ObjectID, amount int64) (domain.Wallet, error)
	Transfer(ctx context.Context, from primitive.ObjectID, toPhone string, amount int64) (domain.Wallet, error)
	Bills(ctx context.Context, uid primitive.ObjectID, page int, size int, typ string) (any, error)
}

type WalletService struct {
	walletRepo repository.WalletRepoInterface
	billRepo   repository.BillRepoInterface
	userRepo   repository.UserRepoInterface
	db         *mongo.Client // 事务
}

func NewWalletService(
	walletRepo repository.WalletRepoInterface,
	billRepo repository.BillRepoInterface,
	userRepo repository.UserRepoInterface,
	db *mongo.Client,
) WalletServiceInterface {
	return &WalletService{
		walletRepo: walletRepo,
		billRepo:   billRepo,
		userRepo:   userRepo,
		db:         db,
	}
}

func (s *WalletService) GetWallet(ctx context.Context, uid primitive.ObjectID) (domain.Wallet, error) {
	_ = s.walletRepo.EnsureWallet(ctx, uid)
	bal, err := s.walletRepo.GetBalance(ctx, uid)
	if err != nil && err != dao.ErrWalletNotFound {
		return domain.Wallet{}, err
	}
	return domain.Wallet{
		Uid:       uid,
		Balance:   bal,
		UpdatedAt: time.Now(),
	}, nil
}

func (s *WalletService) TopUp(ctx context.Context, uid primitive.ObjectID, amount int64) (domain.Wallet, error) {
	if amount <= 0 {
		return domain.Wallet{}, ErrInvalidAmount
	}

	session, err := s.db.StartSession()
	if err != nil {
		return domain.Wallet{}, err
	}
	defer session.EndSession(ctx)

	var newBal int64
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		bal, err := s.walletRepo.UpdateBalance(sc, uid, amount) // 正数入账
		if err != nil {
			return nil, err
		}
		newBal = bal

		bill := domain.Bill{
			Uid:       uid,
			Type:      domain.BillTypeTopUp,
			Amount:    amount, // 正数入账
			Remark:    "钱包充值",
			CreatedAt: time.Now(),
		}
		return nil, s.billRepo.Insert(sc, bill)
	})
	if err != nil {
		return domain.Wallet{}, err
	}
	return domain.Wallet{Uid: uid, Balance: newBal, UpdatedAt: time.Now()}, nil
}

func (s *WalletService) Transfer(ctx context.Context, from primitive.ObjectID, toPhone string, amount int64) (domain.Wallet, error) {
	if amount <= 0 {
		return domain.Wallet{}, ErrInvalidAmount
	}

	toUser, err := s.userRepo.FindOneByPhone(ctx, toPhone)
	if err != nil {
		return domain.Wallet{}, err
	}
	to := toUser.Id
	if to == from {
		return domain.Wallet{}, errors.New("不能给自己转账")
	}

	session, err := s.db.StartSession()
	if err != nil {
		return domain.Wallet{}, err
	}
	defer session.EndSession(ctx)

	var fromBal int64
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// 扣款：delta 为负数
		balFrom, err := s.walletRepo.UpdateBalance(sc, from, -amount)
		if err != nil {
			return nil, err
		}
		fromBal = balFrom

		// 入账
		_, err = s.walletRepo.UpdateBalance(sc, to, amount)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		// 你的 Bill.Amount：正入账，负出账
		fromBill := domain.Bill{
			Uid:       from,
			Type:      domain.BillTypeTrans,
			Amount:    -amount,
			Remark:    fmt.Sprintf("转账给 %s", toPhone),
			CreatedAt: now,
		}
		toBill := domain.Bill{
			Uid:       to,
			Type:      domain.BillTypeTrans,
			Amount:    amount,
			Remark:    "收到转账",
			CreatedAt: now,
		}
		if err := s.billRepo.Insert(sc, fromBill); err != nil {
			return nil, err
		}
		if err := s.billRepo.Insert(sc, toBill); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return domain.Wallet{}, err
	}
	return domain.Wallet{Uid: from, Balance: fromBal, UpdatedAt: time.Now()}, nil
}

func (s *WalletService) Bills(ctx context.Context, uid primitive.ObjectID, page int, size int, typ string) (any, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	list, total, err := s.billRepo.PageByUid(ctx, uid, page, size, typ)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}, nil
}
