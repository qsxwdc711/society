package web

import (
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"

	"github.com/gin-gonic/gin"
)

type UserFavoriteHandler struct {
	svc service.FavoriteInterface
}

func NewUserFavoriteHandler(svc service.FavoriteInterface) *UserFavoriteHandler {
	return &UserFavoriteHandler{svc: svc}
}

func (h *UserFavoriteHandler) RegisterUserFavoriteRouters(server *gin.Engine) {
	user := server.Group("/api/v1/user")
	{
		user.POST("/favorites/:productId", ginx.WrapToken[*domain.UserClaims](h.Add))
		user.DELETE("/favorites/:productId", ginx.WrapToken[*domain.UserClaims](h.Remove))
		user.GET("/favorites", ginx.WrapToken[*domain.UserClaims](h.List))
		user.GET("/favorites/check/:productId", ginx.WrapToken[*domain.UserClaims](h.Check))
	}
}

func (h *UserFavoriteHandler) Add(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	err := h.svc.Add(ctx.Request.Context(), uc.Uid, ctx.Param("productId"))
	if err != nil {
		return ginx.ErrorMess("收藏失败", err.Error()), err
	}
	return ginx.SuccessMess("收藏成功", nil), nil
}

func (h *UserFavoriteHandler) Remove(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	err := h.svc.Remove(ctx.Request.Context(), uc.Uid, ctx.Param("productId"))
	if err != nil {
		return ginx.ErrorMess("取消收藏失败", err.Error()), err
	}
	return ginx.SuccessMess("取消收藏成功", nil), nil
}

func (h *UserFavoriteHandler) List(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	var page ginx.Page
	if err := ctx.ShouldBindQuery(&page); err != nil {
		return ginx.ErrorMess("参数格式错误", err.Error()), err
	}
	data, err := h.svc.List(ctx.Request.Context(), uc.Uid, page)
	if err != nil {
		return ginx.ErrorMess("查询收藏列表失败", err.Error()), err
	}
	return ginx.SuccessMess("查询收藏列表成功", data), nil
}

func (h *UserFavoriteHandler) Check(ctx *gin.Context, uc *domain.UserClaims) (ginx.Response, error) {
	data, err := h.svc.Check(ctx.Request.Context(), uc.Uid, ctx.Param("productId"))
	if err != nil {
		return ginx.ErrorMess("查询收藏状态失败", err.Error()), err
	}
	return ginx.SuccessMess("查询收藏状态成功", data), nil
}
