package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type AdminLoginLogHandler struct {
	svc service.LoginLogInterface
}

func NewAdminLoginLogHandler(svc service.LoginLogInterface) *AdminLoginLogHandler {
	return &AdminLoginLogHandler{svc: svc}
}

func (h *AdminLoginLogHandler) RegisterAdminLoginLogRouters(server *gin.Engine) {
	admin := server.Group("/api/v1/admin")
	{
		admin.GET("/logs/admin-login", ginx.Wrap(h.AdminLoginLogs))
		admin.GET("/logs/user-login", ginx.Wrap(h.UserLoginLogs))
	}
}

func (h *AdminLoginLogHandler) AdminLoginLogs(ctx *gin.Context) (ginx.Response, error) {
	var page ginx.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	res, err := h.svc.ListAdmin(ctx.Request.Context(), page)
	if err != nil {
		return ginx.ErrorMess("查询管理员登录日志失败", err.Error()), err
	}
	return ginx.SuccessMess("查询管理员登录日志成功", res), nil
}

func (h *AdminLoginLogHandler) UserLoginLogs(ctx *gin.Context) (ginx.Response, error) {
	var page ginx.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	res, err := h.svc.ListUser(ctx.Request.Context(), page)
	if err != nil {
		return ginx.ErrorMess("查询用户登录日志失败", err.Error()), err
	}
	return ginx.SuccessMess("查询用户登录日志成功", res), nil
}
