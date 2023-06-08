package controllers

import (
	comm "Lucky_Wheel/common"
	"Lucky_Wheel/models"
)

func (c *IndexController) prizeLarge(
	ip string,
	loginuser *models.ObjLoginUser,
	userInfo *models.LwUser,
	blackipInfo *models.LwBlackip) {

	nowTime := comm.NowUnix()
	balckTime := 30 * 86400
	// 更新用户黑名单信息
	if userInfo == nil || userInfo.Id <= 0 {
		userInfo = &models.LwUser{
			Id:         loginuser.Uid,
			Username:   loginuser.Username,
			Blacktime:  nowTime + balckTime,
			SysCreated: nowTime,
			SysIp:      ip,
		}
		err := c.ServiceUser.Create(userInfo)
		if err != nil {
			return
		}
	} else {
		userInfo = &models.LwUser{
			Id:         loginuser.Uid,
			Blacktime:  nowTime + balckTime,
			SysUpdated: nowTime,
		}
		err := c.ServiceUser.Update(userInfo, nil)
		if err != nil {
			return
		}
	}
	// 更新IP黑名单
	if blackipInfo == nil || blackipInfo.Id <= 0 {
		blackipInfo = &models.LwBlackip{
			Ip:         ip,
			Blacktime:  nowTime + balckTime,
			SysCreated: nowTime,
		}
		err := c.ServiceBlackIp.Create(blackipInfo)
		if err != nil {
			return
		}
	} else {
		blackipInfo = &models.LwBlackip{
			Blacktime:  nowTime + balckTime,
			SysUpdated: nowTime,
		}
		err := c.ServiceBlackIp.Update(blackipInfo, nil)
		if err != nil {
			return
		}
	}
}
