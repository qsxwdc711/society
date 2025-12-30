package web

import (
	"errors"
	"net/http"
	"society/internal/domain"
	"society/internal/service"
	"society/internal/web/payload"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	svc service.UserServiceInterface
}

func NewUserHandler(svc service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("register", ginx.WrapBody[payload.UserRegisterReq](u.Register))
	ug.POST("login", ginx.WrapBody[payload.UserLoginReq](u.Login))
	ug.GET("/profile", u.profile)
	ug.POST("update", ginx.WrapBody[payload.UserUpdateReq](u.Update))
	ug.POST("changePassword", ginx.WrapBodyAndToken[payload.UserChangeReq, *domain.UserClaims](u.ChangePassword))
}
func (u *UserHandler) Register(ctx *gin.Context, req payload.UserRegisterReq) (ginx.Response, error) {
	role, err := primitive.ObjectIDFromHex(req.Role)
	if err != nil {
		return ginx.ErrorMess("系统错误", nil), err
	}
	if req.Avatar == "" {
		req.Avatar = "http://124.220.212.210:9000/shared/33201979-6723-5bc9-af9f-c4cf44ea56bc"
	}
	data, err := u.svc.Register(ctx, domain.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Age:      req.Age,
		Password: req.Password,
		Role:     role,
		Avatar:   req.Avatar,
	})
	if errors.Is(err, service.ErrRegistered) {
		return ginx.ErrorMess("注册失败", "用户已经注册了"), err
	}
	if err != nil {
		return ginx.ErrorMess("注册失败", "系统错误"), err
	}
	return ginx.SuccessMess("注册成功", data), nil
}

func (u *UserHandler) Login(ctx *gin.Context, req payload.UserLoginReq) (ginx.Response, error) {
	data, err := u.svc.Login(ctx, req.Phone, req.Password)
	if errors.Is(err, service.ErrNotRegistered) {
		return ginx.ErrorMess("登录失败", "账户未注册"), err
	}
	if errors.Is(err, service.ErrFalsePassword) {
		return ginx.ErrorMess("登录失败", "密码错误"), err
	}
	if err != nil {
		return ginx.ErrorMess("登录失败", "系统错误"), err
	}
	return ginx.SuccessMess("登录成功", data), nil
}

func (u *UserHandler) profile(ctx *gin.Context) {

	claimsVal, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未获取到用户信息",
		})
		return
	}
	claims, ok := claimsVal.(*domain.UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "用户信息格式错误",
		})
		return
	}
	profile, err := u.svc.GetProfile(ctx.Request.Context(), claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户信息失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}

func (u *UserHandler) Update(ctx *gin.Context, req payload.UserUpdateReq) (ginx.Response, error) {
	role, err := primitive.ObjectIDFromHex(req.Role)
	if err != nil {
		return ginx.ErrorMess("角色id不正确", nil), err
	}
	data, err := u.svc.Update(ctx, domain.User{
		Id:     req.Id,
		Name:   req.Name,
		Phone:  req.Phone,
		Gender: req.Gender,
		Age:    req.Age,
		Avatar: req.Avatar,
		Role:   role,
	})
	if err != nil {
		return ginx.ErrorMess("修改用户信息失败", nil), err
	}
	return ginx.SuccessMess("修改用户信息成功", data), nil

}

func (u *UserHandler) ChangePassword(ctx *gin.Context, req payload.UserChangeReq, user *domain.UserClaims) (ginx.Response, error) {
	err := u.svc.ChangePass(ctx, user.Uid, req.OldPass, req.NewPass)
	if err != nil {
		return ginx.ErrorMess("修改密码失败", nil), err
	}
	return ginx.SuccessMess("修改密码成功", nil), nil
}
