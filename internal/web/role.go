package web

import (
	"net/http"
	"society/internal/domain"
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Api struct {
	Url    string `json:"url" bson:"url"` // 根据业务需求添加字段
	Method string `json:"method" bson:"method"`
}

type AddRoleReq struct {
	Name  string       `json:"name" binding:"required"`
	Apis  []domain.Api `json:"apis"`
	Codes []string     `json:"codes"`
	Desc  string       `json:"desc"` //角色描述

}
type RoleHandler struct {
	svc service.UserRoleInterface
}

func NewRoleHandler(svc service.UserRoleInterface) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (r *RoleHandler) RegisterRoleRouters(server *gin.Engine) {
	role := server.Group("role")
	{
		role.POST("/add", ginx.WrapBody[domain.Role](r.Add))
		role.GET("/get", r.Get)
		role.POST("/update", ginx.WrapBody[domain.Role](r.Update))
	}
}

func (r *RoleHandler) Add(ctx *gin.Context, req domain.Role) (ginx.Response, error) {
	err := r.svc.Add(ctx, req)
	if err != nil {
		return ginx.ErrorMess("创建角色失败", nil), err
	}
	return ginx.SuccessMess("创建角色成功", nil), nil
}

func (r *RoleHandler) Get(c *gin.Context) {
	data, err := r.svc.Get(c)
	if err != nil {
		zap.L().Error("获取角色失败", zap.Error(err))
		c.JSON(http.StatusOK, "获取角色失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  "获取角色成功",
	})
}

func (r *RoleHandler) Update(ctx *gin.Context, req domain.Role) (ginx.Response, error) {
	err := r.svc.Update(ctx, req)
	if err != nil {
		return ginx.ErrorMess("更新角色信息失败", nil), err
	}
	return ginx.SuccessMess("更新角色信息成功", nil), nil
}
