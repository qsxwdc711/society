package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	svc service.AdminOrderInterface
}

func NewAdminOrderHandler(svc service.AdminOrderInterface) *AdminOrderHandler {
	return &AdminOrderHandler{svc: svc}
}

func (h *AdminOrderHandler) RegisterAdminOrderRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.GET("/orders", ginx.Wrap(h.List))
		admin.PUT("/orders/:id/status", ginx.WrapBody[service.UpdateOrderStatusReq](h.UpdateStatus))
	}
}

func (h *AdminOrderHandler) List(ctx *gin.Context) (ginx.Response, error) {
	var req service.AdminOrderListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	if !req.Page.Check() {
		return ginx.ErrorMess("分页参数非法", nil), nil
	}
	res, err := h.svc.List(ctx.Request.Context(), req)
	if err != nil {
		return ginx.ErrorMess("查询订单失败", err.Error()), err
	}
	return ginx.SuccessMess("查询订单成功", res), nil
}

func (h *AdminOrderHandler) UpdateStatus(ctx *gin.Context, req service.UpdateOrderStatusReq) (ginx.Response, error) {
	req.Id = ctx.Param("id")
	err := h.svc.UpdateStatus(ctx.Request.Context(), req)
	if err != nil {
		return ginx.ErrorMess("更新订单状态失败", err.Error()), err
	}
	return ginx.SuccessMess("更新订单状态成功", nil), nil
}
