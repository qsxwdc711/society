//go:build wireinject

package main

import (
	"society/internal/repository"
	"society/internal/repository/dao"
	"society/internal/service"
	"society/internal/web"
	"society/ioc"

	"github.com/google/wire"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitMongodb,
		ioc.InitRedis,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.InitLogger,
		dao.NewUserDao,
		dao.NewRoleDao,
		dao.NewLogDao,
		dao.NewApiDao,
		dao.NewProductDao,
		dao.NewCategoryDao,
		dao.NewPromotionDao,
		dao.NewAdminOrderDao,
		dao.NewStatsDao,
		dao.NewLoginLogDao,
		dao.NewProductViewDao,
		dao.NewUserProductDao,
		repository.NewUserRepo,
		repository.NewRoleRepo,
		repository.NewLogRepository,
		repository.NewApiRepo,
		repository.NewCategoryRepo,
		repository.NewPromotionRepo,
		repository.NewAdminOrderRepo,
		repository.NewProductRepo,
		repository.NewProductViewRepo,
		repository.NewStatsRepo,
		repository.NewLoginLogRepo,
		repository.NewUserProductRepo,
		service.NewUserService,

		service.NewRoleService,
		service.NewLogService,
		service.NewApiService,
		service.NewCategoryService,
		service.NewPromotionService,
		service.NewAdminOrderService,
		service.NewProductAdminService,
		service.NewStatsService,
		service.NewLoginLogService,
		service.NewUserProductService,
		wire.Bind(new(service.UserServiceInterface), new(*service.UserService)),
		web.NewUserHandler,
		web.NewRoleHandler,
		web.NewLogHandler,
		web.NewApiHandler,
		web.NewAdminProductHandler,
		web.NewAdminOrderHandler,
		web.NewAdminCategoryHandler,
		web.NewAdminPromotionHandler,
		web.NewAdminStatsHandler,
		web.NewAdminLoginLogHandler,
		web.NewUserProductHandler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
