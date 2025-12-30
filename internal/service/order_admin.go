package service

import (
	"context"
	"errors"

	"society/internal/domain"
	"society/internal/repository"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminOrderInterface interface {
	List(ctx context.Context, req AdminOrderListReq) (domain.PageResult[domain.Order], error)
	UpdateStatus(ctx context.Context, req UpdateOrderStatusReq) error
}

type AdminOrderListReq struct {
	ginx.Page
	Status string `form:"status"`
}

type UpdateOrderStatusReq struct {
	Id     string             `json:"id"`
	Status domain.OrderStatus `json:"status" binding:"required"` // SHIPPED / CANCELED / DONE
}

type AdminOrderService struct {
	repo repository.AdminOrderRepoInterface
}

func NewAdminOrderService(repo repository.AdminOrderRepoInterface) AdminOrderInterface {
	return &AdminOrderService{repo: repo}
}

func (svc *AdminOrderService) List(ctx context.Context, req AdminOrderListReq) (domain.PageResult[domain.Order], error) {
	if !req.Page.Check() {
		return domain.PageResult[domain.Order]{}, errors.New("分页参数非法")
	}
	var st *domain.OrderStatus
	if req.Status != "" {
		tmp := domain.OrderStatus(req.Status)
		if !isValidOrderStatus(tmp) {
			return domain.PageResult[domain.Order]{}, errors.New("状态参数非法")
		}
		st = &tmp
	}
	return svc.repo.List(ctx, st, req.Page)
}

func (svc *AdminOrderService) UpdateStatus(ctx context.Context, req UpdateOrderStatusReq) error {
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return errors.New("订单ID非法")
	}
	if !isValidOrderStatus(req.Status) {
		return errors.New("状态非法")
	}
	old, err := svc.repo.FindById(ctx, oid)
	if err != nil {
		return err
	}
	if !canTransit(old.Status, req.Status) {
		return errors.New("订单状态流转不合法")
	}
	return svc.repo.UpdateStatus(ctx, oid, req.Status)
}

func isValidOrderStatus(s domain.OrderStatus) bool {
	switch s {
	case domain.OrderCreated, domain.OrderPaid, domain.OrderShipped, domain.OrderDone, domain.OrderCanceled:
		return true
	default:
		return false
	}
}

func canTransit(from, to domain.OrderStatus) bool {
	if from == to {
		return true
	}
	switch from {
	case domain.OrderPaid:
		return to == domain.OrderShipped || to == domain.OrderCanceled
	case domain.OrderShipped:
		return to == domain.OrderDone
	case domain.OrderCreated:
		return to == domain.OrderCanceled
	default:
		return false
	}
}
