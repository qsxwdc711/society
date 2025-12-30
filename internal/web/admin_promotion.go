package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminPromotionHandler struct {
	svc service.PromotionInterface
}

func NewAdminPromotionHandler(svc service.PromotionInterface) *AdminPromotionHandler {
	return &AdminPromotionHandler{svc: svc}
}

func (h *AdminPromotionHandler) RegisterAdminPromotionRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.POST("/promotions", ginx.WrapBody[service.AddPromotionReq](h.Add))
		admin.PUT("/promotions/:id", ginx.WrapBody[service.UpdatePromotionReq](h.Update))
		admin.DELETE("/promotions/:id", ginx.Wrap(h.Delete))
		admin.GET("/promotions", ginx.Wrap(h.List))
		admin.GET("/promotions/:id", ginx.Wrap(h.FindById))

		admin.PUT("/promotions/:id/products", ginx.WrapBody[service.BindPromotionProductsReq](h.BindProducts))
	}
}

func (h *AdminPromotionHandler) Add(ctx *gin.Context, req service.AddPromotionReq) (ginx.Response, error) {
	_, err := h.svc.Add(ctx, req)
	if err != nil {
		return ginx.ErrorMess("新增促销失败", err.Error()), err
	}
	return ginx.SuccessMess("新增促销成功", nil), nil
}

func (h *AdminPromotionHandler) Update(ctx *gin.Context, req service.UpdatePromotionReq) (ginx.Response, error) {
	req.Id = ctx.Param("id")
	err := h.svc.Update(ctx, req)
	if err != nil {
		return ginx.ErrorMess("更新促销失败", err.Error()), err
	}
	return ginx.SuccessMess("更新促销成功", nil), nil
}

func (h *AdminPromotionHandler) Delete(ctx *gin.Context) (ginx.Response, error) {
	err := h.svc.Delete(ctx, ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("删除促销失败", err.Error()), err
	}
	return ginx.SuccessMess("删除促销成功", nil), nil
}

func (h *AdminPromotionHandler) List(ctx *gin.Context) (ginx.Response, error) {
	var page ginx.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	if !page.Check() {
		return ginx.ErrorMess("分页参数非法", nil), nil
	}
	res, err := h.svc.List(ctx.Request.Context(), page)
	if err != nil {
		return ginx.ErrorMess("查询促销失败", err.Error()), err
	}
	return ginx.SuccessMess("查询促销成功", res), nil
}

func (h *AdminPromotionHandler) FindById(ctx *gin.Context) (ginx.Response, error) {
	res, err := h.svc.FindById(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("查询促销失败", err.Error()), err
	}
	return ginx.SuccessMess("查询促销成功", res), nil
}

func (h *AdminPromotionHandler) BindProducts(ctx *gin.Context, req service.BindPromotionProductsReq) (ginx.Response, error) {
	req.Id = ctx.Param("id")
	err := h.svc.BindProducts(ctx, req)
	if err != nil {
		return ginx.ErrorMess("绑定促销商品失败", err.Error()), err
	}
	return ginx.SuccessMess("绑定促销商品成功", nil), nil
}
