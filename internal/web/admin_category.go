package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminCategoryHandler struct {
	svc service.CategoryInterface
}

func NewAdminCategoryHandler(svc service.CategoryInterface) *AdminCategoryHandler {
	return &AdminCategoryHandler{svc: svc}
}

func (h *AdminCategoryHandler) RegisterAdminCategoryRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.POST("/categories", ginx.WrapBody[service.AddCategoryReq](h.Add))
		admin.PUT("/categories/:id", ginx.WrapBody[service.UpdateCategoryReq](h.Update))
		admin.DELETE("/categories/:id", ginx.Wrap(h.Delete))
		admin.GET("/categories", ginx.Wrap(h.List))
	}
}

func (h *AdminCategoryHandler) Add(ctx *gin.Context, req service.AddCategoryReq) (ginx.Response, error) {
	_, err := h.svc.Add(ctx, req)
	if err != nil {
		return ginx.ErrorMess("新增分类失败", err.Error()), err
	}
	return ginx.SuccessMess("新增分类成功", nil), nil
}

func (h *AdminCategoryHandler) Update(ctx *gin.Context, req service.UpdateCategoryReq) (ginx.Response, error) {
	req.Id = ctx.Param("id")
	err := h.svc.Update(ctx, req)
	if err != nil {
		return ginx.ErrorMess("更新分类失败", err.Error()), err
	}
	return ginx.SuccessMess("更新分类成功", nil), nil
}

func (h *AdminCategoryHandler) Delete(ctx *gin.Context) (ginx.Response, error) {
	err := h.svc.Delete(ctx, ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("删除分类失败", err.Error()), err
	}
	return ginx.SuccessMess("删除分类成功", nil), nil
}

func (h *AdminCategoryHandler) List(ctx *gin.Context) (ginx.Response, error) {
	res, err := h.svc.List(ctx.Request.Context())
	if err != nil {
		return ginx.ErrorMess("查询分类失败", err.Error()), err
	}
	return ginx.SuccessMess("查询分类成功", res), nil
}
