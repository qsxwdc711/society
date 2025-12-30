package middleware

import (
	"context"
	"net/http"
	"society/internal/domain"
	"society/internal/repository/dao"
	"strings"

	"gitee.com/qsxwdc711/pkgx/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ApiAuth struct {
	paths   []string
	mongodb *mongo.Database
	log     *zap.Logger
}

func NewApiAuth(mongodb *mongo.Client, log domain.Loggers) *ApiAuth {
	databaseName := viper.GetString("mongo.database")
	return &ApiAuth{
		mongodb: mongodb.Database(databaseName),
		log:     log.Logg,
	}
}

func (l *ApiAuth) IgnorePaths(path string) *ApiAuth {
	l.paths = append(l.paths, path)
	return l
}

func (l *ApiAuth) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 取路由模板（优先 FullPath），用于：忽略匹配 + api 表匹配 + role 权限匹配
		full := c.FullPath()
		if full == "" {
			// 没有匹配到路由模板（可能是 404/OPTIONS/未注册）
			// 建议直接放过，让后续路由/404处理
			c.Next()
			return
		}

		// 2) 忽略路径：支持 精确匹配 + 前缀匹配（以 / 结尾表示前缀）
		for _, path := range l.paths {
			// 精确
			if full == path {
				c.Next()
				return
			}
			// 前缀（用于 /xxx/:id 这种）
			if strings.HasSuffix(path, "/") && strings.HasPrefix(full, path) {
				c.Next()
				return
			}
		}

		// 3) 继续权限校验：用 FullPath 模板与 api 表一致
		url := full
		method := c.Request.Method
		// 查找api
		var api dao.Api
		apiCol := l.mongodb.Collection("api")
		if err := apiCol.FindOne(context.TODO(), bson.M{"url": url, "method": method}).Decode(&api); err != nil {
			c.JSON(http.StatusForbidden, ginx.ErrorMess("验证api：此api不存在", err.Error()))
			c.Abort()
			return
		}

		val, _ := c.Get("claims")
		userClaim, ok := val.(*domain.UserClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 查找用户
		var user dao.User
		userColl := l.mongodb.Collection("user")
		err := userColl.FindOne(context.TODO(), bson.M{"_id": userClaim.Uid}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusForbidden, ginx.ErrorMess("验证api：获取用户失败", err.Error()))
			c.Abort()
			return
		}

		// 查找角色
		var role dao.Role
		roleColl := l.mongodb.Collection("role")
		err = roleColl.FindOne(context.TODO(), bson.M{"_id": user.Role}).Decode(&role)
		if err != nil {
			c.JSON(http.StatusForbidden, ginx.ErrorMess("验证api：获取用户角色失败", err.Error()))
			c.Abort()
			return
		}

		// 权限匹配：同样用 FullPath 模板匹配 role.Apis
		for _, api = range role.Apis {
			if api.Url == url && api.Method == method {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, ginx.ErrorMess("验证api：此用户无访问此api的权限", nil))
		c.Abort()
	}
}
