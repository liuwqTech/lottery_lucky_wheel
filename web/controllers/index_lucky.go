package controllers

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/conf"
	"Lucky_Wheel/models"
	"Lucky_Wheel/web/utils"
	"fmt"
	"log"
)

func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1 验证登录用户
	loginUser := comm.GetLoginUser(c.Ctx.Request())
	if loginUser == nil || loginUser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录再抽奖"
		return rs
	}
	// 2 用户抽奖分布式锁定
	ok := utils.LockLucky(loginUser.Uid)
	if ok {
		defer utils.UnLockLucky(loginUser.Uid)
	} else {
		rs["code"] = 102
		rs["msg"] = "正在抽奖，请稍后重试!"
		return rs
	}
	log.Println("uid:", loginUser.Uid)
	// 3 验证用户参与次数
	ok = c.checkUserDay(loginUser.Uid)
	if !ok {
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完!"
		return rs
	}
	// 4 验证IP今日参与次数
	ip := comm.ClientIp(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "相同IP参与次数太多明天再参与!"
		return rs
	}
	// 5 验证IP黑名单
	limitBlack := false // 黑名单
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	var blackipInfo *models.LwBlackip
	if !limitBlack {
		ok, blackipInfo = c.checkBlackip(ip)
		if !ok {
			fmt.Println("黑名单中的IP", ip, limitBlack)
			limitBlack = true
		}
	}
	var userInfo *models.LwUser
	// 6 验证用户黑名单
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginUser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginUser.Uid, limitBlack)
			limitBlack = true
		}
	}
	// 7 获得抽奖编码
	prizeCode := comm.Random(10000)
	// 8 匹配奖品是否中奖
	prizeGift := c.prize(prizeCode, limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum < 0 ||
		prizeGift.PrizeNum > 0 && prizeGift.LeftNum < 0 {
		rs["code"] = 205
		rs["msg"] = "很遗憾没有中奖!"
		return rs
	}
	// 9 有限制奖品发放
	if prizeGift.PrizeNum > 0 {
		ok := utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum)
		if !ok {
			rs["code"] = 207
			rs["msg"] = "很遗憾没有中奖，请下次再试!"
			return rs
		}
	}
	// 10 不同编码优惠券发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode)
		if code == "" {
			rs["code"] = 208
			rs["msg"] = "很遗憾没有中奖，请下次再试!"
			return rs
		}
		prizeGift.Gdata = code
	}
	// 11 中奖纪录
	result := models.LwResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginUser.Uid,
		Username:   loginUser.Username,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(&result)
	if err != nil {
		log.Println("index_luck.GetLucky ServiceResult.Create", result, err)
		rs["code"] = 209
		rs["msg"] = "很遗憾没有中奖，请下次再试!"
		return rs
	}
	// 12 返回抽奖结果
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 实物大奖，需要设置黑名单用户和IP
		c.prizeLarge(ip, loginUser, userInfo, blackipInfo)
	}
	rs["data"] = prizeGift
	rs["code"] = 200
	rs["msg"] = "恭喜中奖!"
	return rs
}
