package middleware

import (
	"context"
	"net/http"
	"society/internal/domain"
	"society/internal/repository/dao"

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
		// 1) 忽略路径：命中后要 Next()
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 2) 用 FullPath 取“路由模板”，与 engine.Routes() 写入 api 表一致
		url := c.FullPath()
		if url == "" {
			// 如果没匹配到路由（比如 404），直接放过或拦截都行
			c.JSON(http.StatusForbidden, ginx.ErrorMess("验证api：路由未注册", nil))
			c.Abort()
			return
		}
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
