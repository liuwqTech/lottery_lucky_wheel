package routes

import (
	"Lucky_Wheel/bootstrap"
	"Lucky_Wheel/service"
	"Lucky_Wheel/web/controllers"
	"Lucky_Wheel/web/middleware"
	"github.com/kataras/iris/v12/mvc"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := service.NewUserService()
	giftService := service.NewGiftService()
	codeService := service.NewCodeService()
	resultService := service.NewResultService()
	userDayService := service.NewUserDayService()
	blackIpService := service.NewBlackIpService()

	index := mvc.New(b.Party("/"))
	index.Register(
		userService,
		giftService,
		codeService,
		resultService,
		userDayService,
		blackIpService,
	)
	index.Handle(new(controllers.IndexController))

	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(
		userService,
		giftService,
		codeService,
		resultService,
		blackIpService,
		userDayService,
	)
	admin.Handle(new(controllers.AdminController))

	adminGift := admin.Party("/gift")
	adminGift.Register(giftService)
	adminGift.Handle(new(controllers.AdminGiftController))

	adminCode := admin.Party("/code")
	adminCode.Register(codeService)
	adminCode.Handle(new(controllers.AdminCodeController))

	adminResult := admin.Party("/result")
	adminResult.Register(resultService)
	adminResult.Handle(new(controllers.AdminResultController))

	adminUser := admin.Party("/user")
	adminUser.Register(userService)
	adminUser.Handle(new(controllers.AdminUserController))

	adminBlackip := admin.Party("/blackip")
	adminBlackip.Register(blackIpService)
	adminBlackip.Handle(new(controllers.AdminBlackipController))
}
