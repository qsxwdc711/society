package web

import (
	"fmt"
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type UserWalletHandler struct {
	svc service.WalletServiceInterface
}

func NewUserWalletHandler(svc service.WalletServiceInterface) *UserWalletHandler {
	return &UserWalletHandler{svc: svc}
}

func (h *UserWalletHandler) RegisterUserWalletRouters(server *gin.Engine) {
	g := server.Group("/api/v1/user/wallet")
	{
		g.GET("", ginx.WrapToken[*domain.UserClaims](h.GetWallet))
		g.POST("/topup", ginx.WrapBodyAndToken[TopUpReq, *domain.UserClaims](h.TopUp))
		g.POST("/transfer", ginx.WrapBodyAndToken[TransferReq, *domain.UserClaims](h.Transfer))
		g.GET("/bills", ginx.WrapToken[*domain.UserClaims](h.Bills))
	}
}

func (h *UserWalletHandler) GetWallet(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	w, err := h.svc.GetWallet(ctx.Request.Context(), uc.Uid)
	if err != nil {
		return ginx.ErrorMess("查询钱包失败", err.Error()), err
	}
	return ginx.SuccessMess("查询钱包成功", w), nil
}

type TopUpReq struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

func (h *UserWalletHandler) TopUp(ctx *gin.Context, req TopUpReq, uc *domain.UserClaims) (ginx.Response, error) {
	w, err := h.svc.TopUp(ctx.Request.Context(), uc.Uid, req.Amount)
	if err != nil {
		return ginx.ErrorMess("充值失败", err.Error()), err
	}
	return ginx.SuccessMess("充值成功", w), nil
}

type TransferReq struct {
	ToPhone string `json:"toPhone" binding:"required"`
	Amount  int64  `json:"amount" binding:"required,gt=0"`
}

func (h *UserWalletHandler) Transfer(ctx *gin.Context, req TransferReq, uc *domain.UserClaims) (ginx.Response, error) {
	w, err := h.svc.Transfer(ctx.Request.Context(), uc.Uid, req.ToPhone, req.Amount)
	if err != nil {
		return ginx.ErrorMess("转账失败", err.Error()), err
	}
	return ginx.SuccessMess("转账成功", w), nil
}

func (h *UserWalletHandler) Bills(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	page := 1
	size := 10
	typ := ctx.Query("type")

	// 你也可以换成 ginx.Page 的方式；这里为了“最少改动”直接读 query
	if v := ctx.Query("page"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &page)
	}
	if v := ctx.Query("size"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &size)
	}

	data, err := h.svc.Bills(ctx.Request.Context(), uc.Uid, page, size, typ)
	if err != nil {
		return ginx.ErrorMess("查询账单失败", err.Error()), err
	}
	return ginx.SuccessMess("查询账单成功", data), nil
}
