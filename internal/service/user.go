package service

import (
	"context"
	"errors"
	"society/internal/domain"
	"society/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRegistered    = errors.New("用户已经注册了")
	ErrNotRegistered = errors.New("用户没有注册")
	ErrFalsePassword = errors.New("密码错误")
)

type UserServiceInterface interface {
	Register(ctx context.Context, user domain.User) (any, error)
	Login(ctx context.Context, phone string, password string) (any, error)
	GetProfile(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	ChangePass(ctx *gin.Context, uid primitive.ObjectID, pass string, pass2 string) error
}

type UserService struct {
	repo     repository.UserRepoInterface
	roleRepo repository.RoleRepoInterface
	logSvc   LoginLogInterface
}

func NewUserService(
	repo repository.UserRepoInterface,
	roleRepo repository.RoleRepoInterface,
	logSvc LoginLogInterface,
) *UserService {
	return &UserService{
		repo:     repo,
		roleRepo: roleRepo,
		logSvc:   logSvc,
	}
}

func (svc *UserService) Register(ctx context.Context, user domain.User) (any, error) {
	_, err := svc.repo.FindOneByPhone(ctx, user.Phone)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if err == nil {
		return nil, ErrRegistered
	}
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	// 入库
	res, err := svc.repo.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	//重新构造返回给客户端的完整响应对象
	resDTO := domain.User{ // 假设 domain.User 包含了所有非敏感字段
		Id:     res.Id,      // DAO 返回的 ID
		Name:   user.Name,   // 注册请求中的 Name
		Phone:  user.Phone,  // 注册请求中的 Phone
		Gender: user.Gender, // 注册请求中的 Sex
		Age:    user.Age,
		Avatar: user.Avatar, // Avatar
		Role:   user.Role,   // Role
		// 注意：不包含 Password
	}
	//生成token
	tokenString, err := createToken(res.Id)
	if err != nil {
		return nil, err
	}
	resDTO.Token = tokenString
	return resDTO, nil

}

func (svc *UserService) Login(ctx context.Context, phone string, password string) (any, error) {
	user, err := svc.repo.FindOneByPhone(ctx, phone)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotRegistered
	}
	if err != nil {
		return nil, err
	}

	role, err := svc.roleRepo.FindRoleById(ctx, user.Role)
	if err != nil {
		return nil, err
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrFalsePassword
	}

	// 生成 token
	tokenString, err := createToken(user.Id)
	if err != nil {
		return nil, err
	}

	// ===== 新增：记录登录日志（不影响主流程） =====
	if ginCtx, ok := ctx.(*gin.Context); ok {
		_ = svc.logSvc.Add(ginCtx, user.Id, domain.LoginRoleUser)
	}
	// ============================================

	user.Password = ""
	return gin.H{
		"token": "Bearer " + tokenString,
		"user":  user,
		"role":  role,
	}, nil
}

func (svc *UserService) GetProfile(ctx context.Context, id primitive.ObjectID) (domain.User, error) {

	u, err := svc.repo.FindOneById(ctx, id)
	if err != nil {
		zap.L().Error("service.GetProfile repo.FindById error", zap.Error(err))
		return domain.User{}, err
	}
	u.Password = ""
	return u, nil
}

func (svc *UserService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	return svc.repo.Update(ctx, user)
}

func (svc *UserService) ChangePass(ctx *gin.Context, uid primitive.ObjectID, pass string, pass2 string) error {
	user, err := svc.repo.FindOneById(ctx, uid)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return err
	}
	//密码相同可以改密码
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass2), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return svc.repo.ChangePassword(ctx, uid, string(hashPass))
}

func createToken(Id primitive.ObjectID) (tokenString string, err error) {
	// 创建一个我们自己的声明
	secret := viper.GetString("general.jwt")
	claims := domain.UserClaims{
		//设置参数
		RegisteredClaims: jwt.RegisteredClaims{
			//设置7天的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
