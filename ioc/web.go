package ioc

import (
	"net/http"
	"society/internal/domain"
	"society/internal/web"
	"society/internal/web/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitGin(
	mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	roleHdl *web.RoleHandler,
	logHdl *web.LogHandler,
	apiHdl *web.ApiHandler,
	adminProductHdl *web.AdminProductHandler,
	adminOrderHdl *web.AdminOrderHandler,
	adminCategoryHdl *web.AdminCategoryHandler,
	adminPromotionHdl *web.AdminPromotionHandler,
	UserProductHdl *web.UserProductHandler,
	AdminLoginLogHdl *web.AdminLoginLogHandler,
	AdminStatsHdl *web.AdminStatsHandler,
	userFavoriteHdl *web.UserFavoriteHandler,
	userCartHandler *web.UserCartHandler,
	userOrderHandler *web.UserOrderHandler,
	userWalletHdl *web.UserWalletHandler,
) *gin.Engine {

	server := gin.New()        // ❗不使用 gin.Default()
	server.Use(gin.Recovery()) // ❗必须要有 panic 保护
	server.Use(mdls...)        // ❗加载你传入的中间件（包含 RequestLogger）

	userHdl.RegisterUserRoutes(server)
	roleHdl.RegisterRoleRouters(server)
	logHdl.RegisterRoutes(server)
	apiHdl.RegisterRoutes(server)
	adminProductHdl.RegisterAdminProductRouters(server)
	adminOrderHdl.RegisterAdminOrderRouters(server)
	adminCategoryHdl.RegisterAdminCategoryRouters(server)
	adminPromotionHdl.RegisterAdminPromotionRouters(server)
	UserProductHdl.RegisterUserProductRouters(server)
	AdminLoginLogHdl.RegisterAdminLoginLogRouters(server)
	AdminStatsHdl.RegisterAdminStatsRouters(server)
	userFavoriteHdl.RegisterUserFavoriteRouters(server)
	userCartHandler.RegisterUserCartRouters(server)
	userOrderHandler.RegisterUserOrderRouters(server)
	userWalletHdl.RegisterUserWalletRouters(server)
	return server
}

func InitMiddlewares(db *mongo.Client, log domain.Loggers) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		initCors(),
		//身份验证
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/register").
			IgnorePaths("/users/login").
			IgnorePaths("/api/v1/user/products").
			IgnorePaths("/api/v1/user/products/").Build(),

		//配置api验证
		middleware.NewApiAuth(db, log).
			IgnorePaths("/users/register").
			IgnorePaths("/users/login").
			IgnorePaths("/api/v1/user/products"). // 同理
			IgnorePaths("/api/v1/user/products/").
			Build(),
	}
}

func initCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
