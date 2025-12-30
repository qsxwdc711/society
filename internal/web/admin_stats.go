package web

import (
	"strconv"

	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminStatsHandler struct {
	svc service.StatsInterface
}

func NewAdminStatsHandler(svc service.StatsInterface) *AdminStatsHandler {
	return &AdminStatsHandler{svc: svc}
}

func (h *AdminStatsHandler) RegisterAdminStatsRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.GET("/stats/products/sales", ginx.Wrap(h.SalesRank))
		admin.GET("/stats/products/views", ginx.Wrap(h.ViewRank))
	}
}

func (h *AdminStatsHandler) SalesRank(ctx *gin.Context) (ginx.Response, error) {
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	res, err := h.svc.ProductSalesRank(ctx.Request.Context(), limit)
	if err != nil {
		return ginx.ErrorMess("查询销售排行失败", err.Error()), err
	}
	return ginx.SuccessMess("查询销售排行成功", res), nil
}

func (h *AdminStatsHandler) ViewRank(ctx *gin.Context) (ginx.Response, error) {
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	days, _ := strconv.ParseInt(ctx.Query("days"), 10, 64)
	res, err := h.svc.ProductViewRank(ctx.Request.Context(), limit, days)
	if err != nil {
		return ginx.ErrorMess("查询访客排行失败", err.Error()), err
	}
	return ginx.SuccessMess("查询访客排行成功", res), nil
}
