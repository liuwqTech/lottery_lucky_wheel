package controllers

import (
	"Lucky_Wheel/models"
	"Lucky_Wheel/service"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminResultController struct {
	Ctx           iris.Context
	ServiceUser   service.UserService
	ServiceResult service.ResultService
}

func (c *AdminResultController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	uid := c.Ctx.URLParamIntDefault("uid", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	var datalist []models.LwResult
	if giftId > 0 {
		datalist = c.ServiceResult.SearchByGift(giftId, page, size)
	} else if uid > 0 {
		datalist = c.ServiceResult.SearchByUser(uid, page, size)
	} else {
		datalist = c.ServiceResult.GetAll(page, size)
	}
	total := (page-1)*size + len(datalist)
	// 数据总数
	if len(datalist) >= size {
		if giftId > 0 {
			total = int(c.ServiceResult.CountByGift(giftId))
		} else if uid > 0 {
			total = int(c.ServiceResult.CountByUser(uid))
		} else {
			total = int(c.ServiceResult.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/result.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "result",
			"GiftId":   giftId,
			"Uid":      uid,
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminResultController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceResult.Delete(id)
		if err != nil {
			return nil
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminResultController) GetCheat() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceResult.Update(&models.LwResult{
			Id:        id,
			SysStatus: 2,
		}, []string{"sys_status"})
		if err != nil {
			return nil
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminResultController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceResult.Update(&models.LwResult{
			Id:        id,
			SysStatus: 0,
		}, []string{"sys_status"})
		if err != nil {
			return nil
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}
