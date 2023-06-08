package controllers

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/models"
	"Lucky_Wheel/service"
	"fmt"

	"github.com/kataras/iris/v12"
)

type IndexController struct {
	Ctx            iris.Context
	ServiceUser    service.UserService
	ServiceBlackIp service.BlackIpService
	ServiceCode    service.CodeService
	ServiceGift    service.GiftService
	ServiceResult  service.ResultService
	ServiceUserDay service.UserDayService
}

// 首页
// http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to Go 抽奖系统<a href='/public/index.html'>开始抽奖</a>"
}

// 获取礼品信息
// http://localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	dataList := c.ServiceGift.GetAll(true)
	list := make([]models.LwGift, 0)
	for _, data := range dataList {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}
	rs["gifts"] = list
	return rs
}

// http://localhost:8080/newPrize
func (c *IndexController) GetNewPrize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	list := make([]models.LwGift, 0)
	rs["gifts"] = list
	return rs
}

func (c *IndexController) GetLogin() {
	uid := comm.Random(100000)
	loginuser := models.ObjLoginUser{
		Uid:      uid,
		Username: fmt.Sprintf("admin-%d", uid),
		Now:      comm.NowUnix(),
		Ip:       comm.ClientIp(c.Ctx.Request()),
	}
	comm.SetLoginUser(c.Ctx.ResponseWriter(), &loginuser)
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/index.html?from=login")
}

func (c *IndexController) GetLogout() {
	comm.SetLoginUser(c.Ctx.ResponseWriter(), nil)
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/index.html?from=logout")
}
