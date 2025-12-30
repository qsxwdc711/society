package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type UserProductHandler struct {
	svc service.UserProductInterface
}

func NewUserProductHandler(svc service.UserProductInterface) *UserProductHandler {
	return &UserProductHandler{svc: svc}
}

func (h *UserProductHandler) RegisterUserProductRouters(server *gin.Engine) {
	user := server.Group("/api/v1/user")
	{
		// 列表/搜索（不需要登录）
		user.GET("/products", ginx.Wrap(h.List))

		// 详情（允许未登录访问；如果登录了会自动记录 uid）
		user.GET("/products/:id", ginx.Wrap(h.Detail))
	}
}

func (h *UserProductHandler) List(ctx *gin.Context) (ginx.Response, error) {
	var req service.UserProductListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	res, err := h.svc.List(ctx.Request.Context(), req)
	if err != nil {
		return ginx.ErrorMess("查询商品列表失败", err.Error()), err
	}
	return ginx.SuccessMess("查询商品列表成功", res), nil
}

func (h *UserProductHandler) Detail(ctx *gin.Context) (ginx.Response, error) {
	p, err := h.svc.Detail(ctx, ctx.Param("id"))
	if err != nil {
		return ginx.ErrorMess("查询商品详情失败", err.Error()), err
	}
	return ginx.SuccessMess("查询商品详情成功", p), nil
}
