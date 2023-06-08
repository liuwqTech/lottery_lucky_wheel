package controllers

import (
	"Lucky_Wheel/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminController struct {
	Ctx            iris.Context
	ServiceGift    service.GiftService
	ServiceCode    service.CodeService
	ServiceUser    service.UserService
	ServiceResult  service.ResultService
	ServiceUserday service.UserDayService
	ServiceBlackip service.BlackIpService
}

func (c *AdminController) Get() mvc.Result {
	return mvc.View{
		Name: "admin/index.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "",
		},
		Layout: "admin/layout.html",
	}
}
