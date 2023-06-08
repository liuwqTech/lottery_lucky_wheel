package controllers

import (
	"Lucky_Wheel/models"
	"time"
)

func (c *IndexController) checkBlackUser(uid int) (bool, *models.LwUser) {
	info := c.ServiceUser.Get(uid)
	if info != nil && info.Blacktime > int(time.Now().Unix()) {
		// 黑名单存在并且有效
		return false, info
	}
	return true, info
}
