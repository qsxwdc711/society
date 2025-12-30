package web

import (
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type UserOrderHandler struct {
	svc service.UserOrderInterface
}

func NewUserOrderHandler(svc service.UserOrderInterface) *UserOrderHandler {
	return &UserOrderHandler{svc: svc}
}

func (h *UserOrderHandler) RegisterUserOrderRouters(server *gin.Engine) {
	user := server.Group("/api/v1/user")
	{
		user.POST("/orders", ginx.WrapBodyAndToken[service.CreateOrderReq, *domain.UserClaims](h.Create))
		user.POST("/orders/:id/pay", ginx.WrapToken[*domain.UserClaims](h.Pay))
	}
}

func (h *UserOrderHandler) Create(ctx *gin.Context, req service.CreateOrderReq, uc *domain.UserClaims) (ginx.Response, error) {
	order, err := h.svc.CreateFromCart(ctx.Request.Context(), uc.Uid, req)
	if err != nil {
		return ginx.ErrorMess("创建订单失败", err.Error()), err
	}
	return ginx.SuccessMess("创建订单成功", order), nil
}

func (h *UserOrderHandler) Pay(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	order, err := h.svc.Pay(ctx.Request.Context(), uc.Uid, ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("支付失败", err.Error()), err
	}
	return ginx.SuccessMess("支付成功", order), nil
}
