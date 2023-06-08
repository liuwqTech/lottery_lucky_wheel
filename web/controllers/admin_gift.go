package controllers

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/models"
	"Lucky_Wheel/service"
	"Lucky_Wheel/web/utils"
	"Lucky_Wheel/web/viewmodels"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminGiftController struct {
	Ctx            iris.Context
	ServiceGift    service.GiftService
	ServiceCode    service.CodeService
	ServiceUser    service.UserService
	ServiceResult  service.ResultService
	ServiceUserday service.UserDayService
	ServiceBlackip service.BlackIpService
}

func (c *AdminGiftController) Get() mvc.Result {
	datalist := c.ServiceGift.GetAll(true)
	total := c.ServiceGift.CountAll()
	for i, data := range datalist {
		// 发放计划
		prizedata := make([][2]int, 0)
		err := json.Unmarshal([]byte(data.PrizeData), &prizedata)
		if err != nil || len(prizedata) < 1 {
			datalist[i].PrizeData = "[]"
		} else {
			newpd := make([]string, len(prizedata))
			for index, pd := range prizedata {
				ct := comm.FormatFromUnixTime(int64(pd[0]))
				newpd[index] = fmt.Sprintf("【%s】: %d", ct, pd[1])
			}
			str, err := json.Marshal(newpd)
			if err == nil && len(str) > 0 {
				datalist[i].PrizeData = string(str)
			} else {
				datalist[i].PrizeData = "[]"
			}
		}
		//num := utils.GetGiftPoolNum(giftInfo.Id)
		//datalist[i].Title = fmt.Sprintf("【%d】: %s", num, datalist[i].Title)
	}
	return mvc.View{
		Name: "admin/gift.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "gift",
			"Datalist": datalist,
			"Total":    total,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) GetEdit() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServiceGift.Get(id, true)
		giftInfo.Id = data.Id
		giftInfo.Title = data.Title
		giftInfo.PrizeNum = data.PrizeNum
		giftInfo.PrizeCode = data.PrizeCode
		giftInfo.PrizeTime = data.PrizeTime
		giftInfo.Img = data.Img
		giftInfo.Displayorder = data.Displayorder
		giftInfo.Gtype = data.Gtype
		giftInfo.Gdata = data.Gdata
		giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
		giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
	}
	return mvc.View{
		Name: "admin/giftEdit.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) PostSave() mvc.Result {
	data := viewmodels.ViewGift{}
	err := c.Ctx.ReadForm(&data)
	if err != nil {
		fmt.Println("admin_gift.PostSave.ReadForm", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadForm转换异常:%s", err),
		}
	}
	giftInfo := models.LwGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.Img = data.Img
	giftInfo.Displayorder = data.Displayorder
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	t1, err1 := comm.ParseTime(data.TimeBegin)
	giftInfo.TimeBegin = int(t1.Unix())
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		fmt.Println("admin_gift.PostSave.ReadForm", err)
		return mvc.Response{
			Text: fmt.Sprintf("开始时间和结束时间的格式不正确 err1=%s, err2=%s", err1, err2),
		}
	}
	giftInfo.TimeEnd = int(t2.Unix())
	if giftInfo.Id > 0 {
		// 数据更新
		dataInfo := c.ServiceGift.Get(giftInfo.Id, true)
		if dataInfo != nil && dataInfo.Id > 0 {
			if dataInfo.PrizeNum != giftInfo.PrizeNum {
				// 数量发生变化
				giftInfo.LeftNum = dataInfo.LeftNum - dataInfo.PrizeNum - giftInfo.PrizeNum
				if giftInfo.LeftNum <= 0 || giftInfo.PrizeNum <= 0 {
					giftInfo.LeftNum = 0
				}
				// 奖品数量变化
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			if dataInfo.PrizeTime != giftInfo.PrizeTime {
				// 发奖周期发生变化
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			giftInfo.SysUpdated = int(time.Now().Unix())
			err := c.ServiceGift.Update(&giftInfo, []string{})
			if err != nil {
				return nil
			}
		} else {
			giftInfo.Id = 0
		}
	}
	if giftInfo.Id == 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIp(c.Ctx.Request())
		giftInfo.SysCreated = int(time.Now().Unix())
		err := c.ServiceGift.Create(&giftInfo)
		if err != nil {
			return nil
		}
		// 新的奖品，更新奖品的发奖计划
		utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceGift.Delete(id)
		if err != nil {
			return nil
		}
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceGift.Update(&models.LwGift{
			Id:        id,
			SysStatus: 0,
		}, []string{"sys_status"})
		if err != nil {
			return nil
		}
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}
