package web

import (
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type UserCartHandler struct {
	svc service.CartInterface
}

func NewUserCartHandler(svc service.CartInterface) *UserCartHandler {
	return &UserCartHandler{svc: svc}
}

func (h *UserCartHandler) RegisterUserCartRouters(server *gin.Engine) {
	user := server.Group("/api/v1/user")
	{
		user.POST("/cart/:productId", ginx.WrapToken[*domain.UserClaims](h.Add))
		user.PUT("/cart/:productId", ginx.WrapBodyAndToken[UpdateCartQtyReq, *domain.UserClaims](h.UpdateQty))
		user.DELETE("/cart/:productId", ginx.WrapToken[*domain.UserClaims](h.Remove))
		user.GET("/cart", ginx.WrapToken[*domain.UserClaims](h.List))
	}
}

type UpdateCartQtyReq struct {
	Quantity int64 `json:"quantity" binding:"required"`
}

func (h *UserCartHandler) Add(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	if err := h.svc.Add(ctx.Request.Context(), uc.Uid, ctx.Param("productId")); err != nil {
		return ginx.ErrorMess("加入购物车失败", err.Error()), err
	}
	return ginx.SuccessMess("加入购物车成功", nil), nil
}

func (h *UserCartHandler) UpdateQty(ctx *gin.Context, req UpdateCartQtyReq, uc *domain.UserClaims) (ginx.Response, error) {
	if err := h.svc.UpdateQty(ctx.Request.Context(), uc.Uid, ctx.Param("productId"), req.Quantity); err != nil {
		return ginx.ErrorMess("更新数量失败", err.Error()), err
	}
	return ginx.SuccessMess("更新数量成功", nil), nil
}

func (h *UserCartHandler) Remove(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	if err := h.svc.Remove(ctx.Request.Context(), uc.Uid, ctx.Param("productId")); err != nil {
		return ginx.ErrorMess("移除商品失败", err.Error()), err
	}
	return ginx.SuccessMess("移除商品成功", nil), nil
}

func (h *UserCartHandler) List(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	data, err := h.svc.List(ctx.Request.Context(), uc.Uid)
	if err != nil {
		return ginx.ErrorMess("查询购物车失败", err.Error()), err
	}
	return ginx.SuccessMess("查询购物车成功", data), nil
}
