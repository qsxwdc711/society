package web

import (
	"society/internal/service"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	svc service.LogServiceInterface
}

func NewLogHandler(svc service.LogServiceInterface) *LogHandler {
	return &LogHandler{
		svc: svc,
	}
}

func (l *LogHandler) RegisterRoutes(server *gin.Engine) {
	log := server.Group("/log")
	{
		log.GET("/get", ginx.WrapBody[ginx.Page](l.Get))
	}
}

func (l *LogHandler) Get(ctx *gin.Context, req ginx.Page) (ginx.Response, error) {
	level := ctx.DefaultQuery("level", "error")
	data, count, err := l.svc.Get(ctx, req, level)
	if err != nil {
		return ginx.ErrorMess("获取日志记录失败", nil), err
	}
	return ginx.SuccessMess("获取日志记录成功", gin.H{
		"logs":  data,
		"count": count,
	}), nil
}
