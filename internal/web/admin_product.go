package web

import (
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminProductHandler struct {
	svc service.ProductAdminInterface
}

func NewAdminProductHandler(svc service.ProductAdminInterface) *AdminProductHandler {
	return &AdminProductHandler{svc: svc}
}

func (h *AdminProductHandler) RegisterAdminProductRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.POST("/products", ginx.WrapBody[service.AddProductReq](h.Add))
		admin.PUT("/products/:id", ginx.WrapBody[service.UpdateProductReq](h.Update))
		admin.PUT("/products/:id/status", ginx.WrapBody[updateProductStatusReq](h.UpdateStatus))
		admin.DELETE("/products/:id", ginx.Wrap(h.Delete))
		admin.GET("/products/:id", ginx.Wrap(h.FindById))
		admin.GET("/products", ginx.Wrap(h.List))
	}
}

func (h *AdminProductHandler) Add(ctx *gin.Context, req service.AddProductReq) (ginx.Response, error) {
	_, err := h.svc.Add(ctx, req)
	if err != nil {
		return ginx.ErrorMess("新增商品失败", err.Error()), err
	}
	return ginx.SuccessMess("新增商品成功", nil), nil
}

func (h *AdminProductHandler) Update(ctx *gin.Context, req service.UpdateProductReq) (ginx.Response, error) {
	req.Id = ctx.Param("id")
	err := h.svc.Update(ctx, req)
	if err != nil {
		return ginx.ErrorMess("更新商品失败", err.Error()), err
	}
	return ginx.SuccessMess("更新商品成功", nil), nil
}

type updateProductStatusReq struct {
	Status domain.ProductStatus `json:"status" binding:"required"`
}

func (h *AdminProductHandler) UpdateStatus(ctx *gin.Context, req updateProductStatusReq) (ginx.Response, error) {
	err := h.svc.UpdateStatus(ctx, ctx.Param("id"), req.Status)
	if err != nil {
		return ginx.ErrorMess("更新商品状态失败", err.Error()), err
	}
	return ginx.SuccessMess("更新商品状态成功", nil), nil
}

func (h *AdminProductHandler) Delete(ctx *gin.Context) (ginx.Response, error) {
	err := h.svc.Delete(ctx, ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("删除商品失败", err.Error()), err
	}
	return ginx.SuccessMess("删除商品成功", nil), nil
}

func (h *AdminProductHandler) FindById(ctx *gin.Context) (ginx.Response, error) {
	p, err := h.svc.FindById(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("查询商品失败", err.Error()), err
	}
	return ginx.SuccessMess("查询商品成功", p), nil
}

func (h *AdminProductHandler) List(ctx *gin.Context) (ginx.Response, error) {
	var req service.ListProductReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	if !req.Page.Check() {
		return ginx.ErrorMess("分页参数非法", nil), nil
	}
	res, err := h.svc.List(ctx.Request.Context(), req)
	if err != nil {
		return ginx.ErrorMess("查询商品列表失败", err.Error()), err
	}
	return ginx.SuccessMess("查询商品列表成功", res), nil
}
