package controllers

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/conf"
	"Lucky_Wheel/models"
	"Lucky_Wheel/service"
	"Lucky_Wheel/web/utils"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminCodeController struct {
	Ctx         iris.Context
	ServiceGift service.GiftService
	ServiceCode service.CodeService
	ServiceUser service.UserService
}

func (c *AdminCodeController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceCode.Delete(id)
		if err != nil {
			return nil
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminCodeController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		err := c.ServiceCode.Update(&models.LwCode{
			Id:        id,
			SysStatus: 0,
		}, []string{"sys_status"})
		if err != nil {
			return nil
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminCodeController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	var datalist []models.LwCode
	var total int
	var num int
	var cacheNum int
	if giftId > 0 {
		datalist = c.ServiceCode.Search(giftId)
		num, cacheNum = utils.GetCacheCodeNum(giftId, c.ServiceCode)
	} else {
		datalist = c.ServiceCode.GetAll(page, size)
	}
	total = (page - 1) + len(datalist)
	if len(datalist) >= size {
		if giftId > 0 {
			total = int(c.ServiceCode.CountByGift(giftId))
		} else {
			total = int(c.ServiceCode.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	} else {
		pagePrev = fmt.Sprintf("%d", 1)
	}
	return mvc.View{
		Name: "admin/code.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "code",
			"GiftId":   giftId,
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
			"CodeNUm":  num,
			"CacheNum": cacheNum,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminCodeController) PostImport() {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	if giftId < 1 {
		_, err := c.Ctx.Text("没有置顶奖品ID，无法导入 <a href='' onclick='history.go(-1)'></a>")
		if err != nil {
			return
		}
		return
	}
	gift := c.ServiceGift.Get(giftId, true)
	if gift == nil || gift.Gtype != conf.GtypeCodeDiff {
		_, err := c.Ctx.Text("奖品信息不存在或不是差异化奖品类型！")
		if err != nil {
			return
		}
		return
	}
	codes := c.Ctx.PostValue("codes")
	now := comm.NowUnix()
	list := strings.Split(codes, "\n")
	sucNum := 0
	errNum := 0
	for _, code := range list {
		code := strings.TrimSpace(code)
		if code != "" {
			data := &models.LwCode{
				GiftId:     giftId,
				Code:       code,
				SysCreated: now,
			}
			err := c.ServiceCode.Create(data)
			if err != nil {
				errNum++
			} else {
				sucNum++
				// 成功导入数据库，下一步还需要导入到缓存
				ok := utils.ImportCacheCodes(giftId, code)
				if ok {
					sucNum++
				} else {
					errNum++
				}
			}
		}
	}
	_, err := c.Ctx.HTML(fmt.Sprintf("成功导入%d条数据，失败导入%d数据", sucNum, errNum))
	if err != nil {
		return
	}
}

func (c *AdminCodeController) GetRecache() {
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	id, err := c.Ctx.URLParamInt("id")
	if id < 1 || err != nil {
		rs := fmt.Sprintf("没有指定优惠券属于的奖品id，<a href='%s'>呼救</a>", refer)
		_, err := c.Ctx.HTML(rs)
		if err != nil {
			return
		}
	}
	sucNum, errNum := utils.RecacheCodes(id, c.ServiceCode)
	rs := fmt.Sprintf("sucNum=%d, errNum=%d, <a href='%s'></a>", sucNum, errNum, refer)
	_, err = c.Ctx.HTML(rs)
	if err != nil {
		return
	}
}
