package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	svc service.ApiServiceInterface
}

func NewApiHandler(svc service.ApiServiceInterface) *ApiHandler {
	return &ApiHandler{
		svc: svc,
	}
}

func (a *ApiHandler) RegisterRoutes(server *gin.Engine) {
	api := server.Group("api")
	{
		api.GET("get", ginx.Wrap(a.GetAllApis))
		api.GET("page", ginx.WrapBody(a.PageGet))
	}
}

func (a *ApiHandler) GetAllApis(ctx *gin.Context) (ginx.Response, error) {
	data, count, err := a.svc.GetAllApis(ctx)
	if err != nil {
		return ginx.ErrorMess("获取全部api失败", nil), err
	}
	return ginx.SuccessMess("获取全部api成功", gin.H{
		"apis":  data,
		"count": count,
	}), nil
}

func (a *ApiHandler) PageGet(ctx *gin.Context, page ginx.Page) (ginx.Response, error) {
	url := ctx.Query("url")
	method := ctx.Query("method")
	data, count, err := a.svc.PageGet(ctx, page, url, method)
	if err != nil {
		return ginx.ErrorMess("获取api分页失败", nil), err
	}
	return ginx.SuccessMess("获取api分页成功", gin.H{
		"apis":  data,
		"count": count,
	}), nil
}
