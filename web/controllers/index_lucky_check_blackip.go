package controllers

import (
	"Lucky_Wheel/models"
	"time"
)

func (c *IndexController) checkBlackip(ip string) (bool, *models.LwBlackip) {
	info := c.ServiceBlackIp.GetByIp(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		return false, info
	}
	return true, info
}
